package dojo

import (
	"fmt"
	"github.com/gorilla/sessions"
	"net/http"
)

type Session struct {
	Session *sessions.Session
	req     *http.Request
	res     http.ResponseWriter
}

func (s *Session) WithOld(data map[string]interface{}) {
	s.Session.AddFlash(data, FlashOldKey)
	_ = s.Save()
}

func (s *Session) Flash(key string, value interface{}) {
	s.Session.AddFlash(value, key)
	err := s.Save()
	fmt.Println(err)
}

func (s *Session) GetFlash(key string) []interface{} {
	m := s.Session.Flashes(key)
	err := s.Save()
	fmt.Println(err)
	return m
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
