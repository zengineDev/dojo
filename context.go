package dojo

import (
	"context"
	"net/http"
)

type Context interface {
	context.Context
	Response() http.ResponseWriter
	Request() *http.Request
	Session() *Session
	Cookies() *Cookies
	Dojo() *Dojo
	Params() ParamValues
	Param(string) string
	Set(string, interface{})
	Bind(interface{}) error
	Data() map[string]interface{}
	JSON(code int, data interface{}) error
	NoContent(code int) error
}

type ParamValues interface {
	Get(string) string
}
