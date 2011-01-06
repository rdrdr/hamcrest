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

