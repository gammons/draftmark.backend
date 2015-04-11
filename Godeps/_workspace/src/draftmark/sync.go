package draftmark

import (
	dropbox "draftmark/dropbox_client"
	db "draftmark/persistence"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"strings"
)

type Sync struct {
	Db        db.DBClient
	Dropbox   dropbox.DropboxClient
	DoLogging bool
}

func init() {
	setupDotEnv()
}

func NewSync() *Sync {
	dbox := &dropbox.Client{}
	dbox.InitDropbox()
	dbase := db.NewClient()
	return &Sync{Db: dbase, Dropbox: dbox}
}

func (s *Sync) Log(line string) {
	if s.DoLogging {
		log.Println(line)
	}
}

func (s *Sync) DoSync(user db.User, prefix string) {
	s.Log("Getting changes")
	s.Dropbox.SetAccessToken(user.DropboxAccessToken)
	nextCursor, entries := s.Dropbox.GetChanges(&user.DropboxCursor, prefix)

	s.Log("Updating cursor")
	s.Db.UpdateUserCursor(&user, nextCursor)

	s.Log("********* Beginning processing loop")
	for _, entry := range entries {
		s.Log("Processing: " + entry.Path)
		if !strings.HasPrefix(entry.Path, prefix) || !strings.HasSuffix(entry.Path, ".md") {
			continue
		}

		if entry.IsDeleted {
			s.Log("Deleted:" + entry.Path)
			s.deleteEntry(&user, entry)
		} else {
			s.Log("Updated or created:" + entry.Path)
			s.createOrUpdateEntry(&user, entry)
		}
	}
	s.Log("********* Ending  processing loop")
}

func (s *Sync) deleteEntry(user *db.User, entry *dropbox.DropboxEntry) {
	note := db.Note{UserId: user.ID, Path: entry.Path, Mtime: entry.Modified}
	s.Db.DeleteNote(&note)
}

func (s *Sync) createOrUpdateEntry(user *db.User, entry *dropbox.DropboxEntry) {
	content, _ := s.Dropbox.GetContent(entry.Path)
	title := getTitle(content)
	note := db.Note{UserId: user.ID, Path: entry.Path, Mtime: entry.Modified, Content: content, Title: title}
	s.Db.SaveNote(user, &note)
}

func getTitle(content string) string {
	return strings.Split(content, "\n")[0]
}

func setupDotEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file.  Maybe that is ok though.")
	}
}
