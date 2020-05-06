package resource

import (
	"github.com/firmeve/firmeve/kernel/contract"
	"reflect"
	"regexp"
	"strings"

	reflect2 "github.com/firmeve/firmeve/support/reflect"
	strings2 "github.com/firmeve/firmeve/support/strings"
)

type (
	Item struct {
		resource    interface{}
		resolveData contract.ResourceData
		option      *Option
		meta        contract.ResourceMetaData
		link        contract.ResourceLinkData
	}

	mapCache map[string]map[string]string
)

var (
	resourcesFields  = make(map[reflect.Type]mapCache, 0)
	resourcesMethods = make(map[reflect.Type]mapCache, 0)
)

func NewItem(resource interface{}, option *Option) *Item {
	return &Item{
		resource:    resolveResource(resource, option),
		resolveData: make(contract.ResourceData, 0),
		option:      option,
	}
}

func resolveResource(resource interface{}, option *Option) interface{} {
	if option.Transformer != nil {
		option.Transformer.SetResource(resource)
		resource = option.Transformer
	}
	return resource
}

func (i *Item) SetOption(option *Option) *Item {
	i.option = option
	return i
}

func (i *Item) SetMeta(meta contract.ResourceMetaData) {
	i.meta = meta
}
func (i *Item) Meta() contract.ResourceMetaData {
	return i.meta
}

func (i *Item) SetLink(link contract.ResourceLinkData) {
	i.link = link
}

func (i *Item) Link() contract.ResourceLinkData {
	return i.link
}

func (i *Item) resolveFields() []string {
	return i.option.Fields
}

func (i *Item) Data() contract.ResourceData {
	if len(i.resolveData) > 0 {
		return i.resolveData
	}

	i.resolveData = i.resolve()

	return i.resolveData
}

func (i *Item) resolve() contract.ResourceData {
	reflectType := reflect.TypeOf(i.resource)
	reflectValue := reflect.ValueOf(i.resource)
	var data contract.ResourceData

	if _, ok := i.resource.(contract.ResourceTransformer); ok {
		data = i.resolveTransformer(reflectType, reflectValue)
	} else {
		kindType := reflect2.KindElemType(reflectType)
		if kindType == reflect.Map {
			data = i.resolveMap()
		} else if kindType == reflect.Struct {
			data = i.resolveStruct(reflectType, reflectValue)
		} else {
			panic(`type error`)
		}
	}

	return data
}

func (i *Item) resolveMap() contract.ResourceData {
	var alias string
	collection := make(contract.ResourceData, 0)
	for _, field := range i.resolveFields() {
		for k, v := range i.resource.(contract.ResourceData) {
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

func (i *Item) resolveTransformer(reflectType reflect.Type, reflectValue reflect.Value) contract.ResourceData {
	resource := i.resource.(contract.ResourceTransformer).Resource()
	resourceReflectType := reflect.TypeOf(resource)
	resourceReflectValue := reflect.ValueOf(resource)
	fields := i.transpositionFields(resourceReflectType)
	methods := i.transpositionMethods(reflectType)
	collection := make(contract.ResourceData, 0)

	for _, field := range i.resolveFields() {
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

func (i *Item) resolveStruct(reflectType reflect.Type, reflectValue reflect.Value) contract.ResourceData {
	fields := i.transpositionFields(reflectType)
	collection := make(contract.ResourceData, 0)

	for _, field := range i.resolveFields() {
		// method 优先
		if v, ok := fields[field]; ok {
			collection[v[`alias`]] = reflect2.CallFieldValue(reflectValue, v[`name`])
		} else {
			collection[field] = ``
		}
	}
	return collection
}

func (i *Item) transpositionMethods(reflectType reflect.Type) mapCache {
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

func (i *Item) transpositionFields(reflectType reflect.Type) mapCache {
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
