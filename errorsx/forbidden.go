package errorsx

import (
	"encoding/json"
	"github.com/zengineDev/dojo"
	"log"
	"net/http"
)

type ForbiddenError struct {
	err  error
	body interface{}
}

func Forbidden(err error) *ForbiddenError {
	return &ForbiddenError{err: err}
}

func ForbiddenWithBody(body interface{}) *ForbiddenError {
	return &ForbiddenError{body: body}
}

func (e *ForbiddenError) RespondError(ctx dojo.Context, app *dojo.Application) bool {
	if e.body == nil {
		http.Error(ctx.Response(), http.StatusText(http.StatusForbidden), http.StatusForbidden)
	} else {
		ctx.Response().WriteHeader(http.StatusForbidden)

		ctx.Response().Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(ctx.Response()).Encode(e.body)

		if err != nil {
			log.Printf("Failed to encode a response: %v", err)
		}
	}

	return true
}

func (e *ForbiddenError) Error() string {
	return e.err.Error()
}
