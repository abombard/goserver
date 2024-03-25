package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

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

type GamePost struct {
	kweb.Fetcher
	kweb.ReaderJSON
	kweb.WriterJSON
	game.Game
}

func (e GamePost) Do(ctx context.Context) (any, *kweb.Error) {
	fmt.Printf("received POST %+v\n", e.Game)

	return nil, nil
}

type GameGet struct {
	kweb.Reader
	kweb.Fetcher
	kweb.WriterJSON
	Id string `url_param:"id"`
}

func (e GameGet) Do(ctx context.Context) (any, *kweb.Error) {
	return game.Game{
		Id: e.Id,
	}, nil
}

type GameDelete struct {
	kweb.Reader
	kweb.Fetcher
	kweb.Writer
	Id string `url_param:"id"`
}

func (e GameDelete) Do(ctx context.Context) (any, *kweb.Error) {
	return nil, nil
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
