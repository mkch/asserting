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
		c := Equals(nil).SetMessage(err.message)
		if err.fatal {
			c.SetFatal()
		}
		t.Assert(0, c)
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
	t.Helper()
	t.Assert(condition, Equals(true).SetMessage("unexpected false condition"))
}

// AssertNoError asserts the err is nil.
func (t TB) AssertNoError(err error) {
	t.Helper()
	t.Assert(err, Equals(nil).SetMessage(fmt.Sprintf("unexpected error <%v>", err)))
}

// AssertEqual calls t.Assert(v, Equals(expected)).
func (t TB) AssertEqual(v, expected interface{}) {
	t.Helper()
	t.Assert(v, Equals(expected))
}

// AssertEqualSlice calls t.Assert(v, EqualsSlice(expected)).
func (t TB) AssertEqualSlice(v, expected interface{}) {
	t.Helper()
	t.Assert(v, EqualsSlice(expected))
}

// AssertNotEqual calls t.Assert(v, NotEquals(expected)).
func (t TB) AssertNotEqual(v, expected interface{}) {
	t.Helper()
	t.Assert(v, NotEquals(expected))
}

// AssertMatch calls t.Assert(v, Matches(f)).
func (t TB) AssertMatch(v interface{}, f func(v interface{}) bool) {
	t.Helper()
	t.Assert(v, Matches(f))
}

// AssertPanic calls t.Assert(v, Panics(expected)).
func (t TB) AssertPanic(v func(), expected interface{}) {
	t.Helper()
	t.Assert(v, Panics(expected))
}

// AssertPanicMatch calls t.Assert(v, PanicMatches(f)).
func (t TB) AssertPanicMatch(v func(), f func(expected interface{}) bool) {
	t.Helper()
	t.Assert(v, PanicMatches(f))
}

type hasError struct {
	message string
	fatal   bool
}

// ValueError converts v and err to a single value.
//
// If TB.Assert(ValueError(v, err), cond)
// is called,  one of the following 2 things will happen:
//
// 1. If err is not nil, the assertion fails with t.Error("unexpected error ...").
//
// 2. If err is nil, the code is executed the same way as TB.Assert(v, cond)
func ValueError(v interface{}, err error) interface{} {
	if err != nil {
		return &hasError{message: fmt.Sprintf("unexpected error <%v>", err)}
	}
	return v
}

// ValueErrorFatal is equivalent to ValueError except one thing:
// 1. If err is not nil, the assertion fails with t.Fatal("unexpected error ...").
func ValueErrorFatal(v interface{}, err error) interface{} {
	if err != nil {
		return &hasError{message: fmt.Sprintf("unexpected error <%v>", err), fatal: true}
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
	return eq(c.expected, v)
}

func (c *equals) Message(v interface{}) string {
	return formatMsg("expected <%v> but was <%v>", c.expected, v)
}

type notEquals equals

// NotEquals returns a cond which is true if a value does not equal to the expected value.
// The inequality is determined with operator !=
func NotEquals(unexpected interface{}) cond.Cond {
	return cond.New((*notEquals)(&notEquals{expected: unexpected}))
}

func (c *notEquals) Test(v interface{}) bool {
	return !((*equals)(c)).Test(v)
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
// Test() panics if a the tested value is not of type func() when this kind of cond
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
		result = eq(c.expected, c.got)
	}()

	f()

	return
}

func (c *panics) Message(v interface{}) string {
	nilExplain := ""
	if c.got == nil {
		nilExplain = " (didn't panic?)"
	}
	return formatMsg("expected to panic with <%v> but <%v>"+nilExplain, c.expected, c.got)
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
	return formatMsg("expected <%v> but was <%v>", c.expected, v)
}

type untypedInt int64

func (i untypedInt) equals(r interface{}) bool {
	tr := reflect.TypeOf(r)
	if tr == nil {
		return false
	}
	switch tr.Kind() {
	case reflect.Int:
		fallthrough
	case reflect.Int8:
		fallthrough
	case reflect.Int16:
		fallthrough
	case reflect.Int32:
		fallthrough
	case reflect.Int64:
		return int64(i) == reflect.ValueOf(r).Int()
	case reflect.Uint:
		fallthrough
	case reflect.Uint8:
		fallthrough
	case reflect.Uint16:
		fallthrough
	case reflect.Uint32:
		fallthrough
	case reflect.Uint64:
		return int64(i) >= 0 && uint64(i) == reflect.ValueOf(r).Uint()
	case reflect.Float32:
		fallthrough
	case reflect.Float64:
		return float64(i) == reflect.ValueOf(r).Float()
	default:
		return false
	}
}

// UntypedInt returns an untyped integer which equals other integer or float types
// if they have the same value.
func UntypedInt(n int64) interface{} {
	return untypedInt(n)
}

type untypedUint uint64

func (i untypedUint) equals(r interface{}) bool {
	tr := reflect.TypeOf(r)
	if tr == nil {
		return false
	}
	switch tr.Kind() {
	case reflect.Uint:
		fallthrough
	case reflect.Uint8:
		fallthrough
	case reflect.Uint16:
		fallthrough
	case reflect.Uint32:
		fallthrough
	case reflect.Uint64:
		return uint64(i) == reflect.ValueOf(r).Uint()
	case reflect.Int:
		fallthrough
	case reflect.Int8:
		fallthrough
	case reflect.Int16:
		fallthrough
	case reflect.Int32:
		fallthrough
	case reflect.Int64:
		v := reflect.ValueOf(r).Int()
		return v >= 0 && uint64(i) == uint64(v)
	case reflect.Float32:
		fallthrough
	case reflect.Float64:
		return float64(i) == reflect.ValueOf(r).Float()
	default:
		return false
	}
}

type ieq interface {
	equals(r interface{}) bool
}

// UntypedUint returns an untyped integer which is reported by Assert equal to
// values of integer or float types if they have the same value.
func UntypedUint(n int64) interface{} {
	return untypedUint(n)
}

type untypedFloat float64

func (f untypedFloat) equals(r interface{}) bool {
	tr := reflect.TypeOf(r)
	if tr == nil {
		return false
	}
	switch tr.Kind() {
	case reflect.Float32:
		fallthrough
	case reflect.Float64:
		return float64(f) == reflect.ValueOf(r).Float()
	default:
		return false
	}
}

// UntypedFloat returns an untyped float point value which is reported by Assert equal to
// values of float32 or float64 types if they have the same value.
func UntypedFloat(f float64) interface{} {
	return untypedFloat(f)
}

type untypedString string

func (s untypedString) equals(r interface{}) bool {
	tr := reflect.TypeOf(r)
	if tr == nil {
		return false
	}
	switch tr.Kind() {
	case reflect.String:
		return string(s) == reflect.ValueOf(r).String()
	default:
		return false
	}
}

// UntypedString returns an untyped string value which is reported by Assert equal to
// values of string types if they have the same value.
func UntypedString(str string) interface{} {
	return untypedString(str)
}

type untypedComplex complex128

func (c untypedComplex) equals(r interface{}) bool {
	tr := reflect.TypeOf(r)
	if tr == nil {
		return false
	}
	switch tr.Kind() {
	case reflect.Complex64:
		fallthrough
	case reflect.Complex128:
		return complex128(c) == reflect.ValueOf(r).Complex()
	default:
		return false
	}
}

// UntypedComplex returns an untyped complex value which is reported by Assert equal to
// values of complex64 or complex128 types if have the same value.
func UntypedComplex(c complex128) interface{} {
	return untypedComplex(c)
}

func eq(a, b interface{}) bool {
	if a == b {
		return true
	}

	if a == nil {
		return equalsNil(b)
	}

	if b == nil {
		return equalsNil(a)
	}

	if ieq, ok := a.(ieq); ok {
		return ieq.equals(b)
	}

	if ieq, ok := b.(ieq); ok {
		return ieq.equals(a)
	}

	return false
}

// equalsNil tests whether v is a nil interface value or the value of v == nil.
func equalsNil(v interface{}) bool {
	t := reflect.TypeOf(v)
	if t == nil {
		return true
	}
	switch t.Kind() {
	case reflect.Chan:
		fallthrough
	case reflect.Func:
		fallthrough
	case reflect.Interface:
		fallthrough
	case reflect.Map:
		fallthrough
	case reflect.Slice:
		fallthrough
	case reflect.Ptr:
		fallthrough
	case reflect.UnsafePointer:
		return reflect.ValueOf(v).IsNil()
	default:
		return false
	}
}

func formatMsg(format string, arg1, arg2 interface{}) string {
	str1, str2 := fmt.Sprintf("%v", arg1), fmt.Sprintf("%v", arg2)
	if str1 == str2 {
		arg1, arg2 = fmt.Sprintf("%[1]v(%[1]T)", arg1), fmt.Sprintf("%[1]v(%[1]T)", arg2)
	} else {
		arg1, arg2 = str1, str2
	}
	return fmt.Sprintf(format, arg1, arg2)
}
