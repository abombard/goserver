package main

import (
	"net/http"

	"github.com/abombard/goserver/api/web"
)

func main() {
	web.Run(http.Server{
		Addr: ":8080",
	})
}
