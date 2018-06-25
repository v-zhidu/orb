package set

import (
	"sort"
	"sync"
)

//StringSet store unique string values
type StringSet struct {
	sync.RWMutex
	m map[string]bool
}

//New returns a new StringSet struct with items
func New(items ...string) *StringSet {
	s := &StringSet{
		m: make(map[string]bool, len(items)),
	}
	s.Add(items...)

	return s
}

//Add new values to StringSet
func (s *StringSet) Add(items ...string) {
	s.Lock()
	defer s.Unlock()
	for _, v := range items {
		s.m[v] = true
	}
}

//Remove values from StringSet
func (s *StringSet) Remove(items ...string) {
	s.Lock()
	defer s.Unlock()
	for _, v := range items {
		delete(s.m, v)
	}
}

//Contains check whether if StringSet contains values
func (s *StringSet) Contains(items ...string) bool {
	s.RLock()
	defer s.RUnlock()
	for _, v := range items {
		if _, ok := s.m[v]; !ok {
			return false
		}
	}

	return true
}

//Size returns the length of StringSet
func (s *StringSet) Size() int {
	return len(s.m)
}

//Empty returns if the StringSet is empty
func (s *StringSet) Empty() bool {
	return s.Size() == 0
}

//Clear returns an empty StringSet
func (s *StringSet) Clear() {
	s.Lock()
	defer s.Unlock()
	s.m = map[string]bool{}
}

//StringSlice convert StringSet to Slice
func (s *StringSet) StringSlice() []string {
	s.RLock()
	defer s.RUnlock()
	slice := []string{}
	for item := range s.m {
		slice = append(slice, item)
	}

	return slice
}

//SortStringSlice convert StringSet to sorted string slice
func (s *StringSet) SortStringSlice() []string {
	slice := s.StringSlice()
	return sort.StringSlice(slice)
}

//Union 并集
func (s *StringSet) Union(sets ...*StringSet) *StringSet {
	r := New(s.StringSlice()...)
	for _, set := range sets {
		r.Add(set.StringSlice()...)
	}

	return r
}

//Minus 差集
func (s *StringSet) Minus(sets ...*StringSet) *StringSet {
	r := New(s.StringSlice()...)
	for _, set := range sets {
		for item := range set.m {
			if _, ok := s.m[item]; ok {
				delete(r.m, item)
			}
		}
	}

	return r
}

//Intersect 交集
func (s *StringSet) Intersect(sets ...*StringSet) *StringSet {
	r := &StringSet{
		m: map[string]bool{},
	}
	for _, set := range sets {
		for item := range set.m {
			if _, ok := s.m[item]; ok {
				r.m[item] = true
			}
		}
	}

	return r
}
