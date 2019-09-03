package asserting_test

import (
	"errors"
	"testing"

	. "github.com/mkch/asserting"
)

type MockTB struct {
	testing.TB
	ErrorMessages [][]interface{}
	FatalMessages [][]interface{}
	failed        bool
}

func (m *MockTB) Error(args ...interface{}) {
	if m.failed {
		return
	}
	m.ErrorMessages = append(m.ErrorMessages, args)
}

func (m *MockTB) Fatal(args ...interface{}) {
	if m.failed {
		return
	}
	m.FatalMessages = append(m.FatalMessages, args)
	m.failed = true
}

func TestEquals(t1 *testing.T) {
	mock := &MockTB{TB: t1}
	t := TB{mock}

	t.Assert(1, Equals(1))
	if len(mock.ErrorMessages) != 0 || len(mock.FatalMessages) != 0 {
		t1.Fatal()
	}

	t.Assert(1, Equals(2))
	if len(mock.FatalMessages) != 0 {
		t1.Fatal()
	}
	if len(mock.ErrorMessages) != 1 ||
		len(mock.ErrorMessages[0]) != 1 ||
		mock.ErrorMessages[0][0] != "expected <2> but was <1>" {
		t1.Fatal(mock.ErrorMessages)
	}
}

func TestNotEquals(t1 *testing.T) {
	mock := &MockTB{TB: t1}
	t := TB{mock}

	t.Assert(1, NotEquals("abc"))
	if len(mock.ErrorMessages) != 0 || len(mock.FatalMessages) != 0 {
		t1.Fatal()
	}

	t.Assert(1, NotEquals(1))
	if len(mock.FatalMessages) != 0 {
		t1.Fatal()
	}
	if len(mock.ErrorMessages) != 1 ||
		len(mock.ErrorMessages[0]) != 1 ||
		mock.ErrorMessages[0][0] != "unexpected <1>" {
		t1.Fatal(mock.ErrorMessages)
	}
}

func TestMatches(t1 *testing.T) {
	mock := &MockTB{TB: t1}
	t := TB{mock}

	t.Assert(1, Matches(func(v interface{}) bool { return v == 1 }))
	if len(mock.ErrorMessages) != 0 || len(mock.FatalMessages) != 0 {
		t1.Fatal()
	}

	t.Assert(1, Matches(func(v interface{}) bool { return v != 1 }))
	t.Assert("abc", Matches(func(v interface{}) bool { return len(v.(string)) == 0 }))
	if len(mock.FatalMessages) != 0 {
		t1.Fatal()
	}
	if len(mock.ErrorMessages) != 2 ||
		len(mock.ErrorMessages[0]) != 1 ||
		len(mock.ErrorMessages[1]) != 1 ||
		mock.ErrorMessages[0][0] != "unexpected <1>" ||
		mock.ErrorMessages[1][0] != "unexpected <abc>" {
		t1.Fatal(mock.ErrorMessages)
	}
}

func TestPanics(t1 *testing.T) {
	mock := &MockTB{TB: t1}
	t := TB{mock}

	t.Assert(func() { panic(1) }, Panics(1))
	if len(mock.ErrorMessages) != 0 || len(mock.FatalMessages) != 0 {
		t1.Fatal()
	}

	t.Assert(func() { panic(2) }, Panics(1))
	if len(mock.FatalMessages) != 0 {
		t1.Fatal()
	}
	if len(mock.ErrorMessages) != 1 ||
		len(mock.ErrorMessages[0]) != 1 ||
		mock.ErrorMessages[0][0] != "expected to panic with <1> but <2>" {
		t1.Fatal(mock.ErrorMessages)
	}
}

func TestPanicMatches(t1 *testing.T) {
	mock := &MockTB{TB: t1}
	t := TB{mock}

	t.Assert(func() { panic(1) }, PanicMatches(func(v interface{}) bool { return v == 1 }))
	if len(mock.ErrorMessages) != 0 || len(mock.FatalMessages) != 0 {
		t1.Fatal()
	}

	t.Assert(func() { panic(2) }, PanicMatches(func(v interface{}) bool { return v == 1 }))
	if len(mock.FatalMessages) != 0 {
		t1.Fatal()
	}
	if len(mock.ErrorMessages) != 1 ||
		len(mock.ErrorMessages[0]) != 1 ||
		mock.ErrorMessages[0][0] != "unexpected panic <2>" {
		t1.Fatal(mock.ErrorMessages)
	}
}

func TestEqualsSlice(t1 *testing.T) {
	mock := &MockTB{TB: t1}
	t := TB{mock}

	t.Assert([]int{1, 2, 3}, EqualsSlice([]int{1, 2, 3}))
	if len(mock.ErrorMessages) != 0 || len(mock.FatalMessages) != 0 {
		t1.Fatal()
	}

	t.Assert([]int{1, 2, 3}, EqualsSlice([]int{1, 2}))
	if len(mock.FatalMessages) != 0 {
		t1.Fatal()
	}
	if len(mock.ErrorMessages) != 1 ||
		len(mock.ErrorMessages[0]) != 1 ||
		mock.ErrorMessages[0][0] != "expected <[1 2]> but was <[1 2 3]>" {
		t1.Fatal(mock.ErrorMessages)
	}
}

func TestValueError(t1 *testing.T) {
	mock := &MockTB{TB: t1}
	t := TB{mock}

	t.Assert(
		ValueError(func() (int, error) { return 1, nil }()),
		Equals(1))
	if len(mock.ErrorMessages) != 0 || len(mock.FatalMessages) != 0 {
		t1.Fatal()
	}

	t.Assert(
		ValueError(func() (int, error) { return 1, errors.New("error") }()),
		Equals(1))
	if len(mock.FatalMessages) != 0 {
		t1.Fatal()
	}
	if len(mock.ErrorMessages) != 1 ||
		len(mock.ErrorMessages[0]) != 1 ||
		mock.ErrorMessages[0][0] != "unexpected error <error>" {
		t1.Fatal(mock.ErrorMessages)
	}
}

func TestFatal(t1 *testing.T) {
	mock := &MockTB{TB: t1}
	t := TB{mock}

	t.Assert(1, Equals(2).SetFatal())
	t.Assert(1, Equals(3))

	if len(mock.ErrorMessages) != 0 {
		t1.Fatal()
	}
	if len(mock.FatalMessages) != 1 ||
		len(mock.FatalMessages[0]) != 1 ||
		mock.FatalMessages[0][0] != "expected <2> but was <1>" {
		t1.Fatal(mock.FatalMessages)
	}
}
