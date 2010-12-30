package hamcrest

import (
	"fmt"
)

// --------------------------------------------------------------------
// Description
// --------------------------------------------------------------------

// Hamcrest descriptions implement both fmt.Stringer and fmt.Formatter.
type Description struct {
	format string
	args []interface{}
}

// Creates an object that implements fmt.Formatter and fmt.Stringer using
// args with the same meanings as fmt.Fprintf.  Note that this object
// stores its given parameters and evaluates them lazily.
func NewDescription(format string, args...interface{}) *Description {
	return &Description{format:format, args:args}
}

// Implements fmt.Formatter.
func (self *Description) Format(s fmt.State, ch int) {
	fmt.Fprintf(s, self.format, self.args...)
}

// Implements fmt.Stringer.
func (self *Description) String() string {
	return fmt.Sprintf(self.format, self.args...)
}


// --------------------------------------------------------------------
// Result
// --------------------------------------------------------------------

// Self-describing result of applying a Matcher to an input value.
type Result struct {
	matched bool
	matcher *Matcher
	description *Description
	causes []*Result
}
var _NO_CAUSES = make([]*Result, 0)

// Creates a new Result using the given description
func NewResult(matched bool, description *Description) *Result {
	return &Result{
		matched:matched, description:description, causes:_NO_CAUSES}
}

// Returns true if the Match was successful.
func (self *Result) Matched() bool {
	return self.matched
}

// Returns the Matcher that produced this Result.
func (self *Result) Matcher() *Matcher {
	return self.matcher
}

// Implements fmt.Stringer.
func (self *Result) String() string {
	return self.description.String()
}

// Implements fmt.Formatter.
func (self *Result) Format(s fmt.State, ch int) {
	self.description.Format(s, ch)
}

// Returns a slice of sub-Results that caused this Result to
// either match or not match.
func (self *Result) Causes() []*Result {
	causes := make([]*Result, len(self.causes))
	copy(causes, self.causes)
	return causes
}

// Returns a new Result, identical to this one, except with
// the given causes.
func (self *Result) WithCauses(causes... *Result) *Result {
	return &Result{matched:self.matched, description:self.description, causes:causes}
}

// --------------------------------------------------------------------
// Matcher
// --------------------------------------------------------------------

// Self-describing criteria that may match (or not match) an input value.
// Creators of new matchers are *strongly* encouraged not to implement
// Matcher directly, but to create new matchers using the NewMatcher
// factory function.
type Matcher struct {
	description *Description
	match func(v interface{}) *Result
	comments []interface{}
}
var _NO_COMMENTS = make([]interface{}, 0)

// Creates a new Matcher using the given description
func NewMatcher(description *Description, match func(v interface{}) *Result) *Matcher {
	return &Matcher{description:description, match:match, comments:_NO_COMMENTS}
}

// Implementation of SelfDescribing: fmt.Formatter.
func (self *Matcher) Format(s fmt.State, ch int) {
	self.description.Format(s, ch)
}

// Implementation of SelfDescribing: fmt.Stringer.
func (self *Matcher) String() string {
	return self.description.String()
}

// Tests the given input to see if it meets this Matcher's criteria.
func (self *Matcher) Match(input interface{}) *Result {
	result := self.match(input)
	result.matcher = self
	return result
}

// Returns a slice of messages that supplement the description.
func (self *Matcher) Comments() []interface{} {
	comments := make([]interface{}, len(self.comments))
	copy(comments, self.comments)
	return comments
}

// Returns a *new* Matcher similar to this one, but with the
// given additional comments.
func (self *Matcher) AddComments(comments... interface{}) *Matcher {
	all := make([]interface{}, 0, len(self.comments) + len(comments))
	copy(all, self.comments)
	all = append(all, comments...)
	return &Matcher{description:self.description, match:self.match, comments:all}
}




