package kweb

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
)

// Response
type Response struct {
	Status int
	Body   any
}

type Handler interface {
	// Read the request Headers, Params, Body into handleStruct
	Read(r *http.Request, handleStruct Handler) *Error
	// Fetch whatever is needed from the http request to put into the Handler (headers, query params, ..)
	Fetch(r *http.Request) *Error
	// Do executes the api call operational code without access to the writer or the request and returns
	// the response as a struct, or an error
	Do(ctx context.Context) (any, *Error)
	// Write format the response output
	Write(w http.ResponseWriter, r *http.Request, resp any)
}

// NewHandler allocates a new handle and executes each method from the Handler interface
func NewHandler(handleStruct Handler) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		handle := reflect.New(reflect.TypeOf(handleStruct).Elem()).Interface().(Handler)

		var (
			resp any
			err  *Error
		)

		if err = handle.Read(r, handle); err != nil {
			goto onerror
		}

		if err = handle.Fetch(r); err != nil {
			goto onerror
		}

		if resp, err = handle.Do(r.Context()); err != nil {
			goto onerror
		}

		handle.Write(w, r, resp)
		return

	onerror:
		handle.Write(w, r, err)
	}
}

type Handle struct {
	Fetcher
	Reader
	Writer
}

type Fetcher struct{}

func (handle *Fetcher) Fetch(r *http.Request) *Error {
	return nil
}

type Reader struct{}

func (handle *Reader) Read(r *http.Request, handleStruct Handler) *Error {
	t := reflect.TypeOf(handleStruct).Elem()
	v := reflect.ValueOf(handleStruct)

	return DecodeHTTPTags(r, t, v)
}

// ReaderJSON overrides the Read method of Handler, it takes the json
// body from the Request and read it into the handle herited structure.
type ReaderJSON struct {
	Reader
}

func (handle ReaderJSON) Read(r *http.Request, handleStruct Handler) *Error {
	if err := handle.Reader.Read(r, handleStruct); err != nil {
		return err
	}

	// r.Body is never nil but Unmarshal returns io.EOF if Body is empty
	err := json.NewDecoder(r.Body).Decode(handleStruct)
	if err != nil && !errors.Is(err, io.EOF) {
		return &Error{
			Message: fmt.Sprintf("failed to decode request body: %v", err),
			Status:  http.StatusBadRequest,
			Err:     err,
		}
	}

	return nil
}

// Writer
type Writer struct {
	WriteFn func(w http.ResponseWriter, code int, obj any)
}

func (handle *Writer) WriteText(w http.ResponseWriter, code int, obj any) {
	w.WriteHeader(code)

	if obj == nil {
		return
	}

	v := fmt.Sprintf("%v", obj)

	_, err := w.Write([]byte(v))
	if err != nil {
		fmt.Println("kweb: failed to write response:", err)
	}
}

func (handle *Writer) WriteFnJSON(w http.ResponseWriter, code int, obj any) {
	w.WriteHeader(code)

	if obj == nil {
		return
	}

	if err := json.NewEncoder(w).Encode(obj); err != nil {
		handle.WriteText(w, code, obj)
	}
}

func (handle *Writer) Write(w http.ResponseWriter, r *http.Request, resp any) {
	if handle.WriteFn == nil {
		handle.WriteFn = handle.WriteText
	}

	switch cast := resp.(type) {
	case nil:
		handle.WriteFn(w, http.StatusOK, struct{}{})
	case *Response:
		handle.WriteFn(w, cast.Status, cast.Body)
	case *Error:
		handle.WriteFn(w, cast.Status, cast.Message)
	default:
		handle.WriteFn(w, http.StatusOK, resp)
	}
}

type WriterJSON struct {
	Writer
}

func (handle *WriterJSON) Write(w http.ResponseWriter, r *http.Request, response any) {
	handle.WriteFn = handle.WriteFnJSON
	handle.Writer.Write(w, r, response)
}
