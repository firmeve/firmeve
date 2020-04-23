package slices

import (
	"github.com/firmeve/firmeve/support/reflect"
	reflect2 "reflect"
)

func UniqueString(us []string) []string {
	newUs := make([]string, 0)
	for _, v := range UniqueInterface(us) {
		newUs = append(newUs, v.(string))
	}

	return newUs
}

func UniqueInt(us []int) []int {
	newUs := make([]int, 0)
	for _, v := range UniqueInterface(us) {
		newUs = append(newUs, v.(int))
	}

	return newUs
}

func UniqueInterface(v interface{}) []interface{} {
	us := reflect.SliceInterface(reflect2.ValueOf(v))
	keys := make(map[interface{}]int, 0)
	newUs := make([]interface{}, 0)
	for _, v := range us {
		if _, ok := keys[v]; !ok {
			newUs = append(newUs, v)
		}
		keys[v] = 1
	}

	return newUs
}

func InString(iss []string, v string) bool {
	for i := range iss {
		if v == iss[i] {
			return true
		}
	}

	return false
}
