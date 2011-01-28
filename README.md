Introduction
============

Hamcrest is a fluent framework for declarative Matcher objects
that, when applied to input values, produce self-describing
results.

Installation:
========

To install, run `make install` from the same directory as this
`README.md` file.

Packages
========

`hamcrest.go` comes in several packages that you assemble to fit your needs:

*   `hamcrest/base`:  Defines the types `Matcher`, `Result` and `SelfDescribing`
    and provides factory functions to create them.  (Unless you want to define
    your own custom Matchers, you'll only use this indirectly.)

*   `hamcrest/core`:  Defines a set of Matchers for doing basic comparisons,
    equality testing, nil checking, and grouping/composition matchers.

*   `hamcrest/collections`:  Matchers on arrays/slices/maps, such as
    `EachElement`, `EveryElement`, `EachMapElement`, `EveryMapElement`.

*   `hamcrest/reflect`:  Matchers using type reflection, such as `ToType`,
    `SameTypeAs`, `SliceOf`, `MapOf`, etc.

*   `hamcrest/strings`:  Matchers for strings.

*   `hamcrest/asserter`:  Defines an `Asserter` that can be used in conjunction 
    with Hamcrest Matchers to produce helpful logging messages at runtime
    (to stdout, stderr, or any object that implements io.Writer) or in
    unit tests (using `testing.T` from Go's standard `testing` package).

    Note: the `asserter` package isn't *really* part of Hamcrest:  it's just
    a handy way of using the Hamcrest results in conjunction with the
    standard Go testing package.

You may also choose to write your own Matchers (see *Custom matchers*, below).


Example of using Hamcrest at runtime:
=====================================

Create an `Asserter`.  The simplest way to do this is with the factory
method `UsingStderr()`, which returns an `Asserter` that logs problems to
stderr and calls `panic` on `FailNow`:

	import (
		"github.com/rdrdr/hamcrest/asserter"
	)
	
	var we = asserter.UsingStderr()

Use that asserter during init() to make sure globals are properly
initialized:

	import (
		"github.com/rdrdr/hamcrest/asserter"
		"github.com/rdrdr/hamcrest/collections"
		"github.com/rdrdr/hamcrest/core"
		"github.com/rdrdr/hamcrest/strings"
	)
	var we = asserter.UsingStderr()
	
	type Server struct {
		hostname string
		port uint16
	}
	var servers = []Server {
		{ "news.foo.com", 8000 },
		{ "news.bar.com", 8888 },
	}

	func init() {
		EveryElement := collections.EveryElement
		ToHostname := core.Applying(func(s Server) string {
			return s.hostname
		}, "ToHostname")
		IsInOneOfOurDomains := core.AnyOf(strings.HasSuffix(".foo.com"),
		                                  strings.HasSuffix(".bar.com"))
		
		we.FailNowUnless(servers,
			EveryElement(ToHostname(IsInOneOfOurDomains)))
	}

Or use the asserter at runtime to guarantee that a method's
preconditions are met:

	func WriteTo(filename string) bool {
		we.AssertThat(filename, strings.EndsWith(".txt").
			Comment("Files must have txt extension."))
		// Use filename here.
	}

Or use it during development to write your tests in the same file as your code:

	func PigLatin(input string) string {
		...implementation...
	}
	
	func init() {
		we := asserter.UsingStderr()
		we.AssertThat(PigLatin("testing"), EqualTo("esting-tay"))
		we.AssertThat(PigLatin("made"), EqualTo("ade-may"))
		we.AssertThat(PigLatin("easier"), EqualTo("easier-ay"))
	}
	
This makes it easy to cut-and-paste each `init()` block into your
testing suite.  While moving the block over, replace:

	func init() {
		we := asserter.UsingStderr()
		...

With an `Asserter` that uses the testing infrastructure:

	func Test_PigLatin(t *testing.T) {
		we := asserter.Using(t)
		...

How to use hamcrest for testing:
================================

To use Hamcrest matchers, create an `Asserter` and use it to
`Check` or `Assert` that values meet the criteria of those
matchers:

	func TestPoint(t *testing.T) {
		p := Point(3, 4)
		we := asserter.Using(t)
		we.AssertThat(p.X, EqualTo(3).Comment("x coord"))
		we.AssertThat(p.Y, EqualTo(4).Comment("y coord"))
		we.CheckThat(p, ToString(EqualTo("[3, 4]")))
		we.CheckThat(p.Magnitude(), EqualTo(5).Comment("magnitude"))
	}

(`Assert` methods fail immediately, as `testing.T.FailNow()` does,
while `Check` methods defer failure, as `testing.T.Fail()` does.)

The `AssertThat` and `CheckThat` functions are designed to create
conditional checks that read fluently as self-commenting code, and
are self-describing when failures occur.  For example, the above
test might fail with this message:
	
	FAILURE on input &Point{X:3, Y:4}
		Did not match ToString(EqualTo([3, 4]))
		Because: String() was "[4, 3]"
			Did not match EqualTo[[3, 4]]
			Because: "[4, 3]" was not equal to "[3, 4]"

Or:
	FAILURE on input 5
		Did not match EqualTo(5)
		Because: uint 5 could not be compared to int 5
		Comment: magnitude

Note that the majority of the text descriptions are generated
automatically by the matchers.  For typical uses of Hamcrest
matchers, the code is largely self-documenting, and the error
messages are detailed.

Effort invested in good self-describing matchers can be leveraged
across many tests.


A note on library design:
=========================

Hamcrest is designed to be a fluent library.  The Go syntax requires external
symbols be preceded by their package name, which can lead to occasionally
awkward constructions:
	we.CheckThat("foobar", core.AllOf(strings.HasPrefix("foo"), strings.HasSuffix("bar")))

To avoid this clunkiness, Hamcrest matchers are generated by functions that
you can assign to local names:
	AllOf := core.AllOf
	StartsWith := strings.HasPrefix
	EndsWith := strings.HasSuffix
	we.CheckThat("foobar", AllOf(StartsWith("foo"), EndsWith("bar")))


Performance note:
=================
Note:  Hamcrest matchers allocate Description and Result objects to
explain in great detail why they did or did not match. However, these
objects are lazily evaluated.  Users should be generally aware that
there is an *object allocation* cost to using Hamcrest matchers, but
there is (generally) no *string construction* cost unless a Hamcrest
Matcher or Result is explicitly asked to self-describe.

Still, users who are particularly sensitive to performance concerns
may wish to think carefully before using Hamcrest matchers in
performance-critical bottlenecks.


A tour of common matchers
=========================

Hamcrest comes with a library of useful matchers. Here are some of the most
common ones.

  * `Anything` - matches any input
  
  * `True` - only matches bool `true`
  * `False` - only matches bool `false`
  * `Not(matcher)` - logical not of `matcher`.
     
  
  * `If(m1).Then(m2)` - checks that whenever `m1` matches, so does `m2`
  * `IfAndOnlyIf(m1).Then(m2)` - checks that `m1` and `m2` both match/don't match.

  * `EqualTo(y)` - matches any value `x` where `x == y` is legal and true
  * `NotEqualTo(y)` - matches any value `x` where `x != y` is legal and true
  * `DeepEqualTo(y)` - matches any value `x` where `reflect.DeepEquals(x, y)` is true
  * `GreaterThan(y)` - matches any value `x` where `x > y` is legal and true
  * `GreaterThanOrEqualTo(y)` - matches any value `x` where `x <= y` is legal and true
  * `LessThan(y)` - matches any value `x` where `x < y` is legal and true
  * `LessThanOrEqualTo(y)` - matches any value `x` where `x <= y` is legal and true
  
  * `Nil` - matches values of any type with an `IsNil() bool` method that returns true for the given object
  * `NonNil` - matches values of any type with an `IsNil() bool` method that returns false for the given object
  
  * `AnyOf(matchers...)` - short-circuiting n-ary logical Or
  * `AllOf(matchers...)` - short-circuiting n-ary logical And

Syntactic sugar
===============

Hamcrest strives to make your tests as readable as possible. For example,
the `Is` matcher is a wrapper that doesn't add any extra behavior to the
underlying matcher. The following assertions are equivalent:

    we.AssertThat(x, EqualTo(y));
    we.AssertThat(x, Is(EqualTo(y)))
    
Similarly, it is possible to simulate common logical conditions using the `logic`
package for readability:

  * `Is(matcher)` - equivalent to `matcher`
  * `Both(matcher1).And(matcher2)` - short-circuiting logical `And`, equivalent to `AllOf(matcher1, matcher2)`
  * `Either(matcher1).Or(matcher2)` - short-circuiting logical `Or`, equivalent to `AnyOf(matcher1, matcher2)`
  * `Neither(matcher1).Nor(matcher2)` - short-circuiting logical `Nor`
  * `If(matcher1).Then(matcher2)` - short-circuiting logical `If/Then`
  * `Iff(matcher1).Then(matcher2)` - logical `IfAndOnlyIf` (note: iff never short-circuits)
  * `Either(matcher1).Xor(matcher)` - logical `Xor` (note: xor never short-circuits)

    
Custom matchers
===============

Example:

	func IsMultipleOf(k int) *base.Matcher {
		match := func(n int) bool {
			return n % k == 0
		}
		return base.NewMatcherf(match, "multiple of %v", k)
	}
	
And used:

	we.CheckThat(recordSize, IsMultipleOf(8).Comment(
		"profiling suggests better performance than multiple of 4"))

Or, if you want more control over the error messages:

	func IsMultipleOf(k int) *base.Matcher {
		match := func(actual interface{}) *base.Result {
			if n, ok := actual.(int); ok {
				if remainder := n % k; remainder != 0 {
					return base.NewResultf(false,
						"was not a multiple of %v (remainder %v)", k, remainder)
				}
				return base.NewResultf(true, "was a multiple of %v", k)
			}
			return base.NewResultf(false, "was a %T, expected an int!", actual)
		}
		return base.NewMatcherf(match, "multiple of %v", k)
	}



