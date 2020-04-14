package jwt

import (
	"github.com/firmeve/firmeve/kernel/contract"
	"time"
)

type Memory struct {
	jwts      map[string]int64
	audiences map[string][]string
}

func NewMemoryStore() contract.JwtStore {
	return &Memory{
		jwts:      make(map[string]int64, 0),
		audiences: make(map[string][]string, 0),
	}
}

func (m *Memory) Has(id string) bool {
	if v, ok := m.jwts[id]; ok {
		// 判断时间
		if v <= time.Now().Unix() {
			// clear过期数据
			m.Clear(id)
			return false
		}

		return true
	}

	return false
}

func (m *Memory) Put(id string, audience contract.JwtAudience, lifetime time.Time) error {
	m.jwts[id] = lifetime.Unix()
	aud := audience.Audience()
	if _, ok := m.audiences[aud]; !ok {
		m.audiences[aud] = make([]string, 0)
	}
	m.audiences[aud] = append(m.audiences[aud], id)
	return nil
}

func (m *Memory) Forget(audience contract.JwtAudience) error {
	aud := audience.Audience()
	if v, ok := m.audiences[aud]; ok {
		for _, id := range v {
			m.Clear(id)
		}
		delete(m.audiences, aud)
	}

	return nil
}

func (m *Memory) Clear(id string) {
	delete(m.jwts, id)
}
