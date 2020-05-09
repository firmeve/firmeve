package gorm

import (
	"database/sql/driver"
	"encoding/json"
)

type (
	StringArray struct {
		Source []string
		Valid  bool
	}

	IntArray struct {
		Source []int
		Valid  bool
	}
)

func (s *StringArray) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), &s.Source)
}

func (s *StringArray) Value() (driver.Value, error) {
	if s == nil {
		return nil, nil
	}
	if !s.Valid {
		return nil, nil
	}

	return json.Marshal(s.Source)
}

func (s *IntArray) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), &s.Source)
}

func (s *IntArray) Value() (driver.Value, error) {
	if s == nil {
		return nil, nil
	}
	if !s.Valid {
		return nil, nil
	}

	return json.Marshal(s.Source)
}

func NewStringArray(source []string) *StringArray {
	return &StringArray{
		Source: source,
		Valid:  true,
	}
}

func NewIntArray(source []int) *IntArray {
	return &IntArray{
		Source: source,
		Valid:  true,
	}
}
