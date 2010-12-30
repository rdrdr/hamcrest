package strings

import (
	"fmt"
	"hamcrest"
	"hamcrest/asserter"
	"testing"
)

var Not = hamcrest.Not
var Both = hamcrest.Both

func Test_Contains(t *testing.T) {
	we := asserter.Using(t)
	we.CheckThat("abcde", Contains("abc").AddComments("prefix"))
	we.CheckThat("abcde", Contains("bcd").AddComments("middle"))
	we.CheckThat("abcde", Contains("cde").AddComments("suffix"))
	we.CheckThat("abcde", Not(Contains("def")).AddComments("not a full match"))
	we.CheckThat(123, Not(Contains("123")).AddComments("not a string"))
}

func Test_HasPrefix(t *testing.T) {
	we := asserter.Using(t)
	we.CheckThat("abcdef", HasPrefix("abc").AddComments("actual prefix"))
	we.CheckThat("abcdef", Not(HasPrefix("def").AddComments("not a prefix")))
	we.CheckThat(123, Not(HasPrefix("123")).AddComments("not a string"))
	
	failResult := HasPrefix("123").Match("abcdef").String()
	we.CheckThat(failResult, Contains("123"))
	
	passResult := HasPrefix("abc").Match("abcdef").String()
	we.CheckThat(passResult, Contains("abc"))
}

func Test_HasSuffix(t *testing.T) {
	we := asserter.Using(t)
	we.CheckThat("abcdef", HasSuffix("def"))
	we.CheckThat("abcdef", Not(HasSuffix("abc")))
	we.CheckThat(123, Not(HasSuffix("123")))
	
	failResult := fmt.Sprintf("%v", HasSuffix("123").Match("abc"))
	we.CheckThat(failResult, Contains("123"))
	
	passResult := fmt.Sprintf("%v", HasSuffix("def").Match("abcdef"))
	we.CheckThat(passResult, Contains("def"))
}

//func Test_Message(t *testing.T) {
//	we := asserter.Using(t)
//	we.CheckThat(hamcrest.NewDescription("abcdef"),
//		hamcrest.Either(HasPrefix("foo")).Or(HasSuffix("bar")))
//}



