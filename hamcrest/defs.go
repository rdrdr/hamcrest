// Copyright 2011 Mick Killianey.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hamcrest

import (
	"fmt"
	"reflect"
)

// --------------------------------------------------------------------
// Description
// --------------------------------------------------------------------

// Hamcrest descriptions implement both fmt.Stringer and fmt.Formatter.
type SelfDescribing interface {
	fmt.Formatter
	fmt.Stringer
}

type _Description struct {
	format string
	args []interface{}
}

// Implements fmt.Formatter.
func (self *_Description) Format(s fmt.State, ch int) {
	fmt.Fprintf(s, self.format, self.args...)
}

// Implements fmt.Stringer.
func (self *_Description) String() string {
	return fmt.Sprintf(self.format, self.args...)
}


// Creates an object that implements fmt.Formatter and fmt.Stringer using
// args with the same meanings as fmt.Fprintf.  Note that this object
// stores its given parameters and evaluates them lazily.
func Description(format string, args...interface{}) SelfDescribing {
	return &_Description{format:format, args:args}
}


// --------------------------------------------------------------------
// Result
// --------------------------------------------------------------------

// Self-describing result of applying a Matcher to an input value.
type Result struct {
	description SelfDescribing
	matched bool
	value interface{}
	matcher *Matcher
	causes []*Result
}
// Creates a new Result using the given description.
func NewResult(matched bool, description SelfDescribing) *Result {
	return &Result{ matched: matched, description: description }
}

// Creates a new Result using the given format/args as a description.
func NewResultf(matched bool, format string, args...interface{}) *Result {
	return NewResult(matched, Description(format, args...))
}

// Returns true if the Match was successful.
func (self *Result) Matched() bool {
	return self.matched
}

// Returns the Matcher that produced this Result.
func (self *Result) Matcher() *Matcher {
	return self.matcher
}

// Returns the value that was given to the Matcher to produce this Result.
func (self *Result) Value() interface{} {
	return self.value
}

// Implements fmt.Stringer.
func (self *Result) String() string {
	return self.description.String()
}

// Implements fmt.Formatter.
func (self *Result) Format(s fmt.State, ch int) {
	self.description.Format(s, ch)
}

// Returns a slice of sub-Results that caused this Result to
// either match or not match.
func (self *Result) Causes() []*Result {
	if self.causes == nil {
		return nil
	}
	causes := make([]*Result, len(self.causes))
	copy(causes, self.causes)
	return causes
}

// Returns a new Result, identical to this one, except with
// the given causes.
func (self *Result) WithCauses(causes... *Result) *Result {
	return &Result{
		matched:self.matched,
		description:self.description,
		matcher:self.matcher,
		value:self.value,
		causes:causes}
}

// Returns a new Result, identical to this one, except with
// the given matcher and value.
func (self *Result) WithMatcherAndValue(matcher *Matcher, value interface{}) *Result {
	return &Result{
		matched:self.matched,
		description:self.description,
		matcher:matcher,
		value:value,
		causes:self.Causes()}
}

// --------------------------------------------------------------------
// Matcher
// --------------------------------------------------------------------

// Self-describing criteria that may match (or not match) an input value.
// Creators of new matchers are *strongly* encouraged not to implement
// Matcher directly, but to create new matchers using the NewMatcher
// factory function.
type Matcher struct {
	description SelfDescribing
	match func(v interface{}) *Result
	comments []interface{}
}

// Creates a new Matcher using the given matching function, with the
// given description.
//
// Matching functions must have exactly one input parameter.  Matching
// functions must return a *Result or a bool as its first output value.
// If the first output value is a *Result, that must be its only output
// value.  If the first output value is bool, then that value is used
// for Result.Matched(), and if present, any second output value is used
// as the result description.  For example, the following signatures are
// all legal matcher functions:
//
//   func(interface{}) *Result
//   func(interface{}) bool
//   func(interface{}) (bool, os.Error)
//   func(string) *Result
//   func(...*int) bool
//   func(...io.Reader) (bool os.Error)
func NewMatcher(fn interface{}, description SelfDescribing) *Matcher {
	match := normalizeMatchFunction(fn)
	return &Matcher{ match: match, description: description }
}

// Creates a new Matcher using the given matching function, with the
// given format/args as a description.
func NewMatcherf(fn interface{}, format string, args...interface{}) *Matcher {
	match := normalizeMatchFunction(fn)
	return NewMatcher(match, Description(format, args...))
}

func normalizeMatchFunction(fn interface{}) func(interface{}) *Result {
	if match, ok := fn.(func(interface{}) *Result); ok {
		return match
	}
	if funcValue, ok := reflect.NewValue(fn).(*reflect.FuncValue); ok {
		funcType := funcValue.Type().(*reflect.FuncType)
		numIn := funcType.NumIn()
		numOut := funcType.NumOut()
		var constructInputValues func(actual interface{}) ([]reflect.Value, interface{})
		if numIn == 1 {
			constructInputValues = func(actual interface{}) (values []reflect.Value, err interface{}) {
				defer func() {
					err = recover()
				}()
				actualValue := reflect.NewValue(actual)
				inType := funcType.In(0)
				if funcType.DotDotDot() {
					sliceValue := reflect.MakeSlice(inType.(*reflect.SliceType), 1, 1)
					inputValue := sliceValue.Elem(0)
					inputValue.SetValue(actualValue)
					values = []reflect.Value{ sliceValue }
				} else {
					inputValue := reflect.MakeZero(inType)
					inputValue.SetValue(actualValue)
					values = []reflect.Value{ inputValue }
				}
				return
			}
		} else {
			reason := fmt.Sprintf("Can't use %T as a matcher: must have one input, had %v",
				fn, numIn)
			panic(reason)
		}
		outType0 := funcType.Out(0)
		var resultType = reflect.Typeof(&Result{})
		var boolType = reflect.Typeof(true)
		var interpretOutputValues func(values []reflect.Value) *Result
		if numOut == 1 && outType0 == resultType {
			interpretOutputValues = func(values []reflect.Value) (result *Result) {
				defer func() {
					if err := recover(); err != nil {
						result = NewResultf(false,
							"Error:  expected *Result from %T, got %v", fn, values)
					}
				}()
				result = values[0].Interface().(*Result)
				return
			}
		} else if outType0 == boolType && (numOut == 1 || numOut == 2) {
			interpretOutputValues = func(values []reflect.Value) (result *Result) {
				defer func() {
					if err := recover(); err != nil {
						result = NewResultf(false,
							"Error:  expected bool from %T, got %v", fn, values)
					}
				}()
				b := values[0].Interface().(bool)
				if len(values) > 1 {
					result = NewResultf(b, "%v", values[1].Interface())
				} else if b {
					result = NewResultf(true, "Matched")
				} else {
					result = NewResultf(false, "Did not match")
				}
				return
			}
		} else {
			panic(fmt.Sprintf("Can't use %T as a matcher function", fn))
		}
		return func(actual interface{}) *Result {
			inputValues, problem := constructInputValues(actual)
			if problem != nil {
				return NewResultf(false, "Could not apply %T to input of type %T: %v",
					fn, actual, problem)
			}
			outputValues := funcValue.Call(inputValues)
			return interpretOutputValues(outputValues)
		}
	
	}
	return nil
}

// Implementation of SelfDescribing: fmt.Formatter.
func (self *Matcher) Format(s fmt.State, ch int) {
	self.description.Format(s, ch)
}

// Implementation of SelfDescribing: fmt.Stringer.
func (self *Matcher) String() string {
	return self.description.String()
}

// Tests the given input value to see if it meets this Matcher's criteria.
func (self *Matcher) Match(value interface{}) *Result {
	result := self.match(value)
	result.matcher = self
	result.value = value
	return result
}

// Returns a slice of messages that supplement the description.
func (self *Matcher) Comments() []interface{} {
	comments := make([]interface{}, len(self.comments))
	copy(comments, self.comments)
	return comments
}

// Returns a *new* Matcher similar to this one, but with the
// given additional comments.
func (self *Matcher) Comment(comments... interface{}) *Matcher {
	all := make([]interface{}, 0, len(self.comments) + len(comments))
	copy(all, self.comments)
	all = append(all, comments...)
	return &Matcher{description:self.description, match:self.match, comments:all}
}
