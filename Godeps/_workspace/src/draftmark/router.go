package draftmark

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func Router() *httprouter.Router {
	r := httprouter.New()

	h := func(res http.ResponseWriter, req *http.Request, p httprouter.Params) {
		fmt.Fprintf(res, "URL: %s\n", req.URL)
		fmt.Fprintf(res, "ID: %s\n", p.ByName("id"))
	}
	r.GET("/a", h)
	r.GET("/b", h)
	r.GET("/c/:id", h)
	r.GET("/d/*path", h)

	return r
}
