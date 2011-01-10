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
	description := NewDescription("always matches")
	match := func (actual interface{}) *Result {
		return NewResult(true, description)
	}
	_Anything = NewMatcher(NewDescription("matches anything"), match)
}


// Returns a Matcher that matches the boolean value true.
func True() *Matcher {
	return _True
}
var _True *Matcher // singleton
func init() {
	trueDescription := NewDescription("was true")
	falseDescription := NewDescription("was not true")
	match := func (actual interface{}) *Result {
		if b, ok := actual.(bool); ok {
			if b {
				return NewResult(true, trueDescription)
			}
			return NewResult(false, falseDescription)
		}
		return NewResult(false, NewDescription("[%v] was not bool", actual))
	}
	_True = NewMatcher(NewDescription("matches true"), match)
}

// Returns a Matcher that matches the boolean value false.
func False() *Matcher {
	return _False
}
var _False *Matcher // singleton
func init() {
	trueDescription := NewDescription("was not false")
	falseDescription := NewDescription("was false")
	match := func (actual interface{}) *Result {
		if b, ok := actual.(bool); ok {
			if b {
				return NewResult(false, trueDescription)
			}
			return NewResult(true, falseDescription)
		}
		return NewResult(false, NewDescription("[%v] was not bool", actual))
	}
	_False = NewMatcher(NewDescription("matches false"), match)
}


// Returns a Matcher that matches the boolean value true.
func Matched() *Matcher {
	return _Matched
}
var _Matched *Matcher // singleton
func init() {
	description := NewDescription("was a result")
	match := func (actual interface{}) *Result {
		if result, ok := actual.(*Result); ok {
			return NewResult(result.Matched(), description).WithCauses(result)
		}
		return NewResult(false, NewDescription("[%v] was not a result", actual))
	}
	_Matched = NewMatcher(NewDescription("Matched"), match)
}

// Returns a Matcher that matches the boolean value false.
func DidNotMatch() *Matcher {
	return _DidNotMatch
}
var _DidNotMatch *Matcher // singleton
func init() {
	description := NewDescription("was a result")
	match := func (actual interface{}) *Result {
		if result, ok := actual.(*Result); ok {
			return NewResult(!result.Matched(), description).WithCauses(result)
		}
		return NewResult(false, NewDescription("[%v] was not a result", actual))
	}
	_DidNotMatch = NewMatcher(NewDescription("DidNotMatch"), match)
}


// Returns a Matcher that decorates another matcher and only matches
// when the underlying matcher does not match (and vice versa).
func Not(matcher *Matcher) *Matcher {
	match := func (actual interface{}) *Result {
		result := matcher.Match(actual)
		return NewResult(!result.matched, result.description).WithCauses(result)
	}
	return NewMatcher(NewDescription("Not[%v]", matcher), match)
}

// Returns a Matcher that decorates another matcher.
func Is(matcher *Matcher) *Matcher {
	match := func (actual interface{}) *Result {
		result := matcher.Match(actual)
		description := NewDescription("Is[%v]", result.description)
		return NewResult(result.matched, description).WithCauses(result)
	}
	return NewMatcher(NewDescription("Is[%v]", matcher), match)
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
	wasNilDescription := NewDescription("was nil")
	match := func (actual interface{}) *Result {
		if _detectNil(actual) {
			return NewResult(true, wasNilDescription)
		}
		return NewResult(false, NewDescription("[%v] was not nil", actual))
	}
	_Nil = NewMatcher(NewDescription("matches nil"), match)
}

// Returns a Matcher that matches if the actual value is 
// neither nil nor the nil value of its type.  (Note that
// this is *not* equivalent to Not(DeeplyEqualTo(nil)).)
func NonNil() *Matcher {
	return _NonNil
}
var _NonNil *Matcher
func init() {
	wasNilDescription := NewDescription("was nil")
	match := func (actual interface{}) *Result {
		if _detectNil(actual) {
			return NewResult(false, wasNilDescription)
		}
		return NewResult(true, NewDescription("[%v] was not nil", actual))
	}
	_NonNil = NewMatcher(NewDescription("matches non-nil"), match)
}


// Returns a Matcher that checks if the actual value is (deeply)
// equal to the given expected value, using reflect.DeepEqual.
func DeeplyEqualTo(expected interface{}) *Matcher {
	match := func (actual interface{}) *Result {
		if reflect.DeepEqual(expected, actual) {
			return NewResult(true,
				NewDescription("was deeply equal to [%v]", expected))
		}
		return NewResult(false,
			NewDescription("[%v] was not deeply equal to [%v]", actual, expected))
	}
	return NewMatcher(NewDescription("DeeplyEqualTo[%v]", expected), match)
}


// Returns a short-circuiting Matcher that matches whenever all of
// the given matchers match a given input value.  If any component
// matcher fails to match an input value, later matchers are not
// attempted.
func AllOf(matchers...*Matcher) *Matcher {
	descriptions := make([]*Description, len(matchers), len(matchers))
	for index, matcher := range matchers {
		descriptions[index] = NewDescription("[%v: %v]", index+1, matcher)
	}
	description := NewDescription("AllOf%v", descriptions)
	match := func (actual interface{}) *Result {
		var results []*Result
		for index, matcher := range matchers {
			result := matcher.Match(actual)
			results := append(results, result)
			if !result.Matched() {
				because := NewDescription(
					"Failed matcher %v of %v: [%v]",
						index+1, len(matchers), matcher)
				return NewResult(false, because).WithCauses(results...)
			}
		}
		because := NewDescription("Matched all %v matchers", len(matchers))
		return NewResult(true, because).WithCauses(results...)
	}
	return NewMatcher(description, match)
}

// Returns a short-circuiting Matcher that matches whenever all of
// the given matchers match a given input value.  If any component
// matcher fails to match an input value, later matchers are not
// attempted.
func AnyOf(matchers...*Matcher) *Matcher {
	descriptions := make([]*Description, len(matchers), len(matchers))
	for index, matcher := range matchers {
		descriptions[index] = NewDescription("[%v: %v]", index+1, matcher)
	}
	description := NewDescription("AnyOf%v", descriptions)
	match := func (actual interface{}) *Result {
		var results []*Result
		for index, matcher := range matchers {
			result := matcher.Match(actual)
			results := append(results, result)
			if result.Matched() {
				because := NewDescription(
					"Matched on matcher %v of %v: [%v]", index+1, len(matchers), matcher)
				return NewResult(true, because).WithCauses(results...)
			}
		}
		because := NewDescription("Matched none of the %v matchers", len(matchers))
		return NewResult(false, because).WithCauses(results...)
	}
	return NewMatcher(description, match)
}

