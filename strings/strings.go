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
	description := hamcrest.NewDescription("ToString[%v]", matcher)
	match := func(actual interface{}) *hamcrest.Result {
		s := fmt.Sprintf("%v", actual)
		result := matcher.Match(s)
		because := hamcrest.NewDescription("String() was %v", s)
		if result.Matched() {
			return hamcrest.NewResult(true, because).WithCauses(result)
		}
		return hamcrest.NewResult(false, because).WithCauses(result)
	}
	return hamcrest.NewMatcher(description, match)
}


// Applies the given matcher to the result of writing the input object's
// to a string by using fmt.Sprintf("%#v", object).
func ToGoString(matcher *hamcrest.Matcher) *hamcrest.Matcher {
	description := hamcrest.NewDescription("ToGoString[%v]", matcher)
	match := func(actual interface{}) *hamcrest.Result {
		s := fmt.Sprintf("%#v", actual)
		result := matcher.Match(s)
		because := hamcrest.NewDescription("GoString() was %#v", s)
		if result.Matched() {
			return hamcrest.NewResult(true, because).WithCauses(result)
		}
		return hamcrest.NewResult(false, because).WithCauses(result)
	}
	return hamcrest.NewMatcher(description, match)
}

// Matches strings that begin with the given prefix.
func HasPrefix(prefix string) *hamcrest.Matcher {
	description := hamcrest.NewDescription("HasPrefix[\"%v\"]", prefix)
	match := func (actual interface{}) *hamcrest.Result {
		if s, ok := actual.(string); ok {
			if strings.HasPrefix(s, prefix) {
				because := hamcrest.NewDescription(
						"[%v] starts with [%v]", actual, prefix)
				return hamcrest.NewResult(true, because)
			}
			because := hamcrest.NewDescription(
					"[%v] does not start with [%v]", actual, prefix)
			return hamcrest.NewResult(false, because)
		}
		because := hamcrest.NewDescription("[%v] is not a string", actual)
		return hamcrest.NewResult(false, because)
	}
	return hamcrest.NewMatcher(description, match)
}

// Matches strings that end with the given prefix.
func HasSuffix(suffix string) *hamcrest.Matcher {
	description := hamcrest.NewDescription("HasSuffix[\"%v\"]", suffix)
	match := func (actual interface{}) *hamcrest.Result {
		if s, ok := actual.(string); ok {
			if strings.HasSuffix(s, suffix) {
				because := hamcrest.NewDescription(
						"[%v] ends with [%v]", actual, suffix)
				return hamcrest.NewResult(true, because)
			}
			because := hamcrest.NewDescription(
					"[%v] does not end with [%v]", actual, suffix)
			return hamcrest.NewResult(false, because)
		}
		because := hamcrest.NewDescription("[%v] is not a string", actual)
		return hamcrest.NewResult(false, because)
	}
	return hamcrest.NewMatcher(description, match)
}


// Matches strings that contain the given substring.
func Contains(substring string) *hamcrest.Matcher {
	description := hamcrest.NewDescription("Contains[\"%v\"]", substring)
	match := func (actual interface{}) *hamcrest.Result {
		if s, ok := actual.(string); ok {
			if strings.Contains(s, substring) {
				because := hamcrest.NewDescription("contains substring '%v'", substring)
				return hamcrest.NewResult(true, because)
			}
			because := hamcrest.NewDescription("does not contain substring '%v'", substring)
			return hamcrest.NewResult(false, because)
		}
		because := hamcrest.NewDescription("[%v] is not a string", actual)
		return hamcrest.NewResult(false, because)
	}
	return hamcrest.NewMatcher(description, match)
}

// Matches strings that contain the given regexp pattern, using
// the same syntax as the standard regexp package.
func HasPattern(pattern string) *hamcrest.Matcher {
	re := regexp.MustCompile(pattern)
	description := hamcrest.NewDescription("HasPattern[\"%v\"]", pattern)
	match := func (actual interface{}) *hamcrest.Result {
		if s, ok := actual.(string); ok {
			if found := re.FindStringIndex(s); found != nil {
				start, end := found[0], found[1]
				because := hamcrest.NewDescription(
					"pattern '%v' found at substring[%v:%v]='%v'",
						pattern, start, end, s[start:end])
				return hamcrest.NewResult(true, because)
			}
			because := hamcrest.NewDescription(
				"pattern '%v' not found in '%v'", pattern, actual)
			return hamcrest.NewResult(false, because)
		} else if b, ok := actual.([]byte); ok {
			if found := re.FindIndex(b); found != nil {
				start, end := found[0], found[1]
				because := hamcrest.NewDescription(
					"pattern '%v' found at subsequence[%v:%v]='%v'",
						pattern, start, end, b[start:end])
				return hamcrest.NewResult(true, because)
			}
			because := hamcrest.NewDescription(
				"pattern '%v' not found in sequence '%v'", pattern, actual)
			return hamcrest.NewResult(false, because)
		}
		because := hamcrest.NewDescription("[%v] is not a string or byte array", actual)
		return hamcrest.NewResult(false, because)
	}
	return hamcrest.NewMatcher(description, match)
}

