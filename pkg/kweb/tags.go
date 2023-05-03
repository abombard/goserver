package kweb

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

var ErrDev = errors.New("fatal error")

func tagDefault(param, defaultValue string) (string, error) {
	if param == "" {
		param = defaultValue
	}

	return param, nil
}

func tagRequired(param, boolean string) (string, error) {
	isRequired, err := strconv.ParseBool(boolean)
	if err != nil {
		return param, fmt.Errorf("strconv.ParseBool failed on `%s`: %w", boolean, ErrDev)
	}

	if isRequired && param == "" {
		return param, errors.New("missing required parameter")
	}

	return param, nil
}

// keep in mind the order matters
var tagOptions = map[string]func(string, string) (string, error){
	"required": tagRequired,
	"default":  tagDefault,
}

func doTagOptions(tagValues []string, param string) (string, error) {
	if len(tagValues) <= 1 {
		return param, nil
	}

	for _, tagValue := range tagValues[1:] {
		opt := strings.Split(tagValue, "=")
		if len(opt) != 2 {
			return param, fmt.Errorf("invalid tag option `%s`: %w", tagValue, ErrDev)
		}

		optFn, ok := tagOptions[opt[0]]
		if !ok {
			return param, fmt.Errorf("invalid tag option `%s`: %w", opt[0], ErrDev)
		}

		var err error

		param, err = optFn(param, opt[1])
		if err != nil {
			return param, err
		}
	}

	return param, nil
}

type urlParams map[string]string

func (m urlParams) Get(key string) string {
	return m[key]
}

func DecodeHTTPTags(req *http.Request, t reflect.Type, v reflect.Value) *Error {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		tag := ""
		param := ""
		tagValues := []string{}

		maps := map[string]func(string) string{
			"header":      req.Header.Get,
			"query_param": req.URL.Query().Get,
			"url_param":   urlParams(mux.Vars(req)).Get,
		}

		var getParamFn func(string) string

		for tag, getParamFn = range maps {
			if tagValue := field.Tag.Get(tag); tagValue != "" {
				tagValues = strings.Split(tagValue, ",")
				param = getParamFn(tagValues[0])
				break
			}
		}

		var err error

		param, err = doTagOptions(tagValues, strings.Trim(param, " "))
		if err != nil {
			if errors.Is(err, ErrDev) {
				return &Error{
					Status:  http.StatusInternalServerError,
					Message: fmt.Sprintf("failed to decode `%s`: %s", tag, field.Name),
					Err: fmt.Errorf(
						"doTagOptions failed for tag `%s` values `%v` param `%s` on field `%s` in struct `%s`",
						tag, tagValues, param, field.Name, t.Name()),
				}
			}

			return &Error{
				Status:  http.StatusBadRequest,
				Message: fmt.Sprintf("failed to decode `%s` `%s`: %v", tag, field.Name, err),
				Err: fmt.Errorf("failed to decode `%s` `%s` for field `%s` in struct `%s`: %w",
					tag, param, field.Name, t.Name(), err),
			}
		}

		if param == "" {
			continue
		}

		fieldValue := v.Elem().Field(i)

		var value reflect.Value

		switch fieldValue.Interface().(type) {
		case bool:
			var v bool

			v, err = strconv.ParseBool(strings.ToLower(param))
			if err != nil {
				break
			}

			value = reflect.ValueOf(v)
		case int, int64:
			// not gonna work if int = 32
			var v int64

			v, err = strconv.ParseInt(param, 10, 64)
			if err != nil {
				break
			}

			value = reflect.ValueOf(int(v))
		case string:
			value = reflect.ValueOf(param)
		case []int, []int64:
			value = reflect.ValueOf([]int{})
			params := strings.Split(param, ",")

			var v int64

			for _, e := range params {
				e = strings.Trim(e, " ")

				v, err = strconv.ParseInt(e, 10, 64)
				if err != nil {
					break
				}

				value = reflect.Append(value, reflect.ValueOf(v))
			}
		case []string:
			value = reflect.ValueOf([]string{})
			params := strings.Split(param, ",")

			for _, e := range params {
				e = strings.Trim(e, " ")

				value = reflect.Append(value, reflect.ValueOf(e))
			}
		default:
			return &Error{
				Status:  http.StatusInternalServerError,
				Message: fmt.Sprintf("unhandled type `%s` for tag `%s`", field.Type.Name(), tag),
				Err: fmt.Errorf(
					"unhandled type `%s` for tag `%s` on field `%s` in struct `%s`",
					field.Type.Name(), tag, field.Name, t.Name(),
				),
			}
		}

		if err != nil {
			return &Error{
				Status:  http.StatusBadRequest,
				Message: fmt.Sprintf("invalid parameter `%s`", tagValues[0]),
				Err: fmt.Errorf(
					"invalid param type `%s` for tag `%s` on field `%s` in struct `%s`: %v",
					field.Type.Name(), tag, field.Name, t.Name(), err,
				),
			}
		}

		fieldValue.Set(value)
	}

	return nil
}
