package main

import (
	"draftmark"
	db "draftmark/persistence"
	"encoding/json"
	//"fmt"
	"github.com/go-martini/martini"
	"github.com/joho/godotenv"
	"github.com/martini-contrib/sessions"
	"golang.org/x/oauth2"
	//"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

var database = &db.Client{}
var user = &db.User{}
var sync *draftmark.Sync
var conf *oauth2.Config
var session sessions.CookieStore

func setupDatabase() {
	database.InitDB()
	database.Db.LogMode(true)
	//resetDB()

	database.Db.Where("email = ?", "gammons@gmail.com").First(&user)
}

func setupOauth() {
	conf = &oauth2.Config{
		ClientID:     os.Getenv("DROPBOX_KEY"),
		ClientSecret: os.Getenv("DROPBOX_SECRET"),
		RedirectURL:  "http://localhost:3000/redirect",
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://www.dropbox.com/1/oauth2/authorize",
			TokenURL: "https://api.dropbox.com/1/oauth2/token",
		},
	}
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

func oauthInit(res http.ResponseWriter, req *http.Request) {
	url := conf.AuthCodeURL("state")
	http.Redirect(res, req, url, 302)
}

func oauthRedirect(w http.ResponseWriter, r *http.Request, session sessions.Session) string {
	//code := r.URL.Query().Get("code")
	// tok, err := conf.Exchange(oauth2.NoContext, code)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	//tok.TokenType = "Bearer"
	// session.Set("token", tok)
	// log.Println(session.Get("token"))
	// 	client := conf.Client(oauth2.NoContext, tok)
	// 	log.Println(tok)
	// 	resp, err := client.Get("https://api.dropbox.com/1/account/info?locale=en")
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	defer resp.Body.Close()
	// 	body, err := ioutil.ReadAll(resp.Body)
	// 	json.Unmarshal(string(body), )
	//
	// 	return tok.AccessToken
	return "ok"
}

func setupMartini() {
	m := martini.Classic()
	static := martini.Static("public", martini.StaticOptions{Fallback: "/index.html"})
	session = sessions.NewCookieStore([]byte("asdfasdf"))
	m.Use(sessions.Sessions("draftmark_session", session))

	m.Get("/notes.json", listNotes)
	m.Get("/sync", func() {
		go sync.DoSync(*user, "/notes")
	})
	m.Get("/note/:id", getNote)
	m.Get("/authorize", oauthInit)
	m.Get("/redirect", oauthRedirect)
	m.NotFound(static, http.NotFound)
	m.Run()
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	setupDatabase()
	setupOauth()
	log.Println("creating new sync with ", user.DropboxAccessToken)
	sync = draftmark.NewSync(user.DropboxAccessToken)
	sync.DoLogging = true
	setupMartini()
}
