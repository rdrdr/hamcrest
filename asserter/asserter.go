// Copyright 2011 Mick Killianey.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package asserter

import (
	"fmt"
	"io"
	"os"
	"github.com/rdrdr/hamcrest/base"
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
	LogResult(result *base.Result)
	
	// On a successful match, describes the result of applying
	// the given matcher to the given value.
	LogWhen(value interface{}, matcher *base.Matcher)
	
	// On an unsuccessful match, describes the result of applying
	// the given matcher to the given value.
	LogUnless(value interface{}, matcher *base.Matcher)
	
	// On a successful match, describes the result of applying
	// the given matcher to the given value and invokes Fail().
	FailWhen(value interface{}, matcher *base.Matcher)
	
	// On an unsuccessful match, describes the result of applying
	// the given matcher to the given value and invokes Fail().
	FailUnless(value interface{}, matcher *base.Matcher)
	
	// On a successful match, describes the result of applying
	// the given matcher to the given value and invokes FailNow().
	FailNowWhen(value interface{}, matcher *base.Matcher)
	
	// On an unsuccessful match, describes the result of applying
	// the given matcher to the given value and invokes FailNow().
	FailNowUnless(value interface{}, matcher *base.Matcher)
	
	// Equivalent to FailUnless.
	CheckThat(value interface{}, matcher *base.Matcher)
	
	// Equivalent to FailUnless with the True matcher.
	CheckTrue(value bool, messages ...interface{})
	
	// Equivalent to FailUnless with the False matcher.
	CheckFalse(value bool, messages ...interface{})
	
	// Equivalent to FailUnless with the Nil matcher.
	CheckNil(value interface{}, messages ...interface{})
	
	// Equivalent to FailUnless with the NonNil matcher.
	CheckNonNil(value interface{}, messages ...interface{})
	
	// Equivalent to FailNowUnless.
	AssertThat(value interface{}, matcher *base.Matcher)
	
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

// Convenience function to create an Asserter for stderr,
// as per UsingWriter().
func UsingStderr() Asserter {
	return UsingWriter(os.Stderr)
}

// Convenience function to create an Asserter for stdout,
// as per UsingWriter().
func UsingStdout() Asserter {
	return UsingWriter(os.Stdout)
}

// Convenience function to create an Asserter for stdout,
// as per UsingWriter().
func UsingFileNamed(filename string) Asserter {
	fileFlags := os.O_CREATE | os.O_APPEND | os.O_WRONLY | os.O_NONBLOCK
	f, err := os.Open(filename, fileFlags, 0666)
	if err != nil {
		panic("Can't open file named " + filename)
	}
	return UsingWriter(f)
}

// Convenience function to return an Asserter for which every
// operation is a no-op.
func ThatDoesNothing() Asserter {
	return &_NullAsserter{}
}

// --------------------------------------------------------------------
// Implementation
// --------------------------------------------------------------------

type _Flusher1 interface { Flush() }
type _Flusher2 interface { Flush() os.Error }

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

func (self *_LoggerUsingWriter) Flush() {
	switch w := self.writer.(type) {
	case _Flusher1: w.Flush()
	case _Flusher2: w.Flush()
	}
}


type _NullAsserter struct {}

func (self *_NullAsserter) Fail() {}
func (self *_NullAsserter) FailNow() {}
func (self *_NullAsserter) Failed() bool { return false }
func (self *_NullAsserter) LogResult(result *base.Result) { }
func (self *_NullAsserter) LogWhen(value interface{}, matcher *base.Matcher) { }
func (self *_NullAsserter) LogUnless(value interface{}, matcher *base.Matcher) { }
func (self *_NullAsserter) FailWhen(value interface{}, matcher *base.Matcher) { }
func (self *_NullAsserter) FailUnless(value interface{}, matcher *base.Matcher) { }
func (self *_NullAsserter) FailNowWhen(value interface{}, matcher *base.Matcher) { }
func (self *_NullAsserter) FailNowUnless(value interface{}, matcher *base.Matcher) { }
func (self *_NullAsserter) CheckThat(value interface{}, matcher *base.Matcher) { }
func (self *_NullAsserter) CheckTrue(value bool, comments ...interface{}) { }
func (self *_NullAsserter) CheckFalse(value bool, comments ...interface{}) { }
func (self *_NullAsserter) CheckNil(value interface{}, comments ...interface{}) { }
func (self *_NullAsserter) CheckNonNil(value interface{}, comments ...interface{}) { }
func (self *_NullAsserter) AssertThat(value interface{}, matcher *base.Matcher) { }
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

func safeMatch(value interface{}, matcher *base.Matcher) (result *base.Result) {
	defer func() {
		if x := recover(); x != nil {
			result = base.NewResultf(false, "Panic: %v", x).
				WithMatcherAndValue(matcher, value)
		}
	}()
	result = matcher.Match(value)
	return
}


// Insert final newline if needed and indent tabs after internal newlines.
func (self *_Asserter) _LogResult(indent string, result *base.Result) {
	value := result.Value()
	logger := self.logger
	if result.Matched() {
		logger.Logf("%vMATCHED input %v\n", indent, value)
	} else {
		logger.Logf("%vDID NOT MATCH input %v\n", indent, value)
	}
	detailsIndent := indent + "\t"
	matcher := result.Matcher()
	if matcher != nil {
		logger.Logf("%vMatcher: %v\n", detailsIndent, matcher)
	}
	logger.Logf("%vBecause: %v\n", detailsIndent, result)
	if matcher != nil {
		for _, comment := range matcher.Comments() {
			logger.Logf("%vComment: %v\n", detailsIndent, comment)
		}
		if causes := result.Causes(); len(causes) > 0 {
			if len(causes) == 1 {
				logger.Logf("%vCauses: (1 cause)\n", detailsIndent)
			} else {
				logger.Logf("%vCauses: (%v causes)\n", detailsIndent, len(causes))
			}
			for _, cause := range result.Causes() {
				self._LogResult(detailsIndent, cause)
			}
		}
		switch w := self.logger.(type) {
		case _Flusher1: w.Flush()
		case _Flusher2: w.Flush()
		}
	}
}

func (self *_Asserter) LogResult(result *base.Result) {
	self._LogResult("", result)
}

func (self *_Asserter) LogWhen(value interface{}, matcher *base.Matcher) {
	if result := safeMatch(value, matcher); result.Matched() {
		self.LogResult(result)
	}
}

func (self *_Asserter) LogUnless(value interface{}, matcher *base.Matcher) {
	if result := safeMatch(value, matcher); !result.Matched() {
		self.LogResult(result)
	}
}

func (self *_Asserter) FailWhen(value interface{}, matcher *base.Matcher) {
	if result := safeMatch(value, matcher); result.Matched() {
		self.LogResult(result)
		self.Fail()
	}
}
	
func (self *_Asserter) FailUnless(value interface{}, matcher *base.Matcher) {
	if result := safeMatch(value, matcher); !result.Matched() {
		self.LogResult(result)
		self.Fail()
	}
}

func (self *_Asserter) FailNowWhen(value interface{}, matcher *base.Matcher) {
	if result := safeMatch(value, matcher); result.Matched() {
		self.LogResult(result)
		self.FailNow()
	}
}
func (self *_Asserter) FailNowUnless(value interface{}, matcher *base.Matcher) {
	if result := safeMatch(value, matcher); !result.Matched() {
		self.LogResult(result)
		self.FailNow()
	}
}
	
func (self *_Asserter) CheckThat(value interface{}, matcher *base.Matcher) {
	self.FailUnless(value, matcher)
}

func _SelfDescribing(args interface{}) base.SelfDescribing {
	return base.Description("%v", args)
}

func (self *_Asserter) CheckTrue(value bool, comments ...interface{}) {
	self.CheckThat(value, base.True().Comment(comments...))
}

func (self *_Asserter) CheckFalse(value bool, comments ...interface{}) {
	self.CheckThat(value, base.False().Comment(comments...))
}

func (self *_Asserter) CheckNil(value interface{}, comments ...interface{}) {
	self.CheckThat(value, base.Nil().Comment(comments...))
}

func (self *_Asserter) CheckNonNil(value interface{}, comments ...interface{}) {
	self.CheckThat(value, base.NonNil().Comment(comments...))
}

func (self *_Asserter) AssertThat(value interface{}, matcher *base.Matcher) {
	self.FailNowUnless(value, matcher)
}

func (self *_Asserter) AssertTrue(value bool, comments ...interface{}) {
	self.AssertThat(value, base.True().Comment(comments...))
}

func (self *_Asserter) AssertFalse(value bool, comments ...interface{}) {
	self.AssertThat(value, base.False().Comment(comments...))
}

func (self *_Asserter) AssertNil(value interface{}, comments ...interface{}) {
	self.AssertThat(value, base.Nil().Comment(comments...))
}

func (self *_Asserter) AssertNonNil(value interface{}, comments ...interface{}) {
	self.AssertThat(value, base.NonNil().Comment(comments...))
}
