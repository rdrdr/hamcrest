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
    `Is`, `Anything`, `True`, `False`, `Nil`, `EqualTo`, and `DeeplyEqualTo`
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


How to use hamcrest for testing:
================================

To use Hamcrest matchers, one creates an `Asserter`
and uses it to `Check` or `Assert` values.

	func TestPoint(t *testing.T) {
		p := Point(3, 4)
		we := asserter.Using(t)
		we.AssertThat(p.X, EqualTo(3).Comment("x coordinate"))
		we.AssertThat(p.Y, EqualTo(4).Comment("y coordinate"))
		we.CheckThat(p, ToString(EqualTo("[3, 4]")))
		we.CheckThat(p, DeeplyEqualTo(Point(3, 4)))
	}

(`Assert` methods fail immediately, as `testing.T.FailNow`, while `Check`
methods defer failure, as `testing.T.Fail`.)

The `AssertThat` and `CheckThat` functions are designed to create
conditional checks that read fluently as self-commenting code, and
are self-describing when failures occur.  For example, the above
test might fail with this message:
	
	FAILURE: input [4, 3] on EqualTo([12, 34])
		Because: ([34, 12]) was not equal to ([12, 34])

Effort invested in good self-describing matchers can be leveraged
across many tests.


A tour of common matchers
=========================

Hamcrest comes with a library of useful matchers. Here are some of the most
common ones.

  * `Anything` - always matches
  * `EqualTo(obj)` - matches any object `x` where `x==obj` would be true
  * `DeeplyEqualTo(obj)` - matches any object `x` where `reflect.DeepEquals(x, obj)` is true
  * `True` - only matches bool `true`
  * `False` - only matches bool `false`
  * `Not(matcher)` - logical not of `matcher`
  * `Nil` - matches objects whose types have an `IsNil()` method  which returns true for the object
  * `NonNil` - inverse of `Nil` matcher (equivalent to `Not(Nil)`)
  * `AnyOf(matchers...)` - short-circuiting n-ary logical Or
  * `AllOf(matchers...)` - short-circuiting n-ary logical And

Although it is possible to simulate these logical conditions using the above,
separate version are provided to assist readability.

  * `Is(matcher)` - equivalent to `matcher` (see =Syntactic sugar=, below)
  * `Both(matcher1).And(matcher2)` - short-circuiting logical And 
  * `Either(matcher1).Or(matcher2)` - short-circuiting logical Or
  * `Neither(matcher1).Nor(matcher2)` - short-circuiting logical Nor
  * `If(matcher1).Then(matcher2)` - short-circuiting logical If/Then
  * `Iff(matcher1).Then(matcher2)` - logical If-And-Only-If (note: iff never short-circuits)
  * `Either(matcher1).Xor(matcher)` - logical xor (note: xor never short-circuits)

Syntactic sugar
===============

Hamcrest strives to make your tests as readable as possible. For example,
the `Is` matcher is a wrapper that doesn't add any extra behavior to the
underlying matcher. The following assertions are equivalent:

    we.AssertThat(theBiscuit, EqualTo(myBiscuit));
    we.AssertThat(theBiscuit, Is(EqualTo(myBiscuit)));

Custom matchers
===============

Example:

    func IsMultipleOf(k int) *hamcrest.Matcher {
        match := func(actual interface{}) {
            if n, ok := actual.(int); ok {
                if n % k == 0 {
                    because := hamcrest.NewDescription(
                        "%v is divisible by %v", n, k)
                    return hamcrest.NewResult(true, because)
                }
                because := hamcrest.NewDescription(
                    "%v is not divisible by %v", n, k)
                return hamcrest.NewResult(false, because)
            }
            because := hamcrest.NewDescription(
                    "can't convert %T[%v] to int", actual, actual)
            return hamcrest.NewResult(false, because)
        }
        return NewMatcher(hamcrest.NewDescription("multiple of %v", n), match)
    }

And used:
    we.CheckThat(recordSize, IsMultipleOf(16).AddComment(
        "profiling suggests better performance than 8, but 32 is unnecessary"))

        
        