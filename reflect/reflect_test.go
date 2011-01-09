// Copyright 2011 Mick Killianey.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package reflect

import (
	"hamcrest"
	"hamcrest/asserter"
	//"reflect"
	//"runtime"
	"testing"
)

var EqualTo = hamcrest.EqualTo
var Anything = hamcrest.Anything
var Is = hamcrest.Is
var Not = hamcrest.Not

func Test_Bool_And_BoolType(t *testing.T) {
	we := asserter.Using(t)
	checkMatch := func(v interface{}) {
		we.CheckThat(v, Is(Bool()))
		we.CheckThat(v, ToType(Is(BoolType())))
	}
	checkNonMatch := func(v interface{}) {
		we.CheckThat(v, Is(Not(Bool())))
		we.CheckThat(v, ToType(Is(Not(BoolType()))))
	}
	checkMatch(true)
	checkNonMatch("true")
	checkMatch(false)
	checkNonMatch("false")
	checkNonMatch(nil)
	checkNonMatch(interface{}(nil))
}

func Test_Int(t *testing.T) {
	we := asserter.Using(t)
	checkMatch := func(v interface{}) {
		we.CheckThat(v, Is(Int()))
		we.CheckThat(v, ToType(Is(IntType())))
	}
	checkNonMatch := func(v interface{}) {
		we.CheckThat(v, Is(Not(Int())))
		we.CheckThat(v, ToType(Is(Not(IntType()))))
	}
	checkMatch(1)
	checkNonMatch("1")
	checkMatch(0)
	checkNonMatch(0.0)
	checkNonMatch(0i)
	checkMatch(-1)
	checkNonMatch(nil)
	checkNonMatch(true)
	checkNonMatch(interface{}(nil))
}

func Test_Int8(t *testing.T) {
	we := asserter.Using(t)
	we.CheckThat(int8(1), Is(Int8()))
	we.CheckThat(int8(1), ToType(Is(Int8Type())))
	we.CheckThat(int(1), Is(Not(Int8())))
	we.CheckThat(int(1), ToType(Is(Not(Int8Type()))))
}

func Test_Int16(t *testing.T) {
	we := asserter.Using(t)
	we.CheckThat(int16(1), Is(Int16()))
	we.CheckThat(int16(1), ToType(Is(Int16Type())))
	we.CheckThat(int(1), Is(Not(Int16())))
	we.CheckThat(int(1), ToType(Is(Not(Int16Type()))))
}

func Test_Int32(t *testing.T) {
	we := asserter.Using(t)
	we.CheckThat(int32(1), Is(Int32()))
	we.CheckThat(int32(1), ToType(Is(Int32Type())))
	we.CheckThat(int(1), Is(Not(Int32())))
	we.CheckThat(int(1), ToType(Is(Not(Int32Type()))))
}

func Test_Int64(t *testing.T) {
	we := asserter.Using(t)
	we.CheckThat(int64(1), Is(Int64()))
	we.CheckThat(int64(1), ToType(Is(Int64Type())))
	we.CheckThat(int(1), Is(Not(Int64())))
	we.CheckThat(int(1), ToType(Is(Not(Int64Type()))))
}

func Test_Uint(t *testing.T) {
	we := asserter.Using(t)
	we.CheckThat(uint(1), Is(Uint()))
	we.CheckThat(uint(1), ToType(Is(UintType())))
	we.CheckThat(int(1), Is(Not(Uint())))
	we.CheckThat(int(1), ToType(Is(Not(UintType()))))
}

func Test_Uint8(t *testing.T) {
	we := asserter.Using(t)
	we.CheckThat(uint8(1), Is(Uint8()))
	we.CheckThat(uint8(1), ToType(Is(Uint8Type())))
	we.CheckThat(uint(1), Is(Not(Uint8())))
	we.CheckThat(uint(1), ToType(Is(Not(Uint8Type()))))
}

func Test_Uint16(t *testing.T) {
	we := asserter.Using(t)
	we.CheckThat(uint16(1), Is(Uint16()))
	we.CheckThat(uint16(1), ToType(Is(Uint16Type())))
	we.CheckThat(uint(1), Is(Not(Uint16())))
	we.CheckThat(uint(1), ToType(Is(Not(Uint16Type()))))
}

func Test_Uint32(t *testing.T) {
	we := asserter.Using(t)
	we.CheckThat(uint32(1), Is(Uint32()))
	we.CheckThat(uint32(1), ToType(Is(Uint32Type())))
	we.CheckThat(uint(1), Is(Not(Uint32())))
	we.CheckThat(uint(1), ToType(Is(Not(Uint32Type()))))
}

func Test_Uint64(t *testing.T) {
	we := asserter.Using(t)
	we.CheckThat(uint64(1), Is(Uint64()))
	we.CheckThat(uint64(1), ToType(Is(Uint64Type())))
	we.CheckThat(uint(1), Is(Not(Uint64())))
	we.CheckThat(uint(1), ToType(Is(Not(Uint64Type()))))
}

func Test_Float(t *testing.T) {
	we := asserter.Using(t)
	checkMatch := func(v interface{}) {
		we.CheckThat(v, Is(Float()))
		we.CheckThat(v, ToType(Is(FloatType())))
	}
	checkNonMatch := func(v interface{}) {
		we.CheckThat(v, Is(Not(Float())))
		we.CheckThat(v, ToType(Is(Not(FloatType()))))
	}
	checkMatch(1.0)
	checkNonMatch("1.0")
	checkNonMatch(0)
	checkMatch(0.0)
	checkNonMatch(0i)
	checkNonMatch(nil)
	checkNonMatch(true)
	checkNonMatch(interface{}(nil))
}

func Test_Float32(t *testing.T) {
	we := asserter.Using(t)
	we.CheckThat(float32(1.0), Is(Float32()))
	we.CheckThat(float32(1.0), ToType(Is(Float32Type())))
	we.CheckThat(float(1.0), Is(Not(Float32())))
	we.CheckThat(float(1.0), ToType(Is(Not(Float32Type()))))
}

func Test_Float64(t *testing.T) {
	we := asserter.Using(t)
	we.CheckThat(float64(1.0), Is(Float64()))
	we.CheckThat(float64(1.0), ToType(Is(Float64Type())))
	we.CheckThat(float(1.0), Is(Not(Float64())))
	we.CheckThat(float(1.0), ToType(Is(Not(Float64Type()))))
}

func Test_Complex(t *testing.T) {
	we := asserter.Using(t)
	checkMatch := func(v interface{}) {
		we.CheckThat(v, Is(Complex()))
		we.CheckThat(v, ToType(Is(ComplexType())))
	}
	checkNonMatch := func(v interface{}) {
		we.CheckThat(v, Is(Not(Complex())))
		we.CheckThat(v, ToType(Is(Not(ComplexType()))))
	}
	checkMatch(0i)
	checkMatch(0.0i)
	checkMatch(1i)
	checkMatch(1.0i)
	checkMatch(1+1i)
	checkMatch(1.0+1.0i)
	checkNonMatch("1+0i")
	checkNonMatch(1)
	checkNonMatch(1.0)
	checkNonMatch(nil)
	checkNonMatch(true)
	checkNonMatch(interface{}(nil))
}

func Test_Complex64(t *testing.T) {
	we := asserter.Using(t)
	we.CheckThat(complex64(1.0), Is(Complex64()))
	we.CheckThat(complex64(1.0), ToType(Is(Complex64Type())))
	we.CheckThat(complex(1.0), Is(Not(Complex64())))
	we.CheckThat(complex(1.0), ToType(Is(Not(Complex64Type()))))
}

func Test_Complex128(t *testing.T) {
	we := asserter.Using(t)
	we.CheckThat(complex128(1.0i), Is(Complex128()))
	we.CheckThat(complex128(1.0i), ToType(Is(Complex128Type())))
	we.CheckThat(complex(1.0i), Is(Not(Complex128())))
	we.CheckThat(complex(1.0i), ToType(Is(Not(Complex128Type()))))
}

func Test_String(t *testing.T) {
	we := asserter.Using(t)
	checkMatch := func(v interface{}) {
		we.CheckThat(v, Is(String()))
		we.CheckThat(v, ToType(Is(StringType())))
	}
	checkNonMatch := func(v interface{}) {
		we.CheckThat(v, Is(Not(String())))
		we.CheckThat(v, ToType(Is(Not(StringType()))))
	}
	checkMatch("")
	checkMatch("non-empty")
	checkNonMatch(nil)
	checkNonMatch(true)
	checkNonMatch(1)
	checkNonMatch(1.0)
	checkNonMatch(1i)
	checkNonMatch(interface{}(nil))
}

func Test_ArrayTypeOf(t *testing.T) {
	we := asserter.Using(t)
	boolArray := [2]bool {true, false}
	intArray := [3]int {1, 2, 3}
	intArrayArray := [2][3]int { {1, 2, 3}, {4, 5, 6} }
	
	we.CheckThat(boolArray, Is(ArrayOf(Anything())))
	we.CheckThat(boolArray, Is(ArrayOf(BoolType())))
	we.CheckThat(boolArray, Is(Not(ArrayOf(IntType()))))
	we.CheckThat(boolArray, ToType(Is(ArrayTypeOf(Anything()))))
	we.CheckThat(boolArray, ToType(Is(ArrayTypeOf(BoolType()))))
	we.CheckThat(boolArray, ToType(Is(Not(ArrayTypeOf(IntType())))))
	
	we.CheckThat(intArray, Is(ArrayOf(Anything())))
	we.CheckThat(intArray, Is(Not(ArrayOf(BoolType()))))
	we.CheckThat(intArray, Is(ArrayOf(IntType())))
	we.CheckThat(intArray, ToType(Is(ArrayTypeOf(Anything()))))
	we.CheckThat(intArray, ToType(Is(Not(ArrayTypeOf(BoolType())))))
	we.CheckThat(intArray, ToType(Is(ArrayTypeOf(IntType()))))
	
	we.CheckThat(intArrayArray, Is(ArrayOf(Anything())))
	we.CheckThat(intArrayArray, Is(Not(ArrayOf(IntType()))))
	we.CheckThat(intArrayArray, Is(Not(ArrayOf(ArrayOf(IntType())))))
	we.CheckThat(intArrayArray, Is(Not(ArrayTypeOf(ArrayOf(IntType())))))
	we.CheckThat(intArrayArray, Is(ArrayOf(ArrayTypeOf(IntType()))))
	we.CheckThat(intArrayArray, ToType(Is(Not(ArrayOf(Anything())))))
	we.CheckThat(intArrayArray, ToType(Is(ArrayTypeOf(ArrayTypeOf(IntType())))))
	
	intSlice := make([]int, 0, 1)
	we.CheckThat(intSlice, Is(Not(ArrayOf(Anything()))))
}

func Test_ChanOf(t *testing.T) {
	we := asserter.Using(t)
	intChan := make(chan int, 1)
	intChanIn := func(ch chan int) chan<- int { return ch }(intChan)
	intChanOut := func(ch chan int) <-chan int { return ch }(intChan)
	we.CheckThat(intChan, Is(ChannelOf(Anything())))
	we.CheckThat(intChan, Is(ChannelOf(IntType())))
	we.CheckThat(intChan, Is(Not(ChannelOf(StringType()))))
	
	we.CheckThat(intChanIn, Is(ChannelOf(Anything())))
	we.CheckThat(intChanIn, Is(ChannelOf(IntType())))
	we.CheckThat(intChanIn, Is(Not(ChannelOf(StringType()))))
	
	we.CheckThat(intChanOut, Is(ChannelOf(Anything())))
	we.CheckThat(intChanOut, Is(ChannelOf(IntType())))
	we.CheckThat(intChanOut, Is(Not(ChannelOf(StringType()))))
}

func Test_ChanTypeOf(t *testing.T) {
	we := asserter.Using(t)
	intChan := make(chan int, 1)
	intChanIn := func(ch chan int) chan<- int { return ch }(intChan)
	intChanOut := func(ch chan int) <-chan int { return ch }(intChan)
	we.CheckThat(intChan, ToType(Is(ChannelTypeOf(Anything()))))
	we.CheckThat(intChan, ToType(Is(ChannelTypeOf(IntType()))))
	we.CheckThat(intChan, ToType(Is(Not(ChannelTypeOf(StringType())))))
	
	we.CheckThat(intChanIn, ToType(Is(ChannelTypeOf(Anything()))))
	we.CheckThat(intChanIn, ToType(Is(ChannelTypeOf(IntType()))))
	we.CheckThat(intChanIn, ToType(Is(Not(ChannelTypeOf(StringType())))))
	
	we.CheckThat(intChanOut, ToType(Is(ChannelTypeOf(Anything()))))
	we.CheckThat(intChanOut, ToType(Is(ChannelTypeOf(IntType()))))
	we.CheckThat(intChanOut, ToType(Is(Not(ChannelTypeOf(StringType())))))
}

func Test_SliceTypeOf(t *testing.T) {
	we := asserter.Using(t)
	boolSlice := make([]bool, 0, 1)
	intSlice := make([]int, 0, 1)
	intSliceSlice := make([][]int, 0, 1)
	
	we.CheckThat(boolSlice, Is(SliceOf(Anything())))
	we.CheckThat(boolSlice, Is(SliceOf(BoolType())))
	we.CheckThat(boolSlice, Is(Not(SliceOf(IntType()))))
	we.CheckThat(boolSlice, ToType(Is(SliceTypeOf(Anything()))))
	we.CheckThat(boolSlice, ToType(Is(SliceTypeOf(BoolType()))))
	we.CheckThat(boolSlice, ToType(Is(Not(SliceTypeOf(IntType())))))
	
	we.CheckThat(intSlice, Is(SliceOf(Anything())))
	we.CheckThat(intSlice, Is(Not(SliceOf(BoolType()))))
	we.CheckThat(intSlice, Is(SliceOf(IntType())))
	we.CheckThat(intSlice, ToType(Is(SliceTypeOf(Anything()))))
	we.CheckThat(intSlice, ToType(Is(Not(SliceTypeOf(BoolType())))))
	we.CheckThat(intSlice, ToType(Is(SliceTypeOf(IntType()))))
	
	we.CheckThat(intSliceSlice, Is(SliceOf(Anything())))
	we.CheckThat(intSliceSlice, Is(Not(SliceOf(IntType()))))
	we.CheckThat(intSliceSlice, Is(Not(SliceOf(SliceOf(IntType())))))
	we.CheckThat(intSliceSlice, Is(Not(SliceTypeOf(SliceOf(IntType())))))
	we.CheckThat(intSliceSlice, Is(SliceOf(SliceTypeOf(IntType()))))
	we.CheckThat(intSliceSlice, ToType(Is(Not(SliceOf(Anything())))))
	we.CheckThat(intSliceSlice, ToType(Is(SliceTypeOf(SliceTypeOf(IntType())))))
	
	var intArray = [3]int{1, 2, 3}
	we.CheckThat(intArray, Is(Not(SliceOf(Anything()))))
	we.CheckThat(intArray, ToType(Is(Not(SliceTypeOf(Anything())))))
}

func Test_MapOf(t *testing.T) {
	we := asserter.Using(t)
	intStringMap := map[int]string{1: "one", 2: "two"}

	we.CheckThat(intStringMap, Is(MapOf(Anything(), Anything())))
	we.CheckThat(intStringMap, Is(MapOf(IntType(), StringType())))
	we.CheckThat(intStringMap, Is(Not(MapOf(StringType(), StringType()))))
	we.CheckThat(intStringMap, Is(Not(MapOf(IntType(), IntType()))))
	we.CheckThat(intStringMap, Is(Not(MapOf(StringType(), IntType()))))
}

func Test_MapTypeOf(t *testing.T) {
	we := asserter.Using(t)
	intStringMap := map[int]string{1: "one", 2: "two"}

	we.CheckThat(intStringMap, ToType(Is(MapTypeOf(Anything(), Anything()))))
	we.CheckThat(intStringMap, ToType(Is(MapTypeOf(IntType(), StringType()))))
	we.CheckThat(intStringMap, ToType(Is(Not(MapTypeOf(StringType(), StringType())))))
	we.CheckThat(intStringMap, ToType(Is(Not(MapTypeOf(IntType(), IntType())))))
	we.CheckThat(intStringMap, ToType(Is(Not(MapTypeOf(StringType(), IntType())))))
}

func Test_PtrTypeTo(t *testing.T) {
	we := asserter.Using(t)
	intObj := 1
	intPtr := &intObj
	intPtrPtr := &intPtr
	
	we.CheckThat(intObj, Is(Not(PtrTo(Anything()))))

	we.CheckThat(intPtr, Is(PtrTo(Anything())))
	we.CheckThat(intPtr, Is(PtrTo(IntType())))
	we.CheckThat(intPtr, Is(Not(PtrTo(StringType()))))

	we.CheckThat(intPtr, ToType(Is(PtrTypeTo(Anything()))))
	we.CheckThat(intPtr, ToType(Is(PtrTypeTo(IntType()))))
	we.CheckThat(intPtr, ToType(Is(Not(PtrTypeTo(StringType())))))

	we.CheckThat(intPtrPtr, ToType(Is(PtrTypeTo(Anything()))))
	we.CheckThat(intPtrPtr, ToType(Is(Not(PtrTypeTo(IntType())))))
	we.CheckThat(intPtrPtr, ToType(Is(PtrTypeTo(PtrTypeTo(IntType())))))
}

