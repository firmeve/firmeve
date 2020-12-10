package slices

import (
	"fmt"
	"sync"
)

type Slice struct {
	s             []interface{}
	values        map[interface{}][]int
	length        int
	defaultLength int
	lock          sync.Mutex
}

func New() *Slice {
	return &Slice{
		s:      make([]interface{}, 0),
		values: make(map[interface{}][]int, 0),
		length: 0,
		//defaultLength: 10,
	}
}

func NewWith(s []interface{}) *Slice {
	values := make(map[interface{}][]int, 0)
	for i := range s {
		values[s[i]] = append(values[s[i]], i)
	}
	return &Slice{
		s:             s,
		values:        values,
		length:        len(s),
		defaultLength: 0,
	}
}

func (s *Slice) Add(value interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.s = append(s.s, value)
	s.values[value] = append(s.values[value], len(s.s)-1)
	s.length += 1
}

func (s *Slice) DeleteWithValue(value interface{}) {
	s2 := s.values[value]

	newSlice := make([]interface{}, 0)
	start := 0
	end := 0
	for i := range s2 {
		fmt.Println(s2[i])
		end = s2[i]
		newSlice = append(newSlice, s.s[start:end]...)
		start = end
	}
	s.s = append(newSlice, s.s[end:]...)
}

func (s *Slice) Values() []interface{} {
	return s.s
}

func (s *Slice) Len() int {
	return s.length
}

func (s *Slice) Exists(value interface{}) bool {
	for k := range s.values {
		if k == value {
			return true
		}
	}
	return false
}

func (s *Slice) Index(value interface{}) int {
	for k := range s.values {
		if k == value {
			return s.values[k][0]
		}
	}
	return -1
}

func (s *Slice) LastIndex(value interface{}) int {
	for k := range s.values {
		if k == value {
			return s.values[k][len(s.values[k])-1]
		}
	}
	return -1
}
