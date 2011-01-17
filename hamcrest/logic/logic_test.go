// Copyright 2011 Mick Killianey.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package logic

import (
	"hamcrest"
	"reflect"
	"testing"
)

var Anything = hamcrest.Anything
var True = hamcrest.True
var False = hamcrest.False
var Not = hamcrest.Not

func checkResultIsMatching(t *testing.T, result *hamcrest.Result, message string) {
	if !result.Matched() {
		t.Errorf("Expected matching result, was [%v] %v", result, message)
	}
}

func checkResultIsNonMatching(t *testing.T, result *hamcrest.Result, message string) {
	if result.Matched() {
		t.Errorf("Expected non-matching result, was [%v] %v", result, message)
	}
}

type Stringer interface { String() string }

var uninitialized struct {
	_pointer *bool
	_func func()
	_slice []int
	_chan chan int
	_map map[int]bool
	_interface interface{}
}

func checkMatcherIsMatchingOnNils(t *testing.T, matcher *hamcrest.Matcher) {
	checkResultIsMatching(t, matcher.Match(nil), "nil")
	checkResultIsMatching(t, matcher.Match(uninitialized._pointer), "uninitialized pointer")
	checkResultIsMatching(t, matcher.Match(uninitialized._func), "uninitialized func")
	checkResultIsMatching(t, matcher.Match(uninitialized._slice), "uninitialized slice")
	checkResultIsMatching(t, matcher.Match(uninitialized._chan), "uninitialized chan")
	checkResultIsMatching(t, matcher.Match(uninitialized._map), "uninitialized map")
	checkResultIsMatching(t, matcher.Match(uninitialized._interface), "uninitialized interface")
}

func checkMatcherIsNonMatchingOnNils(t *testing.T, matcher *hamcrest.Matcher) {
	checkResultIsNonMatching(t, matcher.Match(nil), "nil")
	checkResultIsNonMatching(t, matcher.Match(uninitialized._pointer), "uninitialized pointer")
	checkResultIsNonMatching(t, matcher.Match(uninitialized._func), "uninitialized func")
	checkResultIsNonMatching(t, matcher.Match(uninitialized._slice), "uninitialized slice")
	checkResultIsNonMatching(t, matcher.Match(uninitialized._chan), "uninitialized chan")
	checkResultIsNonMatching(t, matcher.Match(uninitialized._map), "uninitialized map")
	checkResultIsNonMatching(t, matcher.Match(uninitialized._interface), "uninitialized interface")
}

func logSamples(t *testing.T, matcher *hamcrest.Matcher) {
	t.Logf("Sample results for: %v\n", matcher)
	t.Logf("\ton true: %v\n", matcher.Match(true))
	t.Logf("\ton false: %v\n", matcher.Match(false))
	t.Logf("\ton int: %v\n", matcher.Match(42))
	t.Logf("\ton uint: %v\n", matcher.Match(uint(42)))
	t.Logf("\ton float: %v\n", matcher.Match(42.0))
	t.Logf("\ton string: %v\n", matcher.Match("foobar"))
	t.Logf("\ton struct: %v\n", matcher.Match(struct {Field int} {Field:42}))
	t.Logf("\ton type: %v\n", matcher.Match(reflect.Typeof(uninitialized)))
	
	t.Logf("\ton channel: %v\n", matcher.Match(make(chan int, 1)))
	t.Logf("\ton function: %v\n", matcher.Match(func() int { return 1 }))
	t.Logf("\ton function: %v\n", matcher.Match(interface{}(nil)))
	t.Logf("\ton map: %v\n", matcher.Match(map[int]string{1:"one", 2:"two"}))
	t.Logf("\ton pointer: %v\n", matcher.Match(&struct {Field int} {Field:42}))
	t.Logf("\ton slice: %v\n", matcher.Match([]int{1}))
	
	t.Logf("\ton nil: %v\n", matcher.Match(nil))
	t.Logf("\ton nil channel: %v\n", matcher.Match(uninitialized._chan))
	t.Logf("\ton nil function: %v\n", matcher.Match(uninitialized._func))
	t.Logf("\ton nil interface: %v\n", matcher.Match(uninitialized._interface))
	t.Logf("\ton nil map: %v\n", matcher.Match(uninitialized._map))
	t.Logf("\ton nil pointer: %v\n", matcher.Match(uninitialized._pointer))
	t.Logf("\ton nil slice: %v\n", matcher.Match(uninitialized._slice))
}

// Check Matchers
func Test_BothAnd(t *testing.T) {
	yes, no := Anything(), Not(Anything())
	calledSnoop := false
	snoop := hamcrest.NewMatcherf(func(v interface{}) *hamcrest.Result {
			calledSnoop = true
			return hamcrest.NewResultf(false, "snooped!")
		}, "Snoop")
	
	if result := Both(yes).And(yes).Match(0); !result.Matched() {
		t.Errorf("yes and yes should match, was [%v]", result)
	}
	if result := Both(yes).And(no).Match(0); result.Matched() {
		t.Errorf("yes and no should not match, was [%v]", result)
	}
	result := Both(no).And(snoop).Match(0)
	if calledSnoop {
		t.Errorf("no and snoop should short-circuit before calling snoop")
	}
	if result.Matched() {
		t.Errorf("no and snoop should not match, was [%v]", result)
	}
	logSamples(t, Both(yes).And(yes))
	logSamples(t, Both(yes).And(no))
	logSamples(t, Both(no).And(yes))
	logSamples(t, Both(no).And(no))
}

func Test_EitherOr(t *testing.T) {
	yes, no := Anything(), Not(Anything())
	calledSnoop := false
	snoop := hamcrest.NewMatcherf(func(v interface{}) *hamcrest.Result {
			calledSnoop = true
			return hamcrest.NewResultf(false, "snooped!")
		}, "Snoop")
	
	if result := Either(no).Or(no).Match(0); result.Matched() {
		t.Errorf("no or no should not match, was [%v]", result)
	}
	if result := Either(no).Or(yes).Match(0); !result.Matched() {
		t.Errorf("no or yes should match, was [%v]", result)
	}
	result := Either(yes).Or(snoop).Match(0)
	if calledSnoop {
		t.Errorf("yes or snoop should short-circuit before calling snoop")
	}
	if !result.Matched() {
		t.Errorf("yes or snoop should match, was [%v]", result)
	}
	logSamples(t, Either(yes).Or(yes))
	logSamples(t, Either(yes).Or(no))
	logSamples(t, Either(no).Or(yes))
	logSamples(t, Either(no).Or(no))
}

func Test_EitherXor(t *testing.T) {
	yes, no := Anything(), Not(Anything())
	if result := Either(no).Xor(no).Match(0); result.Matched() {
		t.Errorf("no xor no should not match, was [%v]", result)
	}
	if result := Either(yes).Xor(yes).Match(0); result.Matched() {
		t.Errorf("yes xor yes should not match, was [%v]", result)
	}
	if result := Either(no).Xor(yes).Match(0); !result.Matched() {
		t.Errorf("no xor yes should match, was [%v]", result)
	}
	if result := Either(yes).Xor(no).Match(0); !result.Matched() {
		t.Errorf("yes xor no should match, was [%v]", result)
	}
	logSamples(t, Either(yes).Xor(yes))
	logSamples(t, Either(yes).Xor(no))
	logSamples(t, Either(no).Xor(yes))
	logSamples(t, Either(no).Xor(no))
}

func Test_NeitherNor(t *testing.T) {
	yes, no := Anything(), Not(Anything())
	calledSnoop := false
	snoop := hamcrest.NewMatcherf(func(v interface{}) *hamcrest.Result {
			calledSnoop = true
			return hamcrest.NewResultf(false, "snooped!")
		}, "Snoop")
	
	if result := Neither(no).Nor(no).Match(0); !result.Matched() {
		t.Errorf("no nor no should match, was [%v]", result)
	}
	if result := Neither(no).Nor(yes).Match(0); result.Matched() {
		t.Errorf("no nor yes should not match, was [%v]", result)
	}
	result := Neither(yes).Nor(snoop).Match(0)
	if calledSnoop {
		t.Errorf("yes nor snoop should short-circuit before calling snoop")
	}
	if result.Matched() {
		t.Errorf("yes nor snoop should not match, was [%v]", result)
	}
	logSamples(t, Neither(yes).Nor(yes))
	logSamples(t, Neither(yes).Nor(no))
	logSamples(t, Neither(no).Nor(yes))
	logSamples(t, Neither(no).Nor(no))
}

func Test_IfThen(t *testing.T) {
	yes, no := Anything(), Not(Anything())
	calledSnoop := false
	snoop := hamcrest.NewMatcherf(func(v interface{}) *hamcrest.Result {
			calledSnoop = true
			return hamcrest.NewResultf(false, "snooped!")
		}, "Snoop")
	
	if result := If(yes).Then(yes).Match(0); !result.Matched() {
		t.Errorf("if yes then yes should match, was [%v]", result)
	}
	if result := If(yes).Then(no).Match(0); result.Matched() {
		t.Errorf("if yes then no should not match, was [%v]", result)
	}
	result := If(no).Then(snoop).Match(0)
	if calledSnoop {
		t.Errorf("If-no-then-snoop should short-circuit before calling snoop")
	}
	if !result.Matched() {
		t.Errorf("if-no-then-snoop should match on failing antecedent, was [%v]",
			result)
	}
	logSamples(t, If(yes).Then(yes))
	logSamples(t, If(yes).Then(no))
	logSamples(t, If(no).Then(yes))
	logSamples(t, If(no).Then(no))
}

func Test_IfAndOnlyIfThen(t *testing.T) {
	yes, no := Anything(), Not(Anything())
	Iff := IfAndOnlyIf
	if result := Iff(no).Then(no).Match(0); !result.Matched() {
		t.Errorf("iff no then no should match, was [%v]", result)
	}
	if result := Iff(yes).Then(yes).Match(0); !result.Matched() {
		t.Errorf("iff yes then yes should match, was [%v]", result)
	}
	if result := Iff(no).Then(yes).Match(0); result.Matched() {
		t.Errorf("iff no then yes should not match, was [%v]", result)
	}
	if result := Iff(yes).Then(no).Match(0); result.Matched() {
		t.Errorf("iff yes then no should match, was [%v]", result)
	}
	logSamples(t, Iff(yes).Then(yes))
	logSamples(t, Iff(yes).Then(no))
	logSamples(t, Iff(no).Then(yes))
	logSamples(t, Iff(no).Then(no))
}
