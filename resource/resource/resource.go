package resource

import (
	"reflect"
	"regexp"
	"strings"
	"sync"

	"github.com/firmeve/firmeve/resource"
	reflect2 "github.com/firmeve/firmeve/support/reflect"
	strings2 "github.com/firmeve/firmeve/support/strings"
)

type mapCache map[string]map[string]string

type meta map[string]interface{}

var (
	resourcesFields  = make(map[reflect.Type]mapCache, 0)
	resourcesMethods = make(map[reflect.Type]mapCache, 0)
	mutex            sync.Mutex
)

type ResolveMap map[string]interface{}

type Resolver interface {
	Resolve() ResolveMap
}

type Resource struct {
	source interface{}
	fields []string
	chunks []string
	key    string
	meta   meta
	//transformer interface{}
}

func New(source interface{}) *Resource {
	return &Resource{
		//source:      source,
		source: source,
		key:    `data`,
		meta:   make(meta, 0),
	}
}

func (r *Resource) Fields(fields ...string) *Resource {
	r.fields = fields
	return r
}

func (r *Resource) Chunks(chunks []string) *Resource {
	r.chunks = chunks
	return r
}

func (r *Resource) resolveFields() []string {
	return r.fields
}

func (r *Resource) Resolve() ResolveMap {
	reflectType := reflect.TypeOf(r.source)
	reflectValue := reflect.ValueOf(r.source)
	kindType := reflect2.KindElemType(reflectType)
	var data ResolveMap

	if _, ok := r.source.(resource.Transformer); ok {
		data = r.resolveTransformer(reflectType, reflectValue)
	} else if kindType == reflect.Map {
		data = r.resolveMap(reflectType, reflectValue)
	} else if kindType == reflect.Struct {
		data = r.resolveStruct(reflectType, reflectValue)
	} else {
		panic(`type error`)
	}

	return ResolveMap{r.key: data}
}

func (r *Resource) resolveMap(reflectType reflect.Type, reflectValue reflect.Value) ResolveMap {
	var alias string
	collection := make(ResolveMap, 0)
	for _, field := range r.resolveFields() {
		for k, v := range r.source.(map[string]interface{}) {
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

func (r *Resource) resolveTransformer(reflectType reflect.Type, reflectValue reflect.Value) ResolveMap {
	fields := r.transpositionFields(reflectType)
	methods := r.transpositionMethods(reflectType)
	collection := make(ResolveMap, 0)

	for _, field := range r.resolveFields() {
		// method 优先
		if v, ok := methods[field]; ok {
			collection[v[`alias`]] = reflect2.CallMethodValue(reflectValue, v[`method`])[0] //utils.ReflectValueInterface(utils.ReflectCallMethod(source, v[`method`])[0])
		} else if v, ok := fields[field]; ok {
			if v[`method`] == `` {
				collection[v[`alias`]] = reflect2.CallFieldValue(reflectValue, v[`name`]) // reflect2.InterfaceValue(reflect.Indirect(reflectValue).FieldByName(v[`name`])) //utils.ReflectValueInterface(utils.ReflectCallMethod(source, v[`method`])[0])
			} else {
				collection[v[`alias`]] = reflect2.CallMethodValue(reflectValue, v[`method`])[0]
			}
		} else {
			collection[field] = ``
		}
	}

	return collection
}

func (r *Resource) resolveStruct(reflectType reflect.Type, reflectValue reflect.Value) ResolveMap {
	fields := r.transpositionFields(reflectType)
	collection := make(ResolveMap, 0)

	for _, field := range r.resolveFields() {
		// method 优先
		if v, ok := fields[field]; ok { //reflect.Indirect(reflectValue).FieldByName(v[`name`]).Interface() //
			collection[v[`alias`]] = reflect2.CallFieldValue(reflectValue, v[`name`]) // reflect2.InterfaceValue(reflect.Indirect(reflectValue).FieldByName(v[`name`])) //utils.ReflectValueInterface(utils.ReflectCallMethod(source, v[`method`])[0])
		} else {
			collection[field] = ``
		}
	}
	return collection
}

func (r *Resource) transpositionMethods(reflectType reflect.Type) mapCache {
	if v, ok := resourcesMethods[reflectType]; ok {
		return v
	}

	methods := make(mapCache, 0)

	reflect2.CallMethodType(reflectType, func(i int, method reflect.Method) interface{} {
		name := method.Name
		if regexp.MustCompile("^(.+)Field$").MatchString(name) {
			exceptFieldName := name[0 : len(name)-5]
			methods[exceptFieldName] = map[string]string{`alias`: strings2.SnakeCase(exceptFieldName), `method`: name}
		}
		return nil
	})

	//for name := range reflect2.Methods(reflectType) {
	//	if regexp.MustCompile("^(.+)Field$").MatchString(name) {
	//		exceptFieldName := name[0 : len(name)-5]
	//		methods[exceptFieldName] = map[string]string{`alias`: strings2.SnakeCase(exceptFieldName), `method`: name}
	//	}
	//}

	resourcesMethods[reflectType] = methods

	return methods
}

func (r *Resource) transpositionFields(reflectType reflect.Type) mapCache {
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
	//for name, field := range reflect2.StructFields(reflectType) {
	//
	//}

	resourcesFields[reflectType] = fields

	return fields
}
