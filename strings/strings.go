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

func WithPattern(pattern string) *WithPatternClause {
	return &WithPatternClause{re: regexp.MustCompile(pattern), group: 0 }
}

//
//
//
func (self *WithPatternClause) Group(index int) *WithPatternClause {
	numGroups := self.re.NumSubexp()
	if index < 0 || index > numGroups {
		panic(fmt.Sprintf("Group %v doesn't exist: only %v groups in re %v",
			index, numGroups, self.re))
	}
	return &WithPatternClause{re: self.re, group: index }
}

// Completes a matcher that finds every occurrence of a pattern in the
// given input and applies the matcher to it, only matching if every
// occurrence matches.  For example:
//    i_before_e_except := ToLower(WithPattern(".ei").Each(StartsWith("c")))
// will match:
//    "ceiling receipt"
// but not:
//    "deceiver seizure"
func (self *WithPatternClause) EachMatch(matcher *base.Matcher) *base.Matcher {
	re := self.re
	match := func (s string) *base.Result {
		matches := re.FindAllStringIndex(s, -1)
		if matches == nil {
			return base.NewResultf(true,
				"No occurrences of pattern \"%v\"", re)
		}
		groupOffset := 2 * self.group
		for _, loc := range matches {
			start, end := loc[groupOffset], loc[groupOffset + 1]
			substring := s[start:end]
			result := matcher.Match(substring)
			if !result.Matched() {
				return base.NewResultf(false,
					"did not match substring[%v:%v]=\"%v\" for pattern \"%v\" group %v",
					start, end, substring, re, self.group)
			}
		}
		return base.NewResultf(true,
			"Matched all occurrences of pattern \"%v\"", re)
	}
	return base.NewMatcherf(match, "EachMatch[\"%v\"][%v]", re, matcher)
}

// Completes a matcher that finds every occurrence of a pattern in the
// given input and applies the matcher to it, only matching if at least
// one occurrence matches.  For example:
//    here_kitty := WithPattern(".at").AnyMatch(StartsWith("c"))
// will match:
//    "that cat is phat"
// but not:
//    "Matt spat at a rat"
func (self *WithPatternClause) AnyMatch(matcher *base.Matcher) *base.Matcher {
	re := self.re
	match := func (s string) *base.Result {
		matches := re.FindAllStringIndex(s, -1)
		if matches == nil {
			return base.NewResultf(false, "No occurrences of pattern \"%v\"", re)
		}
		occurrences := 0
		groupOffset := 2 * self.group
		for _, loc := range matches {
			occurrences += 1
			start, end := loc[groupOffset], loc[groupOffset + 1]
			substring := s[start:end]
			result := matcher.Match(substring)
			if result.Matched() {
				return base.NewResultf(true,
					"matched substring[%v:%v]=\"%v\" on pattern \"%v\" group %v",
					start, end, substring, re, self.group)
			}
		}
		return base.NewResultf(false,
			"Matched none of the %v occurrences of pattern \"%v\"",
			occurrences, re)
	}
	return base.NewMatcherf(match, "AnyMatch[\"%v\"][%v]", re, matcher)
}

// Completes a matcher that finds every occurrence of a pattern in
// the given input and applies the matcher to it, only matching if
// there is exactly one occurrence, and the provided matcher matches
// that occurrence.  For example:
//     treasure := WithPattern("x+").TheMatch(ToLen(Equals(1)))
// will match:
//     "..x.."
// but not:
//     "..y.." because there are no occurrences of the pattern
//     "..x..x.." because there are more than one occurrences of the pattern
//     "..xx.." because the ToLen matcher does not match
func (self *WithPatternClause) TheMatch(matcher *base.Matcher) *base.Matcher {
	re := self.re
	match := func (s string) *base.Result {
		matches := re.FindAllStringIndex(s, 2)
		if matches == nil {
			return base.NewResultf(false, "No occurrences of pattern \"%v\"", re)
		}
		if len(matches) > 1 {
			return base.NewResultf(false, "Multiple occurrences of pattern \"%v\"", re)
		}
		loc := matches[0]
		groupOffset := 2 * self.group
		start, end := loc[groupOffset], loc[groupOffset + 1]
		substring := s[start:end]
		result := matcher.Match(substring)
		return base.NewResultf(result.Matched(),
			"Matched substring[%v:%v]=\"%v\" on pattern \"%v\" group %v",
				start, end, substring, re, self.group)
	}
	return base.NewMatcherf(match, "TheMatch[\"%v\"][%v]", re, matcher)
}

