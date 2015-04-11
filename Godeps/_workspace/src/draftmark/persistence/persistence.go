package persistence

import (
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"os"
	"time"
)

type User struct {
	ID                 int `gorm:"primary_key"`
	Email              string
	EncryptedPassword  string
	DropboxAccessToken string
	DropboxCursor      string
	DropboxUserId      string
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type Note struct {
	ID             int       `gorm:"primary_key" json:"id"`
	Title          string    `json:"title"`
	Content        string    `sql:"type:text" json:"-"`
	ContentPreview string    `json:"preview"`
	UserId         int       `json:"user_id"`
	Path           string    `json:"path"`
	Mtime          time.Time `json:"mtime"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type DBClient interface {
	DeleteNote(note *Note) bool
	SaveNote(user *User, note *Note) bool
	UpdateUserCursor(user *User, cursor string) bool
	ListNotes(user *User) []Note
	GetNoteContents(user *User, path string) string
	FindOrCreateUser(user *User) bool
}

type Client struct {
	Db gorm.DB
}

func (c *Client) InitDB() {
	c.Db, _ = gorm.Open("postgres", "dbname="+os.Getenv("DB_TABLENAME")+" sslmode=disable")
	c.Db.LogMode(true)
}

func NewClient() DBClient {
	client := &Client{}
	client.InitDB()
	return client
}

func (c *Client) NoteCount() int {
	var count int
	c.Db.Table("notes").Count(&count)
	return count
}

func (c *Client) ListNotes(user *User) []Note {
	var notes []Note
	c.Db.Select("id, title, left(content, 255) as content_preview, user_id, path, mtime, created_at, updated_at").Where("user_id = ?", user.ID).Limit(100).Order("mtime desc").Find(&notes)
	return notes
}

func (c *Client) GetNoteContents(user *User, path string) string {
	var note Note
	c.Db.Select("content").Where("user_id = ? AND path = ?", user.ID, path).First(&note)
	return note.Content
}

func (c *Client) DeleteNote(note *Note) bool {
	c.Db.Delete(note)
	return true
}

func (c *Client) SaveNote(user *User, note *Note) bool {
	var existingNote Note
	c.Db.Where("user_id = ? AND path = ?", user.ID, note.Path).FirstOrInit(&existingNote)
	note.ID = existingNote.ID
	c.Db.Save(note)
	return true
}

func (c *Client) FindOrCreateUser(user *User) bool {
	c.Db.Where("email = ?", user.Email).FirstOrCreate(&user)
	return true
}

func (c *Client) UpdateUserCursor(user *User, cursor string) bool {
	user.DropboxCursor = cursor
	c.Db.Save(user)
	return true
}

// func (c *Client) UpdateUserAccessToken(user *User, accessToken string) bool {
// 	user.DropboxAccessToken = accessToken
// 	c.Db.Save(user)
// 	return true
// }
