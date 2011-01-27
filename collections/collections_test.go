// Copyright 2011 Mick Killianey.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package collections

import (
	"hamcrest/asserter"
	"hamcrest/comparison"
	"hamcrest/core"
	//"reflect"
	//"runtime"
	"testing"
)

var Is = core.Is
var Not = core.Not
var Anything = core.Anything
var EqualTo = comparison.EqualTo
var LessThan = comparison.LessThan
var GreaterThan = comparison.GreaterThan

func Test_AnyElement_ofArray(t *testing.T) {
	we := asserter.Using(t)
	we.CheckThat([3]int{1, 2, 3}, AnyElement(EqualTo(1)).Comment("first element"))
	we.CheckThat([3]int{1, 2, 3}, AnyElement(EqualTo(2)).Comment("middle element"))
	we.CheckThat([3]int{1, 2, 3}, AnyElement(EqualTo(3)).Comment("last element"))
	we.CheckThat([3]int{1, 2, 3}, Not(AnyElement(EqualTo(4))).Comment("none matching"))
	we.CheckThat([0]int{}, Not(AnyElement(Anything())).Comment("no elements"))
}

func Test_EveryElement_ofArray(t *testing.T) {
	we := asserter.Using(t)
	we.CheckThat([3]int{1, 2, 3}, EveryElement(LessThan(4)).Comment("all elements"))
	we.CheckThat([3]int{2, 1, 1}, Not(EveryElement(LessThan(2))).Comment("all but first"))
	we.CheckThat([3]int{1, 2, 1}, Not(EveryElement(LessThan(2))).Comment("all but middle"))
	we.CheckThat([3]int{1, 1, 2}, Not(EveryElement(LessThan(2))).Comment("all but last"))
	we.CheckThat([0]int{}, EveryElement(Anything()).Comment("no elements"))
}

func Test_AnyElement_ofSlice(t *testing.T) {
	we := asserter.Using(t)
	we.CheckThat([]int{1, 2, 3}, AnyElement(EqualTo(1)).Comment("first element"))
	we.CheckThat([]int{1, 2, 3}, AnyElement(EqualTo(2)).Comment("middle element"))
	we.CheckThat([]int{1, 2, 3}, AnyElement(EqualTo(3)).Comment("last element"))
	we.CheckThat([]int{1, 2, 3}, Not(AnyElement(EqualTo(4))).Comment("none matching"))
	we.CheckThat([]int{}, Not(AnyElement(Anything())).Comment("no elements"))
}

func Test_EveryElement_ofSlice(t *testing.T) {
	we := asserter.Using(t)
	we.CheckThat([]int{1, 2, 3}, EveryElement(LessThan(4)).Comment("all elements"))
	we.CheckThat([]int{2, 1, 1}, Not(EveryElement(LessThan(2))).Comment("all but first"))
	we.CheckThat([]int{1, 2, 1}, Not(EveryElement(LessThan(2))).Comment("all but middle"))
	we.CheckThat([]int{1, 1, 2}, Not(EveryElement(LessThan(2))).Comment("all but last"))
	we.CheckThat([]int{}, EveryElement(Anything()).Comment("no elements"))
}

func Test_AnyMapElement(t *testing.T) {
	we := asserter.Using(t)
	twoMap := map[string]int{ "foo": 1, "bar": 2 }
	emptyMap := map[string]int{}
	we.CheckThat(twoMap, AnyMapElement(EqualTo(1)).
		Comment("entry [foo: 1]"))
	we.CheckThat(twoMap, AnyMapElement(EqualTo(2)).
		Comment("entry [bar: 2]"))
	we.CheckThat(twoMap, Not(AnyMapElement(EqualTo(3))).
		Comment("neither entry matches"))
	we.CheckThat(emptyMap, Not(AnyMapElement(Anything())).
		Comment("no entries"))
}

func Test_EveryMapElement(t *testing.T) {
	we := asserter.Using(t)
	twoMap := map[string]int{ "foo": 1, "bar": 2 }
	emptyMap := map[string]int{}
	we.CheckThat(twoMap, EveryMapElement(GreaterThan(0)).
		Comment("all entries"))
	we.CheckThat(twoMap, Not(EveryMapElement(GreaterThan(1))).
		Comment("not entry [foo: 1]"))
	we.CheckThat(twoMap, Not(EveryMapElement(LessThan(2))).
		Comment("not entry [bar: 2]"))
	we.CheckThat(emptyMap, EveryMapElement(Anything()).
		Comment("no entries"))
}

func Test_ToLen_onArrays(t *testing.T) {
	we := asserter.Using(t)
	empty := [...]string{}
	hasTwo := [...]string{"itsy", "bitsy"}
	we.CheckThat(empty, ToLen(EqualTo(0)))
	we.CheckThat(hasTwo, ToLen(EqualTo(2)))
}

func Test_ToLen_onSlices(t *testing.T) {
	we := asserter.Using(t)
	empty := []string{}
	hasTwo := []string{"itsy", "bitsy"}
	we.CheckThat(empty, ToLen(Is(EqualTo(0))))
	we.CheckThat(hasTwo, ToLen(Is(EqualTo(2))))
}

func Test_ToLen_onMaps(t *testing.T) {
	we := asserter.Using(t)
	empty := map[string]int{}
	hasTwo := map[string]int{ "foo": 1, "bar": 2 }
	we.CheckThat(empty, ToLen(Is(EqualTo(0))))
	we.CheckThat(hasTwo, ToLen(Is(EqualTo(2))))
}

func Test_Empty_onArrays(t *testing.T) {
	we := asserter.Using(t)
	empty := [...]string{}
	hasTwo := [...]string{"itsy", "bitsy"}
	we.CheckThat(empty, Is(Empty()))
	we.CheckThat(hasTwo, Is(Not(Empty())))
}

func Test_Empty_onSlices(t *testing.T) {
	we := asserter.Using(t)
	empty := []string{}
	hasTwo := []string{"itsy", "bitsy"}
	we.CheckThat(empty, Is(Empty()))
	we.CheckThat(hasTwo, Is(Not(Empty())))
}

func Test_Empty_onMaps(t *testing.T) {
	we := asserter.Using(t)
	empty := map[string]int{}
	hasTwo := map[string]int{ "foo": 1, "bar": 2 }
	we.CheckThat(empty, Is(Empty()))
	we.CheckThat(hasTwo, Is(Not(Empty())))
}

