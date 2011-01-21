// Copyright 2011 Mick Killianey.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hamcrest

import (
	"testing"
)

func Test_Description(t *testing.T) {
	description := Description("%v %v", "foo", "bar")
	descriptionString := description.String()
	if descriptionString != "foo bar" {
		t.Errorf("Description should be 'foo bar', was %v", descriptionString)
	}
}

func checkResultIsMatching(t *testing.T, result *Result, message string) {
	if !result.Matched() {
		t.Errorf("Expected matching result from applying %v to %#v, was [%v] %v",
			result.Matcher(), result.Value(), result, message)
	}
}

func checkResultIsNonMatching(t *testing.T, result *Result, message string) {
	if result.Matched() {
		t.Errorf("Expected non-matching result from applying %v to %#v, was [%v].  Message: %v",
			result.Matcher(), result.Value(), result, message)
	}
}


func Test_Matched(t *testing.T) {
	matcher := Matched()
	passResult := NewResultf(true, "pass")
	failResult := NewResultf(false, "fail")
	
	checkResultIsMatching(t, matcher.Match(passResult), "matching")
	checkResultIsNonMatching(t, matcher.Match(failResult), "non-matching")
	checkResultIsNonMatching(t, matcher.Match(nil), "nil")
	checkResultIsNonMatching(t, matcher.Match("foo"), "not a Result")
}

func Test_DidNotMatch(t *testing.T) {
	matcher := DidNotMatch()
	passResult := NewResultf(true, "pass")
	failResult := NewResultf(false, "fail")
	
	checkResultIsNonMatching(t, matcher.Match(passResult), "matching")
	checkResultIsMatching(t, matcher.Match(failResult), "non-matching")
	checkResultIsNonMatching(t, matcher.Match(nil), "nil")
	checkResultIsNonMatching(t, matcher.Match("foo"), "not a Result")
}

