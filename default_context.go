package dojo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Masterminds/formenc/encoding/form"
	"github.com/golang/gddo/httputil/header"
	"github.com/tomasen/realip"
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
	data     *sync.Map
	dojo     *Dojo
}

func (ctx *DefaultContext) Dojo() *Dojo {
	return ctx.dojo
}

// Response returns the original Response for the request.
func (ctx *DefaultContext) Response() http.ResponseWriter {
	return ctx.response
}

// Request returns the original Request.
func (ctx *DefaultContext) Request() *http.Request {
	return ctx.request
}

func (ctx *DefaultContext) Params() ParamValues {
	return ctx.params
}

func (ctx *DefaultContext) Set(key string, value interface{}) {
	ctx.data.Store(key, value)
}

// Value that has previously stored on the context.
func (ctx *DefaultContext) Value(key interface{}) interface{} {
	if k, ok := key.(string); ok {
		if v, ok := ctx.data.Load(k); ok {
			return v
		}
	}
	return ctx.Context.Value(key)
}

func (ctx *DefaultContext) Session() *Session {
	return ctx.session
}

// Cookies for the associated request and response.
func (ctx *DefaultContext) Cookies() *Cookies {
	return &Cookies{ctx.request, ctx.response}
}

func (ctx *DefaultContext) Param(key string) string {
	return ctx.Params().Get(key)
}

func (ctx *DefaultContext) RealIP() string {
	return realip.FromRequest(ctx.Request())
}

func (ctx *DefaultContext) Bind(dst interface{}) error {
	if ctx.Request().Header.Get("Content-Type") != "" {
		value, _ := header.ParseValueAndParams(ctx.Request().Header, "Content-Type")
		if value == "application/json" {
			return decodeJSONBody(ctx.Response(), ctx.Request(), dst)
		} else {
			return decodeFormData(ctx.Request(), dst)
		}
	}
	return nil
}

func (ctx *DefaultContext) Data() map[string]interface{} {
	m := map[string]interface{}{}
	ctx.data.Range(func(k, v interface{}) bool {
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

func (ctx *DefaultContext) writeContentType(value string) {
	headers := ctx.Response().Header()
	if headers.Get(HeaderContentType) == "" {
		headers.Set(HeaderContentType, value)
	}
}

type JsonResponseBody struct {
	Data interface{} `json:"data"`
}

func (ctx *DefaultContext) NoContent(code int) error {
	ctx.response.WriteHeader(code)
	_, err := ctx.Response().Write(nil)
	if err != nil {
		return err
	}
	return nil
}

func (ctx *DefaultContext) JSON(code int, data interface{}) error {
	resBody := JsonResponseBody{
		Data: data,
	}
	respData, err := json.Marshal(resBody)
	if err != nil {
		return err
	}
	ctx.writeContentType(MIMEApplicationJSONCharsetUTF8)
	ctx.response.WriteHeader(code)
	_, err = ctx.Response().Write(respData)
	if err != nil {
		return err
	}
	return nil
}
