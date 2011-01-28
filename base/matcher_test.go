// Copyright 2011 Mick Killianey.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package base

import (
	"fmt"
	"strings"
	"testing"
)

func Test_Description(t *testing.T) {
	description := Description("%v %v", "foo", "bar")
	descriptionString := description.String()
	if descriptionString != "foo bar" {
		t.Errorf("Description should be 'foo bar', was %v", descriptionString)
	}
}

func attemptMatch(matcher *Matcher, input interface{}) (result *Result, err interface{}) {
	defer func() {
		if err = recover(); err != nil {
			result = nil
		}
	}()
	result = matcher.Match(input)
	return
}

func checkResultIsMatching(t *testing.T, matcher *Matcher, input interface{}, message string) *Result {
	result, err := attemptMatch(matcher, input)
	if err != nil {
		t.Errorf("Expected matching result from applying %v to %#v, was panic [%v]",
			matcher, input, err)
		return nil
	}
	if !result.Matched() {
		t.Errorf("Expected matching result from applying %v to %#v, was [%v] %v",
			matcher, input, result, message)
	}
	return result
}

func checkResultIsNonMatching(t *testing.T, matcher *Matcher, input interface{}, message string) *Result {
	result, err := attemptMatch(matcher, input)
	if err != nil {
		t.Errorf("Expected non-matching result from applying %v to %#v, was panic [%v]",
			matcher, input, err)
		return nil
	}
	if result.Matched() {
		t.Errorf("Expected non-matching result from applying %v to %#v, was [%v].  Message: %v",
			result.Matcher(), result.Value(), result, message)
	}
	return result
}

func Test_Matched(t *testing.T) {
	matcher := Matched()
	passResult := NewResultf(true, "pass")
	failResult := NewResultf(false, "fail")
	
	checkResultIsMatching(t, matcher, passResult, "matching")
	checkResultIsNonMatching(t, matcher, failResult, "non-matching")
	checkResultIsNonMatching(t, matcher, nil, "nil")
	checkResultIsNonMatching(t, matcher, "foo", "not a Result")
}

func Test_DidNotMatch(t *testing.T) {
	matcher := DidNotMatch()
	passResult := NewResultf(true, "pass")
	failResult := NewResultf(false, "fail")
	
	checkResultIsNonMatching(t, matcher, passResult, "matching")
	checkResultIsMatching(t, matcher, failResult, "non-matching")
	checkResultIsNonMatching(t, matcher, nil, "nil")
	checkResultIsNonMatching(t, matcher, "foo", "not a Result")
}

func Test_NewMatcher_Func_InterfaceIn_ResultOut(t *testing.T) {
	matcherName := "foo matcher"
	trueMessage := "bar true"
	falseMessage := "baz false"
	function := func(v interface{}) *Result {
		if b, ok := v.(bool); b && ok {
			return NewResultf(true, trueMessage)
		}
		return NewResultf(false, falseMessage)
	}
	matcher := NewMatcherf(function, matcherName)
	if s := fmt.Sprint(matcher); !strings.Contains(s, matcherName) {
		t.Fatalf("String should have contained %v, but was: %v", matcherName, s)
	}
	
	if result := matcher.Match(true); !result.Matched() {
		t.Fatalf("Result should have matched, but was: %v", result)
	} else if s := fmt.Sprint(result); !strings.Contains(s, trueMessage) {
		t.Fatalf("String should have contained %v, but was: %v", trueMessage, s)
	}
	if result := matcher.Match(false);  result.Matched() {
		t.Fatalf("Result should not have matched, but did: %v", result)
	} else if s := fmt.Sprint(result); !strings.Contains(s, falseMessage) {
		t.Fatalf("String should have contained %v, but was: %v", falseMessage, s)
	}
}

func Test_NewMatcher_Func_StringIn_BoolOut(t *testing.T) {
	matcherName := "foo matcher"
	trueMessage := "Matched"
	falseMessage := "Did not match"
	errorMessage := "Could not apply"
	function := func(s string) bool {
		if s == "bar" {
			return true
		}
		return false
	}
	matcher := NewMatcherf(function, matcherName)
	if s := fmt.Sprint(matcher); !strings.Contains(s, matcherName) {
		t.Fatalf("String should have contained %v, but was: %v", matcherName, s)
	}
	
	if result := matcher.Match("bar"); !result.Matched() {
		t.Fatalf("Result should have matched, but was: %v", result)
	} else if s := fmt.Sprint(result); !strings.Contains(s, trueMessage) {
		t.Fatalf("String should have contained %v, but was: %v", trueMessage, s)
	}
	if result := matcher.Match("foo");  result.Matched() {
		t.Fatalf("Result should not have matched, but did: %v", result)
	} else if s := fmt.Sprint(result); !strings.Contains(s, falseMessage) {
		t.Fatalf("String should have contained %v, but was: %v", falseMessage, s)
	}
	if result := matcher.Match(39);  result.Matched() {
		t.Fatalf("Result should not have matched, but did: %v", result)
	} else if s := fmt.Sprint(result); !strings.Contains(s, errorMessage) {
		t.Fatalf("String should have contained %v, but was: %v", errorMessage, s)
	}
}


func Test_NewMatcher_Func_StringDotDotDotIn_BoolStringOut(t *testing.T) {
	matcherName := "foo matcher"
	trueMessage := "bar true"
	falseMessage := "baz false"
	errorMessage := "Could not apply"
	function := func(s... string) (bool, string) {
		if len(s) >= 1 && s[0] == "bar" {
			return true, trueMessage
		}
		return false, falseMessage
	}
	matcher := NewMatcherf(function, matcherName)
	if s := fmt.Sprint(matcher); !strings.Contains(s, matcherName) {
		t.Fatalf("String should have contained %v, but was: %v", matcherName, s)
	}
	
	if result := matcher.Match("bar"); !result.Matched() {
		t.Fatalf("Result should have matched, but was: %v", result)
	} else if s := fmt.Sprint(result); !strings.Contains(s, trueMessage) {
		t.Fatalf("String should have contained %v, but was: %v", trueMessage, s)
	}
	if result := matcher.Match("foo");  result.Matched() {
		t.Fatalf("Result should not have matched, but did: %v", result)
	} else if s := fmt.Sprint(result); !strings.Contains(s, falseMessage) {
		t.Fatalf("String should have contained %v, but was: %v", falseMessage, s)
	}
	if result := matcher.Match(39);  result.Matched() {
		t.Fatalf("Result should not have matched, but did: %v", result)
	} else if s := fmt.Sprint(result); !strings.Contains(s, errorMessage) {
		t.Fatalf("String should have contained %v, but was: %v", errorMessage, s)
	}
}





