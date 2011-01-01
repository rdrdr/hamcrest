package strings

import (
	"fmt"
	"hamcrest"
	"hamcrest/asserter"
	"testing"
)


var Not = hamcrest.Not
var Both = hamcrest.Both

func Test_ToString_builtinTypes(t *testing.T) {
	we := asserter.Using(t)
	we.CheckThat(123, Not(hamcrest.DeeplyEqualTo("123")))
	we.CheckThat(123, ToString(hamcrest.DeeplyEqualTo("123")))
}

type _Stringer struct { s string }
func (self *_Stringer) String() string {
	return string(self.s)
}

func Test_ToString_onTypesThatImplementStringer(t *testing.T) {
	stringer := &_Stringer{s:"123"}
	we := asserter.Using(t)
	we.CheckThat(stringer, Not(hamcrest.DeeplyEqualTo("123")))
	we.CheckThat(stringer, ToString(hamcrest.DeeplyEqualTo("123")))
}

type _Formatter struct { s string }
func (self *_Formatter) Format(s fmt.State, ch int) {
	fmt.Fprint(s, string(self.s))
}

func Test_ToString_onTypesThatImplementFormatter(t *testing.T) {
	formatter := &_Formatter{s:"123"}
	we := asserter.Using(t)
	we.CheckThat(formatter, Not(hamcrest.DeeplyEqualTo("123")))
	we.CheckThat(formatter, ToString(hamcrest.DeeplyEqualTo("123")))
}

func Test_Contains(t *testing.T) {
	we := asserter.Using(t)
	we.CheckThat("abcde", Contains("abc").AddComments("prefix"))
	we.CheckThat("abcde", Contains("bcd").AddComments("middle"))
	we.CheckThat("abcde", Contains("cde").AddComments("suffix"))
	we.CheckThat("abcde", Not(Contains("ace")).AddComments("not a match"))
	we.CheckThat(123, Not(Contains("123")).AddComments("not a string"))
}

func Test_HasPrefix(t *testing.T) {
	we := asserter.Using(t)
	we.CheckThat("abcde", HasPrefix("abc").AddComments("actual prefix"))
	we.CheckThat("abcde", Not(HasPrefix("cde").AddComments("not a prefix")))
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

func Test_HasPattern(t *testing.T) {
	we := asserter.Using(t)
	we.CheckThat("abcd", HasPattern("a+b+").AddComments("at beginning of string"))
	we.CheckThat("abcd", HasPattern("b+c+").AddComments("in middle of string"))
	we.CheckThat("abcd", HasPattern("c+d+").AddComments("at end of string"))
	we.CheckThat("abcd", Not(HasPattern("[xy]+")))
}
