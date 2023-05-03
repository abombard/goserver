package web

import (
	"github.com/abombard/goserver/pkg/kweb"
)

type GamePost struct {
	kweb.Fetcher
	kweb.ReaderJSON
	kweb.WriterJSON
}

type GameGet struct {
	kweb.Reader
	kweb.WriterJSON
}

type GameDelete struct {
	kweb.Reader
	kweb.Writer
}
