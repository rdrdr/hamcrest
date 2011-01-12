package main

import (
	"hamcrest/asserter"
	"hamcrest/strings"
	"hamcrest/logic"
	"hamcrest/collections"
	"hamcrest/comparison"
)






var sayings = []string {
	"Don't count your foobars before they hatch",
	"My dog has fleas",
	"My foobar has fleas",
	"My bar has foos",
	"My food has fleas",
}

var we = asserter.UsingStderr()

func init() {
	AtLeast := comparison.GreaterThanOrEqualTo
	we.AssertThat(len(sayings), AtLeast(3).Comment("Should have at least 3 testcases"))
}

func main() {

	println("This is a demo of hamcrest asserters.")
	println()
	println("Here's an example of a (non-fatal) check failing:")
	we.CheckThat("team", strings.Contains("I").Comment("Coach told me this"))
	println()
	println("Here's an example of a fatal check failing:")
	println()
	
	EveryElement := collections.EveryElement
	If := logic.If
	
	
	we.AssertThat(sayings, EveryElement(
		If(strings.HasPrefix("foo")).Then(strings.Contains("bar"))))
	
}