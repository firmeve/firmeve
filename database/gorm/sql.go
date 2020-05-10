package gorm

import (
	"database/sql"
	"time"
)

func NewString(str string) *sql.NullString {
	return &sql.NullString{
		String: str,
		Valid:  true,
	}
}

func NewInt64(n64 int64) *sql.NullInt64 {
	return &sql.NullInt64{
		Int64: n64,
		Valid: true,
	}
}

func NewInt32(n32 int32) *sql.NullInt32 {
	return &sql.NullInt32{
		Int32: n32,
		Valid: true,
	}
}

func NewFloat64(f64 float64) *sql.NullFloat64 {
	return &sql.NullFloat64{
		Float64: f64,
		Valid:   true,
	}
}

func NewBool(b bool) *sql.NullBool {
	return &sql.NullBool{
		Bool:  b,
		Valid: true,
	}
}

func NewTime(time2 time.Time) *sql.NullTime {
	return &sql.NullTime{
		Time:  time2,
		Valid: true,
	}
}
