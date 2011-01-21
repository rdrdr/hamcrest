// Copyright 2011 Mick Killianey.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hamcrest

import (
	"fmt"
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
// Creates a new Matcher using the given description.
func NewMatcher(match func(v interface{}) *Result, description SelfDescribing) *Matcher {
	return &Matcher{ match: match, description: description }
}

// Creates a new Matcher using the given format/args as a description.
func NewMatcherf(match func(v interface{}) *Result, format string, args...interface{}) *Matcher {
	return NewMatcher(match, Description(format, args...))
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

