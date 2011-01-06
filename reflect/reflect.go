package reflect

import (
	"hamcrest"
	"reflect"
)

func ToType(matcher *hamcrest.Matcher) *hamcrest.Matcher {
	description := hamcrest.NewDescription("ToType(%v)", matcher)
	match := func(actual interface{}) *hamcrest.Result {
		actualType := reflect.Typeof(actual)
		result := matcher.Match(actualType)
		why := hamcrest.NewDescription(
			"reflect.Typeof() returned %v", actualType)
		return hamcrest.NewResult(result.Matched(), why).WithCauses(result)
	}
	return hamcrest.NewMatcher(description, match)
}

func _TypeMatcher(expectedType reflect.Type) *hamcrest.Matcher {
	whyMatched := hamcrest.NewDescription("was a %v", expectedType)
	match := func(actual interface{}) *hamcrest.Result {
		if actual == nil {
			return hamcrest.NewResult(false, hamcrest.NewDescription("was nil"))
		}
		actualType := reflect.Typeof(actual)
		if reflect.DeepEqual(actualType, expectedType) {
			return hamcrest.NewResult(true, whyMatched)
		}
		whyNotMatched := hamcrest.NewDescription(
			"was a %v, not a %v", actualType, expectedType)
		return hamcrest.NewResult(false, whyNotMatched)
	}
	description := hamcrest.NewDescription(expectedType.Name())
	return hamcrest.NewMatcher(description, match)
}

var (
	boolType = reflect.Typeof(false)
	intType = reflect.Typeof(int(0))
	int8Type = reflect.Typeof(int8(0))
	int16Type = reflect.Typeof(int16(0))
	int32Type = reflect.Typeof(int32(0))
	int64Type = reflect.Typeof(int64(0))
	uintType = reflect.Typeof(uint(0))
	uint8Type = reflect.Typeof(uint8(0))
	uint16Type = reflect.Typeof(uint16(0))
	uint32Type = reflect.Typeof(uint32(0))
	uint64Type = reflect.Typeof(uint64(0))
	uintptrType = reflect.Typeof(uintptr(0))
	floatType = reflect.Typeof(float(0))
	float32Type = reflect.Typeof(float32(0))
	float64Type = reflect.Typeof(float64(0))
	complexType = reflect.Typeof(complex(0i))
	complex64Type = reflect.Typeof(complex64(0i))
	complex128Type = reflect.Typeof(complex128(0i))
	stringType = reflect.Typeof("")
	
	boolMatcher = _TypeMatcher(boolType)
	boolTypeMatcher = hamcrest.DeeplyEqualTo(boolType)
	intMatcher = _TypeMatcher(intType)
	intTypeMatcher = hamcrest.DeeplyEqualTo(intType)
	uintMatcher = _TypeMatcher(uintType)
	uintTypeMatcher = hamcrest.DeeplyEqualTo(uintType)
	floatMatcher = _TypeMatcher(floatType)
	floatTypeMatcher = hamcrest.DeeplyEqualTo(floatType)
	complexMatcher = _TypeMatcher(complexType)
	complexTypeMatcher = hamcrest.DeeplyEqualTo(complexType)
	stringMatcher = _TypeMatcher(stringType)
	stringTypeMatcher = hamcrest.DeeplyEqualTo(stringType)
)

func Bool() *hamcrest.Matcher { return boolMatcher }
func Int() *hamcrest.Matcher { return intMatcher }
func Uint() *hamcrest.Matcher { return uintMatcher }
func Float() *hamcrest.Matcher { return floatMatcher }
func Complex() *hamcrest.Matcher { return complexMatcher }
func String() *hamcrest.Matcher { return stringMatcher }

func BoolType() *hamcrest.Matcher { return boolTypeMatcher }
func IntType() *hamcrest.Matcher { return intTypeMatcher }
func UintType() *hamcrest.Matcher { return uintTypeMatcher }
func FloatType() *hamcrest.Matcher { return floatTypeMatcher }
func ComplexType() *hamcrest.Matcher { return complexTypeMatcher }
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
			description := hamcrest.NewDescription(
				"was ArrayType with elements of type %v", elementType.Name())
			result := elementTypeMatcher.Match(elementType)
			return hamcrest.NewResult(
				result.Matched(), description).WithCauses(result)
		}
		why := hamcrest.NewDescription("was of type %T, not an ArrayType", actual)
		return hamcrest.NewResult(false, why)
	}
	description := hamcrest.NewDescription("ArrayTypeOf(%v)", elementTypeMatcher)
	return hamcrest.NewMatcher(description, match)
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
			description := hamcrest.NewDescription(
				"was array with elements of type %v", elementType)
			result := elementTypeMatcher.Match(elementType)
			return hamcrest.NewResult(
				result.Matched(), description).WithCauses(result)
		}
		why := hamcrest.NewDescription("was of type %T, not an array", actual)
		return hamcrest.NewResult(false, why)
	}
	description := hamcrest.NewDescription("ArrayOf(%v)", elementTypeMatcher)
	return hamcrest.NewMatcher(description, match)
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
			description := hamcrest.NewDescription(
				"was SliceType with elements of type %v", elementType.Name())
			result := elementTypeMatcher.Match(elementType)
			return hamcrest.NewResult(
				result.Matched(), description).WithCauses(result)
		}
		why := hamcrest.NewDescription("was of type %T, not a slice", actual)
		return hamcrest.NewResult(false, why)
	}
	description := hamcrest.NewDescription("SliceTypeOf(%v)", elementTypeMatcher)
	return hamcrest.NewMatcher(description, match)
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
			description := hamcrest.NewDescription(
				"was slice with elements of type %v", elementType)
			result := elementTypeMatcher.Match(elementType)
			return hamcrest.NewResult(
				result.Matched(), description).WithCauses(result)
		}
		why := hamcrest.NewDescription("was of type %T, not a slice", actual)
		return hamcrest.NewResult(false, why)
	}
	description := hamcrest.NewDescription("SliceOf(%v)", elementTypeMatcher)
	return hamcrest.NewMatcher(description, match)
}

// Returns a new matcher that, on any input that is a *reflect.MapType,
// extracts the type of key element and matches it against the given matcher.
//
// If the given input is not an *reflect.MapType, this fails to match.
// Note:  this matches map *types*, not maps.  (See MapWithKeyType.)
func MapTypeWithKeyType(keyTypeMatcher *hamcrest.Matcher) *hamcrest.Matcher {
	match := func(actual interface{}) *hamcrest.Result {
		if mapType, ok := actual.(*reflect.MapType); ok {
			keyType := mapType.Key()
			description := hamcrest.NewDescription(
				"was MapType with keys of type %v", keyType)
			result := keyTypeMatcher.Match(keyType)
			return hamcrest.NewResult(
				result.Matched(), description).WithCauses(result)
		}
		why := hamcrest.NewDescription("was of type %T, not a MapType", actual)
		return hamcrest.NewResult(false, why)
	}
	description := hamcrest.NewDescription("MapTypeWithKeyType(%v)", keyTypeMatcher)
	return hamcrest.NewMatcher(description, match)
}

// Returns a new matcher that, on any input that is a map, extracts the
// type of key element and matches it against the given matcher.
//
// If the given input is not an map, this fails to match.
// Note:  this matches maps, not map *types*.  (See MapTypeWithKeyType.)
func MapWithKeyType(keyTypeMatcher *hamcrest.Matcher) *hamcrest.Matcher {
	match := func(actual interface{}) *hamcrest.Result {
		actualType := reflect.Typeof(actual)
		if mapType, ok := actualType.(*reflect.MapType); ok {
			keyType := mapType.Key()
			description := hamcrest.NewDescription(
				"was map with keys of type %v", keyType)
			result := keyTypeMatcher.Match(keyType)
			return hamcrest.NewResult(
				result.Matched(), description).WithCauses(result)
		}
		why := hamcrest.NewDescription("was of type %T, not a map", actual)
		return hamcrest.NewResult(false, why)
	}
	description := hamcrest.NewDescription("MapWithKeyType(%v)", keyTypeMatcher)
	return hamcrest.NewMatcher(description, match)
}

// Returns a new matcher that, on any input that is a *reflect.MapType,
// extracts the type of element and matches it against the given matcher.
//
// If the given input is not an *reflect.MapType, this fails to match.
// Note:  this matches map *types*, not maps.  (See MapWithElementType.)
func MapTypeWithElementType(elementTypeMatcher *hamcrest.Matcher) *hamcrest.Matcher {
	match := func(actual interface{}) *hamcrest.Result {
		if mapType, ok := actual.(*reflect.MapType); ok {
			elementType := mapType.Elem()
			description := hamcrest.NewDescription(
				"was MapType with elements of type %v", elementType)
			result := elementTypeMatcher.Match(elementType)
			return hamcrest.NewResult(
				result.Matched(), description).WithCauses(result)
		}
		why := hamcrest.NewDescription("was of type %T, not a MapType", actual)
		return hamcrest.NewResult(false, why)
	}
	description := hamcrest.NewDescription(
		"MapTypeWithElementType(%v)", elementTypeMatcher)
	return hamcrest.NewMatcher(description, match)
}

// Returns a new matcher that, on any input that is a *reflect.MapType,
// extracts the type of key element and matches it against the given matcher.
//
// If the given input is not an *reflect.MapType, this fails to match.
// Note:  this matches map *types*, not maps.  (See MapTypeWithElementType.)
func MapWithElementType(elementTypeMatcher *hamcrest.Matcher) *hamcrest.Matcher {
	match := func(actual interface{}) *hamcrest.Result {
		actualType := reflect.Typeof(actual)
		if mapType, ok := actualType.(*reflect.MapType); ok {
			elementType := mapType.Elem()
			description := hamcrest.NewDescription(
				"was map with elements of type %v", elementType)
			result := elementTypeMatcher.Match(elementType)
			return hamcrest.NewResult(
				result.Matched(), description).WithCauses(result)
		}
		why := hamcrest.NewDescription("was of type %T, not a map", actual)
		return hamcrest.NewResult(false, why)
	}
	description := hamcrest.NewDescription(
		"MapWithElementType(%v)", elementTypeMatcher)
	return hamcrest.NewMatcher(description, match)
}





