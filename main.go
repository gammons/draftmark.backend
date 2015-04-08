package main

import (
	"draftmark"
	db "draftmark/persistence"
	"github.com/codegangsta/negroni"
	"github.com/unrolled/render"
	//"github.com/goincremental/negroni-sessions"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"strconv"
)

var database = &db.Client{}
var user = &db.User{}
var sync *draftmark.Sync
var rndr = render.New()

func setupDatabase() {
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

func listNotes(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	notes := database.ListNotes(user)
	res.Header().Set("Content-Type", "application/json; charset=UTF-8")
	rndr.JSON(res, http.StatusOK, notes)
}

func getNote(res http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	noteId, _ := strconv.Atoi(ps.ByName("id"))
	fmt.Fprintf(res, database.GetNoteContents(noteId))
}

func doSync(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	go sync.DoSync(*user, "/notes")
}

func setupNegroni() {
	n := negroni.Classic()
	router := httprouter.New()
	router.GET("/notes.json", listNotes)
	router.GET("/notes/:id/content.json", getNote)

	router.GET("/authorize", oauthInit)
	router.GET("/redirect", oauthRedirect)
	router.GET("/sync", doSync)

	n.UseHandler(router)
	n.Run(":3000")

	// n.Get("/notes.json", listNotes)
	// n.Get("/sync", func() {
	// 	go sync.DoSync(*user, "/notes")
	// })
	// n.Get("/notes/:id/content.json", getNote)
	// n.Get("/authorize", oauthInit)
	// n.Get("/redirect", oauthRedirect)
	// n.NotFound(static, http.NotFound)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	setupDatabase()
	setupOauth()
	sync = draftmark.NewSync()
	sync.DoLogging = true
	setupNegroni()
}
