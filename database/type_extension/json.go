package type_extension

import (
	"database/sql/driver"
	"encoding/json"
)

type (
	Json struct {
		Source map[string]interface{}
		Valid  bool
	}
)

func (s *Json) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), &s.Source)
}

func (s *Json) Value() (driver.Value, error) {
	if s == nil {
		return nil, nil
	}
	if !s.Valid {
		return nil, nil
	}

	return json.Marshal(s.Source)
}
