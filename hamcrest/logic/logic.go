// Copyright 2011 Mick Killianey.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package logic

import (
	"hamcrest"
)

var NewDescription = hamcrest.NewDescription
var NewMatcher = hamcrest.NewMatcher
var NewResult = hamcrest.NewResult

// First part of a builder for a short-circuiting both/and matcher:
//     matcher := Both(Matcher1).And(Matcher2)
func Both(matcher *hamcrest.Matcher) *BothClause {
	return &BothClause{matcher:matcher}
}

// Intermediate state in the construction of a Both/And clause.
type BothClause struct {
	matcher *hamcrest.Matcher
}

// Second part of a builder for a short-circuiting both/and matcher.
func (self *BothClause) And(matcher2 *hamcrest.Matcher) *hamcrest.Matcher {
	matcher1 := self.matcher
	description := hamcrest.NewDescription("both [%v] and [%v]", matcher1, matcher2)
	match := func(actual interface{}) *hamcrest.Result {
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
func Either(matcher *hamcrest.Matcher) *EitherClause {
	return &EitherClause{matcher:matcher}
}

// Intermediate state in the construction of an Either/Or 
// or Either/Xor clause.
type EitherClause struct {
	matcher *hamcrest.Matcher
}

// Second part of a builder for a short-circuiting either/or matcher:
//     matcher := Either(matcher1).Or(matcher2)
// This matcher short-circuits without invoking the second matcher if
// the first matcher successfully matches, and matches whenever either
// of the two component matches successfully matches.
func (self *EitherClause) Or(matcher2 *hamcrest.Matcher) *hamcrest.Matcher {
	matcher1 := self.matcher
	description := NewDescription("either [%v] or [%v]", matcher1, matcher2)
	match := func(actual interface{}) *hamcrest.Result {
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
func (self *EitherClause) Xor(matcher2 *hamcrest.Matcher) *hamcrest.Matcher {
	matcher1 := self.matcher
	description := NewDescription("either [%v] xor [%v]", matcher1, matcher2)
	match := func(actual interface{}) *hamcrest.Result {
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
func Neither(matcher *hamcrest.Matcher) *NeitherClause {
	return &NeitherClause{matcher:matcher}
}

// Intermediate state in the construction of a Neither/Nor clause. 
type NeitherClause struct {
	matcher *hamcrest.Matcher
}

// Creates a matcher that passes when neither this matcher nor the
// other matcher pass.  This operation is short-circuiting, so that
// if the first matcher matches, the second is not attempted.
//  Note that this is logically equivalent to:
//     Both(Not(matcher1)).And(Not(matcher2))
// But may be more readable in practice.
func (self *NeitherClause) Nor(matcher2 *hamcrest.Matcher) *hamcrest.Matcher {
	matcher1 := self.matcher
	description := NewDescription("neither [%v] nor [%v]", matcher1, matcher2)
	match := func(actual interface{}) *hamcrest.Result {
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
func If(antecedent *hamcrest.Matcher) *IfClause {
	return &IfClause{antecedent:antecedent}
}

// Temporary builder state in the middle of constructing
// an If/Then clause.
type IfClause struct {
	antecedent *hamcrest.Matcher
}

// Constructs a short-circuiting if/then matcher:
//     matcher := If(AntecedentMatcher).Then(ConsequentMatcher)
// such that the consequent is only tested when the antecedent
// matches, and the resulting matcher only fails to match when the
// consequent fails to match. Note that this is logically
// equivalent to:
//     Either(Not(AntecedentMatcher)).Or(ConsequentMatcher)
// But may be more readable in practice.
func (self *IfClause) Then(consequent *hamcrest.Matcher) *hamcrest.Matcher {
	antecedent := self.antecedent
	description := NewDescription("if [%v] then [%v]", antecedent, consequent)
	match := func(actual interface{}) *hamcrest.Result {
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
func IfAndOnlyIf(antecedent *hamcrest.Matcher) *IfAndOnlyIfClause {
	return &IfAndOnlyIfClause{antecedent:antecedent}
}

// Temporary builder state in the middle of constructing
// an IfAndOnlyIf/Then clause.
type IfAndOnlyIfClause struct {
	antecedent *hamcrest.Matcher
}

// Constructs an if-and-only-if/then matcher:
//     matcher := IfAndOnlyIf(AntecedentMatcher).Then(ConsequentMatcher)
// that matches when both or neither of the Antecedent and the
// Consequent match.  Note that this is logically equivalent to:
//     Either(Not(AntecedentMatcher)).Xor(ConsequentMatcher)
// But may be more readable in practice.
func (self *IfAndOnlyIfClause) Then(consequent *hamcrest.Matcher) *hamcrest.Matcher {
	antecedent := self.antecedent
	description := NewDescription("if and only if [%v] then [%v]", antecedent, consequent)
	match := func(actual interface{}) *hamcrest.Result {
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

