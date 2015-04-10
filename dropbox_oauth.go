package main

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/oauth2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var conf *oauth2.Config

type dropboxUser struct {
	DisplayName        string `json:"display_name"`
	DropboxId          int    `json:"uid"`
	Email              string `json:"email"`
	DropboxAccessToken string
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

func oauthInit(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	url := conf.AuthCodeURL("state")
	http.Redirect(res, req, url, 302)
}

func oauthRedirect(res http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	code := req.URL.Query().Get("code")
	tok, err := conf.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatal(err)
	}
	dUser := getUserInfo(tok)
	setCurrentUser(res, req, dUser)
	http.Redirect(res, req, "/", 301)
}

func getUserInfo(tok *oauth2.Token) *dropboxUser {
	tok.TokenType = "Bearer"
	client := conf.Client(oauth2.NoContext, tok)
	resp, _ := client.Get("https://api.dropbox.com/1/account/info")
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	var dropboxUser dropboxUser
	json.Unmarshal(body, &dropboxUser)
	return &dropboxUser
}
