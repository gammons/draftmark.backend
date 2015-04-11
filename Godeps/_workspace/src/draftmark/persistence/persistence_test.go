package persistence_test

import (
	database "draftmark/persistence"
	//"fmt"
	"github.com/joho/godotenv"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"log"
)

var db = &database.Client{}
var _ = BeforeSuite(func() {
	err := godotenv.Load("../test.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db.InitDB()
	//db.Db.LogMode(true)
	db.Db.CreateTable(&database.User{})
	db.Db.CreateTable(&database.Note{})
})

var _ = AfterSuite(func() {
	db.Db.DropTable(&database.User{})
	db.Db.DropTable(&database.Note{})
})

var user = database.User{Email: "gammons@gmail.com"}

var _ = Describe("ListNotes", func() {
	var note1 = &database.Note{Path: "/notes/test.md", UserId: 1}
	var note2 = &database.Note{Path: "/notes/test2.md", UserId: 1}
	var _ = &database.Note{Path: "/notes/test3.md", UserId: 2}

	It("Lists notes for a user", func() {
		db.Db.Create(&user)
		db.Db.Create(note1)
		db.Db.Create(note2)
		notes := db.ListNotes(&user)

		var paths []string
		for _, note := range notes {
			paths = append(paths, note.Path)
		}
		var expected = []string{"/notes/test.md", "/notes/test2.md"}
		Expect(paths).To(Equal(expected))
	})

})

var _ = Describe("GetNoteContents", func() {
	var note = &database.Note{Path: "/notes/test.md", UserId: 5, Content: "this is the content"}
	It("Gets the note contents", func() {
		db.Db.Create(note)
		Expect(db.GetNoteContents(note)).To(Equal("this is the content"))
	})
})
