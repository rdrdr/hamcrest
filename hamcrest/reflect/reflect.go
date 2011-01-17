// Copyright 2011 Mick Killianey.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package reflect

import (
	"hamcrest"
	"reflect"
)

// Returns a new matcher that applies the type of its input
// element to the given matcher.
func ToType(matcher *hamcrest.Matcher) *hamcrest.Matcher {
	match := func(actual interface{}) *hamcrest.Result {
		actualType := reflect.Typeof(actual)
		result := matcher.Match(actualType)
		return hamcrest.NewResultf(result.Matched(),
			"reflect.Typeof() returned %v", actualType).
			WithCauses(result)
	}
	return hamcrest.NewMatcherf(match, "ToType(%v)", matcher)
}

// Returns a matcher that matches any object with the same
// type as the given example.
func SameTypeAs(example interface{}) *hamcrest.Matcher {
	exampleType := reflect.Typeof(example)
	return _TypeMatcher(exampleType.Name(), exampleType)
}

func _TypeMatcher(name string, expectedType reflect.Type) *hamcrest.Matcher {
	match := func(actual interface{}) *hamcrest.Result {
		if actual == nil {
			return hamcrest.NewResultf(false, "was nil")
		}
		actualType := reflect.Typeof(actual)
		if reflect.DeepEqual(actualType, expectedType) {
			return hamcrest.NewResultf(true, "was of type %v", expectedType)
		}
		return hamcrest.NewResultf(false,
			"was a %v, not a %v", actualType, expectedType)
	}
	return hamcrest.NewMatcherf(match, "Typeof[%v]", expectedType)
}

var (
	boolType = reflect.Typeof(false)
	boolMatcher = _TypeMatcher("Bool", boolType)
	boolTypeMatcher = hamcrest.DeeplyEqualTo(boolType)
	
	intType = reflect.Typeof(int(0))
	intMatcher = _TypeMatcher("Int", intType)
	intTypeMatcher = hamcrest.DeeplyEqualTo(intType)
	
	int8Type = reflect.Typeof(int8(0))
	int8Matcher = _TypeMatcher("Int8", int8Type)
	int8TypeMatcher = hamcrest.DeeplyEqualTo(int8Type)
	
	int16Type = reflect.Typeof(int16(0))
	int16Matcher = _TypeMatcher("Int16", int16Type)
	int16TypeMatcher = hamcrest.DeeplyEqualTo(int16Type)

	int32Type = reflect.Typeof(int32(0))
	int32Matcher = _TypeMatcher("Int32", int32Type)
	int32TypeMatcher = hamcrest.DeeplyEqualTo(int32Type)
	
	int64Type = reflect.Typeof(int64(0))
	int64Matcher = _TypeMatcher("Int64", int64Type)
	int64TypeMatcher = hamcrest.DeeplyEqualTo(int64Type)
	
	uintType = reflect.Typeof(uint(0))
	uintMatcher = _TypeMatcher("Uint", uintType)
	uintTypeMatcher = hamcrest.DeeplyEqualTo(uintType)
	
	uint8Type = reflect.Typeof(uint8(0))
	uint8Matcher = _TypeMatcher("Uint8", uint8Type)
	uint8TypeMatcher = hamcrest.DeeplyEqualTo(uint8Type)
	
	uint16Type = reflect.Typeof(uint16(0))
	uint16Matcher = _TypeMatcher("Uint16", uint16Type)
	uint16TypeMatcher = hamcrest.DeeplyEqualTo(uint16Type)

	uint32Type = reflect.Typeof(uint32(0))
	uint32Matcher = _TypeMatcher("Uint32", uint32Type)
	uint32TypeMatcher = hamcrest.DeeplyEqualTo(uint32Type)
	
	uint64Type = reflect.Typeof(uint64(0))
	uint64Matcher = _TypeMatcher("Uint64", uint64Type)
	uint64TypeMatcher = hamcrest.DeeplyEqualTo(uint64Type)

	uintptrType = reflect.Typeof(uintptr(0))
	uintptrMatcher = _TypeMatcher("Uintptr", uintptrType)
	uintptrTypeMatcher = hamcrest.DeeplyEqualTo(uintptrType)
	
	floatType = reflect.Typeof(float(0))
	floatMatcher = _TypeMatcher("Float", floatType)
	floatTypeMatcher = hamcrest.DeeplyEqualTo(floatType)
	
	float32Type = reflect.Typeof(float32(0))
	float32Matcher = _TypeMatcher("Float32", float32Type)
	float32TypeMatcher = hamcrest.DeeplyEqualTo(float32Type)

	float64Type = reflect.Typeof(float64(0))
	float64Matcher = _TypeMatcher("Float64", float64Type)
	float64TypeMatcher = hamcrest.DeeplyEqualTo(float64Type)

	complexType = reflect.Typeof(complex(0i))
	complexMatcher = _TypeMatcher("Complex", complexType)
	complexTypeMatcher = hamcrest.DeeplyEqualTo(complexType)

	complex64Type = reflect.Typeof(complex64(0i))
	complex64Matcher = _TypeMatcher("Complex64", complex64Type)
	complex64TypeMatcher = hamcrest.DeeplyEqualTo(complex64Type)

	complex128Type = reflect.Typeof(complex128(0i))
	complex128Matcher = _TypeMatcher("Complex128", complex128Type)
	complex128TypeMatcher = hamcrest.DeeplyEqualTo(complex128Type)

	stringType = reflect.Typeof("")
	stringMatcher = _TypeMatcher("String", stringType)
	stringTypeMatcher = hamcrest.DeeplyEqualTo(stringType)
)

func Bool() *hamcrest.Matcher { return boolMatcher }
func Int() *hamcrest.Matcher { return intMatcher }
func Int8() *hamcrest.Matcher { return int8Matcher }
func Int16() *hamcrest.Matcher { return int16Matcher }
func Int32() *hamcrest.Matcher { return int32Matcher }
func Int64() *hamcrest.Matcher { return int64Matcher }
func Uint() *hamcrest.Matcher { return uintMatcher }
func Uint8() *hamcrest.Matcher { return uint8Matcher }
func Uint16() *hamcrest.Matcher { return uint16Matcher }
func Uint32() *hamcrest.Matcher { return uint32Matcher }
func Uint64() *hamcrest.Matcher { return uint64Matcher }
func Float() *hamcrest.Matcher { return floatMatcher }
func Float32() *hamcrest.Matcher { return float32Matcher }
func Float64() *hamcrest.Matcher { return float64Matcher }
func Complex() *hamcrest.Matcher { return complexMatcher }
func Complex64() *hamcrest.Matcher { return complex64Matcher }
func Complex128() *hamcrest.Matcher { return complex128Matcher }
func String() *hamcrest.Matcher { return stringMatcher }

func BoolType() *hamcrest.Matcher { return boolTypeMatcher }
func IntType() *hamcrest.Matcher { return intTypeMatcher }
func Int8Type() *hamcrest.Matcher { return int8TypeMatcher }
func Int16Type() *hamcrest.Matcher { return int16TypeMatcher }
func Int32Type() *hamcrest.Matcher { return int32TypeMatcher }
func Int64Type() *hamcrest.Matcher { return int64TypeMatcher }
func UintType() *hamcrest.Matcher { return uintTypeMatcher }
func Uint8Type() *hamcrest.Matcher { return uint8TypeMatcher }
func Uint16Type() *hamcrest.Matcher { return uint16TypeMatcher }
func Uint32Type() *hamcrest.Matcher { return uint32TypeMatcher }
func Uint64Type() *hamcrest.Matcher { return uint64TypeMatcher }
func FloatType() *hamcrest.Matcher { return floatTypeMatcher }
func Float32Type() *hamcrest.Matcher { return float32TypeMatcher }
func Float64Type() *hamcrest.Matcher { return float64TypeMatcher }
func ComplexType() *hamcrest.Matcher { return complexTypeMatcher }
func Complex64Type() *hamcrest.Matcher { return complex64TypeMatcher }
func Complex128Type() *hamcrest.Matcher { return complex128TypeMatcher }
func StringType() *hamcrest.Matcher { return stringTypeMatcher }


// Returns a new matcher that, on any input that is a *reflect.ArrayType,
// extracts the type of element and matches it against the given matcher.
//
// If the given input is not an *reflect.ArrayType, this fails to match.
// Note:  this matches array *types*, not arrays. (See ArrayOf.)
func ArrayTypeOf(elementTypeMatcher *hamcrest.Matcher) *hamcrest.Matcher {
	match := func(actual interface{}) *hamcrest.Result {
		if arrayType, ok := actual.(*reflect.ArrayType); ok {
			elementType := arrayType.Elem()
			result := elementTypeMatcher.Match(elementType)
			return hamcrest.NewResultf(
				result.Matched(),
				"was ArrayType with elements of type %v", elementType.Name()).
				WithCauses(result)
		}
		return hamcrest.NewResultf(false,
			"was of type %T, not an ArrayType", actual)
	}
	return hamcrest.NewMatcherf(match, "ArrayTypeOf(%v)", elementTypeMatcher)
}

// Returns a new matcher that, on any input that is an array, extracts
// its type and matches it against the given matcher.
//
// If the given input is not an array, this fails to match.
// Note: this matches *arrays*, not array *types*. (See ArrayTypeOf.)
func ArrayOf(elementTypeMatcher *hamcrest.Matcher) *hamcrest.Matcher {
	match := func(actual interface{}) *hamcrest.Result {
		actualType := reflect.Typeof(actual)
		if arrayType, ok := actualType.(*reflect.ArrayType); ok {
			elementType := arrayType.Elem()
			result := elementTypeMatcher.Match(elementType)
			return hamcrest.NewResultf(
				result.Matched(),
				"was array with elements of type %v", elementType).
				WithCauses(result)
		}
		return hamcrest.NewResultf(false, "was of type %T, not an array", actual)
	}
	return hamcrest.NewMatcherf(match, "ArrayOf(%v)", elementTypeMatcher)
}

// Returns a new matcher that, on any input that is a *reflect.ChanType,
// extracts its element type and matches it against the given matcher.
//
// If the given input is not a *reflect.ChanType, this fails to match.
// Note: this matches channel *types*, not *channels*. (See ChannelOf.)
func ChannelTypeOf(elementTypeMatcher *hamcrest.Matcher) *hamcrest.Matcher {
	match := func(actual interface{}) *hamcrest.Result {
		if channelType, ok := actual.(*reflect.ChanType); ok {
			elementType := channelType.Elem()
			result := elementTypeMatcher.Match(elementType)
			return hamcrest.NewResultf(
				result.Matched(),
				"was *reflect.ChanType with elements of type %v", elementType).
				WithCauses(result)
		}
		return hamcrest.NewResultf(false,
			"was of type %T, not a *reflect.ChanType", actual)
	}
	return hamcrest.NewMatcherf(match, "ChannelTypeOf(%v)", elementTypeMatcher)
}

// Returns a new matcher that, on any input that is a channel, extracts
// its type and matches it against the given matcher.
//
// If the given input is not a channel, this fails to match.
// Note: this matches *channels*, not channel *types*. (See ChannelTypeOf.)
func ChannelOf(elementTypeMatcher *hamcrest.Matcher) *hamcrest.Matcher {
	match := func(actual interface{}) *hamcrest.Result {
		actualType := reflect.Typeof(actual)
		if channelType, ok := actualType.(*reflect.ChanType); ok {
			elementType := channelType.Elem()
			result := elementTypeMatcher.Match(elementType)
			return hamcrest.NewResultf(result.Matched(),
				"was channel with elements of type %v",
				elementType).
				WithCauses(result)
		}
		return hamcrest.NewResultf(false, "was of type %T, not a channel", actual)
	}
	return hamcrest.NewMatcherf(match, "ChannelOf(%v)", elementTypeMatcher)
}

// Returns a new matcher that, on any input that is a *reflect.SliceType,
// extracts the type of element and matches it against the given matcher.
//
// If the given input is not an *reflect.SliceType, this fails to match.
// Note:  this matches slice *types*, not slices.  (See SliceOf.)
func SliceTypeOf(elementTypeMatcher *hamcrest.Matcher) *hamcrest.Matcher {
	match := func(actual interface{}) *hamcrest.Result {
		if sliceType, ok := actual.(*reflect.SliceType); ok {
			elementType := sliceType.Elem()
			result := elementTypeMatcher.Match(elementType)
			return hamcrest.NewResultf(
				result.Matched(),
				"was SliceType with elements of type %v", elementType.Name()).
				WithCauses(result)
		}
		return hamcrest.NewResultf(false, "was of type %T, not a slice", actual)
	}
	return hamcrest.NewMatcherf(match, "SliceTypeOf(%v)", elementTypeMatcher)
}

// Returns a new matcher that, on any input that is an array, extracts
// its type and matches it against the given matcher.
//
// If the given input is not an array, this fails to match.
// Note: this matches *slices*, not slice *types*.  (See SliceTypeOf.)
func SliceOf(elementTypeMatcher *hamcrest.Matcher) *hamcrest.Matcher {
	match := func(actual interface{}) *hamcrest.Result {
		actualType := reflect.Typeof(actual)
		if sliceType, ok := actualType.(*reflect.SliceType); ok {
			elementType := sliceType.Elem()
			result := elementTypeMatcher.Match(elementType)
			return hamcrest.NewResultf(
				result.Matched(),
				"was slice with elements of type %v", elementType).
				WithCauses(result)
		}
		return hamcrest.NewResultf(false, "was of type %T, not a slice", actual)
	}
	return hamcrest.NewMatcherf(match, "SliceOf(%v)", elementTypeMatcher)
}

// Returns a new matcher that, on any input that is a *reflect.MapType,
// extracts the type of keys and element and matches them against two
// given matchers.
//
// If the given input is not an *reflect.MapType, this fails to match.
// Note:  this matches map *types*, not maps.  (See MapOf.)
func MapTypeOf(keyTypeMatcher, elementTypeMatcher *hamcrest.Matcher) *hamcrest.Matcher {
	match := func(actual interface{}) *hamcrest.Result {
		if mapType, ok := actual.(*reflect.MapType); ok {
			keyType := mapType.Key()
			elementType := mapType.Elem()
			keyResult := keyTypeMatcher.Match(keyType)
			if !keyResult.Matched() {
				return hamcrest.NewResultf(false,
					"was MapType with keys of type %v", keyType).
					WithCauses(keyResult)
			}
			elementResult := elementTypeMatcher.Match(elementType)
			return hamcrest.NewResultf(elementResult.Matched(),
				"was MapType with keys/elements of type %v/%v",
				keyType, elementType).
				WithCauses(keyResult, elementResult)
		}
		return hamcrest.NewResultf(false, "was of type %T, not a MapType", actual)
	}
	return hamcrest.NewMatcherf(match,
		"MapTypeOf(%v, %v)", keyTypeMatcher, elementTypeMatcher)
}

// Returns a new matcher that, on any input that is a map, extracts the
// type of keys and elements and matches them against the given matchers.
//
// If the given input is not an map, this fails to match.
// Note:  this matches maps, not map *types*.  (See MapTypeOf.)
//
// This matcher is logically equivalent to:
//    AllOf(MapWithKeyType(keyTypeMatcher), MapWithElementType(elementTypeMatcher))
// but may be easier to read/type.
func MapOf(keyTypeMatcher, elementTypeMatcher *hamcrest.Matcher) *hamcrest.Matcher {
	return ToType(MapTypeOf(keyTypeMatcher, elementTypeMatcher))
}



// Returns a new matcher that, on any input that is a *reflect.PtrType,
// extracts the type of object that it thinks it's pointing to (the
// "pointee") and matches it against the given matcher.
//
// If the given input is not an *reflect.PtrType, this fails to match.
// Note:  this matches pointer *types*, not pointers. (See PointerOf.)
func PtrTypeTo(pointeeTypeMatcher *hamcrest.Matcher) *hamcrest.Matcher {
	match := func(actual interface{}) *hamcrest.Result {
		if ptrType, ok := actual.(*reflect.PtrType); ok {
			elementType := ptrType.Elem()
			result := pointeeTypeMatcher.Match(elementType)
			return hamcrest.NewResultf(
				result.Matched(),
				"was PtrType pointing to type %v", elementType.Name()).
				WithCauses(result)
		}
		return hamcrest.NewResultf(false,
			"was type %T, not a PtrType", actual)
	}
	return hamcrest.NewMatcherf(match,
		"PtrTypeTo(%v)", pointeeTypeMatcher)
}

// Returns a new matcher that, on any input that is a pointer, extracts the
// type of object that it thinks it's pointing to (the "pointee") and
// matches it against the given matcher.
//
// If the given input is not an pointer, this fails to match.
// Note:  this matches *pointers*, not pointer *types*. (See PtrTypeTo.)
func PtrTo(pointeeTypeMatcher *hamcrest.Matcher) *hamcrest.Matcher {
	match := func(actual interface{}) *hamcrest.Result {
		actualType := reflect.Typeof(actual)
		if ptrType, ok := actualType.(*reflect.PtrType); ok {
			elementType := ptrType.Elem()
			result := pointeeTypeMatcher.Match(elementType)
			return hamcrest.NewResultf(
				result.Matched(), "was PtrType to type %v", elementType).
				WithCauses(result)
		}
		return hamcrest.NewResultf(false,
			"was type %T, not a pointer", actual)
	}
	return hamcrest.NewMatcherf(match,
		"PtrTo(%v)", pointeeTypeMatcher)
}




