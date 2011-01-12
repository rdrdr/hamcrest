// Copyright 2011 Mick Killianey.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package collections

import (
	"hamcrest"
	"reflect"
)

type _ElemAndLen interface {
	Elem(i int) reflect.Value
	Len() int
}

func AnyElement(matcher *hamcrest.Matcher) *hamcrest.Matcher {
	description := hamcrest.NewDescription("AnyElement[%v]", matcher)
	match := func(actual interface{}) *hamcrest.Result {
		v := reflect.NewValue(actual)
		if value, ok := v.(_ElemAndLen); ok {
			n := value.Len()
			for i := 0; i < n; i++ {
				elem := value.Elem(i).Interface()
				result := matcher.Match(elem)
				if result.Matched() {
					why := hamcrest.NewDescription(
						"Matched element %v of %v: %v", i+1, n, elem)
					return hamcrest.NewResult(true, why)
				}
			}
			why := hamcrest.NewDescription(
				"Matched none of the %v elements", n)
			return hamcrest.NewResult(false, why)
		}
		return matcher.Match(v)
	}
	return hamcrest.NewMatcher(description, match)
}

func EveryElement(matcher *hamcrest.Matcher) *hamcrest.Matcher {
	description := hamcrest.NewDescription("EveryElement[%v]", matcher)
	match := func(actual interface{}) *hamcrest.Result {
		v := reflect.NewValue(actual)
		if value, ok := v.(_ElemAndLen); ok {
			n := value.Len()
			for i := 0; i < n; i++ {
				elem := value.Elem(i).Interface()
				result := matcher.Match(elem)
				if !result.Matched() {
					why := hamcrest.NewDescription(
						"Failed to match element %v of %v: %v", i+1, n, elem)
					return hamcrest.NewResult(false, why)
				}
			}
			why := hamcrest.NewDescription(
				"Matched all of the %v elements", n)
			return hamcrest.NewResult(true, why)
		}
		return matcher.Match(v)
	}
	return hamcrest.NewMatcher(description, match)
}

func AnyMapElement(matcher *hamcrest.Matcher) *hamcrest.Matcher {
	description := hamcrest.NewDescription("AnyMapElement[%v]", matcher)
	match := func(actual interface{}) *hamcrest.Result {
		v := reflect.NewValue(actual)
		if value, ok := v.(*reflect.MapValue); ok {
			keys := value.Keys()
			for i, key := range value.Keys() {
				elem := value.Elem(key).Interface()
				result := matcher.Match(elem)
				if result.Matched() {
					why := hamcrest.NewDescription(
						"Matched element %v of %v: %v", i+1, len(keys), elem)
					return hamcrest.NewResult(true, why)
				}
			}
			why := hamcrest.NewDescription(
				"Matched none of the %v elements", len(keys))
			return hamcrest.NewResult(false, why)
		}
		return matcher.Match(v)
	}
	return hamcrest.NewMatcher(description, match)
}

func EveryMapElement(matcher *hamcrest.Matcher) *hamcrest.Matcher {
	description := hamcrest.NewDescription("EveryMapElement[%v]", matcher)
	match := func(actual interface{}) *hamcrest.Result {
		v := reflect.NewValue(actual)
		if value, ok := v.(*reflect.MapValue); ok {
			keys := value.Keys()
			for _, key := range keys {
				elem := value.Elem(key).Interface()
				result := matcher.Match(elem)
				if !result.Matched() {
					why := hamcrest.NewDescription(
						"Failed to match map element at key[%v]: %v", key, elem)
					return hamcrest.NewResult(false, why)
				}
			}
			why := hamcrest.NewDescription(
				"Matched all of the %v map elements", len(keys))
			return hamcrest.NewResult(true, why)
		}
		return matcher.Match(v)
	}
	return hamcrest.NewMatcher(description, match)
}


