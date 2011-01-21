// Copyright 2011 Mick Killianey.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package asserter

import (
	"bytes"
	"strings"
	"testing"
	"hamcrest"
	"hamcrest/core"
)

func newBuffer() *bytes.Buffer{
	return bytes.NewBuffer(make([]byte, 0, 1000))
}

func checkResultIsMatching(t *testing.T, result *hamcrest.Result, message string) {
	if !result.Matched() {
		t.Errorf("Expected matching result, was [%v] %v", result, message)
	}
}

func checkResultIsNonMatching(t *testing.T, result *hamcrest.Result, message string) {
	if result.Matched() {
		t.Errorf("Expected non-matching result, was [%v] %v", result, message)
	}
}

func checkBufferIsEmpty(t *testing.T, buffer *bytes.Buffer) {
	if buffer.Len() > 0 {
		t.Errorf("Should not have written to buffer, was (len: %v):\n%v",
			buffer.Len(), buffer.String())
	}
}

func checkAsserterFailed(t *testing.T, asserter Asserter) {
	if !asserter.Failed() {
		t.Errorf("Should have failed")
	}
}

func checkAsserterDidNotFail(t *testing.T, asserter Asserter) {
	if asserter.Failed() {
		t.Errorf("Should not have failed")
	}
}

func checkBufferContainsStrings(
		t *testing.T, buffer *bytes.Buffer, pieces...string) {
	src := buffer.String()
	for _, piece := range pieces {
		if !strings.Contains(src, piece) {
			t.Errorf("Expected string to contain [%v], was [%v]", piece, src)
		}
	}
}

func checkBufferContainsMatchingStrings(t *testing.T, buffer *bytes.Buffer) {
	checkBufferContainsStrings(t, buffer,
			MATCHING_VALUE, MATCHING_RESULT, MATCHER_DESCRIPTION)
	checkBufferContainsStrings(t, buffer, MATCHER_COMMENTS...)
}

func checkBufferContainsNonMatchingStrings(t *testing.T, buffer *bytes.Buffer) {
	checkBufferContainsStrings(t, buffer,
			NONMATCHING_VALUE, NONMATCHING_RESULT, MATCHER_DESCRIPTION)
	checkBufferContainsStrings(t, buffer, MATCHER_COMMENTS...)
}

// Sample matcher and values for tests below
var MATCHING_VALUE = "expected_value"
var NONMATCHING_VALUE = "different_value"
var MATCHING_RESULT = "was_matching"
var NONMATCHING_RESULT = "was_nonmatching"
var MATCHER_DESCRIPTION = "matcher_message"
var MATCHER_COMMENT1 = "matcher_message_1_with_%v" // catch accidental expansion
var MATCHER_COMMENT2 = 2
var MATCHER_COMMENTS = []string{ MATCHER_COMMENT1, "2" }
var MATCHER = hamcrest.NewMatcherf(
	func (actual interface{}) *hamcrest.Result {
		if actual == MATCHING_VALUE {
			return hamcrest.NewResultf(true, MATCHING_RESULT)
		}
		return hamcrest.NewResultf(false, NONMATCHING_RESULT)
	}, MATCHER_DESCRIPTION).Comment(MATCHER_COMMENT1, MATCHER_COMMENT2)

func Test_LogWhen_onNonMatchingResult(t *testing.T) {
	buffer := newBuffer()
	asserter := UsingWriter(buffer)
	asserter.LogWhen(NONMATCHING_VALUE, MATCHER)
	checkBufferIsEmpty(t, buffer)
}

func Test_LogWhen_onMatchingResult(t *testing.T) {
	buffer := newBuffer()
	asserter := UsingWriter(buffer)
	asserter.LogWhen(MATCHING_VALUE, MATCHER)
	checkBufferContainsMatchingStrings(t, buffer)
}

func Test_LogUnless_onMatchingResult(t *testing.T) {
	buffer := newBuffer()
	asserter := UsingWriter(buffer)
	asserter.LogUnless(MATCHING_VALUE, MATCHER)
	checkBufferIsEmpty(t, buffer)
}

func Test_LogUnless_onNonMatchingResult(t *testing.T) {
	buffer := newBuffer()
	asserter := UsingWriter(buffer)
	asserter.LogUnless(NONMATCHING_VALUE, MATCHER)
	checkBufferContainsNonMatchingStrings(t, buffer)

}

func Test_FailWhen_onNonMatchingResult(t *testing.T) {
	buffer := newBuffer()
	asserter := UsingWriter(buffer)
	asserter.FailWhen(false, core.True())
	checkBufferIsEmpty(t, buffer)
	checkAsserterDidNotFail(t, asserter)
}
func Test_FailWhen_onMatchingResult(t *testing.T) {
	buffer := newBuffer()
	asserter := UsingWriter(buffer)
	asserter.FailWhen(MATCHING_VALUE, MATCHER)
	checkBufferContainsMatchingStrings(t, buffer)
	checkAsserterFailed(t, asserter)
}

func Test_FailUnless_onMatchingResult(t *testing.T) {
	buffer := newBuffer()
	asserter := UsingWriter(buffer)
	asserter.FailUnless(MATCHING_VALUE, MATCHER)
	checkBufferIsEmpty(t, buffer)
	checkAsserterDidNotFail(t, asserter)
}
func Test_FailUnless_onNonMatchingResult(t *testing.T) {
	buffer := newBuffer()
	asserter := UsingWriter(buffer)
	asserter.FailUnless(NONMATCHING_VALUE, MATCHER)
	checkBufferContainsNonMatchingStrings(t, buffer)
	checkAsserterFailed(t, asserter)
}


func Test_FailNowWhen_onNonMatchingResult(t *testing.T) {
	buffer := newBuffer()
	calledFailNow := false
	asserter := UsingWriterAndFailNow(buffer, func() { calledFailNow = true})
	asserter.FailNowWhen(NONMATCHING_VALUE, MATCHER)
	checkBufferIsEmpty(t, buffer)
	checkAsserterDidNotFail(t, asserter)
	if calledFailNow {
		t.Error("Should not have called failNow")
	}
}
func Test_FailNowWhen_onMatchingResult(t *testing.T) {
	buffer := newBuffer()
	calledFailNow := false
	asserter := UsingWriterAndFailNow(buffer, func() { calledFailNow = true})
	asserter.FailNowWhen(MATCHING_VALUE, MATCHER)
	checkBufferContainsMatchingStrings(t, buffer)
	checkAsserterFailed(t, asserter)
	if !calledFailNow {
		t.Error("Should have called failNow")
	}
}

func Test_FailNowUnless_onMatchingResult(t *testing.T) {
	buffer := newBuffer()
	calledFailNow := false
	asserter := UsingWriterAndFailNow(buffer, func() { calledFailNow = true})
	asserter.FailNowUnless(MATCHING_VALUE, MATCHER)
	checkBufferIsEmpty(t, buffer)
	checkAsserterDidNotFail(t, asserter)
	if calledFailNow {
		t.Error("Should not have called failNow")
	}
}
func Test_FailNowUnless_onNonMatchingResult(t *testing.T) {
	buffer := newBuffer()
	calledFailNow := false
	asserter := UsingWriterAndFailNow(buffer, func() { calledFailNow = true})
	asserter.FailNowUnless(NONMATCHING_VALUE, MATCHER)
	checkBufferContainsNonMatchingStrings(t, buffer)
	checkAsserterFailed(t, asserter)
	if !calledFailNow {
		t.Error("Should have called failNow")
	}
}

func Test_CheckThat_onMatchingResult(t *testing.T) {
	buffer := newBuffer()
	calledFailNow := false
	asserter := UsingWriterAndFailNow(buffer, func() { calledFailNow = true} )
	asserter.CheckThat(MATCHING_VALUE, MATCHER)
	checkBufferIsEmpty(t, buffer)
	checkAsserterDidNotFail(t, asserter)
	if calledFailNow {
		t.Error("Should not have called failNow")
	}
}

func Test_CheckThat_onNonMatchingResult(t *testing.T) {
	buffer := newBuffer()
	calledFailNow := false
	asserter := UsingWriterAndFailNow(buffer, func() { calledFailNow = true} )
	asserter.CheckThat(NONMATCHING_VALUE, MATCHER)
	checkBufferContainsNonMatchingStrings(t, buffer)
	checkAsserterFailed(t, asserter)
	if calledFailNow {
		t.Error("Should not have called failNow")
	}
}

func Test_AssertThat(t *testing.T) {
	buffer := newBuffer()
	calledFailNow := false
	asserter := UsingWriterAndFailNow(buffer, func() { calledFailNow = true} )
	asserter.AssertThat(NONMATCHING_VALUE, MATCHER)
	checkBufferContainsNonMatchingStrings(t, buffer)
	checkAsserterFailed(t, asserter)
	if !calledFailNow {
		t.Error("Should have called failNow")
	}
}

func Test_NullAsserter(t *testing.T) {
	asserter := ThatDoesNothing()
	snooped := false
	snoopMatcher := hamcrest.NewMatcherf(func(v interface{}) *hamcrest.Result { 
			snooped = true
			return hamcrest.NewResultf(true, "snooped!")
		}, "Snoop")
	asserter.AssertThat("x", snoopMatcher)
	if snooped {
		t.Fatal("Calling AssertThat() should not invoke matcher!")
	}
	asserter.AssertTrue(false, "Should ignore attempts to AssertTrue")
	asserter.AssertFalse(true, "Should ignore attempts to AssertFalse")
	asserter.AssertNil("ha!", "Should ignore attempts to AssertNil")
	asserter.AssertNonNil(nil, "Should ignore attempts to AssertNonNil")
}