package parser

import "encoding/json"

type (
	JSONData map[string]interface{}

	JSON struct {
		original []byte
		data     JSONData
	}
)

func (j *JSON) Bind(object interface{}) error {
	return json.Unmarshal(j.original, object)
}

func (j *JSON) Has(key string) bool {
	_, ok := j.data[key]
	return ok
}

func (j *JSON) Get(key string) interface{} {
	if j.Has(key) {
		return j.data[key]
	}

	return ``
}

func (j *JSON) parse() {
	if err := json.Unmarshal(j.original, &j.data); err != nil {
		panic(err)
	}
}

func NewJSON(data []byte) *JSON {
	j := &JSON{
		original: data,
		data:     make(JSONData, 0),
	}
	j.parse()
	return j
}
