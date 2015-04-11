package draftmark_test

import (
	. "draftmark"
	dropbox "draftmark/dropbox_client"
	db "draftmark/persistence"
	"github.com/joho/godotenv"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"log"
	"os"
	"time"
)

var database = &db.Client{}
var dbox = &dropbox.Client{}
var tmpFolder = "/__draftmarktest__"

var _ = BeforeSuite(func() {
	_ = godotenv.Load("test.env")
	setupDatabase()
	setupDropbox()
})

var _ = AfterSuite(func() {
	log.Println(dbox.Dbox.Delete(tmpFolder))
	dbox.Dbox.Delete(tmpFolder)
	database.Db.DropTable(&db.User{})
	database.Db.DropTable(&db.Note{})
})

var _ = Describe("Sync Integration", func() {
	var user = db.User{ID: 1, Email: "gammons@gmail.com", DropboxCursor: "", DropboxAccessToken: os.Getenv("DROPBOX_ACCESS_TOKEN"), DropboxUserId: "123", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	var sync = &Sync{Db: database, Dropbox: dbox}

	BeforeEach(func() {
		database.Db.Delete(&db.Note{})
		database.Db.Delete(&db.User{})
	})

	Context("A new file is added in dropbox", func() {
		BeforeEach(func() {
			cursor, _ := dbox.Dbox.LatestCursor(tmpFolder, false)

			// ensure you delete test.md in your real dropbox folder for this test to pass.
			dbox.Dbox.UploadFile("test.md", tmpFolder+"/test.md", true, "")
			user.DropboxCursor = cursor.Cursor
		})

		It("Adds the file to the db", func() {
			sync.DoSync(user, tmpFolder)
			Expect(database.NoteCount()).To(Equal(1))
			dbox.Dbox.Delete(tmpFolder + "/test.md")
		})

		It("Sets the title", func() {
			sync.DoSync(user, tmpFolder)
			var note db.Note
			database.Db.First(&note)
			Expect(note.Title).To(Equal("# this is a test file"))
		})
	})
	Context("A file was changed", func() {
		It("Removes the file from the db and replaces it with a new entry", func() {
		})
	})

	Context("A file was deleted from dropbox", func() {
		BeforeEach(func() {
			// get the latest cursor
			cursor, _ := dbox.Dbox.LatestCursor(tmpFolder, false)
			user.DropboxCursor = cursor.Cursor

			// upload the file we want to delete, and ensure it gets put into our local db
			dbox.Dbox.UploadFile("test.md", tmpFolder+"/deleted.md", true, "")
			sync.DoSync(user, tmpFolder)

			// set the cursor for the user
			cursor, _ = dbox.Dbox.LatestCursor(tmpFolder, false)
			user.DropboxCursor = cursor.Cursor

			log.Println(dbox.Dbox.Delete(tmpFolder + "/deleted.md"))
		})

		It("Removes the file to the db", func() {
			sync.DoSync(user, tmpFolder)
			Expect(database.NoteCount()).To(Equal(0))
		})
	})
})

func setupDatabase() {
	database.InitDB()
	database.Db.DropTable(&db.User{})
	database.Db.DropTable(&db.Note{})
	database.Db.CreateTable(&db.User{})
	database.Db.CreateTable(&db.Note{})
}

func setupDropbox() {
	dbox.InitDropbox()
	dbox.Dbox.CreateFolder(tmpFolder)
}
