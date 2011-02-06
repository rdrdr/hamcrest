// Copyright 2011 Mick Killianey.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package strings

import (
	"fmt"
	"github.com/rdrdr/hamcrest/base"
	"regexp"
	"strings"
)

// Applies the given matcher to the result of writing the input object's
// to a string by using fmt.Sprintf("%v", object).
func ToString(matcher *base.Matcher) *base.Matcher {
	match := func(actual interface{}) *base.Result {
		if stringer, ok := actual.(fmt.Stringer); ok {
			s := stringer.String()
			result := matcher.Match(s)
			return base.NewResultf(
				result.Matched(), "String() returned %v", s).
				WithCauses(result)
		}
		s :=  fmt.Sprintf("%v", actual)
		result := matcher.Match(s)
		return base.NewResultf(result.Matched(),
			"Not a fmt.Stringer, but prints as %v", s).
			WithCauses(result)
	}
	return base.NewMatcherf(match, "ToString(%v)", matcher)
}


// Applies the given matcher to the result of writing the input object's
// to a string by using fmt.Sprintf("%#v", object).
func ToGoString(matcher *base.Matcher) *base.Matcher {
	match := func(actual interface{}) *base.Result {
		if gostringer, ok := actual.(fmt.GoStringer); ok {
			s := gostringer.GoString()
			result := matcher.Match(s)
			return base.NewResultf(result.Matched(),
				"GoString() returned %v", s).
				WithCauses(result)
		}
		s := fmt.Sprintf("%#v", actual)
		result := matcher.Match(s)
		return base.NewResultf(result.Matched(),
			"Not a fmt.GoStringer, but prints as %v", s).
			WithCauses(result)
	}
	return base.NewMatcherf(match, "ToGoString(%v)", matcher)
}


// Creates a new matcher that applies the given matcher to the result of
// converting an input string to lowercase (using strings.ToLower).
// If the input value is not a string, the matcher fails to match.
func ToLower(matcher *base.Matcher) *base.Matcher {
	match := func(s string) *base.Result {
		lower := strings.ToLower(s)
		result := matcher.Match(lower)
		return base.NewResultf(result.Matched(),
			"ToLower is %v", lower).
			WithCauses(result)
	}
	return base.NewMatcherf(match, "ToLower(%v)", matcher)
}


// Creates a new matcher that applies the given matcher to the result of
// converting an input string to uppercase (using strings.ToUpper).
// If the input value is not a string, the matcher fails to match.
func ToUpper(matcher *base.Matcher) *base.Matcher {
	match := func(s string) *base.Result {
		upper := strings.ToUpper(s)
		result := matcher.Match(upper)
		return base.NewResultf(result.Matched(),
			"ToUpper is %v", upper).
			WithCauses(result)
	}
	return base.NewMatcherf(match, "ToUpper(%v)", matcher)
}

func EqualToIgnoringCase(expected string) *base.Matcher {
	expectedToLower := strings.ToLower(expected)
	match := func(actual string) *base.Result {
		actualToLower := strings.ToLower(actual)
		if actualToLower == expectedToLower {
			return base.NewResultf(true,
				"\"%v\" matches \"%v\" (ignoring case)",
				actual, expected)
		}
		return base.NewResultf(false,
			"\"%v\" differs from \"%v\" (ignoring case)",
			actual, expected)
	}
	return base.NewMatcherf(match, "EqualToIgnoringCase(\"%v\")", expected)
}



// Creates a new matcher that applies the given matcher to the result of
// converting an input string its length. (using the `len()` builtin).
// If the input value is not a string, the matcher fails to match.
func ToLen(matcher *base.Matcher) *base.Matcher {
	match := func(s string) *base.Result {
		length := len(s)
		result := matcher.Match(length)
		return base.NewResultf(result.Matched(),
			"length is %v", length).
			WithCauses(result)
	}
	return base.NewMatcherf(match, "ToLen(%v)", matcher)
}


// Matches strings that begin with the given prefix.
func HasPrefix(prefix string) *base.Matcher {
	maxLength := len(prefix) + 8 // arbitrary extra amount
	match := func (s string) *base.Result {
		continued := ""
		if len(s) > maxLength {
			s, continued = s[:maxLength], "..."
		}
		if strings.HasPrefix(s, prefix) {
			return base.NewResultf(true, 
				"\"%v%v\" starts with \"%v\"", s, continued, prefix)
		}
		return base.NewResultf(false,
			"\"%v%v\" does not start with \"%v\"", s, continued, prefix)
	}
	return base.NewMatcherf(match, "HasPrefix(\"%v\")", prefix)
}

// Matches strings that end with the given prefix.
func HasSuffix(suffix string) *base.Matcher {
	maxLength := len(suffix) + 8 // arbitrary extra amount
	match := func (s string) *base.Result {
		continued := ""
		if len(s) > maxLength {
			continued, s = "...", s[len(s) - maxLength:]
		}
		if strings.HasSuffix(s, suffix) {
			return base.NewResultf(true,
				"\"%v%v\" ends with \"%v\"", s, continued, suffix)
		}
		return base.NewResultf(false,
			"\"%v%v\" does not end with \"%v\"", s, continued, suffix)
	}
	return base.NewMatcherf(match, "HasSuffix(\"%v\")", suffix)
}


// Matches strings that contain the given substring.
func Contains(substring string) *base.Matcher {
	match := func (s string) *base.Result {
		extra := 8
		if foundStart := strings.Index(s, substring); foundStart >= 0 {
			foundEnd := foundStart + len(substring)
			start, end := foundStart - extra, foundEnd + extra
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
			return base.NewResultf(true,
				"substring \"%v\" appears in \"%v%v[%v]%v%v\"", substring,
				prefix, s[start:foundStart], substring, s[foundEnd:end], suffix)
		}
		return base.NewResultf(false,
			"substring \"%v\" does not appear in \"%v\"",
			substring, s)
	}
	return base.NewMatcherf(match, "Contains(\"%v\")", substring)
}

// Matches strings that contain the given regexp pattern, using
// the same syntax as the standard regexp package.
func HasPattern(pattern string) *base.Matcher {
	re := regexp.MustCompile(pattern)
	match := func (s string) *base.Result {
		if found := re.FindStringIndex(s); found != nil {
			start, end := found[0], found[1]
			return base.NewResultf(true,
				"pattern \"%v\" matched substring[%v:%v]=\"%v\"",
				pattern, start, end, s[start:end])
		}
		return base.NewResultf(false,
			"pattern \"%v\" not found in \"%v\"", pattern, s)
	}
	return base.NewMatcherf(match, "HasPattern[\"%v\"]", pattern)
}

type WithPatternClause struct {
	re *regexp.Regexp
	group int
}

// Returns a short-circuiting function that applies the matcher to each
// occurrence of the pattern in an input string, until a failing pattern
// is found (in which case the output matcher fails to match) or all
// matching instances are exhausted (in which case the output matcher
// successfully matches).
//
// For example:
//    EachPattern("q.")(EqualTo("qu"))
// would match:
//    "quick quack mq" (two matches of "q.", both equal to "qu")
// but not:
//    "quick qs for mq" (two matches of "q.", second is not "qu")
func EachPattern(pattern string) func(matcher *base.Matcher) *base.Matcher  {
	re := regexp.MustCompile(pattern)
	return func(matcher *base.Matcher) *base.Matcher {
		match := func(s string) *base.Result {
			matches := re.FindAllStringIndex(s, -1)
			if matches == nil {
				return base.NewResultf(true,
					"No occurrences of pattern \"%v\"", pattern)
			}
			for index, loc := range matches {
				start, end := loc[0], loc[1]
				substring := s[start:end]
				result := matcher.Match(substring)
				if !result.Matched() {
					return base.NewResultf(false,
						"did not match substring[%v:%v]=\"%v\", occurrence #%v (of %v) of pattern \"%v\"",
						start, end, substring, index+1, len(matches), pattern)
				}
			}
			return base.NewResultf(true,
				"Matched every occurrence (all %v) of pattern \"%v\"",
				len(matches), pattern)
		}
		return base.NewMatcherf(match,
			"EachPattern[\"%v\"][%v]", pattern, matcher)
	}
}

// Returns a short-circuiting function that applies the matcher to each
// occurrence of the pattern in an input string, until a failing pattern
// is found (in which case the output matcher fails to match) or all
// matching instances are exhausted (in which case the output matcher
// successfully matches).
//
// For example:
//    EachPattern("q.")(EqualTo("qu"))
// would match:
//    "quick quack mq" (two matches of "q.", both equal to "qu")
// but not:
//    "quick qs for mq" (two matches of "q.", second is not "qu")
func EachPatternGroup(pattern string, group int) func(matcher *base.Matcher) *base.Matcher  {
	re := regexp.MustCompile(pattern)
	if num := re.NumSubexp(); num < group {
		println("Illegal group #", group, ": there are only ", num, "groups.") 
		panic("Group index out of bounds.")
	}
	return func(matcher *base.Matcher) *base.Matcher {
		match := func(s string) *base.Result {
			matches := re.FindAllStringSubmatchIndex(s, -1)
			if matches == nil {
				return base.NewResultf(true,
					"No occurrences of pattern \"%v\"", pattern)
			}
			for index, loc := range matches {
				start, end := loc[2*group], loc[2*group+1]
				substring := s[start:end]
				result := matcher.Match(substring)
				if !result.Matched() {
					prefix, suffix := s[loc[0]:start], s[end:loc[1]]
					return base.NewResultf(false,
						"did not match substring[%v:%v]=\"%v[%v]%v\", occurrence #%v (of %v) of pattern \"%v\"",
						start, end, prefix, substring, suffix, index+1, len(matches), pattern)
				}
			}
			return base.NewResultf(true,
				"Matched every occurrence (all %v) of pattern \"%v\", group %v",
				len(matches), pattern, group)
		}
		return base.NewMatcherf(match,
			"EachPatternGroup[\"%v\", %v][%v]", pattern, group, matcher)
	}
}


// Returns a short-circuiting function that applies the matcher to each
// occurrence of the pattern in an input string, until a matching pattern
// is found (in which case the matcher successfully matches) or all
// matching instances are exhausted (in which case the output matcher
// fails to match).
//
// For example:
//    AnyPattern("x.")(EqualTo("xy"))
// would match:
//    "six sax are sexy" (three matches of "x.", third is "xy")
// but not:
//    "pox pix are pixelated" (three matches of "x.", none is "xy")
func AnyPattern(pattern string) func(matcher *base.Matcher) *base.Matcher  {
	re := regexp.MustCompile(pattern)
	return func(matcher *base.Matcher) *base.Matcher {
		match := func(s string) *base.Result {
			matches := re.FindAllStringIndex(s, -1)
			if matches == nil {
				return base.NewResultf(false,
					"No occurrences of pattern \"%v\"", pattern)
			}
			
			for index, loc := range matches {
				start, end := loc[0], loc[1]
				substring := s[start:end]
				result := matcher.Match(substring)
				if result.Matched() {
					return base.NewResultf(true,
						"matched substring[%v:%v]=\"%v\", occurrence #%v (of %v) of pattern \"%v\"",
						start, end, substring, index+1, len(matches), pattern)
				}
			}
			return base.NewResultf(false,
				"Did not match any occurrence (of %v) of pattern \"%v\"",
				len(matches), pattern)
		}
		return base.NewMatcherf(match,
			"AnyPattern[\"%v\"][%v]", pattern, matcher)
	}
}

// Returns a function that applies the matcher to the first occurrence of
// the pattern in an input string.
//
// For example:
//    FirstInstanceOf("h.s")(EqualTo("his"))
// would match:
//    "hers and his" (because the first instance of "h.s" is equal to "his")
// but none of:
//    "just hers" (no instances of "h.s")
//    "has chisel" (the first instance of "h.s" is not "his")
func FirstInstanceOf(pattern string) func(matcher *base.Matcher) *base.Matcher  {
	re := regexp.MustCompile(pattern)
	return func(matcher *base.Matcher) *base.Matcher {
		match := func(s string) *base.Result {
			matches := re.FindStringIndex(s)
			if matches == nil {
				return base.NewResultf(false,
					"No occurrences of pattern \"%v\"", pattern)
			}
			start, end := matches[0], matches[1]
			substring := s[start:end]
			result := matcher.Match(substring)
			return base.NewResultf(result.Matched(),
				"Found substring[%v:%v]=\"%v\" for pattern \"%v\"",
				start, end, substring, pattern).WithCauses(result)
		}
		return base.NewMatcherf(match,
			"FirstInstanceOf[\"%v\"][%v]", pattern, matcher)
	}
}

