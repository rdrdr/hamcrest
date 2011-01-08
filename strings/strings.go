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
	description := hamcrest.NewDescription("ToString(%v)", matcher)
	match := func(actual interface{}) *hamcrest.Result {
		var s string
		var why *hamcrest.Description
		if stringer, ok := actual.(fmt.Stringer); ok {
			s = stringer.String()
			why = hamcrest.NewDescription("String() returned %v", s)
		} else {
			s = fmt.Sprintf("%v", actual)
			why = hamcrest.NewDescription("Not a fmt.Stringer, but prints as %v", s)
		}
		result := matcher.Match(s)
		return hamcrest.NewResult(result.Matched(), why).WithCauses(result)
	}
	return hamcrest.NewMatcher(description, match)
}


// Applies the given matcher to the result of writing the input object's
// to a string by using fmt.Sprintf("%#v", object).
func ToGoString(matcher *hamcrest.Matcher) *hamcrest.Matcher {
	description := hamcrest.NewDescription("ToGoString(%v)", matcher)
	match := func(actual interface{}) *hamcrest.Result {
		var s string
		var why *hamcrest.Description
		if gostringer, ok := actual.(fmt.GoStringer); ok {
			s = gostringer.GoString()
			why = hamcrest.NewDescription("GoString() returned %v", s)
		} else {
			s = fmt.Sprintf("%#v", actual)
			why = hamcrest.NewDescription("Not a fmt.GoStringer, but prints as %v", s)
		}
		result := matcher.Match(s)
		if result.Matched() {
			return hamcrest.NewResult(true, why).WithCauses(result)
		}
		return hamcrest.NewResult(false, why).WithCauses(result)
	}
	return hamcrest.NewMatcher(description, match)
}


// Creates a new matcher (as hamcrest.NewMatcher), that only matches
// strings.  Note that its match function has an (appropriately) different
// signature from the match function required by hamcrest.NewMatcher.
//
// If the input value is not a string, the matcher will fail to match.
func NewStringMatcher(description *hamcrest.Description, matchString func(s string) *hamcrest.Result) *hamcrest.Matcher {
	match := func (v interface{}) *hamcrest.Result {
		if s, ok := v.(string); ok {
			return matchString(s)
		}
		why := hamcrest.NewDescription("%v is a %T, not a string", v, v)
		return hamcrest.NewResult(false, why)
	}
	return hamcrest.NewMatcher(description, match)
}


// Creates a new matcher that applies the given matcher to the result of
// converting an input string to lowercase (using strings.ToLower).
// If the input value is not a string, the matcher fails to match.
func ToLower(matcher *hamcrest.Matcher) *hamcrest.Matcher {
	description := hamcrest.NewDescription("ToLower(%v)", matcher)
	match := func(s string) *hamcrest.Result {
		lower := strings.ToLower(s)
		why := hamcrest.NewDescription("to lower is %v", lower)
		result := matcher.Match(lower)
		return hamcrest.NewResult(result.Matched(), why).WithCauses(result)
	}
	return NewStringMatcher(description, match)
}


// Creates a new matcher that applies the given matcher to the result of
// converting an input string to uppercase (using strings.ToUpper).
// If the input value is not a string, the matcher fails to match.
func ToUpper(matcher *hamcrest.Matcher) *hamcrest.Matcher {
	description := hamcrest.NewDescription("ToUpper(%v)", matcher)
	match := func(s string) *hamcrest.Result {
		upper := strings.ToUpper(s)
		why := hamcrest.NewDescription("to upper is %v", upper)
		result := matcher.Match(upper)
		return hamcrest.NewResult(result.Matched(), why).WithCauses(result)
	}
	return NewStringMatcher(description, match)
}

func EqualToIgnoringCase(expected string) *hamcrest.Matcher {
	description := hamcrest.NewDescription("EqualToIgnoringCase(\"%v\")", expected)
	expectedToLower := strings.ToLower(expected)
	match := func(actual string) *hamcrest.Result {
		actualToLower := strings.ToLower(actual)
		if actualToLower == expectedToLower {
			why := hamcrest.NewDescription("\"%v\" matches \"%v\" (ignoring case)",
				actual, expected)
			return hamcrest.NewResult(true, why)
		}
		why := hamcrest.NewDescription("\"%v\" differs from \"%v\" (ignoring case)",
			actual, expected)
		return hamcrest.NewResult(false, why)
	}
	return NewStringMatcher(description, match)
}



// Creates a new matcher that applies the given matcher to the result of
// converting an input string its length. (using the `len()` builtin).
// If the input value is not a string, the matcher fails to match.
func ToLen(matcher *hamcrest.Matcher) *hamcrest.Matcher {
	description := hamcrest.NewDescription("ToLen(%v)", matcher)
	match := func(s string) *hamcrest.Result {
		length := len(s)
		why := hamcrest.NewDescription("length is %v", length)
		result := matcher.Match(length)
		return hamcrest.NewResult(result.Matched(), why).WithCauses(result)
	}
	return NewStringMatcher(description, match)
}


// Matches strings that begin with the given prefix.
func HasPrefix(prefix string) *hamcrest.Matcher {
	description := hamcrest.NewDescription("HasPrefix(\"%v\")", prefix)
	maxLength := len(prefix) + 8 // arbitrary extra amount
	match := func (s string) *hamcrest.Result {
		continued := ""
		if len(s) > maxLength {
			s, continued = s[:maxLength], "..."
		}
		if strings.HasPrefix(s, prefix) {
			why := hamcrest.NewDescription(
				"\"%v%v\" starts with \"%v\"", s, continued, prefix)
			return hamcrest.NewResult(true, why)
		}
		why := hamcrest.NewDescription(
			"\"%v%v\" does not start with \"%v\"", s, continued, prefix)
		return hamcrest.NewResult(false, why)
	}
	return NewStringMatcher(description, match)
}

// Matches strings that end with the given prefix.
func HasSuffix(suffix string) *hamcrest.Matcher {
	description := hamcrest.NewDescription("HasSuffix(\"%v\")", suffix)
	maxLength := len(suffix) + 8 // arbitrary extra amount
	match := func (s string) *hamcrest.Result {
		continued := ""
		if len(s) > maxLength {
			continued, s = "...", s[len(s) - maxLength:]
		}
		if strings.HasSuffix(s, suffix) {
			why := hamcrest.NewDescription(
				"\"%v%v\" ends with \"%v\"", s, continued, suffix)
			return hamcrest.NewResult(true, why)
		}
		why := hamcrest.NewDescription(
				"\"%v%v\" does not end with \"%v\"", s, continued, suffix)
		return hamcrest.NewResult(false, why)
	}
	return NewStringMatcher(description, match)
}


// Matches strings that contain the given substring.
func Contains(substring string) *hamcrest.Matcher {
	description := hamcrest.NewDescription("Contains(\"%v\")", substring)
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
			why := hamcrest.NewDescription("substring \"%v\" appears in \"%v%v%v\"",
				substring, prefix, s[start:end], suffix)
			return hamcrest.NewResult(true, why)
		}
		why := hamcrest.NewDescription(
			"substring \"%v\" does not appear in \"%v\"", substring, s)
		return hamcrest.NewResult(false, why)
	}
	return NewStringMatcher(description, match)
}

// Matches strings that contain the given regexp pattern, using
// the same syntax as the standard regexp package.
func HasPattern(pattern string) *hamcrest.Matcher {
	re := regexp.MustCompile(pattern)
	description := hamcrest.NewDescription("HasPattern[\"%v\"]", pattern)
	match := func (s string) *hamcrest.Result {
		if found := re.FindStringIndex(s); found != nil {
			start, end := found[0], found[1]
			why := hamcrest.NewDescription(
				"pattern \"%v\" matched substring[%v:%v]=\"%v\"",
					pattern, start, end, s[start:end])
			return hamcrest.NewResult(true, why)
		}
		why := hamcrest.NewDescription(
			"pattern \"%v\" not found in \"%v\"", pattern, s)
		return hamcrest.NewResult(false, why)
	}
	return NewStringMatcher(description, match)
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
	description := hamcrest.NewDescription("Each[\"%v\"][%v]", re, matcher)
	match := func (s string) *hamcrest.Result {
		matches := re.FindAllStringIndex(s, -1)
		if matches == nil {
			why := hamcrest.NewDescription("No occurrences of pattern")
			return hamcrest.NewResult(true, why)
		}
		for _, loc := range matches {
			start, end := loc[0], loc[1]
			substring := s[start:end]
			result := matcher.Match(substring)
			if !result.Matched() {
				why := hamcrest.NewDescription("did not match substring[%v:%v]=%v",
					start, end, substring)
				return hamcrest.NewResult(false, why)
			}
		}
		why := hamcrest.NewDescription("Matched all occurrences of pattern")
		return hamcrest.NewResult(true, why)
	}
	return NewStringMatcher(description, match)
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
	description := hamcrest.NewDescription("Any[\"%v\"][%v]", re, matcher)
	match := func (s string) *hamcrest.Result {
		matches := re.FindAllStringIndex(s, -1)
		if matches == nil {
			why := hamcrest.NewDescription("No occurrences of pattern")
			return hamcrest.NewResult(false, why)
		}
		occurrences := 0
		for _, loc := range matches {
			occurrences += 1
			start, end := loc[0], loc[1]
			substring := s[start:end]
			result := matcher.Match(substring)
			if result.Matched() {
				why := hamcrest.NewDescription("matched substring[%v:%v]=%v",
					start, end, substring)
				return hamcrest.NewResult(true, why)
			}
		}
		why := hamcrest.NewDescription(
			"Matched none of the %v occurrences of pattern", occurrences)
		return hamcrest.NewResult(false, why)
	}
	return NewStringMatcher(description, match)
}

