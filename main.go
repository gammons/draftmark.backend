package main

import (
	"draftmark"
	db "draftmark/persistence"
	"encoding/json"
	"fmt"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
	"github.com/unrolled/render"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

var database = &db.Client{}
var user = &db.User{}
var sync *draftmark.Sync
var rndr = render.New()
var store = sessions.NewCookieStore([]byte("dingleton"))

func setupDatabase() {
	database.InitDB()
	database.Db.LogMode(true)
}

func resetDB(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	database.Db.DropTable(&db.User{})
	database.Db.DropTable(&db.Note{})
	database.Db.CreateTable(&db.User{})
	database.Db.CreateTable(&db.Note{})
	database.Db.AutoMigrate(&db.User{}, &db.Note{})
	database.Db.Model(&db.Note{}).AddIndex("idx_user_notes", "user_id")
}

func listNotes(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	notes := database.ListNotes(currentUser(req))
	res.Header().Set("Content-Type", "application/json; charset=UTF-8")
	rndr.JSON(res, http.StatusOK, notes)
}

func getNote(res http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	path := ps.ByName("filepath")
	fmt.Fprintf(res, database.GetNoteContents(currentUser(req), path))
}

type SyncUsers struct {
	Delta struct {
		Users []int `json:"users"`
	} `json:"delta"`
}

func doSync(res http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	if challenge := req.URL.Query().Get("challenge"); challenge != "" {
		fmt.Fprintf(res, challenge)
		return
	}
	body, _ := ioutil.ReadAll(req.Body)
	var syncUsers SyncUsers
	json.Unmarshal(body, &syncUsers)
	for _, userId := range syncUsers.Delta.Users {
		var user = db.User{}
		database.Db.Where("dropbox_user_id = ?", userId).First(&user)
		go sync.DoSync(user, "/notes")
	}
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

	router.GET("/authorize", oauthInit)
	router.GET("/redirect", oauthRedirect)
	router.GET("/sync", doSync)
	router.POST("/sync", doSync)
	router.GET("/reset__", resetDB)

	n.UseHandler(router)
	n.Run(":3000")
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file.  This might be ok though if we're in prod.")
	}
	setupDatabase()
	setupOauth()
	sync = draftmark.NewSync()
	sync.DoLogging = true
	setupNegroni()
}
