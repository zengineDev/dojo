package dojo

import (
	"context"
	"net/http"
	"net/url"
	"sync"
)

type DefaultContext struct {
	context.Context
	response http.ResponseWriter
	request  *http.Request
	params   url.Values
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
