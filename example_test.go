package asserting_test

import (
	"strconv"
	"testing"

	. "github.com/mkch/asserting"
)

func Example() {
	// See the test functions below.
	// The name of test functions begin with underscores to be shown as examples.
}

func _TestAdd(t1 *testing.T) {
	t := TB{t1}
	// Asserts 1+1 == 2
	t.Assert(1+1, Equals(2))
	// Asserts 1+1 != 0 with custom failure message.
	t.Assert(1+1, NotEquals(0).SetMessage("1+1 != 0"))
}

func SomeOddNumber() int {
	return 1333
}

func _TestSomeOddNumber(t1 *testing.T) {
	t := TB{t1}
	// Asserts SomeOddNumber() returns an odd number.
	t.Assert(SomeOddNumber(), Matches(
		func(v interface{}) bool {
			return v.(int)%2 != 0
		}).
		SetMessage("Not an odd number"))
}

func PanicWith100() {
	panic(100)
}

func Div(a, b int) int {
	if b == 0 {
		panic("can' div by 0")
	}
	return a / b
}

func _TestPanicWith100(t1 *testing.T) {
	t := TB{t1}
	// Asserts calling a function must panic with 100.
	t.Assert(PanicWith100, Panics(100))
	// Asserts calling a function must panic with a string.
	t.Assert(func() { Div(1, 0) },
		PanicMatches(
			func(v interface{}) bool {
				_, ok := v.(string)
				return ok
			}))
}

func _TestAtoi(t1 *testing.T) {
	t := TB{t1}
	// Test strconv.Atoi who returns an int and an error.
	// If the error value of Atoi is not nil, or the int
	// value is not 1, the assertion fails.
	t.Assert(ValueError(strconv.Atoi("1")), Equals(1))
}
