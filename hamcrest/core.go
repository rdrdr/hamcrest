// Copyright 2011 Mick Killianey.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hamcrest

import (
	"fmt"
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

// Returns a Matcher that matches on values that cause the given
// functionOrMatcher to panic.
//
// functionOrMatcher should either be a function that accepts one
// parameter or a Matcher
func PanicWhenApplying(functionOrMatcher interface{}, name string) *Matcher {
	var doSomething func(interface{})
	if matcher, ok := functionOrMatcher.(*Matcher); ok {
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
	match := func (actual interface{}) (result *Result) {
		defer func() {
			if recover() != nil {
				result = NewResultf(true, "Panicked")
			}
		}()
		doSomething(actual)
		result = NewResultf(false, "Did not panic")
		return
	}
	return NewMatcherf(match, "PanicWhenApplying[%v]", name)
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


// Returns a function that composes the given function with a Matcher, such as:
//    ToLength := Composer(func(s string) int { return len(s) }, "Len")
// And then:
//    HasLengthThree := ToLength(Is(EqualTo(3)))
//    HasLengthThree.Match("no").Matched() // false
//    HasLengthThree.Match("yes").Matched() // true
//
// The given function must be able to accept a single argument and
// return a single argument.
func Composer(function interface{}, name string) func(*Matcher) *Matcher  {
	funcValue := reflect.NewValue(function).(*reflect.FuncValue)
	funcType := funcValue.Type().(*reflect.FuncType)
	numIn := funcType.NumIn()
	switch {
	case numIn == 1:
	case numIn == 2 && funcType.DotDotDot():
	default:
		panic(fmt.Sprintf("func must accept a single arg, was %T", function))
	}
	inType := funcType.In(0)
	numOut := funcType.NumOut()
	if numOut != 1 {
		panic(fmt.Sprintf("function must return exactly one value, was %v by function %v", numOut, function))
	}
	return func(matcher *Matcher) *Matcher {
		match := func (actual interface{}) *Result {
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
					inSlice := reflect.MakeSlice(inType2.(*reflect.SliceType), 0, 0)
					argValues[1] = inSlice
				}
			}
			outValues := funcValue.Call(argValues)
			outValue := outValues[0]
			out := outValue.Interface()
			result := matcher.Match(out)
			return NewResultf(result.Matched(),
				"%v(%#v) = %v", name, actual, out).
				WithCauses(result)
		}
		return NewMatcherf(match, "%v[%v]", name, matcher)
	}
}

