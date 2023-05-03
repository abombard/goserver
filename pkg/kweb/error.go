package kweb

import (
	"fmt"
	"net/http"
)

type Error struct {
	Status  int    `json:"-"`
	Message string `json:"error,omitempty"`
	Err     error  `json:"-"`
	Kind    string `json:"-"`
}

func (e Error) Error() string {
	return fmt.Sprintf("%d %s: %s", e.Status, http.StatusText(e.Status), e.Message)
}

func ErrorWrap(err error) *Error {
	if err == nil {
		return nil
	}

	switch cast := err.(type) {
	case *Error:
		return cast
	default:
	}

	return &Error{
		Status:  http.StatusInternalServerError,
		Message: err.Error(),
		Err:     err,
	}
}

func ErrorWrapMsg(err error, msg string) *Error {
	e := ErrorWrap(err)

	e.Message = msg

	return e
}
