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
	
	if yesAndYes := Both(yes).And(yes).Match(0); !yesAndYes.Matched() {
		t.Errorf("yes and yes should match, was [%v]", yesAndYes)
	}
	if yesAndNo := Both(yes).And(no).Match(0); yesAndNo.Matched() {
		t.Errorf("yes and no should not match, was [%v]", yesAndNo)
	}
	noAndSnoop := Both(no).And(snoop).Match(0)
	if calledSnoop {
		t.Errorf("no and snoop should short-circuit before calling snoop")
	}
	if noAndSnoop.Matched() {
		t.Errorf("no and snoop should not match, was [%v]", noAndSnoop)
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
	
	if noOrNo := Either(no).Or(no).Match(0); noOrNo.Matched() {
		t.Errorf("no or no should not match, was [%v]", noOrNo)
	}
	if noOrYes := Either(no).Or(yes).Match(0); !noOrYes.Matched() {
		t.Errorf("no or yes should match, was [%v]", noOrYes)
	}
	yesOrSnoop := Either(yes).Or(snoop).Match(0)
	if calledSnoop {
		t.Errorf("yes or snoop should short-circuit before calling snoop")
	}
	if !yesOrSnoop.Matched() {
		t.Errorf("yes or snoop should match, was [%v]", yesOrSnoop)
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
	
	if noNorNo := Neither(no).Nor(no).Match(0); !noNorNo.Matched() {
		t.Errorf("no nor no should match, was [%v]", noNorNo)
	}
	if noNorYes := Neither(no).Nor(yes).Match(0); noNorYes.Matched() {
		t.Errorf("no nor yes should not match, was [%v]", noNorYes)
	}
	yesNorSnoop := Neither(yes).Nor(snoop).Match(0)
	if calledSnoop {
		t.Errorf("yes nor snoop should short-circuit before calling snoop")
	}
	if yesNorSnoop.Matched() {
		t.Errorf("yes nor snoop should not match, was [%v]", yesNorSnoop)
	}
}

func Test_EitherXor(t *testing.T) {
	yes, no := Anything(), Not(Anything())
	if result := Either(no).Xor(no).Match(0); result.Matched() {
		t.Errorf("no xor no should not match, was [%v]", result)
	}
	if result := Either(no).Xor(yes).Match(0); !result.Matched() {
		t.Errorf("no xor yes should match, was [%v]", result)
	}
	if result := Either(yes).Xor(no).Match(0); !result.Matched() {
		t.Errorf("yes xor no should match, was [%v]", result)
	}
	if result := Either(yes).Xor(yes).Match(0); result.Matched() {
		t.Errorf("yes xor yes should not match, was [%v]", result)
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
	ifNoThenSnoop := If(no).Then(snoop).Match(0)
	if calledSnoop {
		t.Errorf("If-no-then-snoop should short-circuit before calling snoop")
	}
	if !ifNoThenSnoop.Matched() {
		t.Errorf("if-no-then-snoop should match on failing antecedent, was [%v]",
			ifNoThenSnoop)
	}
}



