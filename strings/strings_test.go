// Copyright 2011 Mick Killianey.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package strings

import (
	"fmt"
	"github.com/rdrdr/hamcrest/asserter"
	"github.com/rdrdr/hamcrest/base"
	. "github.com/rdrdr/hamcrest/core"
	"testing"
)

var Matched = base.Matched()
var DidNotMatch = base.DidNotMatch()

type _Stringer struct { s string }
func (self *_Stringer) String() string {
	return string(self.s)
}

type _Formatter struct { s string }
func (self *_Formatter) Format(s fmt.State, ch int) {
	fmt.Fprint(s, string(self.s))
}

type _GoStringer struct { s string }
func (self *_GoStringer) GoString() string {
	return string(self.s)
}

func Test_ToString_builtinTypes(t *testing.T) {
	we := asserter.Using(t)
	we.CheckThat(EqualTo("123").Match(123), DidNotMatch)
	we.CheckThat(ToString(EqualTo("123")).Match(123), Matched)
}

func Test_ToString_onTypesThatImplementStringer(t *testing.T) {
	we := asserter.Using(t)
	stringer := &_Stringer{s:"123"}
	we.CheckThat(EqualTo("123").Match(stringer), DidNotMatch)
	we.CheckThat(ToString(EqualTo("123")).Match(stringer), Matched)
}

func Test_ToString_onTypesThatImplementFormatter(t *testing.T) {
	we := asserter.Using(t)
	formatter := &_Formatter{s:"123"}
	we.CheckThat(EqualTo("123").Match(formatter), DidNotMatch)
	we.CheckThat(ToString(EqualTo("123")).Match(formatter), Matched)
}

func Test_ToGoString_onTypesThatImplementGoStringer(t *testing.T) {
	we := asserter.Using(t)
	gostringer := &_GoStringer{s:"123"}
	we.CheckThat(EqualTo("123").Match(gostringer), DidNotMatch)
	we.CheckThat(ToGoString(EqualTo("123")).Match(gostringer), Matched)
}

func Test_ToLower(t *testing.T) {
	we := asserter.Using(t)
	we.CheckThat(ToLower(EqualTo("shout")).Match("shout"), Matched)
	we.CheckThat(ToLower(EqualTo("shout")).Match("SHOUT"), Matched)
	we.CheckThat(ToLower(EqualTo("SHOUT")).Match("shout"), DidNotMatch)
	we.CheckThat(ToLower(EqualTo("SHOUT")).Match("SHOUT"), DidNotMatch)
	we.CheckThat(ToLower(EqualTo("123")).Match("123"), Matched)
	we.CheckThat(ToLower(EqualTo("123")).Match(123), DidNotMatch)
}

func Test_ToUpper(t *testing.T) {
	we := asserter.Using(t)
	we.CheckThat(ToUpper(EqualTo("shout")).Match("shout"), DidNotMatch)
	we.CheckThat(ToUpper(EqualTo("shout")).Match("SHOUT"), DidNotMatch)
	we.CheckThat(ToUpper(EqualTo("SHOUT")).Match("shout"), Matched)
	we.CheckThat(ToUpper(EqualTo("SHOUT")).Match("SHOUT"), Matched)
	we.CheckThat(ToUpper(EqualTo("123")).Match("123"), Matched)
	we.CheckThat(ToUpper(EqualTo("123")).Match(123), DidNotMatch)
}

func Test_ToLen(t *testing.T) {
	we := asserter.Using(t)
	we.CheckThat(ToLen(EqualTo(5)).Match("SHOUT"), Matched)
	we.CheckThat(ToLen(EqualTo(5)).Match("OOPS"), DidNotMatch)
	we.CheckThat(ToLen(Anything()).Match(0), DidNotMatch.
		Comment("should fail when can't determine length"))
}

func Test_EqualToIgnoringCase(t *testing.T) {
	we := asserter.Using(t)
	we.CheckThat(EqualToIgnoringCase("one").Match("one"), Matched)
	we.CheckThat(EqualToIgnoringCase("one").Match("ONE"), Matched)
	we.CheckThat(EqualToIgnoringCase("ONE").Match("one"), Matched)
	we.CheckThat(EqualToIgnoringCase("ONE").Match("ONE"), Matched)
	we.CheckThat(EqualToIgnoringCase("oNe").Match("OnE"), Matched)
	we.CheckThat(EqualToIgnoringCase("one").Match("two"), DidNotMatch)
	we.CheckThat(EqualToIgnoringCase("one").Match(1), DidNotMatch)
	we.CheckThat(EqualToIgnoringCase("one").Match(nil), DidNotMatch)
}


func Test_Contains(t *testing.T) {
	we := asserter.Using(t)
	we.CheckThat(Contains("abc").Match("abcde"), Matched.Comment("prefix"))
	we.CheckThat(Contains("bcd").Match("abcde"), Matched.Comment("middle"))
	we.CheckThat(Contains("cde").Match("abcde"), Matched.Comment("suffix"))
	we.CheckThat(Contains("ace").Match("abcde"), DidNotMatch)
	we.CheckThat(Contains("123").Match(123), DidNotMatch)
	we.CheckThat(Contains("123").Match(nil), DidNotMatch)
}

var alphabet = "abcdefghijklmnopqrstuvwxyz"

func Test_HasPrefix(t *testing.T) {
	we := asserter.Using(t)
	we.CheckThat(HasPrefix("abc").Match(alphabet), Matched)
	we.CheckThat(HasPrefix("cde").Match(alphabet), DidNotMatch)
	we.CheckThat(HasPrefix("123").Match(123), DidNotMatch)
	we.CheckThat(HasPrefix("123").Match(nil), DidNotMatch)

	failResultString := HasPrefix("123").Match(alphabet).String()
	we.CheckThat(failResultString, Contains("123"))
	we.CheckThat(failResultString, Contains("abcdef"))
	we.CheckThat(failResultString, Not(Contains("xyz")).
		Comment("Should truncate input on failResult"))
	
	passResultString := HasPrefix("abc").Match(alphabet).String()
	we.CheckThat(passResultString, Contains("abcdef"))
	we.CheckThat(passResultString, Not(Contains("xyz")).
		Comment("Should truncate input on passResult"))
}

func Test_HasSuffix(t *testing.T) {
	we := asserter.Using(t)
	we.CheckThat(HasSuffix("xyz").Match(alphabet), Matched)
	we.CheckThat(HasSuffix("wxy").Match(alphabet), DidNotMatch)
	we.CheckThat(HasSuffix("123").Match(123), DidNotMatch)
	we.CheckThat(HasSuffix("123").Match(nil), DidNotMatch)
	
	failResultString := HasSuffix("123").Match(alphabet).String()
	we.CheckThat(failResultString, Contains("123"))
	we.CheckThat(failResultString, Contains("uvwxyz"))
	we.CheckThat(failResultString, Not(Contains("abc")).
		Comment("Should truncate input on failResult"))
	
	passResultString := HasSuffix("xyz").Match(alphabet).String()
	we.CheckThat(passResultString, Contains("uvwxyz"))
	we.CheckThat(passResultString, Not(Contains("abc")).
		Comment("Should truncate input on passResult"))
}

func Test_HasPattern(t *testing.T) {
	we := asserter.Using(t)
	we.CheckThat(HasPattern("a+b+").Match("abcd"), Matched.Comment("beginning"))
	we.CheckThat(HasPattern("b+c+").Match("abcd"), Matched.Comment("middle"))
	we.CheckThat(HasPattern("c+d+").Match("abcd"), Matched.Comment("end"))
	we.CheckThat(HasPattern("[xy]").Match("abcd"), DidNotMatch)
}

func Test_EachPattern(t *testing.T) {
	we := asserter.Using(t)
	eachGoPlusIsGoo := EachPattern("go+")(Is(EqualTo("goo")))
	we.CheckThat(eachGoPlusIsGoo.Match("stop stop stop"), Matched)
	we.CheckThat(eachGoPlusIsGoo.Match("goo goo goo"), Matched)
	we.CheckThat(eachGoPlusIsGoo.Match("goo go goo"), DidNotMatch)
	we.CheckThat(eachGoPlusIsGoo.Match("go go go"), DidNotMatch)
	we.CheckThat(eachGoPlusIsGoo.Match(123), DidNotMatch)
	we.CheckThat(eachGoPlusIsGoo.Match(nil), DidNotMatch)
	
	i_before_e := EachPattern("[^aeiou]ei")(HasPrefix("c"))
	we.CheckThat("ceiling receipt", i_before_e)
	we.CheckThat("deceiver seizure", Not(i_before_e))
}

func Test_EachPatternGroup(t *testing.T) {
	we := asserter.Using(t)
	eachQHasU := EachPatternGroup("([qQ])(.)", 2)(ToLower(EqualTo("u")))
	we.CheckThat(eachQHasU.Match("Quick quack MQ"), Matched)
	we.CheckThat(eachQHasU.Match("Quick FAQ for MQ"), DidNotMatch)
	we.CheckThat(eachQHasU.Match(123), DidNotMatch)
	we.CheckThat(eachQHasU.Match(nil), DidNotMatch)
	
	ei_after_c := EachPatternGroup("[cC](ei|ie)", 1)(EqualTo("ei"))
	we.CheckThat(ei_after_c.Match("Ceiling receipt"), Matched)
	we.CheckThat(ei_after_c.Match("receive Cienfuegos"), DidNotMatch)
}

func Test_AnyPattern(t *testing.T) {
	we := asserter.Using(t)
	
	anyGoPlusIsGoo := AnyPattern("go+")(Is(EqualTo("goo")))
	we.CheckThat(anyGoPlusIsGoo.Match("stop stop stop"), DidNotMatch)
	we.CheckThat(anyGoPlusIsGoo.Match("goo goo goo"), Matched)
	we.CheckThat(anyGoPlusIsGoo.Match("goo go goo"), Matched)
	we.CheckThat(anyGoPlusIsGoo.Match("go go go"), DidNotMatch)
	we.CheckThat(anyGoPlusIsGoo.Match(123), DidNotMatch)
	we.CheckThat(anyGoPlusIsGoo.Match(nil), DidNotMatch)
	
	has_cat_word := AnyPattern("[a-z]+at")(HasPrefix("c"))
	we.CheckThat("that cravat is phat", has_cat_word)
	we.CheckThat("Matt spat at a rat", Not(has_cat_word))
}

func Test_AnyPatternGroup(t *testing.T) {
	we := asserter.Using(t)
	
	isHappy := AnyPatternGroup("[:;B]-?(.)", 1)(Is(EqualTo(")")))
	we.CheckThat(isHappy.Match("x :- y+z :-)"), Matched)
	we.CheckThat(isHappy.Match(":-P :-/ ;-}"), DidNotMatch)
	we.CheckThat(isHappy.Match("<>v<> o rly?"), DidNotMatch)
	we.CheckThat(isHappy.Match(123), DidNotMatch)
	we.CheckThat(isHappy.Match(nil), DidNotMatch)
}

func Test_OnPattern(t *testing.T) {
	we := asserter.Using(t)
	
	firstGoPlusIsGoo := OnPattern("Go+")(Is(EqualTo("Goo")))
	we.CheckThat(firstGoPlusIsGoo.Match("take Google to Go"), Matched)
	we.CheckThat(firstGoPlusIsGoo.Match("take Go to Google"), DidNotMatch)
	we.CheckThat(firstGoPlusIsGoo.Match("Goooooooal"), DidNotMatch)
	we.CheckThat(firstGoPlusIsGoo.Match("no big-G goo"), DidNotMatch)
	we.CheckThat(firstGoPlusIsGoo.Match(123), DidNotMatch)
	we.CheckThat(firstGoPlusIsGoo.Match(nil), DidNotMatch)
	
	his := OnPattern("h.s")(EqualTo("his"))
	we.CheckThat(his.Match("hers and his"), Matched)
	we.CheckThat(his.Match("just hers"), DidNotMatch.
		Comment("no instances"))
	we.CheckThat(his.Match("has chisel"), DidNotMatch.
		Comment("second instance matches, but not first"))
}

func Test_OnPatternGroup(t *testing.T) {
	we := asserter.Using(t)
	
	americanComedy := OnPatternGroup("hum(ou?)r", 1)(Is(EqualTo("o")))
	we.CheckThat(americanComedy.Match("scatological humor"), Matched)
	we.CheckThat(americanComedy.Match("spamalot humour"), DidNotMatch)
	we.CheckThat(americanComedy.Match("humor I mean humour"), Matched)
	we.CheckThat(americanComedy.Match("humour I mean humor"), DidNotMatch)
	we.CheckThat(americanComedy.Match("not funny"), DidNotMatch)
	we.CheckThat(americanComedy.Match(123), DidNotMatch)
	we.CheckThat(americanComedy.Match(nil), DidNotMatch)
}

