// Copyright 2011 Mick Killianey.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package collections

import (
	"github.com/rdrdr/hamcrest/base"
	"reflect"
)

type _ElemAndLen interface {
	Elem(i int) reflect.Value
	Len() int
}

// Returns a matcher that matches on any array or slice input value
// if the given matcher matches at least one element of that array
// or slice.
//
// The returned matcher does not match any non-array-or-slice value.
func AnyElement(matcher *base.Matcher) *base.Matcher {
	match := func(actual interface{}) *base.Result {
		v := reflect.NewValue(actual)
		if value, ok := v.(_ElemAndLen); ok {
			n := value.Len()
			for i := 0; i < n; i++ {
				elem := value.Elem(i).Interface()
				result := matcher.Match(elem)
				if result.Matched() {
					return base.NewResultf(true,
						"Matched element %v of %v: %v", i+1, n, elem).
						WithCauses(result)
				}
			}
			return base.NewResultf(false,
				"Matched none of the %v elements", n)
		}
		return matcher.Match(v)
	}
	return base.NewMatcherf(match, "AnyElement[%v]", matcher)
}

// Returns a matcher that matches on any array or slice input value
// if the given matcher matches every element of that array or slice.
//
// The returned matcher does not match any non-array-or-slice value.
func EveryElement(matcher *base.Matcher) *base.Matcher {
	match := func(actual interface{}) *base.Result {
		v := reflect.NewValue(actual)
		var value _ElemAndLen
		var ok bool
		value, ok = v.(*reflect.ArrayValue)
		if !ok {
			value, ok = v.(*reflect.SliceValue)
		}
		if !ok {
			return base.NewResultf(false,
				"Was not array or slice: was type %T", actual)
		}
		n := value.Len()
		for i := 0; i < n; i++ {
			elem := value.Elem(i).Interface()
			result := matcher.Match(elem)
			if !result.Matched() {
				return base.NewResultf(false,
					"Failed to match element %v of %v: %v",
					i+1, n, elem).
					WithCauses(result)
			}
		}
		return base.NewResultf(true,
			"Matched all of the %v elements", n)
	}
	return base.NewMatcherf(match, "EveryElement[%v]", matcher)
}

func AnyMapElement(matcher *base.Matcher) *base.Matcher {
	match := func(actual interface{}) *base.Result {
		v := reflect.NewValue(actual)
		value, ok := v.(*reflect.MapValue)
		if !ok {
			return base.NewResultf(false,
				"Was not map: was type %T", actual)
		}
		keys := value.Keys()
		for i, keyValue := range value.Keys() {
			elem := value.Elem(keyValue).Interface()
			result := matcher.Match(elem)
			if result.Matched() {
				return base.NewResultf(true,
					"Matched map element [%v/%v] with key [%v]: %v",
					i+1, len(keys), keyValue.Interface(), elem).
					WithCauses(result)
			}
		}
		return base.NewResultf(false,
			"Matched none of the %v elements", len(keys))
	}
	return base.NewMatcherf(match, "AnyMapElement[%v]", matcher)
}

func EveryMapElement(matcher *base.Matcher) *base.Matcher {
	match := func(actual interface{}) *base.Result {
		v := reflect.NewValue(actual)
		value, ok := v.(*reflect.MapValue)
		if !ok {
			return base.NewResultf(false,
				"Was not map: was type %T", actual)
		}
		keys := value.Keys()
		for i, keyValue := range keys {
			elem := value.Elem(keyValue).Interface()
			result := matcher.Match(elem)
			if !result.Matched() {
				return base.NewResultf(false,
					"Failed to match map element [%v/%v] with key[%v]: %v",
					i+1, len(keys), keyValue.Interface(), elem).
					WithCauses(result)
			}
		}
		return base.NewResultf(true,
			"Matched all of the %v map elements", len(keys))
	}
	return base.NewMatcherf(match, "EveryMapElement[%v]", matcher)
}

type _HasLen interface { Len() int }

// Applies the given matcher to the length of the input element,
// if the input element is an array, slice, or map.
func ToLen(matcher *base.Matcher) *base.Matcher {
	match := func(actual interface{}) *base.Result {
		value := reflect.NewValue(actual)
		if hasLen, ok := value.(_HasLen); ok {
			length := hasLen.Len()
			result := matcher.Match(length)
			return base.NewResultf(result.Matched(), "Len() returned %v", length)
		}
		return base.NewResultf(false,
			"Can't determine Len() for %T", actual)
	}
	return base.NewMatcherf(match, "ToLen[%v]", matcher)
}


// Matches any input element that is an empty array, slice, or map.
func Empty() *base.Matcher {
	match := func(actual interface{}) *base.Result {
		value := reflect.NewValue(actual)
		if hasLen, ok := value.(_HasLen); ok {
			length := hasLen.Len()
			return base.NewResultf(length == 0,
				"Len() returned %v", length)
		}
		return base.NewResultf(false, "Can't determine length of type %T", actual)
	}
	return base.NewMatcherf(match, "Empty")
}

