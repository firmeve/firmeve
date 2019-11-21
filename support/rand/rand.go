package rand

import (
	"math/rand"
	"time"
)

func RangeInt(min, max int) int {
	source := rand.New(rand.NewSource(time.Now().UnixNano()))
	return source.Intn(max-min) + min
}
