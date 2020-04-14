package http

import (
	"github.com/firmeve/firmeve/kernel/contract"
	"github.com/gorilla/sessions"
	"net/http"
)

type (
	httpSession struct {
		//Store sessions.Store `inject:"http.session.store"`
		store    sessions.Store
		request  *http.Request
		response http.ResponseWriter
		session  *sessions.Session
	}
)

func NewSession(store sessions.Store, request *http.Request, response http.ResponseWriter) contract.Session {
	s := &httpSession{
		store:    store,
		request:  request,
		response: response,
	}

	session, _ := store.Get(request, "session-name")
	s.session = session

	return s
}

func (s *httpSession) Id() string {
	return s.session.ID
}

func (s *httpSession) Get(key string) interface{} {
	if v, ok := s.session.Values[key]; ok {
		return v
	}
	return nil
}

func (s *httpSession) GetString(key string) string {
	return s.Get(key).(string)
}

func (s *httpSession) GetInt(key string) int {
	return s.Get(key).(int)
}

func (s *httpSession) GetBool(key string) bool {
	return s.Get(key).(bool)
}

func (s *httpSession) GetFloat(key string) float64 {
	return s.Get(key).(float64)
}

func (s *httpSession) Put(key string, value interface{}) error {
	s.session.Values[key] = value
	return s.session.Save(s.request, s.response)
}

func (s *httpSession) Flush() {
	s.session.Values = make(map[interface{}]interface{}, 0)
}

func (s *httpSession) Delete(key string) {
	delete(s.session.Values, key)
}
