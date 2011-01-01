package hamcrest

import (
	"reflect"
)

// Returns a Matcher that matches any input value.
func Anything() *Matcher {
	return _Anything
}
var _Anything *Matcher // singleton
func init() {
	description := NewDescription("always matches")
	match := func (actual interface{}) *Result {
		return NewResult(true, description)
	}
	_Anything = NewMatcher(NewDescription("matches anything"), match)
}


// Returns a Matcher that matches the boolean value true.
func True() *Matcher {
	return _True
}
var _True *Matcher // singleton
func init() {
	trueDescription := NewDescription("was true")
	falseDescription := NewDescription("was not true")
	errorDescription := NewDescription("was not bool")
	match := func (actual interface{}) *Result {
		if b, ok := actual.(bool); ok {
			if b {
				return NewResult(true, trueDescription)
			}
			return NewResult(false, falseDescription)
		}
		return NewResult(false, errorDescription)
	}
	_True = NewMatcher(NewDescription("matches true"), match)
}

// Returns a Matcher that matches the boolean value false.
func False() *Matcher {
	return _False
}
var _False *Matcher // singleton
func init() {
	trueDescription := NewDescription("was not false")
	falseDescription := NewDescription("was false")
	errorDescription := NewDescription("was not bool")
	match := func (actual interface{}) *Result {
		if b, ok := actual.(bool); ok {
			if b {
				return NewResult(false, trueDescription)
			}
			return NewResult(true, falseDescription)
		}
		return NewResult(false, errorDescription)
	}
	_False = NewMatcher(NewDescription("matches false"), match)
}


// Returns a Matcher that decorates another matcher and only matches
// when the underlying matcher does not match (and vice versa).
func Not(matcher *Matcher) *Matcher {
	match := func (actual interface{}) *Result {
		result := matcher.Match(actual)
		description := NewDescription("Not[%v]", result.description)
		return NewResult(!result.matched, description).WithCauses(result)
	}
	return NewMatcher(NewDescription("Not[%v]", matcher), match)
}

// Returns a Matcher that decorates another matcher.
func Is(matcher *Matcher) *Matcher {
	match := func (actual interface{}) *Result {
		result := matcher.Match(actual)
		description := NewDescription("Is[%v]", result.description)
		return NewResult(result.matched, description).WithCauses(result)
	}
	return NewMatcher(NewDescription("Is[%v]", matcher), match)
}


// Helper function for Nil/NonNil
func _detectNil(actual interface{}) bool {
	if actual == nil {
		return true
	}
	if ptr, ok := reflect.NewValue(actual).(*reflect.PtrValue); ok {
		return ptr.IsNil()
	}
	return false
}
type _CanAskIsNil interface { IsNil() bool }
type _CanAskAddr interface { Addr() uintptr }


// Returns a Matcher that matches if the actual value is nil
// or the nil value of its type.  (Note that this is *not*
// equivalent to EqualTo(nil).)
func Nil() *Matcher {
	return _Nil
}
var _Nil *Matcher // singleton
func init() {
	wasNilDescription := NewDescription("was nil")
	match := func (actual interface{}) *Result {
		if _detectNil(actual) {
			return NewResult(true, wasNilDescription)
		}
		return NewResult(false, NewDescription("[%#v] was not nil", actual))
	}
	_Nil = NewMatcher(NewDescription("matches nil"), match)
}

// Returns a Matcher that matches if the actual value is 
// neither nil nor the nil value of its type.  (Note that
// this is *not* equivalent to Not(EqualTo(nil)).)
func NonNil() *Matcher {
	return _NonNil
}
var _NonNil *Matcher
func init() {
	wasNilDescription := NewDescription("was nil")
	match := func (actual interface{}) *Result {
		if _detectNil(actual) {
			return NewResult(false, wasNilDescription)
		}
		return NewResult(true, NewDescription("[%#v] was not nil", actual))
	}
	_NonNil = NewMatcher(NewDescription("matches non-nil"), match)
}


// Returns a Matcher that checks if the actual value is (shallowly)
// equal to the given expected value, as if compared using ==.
func EqualTo(expected interface{}) *Matcher {
	expectedType := reflect.Typeof(expected)
	match := func (actual interface{}) *Result {
		actualType := reflect.Typeof(expected)
		if actualType == expectedType {
			if actual == expected {
				return NewResult(true, 
					NewDescription("was equal to [%v]", expected))
			}
			return NewResult(false,
				NewDescription("[%v] was not equal to [%v]", actual, expected))
		}
		return NewResult(false,
			NewDescription("[%v] could not be compared to [%v]", actual, expected))
	}
	return NewMatcher(NewDescription("EqualTo[%v]", expected), match)
}

// Returns a Matcher that checks if the actual value is (deeply)
// equal to the given expected value, using reflect.DeepEqual.
func DeeplyEqualTo(expected interface{}) *Matcher {
	match := func (actual interface{}) *Result {
		if reflect.DeepEqual(expected, actual) {
			return NewResult(true,
				NewDescription("was deeply equal to [%v]", expected))
		}
		return NewResult(false,
			NewDescription("[%v] was not deeply equal to [%v]", expected))
	}
	return NewMatcher(NewDescription("DeeplyEqualTo[%v]", expected), match)
}


// First part of a builder for a short-circuiting both/and matcher:
//     matcher := Both(Matcher1).And(Matcher2)
func Both(matcher *Matcher) *BothClause {
	return &BothClause{matcher:matcher}
}

// Intermediate state in the construction of a Both/And clause.
type BothClause struct {
	matcher *Matcher
}

// Second part of a builder for a short-circuiting both/and matcher.
func (self *BothClause) And(matcher2 *Matcher) *Matcher {
	matcher1 := self.matcher
	description := NewDescription("both [%v] and [%v]", matcher1, matcher2)
	match := func(v interface{}) *Result {
		result1 := matcher1.Match(v)
		if !result1.Matched() {
			because := NewDescription("first part of 'Both/And' did not match")
			return NewResult(false, because).WithCauses(result1)
		}
		result2 := matcher2.Match(v)
		if !result2.Matched() {
			because := NewDescription("second part of 'Both/And' did not match")
			return NewResult(false, because).WithCauses(result2)
		}
		because := NewDescription("both parts of 'Both/And' matched")
		return NewResult(true, because).WithCauses(result1, result2)
	}
	return NewMatcher(description, match)
}


// First part of a builder for a short-circuiting either/or matcher or
// a either/xor matcher, such as:
//     matcher := Either(matcher1).Or(matcher2)
// or:
//     matcher := Either(matcher1).Xor(matcher2)
func Either(matcher *Matcher) *EitherClause {
	return &EitherClause{matcher:matcher}
}

// Intermediate state in the construction of an Either/Or 
// or Either/Xor clause.
type EitherClause struct {
	matcher *Matcher
}

// Second part of a builder for a short-circuiting either/or matcher:
//     matcher := Either(matcher1).Or(matcher2)
// This matcher short-circuits without invoking the second matcher if
// the first matcher successfully matches, and matches whenever either
// of the two component matches successfully matches.
func (self *EitherClause) Or(matcher2 *Matcher) *Matcher {
	matcher1 := self.matcher
	description := NewDescription("either [%v] or [%v]", matcher1, matcher2)
	match := func(v interface{}) *Result {
		result1 := matcher1.Match(v)
		if result1.Matched() {
			because := NewDescription("first part of 'Either/Or' matched")
			return NewResult(true, because).WithCauses(result1)
		}
		result2 := matcher2.Match(v)
		if result2.Matched() {
			because := NewDescription("second part of 'Either/Or' matched")
			return NewResult(true, because).WithCauses(result2)
		}
		because := NewDescription("neither part of 'Either/Or' matched")
		return NewResult(false, because).WithCauses(result1, result2)
	}
	return NewMatcher(description, match)
}

// Second part of a builder for an either/xor matcher:
//     matcher := Either(matcher1).Xor(matcher2)
// This matcher matches when exactly one of the two matchers matches
// a given value;  if both or neither of the matchers is successful,
// xor fails to match.  Note that this is *never* a short-circuiting
// operation.
func (self *EitherClause) Xor(matcher2 *Matcher) *Matcher {
	matcher1 := self.matcher
	description := NewDescription("either [%v] xor [%v]", matcher1, matcher2)
	match := func(v interface{}) *Result {
		result1 := matcher1.Match(v)
		result2 := matcher2.Match(v)
		if result1.Matched() {
			if result2.Matched() {
				because := NewDescription("both parts of 'Either/Xor' matched")
				return NewResult(false, because).WithCauses(result1, result2)
			}
			because := NewDescription("only the first part of 'Either/Xor' matched")
			return NewResult(true, because).WithCauses(result1, result2)
		}
		if result2.Matched() {
			because := NewDescription("only the second part of 'Either/Xor' matched")
			return NewResult(true, because).WithCauses(result1, result2)
		}
		because := NewDescription("neither part of 'Either/Xor' matched")
		return NewResult(false, because).WithCauses(result1, result2)
	}
	return NewMatcher(description, match)
}

// First part of a builder for a short-circuiting neither/nor matcher:
//     matcher := Neither(matcher1).Or(matcher2)
// such that the second matcher is only tested if the first matcher
// fails to match, and the resulting matcher matches if either matches.
func Neither(matcher *Matcher) *NeitherClause {
	return &NeitherClause{matcher:matcher}
}

// Intermediate state in the construction of a Neither/Nor clause. 
type NeitherClause struct {
	matcher *Matcher
}

// Creates a matcher that passes when neither this matcher nor the
// other matcher pass.  This operation is short-circuiting, so that
// if the first matcher matches, the second is not attempted.
func (self *NeitherClause) Nor(matcher2 *Matcher) *Matcher {
	matcher1 := self.matcher
	description := NewDescription("neither [%v] nor [%v]", matcher1, matcher2)
	match := func(v interface{}) *Result {
		result1 := matcher1.Match(v)
		if result1.Matched() {
			because := NewDescription("first part of 'Nor' matched")
			return NewResult(false, because).WithCauses(result1)
		}
		result2 := matcher2.Match(v)
		if result2.Matched() {
			because := NewDescription("second part of 'Nor' matched")
			return NewResult(false, because).WithCauses(result2)
		}
		because := NewDescription("neither part of 'Nor' matched")
		return NewResult(true, because).WithCauses(result1, result2)
	}
	return NewMatcher(description, match)
}



// First part of a builder for a short-circuiting if/then matcher:
//     matcher := If(AntecedentMatcher).Then(ConsequentMatcher)
// such that the consequent is only tested when the antecedent
// matches, and the resulting matcher only fails to match when the
// consequent fails to match. Note that this is logically
// equivalent to:
//     Not(AntecedentMatcher).Or(ConsequentMatcher)
// But may be more readable in practice.
func If(antecedent *Matcher) *IfClause {
	return &IfClause{antecedent:antecedent}
}

// Temporary builder state in the middle of constructing
// an If/Then clause.
type IfClause struct {
	antecedent *Matcher
}

// Constructs a short-circuiting if/then matcher:
//     matcher := If(AntecedentMatcher).Then(ConsequentMatcher)
// such that the consequent is only tested when the antecedent
// matches, and the resulting matcher only fails to match when the
// consequent fails to match. Note that this is logically
// equivalent to:
//     Either(Not(AntecedentMatcher)).Or(ConsequentMatcher)
// But may be more readable in practice.
func (self *IfClause) Then(consequent *Matcher) *Matcher {
	antecedent := self.antecedent
	description := NewDescription("if [%v] then [%v]", antecedent, consequent)
	match := func(v interface{}) *Result {
		result1 := antecedent.Match(v)
		if !result1.Matched() {
			because := NewDescription("'If/Then' matched because antecedent failed")
			return NewResult(true, because).WithCauses(result1)
		}
		result2 := consequent.Match(v)
		if result2.Matched() {
			because := NewDescription("'If/Then' matched because consequent matched")
			return NewResult(true, because).WithCauses(result2)
		}
		because := NewDescription("'If/Then' failed")
		return NewResult(false, because).WithCauses(result1, result2)
	}
	return NewMatcher(description, match)
}


// First part of a builder for an if-and-only-if expression:
//     matcher := IfAndOnlyIf(AntecedentMatcher).Then(ConsequentMatcher)
// This is logically equivalent to:
//     Either(Not(AntecedentMatcher)).Xor(ConsequentMatcher)
// But may be more readable in practice.
func IfAndOnlyIf(antecedent *Matcher) *IfAndOnlyIfClause {
	return &IfAndOnlyIfClause{antecedent:antecedent}
}

// Temporary builder state in the middle of constructing
// an IfAndOnlyIf/Then clause.
type IfAndOnlyIfClause struct {
	antecedent *Matcher
}

// Constructs an if-and-only-if/then matcher:
//     matcher := IfAndOnlyIf(AntecedentMatcher).Then(ConsequentMatcher)
// that matches when both or neither of the Antecedent and the
// Consequent match.  Note that this is logically equivalent to:
//     Either(Not(AntecedentMatcher)).Xor(ConsequentMatcher)
// But may be more readable in practice.
func (self *IfAndOnlyIfClause) Then(consequent *Matcher) *Matcher {
	antecedent := self.antecedent
	description := NewDescription("if and only if [%v] then [%v]", antecedent, consequent)
	match := func(v interface{}) *Result {
		result1 := antecedent.Match(v)
		result2 := consequent.Match(v)
		if result1.Matched() {
			if result2.Matched() {
				because := NewDescription("Matched because both parts of 'IfAndOnlyIf/Then' matched")
				return NewResult(true, because).WithCauses(result1, result2)
			}
			because := NewDescription("Failed because only the first part of 'IfAndOnlyIf/Then' matched")
			return NewResult(false, because).WithCauses(result1, result2)
		}
		if result2.Matched() {
			because := NewDescription("Failed because only the second part of 'IfAndOnlyIf/Then' matched")
			return NewResult(false, because).WithCauses(result1, result2)
		}
		because := NewDescription("Matched because neither part of 'IfAndOnlyIf/Then' matched")
		return NewResult(true, because).WithCauses(result1, result2)
	}
	return NewMatcher(description, match)
}


