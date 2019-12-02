package input

import (
	"github.com/firmeve/firmeve/input/parser"
)

type (
	Input struct {
		Parser parser.IParser
	}
)

func New(parser parser.IParser) *Input {
	return &Input{
		Parser: parser,
	}
}

func (i *Input) Bind(v interface{}) error {
	return i.Parser.Bind(v)
}

//func (i *Input) Input(key string) interface{} {
//	return i.parser.Get(key)
//}

func (i *Input) Has(key string) bool {
	return i.Parser.Has(key)
}
func (i *Input) Get(key string) interface{} {
	return i.Parser.Get(key)
}

func (i *Input) GetString(key string) string {
	return i.Get(key).(string)
}

func (i *Input) GetInt(key string) int {
	return i.Get(key).(int)
}

func (i *Input) GetUInt(key string) uint {
	//v := i.Get(key)
	//if value, ok := v.(string); ok {
	//	intValue, _ := strconv.Atoi(value)
	//	return uint(intValue)
	//}

	return i.Get(key).(uint)
}

func (i *Input) GetFloat(key string) float64 {
	return i.Get(key).(float64)
}

func (i *Input) GetSliceString(key string) []string  {
	return i.Get(key).([]string)
}

func (i *Input) GetSliceInt(key string) []int  {
	return i.Get(key).([]int)
}

func (i *Input) GetSliceFloat(key string) []float64  {
	return i.Get(key).([]float64)
}

func (i *Input) GetMapString(key string) map[string]string  {
	return i.Get(key).(map[string]string)
}

func (i *Input) GetMapInterface(key string) map[string]interface{}  {
	return i.Get(key).(map[string]interface{})
}
