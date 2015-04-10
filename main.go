package main

import (
	"draftmark"
	db "draftmark/persistence"
	"fmt"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
	"github.com/unrolled/render"
	"log"
	"net/http"
	"strconv"
)

var database = &db.Client{}
var user = &db.User{}
var sync *draftmark.Sync
var rndr = render.New()
var store = sessions.NewCookieStore([]byte("something-very-secret"))

func setupDatabase() {
	database.InitDB()
	database.Db.LogMode(true)
	//resetDB()

	//database.Db.Where("email = ?", "gammons@gmail.com").First(&user)
}

// func resetDB() {
// 	database.Db.DropTable(&db.User{})
// 	database.Db.DropTable(&db.Note{})
// 	database.Db.CreateTable(&db.User{})
// 	database.Db.CreateTable(&db.Note{})
// 	database.Db.AutoMigrate(&db.User{}, &db.Note{})
// 	database.Db.Model(&db.Note{}).AddIndex("idx_user_notes", "user_id")
// 	user := &db.User{Email: "gammons@gmail.com", DropboxAccessToken: "RzfZv3hAoIYAAAAAAAAJ2ue-DKJPep3jvHF3XNGvvjJk-gDHkgUvOUyOcxH4XG_V", DropboxCursor: ""}
// 	database.Db.Create(&user)
//
// }
//
func listNotes(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	notes := database.ListNotes(currentUser(req))
	res.Header().Set("Content-Type", "application/json; charset=UTF-8")
	rndr.JSON(res, http.StatusOK, notes)
}

func getNote(res http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	path := ps.ByName("filepath")
	fmt.Fprintf(res, database.GetNoteContents(currentUser(req), path))
}

func doSync(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	go sync.DoSync(*currentUser(req), "/notes")
}

func AuthRequired(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		if currentUser(r).ID == 0 {
			http.Redirect(w, r, "/notloggedin", 301)
		} else {
			h(w, r, ps)
		}
	}
}

func currentUser(r *http.Request) *db.User {
	session, _ := store.Get(r, "draftmark")
	var user = &db.User{}
	database.Db.Where("id = ?", session.Values["userId"]).First(&user)
	return user

}

func setCurrentUser(w http.ResponseWriter, r *http.Request, dUser *dropboxUser) {
	session, _ := store.Get(r, "draftmark")
	var user = db.User{}
	database.Db.Where("email = ?", dUser.Email).FirstOrInit(&user)
	user.Email = dUser.Email
	user.DropboxUserId = strconv.Itoa(dUser.DropboxId)
	user.DropboxAccessToken = dUser.DropboxAccessToken
	database.Db.Save(&user)
	session.Values["userId"] = user.ID
	session.Save(r, w)
}

func setupNegroni() {
	n := negroni.Classic()
	router := httprouter.New()
	router.GET("/notes.json", AuthRequired(listNotes))
	router.GET("/content/*filepath", AuthRequired(getNote))
	router.GET("/sync", AuthRequired(doSync))

	router.GET("/authorize", oauthInit)
	router.GET("/redirect", oauthRedirect)

	n.UseHandler(router)
	n.Run(":3000")
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
