// Copyright 2011 Mick Killianey.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

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
	match := func (actual interface{}) *Result {
		if b, ok := actual.(bool); ok {
			if b {
				return NewResult(true, trueDescription)
			}
			return NewResult(false, falseDescription)
		}
		return NewResult(false, NewDescription("[%v] was not bool", actual))
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
	match := func (actual interface{}) *Result {
		if b, ok := actual.(bool); ok {
			if b {
				return NewResult(false, trueDescription)
			}
			return NewResult(true, falseDescription)
		}
		return NewResult(false, NewDescription("[%v] was not bool", actual))
	}
	_False = NewMatcher(NewDescription("matches false"), match)
}


// Returns a Matcher that matches the boolean value true.
func Matched() *Matcher {
	return _Matched
}
var _Matched *Matcher // singleton
func init() {
	description := NewDescription("was a result")
	match := func (actual interface{}) *Result {
		if result, ok := actual.(*Result); ok {
			return NewResult(result.Matched(), description).WithCauses(result)
		}
		return NewResult(false, NewDescription("[%v] was not a result", actual))
	}
	_Matched = NewMatcher(NewDescription("Matched"), match)
}

// Returns a Matcher that matches the boolean value false.
func DidNotMatch() *Matcher {
	return _DidNotMatch
}
var _DidNotMatch *Matcher // singleton
func init() {
	description := NewDescription("was a result")
	match := func (actual interface{}) *Result {
		if result, ok := actual.(*Result); ok {
			return NewResult(!result.Matched(), description).WithCauses(result)
		}
		return NewResult(false, NewDescription("[%v] was not a result", actual))
	}
	_DidNotMatch = NewMatcher(NewDescription("DidNotMatch"), match)
}


// Returns a Matcher that decorates another matcher and only matches
// when the underlying matcher does not match (and vice versa).
func Not(matcher *Matcher) *Matcher {
	match := func (actual interface{}) *Result {
		result := matcher.Match(actual)
		return NewResult(!result.matched, result.description).WithCauses(result)
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
	if value, ok := reflect.NewValue(actual).(_CanAskIsNil); ok {
		return value.IsNil()
	}
	return false
}
type _CanAskIsNil interface { IsNil() bool }


// Returns a Matcher that matches if the actual value is nil
// or the nil value of its type.  (Note that this is *not*
// equivalent to DeeplyEqualTo(nil).)
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
		return NewResult(false, NewDescription("[%v] was not nil", actual))
	}
	_Nil = NewMatcher(NewDescription("matches nil"), match)
}

// Returns a Matcher that matches if the actual value is 
// neither nil nor the nil value of its type.  (Note that
// this is *not* equivalent to Not(DeeplyEqualTo(nil)).)
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
		return NewResult(true, NewDescription("[%v] was not nil", actual))
	}
	_NonNil = NewMatcher(NewDescription("matches non-nil"), match)
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
			NewDescription("[%v] was not deeply equal to [%v]", actual, expected))
	}
	return NewMatcher(NewDescription("DeeplyEqualTo[%v]", expected), match)
}


// Returns a short-circuiting Matcher that matches whenever all of
// the given matchers match a given input value.  If any component
// matcher fails to match an input value, later matchers are not
// attempted.
func AllOf(matchers...*Matcher) *Matcher {
	descriptions := make([]*Description, len(matchers), len(matchers))
	for index, matcher := range matchers {
		descriptions[index] = NewDescription("[%v: %v]", index+1, matcher)
	}
	description := NewDescription("AllOf%v", descriptions)
	match := func (actual interface{}) *Result {
		var results []*Result
		for index, matcher := range matchers {
			result := matcher.Match(actual)
			results := append(results, result)
			if !result.Matched() {
				because := NewDescription(
					"Failed matcher %v of %v: [%v]",
						index+1, len(matchers), matcher)
				return NewResult(false, because).WithCauses(results...)
			}
		}
		because := NewDescription("Matched all %v matchers", len(matchers))
		return NewResult(true, because).WithCauses(results...)
	}
	return NewMatcher(description, match)
}

// Returns a short-circuiting Matcher that matches whenever all of
// the given matchers match a given input value.  If any component
// matcher fails to match an input value, later matchers are not
// attempted.
func AnyOf(matchers...*Matcher) *Matcher {
	descriptions := make([]*Description, len(matchers), len(matchers))
	for index, matcher := range matchers {
		descriptions[index] = NewDescription("[%v: %v]", index+1, matcher)
	}
	description := NewDescription("AnyOf%v", descriptions)
	match := func (actual interface{}) *Result {
		var results []*Result
		for index, matcher := range matchers {
			result := matcher.Match(actual)
			results := append(results, result)
			if result.Matched() {
				because := NewDescription(
					"Matched on matcher %v of %v: [%v]", index+1, len(matchers), matcher)
				return NewResult(true, because).WithCauses(results...)
			}
		}
		because := NewDescription("Matched none of the %v matchers", len(matchers))
		return NewResult(false, because).WithCauses(results...)
	}
	return NewMatcher(description, match)
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
	match := func(actual interface{}) *Result {
		result1 := matcher1.Match(actual)
		if !result1.Matched() {
			because := NewDescription("first part of 'Both/And' did not match [%v]", actual)
			return NewResult(false, because).WithCauses(result1)
		}
		result2 := matcher2.Match(actual)
		if !result2.Matched() {
			because := NewDescription("second part of 'Both/And' did not match [%v]", actual)
			return NewResult(false, because).WithCauses(result2)
		}
		because := NewDescription("both parts of 'Both/And' matched [%v]", actual)
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
	match := func(actual interface{}) *Result {
		result1 := matcher1.Match(actual)
		if result1.Matched() {
			because := NewDescription("first part of 'Either/Or' matched [%v]", actual)
			return NewResult(true, because).WithCauses(result1)
		}
		result2 := matcher2.Match(actual)
		if result2.Matched() {
			because := NewDescription("second part of 'Either/Or' matched [%v]", actual)
			return NewResult(true, because).WithCauses(result2)
		}
		because := NewDescription("neither part of 'Either/Or' matched [%v]", actual)
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
	match := func(actual interface{}) *Result {
		result1 := matcher1.Match(actual)
		result2 := matcher2.Match(actual)
		if result1.Matched() {
			if result2.Matched() {
				because := NewDescription("both parts of 'Either/Xor' matched [%v]", actual)
				return NewResult(false, because).WithCauses(result1, result2)
			}
			because := NewDescription("only the first part of 'Either/Xor' matched [%v]", actual)
			return NewResult(true, because).WithCauses(result1, result2)
		}
		if result2.Matched() {
			because := NewDescription("only the second part of 'Either/Xor' matched [%v]", actual)
			return NewResult(true, because).WithCauses(result1, result2)
		}
		because := NewDescription("neither part of 'Either/Xor' matched [%v]", actual)
		return NewResult(false, because).WithCauses(result1, result2)
	}
	return NewMatcher(description, match)
}

// First part of a builder for a short-circuiting neither/nor matcher:
//     matcher := Neither(matcher1).Nor(matcher2)
// such that the second matcher is only tested if the first matcher
// fails to match, and the resulting matcher matches if either matches.
// Note that the expression is logically equivalent to:
//     Both(Not(matcher1)).And(Not(matcher2))
// But may be more readable in practice.
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
//  Note that this is logically equivalent to:
//     Both(Not(matcher1)).And(Not(matcher2))
// But may be more readable in practice.
func (self *NeitherClause) Nor(matcher2 *Matcher) *Matcher {
	matcher1 := self.matcher
	description := NewDescription("neither [%v] nor [%v]", matcher1, matcher2)
	match := func(actual interface{}) *Result {
		result1 := matcher1.Match(actual)
		if result1.Matched() {
			because := NewDescription("first part of 'Nor' matched [%v]", actual)
			return NewResult(false, because).WithCauses(result1)
		}
		result2 := matcher2.Match(actual)
		if result2.Matched() {
			because := NewDescription("second part of 'Nor' matched [%v]", actual)
			return NewResult(false, because).WithCauses(result2)
		}
		because := NewDescription("neither part of 'Nor' matched [%v]", actual)
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
//     Either(Not(AntecedentMatcher)).Or(ConsequentMatcher)
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
	match := func(actual interface{}) *Result {
		result1 := antecedent.Match(actual)
		if !result1.Matched() {
			because := NewDescription("'If/Then' matched because antecedent failed on [%v]", actual)
			return NewResult(true, because).WithCauses(result1)
		}
		result2 := consequent.Match(actual)
		if result2.Matched() {
			because := NewDescription("'If/Then' matched because consequent matched on [%v]", actual)
			return NewResult(true, because).WithCauses(result2)
		}
		because := NewDescription("'If/Then' failed on [%v]", actual)
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
	match := func(actual interface{}) *Result {
		result1 := antecedent.Match(actual)
		result2 := consequent.Match(actual)
		if result1.Matched() {
			if result2.Matched() {
				because := NewDescription("Matched because both parts of 'Iff/Then' matched on [%v]", actual)
				return NewResult(true, because).WithCauses(result1, result2)
			}
			because := NewDescription("Failed because only the first part of 'Iff/Then' matched on [%v]", actual)
			return NewResult(false, because).WithCauses(result1, result2)
		}
		if result2.Matched() {
			because := NewDescription("Failed because only the second part of 'IFf/Then' matched on [%v]", actual)
			return NewResult(false, because).WithCauses(result1, result2)
		}
		because := NewDescription("Matched because neither part of 'Iff/Then' matched on [%v]", actual)
		return NewResult(true, because).WithCauses(result1, result2)
	}
	return NewMatcher(description, match)
}

