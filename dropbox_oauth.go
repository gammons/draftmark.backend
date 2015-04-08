package main

import (
	"github.com/julienschmidt/httprouter"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"os"
)

var conf *oauth2.Config

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
	code := ps.ByName("code")
	tok, err := conf.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatal(err)
	}
	user.DropboxAccessToken = tok.AccessToken
	sync.Db.UpdateUserAccessToken(user, tok.AccessToken)
}
