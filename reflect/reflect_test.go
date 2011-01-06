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

func Test_ArrayOfType(t *testing.T) {
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

func Test_SliceOfType(t *testing.T) {
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

func Test_MapTypeWithKeyType(t *testing.T) {
	we := asserter.Using(t)
	stringIntMap := map[string]int{"one": 1, "two": 2}
	intStringMap := map[int]string{1: "one", 2: "two"}
	
	we.CheckThat(intStringMap, Is(MapWithKeyType(Anything())))
	we.CheckThat(intStringMap, Is(MapWithKeyType(IntType())))
	we.CheckThat(intStringMap, Is(Not(MapWithKeyType(StringType()))))
	we.CheckThat(stringIntMap, Is(Not(MapWithKeyType(IntType()))))
	we.CheckThat(stringIntMap, Is(MapWithKeyType(StringType())))

	we.CheckThat(intStringMap, ToType(Is(MapTypeWithKeyType(Anything()))))
	we.CheckThat(intStringMap, ToType(Is(MapTypeWithKeyType(IntType()))))
	we.CheckThat(intStringMap, ToType(Is(Not(MapTypeWithKeyType(StringType())))))
	we.CheckThat(stringIntMap, ToType(Is(Not(MapTypeWithKeyType(IntType())))))
	we.CheckThat(stringIntMap, ToType(Is(MapTypeWithKeyType(StringType()))))
}

func Test_MapTypeWithElementType(t *testing.T) {
	we := asserter.Using(t)
	stringIntMap := map[string]int{"one": 1, "two": 2}
	intStringMap := map[int]string{1: "one", 2: "two"}
	
	we.CheckThat(stringIntMap, Is(MapWithElementType(Anything())))
	we.CheckThat(stringIntMap, Is(MapWithElementType(IntType())))
	we.CheckThat(stringIntMap, Is(Not(MapWithElementType(StringType()))))
	we.CheckThat(intStringMap, Is(Not(MapWithElementType(IntType()))))
	we.CheckThat(intStringMap, Is(MapWithElementType(StringType())))

	we.CheckThat(stringIntMap, ToType(Is(MapTypeWithElementType(Anything()))))
	we.CheckThat(stringIntMap, ToType(Is(MapTypeWithElementType(IntType()))))
	we.CheckThat(stringIntMap, ToType(Is(Not(MapTypeWithElementType(StringType())))))
	we.CheckThat(intStringMap, ToType(Is(Not(MapTypeWithElementType(IntType())))))
	we.CheckThat(intStringMap, ToType(Is(MapTypeWithElementType(StringType()))))
}

