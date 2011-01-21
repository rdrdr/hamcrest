// Copyright 2011 Mick Killianey.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import (
	"hamcrest"
	"fmt"
	"reflect"
	"testing"
)

func _LogResult(t *testing.T, indent string, result *hamcrest.Result) {
	s := "Matched"
	if !result.Matched() {
		s = "Nonmatch"
	}
	t.Logf("%v%v: %v", indent, s, result)
	indent = "\t" + indent
	for _, cause := range result.Causes() {
		_LogResult(t, indent, cause)
	}
}

func checkResultIsMatching(t *testing.T, result *hamcrest.Result, message string) {
	if !result.Matched() {
		t.Errorf("Expected matching result from applying %v to %#v, was [%v] %v",
			result.Matcher(), result.Value(), result, message)
	}
	_LogResult(t, "", result)
}

func checkResultIsNonMatching(t *testing.T, result *hamcrest.Result, message string) {
	if result.Matched() {
		t.Errorf("Expected non-matching result from applying %v to %#v, was [%v].  Message: %v",
			result.Matcher(), result.Value(), result, message)
	}
	_LogResult(t, "", result)
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
func TestAnything(t *testing.T) {
	matcher := Anything()
	checkResultIsMatching(t, matcher.Match(nil), "Should match nil")
	checkResultIsMatching(t, matcher.Match("non-nil"), "Should match non-nil")
	checkResultIsMatching(t, matcher.Match(true), "Should match true")
	checkResultIsMatching(t, matcher.Match(false), "Should match false")
	checkResultIsMatching(t, matcher.Match(make([]int, 1)), "Should match slice")
	checkMatcherIsMatchingOnNils(t, matcher)
	logSamples(t, matcher)
}

func Test_True(t *testing.T) {
	matcher := True()
	checkResultIsMatching(t, matcher.Match(true), "bool true")
	checkResultIsNonMatching(t, matcher.Match(false), "bool false")
	checkResultIsNonMatching(t, matcher.Match("foo"), "string")
	checkMatcherIsNonMatchingOnNils(t, matcher)
	logSamples(t, matcher)
}

func Test_False(t *testing.T) {
	matcher := False()
	checkResultIsNonMatching(t, matcher.Match(true), "bool true")
	checkResultIsMatching(t, matcher.Match(false), "bool false")
	checkResultIsNonMatching(t, matcher.Match("foo"), "string")
	checkMatcherIsNonMatchingOnNils(t, matcher)
	logSamples(t, matcher)
}

func Test_Not(t *testing.T) {
	matcher := Not(False())
	
	checkResultIsMatching(t, matcher.Match(true), "bool true")
	checkResultIsNonMatching(t, matcher.Match(false), "bool false")
	checkResultIsMatching(t, matcher.Match("foo"), "string")
	checkMatcherIsMatchingOnNils(t, matcher)
	logSamples(t, matcher)
}

func Test_Is(t *testing.T) {
	matcher := Is(True())
	checkResultIsMatching(t, matcher.Match(true), "bool true")
	checkResultIsNonMatching(t, matcher.Match(false), "bool false")
	checkResultIsNonMatching(t, matcher.Match("foo"), "string")
	checkMatcherIsNonMatchingOnNils(t, matcher)
	logSamples(t, matcher)
}

func Test_Nil(t *testing.T) {
	matcher := Nil()
	checkResultIsNonMatching(t, matcher.Match(true), "bool")
	checkResultIsNonMatching(t, matcher.Match(0), "number")
	checkResultIsNonMatching(t, matcher.Match("foo"), "string")
	checkMatcherIsMatchingOnNils(t, matcher)
	logSamples(t, matcher)
}

func Test_NonNil(t *testing.T) {
	matcher := NonNil()
	checkResultIsMatching(t, matcher.Match(true), "bool")
	checkResultIsMatching(t, matcher.Match(0), "number")
	checkResultIsMatching(t, matcher.Match("foo"), "string")
	checkMatcherIsNonMatchingOnNils(t, matcher)
	logSamples(t, matcher)
}

type _DeeplyEqualsType struct { x int }
func Test_DeeplyEqualTo(t *testing.T) {
	data := []interface{} {
		nil, true, false,
		int(42), uint(42), float(42), complex(42),
		struct { x int } { x: 42 },
		struct { x int } { x: 42 },
		&struct { x int } { x: 42 },
		struct { y int } { y: 42 },
		_DeeplyEqualsType { x: 42 },
		&_DeeplyEqualsType { x: 42 },
		[]int { 42 },
		[]int { 42 },
		map[int]int{ 42: 42 },
		map[int]int{ 42: 42 },
		make(chan int, 42),
		make(chan int, 42),
	}
	for _, x := range data {
		matcher := DeeplyEqualTo(x)
		for _, y := range data {
			message := fmt.Sprintf("%T[%v] and %T[%v]", x, x, y, y)
			if reflect.DeepEqual(x, y) {
				checkResultIsMatching(t, matcher.Match(y), message)
			} else {
				checkResultIsNonMatching(t, matcher.Match(y), message)
			}
		}
	}
	logSamples(t, DeeplyEqualTo(42))
}

func Test_AllOf(t *testing.T) {
	yes, no := Anything(), Not(Anything())
	calledSnoop := false
	snoop := hamcrest.NewMatcherf(func(v interface{}) *hamcrest.Result {
			calledSnoop = true
			return hamcrest.NewResultf(false, "snooped!")
		}, "Snoop")
	checkResultIsMatching(t, AllOf(yes, yes, yes).Match(true), "all yes")
	checkResultIsNonMatching(t, AllOf(yes, yes, no).Match(false), "not all yes")
	checkResultIsNonMatching(t, AllOf(yes, no, snoop).Match(false),
		"should short-circuit on no without snooping")
	if calledSnoop {
		t.Errorf("AllOf should short-circuit before calling snoop")
	}
	logSamples(t, AllOf(Anything(), NonNil(), DeeplyEqualTo(42)))
}

func Test_AnyOf(t *testing.T) {
	yes, no := Anything(), Not(Anything())
	calledSnoop := false
	snoop := hamcrest.NewMatcherf(func(v interface{}) *hamcrest.Result {
			calledSnoop = true
			return hamcrest.NewResultf(false, "snooped!")
		}, "Snoop")
	checkResultIsNonMatching(t, AnyOf(no, no, no).Match(true), "all no")
	checkResultIsMatching(t, AnyOf(no, no, yes).Match(false), "one yes")
	checkResultIsMatching(t, AnyOf(no, yes, snoop).Match(false),
		"should short-circuit on yes without snooping")
	if calledSnoop {
		t.Errorf("AnyOf should short-circuit before calling snoop")
	}
	logSamples(t, AnyOf(True(), Nil(), DeeplyEqualTo(42)))
}

func Test_Applying_onFunction_FromType_ToType(t *testing.T) {
	CheckEven := Applying(func(n int) bool { return n&1 == 0 }, "CheckEven")
	checkResultIsMatching(t, CheckEven(Is(True())).Match(1234), "is even")
	checkResultIsMatching(t, CheckEven(Is(False())).Match(123), "is not even")
	checkResultIsMatching(t, Not(CheckEven(Is(False()))).Match(1234), "not is not even")
	checkResultIsMatching(t, Not(CheckEven(Is(True()))).Match(123), "not is even")
	Even := CheckEven(Is(True()))
	Length := Applying(func(s string) int { return len(s) }, "Length")
	checkResultIsMatching(t, Length(Is(Even)).Match("1234"), "length is even")
	checkResultIsMatching(t, Length(Not(Even)).Match("123"), "length is not even")
	checkResultIsMatching(t, Not(Length(Even)).Match("123"), "not length is even")
}

func Test_Applying_onFunction_FromTypeDotDotDot_ToType(t *testing.T) {
	CheckEven := Applying(func(n...int) bool { return n[0]&1 == 0 }, "CheckEven")
	checkResultIsMatching(t, CheckEven(Is(True())).Match(1234), "is even")
	checkResultIsMatching(t, CheckEven(Is(False())).Match(123), "is not even")
	checkResultIsMatching(t, Not(CheckEven(Is(False()))).Match(1234), "not is not even")
	checkResultIsMatching(t, Not(CheckEven(Is(True()))).Match(123), "not is even")
	Even := CheckEven(Is(True()))
	Length := Applying(func(s...string) int { return len(s[0]) }, "Length")
	checkResultIsMatching(t, Length(Is(Even)).Match("1234"), "length is even")
	checkResultIsMatching(t, Length(Not(Even)).Match("123"), "length is not even")
	checkResultIsMatching(t, Not(Length(Even)).Match("123"), "not length is even")
}

func Test_Applying_onFunction_FromTypeTypeDotDotDot_ToType(t *testing.T) {
	CheckEven := Applying(func(n int, other...string) bool { return n&1 == 0 }, "CheckEven")
	checkResultIsMatching(t, CheckEven(Is(True())).Match(1234), "is even")
	checkResultIsMatching(t, CheckEven(Is(False())).Match(123), "is not even")
	checkResultIsMatching(t, Not(CheckEven(Is(False()))).Match(1234), "not is not even")
	checkResultIsMatching(t, Not(CheckEven(Is(True()))).Match(123), "not is even")
	Even := CheckEven(Is(True()))
	Length := Applying(func(s string, other...int) int { return len(s) }, "Length")
	checkResultIsMatching(t, Length(Is(Even)).Match("1234"), "length is even")
	checkResultIsMatching(t, Length(Not(Even)).Match("123"), "length is not even")
	checkResultIsMatching(t, Not(Length(Even)).Match("123"), "not length is even")
}

func Test_PanicWhen_onFunctionAcceptingInterface(t *testing.T) {
	panicOnBools := PanicWhenApplying(func (v interface{}) {
		if _, ok := v.(bool); ok {
			panic("no bools!")
		}
	}, "DisallowBools")

	checkResultIsNonMatching(t, panicOnBools.Match("true"), "shouldn't panic on string")
	checkResultIsMatching(t, panicOnBools.Match(true), "should panic on bool")
}

func Test_PanicWhen_onFunctionAcceptingBool(t *testing.T) {
	panicOnNonTruthInvoked := false
	panicOnNonTruth := PanicWhenApplying(func (b bool) {
		panicOnNonTruthInvoked = true
		if !b { panic("Must be true") }
	}, "RequireTruth")

	checkResultIsMatching(t, panicOnNonTruth.Match("true"), "should panic when can't invoke")
	checkResultIsNonMatching(t, panicOnNonTruth.Match(true), "should invoke and not panic")
	if !panicOnNonTruthInvoked {
		t.Error("Didn't invoke panicOnNonTruth on true")
	}
	panicOnNonTruthInvoked = false
	checkResultIsMatching(t, panicOnNonTruth.Match(false), "should invoke and panic")
	if !panicOnNonTruthInvoked {
		t.Error("Didn't invoke panicOnNonTruth on false")
	}
}

func Test_PanicWhen_onFunctionAcceptingOneArgDotDotDot(t *testing.T) {
	panicOn13Invoked := false
	panicOn13 := PanicWhenApplying(func (args...int) {
		panicOn13Invoked = true
		if args[0] == 13 {
			panic("Superstition")
		}
	}, "Disallow13")

	checkResultIsMatching(t, panicOn13.Match("thirteen"), "should panic when can't invoke")
	checkResultIsNonMatching(t, panicOn13.Match(12), "should invoke and not panic")
	if !panicOn13Invoked {
		t.Error("Didn't invoke panicOn13 on 12")
	}
	panicOn13Invoked = false
	checkResultIsMatching(t, panicOn13.Match(13), "should invoke and panic")
	if !panicOn13Invoked {
		t.Error("Didn't invoke panicOn13 on 13")
	}
}

func Test_PanicWhen_onFunctionAcceptingTwoArgsDotDotDot(t *testing.T) {
	panicOn13Invoked := false
	panicOn13 := PanicWhenApplying(func (arg int, why...string) {
		panicOn13Invoked = true
		if arg == 13 {
			panic("Superstition")
		}
	}, "Disallow13")

	checkResultIsMatching(t, panicOn13.Match("thirteen"), "should panic when can't invoke")
	checkResultIsNonMatching(t, panicOn13.Match(12), "should invoke and not panic")
	if !panicOn13Invoked {
		t.Error("Didn't invoke panicOn13 on 12")
	}
	panicOn13Invoked = false
	checkResultIsMatching(t, panicOn13.Match(13), "should invoke and panic")
	if !panicOn13Invoked {
		t.Error("Didn't invoke panicOn13 on 13")
	}
}


