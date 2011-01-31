// Copyright 2010 Mick Killianey.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package base

import (
	"reflect"
)


type _Comparison int8

const (
	_LESS_THAN = _Comparison(iota)
	_ORDERED_EQUAL_TO
	_GREATER_THAN
	_UNORDERED_EQUAL_TO
	_UNORDERED_NOT_EQUAL_TO
	_INCOMPARABLE_TYPES
)

func (c _Comparison) _Describe(x, y interface{}) SelfDescribing {
	switch c {
	case _LESS_THAN:
		return Description("%v was less than %v", x, y)
	case _ORDERED_EQUAL_TO:
		return Description("%v was equal to %v", x, y)
	case _GREATER_THAN:
		return Description("%v was greater than %v", x, y)
	case _UNORDERED_EQUAL_TO:
		return Description("%v was (unordered) equal to %v", x, y)
	case _UNORDERED_NOT_EQUAL_TO:
		return Description("%v was (unordered) not equal to %v", x, y)
	case _INCOMPARABLE_TYPES:
		return Description("types %T and %T cannot be compared", x, y)
	}
	return Description("Unrecognized %v on %v and %v",int8(c), x, y)
}

// All the messiness of determining whether or not two objects
// can be compared (and if they can be, what the result is).
func _Compare(x interface{}, y interface{}) _Comparison {
	switch i := x.(type) {
	case int:
		switch j := y.(type) {
		case int:
			if i < j { return _LESS_THAN
			} else if i > j { return _GREATER_THAN
			} else { return _ORDERED_EQUAL_TO
			}
		default: return _INCOMPARABLE_TYPES
		}
	case int8:
		switch j := y.(type) {
		case int8:
			if i < j { return _LESS_THAN
			} else if i > j { return _GREATER_THAN
			} else { return _ORDERED_EQUAL_TO
			}
		default: return _INCOMPARABLE_TYPES
		}
	case int16:
		switch j := y.(type) {
		case int16:
			if i < j { return _LESS_THAN
			} else if i > j { return _GREATER_THAN
			} else { return _ORDERED_EQUAL_TO
			}
		default: return _INCOMPARABLE_TYPES
		}
	case int32:
		switch j := y.(type) {
		case int32:
			if i < j { return _LESS_THAN
			} else if i > j { return _GREATER_THAN
			} else { return _ORDERED_EQUAL_TO
			}
		default: return _INCOMPARABLE_TYPES
		}
	case int64:
		switch j := y.(type) {
		case int64:
			if i < j { return _LESS_THAN
			} else if i > j { return _GREATER_THAN
			} else { return _ORDERED_EQUAL_TO
			}
		default: return _INCOMPARABLE_TYPES
		}
	case uint:
		switch j := y.(type) {
		case uint:
			if i < j { return _LESS_THAN
			} else if i > j { return _GREATER_THAN
			} else { return _ORDERED_EQUAL_TO
			}
		default: return _INCOMPARABLE_TYPES
		}
	case uint8:
		switch j := y.(type) {
		case uint8:
			if i < j { return _LESS_THAN
			} else if i > j { return _GREATER_THAN
			} else { return _ORDERED_EQUAL_TO
			}
		default: return _INCOMPARABLE_TYPES
		}
	case uint16:
		switch j := y.(type) {
		case uint16:
			if i < j { return _LESS_THAN
			} else if i > j { return _GREATER_THAN
			} else { return _ORDERED_EQUAL_TO
			}
		default: return _INCOMPARABLE_TYPES
		}
	case uint32:
		switch j := y.(type) {
		case uint32:
			if i < j { return _LESS_THAN
			} else if i > j { return _GREATER_THAN
			} else { return _ORDERED_EQUAL_TO
			}
		default: return _INCOMPARABLE_TYPES
		}
	case uint64:
		switch j := y.(type) {
		case uint64:
			if i < j { return _LESS_THAN
			} else if i > j { return _GREATER_THAN
			} else { return _ORDERED_EQUAL_TO
			}
		default: return _INCOMPARABLE_TYPES
		}
	case float32:
		switch j := y.(type) {
		case float32:
			if i < j { return _LESS_THAN
			} else if i > j { return _GREATER_THAN
			} else if i == j { return _ORDERED_EQUAL_TO
			} else { return _UNORDERED_NOT_EQUAL_TO
			}
		default: return _INCOMPARABLE_TYPES
		}
	case float64:
		switch j := y.(type) {
		case float64:
			if i < j { return _LESS_THAN
			} else if i > j { return _GREATER_THAN
			} else if i == j { return _ORDERED_EQUAL_TO
			} else { return _UNORDERED_NOT_EQUAL_TO
			}
		default: return _INCOMPARABLE_TYPES
		}
	case string:
		switch j := y.(type) {
		case string:
			if i < j { return _LESS_THAN
			} else if i > j { return _GREATER_THAN
			} else if i == j { return _ORDERED_EQUAL_TO
			} else { return _UNORDERED_NOT_EQUAL_TO
			}
		}
	default:
		if reflect.Typeof(x) == reflect.Typeof(y) {
			if x == y {
				return _UNORDERED_EQUAL_TO
			}
			return _UNORDERED_NOT_EQUAL_TO
		}
	}
	return _INCOMPARABLE_TYPES
}

// Returns a matcher that matches values that are greater-than the given
// expected value, using the greater-than (<) operator.
func GreaterThan(expected interface{}) *Matcher {
	match := func(actual interface{}) *Result {
		c := _Compare(actual, expected)
		switch c {
		case _GREATER_THAN:
			return NewResult(true, c._Describe(actual, expected))
		default:
			return NewResult(false, c._Describe(actual, expected))
		}
		panic("every case should have a return")
	}
	return NewMatcherf(match, "GreaterThan(%v)", expected)
}

// Returns a matcher that matches values that are greater-than-or-equal-to
// the given expected value, using the greater-than-or-equal-to (>=) operator.
func GreaterThanOrEqualTo(expected interface{}) *Matcher {
	match := func(actual interface{}) *Result {
		c := _Compare(actual, expected)
		switch c {
		case _GREATER_THAN, _ORDERED_EQUAL_TO:
			return NewResult(true, c._Describe(actual, expected))
		default:
			return NewResult(false, c._Describe(actual, expected))
		}
		panic("every case should have a return")
	}
	return NewMatcherf(match, "GreaterThanOrEqualTo(%v)", expected)
}

// Returns a matcher that matches values that are less-than the given
// expected value, using the less-than (<) operator.
func LessThan(expected interface{}) *Matcher {
	match := func(actual interface{}) *Result {
		c := _Compare(actual, expected)
		switch c {
		case _LESS_THAN:
			return NewResult(true, c._Describe(actual, expected))
		default:
			return NewResult(false, c._Describe(actual, expected))
		}
		panic("every case should have a return")
	}
	return NewMatcherf(match, "LessThan(%v)", expected)
}

// Returns a matcher that matches values that are less-than-or-equal-to
// the given expected value, using the less-than-or-equal-to (<=) operator.
func LessThanOrEqualTo(expected interface{}) *Matcher {
	match := func(actual interface{}) *Result {
		c := _Compare(actual, expected)
		switch c {
		case _LESS_THAN, _ORDERED_EQUAL_TO:
			return NewResult(true, c._Describe(actual, expected))
		default:
			return NewResult(false, c._Describe(actual, expected))
		}
		panic("every case should have a return")
	}
	return NewMatcherf(match, "LessThanOrEqualTo(%v)", expected)
}

// Returns a matcher that matches values that are equal to the
// given expected value, using the equality (==) operator.
func EqualTo(expected interface{}) *Matcher {
	match := func(actual interface{}) *Result {
		c := _Compare(actual, expected)
		switch c {
		case _ORDERED_EQUAL_TO, _UNORDERED_EQUAL_TO:
			return NewResult(true, c._Describe(actual, expected))
		default:
			return NewResult(false, c._Describe(actual, expected))
		}
		panic("every case should have a return")
	}
	return NewMatcherf(match, "EqualTo(%v)", expected)
}

// Returns a matcher that matches values that are not equal to the
// given expected value, using the inequality (!=) operator.
func NotEqualTo(expected interface{}) *Matcher {
	match := func(actual interface{}) *Result {
		c := _Compare(actual, expected)
		switch c {
		case _ORDERED_EQUAL_TO, _UNORDERED_EQUAL_TO:
			return NewResult(false, c._Describe(actual, expected))
		default:
			return NewResult(true, c._Describe(actual, expected))
		}
		panic("every case should have a return")
	}
	return NewMatcherf(match, "NotEqualTo(%v)", expected)
}

