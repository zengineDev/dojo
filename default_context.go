package dojo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Masterminds/formenc/encoding/form"
	"github.com/golang/gddo/httputil/header"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

type DefaultContext struct {
	context.Context
	response http.ResponseWriter
	request  *http.Request
	params   url.Values
	session  *Session
	flash    *Flash
	data     *sync.Map
}

// Response returns the original Response for the request.
func (d *DefaultContext) Response() http.ResponseWriter {
	return d.response
}

// Request returns the original Request.
func (d *DefaultContext) Request() *http.Request {
	return d.request
}

func (d *DefaultContext) Params() ParamValues {
	return d.params
}

func (d *DefaultContext) Set(key string, value interface{}) {
	d.data.Store(key, value)
}

// Value that has previously stored on the context.
func (d *DefaultContext) Value(key interface{}) interface{} {
	if k, ok := key.(string); ok {
		if v, ok := d.data.Load(k); ok {
			return v
		}
	}
	return d.Context.Value(key)
}

func (d *DefaultContext) Session() *Session {
	return d.session
}

// Cookies for the associated request and response.
func (d *DefaultContext) Cookies() *Cookies {
	return &Cookies{d.request, d.response}
}

// Flash messages for the associated Request.
func (d *DefaultContext) Flash() *Flash {
	return d.flash
}

func (d *DefaultContext) Param(key string) string {
	return d.Params().Get(key)
}

func (d *DefaultContext) Bind(dst interface{}) error {
	if d.Request().Header.Get("Content-Type") != "" {
		value, _ := header.ParseValueAndParams(d.Request().Header, "Content-Type")
		if value == "application/json" {
			return decodeJSONBody(d.Response(), d.Request(), dst)
		} else {
			return decodeFormData(d.Request(), dst)
		}
	}
	return nil
}

func (d *DefaultContext) Data() map[string]interface{} {
	m := map[string]interface{}{}
	d.data.Range(func(k, v interface{}) bool {
		s, ok := k.(string)
		if !ok {
			return false
		}
		m[s] = v
		return true
	})
	return m
}

type malformedRequest struct {
	status int
	msg    string
}

func (mr *malformedRequest) Error() string {
	return mr.msg
}

func decodeFormData(r *http.Request, dst interface{}) error {
	if err := form.Unmarshal(r.PostForm, dst); err != nil {
		return err
	}
	return nil
}

func decodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	if r.Header.Get("Content-Type") != "" {
		value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
		if value != "application/json" {
			msg := "Content-Type header is not application/json"
			return &malformedRequest{status: http.StatusUnsupportedMediaType, msg: msg}
		}
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case errors.Is(err, io.ErrUnexpectedEOF):
			msg := "Request body contains badly-formed JSON"
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case errors.Is(err, io.EOF):
			msg := "Request body must not be empty"
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case err.Error() == "http: request body too large":
			msg := "Request body must not be larger than 1MB"
			return &malformedRequest{status: http.StatusRequestEntityTooLarge, msg: msg}

		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		msg := "Request body must only contain a single JSON object"
		return &malformedRequest{status: http.StatusBadRequest, msg: msg}
	}

	return nil
}
