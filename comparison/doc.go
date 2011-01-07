// Copyright 2010 Mick Killianey.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
	Hamcrest Matchers based on the standard comparison operators.
	
	Comparison operators are defined:
		http://golang.org/doc/go_spec.html#Comparison_operators

	Note that objects that are unordered cannot be compared using
	the ordered operators.  Therefore, equality (==) and inequality (>=) 
	do not necessarily imply behavior of the ordered operators
	greater-than(>), greater-than-or-equal-to (>=), less-than (<), or
	less-than-or-equal-to(<=).
	
	For example, whenever one of the following matchers are successful:
		Both(GreaterThanOrEqual(x)).And(LessThanOrEqualTo(x))
		Either(GreaterThan(x)).Or(LessThan(x))
	they imply that their respective equality/inequality matcher *must* match:
		EqualTo(x)
		NotEqualTo(x)
	But the converse does not hold.
*/
package comparison
