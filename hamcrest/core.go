// Copyright 2011 Mick Killianey.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hamcrest

import (
	"reflect"
)

// Returns a Matcher that matches any input value.
func Anything() *Matcher {
	return _Anything
}
var _Anything *Matcher // singleton
func init() {
	match := func (actual interface{}) *Result {
		return NewResultf(true, "always matches")
	}
	_Anything = NewMatcherf(match, "Anything")
}


// Returns a Matcher that matches the boolean value true.
func True() *Matcher {
	return _True
}
var _True *Matcher // singleton
func init() {
	match := func (actual interface{}) *Result {
		if b, ok := actual.(bool); ok {
			if b {
				return NewResultf(true, "was true")
			}
			return NewResultf(false, "was not true")
		}
		return NewResultf(false, "[%v] was not bool", actual)
	}
	_True = NewMatcherf(match, "True")
}

// Returns a Matcher that matches the boolean value false.
func False() *Matcher {
	return _False
}
var _False *Matcher // singleton
func init() {
	match := func (actual interface{}) *Result {
		if b, ok := actual.(bool); ok {
			if !b {
				return NewResultf(true, "was false")
			}
			return NewResultf(false, "was not false")
		}
		return NewResultf(false, "[%v] was not bool", actual)
	}
	_False = NewMatcherf(match, "False")
}


// Returns a Matcher that matches the boolean value true.
func Matched() *Matcher {
	return _Matched
}
var _Matched *Matcher // singleton
func init() {
	match := func (actual interface{}) *Result {
		if result, ok := actual.(*Result); ok {
			return NewResultf(result.Matched(), "was a result").WithCauses(result)
		}
		return NewResultf(false, "[%v] was not a result", actual)
	}
	_Matched = NewMatcherf(match, "Matched")
}

// Returns a Matcher that matches the boolean value false.
func DidNotMatch() *Matcher {
	return _DidNotMatch
}
var _DidNotMatch *Matcher // singleton
func init() {
	match := func (actual interface{}) *Result {
		if result, ok := actual.(*Result); ok {
			return NewResultf(!result.Matched(), "was a result").WithCauses(result)
		}
		return NewResultf(false, "[%v] was not a result", actual)
	}
	_DidNotMatch = NewMatcherf(match, "DidNotMatch")
}


// Returns a Matcher that decorates another matcher and only matches
// when the underlying matcher does not match (and vice versa).
func Not(matcher *Matcher) *Matcher {
	match := func (actual interface{}) *Result {
		result := matcher.Match(actual)
		return NewResult(!result.matched, result.description).WithCauses(result)
	}
	return NewMatcherf(match, "Not[%v]", matcher)
}

// Returns a Matcher that decorates another matcher.
func Is(matcher *Matcher) *Matcher {
	match := func (actual interface{}) *Result {
		result := matcher.Match(actual)
		return NewResult(result.matched, result.description).WithCauses(result.Causes()...)
	}
	return NewMatcherf(match, "Is[%v]", matcher)
}


// Helper function for Nil/NonNil
func _detectNil(actual interface{}) bool {
	if actual == nil {
		return true
	}
	if value, ok := reflect.NewValue(actual).(_CanAskIsNil); ok {
		return value.IsNil()
	}
	return false
}
type _CanAskIsNil interface { IsNil() bool }


// Returns a Matcher that matches if the actual value is nil
// or the nil value of its type.  (Note that this is *not*
// equivalent to DeeplyEqualTo(nil).)
func Nil() *Matcher {
	return _Nil
}
var _Nil *Matcher // singleton
func init() {
	match := func (actual interface{}) *Result {
		if _detectNil(actual) {
			return NewResultf(true, "was nil")
		}
		return NewResultf(false, "[%v] was not nil", actual)
	}
	_Nil = NewMatcherf(match, "matches nil")
}

// Returns a Matcher that matches if the actual value is 
// neither nil nor the nil value of its type.  (Note that
// this is *not* equivalent to Not(DeeplyEqualTo(nil)).)
func NonNil() *Matcher {
	return _NonNil
}
var _NonNil *Matcher
func init() {
	match := func (actual interface{}) *Result {
		if _detectNil(actual) {
			return NewResultf(false, "was nil")
		}
		return NewResultf(true, "[%v] was not nil", actual)
	}
	_NonNil = NewMatcherf(match, "matches non-nil")
}


// Returns a Matcher that checks if the actual value is (deeply)
// equal to the given expected value, using reflect.DeepEqual.
func DeeplyEqualTo(expected interface{}) *Matcher {
	match := func (actual interface{}) *Result {
		if reflect.DeepEqual(expected, actual) {
			return NewResultf(true,
				"was deeply equal to [%v]", expected)
		}
		return NewResultf(false,
			"[%v] was not deeply equal to [%v]", actual, expected)
	}
	return NewMatcherf(match, "DeeplyEqualTo[%v]", expected)
}


// Returns a short-circuiting Matcher that matches whenever all of
// the given matchers match a given input value.  If any component
// matcher fails to match an input value, later matchers are not
// attempted.
func AllOf(matchers...*Matcher) *Matcher {
	match := func (actual interface{}) *Result {
		var results []*Result
		for index, matcher := range matchers {
			result := matcher.Match(actual)
			results := append(results, result)
			if !result.Matched() {
				return NewResultf(false,
					"Failed matcher %v of %v: [%v]",
					index+1, len(matchers), matcher).
					WithCauses(results...)
			}
		}
		return NewResultf(true,
			"Matched all %v matchers", len(matchers)).
			WithCauses(results...)
	}
	descriptions := make([]SelfDescribing, len(matchers), len(matchers))
	for index, matcher := range matchers {
		descriptions[index] = Description("[#%v: %v]", index+1, matcher)
	}
	return NewMatcherf(match, "AllOf%v", descriptions)
}

// Returns a short-circuiting Matcher that matches whenever all of
// the given matchers match a given input value.  If any component
// matcher fails to match an input value, later matchers are not
// attempted.
func AnyOf(matchers...*Matcher) *Matcher {
	match := func (actual interface{}) *Result {
		var results []*Result
		for index, matcher := range matchers {
			result := matcher.Match(actual)
			results := append(results, result)
			if result.Matched() {
				return NewResultf(true,
					"Matched on matcher %v of %v: [%v]",
					index+1, len(matchers), matcher).
					WithCauses(results...)
			}
		}
		return NewResultf(false,
			"Matched none of the %v matchers", len(matchers)).
			WithCauses(results...)
	}
	descriptions := make([]SelfDescribing, len(matchers), len(matchers))
	for index, matcher := range matchers {
		descriptions[index] = Description("[#%v: %v]", index+1, matcher)
	}
	return NewMatcherf(match, "AnyOf%v", descriptions)
}

