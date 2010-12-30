package strings

import (
	"fmt"
	"hamcrest"
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

// Matches strings that begin with the given prefix.
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