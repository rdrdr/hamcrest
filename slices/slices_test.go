// Copyright 2011 Mick Killianey.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slices

import (
	"github.com/rdrdr/hamcrest/asserter"
	. "github.com/rdrdr/hamcrest/core"
	"github.com/rdrdr/hamcrest/base"
	"testing"
)

var Matched = base.Matched()
var DidNotMatch = base.DidNotMatch()

func Test_AnyElem(t *testing.T) {
	we := asserter.Using(t)
	we.CheckThat(AnyElem(EqualTo(1)).Match([]int{1, 2, 3}),
		Matched.Comment("first element"))
	we.CheckThat(AnyElem(EqualTo(2)).Match([]int{1, 2, 3}),
		Matched.Comment("middle element"))
	we.CheckThat(AnyElem(EqualTo(3)).Match([]int{1, 2, 3}),
		Matched.Comment("last element"))
	we.CheckThat(AnyElem(EqualTo(4)).Match([]int{1, 2, 3}),
		DidNotMatch.Comment("None of the three match"))
	we.CheckThat(AnyElem(Anything()).Match([]int{}),
		DidNotMatch.Comment("no elements"))
}

func Test_EachElem(t *testing.T) {
	we := asserter.Using(t)
	we.CheckThat(EachElem(NotEqualTo(1)).Match([]int{1, 2, 3}),
		DidNotMatch.Comment("all but first element"))
	we.CheckThat(EachElem(NotEqualTo(2)).Match([]int{1, 2, 3}),
		DidNotMatch.Comment("all but middle element"))
	we.CheckThat(EachElem(NotEqualTo(3)).Match([]int{1, 2, 3}),
		DidNotMatch.Comment("all but last element"))
	we.CheckThat(EachElem(NotEqualTo(4)).Match([]int{1, 2, 3}),
		Matched.Comment("all match"))
	we.CheckThat(EachElem(Anything()).Match([]int{}),
		Matched.Comment("no elements"))
}

func Test_ToLen(t *testing.T) {
	we := asserter.Using(t)
	IsLength2 := ToLen(Is(EqualTo(2)))
	we.CheckThat(IsLength2.Match([]int{}), DidNotMatch.Comment("no elements"))
	we.CheckThat(IsLength2.Match([]int{7}), DidNotMatch.Comment("too few"))
	we.CheckThat(IsLength2.Match([]int{7, 8}), Matched.Comment("just right"))
	we.CheckThat(IsLength2.Match([]int{7, 8, 9}), DidNotMatch.Comment("too many"))
}

func Test_Empty(t *testing.T) {
	we := asserter.Using(t)
	we.CheckThat(Empty().Match([]int{}), Matched.Comment("no ints"))
	we.CheckThat(Empty().Match([]string{}), Matched.Comment("no strings"))
	we.CheckThat(Empty().Match([]string{"not", "empty"}), DidNotMatch)
}
