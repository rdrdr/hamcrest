// Copyright 2010 Mick Killianey.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package comparison

import (
	"hamcrest"
)


// Returns a matcher that matches values that are greater-than the given
// expected value, using the greater-than (<) operator.
func GreaterThan(expected interface{}) *hamcrest.Matcher {
	return hamcrest.GreaterThan(expected)
}

// Returns a matcher that matches values that are greater-than-or-equal-to
// the given expected value, using the greater-than-or-equal-to (>=) operator.
func GreaterThanOrEqualTo(expected interface{}) *hamcrest.Matcher {
	return hamcrest.GreaterThanOrEqualTo(expected)
}

// Returns a matcher that matches values that are less-than the given
// expected value, using the less-than (<) operator.
func LessThan(expected interface{}) *hamcrest.Matcher {
	return hamcrest.LessThan(expected)
}

// Returns a matcher that matches values that are less-than-or-equal-to
// the given expected value, using the less-than-or-equal-to (<=) operator.
func LessThanOrEqualTo(expected interface{}) *hamcrest.Matcher {
	return hamcrest.LessThanOrEqualTo(expected)
}

// Returns a matcher that matches values that are equal to the
// given expected value, using the equality (==) operator.
func EqualTo(expected interface{}) *hamcrest.Matcher {
	return hamcrest.EqualTo(expected)
}

// Returns a matcher that matches values that are not equal to the
// given expected value, using the inequality (!=) operator.
func NotEqualTo(expected interface{}) *hamcrest.Matcher {
	return hamcrest.NotEqualTo(expected)
}

