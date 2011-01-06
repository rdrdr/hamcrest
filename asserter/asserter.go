package asserter

import (
	"fmt"
	"io"
	"os"
	"hamcrest"
)

var _True = hamcrest.True
var _False = hamcrest.False
var _Nil = hamcrest.Nil
var _NonNil = hamcrest.NonNil


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
	
	// On a successful match, describes the result of applying
	// the given matcher to the given value.
	LogWhen(value interface{}, matcher *hamcrest.Matcher) *hamcrest.Result
	
	// On an unsuccessful match, describes the result of applying
	// the given matcher to the given value.
	LogUnless(value interface{}, matcher *hamcrest.Matcher) *hamcrest.Result
	
	// On a successful match, describes the result of applying
	// the given matcher to the given value and invokes Fail().
	FailWhen(value interface{}, matcher *hamcrest.Matcher) *hamcrest.Result
	
	// On an unsuccessful match, describes the result of applying
	// the given matcher to the given value and invokes Fail().
	FailUnless(value interface{}, matcher *hamcrest.Matcher) *hamcrest.Result
	
	// On a successful match, describes the result of applying
	// the given matcher to the given value and invokes FailNow().
	FailNowWhen(value interface{}, matcher *hamcrest.Matcher)
	
	// On an unsuccessful match, describes the result of applying
	// the given matcher to the given value and invokes FailNow().
	FailNowUnless(value interface{}, matcher *hamcrest.Matcher)
	
	// Equivalent to FailUnless.
	CheckThat(value interface{}, matcher *hamcrest.Matcher) *hamcrest.Result
	
	// Equivalent to FailUnless with the True matcher.
	CheckTrue(value bool, messages ...interface{}) *hamcrest.Result
	
	// Equivalent to FailUnless with the False matcher.
	CheckFalse(value bool, messages ...interface{}) *hamcrest.Result
	
	// Equivalent to FailUnless with the Nil matcher.
	CheckNil(value interface{}, messages ...interface{}) *hamcrest.Result
	
	// Equivalent to FailUnless with the NonNil matcher.
	CheckNonNil(value interface{}, messages ...interface{}) *hamcrest.Result
	
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
			result = hamcrest.NewResult(false,
				hamcrest.NewDescription("Panic: %v", x))
		}
	}()
	result = matcher.Match(value)
	return
}


// Insert final newline if needed and indent tabs after internal newlines.
func (self *_Asserter) _LogWithCauses(previousIndent string, result *hamcrest.Result) {
	indent := previousIndent + "\t"
	matcher := result.Matcher()
	if result.Matched() {
		self.logger.Logf("%vMatched %v\n", indent, matcher)
	} else {
		self.logger.Logf("%vDid not match %v\n", indent, matcher)
	}
	self.logger.Logf("%vBecause: %v\n", indent, result)
	for _, comment := range matcher.Comments() {
		self.logger.Logf("%vComment: %v\n", indent, comment)
	}
	for _, cause := range result.Causes() {
		self._LogWithCauses(indent, cause)
	}
}

func (self *_Asserter) LogWhen(value interface{}, matcher *hamcrest.Matcher) *hamcrest.Result {
	result := safeMatch(value, matcher)
	if result.Matched() {
		self.logger.Logf("MATCHED on input %#v\n", value)
		self._LogWithCauses("", result)
	}
	return result
}

func (self *_Asserter) LogUnless(value interface{}, matcher *hamcrest.Matcher) *hamcrest.Result {
	result := safeMatch(value, matcher)
	if !result.Matched() {
		self.logger.Logf("FAILURE on input %#v\n", value)
		self._LogWithCauses("", result)
	}
	return result
}

func (self *_Asserter) FailWhen(value interface{}, matcher *hamcrest.Matcher) *hamcrest.Result {
	result := self.LogWhen(value, matcher)
	if result.Matched() {
		self.Fail()
	}
	return result
}
	
func (self *_Asserter) FailUnless(value interface{}, matcher *hamcrest.Matcher) *hamcrest.Result {
	result := self.LogUnless(value, matcher)
	if !result.Matched() {
		self.Fail()
	}
	return result
}

func (self *_Asserter) FailNowWhen(value interface{}, matcher *hamcrest.Matcher) {
	if self.LogWhen(value, matcher).Matched() {
		self.FailNow()
	}
}
func (self *_Asserter) FailNowUnless(value interface{}, matcher *hamcrest.Matcher) {
	if !self.LogUnless(value, matcher).Matched() {
		self.FailNow()
	}
}
	
func (self *_Asserter) CheckThat(value interface{}, matcher *hamcrest.Matcher) *hamcrest.Result {
	return self.FailUnless(value, matcher)
}

func (self *_Asserter) CheckTrue(value bool, comments ...interface{}) *hamcrest.Result {
	return self.CheckThat(value, hamcrest.True().Comment(comments...))
}

func (self *_Asserter) CheckFalse(value bool, comments ...interface{}) *hamcrest.Result {
	return self.CheckThat(value, hamcrest.False().Comment(comments...))
}

func (self *_Asserter) CheckNil(value interface{}, comments ...interface{}) *hamcrest.Result {
	return self.CheckThat(value, hamcrest.Nil().Comment(comments...))
}

func (self *_Asserter) CheckNonNil(value interface{}, comments ...interface{}) *hamcrest.Result {
	return self.CheckThat(value, hamcrest.NonNil().Comment(comments...))
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
