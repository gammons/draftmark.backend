package main

import (
	"draftmark"
	db "draftmark/persistence"
	"encoding/json"
	"github.com/go-martini/martini"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"strconv"
)

var database = &db.Client{}
var user = &db.User{}
var sync *draftmark.Sync

func setupDatabase() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	database.InitDB()
	database.Db.LogMode(true)
	//resetDB()

	database.Db.Where("email = ?", "gammons@gmail.com").First(&user)
}

func resetDB() {
	database.Db.DropTable(&db.User{})
	database.Db.DropTable(&db.Note{})
	database.Db.CreateTable(&db.User{})
	database.Db.CreateTable(&db.Note{})
	database.Db.AutoMigrate(&db.User{}, &db.Note{})
	database.Db.Model(&db.Note{}).AddIndex("idx_user_notes", "user_id")
	user := &db.User{Email: "gammons@gmail.com", DropboxAccessToken: "RzfZv3hAoIYAAAAAAAAJ2ue-DKJPep3jvHF3XNGvvjJk-gDHkgUvOUyOcxH4XG_V", DropboxCursor: ""}
	database.Db.Create(&user)

}

func listNotes() ([]byte, error) {
	notes := database.ListNotes(user)
	return json.Marshal(notes)
}

func getNote(params martini.Params) string {
	noteId, _ := strconv.Atoi(params["id"])
	return database.GetNoteContents(noteId)
}

func setupMartini() {
	m := martini.Classic()
	static := martini.Static("public", martini.StaticOptions{Fallback: "/index.html"})
	m.Get("/notes", listNotes)
	m.Get("/sync", func() {
		go sync.DoSync(*user, "/notes")
	})
	m.Get("/note/:id", getNote)
	m.NotFound(static, http.NotFound)
	m.Run()
}

func main() {
	setupDatabase()
	log.Println("creating new sync with ", user.DropboxAccessToken)
	sync = draftmark.NewSync(user.DropboxAccessToken)
	sync.DoLogging = true
	setupMartini()
}
