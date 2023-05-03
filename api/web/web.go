package web

import (
	"context"
	"log"
	"net/http"

	"github.com/abombard/goserver/pkg/kweb"
	"github.com/gorilla/mux"
)

type Ping struct {
	kweb.Handle
}

func (wctx Ping) Do(ctx context.Context) (any, *kweb.Error) {
	return "pong", nil
}

var routes = []struct {
	url     string
	method  string
	handler kweb.Handler
}{
	{"/ping", http.MethodGet, (*Ping)(nil)},
	{"/game", http.MethodPost, (*GamePost)(nil)},
	{"/game", http.MethodPut, (*GamePost)(nil)},
	{"/game", http.MethodGet, (*GameGet)(nil)},
	{"/game", http.MethodDelete, (*GameDelete)(nil)},
}

func Run(srv http.Server) {
	r := mux.NewRouter()

	for _, route := range routes {
		r.HandleFunc(
			route.url,
			kweb.NewHandler(route.handler),
		).Methods(route.method)
	}

	srv.Handler = r

	log.Fatal(srv.ListenAndServe())
}
