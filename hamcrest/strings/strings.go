// Copyright 2011 Mick Killianey.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package strings

import (
	"fmt"
	"hamcrest"
	"regexp"
	"strings"
)

// Applies the given matcher to the result of writing the input object's
// to a string by using fmt.Sprintf("%v", object).
func ToString(matcher *hamcrest.Matcher) *hamcrest.Matcher {
	match := func(actual interface{}) *hamcrest.Result {
		if stringer, ok := actual.(fmt.Stringer); ok {
			s := stringer.String()
			result := matcher.Match(s)
			return hamcrest.NewResultf(
				result.Matched(), "String() returned %v", s).
				WithCauses(result)
		}
		s :=  fmt.Sprintf("%v", actual)
		result := matcher.Match(s)
		return hamcrest.NewResultf(result.Matched(),
			"Not a fmt.Stringer, but prints as %v", s).
			WithCauses(result)
	}
	return hamcrest.NewMatcherf(match, "ToString(%v)", matcher)
}


// Applies the given matcher to the result of writing the input object's
// to a string by using fmt.Sprintf("%#v", object).
func ToGoString(matcher *hamcrest.Matcher) *hamcrest.Matcher {
	match := func(actual interface{}) *hamcrest.Result {
		if gostringer, ok := actual.(fmt.GoStringer); ok {
			s := gostringer.GoString()
			result := matcher.Match(s)
			return hamcrest.NewResultf(result.Matched(),
				"GoString() returned %v", s).
				WithCauses(result)
		}
		s := fmt.Sprintf("%#v", actual)
		result := matcher.Match(s)
		return hamcrest.NewResultf(result.Matched(),
			"Not a fmt.GoStringer, but prints as %v", s).
			WithCauses(result)
	}
	return hamcrest.NewMatcherf(match, "ToGoString(%v)", matcher)
}


// Creates a new matcher (as hamcrest.NewMatcher), that only matches
// strings.  Note that its match function has an (appropriately) different
// signature from the match function required by hamcrest.NewMatcher.
//
// If the input value is not a string, the matcher will fail to match.
func NewStringMatcher(match func(s string) *hamcrest.Result, description hamcrest.SelfDescribing) *hamcrest.Matcher {
	return hamcrest.NewMatcher(func (v interface{}) *hamcrest.Result {
		if s, ok := v.(string); ok {
			return match(s)
		}
		return hamcrest.NewResultf(false, "was a %T, not a string", v)
	}, description)
}

// Creates a new matcher (as hamcrest.NewMatcher), that only matches
// strings.  Note that its match function has an (appropriately) different
// signature from the match function required by hamcrest.NewMatcher.
//
// If the input value is not a string, the matcher will fail to match.
func NewStringMatcherf(match func(s string) *hamcrest.Result, format string, args...interface{}) *hamcrest.Matcher {
	return NewStringMatcher(match, hamcrest.Description(format, args...))
}

// Creates a new matcher that applies the given matcher to the result of
// converting an input string to lowercase (using strings.ToLower).
// If the input value is not a string, the matcher fails to match.
func ToLower(matcher *hamcrest.Matcher) *hamcrest.Matcher {
	match := func(s string) *hamcrest.Result {
		lower := strings.ToLower(s)
		result := matcher.Match(lower)
		return hamcrest.NewResultf(result.Matched(),
			"ToLower is %v", lower).
			WithCauses(result)
	}
	return NewStringMatcherf(match, "ToLower(%v)", matcher)
}


// Creates a new matcher that applies the given matcher to the result of
// converting an input string to uppercase (using strings.ToUpper).
// If the input value is not a string, the matcher fails to match.
func ToUpper(matcher *hamcrest.Matcher) *hamcrest.Matcher {
	match := func(s string) *hamcrest.Result {
		upper := strings.ToUpper(s)
		result := matcher.Match(upper)
		return hamcrest.NewResultf(result.Matched(),
			"ToUpper is %v", upper).
			WithCauses(result)
	}
	return NewStringMatcherf(match, "ToUpper(%v)", matcher)
}

func EqualToIgnoringCase(expected string) *hamcrest.Matcher {
	expectedToLower := strings.ToLower(expected)
	match := func(actual string) *hamcrest.Result {
		actualToLower := strings.ToLower(actual)
		if actualToLower == expectedToLower {
			return hamcrest.NewResultf(true,
				"\"%v\" matches \"%v\" (ignoring case)",
				actual, expected)
		}
		return hamcrest.NewResultf(false,
			"\"%v\" differs from \"%v\" (ignoring case)",
			actual, expected)
	}
	return NewStringMatcherf(match, "EqualToIgnoringCase(\"%v\")", expected)
}



// Creates a new matcher that applies the given matcher to the result of
// converting an input string its length. (using the `len()` builtin).
// If the input value is not a string, the matcher fails to match.
func ToLen(matcher *hamcrest.Matcher) *hamcrest.Matcher {
	match := func(s string) *hamcrest.Result {
		length := len(s)
		result := matcher.Match(length)
		return hamcrest.NewResultf(result.Matched(),
			"length is %v", length).
			WithCauses(result)
	}
	return NewStringMatcherf(match, "ToLen(%v)", matcher)
}


// Matches strings that begin with the given prefix.
func HasPrefix(prefix string) *hamcrest.Matcher {
	maxLength := len(prefix) + 8 // arbitrary extra amount
	match := func (s string) *hamcrest.Result {
		continued := ""
		if len(s) > maxLength {
			s, continued = s[:maxLength], "..."
		}
		if strings.HasPrefix(s, prefix) {
			return hamcrest.NewResultf(true, 
				"\"%v%v\" starts with \"%v\"", s, continued, prefix)
		}
		return hamcrest.NewResultf(false,
			"\"%v%v\" does not start with \"%v\"", s, continued, prefix)
	}
	return NewStringMatcherf(match, "HasPrefix(\"%v\")", prefix)
}

// Matches strings that end with the given prefix.
func HasSuffix(suffix string) *hamcrest.Matcher {
	maxLength := len(suffix) + 8 // arbitrary extra amount
	match := func (s string) *hamcrest.Result {
		continued := ""
		if len(s) > maxLength {
			continued, s = "...", s[len(s) - maxLength:]
		}
		if strings.HasSuffix(s, suffix) {
			return hamcrest.NewResultf(true,
				"\"%v%v\" ends with \"%v\"", s, continued, suffix)
		}
		return hamcrest.NewResultf(false,
			"\"%v%v\" does not end with \"%v\"", s, continued, suffix)
	}
	return NewStringMatcherf(match, "HasSuffix(\"%v\")", suffix)
}


// Matches strings that contain the given substring.
func Contains(substring string) *hamcrest.Matcher {
	match := func (s string) *hamcrest.Result {
		extra := 8
		if found := strings.Index(s, substring); found >= 0 {
			start, end := found - extra, found + len(substring) + extra
			prefix, suffix := "", ""
			if start <= 0 {
				start = 0
			} else {
				prefix = "..."
			}
			if end >= len(s) {
				end = len(s)
			} else {
				suffix = "..."
			}
			return hamcrest.NewResultf(true,
				"substring \"%v\" appears in \"%v%v%v\"",
				substring, prefix, s[start:end], suffix)
		}
		return hamcrest.NewResultf(false,
			"substring \"%v\" does not appear in \"%v\"",
			substring, s)
	}
	return NewStringMatcherf(match, "Contains(\"%v\")", substring)
}

// Matches strings that contain the given regexp pattern, using
// the same syntax as the standard regexp package.
func HasPattern(pattern string) *hamcrest.Matcher {
	re := regexp.MustCompile(pattern)
	match := func (s string) *hamcrest.Result {
		if found := re.FindStringIndex(s); found != nil {
			start, end := found[0], found[1]
			return hamcrest.NewResultf(true,
				"pattern \"%v\" matched substring[%v:%v]=\"%v\"",
				pattern, start, end, s[start:end])
		}
		return hamcrest.NewResultf(false,
			"pattern \"%v\" not found in \"%v\"", pattern, s)
	}
	return NewStringMatcherf(match, "HasPattern[\"%v\"]", pattern)
}

type ExtractPatternClause struct {
	re *regexp.Regexp
}

func ExtractPattern(pattern string) *ExtractPatternClause {
	return &ExtractPatternClause{ re : regexp.MustCompile(pattern) }
}

// Completes a matcher that finds every occurrence of a pattern in the
// given input and applies the matcher to it, only matching if every
// occurrence matches.  For example:
//    i_before_e_except := ToLower(ExtractPattern(".ei").Each(StartsWith("c")))
// will match:
//    "ceiling receipt"
// but not:
//    "deceiver seizure"
func (self *ExtractPatternClause) Each(matcher *hamcrest.Matcher) *hamcrest.Matcher {
	re := self.re
	match := func (s string) *hamcrest.Result {
		matches := re.FindAllStringIndex(s, -1)
		if matches == nil {
			return hamcrest.NewResultf(true,
				"No occurrences of pattern \"%v\"", re)
		}
		for _, loc := range matches {
			start, end := loc[0], loc[1]
			substring := s[start:end]
			result := matcher.Match(substring)
			if !result.Matched() {
				return hamcrest.NewResultf(false,
					"did not match substring[%v:%v]=\"%v\" for pattern \"%v\"",
					start, end, substring, re)
			}
		}
		return hamcrest.NewResultf(true,
			"Matched all occurrences of pattern \"%v\"", re)
	}
	return NewStringMatcherf(match, "Each[\"%v\"][%v]", re, matcher)
}

// Completes a matcher that finds every occurrence of a pattern in the
// given input and applies the matcher to it, only matching if at least
// one occurrence matches.  For example:
//    here_kitty := ExtractPattern(".at").Any(StartsWith("c"))
// will match:
//    "that cat is phat"
// but not:
//    "Matt spat at a rat"
func (self *ExtractPatternClause) Any(matcher *hamcrest.Matcher) *hamcrest.Matcher {
	re := self.re
	match := func (s string) *hamcrest.Result {
		matches := re.FindAllStringIndex(s, -1)
		if matches == nil {
			return hamcrest.NewResultf(false, "No occurrences of pattern \"%v\"", re)
		}
		occurrences := 0
		for _, loc := range matches {
			occurrences += 1
			start, end := loc[0], loc[1]
			substring := s[start:end]
			result := matcher.Match(substring)
			if result.Matched() {
				return hamcrest.NewResultf(true,
					"matched substring[%v:%v]=\"%v\" on pattern \"%v\"",
					start, end, substring, re)
			}
		}
		return hamcrest.NewResultf(false,
			"Matched none of the %v occurrences of pattern \"%v\"",
			occurrences, re)
	}
	return NewStringMatcherf(match, "Any[\"%v\"][%v]", re, matcher)
}

