package errorsx

import (
	"encoding/json"
	"github.com/zengineDev/dojo"
	"log"
	"net/http"
)

type BadRequestError struct {
	err  error
	body interface{}
}

func BadRequest(err error) *BadRequestError {
	return &BadRequestError{err: err}
}

func BadRequestWithBody(body interface{}) *BadRequestError {
	return &BadRequestError{body: body}
}

func (e *BadRequestError) RespondError(ctx dojo.Context, application *dojo.Application) bool {
	if e.body == nil {
		http.Error(ctx.Response(), http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	} else {
		ctx.Response().WriteHeader(http.StatusBadRequest)

		ctx.Response().Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(ctx.Response()).Encode(e.body)

		if err != nil {
			log.Printf("Failed to encode a response: %v", err)
		}
	}

	return true
}

func (e *BadRequestError) Error() string {
	return e.err.Error()
}
