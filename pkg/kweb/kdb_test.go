package kweb_test

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/abombard/goserver/pkg/kweb"
)

type Pokemon struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Age  int    `json:"age"`
}

type PokemonPost struct {
	kweb.Handle
	kweb.ReaderJSON
	Pokemon
}

func (e PokemonPost) Do(ctx context.Context) (any, *kweb.Error) {
	fmt.Printf("received pokemon %+v\n", e.Pokemon)

	return nil, nil
}

type PokemonGet struct {
	kweb.Handle
	kweb.WriterJSON
	Name string `query_param:"name,required=true"`
	Type string `query_param:"type,default=plant"`
	Age  int    `query_param:"age"`
}

func (e PokemonGet) Do(ctx context.Context) (any, *kweb.Error) {
	p := Pokemon{
		Name: "pikachu",
		Type: "electric",
		Age:  12,
	}

	if e.Age != 0 {
		p.Age = e.Age
	}

	if e.Name != "" {
		p.Name = e.Name
	}

	if e.Type != "" {
		p.Type = e.Type
	}

	fmt.Printf("sending pokemon %+v\n", p)

	return p, nil
}

type PokemonFetch struct {
	kweb.Handle
	PokemonFetched []Pokemon
}

func (e *PokemonFetch) Fetch(r *http.Request) *kweb.Error {
	for h, values := range r.Header {
		if !strings.HasPrefix(h, "pokemon") {
			continue
		}

		for _, name := range values {
			e.PokemonFetched = append(e.PokemonFetched, Pokemon{
				Name: name,
			})
		}
	}

	return nil
}

func (e PokemonFetch) Do(ctx context.Context) (any, *kweb.Error) {
	fmt.Printf("fetched pokemons %+v\n", e.PokemonFetched)

	return nil, nil
}

type PokemonHeaders struct {
	kweb.Handle
	kweb.ReaderJSON
	kweb.WriterJSON
	HeaderName string `header:"name,required=true"`
	HeaderAge  int    `header:"age"`
	HeaderType int    `header:"type"`
}

func (e PokemonHeaders) Do(ctx context.Context) (any, *kweb.Error) {
	fmt.Printf("received headers name=%s age=%d\n", e.HeaderName, e.HeaderAge)

	return nil, nil
}
