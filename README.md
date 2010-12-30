Introduction
============

Hamcrest is a fluent framework for declarative Matcher objects
that, when applied to input values, produce self-describing
results.

Packages
========

`hamcrest.go` comes in several packages that you assemble to fit your needs:

*  `hamcrest`:  Defines the types `Matcher` and `Result`, provides factory
    functions to create them, and defines several core Matchers:
    `Is`, `Anything`, `True`, `False`, `Nil`, `EqualTo`, and `DeepEqualTo`
    plus common logical compositions: `Not`, `And`, `Or`, `Nor`, `Xor`,
    and `If/Then`.

*   `hamcrest/strings`:  Matchers for strings.

*   `hamcrest/asserter`:  Defines an `Asserter` that can be used in conjunction 
    with Hamcrest Matchers to produce helpful logging messages at runtime
    (to stdout, stderr, or any object that implements io.Writer) or in
    unit tests (using `testing.T` from Go's standard `testing` package).

    Note: this package isn't *really* part of Hamcrest:  it's just a handy
    way of using the Hamcrest results in conjunction with the standard Go
    testing package.

Expected future packages:

*   `hamcrest/files`:  Matchers for files.
*   `hamcrest/maps`:  Matchers for maps.
*   `hamcrest/numbers`:  Matchers for numerical values.
*   `hamcrest/slices`:  Matchers for slices.
*   `hamcrest/types`:  Matchers using type reflection.


How to use hamcrest in tests:
=============================

Given a unit test that looks like this:

	package foo
	
	import "testing"
	
	func TestCoord(t *testing.T) {
		coord := Coord(12, 34)
		if coord.X != 12 {
			t.Fatalf("Expected X to be 12, was %v", coord.X)
		}
		if coord.Y != 34 {
			t.Fatalf("Expected Y to be 34, was %v", coord.Y)
		}
		if coord.String() != "[12, 34]" {
			t.Errorf("Expected String() to be '[12, 34]', was %v", coord.String())
		}
		if other := Coord(12, 34); !reflect.DeepEquals(coord, other) {
			t.Errorf("Similarly constructed coord %v should be deeply equal to %v",
				other, coord)
		}
	}
	
To convert this example to Hamcrest matchers, one creates an `Asserter`
and uses it to Check or Assert .

	package foo
	
	import (
		"testing"
		"github.com/rdrdr/hamcrest.go/hamcrest"
		"github.com/rdrdr/hamcrest.go/hamcrest/asserter"
                "github.com/rdrdr/hamcrest.go/hamcrest/strings"
	)
	
	func TestCoord(t *testing.T) {
		Equals := hamcrest.DeeplyEqualTo
		ToString := strings.ToString
		
		we := asserter.Using(t)
		coord := Coord(12, 34)
		we.AssertThat(coord.X(), Equals(12).AddComment("X field"))
		we.AssertThat(coord.Y(), Equals(34).AddComment("Y field"))
		we.CheckThat(coord, ToString(Equals("[12, 34]"))))
		we.CheckThat(coord, Equals(Coord(12, 34)).Comment("Similarly constructed coord"))
	}

The `AssertThat` and `CheckThat` methods lead to stylized sentences for making
test assertions.

A tour of common matchers
=========================

Hamcrest comes with a library of useful matchers. Here are some of the most
important ones.

  * `Anything` - always matches
  * `EqualTo(obj)` - matches any object `x` where `x==obj` would be true
  * `DeeplyEqualTo(obj)` - matches any object `x` where `reflect.DeepEquals(x, obj)` is true
  * `True` - only matches bool `true`
  * `False` - only matches bool `false`
  * `Not(matcher)` - logical not of `matcher`
  * `Nil` - matches objects whose types have an `IsNil()` method  which returns true for the object
  * `NonNil` - inverse of `Nil` matcher (equivalent to `Not(Nil)`)
  * `Both(matcher1).And(matcher2)` - short-circuiting logical and 
  * `Either(matcher1).Or(matcher2)` - short-circuiting logical or
  * `Neither(matcher1).Nor(matcher2)` - short-circuiting logical nor
  * `If(matcher1).Then(matcher2)` - short-circuiting logical if/then
  * `Either(matcher1).Xor(matcher)` - logical xor (note: xor is never short-circuiting)

Syntactic sugar
===============

Hamcrest strives to make your tests as readable as possible. For example,
the `Is` matcher is a wrapper that doesn't add any extra behavior to the
underlying matcher. The following assertions are equivalent:

    we.AssertThat(theBiscuit, EqualTo(myBiscuit));
    we.AssertThat(theBiscuit, Is(EqualTo(myBiscuit)));

TODO:  Make this statement true:
The last form is allowed since ``is`` wraps non-matcher arguments with
``equalTo``. Other matchers that take matchers as arguments provide similar
shortcuts, wrapping non-matcher arguments in ``equalTo``.


Custom matchers
===============

Example:

    func IsMultipleOf(k int) *hamcrest.Matcher {
        match := func(actual interface{}) {
            if n, ok := actual.(int); ok {
                if n % k == 0 {
                    because := hamcrest.NewDescription("%v is divisible by %v", n, k)
                    return hamcrest.NewResult(true, because)
                }
                because := hamcrest.NewDescription("%v is not divisible by %v", n, k)
                return hamcrest.NewResult(false, because)
            }
            because := hamcrest.NewDescription("can't convert %v to int", actual)
            return hamcrest.NewResult(false, because)
        }
        return NewMatcher(hamcrest.NewDescription("multiple of %v", n), match)        
    }

And used:

    we.CheckThat(13, IsMultipleOf(12).AddComment("A Bakers' Dozen"))
