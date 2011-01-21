package main

import (
	"hamcrest"
	"hamcrest/asserter"
	"hamcrest/strings"
	"hamcrest/logic"
	"hamcrest/collections"
	"hamcrest/comparison"
)


var we = asserter.UsingStderr()

func init() {
	AtLeast := comparison.GreaterThanOrEqualTo
	we.AssertThat(len(sayings), AtLeast(3).Comment("Should have at least 3 testcases"))
}

var sayings = map[string][]string {
	"Hamlet": {
		"To be or not, whatever",
		"Get thee to a foobarery",
	},
	"ConFOOcious": {
		"Don't count your foobars before they hatch",
		"My dog has fleas",
		"My foobar has fleas",
		"My bar has foos",
		"My food has fleas",
	},
	"Hall of Presidents": {
		"There is foo in my bar.",
		"Well, foo, there you bar again.",
		"No new foobars.",
		"I never fooed that bar.",
		"I'm the foobarerer.",
		"A vote for me is a vote for foo and bar.",
	},
}

func init() {
	ForEverySaying := func(m *hamcrest.Matcher) *hamcrest.Matcher {
		return collections.EveryMapElement(collections.EveryElement(m))
	}
	If := logic.If
	Contains := strings.Contains
	we.AssertThat(sayings, ForEverySaying(If(Contains("foo")).Then(Contains("bar"))))
}

func main() {}
