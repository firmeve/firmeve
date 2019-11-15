package resource

import (
	"reflect"
	"regexp"
	"strings"

	"github.com/firmeve/firmeve/converter/transform"

	"github.com/firmeve/firmeve/support"

	reflect2 "github.com/firmeve/firmeve/support/reflect"
	strings2 "github.com/firmeve/firmeve/support/strings"
)

type mapCache map[string]map[string]string

type baseResource struct {
	//source   Resolver
	resource interface{}
	fields   []string
	//key      string
	meta Meta
	link Link
}

var (
	resourcesFields  = make(map[reflect.Type]mapCache, 0)
	resourcesMethods = make(map[reflect.Type]mapCache, 0)
	//mutex            sync.Mutex
)

//type option struct {
//	transformer transform.Transformer
//}

//func WithTransformer(t transform.Transformer) support.Option {
//	return func(object support.Object) {
//		object.(*option).transformer = t
//	}
//}

func newBaseResource(resource interface{}, options ...support.Option) *baseResource {
	return &baseResource{
		//source:   source,
		resource: resource,
		//key:      `data`,
		meta: make(Meta, 0),
		link: make(Link, 0),
	}
}

//func (r *baseResource) Fields(fields ...string) *baseResource {
//	r.fields = fields
//	return r
//}
//
//func (r *baseResource) Meta(meta Meta) *baseResource {
//	r.meta = meta
//	return r
//}
//
//func (r *baseResource) Key(key string) *baseResource {
//	r.key = key
//	return r
//}

//func (r *baseResource) SetMeta(meta Meta) IMeta {
//	r.meta = meta
//	return r
//}
//func (r *baseResource) Meta() Meta {
//	return r.meta
//}

func (r *baseResource) resolve() Data {
	reflectType := reflect.TypeOf(r.resource)
	reflectValue := reflect.ValueOf(r.resource)
	var data Data

	if _, ok := r.resource.(transform.Transformer); ok {
		data = r.resolveTransformer(reflectType, reflectValue)
	} else {
		kindType := reflect2.KindElemType(reflectType)
		if kindType == reflect.Map {
			data = r.resolveMap(reflectType, reflectValue)
		} else if kindType == reflect.Struct {
			data = r.resolveStruct(reflectType, reflectValue)
		} else {
			panic(`type error`)
		}
	}

	return data
}

func (r *baseResource) resolveFields() []string {
	return r.fields
}

func (r *baseResource) resolveMap(reflectType reflect.Type, reflectValue reflect.Value) Data {
	var alias string
	collection := make(Data, 0)
	for _, field := range r.resolveFields() {
		for k, v := range r.resource.(map[string]interface{}) {
			alias = strings2.SnakeCase(k)
			if field != alias {
				continue
			}
			if reflect2.KindElemType(reflect.TypeOf(v)) == reflect.Func {
				collection[alias] = reflect2.CallFuncValue(reflect.ValueOf(v))[0]
			} else {
				collection[alias] = v
			}
		}
	}
	return collection
}

func (r *baseResource) resolveTransformer(reflectType reflect.Type, reflectValue reflect.Value) Data {
	resource := r.resource.(transform.Transformer).Resource()
	resourceReflectType := reflect.TypeOf(resource)
	resourceReflectValue := reflect.ValueOf(resource)
	fields := r.transpositionFields(resourceReflectType)
	methods := r.transpositionMethods(reflectType)
	collection := make(Data, 0)

	for _, field := range r.resolveFields() {
		// method 优先
		if v, ok := methods[field]; ok {
			collection[v[`alias`]] = reflect2.CallMethodValue(reflectValue, v[`method`])[0]
		} else if v, ok := fields[field]; ok {
			if v[`method`] == `` {
				collection[v[`alias`]] = reflect2.CallFieldValue(resourceReflectValue, v[`name`])
			} else {
				collection[v[`alias`]] = reflect2.CallMethodValue(reflectValue, v[`method`])[0]
			}
		} else {
			collection[field] = ``
		}
	}

	return collection
}

func (r *baseResource) resolveStruct(reflectType reflect.Type, reflectValue reflect.Value) Data {
	fields := r.transpositionFields(reflectType)
	collection := make(Data, 0)

	for _, field := range r.resolveFields() {
		// method 优先
		if v, ok := fields[field]; ok {
			collection[v[`alias`]] = reflect2.CallFieldValue(reflectValue, v[`name`])
		} else {
			collection[field] = ``
		}
	}
	return collection
}

func (r *baseResource) transpositionMethods(reflectType reflect.Type) mapCache {
	if v, ok := resourcesMethods[reflectType]; ok {
		return v
	}

	methods := make(mapCache, 0)

	reflect2.CallMethodType(reflectType, func(i int, method reflect.Method) interface{} {
		name := method.Name
		if regexp.MustCompile("^(.+)Field$").MatchString(name) {
			alias := strings2.SnakeCase(name[0 : len(name)-5])
			methods[alias] = map[string]string{`alias`: alias, `method`: name}
		}
		return nil
	})

	resourcesMethods[reflectType] = methods

	return methods
}

func (r *baseResource) transpositionFields(reflectType reflect.Type) mapCache {
	if v, ok := resourcesFields[reflectType]; ok {
		return v
	}

	fields := make(mapCache, 0)
	var alias, method string

	reflect2.CallFieldType(reflectType, func(i int, field reflect.StructField) interface{} {
		method = ``
		name := field.Name
		if field.Tag.Get(`resource`) != `` {
			tagNames := strings.Split(field.Tag.Get(`resource`), `,`)
			alias = tagNames[0]
			if len(tagNames) >= 2 {
				method = tagNames[1]
			}
		} else { //method
			alias = strings2.SnakeCase(name)
		}

		fields[alias] = map[string]string{`alias`: alias, `method`: method, `name`: name}

		return nil
	})

	resourcesFields[reflectType] = fields

	return fields
}
