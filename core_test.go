package hamcrest

import (
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

// Check Matchers
func TestAnything(t *testing.T) {
	matcher := Anything()
	checkResultIsMatching(t, matcher.Match(nil), "Should match nil")
	checkResultIsMatching(t, matcher.Match("non-nil"), "Should match non-nil")
	checkResultIsMatching(t, matcher.Match(true), "Should match true")
	checkResultIsMatching(t, matcher.Match(false), "Should match false")
	checkResultIsMatching(t, matcher.Match(make([]int, 1)), "Should match slice")
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
}

func Test_Nil_onNil(t *testing.T) {
	if result := Nil().Match(nil); !result.Matched() {
		t.Errorf("Nil should match nil")
	}
}

func Test_Nil_onNilReference(t *testing.T) {
	var nilReference *int = nil
	if result := Nil().Match(nilReference); !result.Matched() {
		t.Errorf("Nil should match a nil reference: %v", nilReference)
	}
}

func Test_Nil_onInterface(t *testing.T) {
	var nilInterface interface{} = nil
	if result := Nil().Match(nilInterface); !result.Matched() {
		t.Errorf("Nil should match a nil interface: %v", result)
	}
}
type _SampleStruct struct { Field int }
func Test_Nil_onStructPtr(t *testing.T) {
	var nilStructPtr *_SampleStruct = nil
	if result := Nil().Match(nilStructPtr); !result.Matched() {
		t.Errorf("Nil should match a nil struct ptr: %v", result)
	}
}
