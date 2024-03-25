package main

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/abombard/goserver/model/game"
	"github.com/abombard/goserver/pkg/kweb"
	"github.com/gorilla/mux"
)

type Ping struct {
	kweb.Handle
}

func (wctx Ping) Do(ctx context.Context) (any, *kweb.Error) {
	return "pong", nil
}

var games = map[string]game.Game{}

type GamePost struct {
	kweb.Fetcher
	kweb.ReaderJSON
	kweb.WriterJSON
	game.Game
}

func (e GamePost) Do(ctx context.Context) (any, *kweb.Error) {
	games[e.Id] = e.Game

	return kweb.Response{Status: http.StatusCreated}, nil
}

type GameGet struct {
	kweb.Reader
	kweb.Fetcher
	kweb.WriterJSON
	Id string `url_param:"id"`
}

func (e GameGet) Do(ctx context.Context) (any, *kweb.Error) {
	if g, ok := games[e.Id]; ok {
		return g, nil
	}

	return kweb.Error{
		Status:  http.StatusNotFound,
		Message: "game does not exist",
	}, nil
}

type GameDelete struct {
	kweb.Reader
	kweb.Fetcher
	kweb.Writer
	Id string `url_param:"id"`
}

func (e GameDelete) Do(ctx context.Context) (any, *kweb.Error) {
	if _, ok := games[e.Id]; ok {
		delete(games, e.Id)

		return kweb.Response{Status: http.StatusOK}, nil
	}

	return kweb.Error{
		Status:  http.StatusNotFound,
		Message: "game does not exist",
	}, nil
}

type Authorization struct {
	Header string `header:"Authorization"`
	token  string
}

func (auth *Authorization) Fetch(r *http.Request) *kweb.Error {
	auth.token = strings.TrimPrefix(auth.Header, "Bearer ")

	// decode token and check auth

	return nil
}

var routes = []struct {
	url     string
	method  string
	handler kweb.Handler
}{
	{"/ping", http.MethodGet, (*Ping)(nil)},
	{"/game", http.MethodPost, (*GamePost)(nil)},
	{"/game/{id}", http.MethodPut, (*GamePost)(nil)},
	{"/game/{id}", http.MethodGet, (*GameGet)(nil)},
	{"/game/{id}", http.MethodDelete, (*GameDelete)(nil)},
}

func Run(srv *http.Server) {
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

func main() {
	Run(&http.Server{
		Addr: ":8080",
	})
}
