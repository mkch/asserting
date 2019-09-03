// Package asserting is a utility package to do unit testing.
package asserting

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/mkch/asserting/cond"
)

// TB is a wrapper of testing.TB to do assertion.
type TB struct {
	testing.TB
}

// NewTB creates a TB.
func NewTB(t testing.TB) TB {
	return TB{t}
}

// Assert asserts v meets the condition c.
// If v does not meet c, the assertion fails and a failure message
// is reported. See the document of cond.Cond.
func (t TB) Assert(v interface{}, c cond.Cond) {
	t.Helper()
	if err, ok := v.(*hasError); ok {
		t.Assert(0, Equals(nil).SetMessage(err.message))
		return
	}
	if !c.Test(v) {
		f := t.Error
		if cond.Fatal(c) {
			f = t.Fatal
		}
		f(cond.Message(c, v))
	}
}

// AssertTrue asserts the condition is true.
func (t TB) AssertTrue(condition bool) {
	t.Assert(condition, Equals(true).SetMessage("unexpected false condition"))
}

// AssertNoError asserts the err is nil.
func (t TB) AssertNoError(err error) {
	t.Assert(err, Equals(nil).SetMessage(fmt.Sprintf("unexpected error <%v>", err)))
}

type hasError struct {
	message string
}

// ValueError converts v and err to a single value.
//
// If TB.Assert(ValueError(v, err), cond)
// is called,  one of the following 2 things will happen:
//
// 1. If err is not nil, the assertion fails with message "unexpected error".
//
// 2. If err is not nil, the code is executed the same way as TB.Assert(v, cond)
func ValueError(v interface{}, err error) interface{} {
	if err != nil {
		return &hasError{fmt.Sprintf("unexpected error <%v>", err)}
	}
	return v
}

type equals struct {
	expected interface{}
}

// Equals returns a cond which is true if a value equals to the expected value.
// The equality is determined with operator ==.
func Equals(expected interface{}) cond.Cond {
	return cond.New(&equals{expected: expected})
}

func (c *equals) Test(v interface{}) bool {
	return v == c.expected
}

func (c *equals) Message(v interface{}) string {
	return fmt.Sprintf("expected <%v> but was <%v>", c.expected, v)
}

type notEquals struct {
	unexpected interface{}
}

// NotEquals returns a cond which is true if a value does not equal to the expected value.
// The inequality is determined with operator !=
func NotEquals(unexpected interface{}) cond.Cond {
	return cond.New(&notEquals{unexpected: unexpected})
}

func (c *notEquals) Test(v interface{}) bool {
	return v != c.unexpected
}

func (c *notEquals) Message(v interface{}) string {
	return fmt.Sprintf("unexpected <%v>", v)
}

type matches struct {
	f func(v interface{}) bool
}

// Matches returns a cond which is true if a value passes the test of function f.
func Matches(f func(v interface{}) bool) cond.Cond {
	return cond.New(&matches{f: f})
}

func (c *matches) Test(v interface{}) bool {
	return c.f(v)
}

func (c *matches) Message(v interface{}) string {
	return fmt.Sprintf("unexpected <%v>", v)
}

type panics struct {
	expected interface{}
	got      interface{} // The actual recovered value.
}

// Panics returns a cond which is true if the tested function panics with the expected value.
// TB.Assert() panics if a the tested value is not of type func() when this kind of cond
// is used.
func Panics(expected interface{}) cond.Cond {
	return cond.New(&panics{expected: expected})
}

func (c *panics) Test(v interface{}) (result bool) {
	f, ok := v.(func())
	if !ok {
		panic(fmt.Sprintf("<%v> is not a func()", v))
	}

	defer func() {
		c.got = recover()
		result = c.expected == c.got
	}()

	f()

	return
}

func (c *panics) Message(v interface{}) string {
	nilExplain := ""
	if c.got == nil {
		nilExplain = " (didn't panic?)"
	}
	return fmt.Sprintf("expected to panic with <%v> but <%v>"+nilExplain, c.expected, c.got)
}

type panicMatches struct {
	got interface{} // The actual recovered value.
	f   func(interface{}) bool
}

// PanicMatches returns a cond which is true if the tested function panics with a value that passes
// the test of function f.
// TB.Assert() panics if a the tested value is not of type func() when this kind of cond
// is used.
func PanicMatches(f func(interface{}) bool) cond.Cond {
	return cond.New(&panicMatches{f: f})
}

func (c *panicMatches) Test(v interface{}) (result bool) {
	f, ok := v.(func())
	if !ok {
		panic(fmt.Sprintf("<%v> is not a func()", v))
	}

	defer func() {
		c.got = recover()
		result = c.f(c.got)
	}()

	f()

	return
}

func (c *panicMatches) Message(v interface{}) string {
	nilExplain := ""
	if c.got == nil {
		nilExplain = " (didn't panic?)"
	}
	return fmt.Sprintf("unexpected panic <%v>"+nilExplain, c.got)
}

type equalsSlice struct {
	expected interface{}
}

// EqualsSlice returns a cond which is true if the tested slice equals to the expected slice.
// TB.Assert() panics if a the tested value and the expected value are not of the same slice
// type or nil when this kind of cond is used.
// The equality is defined by the following 2 rules:
//
// nil equals to empty slice.
//
// 2 non nil slices a and b equals to each other if reflect.DeepEqual(a, b) returns true.
func EqualsSlice(expected interface{}) cond.Cond {
	return cond.New(&equalsSlice{expected: expected})
}

func (c *equalsSlice) Test(v interface{}) bool {
	t1 := reflect.TypeOf(v)
	if t1 != nil && t1.Kind() != reflect.Slice {
		panic(fmt.Sprintf("testing2: <%[1]v(%[1]T)> is not a slice", v))
	}

	t2 := reflect.TypeOf(c.expected)
	if t2 != nil && t2.Kind() != reflect.Slice {
		panic(fmt.Sprintf("testing2: <%[1]v(%[1]T)> is not a slice", c.expected))
	}

	v1 := reflect.ValueOf(v)
	v2 := reflect.ValueOf(c.expected)

	if t1 == nil {
		if t2 != nil && v2.Len() != 0 {
			return false
		}
		return true
	}

	if t2 == nil {
		if t1 != nil && v1.Len() != 0 {
			return false
		}
		return true
	}

	if t1 != t2 {
		panic(fmt.Sprintf("type mismatch: <%v> and <%v>", t1, t2))
	}

	return reflect.DeepEqual(v, c.expected)
}

func (c *equalsSlice) Message(v interface{}) string {
	return fmt.Sprintf("expected <%v> but was <%v>", c.expected, v)
}
