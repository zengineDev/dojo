package dojo

import (
	"github.com/gorilla/sessions"
	"net/http"
)

type Session struct {
	Session *sessions.Session
	req     *http.Request
	res     http.ResponseWriter
}

func (s *Session) Flash(key string, value interface{}) {
	s.Session.AddFlash(value, key)
}

func (s *Session) GetFlash(key string) []interface{} {
	return s.Session.Flashes(key)
}

func (s *Session) Save() error {
	return s.Session.Save(s.req, s.res)
}

func (s *Session) Get(name interface{}) interface{} {
	return s.Session.Values[name]
}

func (s *Session) GetOnce(name interface{}) interface{} {
	if x, ok := s.Session.Values[name]; ok {
		s.Delete(name)
		return x
	}
	return nil
}

func (s *Session) Set(name, value interface{}) {
	s.Session.Values[name] = value
}

func (s *Session) Delete(name interface{}) {
	delete(s.Session.Values, name)
}

func (s *Session) Clear() {
	for k := range s.Session.Values {
		s.Delete(k)
	}
}

func (app *Application) getSession(r *http.Request, w http.ResponseWriter) *Session {
	session, _ := app.SessionStore.Get(r, app.Configuration.Session.Name)
	return &Session{
		Session: session,
		req:     r,
		res:     w,
	}
}
