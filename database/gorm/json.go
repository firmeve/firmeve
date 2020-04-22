package gorm

import (
	"database/sql/driver"
	"encoding/json"
)

type (
	JSON struct {
		Source map[string]interface{}
		Valid  bool
	}
)

func (s *JSON) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), &s.Source)
}

func (s *JSON) Value() (driver.Value, error) {
	if s == nil {
		return nil, nil
	}
	if !s.Valid {
		return nil, nil
	}

	return json.Marshal(s.Source)
}

func NewJSON(source map[string]interface{}) *JSON {
	return &JSON{
		Source: source,
		Valid:  true,
	}
}
