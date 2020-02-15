package contract

import "time"

type (
	Configuration interface {
		Get(key string) interface{}
		GetBool(key string) bool
		GetFloat(key string) float64
		GetInt(key string) int
		GetIntSlice(key string) []int
		GetString(key string) string
		GetStringMap(key string) map[string]interface{}
		GetStringMapString(key string) map[string]string
		GetStringSlice(key string) []string
		GetTime(key string) time.Time
		GetDuration(key string) time.Duration
		Exists(key string) bool
		Set(key string, value interface{})
		SetDefault(key string, value interface{})
	}
)
