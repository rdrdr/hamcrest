// Copyright 2011 Mick Killianey.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import (
	"fmt"
	"github.com/rdrdr/hamcrest/base"
	"reflect"
)

// Returns a Matcher that matches any input value.
func Anything() *base.Matcher {
	return _Anything
}
var _Anything *base.Matcher // singleton
func init() {
	match := func (actual interface{}) *base.Result {
		return base.NewResultf(true, "always matches")
	}
	_Anything = base.NewMatcherf(match, "Anything")
}


// Returns a Matcher that matches the boolean value true.
func True() *base.Matcher {
	return base.True()
}

// Returns a Matcher that matches the boolean value false.
func False() *base.Matcher {
	return base.False()
}

// Returns a Matcher that matches on values that cause the given
// functionOrMatcher to panic.
//
// functionOrMatcher should either be a function that accepts one
// parameter or a Matcher.
func PanicWhenApplying(functionOrMatcher interface{}, name string) *base.Matcher {
	var doSomething func(interface{})
	if matcher, ok := functionOrMatcher.(*base.Matcher); ok {
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
	match := func (actual interface{}) (result *base.Result) {
		defer func() {
			if recover() != nil {
				result = base.NewResultf(true, "Panicked")
			}
		}()
		doSomething(actual)
		result = base.NewResultf(false, "Did not panic")
		return
	}
	return base.NewMatcherf(match,
		"PanicWhenApplying[%v]", name)
}


// Returns a Matcher that decorates another matcher and only matches
// when the underlying matcher does not match (and vice versa).
func Not(matcher *base.Matcher) *base.Matcher {
	match := func (actual interface{}) *base.Result {
		result := matcher.Match(actual)
		if result.Matched() {
			return base.NewResultf(false,
				"'Not' failed because inner matcher passed: %v", matcher).
				WithCauses(result)
		}
		return base.NewResultf(true,
			"'Not' passed because inner matcher failed: %v", matcher).
			WithCauses(result)
	}
	return base.NewMatcherf(match, "Not[%v]", matcher)
}

// Returns a Matcher that decorates another matcher.
func Is(matcher *base.Matcher) *base.Matcher {
	match := func (actual interface{}) *base.Result {
		result := matcher.Match(actual)
		return base.NewResult(result.Matched(), result).
			WithCauses(result.Causes()...)
	}
	return base.NewMatcherf(match, "Is[%v]", matcher)
}


// Returns a Matcher that matches if the actual value is nil
// or the nil value of its type.  (Note that this is *not*
// equivalent to DeeplyEqualTo(nil).)
func Nil() *base.Matcher {
	return base.Nil()
}

// Returns a Matcher that matches if the actual value is 
// neither nil nor the nil value of its type.  (Note that
// this is *not* equivalent to Not(DeeplyEqualTo(nil)).)
func NonNil() *base.Matcher {
	return base.NonNil()
}

// Returns a Matcher that checks if the actual value is (deeply)
// equal to the given expected value, using reflect.DeepEqual.
//
// For an equality test equivalent to `==`, see the
// `hamcrest/comparison` package.
func DeepEqualTo(expected interface{}) *base.Matcher {
	return base.DeepEqualTo(expected)
}

// Returns a matcher that matches values that are greater-than the given
// expected value, using the greater-than (<) operator.
func GreaterThan(expected interface{}) *base.Matcher {
	return base.GreaterThan(expected)
}

// Returns a matcher that matches values that are greater-than-or-equal-to
// the given expected value, using the greater-than-or-equal-to (>=) operator.
func GreaterThanOrEqualTo(expected interface{}) *base.Matcher {
	return base.GreaterThanOrEqualTo(expected)
}

// Returns a matcher that matches values that are less-than the given
// expected value, using the less-than (<) operator.
func LessThan(expected interface{}) *base.Matcher {
	return base.LessThan(expected)
}

// Returns a matcher that matches values that are less-than-or-equal-to
// the given expected value, using the less-than-or-equal-to (<=) operator.
func LessThanOrEqualTo(expected interface{}) *base.Matcher {
	return base.LessThanOrEqualTo(expected)
}

// Returns a matcher that matches values that are equal to the
// given expected value, using the equality (==) operator.
func EqualTo(expected interface{}) *base.Matcher {
	return base.EqualTo(expected)
}

// Returns a matcher that matches values that are not equal to the
// given expected value, using the inequality (!=) operator.
func NotEqualTo(expected interface{}) *base.Matcher {
	return base.NotEqualTo(expected)
}

// Returns a short-circuiting Matcher that matches whenever all of
// the given matchers match a given input value.  If any component
// matcher fails to match an input value, later matchers are not
// attempted.
func AllOf(matchers...*base.Matcher) *base.Matcher {
	match := func (actual interface{}) *base.Result {
		var results []*base.Result
		for index, matcher := range matchers {
			result := matcher.Match(actual)
			results := append(results, result)
			if !result.Matched() {
				return base.NewResultf(false,
					"Failed matcher %v of %v: [%v]",
					index+1, len(matchers), matcher).
					WithCauses(results...)
			}
		}
		return base.NewResultf(true,
			"Matched all %v matchers", len(matchers)).
			WithCauses(results...)
	}
	descriptions := make([]interface{}, len(matchers), len(matchers))
	for index, matcher := range matchers {
		descriptions[index] = base.Description("[#%v: %v]", index+1, matcher)
	}
	return base.NewMatcherf(match, "AllOf%v", descriptions)
}

// Returns a short-circuiting Matcher that matches whenever all of
// the given matchers match a given input value.  If any component
// matcher fails to match an input value, later matchers are not
// attempted.
func AnyOf(matchers...*base.Matcher) *base.Matcher {
	match := func (actual interface{}) *base.Result {
		var results []*base.Result
		for index, matcher := range matchers {
			result := matcher.Match(actual)
			results := append(results, result)
			if result.Matched() {
				return base.NewResultf(true,
					"Matched on matcher %v of %v: [%v]",
					index+1, len(matchers), matcher).
					WithCauses(results...)
			}
		}
		return base.NewResultf(false,
			"Matched none of the %v matchers", len(matchers)).
			WithCauses(results...)
	}
	descriptions := make([]interface{}, len(matchers), len(matchers))
	for index, matcher := range matchers {
		descriptions[index] = base.Description("[#%v: %v]", index+1, matcher)
	}
	return base.NewMatcherf(match, "AnyOf%v", descriptions)
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
func Applying(function interface{}, name string) func(*base.Matcher) *base.Matcher  {
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
	return func(matcher *base.Matcher) *base.Matcher {
		match := func (actual interface{}) *base.Result {
			Assign := func(dst reflect.Value, src reflect.Value) (ok bool) {
				defer func() { recover(); }()
				dst.SetValue(src)
				ok = true
				return
			}
			actualValue := reflect.NewValue(actual)
			argValues := make([]reflect.Value, numIn, numIn)
			if numIn > 0 {
				inType := funcType.In(0)
				if numIn == 1 && funcType.DotDotDot() {
					inSlice := reflect.MakeSlice(inType.(*reflect.SliceType), 1, 1)
					if !Assign(inSlice.Elem(0), actualValue) {
						return base.NewResultf(false,
							"Cannot use %T as input to %T", actual, function)
					}
					argValues[0] = inSlice
				} else {
					inValue := reflect.MakeZero(inType)
					if !Assign(inValue, actualValue) {
						return base.NewResultf(false,
							"Cannot use %T as input to %T", actual, function)
					}
					argValues[0] = inValue
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
			return base.NewResultf(result.Matched(),
				"%v(%#v) = %v", name, actual, out).
				WithCauses(result)
		}
		return base.NewMatcherf(match, "%v[%v]", name, matcher)
	}
}
