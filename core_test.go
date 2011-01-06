package hamcrest

import (
	"fmt"
	"reflect"
	"testing"
)

func checkResultIsMatching(t *testing.T, result *Result, message string) {
	if !result.Matched() {
		t.Errorf("Expected matching result, was [%v] %v", result, message)
	}
}

func checkResultIsNonMatching(t *testing.T, result *Result, message string) {
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

func checkMatcherIsMatchingOnNils(t *testing.T, matcher *Matcher) {
	checkResultIsMatching(t, matcher.match(nil), "nil")
	checkResultIsMatching(t, matcher.match(uninitialized._pointer), "uninitialized pointer")
	checkResultIsMatching(t, matcher.match(uninitialized._func), "uninitialized func")
	checkResultIsMatching(t, matcher.match(uninitialized._slice), "uninitialized slice")
	checkResultIsMatching(t, matcher.match(uninitialized._chan), "uninitialized chan")
	checkResultIsMatching(t, matcher.match(uninitialized._map), "uninitialized map")
	checkResultIsMatching(t, matcher.match(uninitialized._interface), "uninitialized interface")
}

func checkMatcherIsNonMatchingOnNils(t *testing.T, matcher *Matcher) {
	checkResultIsNonMatching(t, matcher.match(nil), "nil")
	checkResultIsNonMatching(t, matcher.match(uninitialized._pointer), "uninitialized pointer")
	checkResultIsNonMatching(t, matcher.match(uninitialized._func), "uninitialized func")
	checkResultIsNonMatching(t, matcher.match(uninitialized._slice), "uninitialized slice")
	checkResultIsNonMatching(t, matcher.match(uninitialized._chan), "uninitialized chan")
	checkResultIsNonMatching(t, matcher.match(uninitialized._map), "uninitialized map")
	checkResultIsNonMatching(t, matcher.match(uninitialized._interface), "uninitialized interface")
}

func logSamples(t *testing.T, matcher *Matcher) {
	t.Logf("Sample results for: %v\n", matcher)
	t.Logf("\ton true: %v\n", matcher.match(true))
	t.Logf("\ton false: %v\n", matcher.match(false))
	t.Logf("\ton int: %v\n", matcher.match(42))
	t.Logf("\ton uint: %v\n", matcher.match(uint(42)))
	t.Logf("\ton float: %v\n", matcher.match(42.0))
	t.Logf("\ton string: %v\n", matcher.match("foobar"))
	t.Logf("\ton struct: %v\n", matcher.match(struct {Field int} {Field:42}))
	t.Logf("\ton type: %v\n", matcher.match(reflect.Typeof(uninitialized)))
	
	t.Logf("\ton channel: %v\n", matcher.match(make(chan int, 1)))
	t.Logf("\ton function: %v\n", matcher.match(func() int { return 1 }))
	t.Logf("\ton function: %v\n", matcher.match(interface{}(nil)))
	t.Logf("\ton map: %v\n", matcher.match(map[int]string{1:"one", 2:"two"}))
	t.Logf("\ton pointer: %v\n", matcher.match(&struct {Field int} {Field:42}))
	t.Logf("\ton slice: %v\n", matcher.match([]int{1}))
	
	t.Logf("\ton nil: %v\n", matcher.match(nil))
	t.Logf("\ton nil channel: %v\n", matcher.match(uninitialized._chan))
	t.Logf("\ton nil function: %v\n", matcher.match(uninitialized._func))
	t.Logf("\ton nil interface: %v\n", matcher.match(uninitialized._interface))
	t.Logf("\ton nil map: %v\n", matcher.match(uninitialized._map))
	t.Logf("\ton nil pointer: %v\n", matcher.match(uninitialized._pointer))
	t.Logf("\ton nil slice: %v\n", matcher.match(uninitialized._slice))
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
	checkResultIsMatching(t, matcher.match(true), "bool true")
	checkResultIsNonMatching(t, matcher.match(false), "bool false")
	checkResultIsNonMatching(t, matcher.match("foo"), "string")
	checkMatcherIsNonMatchingOnNils(t, matcher)
	logSamples(t, matcher)
}

func Test_False(t *testing.T) {
	matcher := False()
	checkResultIsNonMatching(t, matcher.match(true), "bool true")
	checkResultIsMatching(t, matcher.match(false), "bool false")
	checkResultIsNonMatching(t, matcher.match("foo"), "string")
	checkMatcherIsNonMatchingOnNils(t, matcher)
	logSamples(t, matcher)
}

func Test_Matched(t *testing.T) {
	matcher := Matched()
	passResult := NewResult(true, NewDescription("pass"))
	failResult := NewResult(false, NewDescription("fail"))
	
	checkResultIsMatching(t, matcher.Match(passResult), "matching")
	checkResultIsNonMatching(t, matcher.Match(failResult), "non-matching")
	checkResultIsNonMatching(t, matcher.Match(nil), "nil")
	checkResultIsNonMatching(t, matcher.Match("foo"), "not a Result")
	checkMatcherIsNonMatchingOnNils(t, matcher)
	logSamples(t, matcher)
}

func Test_DidNotMatch(t *testing.T) {
	matcher := DidNotMatch()
	passResult := NewResult(true, NewDescription("pass"))
	failResult := NewResult(false, NewDescription("fail"))
	
	checkResultIsNonMatching(t, matcher.Match(passResult), "matching")
	checkResultIsMatching(t, matcher.Match(failResult), "non-matching")
	checkResultIsNonMatching(t, matcher.Match(nil), "nil")
	checkResultIsNonMatching(t, matcher.Match("foo"), "not a Result")
	checkMatcherIsNonMatchingOnNils(t, matcher)
	logSamples(t, matcher)
}

func Test_Not(t *testing.T) {
	matcher := Not(False())
	
	checkResultIsMatching(t, matcher.match(true), "bool true")
	checkResultIsNonMatching(t, matcher.match(false), "bool false")
	checkResultIsMatching(t, matcher.match("foo"), "string")
	checkMatcherIsMatchingOnNils(t, matcher)
	logSamples(t, matcher)
}

func Test_Is(t *testing.T) {
	matcher := Is(True())
	checkResultIsMatching(t, matcher.match(true), "bool true")
	checkResultIsNonMatching(t, matcher.match(false), "bool false")
	checkResultIsNonMatching(t, matcher.match("foo"), "string")
	checkMatcherIsNonMatchingOnNils(t, matcher)
	logSamples(t, matcher)
}

func Test_Nil(t *testing.T) {
	matcher := Nil()
	checkResultIsNonMatching(t, matcher.match(true), "bool")
	checkResultIsNonMatching(t, matcher.match(0), "number")
	checkResultIsNonMatching(t, matcher.match("foo"), "string")
	checkMatcherIsMatchingOnNils(t, matcher)
	logSamples(t, matcher)
}

func Test_NonNil(t *testing.T) {
	matcher := NonNil()
	checkResultIsMatching(t, matcher.match(true), "bool")
	checkResultIsMatching(t, matcher.match(0), "number")
	checkResultIsMatching(t, matcher.match("foo"), "string")
	checkMatcherIsNonMatchingOnNils(t, matcher)
	logSamples(t, matcher)
}

func Test_EqualTo(t *testing.T) {
	matcher := EqualTo(42)
	checkResultIsNonMatching(t, matcher.match(true), "bool")
	checkResultIsMatching(t, matcher.match(42), "int")
	checkResultIsNonMatching(t, matcher.match(uint(42)), "uint")
	checkResultIsNonMatching(t, matcher.match(42.0), "float")
	checkResultIsNonMatching(t, matcher.match(42+0i), "complex")
	checkResultIsNonMatching(t, matcher.match("42"), "string")
	checkResultIsNonMatching(t, matcher.match([]int{42}), "slice")
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
				checkResultIsMatching(t, matcher.match(y), message)
			} else {
				checkResultIsNonMatching(t, matcher.match(y), message)
			}
		}
	}
	logSamples(t, DeeplyEqualTo(42))
}

func Test_AllOf(t *testing.T) {
	yes, no := Anything(), Not(Anything())
	calledSnoop := false
	snoop := NewMatcher(NewDescription("snoop"), 
		func(v interface{}) *Result {
			calledSnoop = true
			return NewResult(false, NewDescription("snooped!"))
		})
	checkResultIsMatching(t, AllOf(yes, yes, yes).match(true), "all yes")
	checkResultIsNonMatching(t, AllOf(yes, yes, no).match(false), "not all yes")
	checkResultIsNonMatching(t, AllOf(yes, no, snoop).match(false),
		"should short-circuit on no without snooping")
	if calledSnoop {
		t.Errorf("AllOf should short-circuit before calling snoop")
	}
	logSamples(t, AllOf(Anything(), NonNil(), EqualTo(42)))
}

func Test_AnyOf(t *testing.T) {
	yes, no := Anything(), Not(Anything())
	calledSnoop := false
	snoop := NewMatcher(NewDescription("snoop"), 
		func(v interface{}) *Result {
			calledSnoop = true
			return NewResult(false, NewDescription("snooped!"))
		})
	checkResultIsNonMatching(t, AnyOf(no, no, no).match(true), "all no")
	checkResultIsMatching(t, AnyOf(no, no, yes).match(false), "one yes")
	checkResultIsMatching(t, AnyOf(no, yes, snoop).match(false),
		"should short-circuit on yes without snooping")
	if calledSnoop {
		t.Errorf("AnyOf should short-circuit before calling snoop")
	}
	logSamples(t, AnyOf(True(), Nil(), EqualTo(42)))
}

func Test_BothAnd(t *testing.T) {
	yes, no := Anything(), Not(Anything())
	calledSnoop := false
	snoop := NewMatcher(NewDescription("snoop"), 
		func(v interface{}) *Result {
			calledSnoop = true
			return NewResult(false, NewDescription("snooped!"))
		})
	
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
	snoop := NewMatcher(NewDescription("snoop"), 
		func(v interface{}) *Result {
			calledSnoop = true
			return NewResult(false, NewDescription("snooped!"))
		})
	
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
	snoop := NewMatcher(NewDescription("snoop"), 
		func(v interface{}) *Result {
			calledSnoop = true
			return NewResult(false, NewDescription("snooped!"))
		})
	
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
	snoop := NewMatcher(NewDescription("snoop"), 
		func(v interface{}) *Result {
			calledSnoop = true
			return NewResult(false, NewDescription("snooped!"))
		})
	
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
