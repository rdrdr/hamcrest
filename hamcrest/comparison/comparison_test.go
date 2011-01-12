// Copyright 2010 Mick Killianey.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package comparison

import (
	"hamcrest"
	"hamcrest/asserter"
	"testing"
)

var Not = hamcrest.Not

func Test_GreaterThan(t *testing.T) {
	we := asserter.Using(t)
	we.CheckThat(3, GreaterThan(2))
	we.CheckThat(3, Not(GreaterThan(3)))
	we.CheckThat(3, Not(GreaterThan(4)))
	
	we.CheckThat(3, Not(GreaterThan("2")))
	we.CheckThat("3", Not(GreaterThan(2)))
	we.CheckThat(t, Not(GreaterThan(t)))
	
	we.CheckThat(nil, Not(GreaterThan(3)))
	we.CheckThat(3, Not(GreaterThan(nil)))
}

func Test_GreaterThanOrEqualTo(t *testing.T) {
	we := asserter.Using(t)
	
	we.CheckThat(3, GreaterThanOrEqualTo(2))
	we.CheckThat(3, GreaterThanOrEqualTo(3))
	we.CheckThat(3, Not(GreaterThanOrEqualTo(4)))
	
	we.CheckThat(3, Not(GreaterThanOrEqualTo("2")))
	we.CheckThat("3", Not(GreaterThanOrEqualTo(2)))
	we.CheckThat(t, Not(GreaterThanOrEqualTo(t)))
	
	we.CheckThat(nil, Not(GreaterThanOrEqualTo(3)))
	we.CheckThat(3, Not(GreaterThanOrEqualTo(nil)))
}

func Test_LessThan(t *testing.T) {
	we := asserter.Using(t)
	
	we.CheckThat(3, Not(LessThan(2)))
	we.CheckThat(3, Not(LessThan(3)))
	we.CheckThat(3, LessThan(4))
	
	we.CheckThat(3, Not(LessThan("4")))
	we.CheckThat("3", Not(LessThan(4)))
	we.CheckThat(t, Not(LessThan(t)))
	
	we.CheckThat(nil, Not(LessThan(3)))
	we.CheckThat(3, Not(LessThan(nil)))
}

func Test_LessThanOrEqualTo(t *testing.T) {
	we := asserter.Using(t)
	
	we.CheckThat(3, Not(LessThanOrEqualTo(2)))
	we.CheckThat(3, LessThanOrEqualTo(3))
	we.CheckThat(3, LessThanOrEqualTo(4))
	
	we.CheckThat(3, Not(LessThanOrEqualTo("4")))
	we.CheckThat("3", Not(LessThanOrEqualTo(4)))
	we.CheckThat(t, Not(LessThanOrEqualTo(t)))
	
	we.CheckThat(nil, Not(LessThanOrEqualTo(3)))
	we.CheckThat(3, Not(LessThanOrEqualTo(nil)))
}

func Test_EqualTo(t *testing.T) {
	we := asserter.Using(t)
	
	we.CheckThat(3, Not(EqualTo(2)))
	we.CheckThat(3, EqualTo(3))
	
	we.CheckThat(3, Not(EqualTo("3")))
	we.CheckThat(int(3), Not(EqualTo(uint(3))))
	
	we.CheckThat(3, Not(EqualTo(nil)))
	we.CheckThat(nil, Not(EqualTo(3)))
}

func Test_NotEqualTo(t *testing.T) {
	we := asserter.Using(t)
	
	we.CheckThat(3, NotEqualTo(2))
	we.CheckThat(3, Not(NotEqualTo(3)))
	
	we.CheckThat(3, NotEqualTo("3"))
	we.CheckThat(int(3), NotEqualTo(uint(3)))
	
	we.CheckThat(3, NotEqualTo(nil))
	we.CheckThat(nil, NotEqualTo(3))
}

func checkOrderingOfOneAndTwo(we asserter.Asserter, one interface{}, two interface{}) {
	we.CheckThat(one, LessThan(two))
	we.CheckThat(one, Not(LessThan(one)))
	we.CheckThat(two, Not(LessThan(one)))

	we.CheckThat(one, LessThanOrEqualTo(two))
	we.CheckThat(one, LessThanOrEqualTo(one))
	we.CheckThat(two, Not(LessThanOrEqualTo(one)))

	we.CheckThat(one, Not(GreaterThan(two)))
	we.CheckThat(one, Not(GreaterThan(one)))
	we.CheckThat(two, GreaterThan(one))

	we.CheckThat(one, Not(GreaterThanOrEqualTo(two)))
	we.CheckThat(one, GreaterThanOrEqualTo(one))
	we.CheckThat(two, GreaterThanOrEqualTo(one))
	
	we.CheckThat(one, EqualTo(one))
	we.CheckThat(one, Not(EqualTo(two)))
	we.CheckThat(one, NotEqualTo(two))
	we.CheckThat(one, Not(NotEqualTo(one)))
}

func TestOrderingOfTypes(t *testing.T) {
	we := asserter.Using(t)
	checkOrderingOfOneAndTwo(we, int(1), int(2))
	checkOrderingOfOneAndTwo(we, int8(1), int8(2))
	checkOrderingOfOneAndTwo(we, int16(1), int16(2))
	checkOrderingOfOneAndTwo(we, int32(1), int32(2))
	checkOrderingOfOneAndTwo(we, int64(1), int64(2))
	checkOrderingOfOneAndTwo(we, uint(1), uint(2))
	checkOrderingOfOneAndTwo(we, uint8(1), uint8(2))
	checkOrderingOfOneAndTwo(we, uint16(1), uint16(2))
	checkOrderingOfOneAndTwo(we, uint32(1), uint32(2))
	checkOrderingOfOneAndTwo(we, uint64(1), uint64(2))
	checkOrderingOfOneAndTwo(we, float(1), float(2))
	checkOrderingOfOneAndTwo(we, float32(1), float32(2))
	checkOrderingOfOneAndTwo(we, float64(1), float64(2))
	checkOrderingOfOneAndTwo(we, string("1"), string("2"))
}



