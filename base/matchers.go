// Copyright 2011 Mick Killianey.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package base

import (
	"reflect"
)


// --------------------------------------------------------------------
// Matchers
// --------------------------------------------------------------------

// Returns a Matcher that matches the boolean value true.
func Matched() *Matcher {
	return _Matched
}
var _Matched *Matcher // singleton
func init() {
	match := func (actual interface{}) *Result {
		if result, ok := actual.(*Result); ok {
			if result.Matched() {
				return NewResultf(true, "was a matching result").WithCauses(result)
			}
			return NewResultf(false, "was a result that did not match").WithCauses(result)
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
			if result.Matched() {
				return NewResultf(false, "was a matching result").WithCauses(result)
			}
			return NewResultf(true, "was a result that did not match").WithCauses(result)
		}
		return NewResultf(false, "[%v] was not a result", actual)
	}
	_DidNotMatch = NewMatcherf(match, "DidNotMatch")
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
	_Nil = NewMatcherf(match, "Nil")
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
	_NonNil = NewMatcherf(match, "NonNil")
}

// Returns a Matcher that checks if the actual value is (deeply)
// equal to the given expected value, using `reflect.DeepEqual`.
//
// For an equality test equivalent to `==`, see the
// `hamcrest/comparison` package.
func DeepEqualTo(expected interface{}) *Matcher {
	match := func (actual interface{}) *Result {
		if reflect.DeepEqual(expected, actual) {
			return NewResultf(true,
				"was deeply equal to [%v]", expected)
		}
		return NewResultf(false,
			"[%v] was not deeply equal to [%v]", actual, expected)
	}
	return NewMatcherf(match, "DeepEqualTo[%v]", expected)
}



