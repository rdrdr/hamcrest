// Copyright 2011 Mick Killianey.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import (
	"fmt"
	"hamcrest"
	"reflect"
)

// Returns a Matcher that matches any input value.
func Anything() *hamcrest.Matcher {
	return _Anything
}
var _Anything *hamcrest.Matcher // singleton
func init() {
	match := func (actual interface{}) *hamcrest.Result {
		return hamcrest.NewResultf(true, "always matches")
	}
	_Anything = hamcrest.NewMatcherf(match, "Anything")
}


// Returns a Matcher that matches the boolean value true.
func True() *hamcrest.Matcher {
	return _True
}
var _True *hamcrest.Matcher // singleton
func init() {
	match := func (actual interface{}) *hamcrest.Result {
		if b, ok := actual.(bool); ok {
			if b {
				return hamcrest.NewResultf(true, "was true")
			}
			return hamcrest.NewResultf(false, "was not true")
		}
		return hamcrest.NewResultf(false, "[%v] was not bool", actual)
	}
	_True = hamcrest.NewMatcherf(match, "True")
}

// Returns a Matcher that matches the boolean value false.
func False() *hamcrest.Matcher {
	return _False
}
var _False *hamcrest.Matcher // singleton
func init() {
	match := func (actual interface{}) *hamcrest.Result {
		if b, ok := actual.(bool); ok {
			if !b {
				return hamcrest.NewResultf(true, "was false")
			}
			return hamcrest.NewResultf(false, "was not false")
		}
		return hamcrest.NewResultf(false, "[%v] was not bool", actual)
	}
	_False = hamcrest.NewMatcherf(match, "False")
}


// Returns a Matcher that matches on values that cause the given
// functionOrMatcher to panic.
//
// functionOrMatcher should either be a function that accepts one
// parameter or a Matcher.
func PanicWhenApplying(functionOrMatcher interface{}, name string) *hamcrest.Matcher {
	var doSomething func(interface{})
	if matcher, ok := functionOrMatcher.(*hamcrest.Matcher); ok {
		doSomething = func(actual interface{}) { matcher.Match(actual) }
	} else {
		value := reflect.NewValue(functionOrMatcher)
		if funcValue, ok := value.(*reflect.FuncValue); ok {
			funcType := funcValue.Type().(*reflect.FuncType)
			numIn := funcType.NumIn()
			if numIn == 0 {
				panic(fmt.Sprintf("func must accept a single arg, was %T", functionOrMatcher))
			}
			inType := funcType.In(0)
			switch {
			case numIn == 1: // always ok
			case numIn == 2 && funcType.DotDotDot(): // ok
			default:
				panic(fmt.Sprintf("func must accept a single arg, was %T", functionOrMatcher))
			}
			doSomething = func(actual interface{}) {
				actualValue := reflect.NewValue(actual)
				argValues := make([]reflect.Value, numIn, numIn)
				if numIn == 1 && funcType.DotDotDot() {
					inSlice := reflect.MakeSlice(inType.(*reflect.SliceType), 1, 1)
					inValue := inSlice.Elem(0)
					inValue.SetValue(actualValue)
					argValues[0] = inSlice
				} else {
					inValue := reflect.MakeZero(inType)
					inValue.SetValue(actualValue)
					argValues[0] = inValue
					if numIn == 2 && funcType.DotDotDot() {
						inType2 := funcType.In(1)
						argValues[1] = reflect.MakeSlice(inType2.(*reflect.SliceType), 0, 0)
					}
				}
				funcValue.Call(argValues)
			}
		}
	}
	match := func (actual interface{}) (result *hamcrest.Result) {
		defer func() {
			if recover() != nil {
				result = hamcrest.NewResultf(true, "Panicked")
			}
		}()
		doSomething(actual)
		result = hamcrest.NewResultf(false, "Did not panic")
		return
	}
	return hamcrest.NewMatcherf(match,
		"PanicWhenApplying[%v]", name)
}


// Returns a Matcher that decorates another matcher and only matches
// when the underlying matcher does not match (and vice versa).
func Not(matcher *hamcrest.Matcher) *hamcrest.Matcher {
	match := func (actual interface{}) *hamcrest.Result {
		result := matcher.Match(actual)
		return hamcrest.NewResult(!result.Matched(), result).
			WithCauses(result.Causes()...)
	}
	return hamcrest.NewMatcherf(match, "Not[%v]", matcher)
}

// Returns a Matcher that decorates another matcher.
func Is(matcher *hamcrest.Matcher) *hamcrest.Matcher {
	match := func (actual interface{}) *hamcrest.Result {
		result := matcher.Match(actual)
		return hamcrest.NewResult(result.Matched(), result).
			WithCauses(result.Causes()...)
	}
	return hamcrest.NewMatcherf(match, "Is[%v]", matcher)
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
func Nil() *hamcrest.Matcher {
	return _Nil
}
var _Nil *hamcrest.Matcher // singleton
func init() {
	match := func (actual interface{}) *hamcrest.Result {
		if _detectNil(actual) {
			return hamcrest.NewResultf(true, "was nil")
		}
		return hamcrest.NewResultf(false, "[%v] was not nil", actual)
	}
	_Nil = hamcrest.NewMatcherf(match, "Nil")
}

// Returns a Matcher that matches if the actual value is 
// neither nil nor the nil value of its type.  (Note that
// this is *not* equivalent to Not(DeeplyEqualTo(nil)).)
func NonNil() *hamcrest.Matcher {
	return _NonNil
}
var _NonNil *hamcrest.Matcher
func init() {
	match := func (actual interface{}) *hamcrest.Result {
		if _detectNil(actual) {
			return hamcrest.NewResultf(false, "was nil")
		}
		return hamcrest.NewResultf(true, "[%v] was not nil", actual)
	}
	_NonNil = hamcrest.NewMatcherf(match, "NonNil")
}


// Returns a Matcher that checks if the actual value is (deeply)
// equal to the given expected value, using reflect.DeepEqual.
//
// For an equality test equivalent to `==`, see the
// `hamcrest/comparison` package.
func DeeplyEqualTo(expected interface{}) *hamcrest.Matcher {
	match := func (actual interface{}) *hamcrest.Result {
		if reflect.DeepEqual(expected, actual) {
			return hamcrest.NewResultf(true,
				"was deeply equal to [%v]", expected)
		}
		return hamcrest.NewResultf(false,
			"[%v] was not deeply equal to [%v]", actual, expected)
	}
	return hamcrest.NewMatcherf(match, "DeeplyEqualTo[%v]", expected)
}


// Returns a short-circuiting Matcher that matches whenever all of
// the given matchers match a given input value.  If any component
// matcher fails to match an input value, later matchers are not
// attempted.
func AllOf(matchers...*hamcrest.Matcher) *hamcrest.Matcher {
	match := func (actual interface{}) *hamcrest.Result {
		var results []*hamcrest.Result
		for index, matcher := range matchers {
			result := matcher.Match(actual)
			results := append(results, result)
			if !result.Matched() {
				return hamcrest.NewResultf(false,
					"Failed matcher %v of %v: [%v]",
					index+1, len(matchers), matcher).
					WithCauses(results...)
			}
		}
		return hamcrest.NewResultf(true,
			"Matched all %v matchers", len(matchers)).
			WithCauses(results...)
	}
	descriptions := make([]interface{}, len(matchers), len(matchers))
	for index, matcher := range matchers {
		descriptions[index] = hamcrest.Description("[#%v: %v]", index+1, matcher)
	}
	return hamcrest.NewMatcherf(match, "AllOf%v", descriptions)
}

// Returns a short-circuiting Matcher that matches whenever all of
// the given matchers match a given input value.  If any component
// matcher fails to match an input value, later matchers are not
// attempted.
func AnyOf(matchers...*hamcrest.Matcher) *hamcrest.Matcher {
	match := func (actual interface{}) *hamcrest.Result {
		var results []*hamcrest.Result
		for index, matcher := range matchers {
			result := matcher.Match(actual)
			results := append(results, result)
			if result.Matched() {
				return hamcrest.NewResultf(true,
					"Matched on matcher %v of %v: [%v]",
					index+1, len(matchers), matcher).
					WithCauses(results...)
			}
		}
		return hamcrest.NewResultf(false,
			"Matched none of the %v matchers", len(matchers)).
			WithCauses(results...)
	}
	descriptions := make([]interface{}, len(matchers), len(matchers))
	for index, matcher := range matchers {
		descriptions[index] = hamcrest.Description("[#%v: %v]", index+1, matcher)
	}
	return hamcrest.NewMatcherf(match, "AnyOf%v", descriptions)
}


// Returns a function that composes the given function with a Matcher, such as:
//    ToLength := Composer(func(s string) int { return len(s) }, "ToLength")
// And then:
//    HasLengthThree := ToLength(Is(EqualTo(3)))
//    HasLengthThree.Match("no").Matched() // false
//    HasLengthThree.Match("yes").Matched() // true
//
// The given function must be able to accept a single argument and
// return a single argument.
func Applying(function interface{}, name string) func(*hamcrest.Matcher) *hamcrest.Matcher  {
	funcValue := reflect.NewValue(function).(*reflect.FuncValue)
	funcType := funcValue.Type().(*reflect.FuncType)
	numIn := funcType.NumIn()
	numOut := funcType.NumOut()
	if numIn == 0 {
		panic(fmt.Sprintf("function must accept at least one value, was %v by function %v", numIn, function))
	}
	if numOut == 0 {
		panic(fmt.Sprintf("function must return at least one value, was %v by function %v", numOut, function))
	}
	return func(matcher *hamcrest.Matcher) *hamcrest.Matcher {
		match := func (actual interface{}) *hamcrest.Result {
			actualValue := reflect.NewValue(actual)
			argValues := make([]reflect.Value, numIn, numIn)
			if numIn > 0 {
				inType := funcType.In(0)
				if numIn == 1 && funcType.DotDotDot() {
					inSlice := reflect.MakeSlice(inType.(*reflect.SliceType), 1, 1)
					inSlice.Elem(0).SetValue(actualValue)
					argValues[0] = inSlice
				} else {
					argValues[0] = actualValue
				}
				for i := 1; i < numIn; i++ {
					inType = funcType.In(i)
					argValues[i] = reflect.MakeZero(inType)
				}
			}
			outValues := funcValue.Call(argValues)
			outValue := outValues[0]
			out := outValue.Interface()
			result := matcher.Match(out)
			return hamcrest.NewResultf(result.Matched(),
				"%v(%#v) = %v", name, actual, out).
				WithCauses(result)
		}
		return hamcrest.NewMatcherf(match, "%v[%v]", name, matcher)
	}
}
