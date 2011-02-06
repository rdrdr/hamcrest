// Copyright 2011 Mick Killianey.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import (
	"github.com/rdrdr/hamcrest/base"
	"github.com/rdrdr/hamcrest/asserter"
	"fmt"
	"reflect"
	"testing"
)

var Matched = base.Matched()
var DidNotMatch = base.DidNotMatch()

type Stringer interface { String() string }

var uninitialized struct {
	_pointer *bool
	_func func()
	_slice []int
	_chan chan int
	_map map[int]bool
	_interface interface{}
}

func checkMatcherIsMatchingOnNils(t *testing.T, matcher *base.Matcher) {
	we := asserter.Using(t)
	we.CheckThat(matcher.Match(nil), Matched.Comment("nil"))
	we.CheckThat(matcher.Match(uninitialized._pointer), Matched.Comment("uninitialized pointer"))
	we.CheckThat(matcher.Match(uninitialized._func), Matched.Comment("uninitialized func"))
	we.CheckThat(matcher.Match(uninitialized._slice), Matched.Comment("uninitialized slice"))
	we.CheckThat(matcher.Match(uninitialized._chan), Matched.Comment("uninitialized chan"))
	we.CheckThat(matcher.Match(uninitialized._map), Matched.Comment("uninitialized map"))
	we.CheckThat(matcher.Match(uninitialized._interface), Matched.Comment("uninitialized interface"))
}

func checkMatcherIsNonMatchingOnNils(t *testing.T, matcher *base.Matcher) {
	we := asserter.Using(t)
	we.CheckThat(matcher.Match(nil), DidNotMatch.Comment("nil"))
	we.CheckThat(matcher.Match(uninitialized._pointer), DidNotMatch.Comment("uninitialized pointer"))
	we.CheckThat(matcher.Match(uninitialized._func), DidNotMatch.Comment("uninitialized func"))
	we.CheckThat(matcher.Match(uninitialized._slice), DidNotMatch.Comment("uninitialized slice"))
	we.CheckThat(matcher.Match(uninitialized._chan), DidNotMatch.Comment("uninitialized chan"))
	we.CheckThat(matcher.Match(uninitialized._map), DidNotMatch.Comment("uninitialized map"))
	we.CheckThat(matcher.Match(uninitialized._interface), DidNotMatch.Comment("uninitialized interface"))
}

var sampleValues []interface{}
func init() {
	sampleValues = []interface{}{
		true,
		false,
		int(42), int8(42), int16(42), int32(42), int64(42),
		uint(42), uint8(42), uint16(42), uint32(42), uint64(42), 
		float32(42), float64(42),
		complex(42,42), complex64(42), complex128(42),
		"42",
		struct {Field int} {Field:42},
		&struct {Field int} {Field:42},
		reflect.Typeof(struct {Field int} {Field:42}),
		make(chan int, 42),
		func() int { return 42 },
		map[string]int{ "forty":40, "two":2, "forty-two":42 },
		[]int{40, 41, 42},
		[...]int{40, 41, 42},
		[42]int{2:40, 40:2},
		nil,
		uninitialized,
		uninitialized._chan,
		uninitialized._func,
		uninitialized._interface,
		uninitialized._map,
		uninitialized._pointer,
		uninitialized._slice,
	}
}

func logSamples(t *testing.T, matcher *base.Matcher) {
	t.Logf("Sample results for: %v\n", matcher)
	we := asserter.Using(t)
	for index, value := range sampleValues {
		t.Logf("Sample #%v: %T[value: %v]\n", index+1, value, value)
		we.LogResult(matcher.Match(value))
	}
}

// Check Matchers
func TestAnything(t *testing.T) {
	we := asserter.Using(t)
	matcher := Anything()
	we.CheckThat(matcher.Match(true), Matched)
	we.CheckThat(matcher.Match(false), Matched)
	we.CheckThat(matcher.Match("foo"), Matched)
	checkMatcherIsMatchingOnNils(t, matcher)
	logSamples(t, matcher)
}

func Test_True(t *testing.T) {
	we := asserter.Using(t)
	matcher := True()
	we.CheckThat(matcher.Match(true), Matched)
	we.CheckThat(matcher.Match(false), DidNotMatch)
	we.CheckThat(matcher.Match("true"), DidNotMatch)
	we.CheckThat(matcher.Match(1), DidNotMatch)
	checkMatcherIsNonMatchingOnNils(t, matcher)
	logSamples(t, matcher)
}

func Test_False(t *testing.T) {
	we := asserter.Using(t)
	matcher := False()
	we.CheckThat(matcher.Match(true), DidNotMatch)
	we.CheckThat(matcher.Match(false), Matched)
	we.CheckThat(matcher.Match("false"), DidNotMatch)
	we.CheckThat(matcher.Match(0), DidNotMatch)
	checkMatcherIsNonMatchingOnNils(t, matcher)
	logSamples(t, matcher)
}

func Test_Not(t *testing.T) {
	we := asserter.Using(t)
	matcher := Not(True())
	we.CheckThat(matcher.Match(true), DidNotMatch)
	we.CheckThat(matcher.Match(false), Matched)
	logSamples(t, matcher)
}

func Test_Is(t *testing.T) {
	we := asserter.Using(t)
	matcher := Is(True())
	we.CheckThat(matcher.Match(true), Matched)
	we.CheckThat(matcher.Match(false), DidNotMatch)
	logSamples(t, matcher)
}

func Test_Nil(t *testing.T) {
	we := asserter.Using(t)
	matcher := Nil()
	we.CheckThat(matcher.Match(false), DidNotMatch)
	we.CheckThat(matcher.Match(0), DidNotMatch)
	we.CheckThat(matcher.Match("nil"), DidNotMatch)
	checkMatcherIsMatchingOnNils(t, matcher)
	logSamples(t, matcher)
}

func Test_NonNil(t *testing.T) {
	we := asserter.Using(t)
	matcher := NonNil()
	we.CheckThat(matcher.Match(false), Matched)
	we.CheckThat(matcher.Match(0), Matched)
	we.CheckThat(matcher.Match("nil"), Matched)
	checkMatcherIsNonMatchingOnNils(t, matcher)
	logSamples(t, matcher)
}

type _DeepEqualType struct { x int }
func Test_DeepEqualTo(t *testing.T) {
	we := asserter.Using(t)
	data := []interface{} {
		nil, true, false,
		int(42), uint(42), float64(42), complex128(42),
		struct { x int } { x: 42 },
		struct { x int } { x: 42 },
		&struct { x int } { x: 42 },
		struct { y int } { y: 42 },
		_DeepEqualType { x: 42 },
		&_DeepEqualType { x: 42 },
		[]int { 42 },
		[]int { 42 },
		map[int]int{ 42: 42 },
		map[int]int{ 42: 42 },
		make(chan int, 42),
		make(chan int, 42),
	}
	for _, x := range data {
		matcher := DeepEqualTo(x)
		for _, y := range data {
			message := fmt.Sprintf("%T[%v] and %T[%v]", x, x, y, y)
			if reflect.DeepEqual(x, y) {
				we.CheckThat(matcher.Match(y), Matched.Comment(message))
			} else {
				we.CheckThat(matcher.Match(y), DidNotMatch.Comment(message))
			}
		}
	}
	logSamples(t, DeepEqualTo(42))
}

func Test_AllOf(t *testing.T) {
	we := asserter.Using(t)
	yes, no := Anything(), Not(Anything())
	calledSnoop := false
	snoop := base.NewMatcherf(func(v interface{}) *base.Result {
			calledSnoop = true
			return base.NewResultf(false, "snooped!")
		}, "Snoop")
	
	we.CheckThat(AllOf(yes, yes, yes).Match(0), Matched.Comment("all matched"))
	we.CheckThat(AllOf(yes, yes, no).Match(0), DidNotMatch.Comment("not all matched"))
	we.CheckThat(AllOf(yes).Match(0), Matched.Comment("can pass one matcher"))
	we.CheckThat(AllOf(no).Match(0), DidNotMatch.Comment("can fail one matcher"))
	we.CheckThat(AllOf(yes, no, snoop).Match(0), DidNotMatch.Comment("can short-circuit"))
	we.CheckFalse(calledSnoop, "AllOf should short-circuit on first non-match")
	logSamples(t, AllOf(Not(True()), NonNil(), EqualTo(42)))
}

func Test_AnyOf(t *testing.T) {
	we := asserter.Using(t)
	yes, no := Anything(), Not(Anything())
	calledSnoop := false
	snoop := base.NewMatcherf(func(v interface{}) *base.Result {
			calledSnoop = true
			return base.NewResultf(false, "snooped!")
		}, "Snoop")
	we.CheckThat(AnyOf(no, no, no).Match(0), DidNotMatch.Comment("none matched"))
	we.CheckThat(AnyOf(no, no, yes).Match(0), Matched.Comment("one matched"))
	we.CheckThat(AnyOf(yes).Match(0), Matched.Comment("can pass one matcher"))
	we.CheckThat(AnyOf(no).Match(0), DidNotMatch.Comment("can fail one matcher"))
	we.CheckThat(AnyOf(no, yes, snoop).Match(0), Matched.Comment("can short-circuit"))
	we.CheckFalse(calledSnoop, "AnyOf should short-circuit on first match")
	logSamples(t, AnyOf(True(), Nil(), EqualTo(42)))
}

func Test_Applying_onFunction_FromType_ToType(t *testing.T) {
	we := asserter.Using(t)
	IsEven := Applying(func(n int) bool { return n&1 == 0 }, "IsEven")
	ToLength := Applying(func(s string) int { return len(s) }, "ToLength")
	ToString := Applying(func(i int) string { return fmt.Sprint(i) }, "ToString")
	
	ValueIsEven := IsEven(Is(True()))
	LengthIsEven := ToString(ToLength(ValueIsEven))
	
	we.CheckThat(ValueIsEven.Match(123), DidNotMatch)
	we.CheckThat(LengthIsEven.Match(123), DidNotMatch)
	
	we.CheckThat(ValueIsEven.Match(1234), Matched)
	we.CheckThat(LengthIsEven.Match(1234), Matched)
	
	we.CheckThat(ValueIsEven.Match(124), Matched)
	we.CheckThat(LengthIsEven.Match(124), DidNotMatch)
	
	we.CheckThat(ValueIsEven.Match(1233), DidNotMatch)
	we.CheckThat(LengthIsEven.Match(1233), Matched)
	
	logSamples(t, ValueIsEven)
	logSamples(t, LengthIsEven)
}

func Test_Applying_onFunction_FromTypeDotDotDot_ToType(t *testing.T) {
	we := asserter.Using(t)
	IsEven := Applying(func(n...int) bool { return n[0]&1 == 0 }, "IsEven")
	ToLength := Applying(func(s... string) int { return len(s[0]) }, "ToLength")
	ToString := Applying(func(i... int) string { return fmt.Sprint(i[0]) }, "ToString")
	
	ValueIsEven := IsEven(Is(True()))
	LengthIsEven := ToString(ToLength(ValueIsEven))
	
	we.CheckThat(ValueIsEven.Match(123), DidNotMatch)
	we.CheckThat(LengthIsEven.Match(123), DidNotMatch)
	
	we.CheckThat(ValueIsEven.Match(1234), Matched)
	we.CheckThat(LengthIsEven.Match(1234), Matched)
	
	we.CheckThat(ValueIsEven.Match(124), Matched)
	we.CheckThat(LengthIsEven.Match(124), DidNotMatch)
	
	we.CheckThat(ValueIsEven.Match(1233), DidNotMatch)
	we.CheckThat(LengthIsEven.Match(1233), Matched)
	
	logSamples(t, ValueIsEven)
	logSamples(t, LengthIsEven)
}

func Test_Applying_onFunction_FromTypeTypeDotDotDot_ToType(t *testing.T) {
	we := asserter.Using(t)
	IsEven := Applying(func(n int, other...string) bool { return n&1 == 0 }, "IsEven")
	ToLength := Applying(func(s string, other...int) int { return len(s) }, "ToLength")
	ToString := Applying(func(i int, other... string) string { return fmt.Sprint(i) }, "ToString")
	
	ValueIsEven := IsEven(Is(True()))
	LengthIsEven := ToString(ToLength(ValueIsEven))
	
	we.CheckThat(ValueIsEven.Match(123), DidNotMatch)
	we.CheckThat(LengthIsEven.Match(123), DidNotMatch)
	
	we.CheckThat(ValueIsEven.Match(1234), Matched)
	we.CheckThat(LengthIsEven.Match(1234), Matched)
	
	we.CheckThat(ValueIsEven.Match(124), Matched)
	we.CheckThat(LengthIsEven.Match(124), DidNotMatch)
	
	we.CheckThat(ValueIsEven.Match(1233), DidNotMatch)
	we.CheckThat(LengthIsEven.Match(1233), Matched)
	
	logSamples(t, ValueIsEven)
	logSamples(t, LengthIsEven)
}

func Test_PanicWhen_onFunctionAcceptingInterface(t *testing.T) {
	we := asserter.Using(t)
	panicOnBools := PanicWhenApplying(func (v interface{}) {
		if _, ok := v.(bool); ok {
			panic("no bools!")
		}
	}, "DisallowBools")

	we.CheckThat(panicOnBools.Match("true"), DidNotMatch)
	we.CheckThat(panicOnBools.Match(true), Matched)
	we.CheckThat(panicOnBools.Match(nil), DidNotMatch)
	
	logSamples(t, panicOnBools)
}

func Test_PanicWhen_onFunctionAcceptingBool(t *testing.T) {
	we := asserter.Using(t)
	
	var functionInvoked bool
	PanicOnFalse := PanicWhenApplying(func (b bool) {
		functionInvoked = true
		if !b { panic("Must be true") }
	}, "PanicOnFalse")

	functionInvoked = false
	we.CheckThat(PanicOnFalse.Match("true"), Matched.
		Comment("Should panic when can't invoke function"))
	we.CheckFalse(functionInvoked, "Shouldn't have invoked function")
	
	functionInvoked = false
	we.CheckThat(PanicOnFalse.Match(nil), Matched.
		Comment("Can't invoke function with string"))
	we.CheckFalse(functionInvoked, "Shouldn't have invoked function")
	
	functionInvoked = false
	we.CheckThat(PanicOnFalse.Match(true), DidNotMatch)
	we.CheckTrue(functionInvoked, "Should have invoked function")
	
	functionInvoked = false
	we.CheckThat(PanicOnFalse.Match(false), Matched)
	we.CheckTrue(functionInvoked, "Should have invoked function")

	logSamples(t, PanicOnFalse)
}

func Test_PanicWhen_onFunctionAcceptingOneArgDotDotDot(t *testing.T) {
	we := asserter.Using(t)

	var functionInvoked bool
	PanicOn13 := PanicWhenApplying(func (args...int) {
		functionInvoked = true
		if args[0] == 13 {
			panic("Superstition")
		}
	}, "Disallow13")

	functionInvoked = false
	we.CheckThat(PanicOn13.Match("thirteen"), Matched.
		Comment("Should panic when can't invoke function"))
	we.CheckFalse(functionInvoked, "Shouldn't have invoked function")
	
	functionInvoked = false
	we.CheckThat(PanicOn13.Match(12), DidNotMatch)
	we.CheckTrue(functionInvoked, "Should have invoked function")
	
	functionInvoked = false
	we.CheckThat(PanicOn13.Match(13), Matched)
	we.CheckTrue(functionInvoked, "Should have invoked function")

	logSamples(t, PanicOn13)
}

func Test_PanicWhen_onFunctionAcceptingTwoArgsDotDotDot(t *testing.T) {
	we := asserter.Using(t)

	var functionInvoked bool
	PanicOn13 := PanicWhenApplying(func (arg int, why...string) {
		functionInvoked = true
		if arg == 13 {
			panic("Superstition")
		}
	}, "Disallow13")

	functionInvoked = false
	we.CheckThat(PanicOn13.Match("thirteen"), Matched.
		Comment("Should panic when can't invoke function"))
	we.CheckFalse(functionInvoked, "Shouldn't have invoked function")
	
	functionInvoked = false
	we.CheckThat(PanicOn13.Match(12), DidNotMatch)
	we.CheckTrue(functionInvoked, "Should have invoked function")
	
	functionInvoked = false
	we.CheckThat(PanicOn13.Match(13), Matched)
	we.CheckTrue(functionInvoked, "Should have invoked function")

	logSamples(t, PanicOn13)
}


