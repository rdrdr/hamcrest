// Copyright 2011 Mick Killianey.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package asserter

import (
	"fmt"
	"io"
	"os"
	"hamcrest"
)


// Used by an Asserter to write log/error messages.
//
// A bit hacky...this was reverse-engineered to conform to testing.T.
type Logger interface {
	Logf(format string, args ...interface{})
	Failed() bool
	Fail()
	FailNow()
}

// Applies matchers to values, writing descriptions of the
// results to a Logger.
type Asserter interface {
	// Returns true if Fail() has been called.
	Failed() bool
	
	// Marks this asserter as having a failed assertion
	// but continues execution.
	Fail()
	
	// Marks this asserter as having a failed assertion
	// and invokes the immediate failure action.
	FailNow()
	
	// Writes the given result to the underlying logger.
	LogResult(result *hamcrest.Result)
	
	// On a successful match, describes the result of applying
	// the given matcher to the given value.
	LogWhen(value interface{}, matcher *hamcrest.Matcher)
	
	// On an unsuccessful match, describes the result of applying
	// the given matcher to the given value.
	LogUnless(value interface{}, matcher *hamcrest.Matcher)
	
	// On a successful match, describes the result of applying
	// the given matcher to the given value and invokes Fail().
	FailWhen(value interface{}, matcher *hamcrest.Matcher)
	
	// On an unsuccessful match, describes the result of applying
	// the given matcher to the given value and invokes Fail().
	FailUnless(value interface{}, matcher *hamcrest.Matcher)
	
	// On a successful match, describes the result of applying
	// the given matcher to the given value and invokes FailNow().
	FailNowWhen(value interface{}, matcher *hamcrest.Matcher)
	
	// On an unsuccessful match, describes the result of applying
	// the given matcher to the given value and invokes FailNow().
	FailNowUnless(value interface{}, matcher *hamcrest.Matcher)
	
	// Equivalent to FailUnless.
	CheckThat(value interface{}, matcher *hamcrest.Matcher)
	
	// Equivalent to FailUnless with the True matcher.
	CheckTrue(value bool, messages ...interface{})
	
	// Equivalent to FailUnless with the False matcher.
	CheckFalse(value bool, messages ...interface{})
	
	// Equivalent to FailUnless with the Nil matcher.
	CheckNil(value interface{}, messages ...interface{})
	
	// Equivalent to FailUnless with the NonNil matcher.
	CheckNonNil(value interface{}, messages ...interface{})
	
	// Equivalent to FailNowUnless.
	AssertThat(value interface{}, matcher *hamcrest.Matcher)
	
	// Equivalent to FailNowUnless with the True matcher.
	AssertTrue(value bool, messages ...interface{})
	
	// Equivalent to FailNowUnless with the False matcher.
	AssertFalse(value bool, messages ...interface{})
	
	// Equivalent to FailNowUnless with the Nil matcher.
	AssertNil(value interface{}, messages ...interface{})
	
	// Equivalent to FailNowUnless with the NonNil matcher.
	AssertNonNil(value interface{}, messages ...interface{})
}

// Convenience function to create an Asserter from a Logger.
// Note that (*testing.T) satisfies Logger, and can be used here.
func Using(logger Logger) Asserter {
	return &_Asserter{logger:logger}
}

// Convenience function to create an Asserter from an io.Writer
// and invoke panic() when the logger is asked to FailNow().
func UsingWriter(writer io.Writer) Asserter {
	failNow := func() { panic("Invoked FailNow()") }
	return UsingWriterAndFailNow(writer, failNow)
}

// Convenience function to create an Asserter from an io.Writer,
// and a custom FailNow() function.
func UsingWriterAndFailNow(writer io.Writer, failNow func()) Asserter {
	return Using(&_LoggerUsingWriter{writer:writer, failNow:failNow, failed:false})
}

// Convenience function to create an Asserter for standard
// error, as per UsingWriter().
func UsingStderr() Asserter {
	return UsingWriter(os.Stderr)
}

// Convenience function to return an Asserter for which every
// operation is a no-op.
func ThatDoesNothing() Asserter {
	return &_NullAsserter{}
}

// --------------------------------------------------------------------
// Implementation
// --------------------------------------------------------------------

type _LoggerUsingWriter struct {
	writer io.Writer
	failNow func()
	failed bool
}

func (self *_LoggerUsingWriter) Logf(format string, messages ...interface{}) {
	fmt.Fprintf(self.writer, format, messages...)
}
func (self *_LoggerUsingWriter) Failed() bool {
	return self.failed
}
func (self *_LoggerUsingWriter) Fail() {
	self.failed = true
}
func (self *_LoggerUsingWriter) FailNow() {
	self.failed = true
	self.failNow()
}


type _NullAsserter struct {}

func (self *_NullAsserter) Fail() {}
func (self *_NullAsserter) FailNow() {}
func (self *_NullAsserter) Failed() bool { return false }
func (self *_NullAsserter) LogResult(result *hamcrest.Result) { }
func (self *_NullAsserter) LogWhen(value interface{}, matcher *hamcrest.Matcher) { }
func (self *_NullAsserter) LogUnless(value interface{}, matcher *hamcrest.Matcher) { }
func (self *_NullAsserter) FailWhen(value interface{}, matcher *hamcrest.Matcher) { }
func (self *_NullAsserter) FailUnless(value interface{}, matcher *hamcrest.Matcher) { }
func (self *_NullAsserter) FailNowWhen(value interface{}, matcher *hamcrest.Matcher) { }
func (self *_NullAsserter) FailNowUnless(value interface{}, matcher *hamcrest.Matcher) { }
func (self *_NullAsserter) CheckThat(value interface{}, matcher *hamcrest.Matcher) { }
func (self *_NullAsserter) CheckTrue(value bool, comments ...interface{}) { }
func (self *_NullAsserter) CheckFalse(value bool, comments ...interface{}) { }
func (self *_NullAsserter) CheckNil(value interface{}, comments ...interface{}) { }
func (self *_NullAsserter) CheckNonNil(value interface{}, comments ...interface{}) { }
func (self *_NullAsserter) AssertThat(value interface{}, matcher *hamcrest.Matcher) { }
func (self *_NullAsserter) AssertTrue(value bool, comments ...interface{}) { }
func (self *_NullAsserter) AssertFalse(value bool, comments ...interface{}) { }
func (self *_NullAsserter) AssertNil(value interface{}, comments ...interface{}) { }
func (self *_NullAsserter) AssertNonNil(value interface{}, comments ...interface{}) { }

type _Asserter struct {
	logger Logger
}

func (self *_Asserter) Fail() {
	self.logger.Fail()
}

func (self *_Asserter) FailNow() {
	self.logger.FailNow()
}

func (self *_Asserter) Failed() bool {
	return self.logger.Failed()
}

func safeMatch(value interface{}, matcher *hamcrest.Matcher) (result *hamcrest.Result) {
	defer func() {
		if x := recover(); x != nil {
			result = hamcrest.NewResultf(false, "Panic: %v", x)
		}
	}()
	result = matcher.Match(value)
	return
}


// Insert final newline if needed and indent tabs after internal newlines.
func (self *_Asserter) _LogResult(indent string, result *hamcrest.Result) {
	matcher := result.Matcher()
	value := result.Value()
	if result.Matched() {
		self.logger.Logf("%vMATCHED input %v\n", indent, value)
	} else {
		self.logger.Logf("%vDID NOT MATCH input %v\n", indent, value)
	}
	self.logger.Logf("%vMatcher: %v\n", indent, matcher)
	self.logger.Logf("%vBecause: %v\n", indent, result)
	for _, comment := range matcher.Comments() {
		self.logger.Logf("%vComment: %v\n", indent, comment)
	}
	causeIndent := indent + "\t"
	for _, cause := range result.Causes() {
		self._LogResult(causeIndent, cause)
	}
}

func (self *_Asserter) LogResult(result *hamcrest.Result) {
	self._LogResult("", result)
}

func (self *_Asserter) LogWhen(value interface{}, matcher *hamcrest.Matcher) {
	if result := safeMatch(value, matcher); result.Matched() {
		self.LogResult(result)
	}
}

func (self *_Asserter) LogUnless(value interface{}, matcher *hamcrest.Matcher) {
	if result := safeMatch(value, matcher); !result.Matched() {
		self.LogResult(result)
	}
}

func (self *_Asserter) FailWhen(value interface{}, matcher *hamcrest.Matcher) {
	if result := safeMatch(value, matcher); result.Matched() {
		self.LogResult(result)
		self.Fail()
	}
}
	
func (self *_Asserter) FailUnless(value interface{}, matcher *hamcrest.Matcher) {
	if result := safeMatch(value, matcher); !result.Matched() {
		self.LogResult(result)
		self.Fail()
	}
}

func (self *_Asserter) FailNowWhen(value interface{}, matcher *hamcrest.Matcher) {
	if result := safeMatch(value, matcher); result.Matched() {
		self.LogResult(result)
		self.FailNow()
	}
}
func (self *_Asserter) FailNowUnless(value interface{}, matcher *hamcrest.Matcher) {
	if result := safeMatch(value, matcher); !result.Matched() {
		self.LogResult(result)
		self.FailNow()
	}
}
	
func (self *_Asserter) CheckThat(value interface{}, matcher *hamcrest.Matcher) {
	self.FailUnless(value, matcher)
}

func (self *_Asserter) CheckTrue(value bool, comments ...interface{}) {
	self.CheckThat(value, hamcrest.True().Comment(comments...))
}

func (self *_Asserter) CheckFalse(value bool, comments ...interface{}) {
	self.CheckThat(value, hamcrest.False().Comment(comments...))
}

func (self *_Asserter) CheckNil(value interface{}, comments ...interface{}) {
	self.CheckThat(value, hamcrest.Nil().Comment(comments...))
}

func (self *_Asserter) CheckNonNil(value interface{}, comments ...interface{}) {
	self.CheckThat(value, hamcrest.NonNil().Comment(comments...))
}

func (self *_Asserter) AssertThat(value interface{}, matcher *hamcrest.Matcher) {
	self.FailNowUnless(value, matcher)
}

func (self *_Asserter) AssertTrue(value bool, comments ...interface{}) {
	self.AssertThat(value, hamcrest.True().Comment(comments...))
}

func (self *_Asserter) AssertFalse(value bool, comments ...interface{}) {
	self.AssertThat(value, hamcrest.False().Comment(comments...))
}

func (self *_Asserter) AssertNil(value interface{}, comments ...interface{}) {
	self.AssertThat(value, hamcrest.Nil().Comment(comments...))
}

func (self *_Asserter) AssertNonNil(value interface{}, comments ...interface{}) {
	self.AssertThat(value, hamcrest.NonNil().Comment(comments...))
}
