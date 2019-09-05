package asserting_test

import (
	"errors"
	"testing"
	"unsafe"

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

func TestAssertTrue(t1 *testing.T) {
	mock := &MockTB{TB: t1}
	t := NewTB(mock)

	t.AssertTrue(true)

	if len(mock.ErrorMessages) != 0 || len(mock.FatalMessages) != 0 {
		t1.Fatal()
	}

	t.AssertTrue(1 == 0)
	if len(mock.FatalMessages) != 0 {
		t1.Fatal()
	}
	if len(mock.ErrorMessages) != 1 ||
		len(mock.ErrorMessages[0]) != 1 ||
		mock.ErrorMessages[0][0] != "unexpected false condition" {
		t1.Fatal(mock.ErrorMessages)
	}
}

func TestAssertNoError(t1 *testing.T) {
	mock := &MockTB{TB: t1}
	t := NewTB(mock)

	t.AssertNoError(nil)

	if len(mock.ErrorMessages) != 0 || len(mock.FatalMessages) != 0 {
		t1.Fatal()
	}

	t.AssertNoError(errors.New("err"))
	if len(mock.FatalMessages) != 0 {
		t1.Fatal()
	}
	if len(mock.ErrorMessages) != 1 ||
		len(mock.ErrorMessages[0]) != 1 ||
		mock.ErrorMessages[0][0] != "unexpected error <err>" {
		t1.Fatal(mock.ErrorMessages)
	}
}

func TestAssertNil(t1 *testing.T) {
	mock := &MockTB{TB: t1}
	t := NewTB(mock)

	t.Assert(nil, Equals(nil))

	if len(mock.ErrorMessages) != 0 || len(mock.FatalMessages) != 0 {
		t1.Fatal()
	}

	t.Assert((*int)(nil), Equals(nil))

	if len(mock.ErrorMessages) != 0 || len(mock.FatalMessages) != 0 {
		t1.Fatal()
	}

	t.Assert(([]int)(nil), Equals(nil))

	if len(mock.ErrorMessages) != 0 || len(mock.FatalMessages) != 0 {
		t1.Fatal()
	}

	t.Assert((map[int]int)(nil), Equals(nil))

	if len(mock.ErrorMessages) != 0 || len(mock.FatalMessages) != 0 {
		t1.Fatal()
	}

	t.Assert((func())(nil), Equals(nil))

	if len(mock.ErrorMessages) != 0 || len(mock.FatalMessages) != 0 {
		t1.Fatal()
	}

	t.Assert((unsafe.Pointer)(nil), Equals(nil))

	if len(mock.ErrorMessages) != 0 || len(mock.FatalMessages) != 0 {
		t1.Fatal()
	}

	t.Assert((chan int)(nil), Equals(nil))

	if len(mock.ErrorMessages) != 0 || len(mock.FatalMessages) != 0 {
		t1.Fatal()
	}

	t.Assert((error)(nil), Equals(nil))

	if len(mock.ErrorMessages) != 0 || len(mock.FatalMessages) != 0 {
		t1.Fatal()
	}

	t.Assert(nil, Equals((*int)(nil)))

	if len(mock.ErrorMessages) != 0 || len(mock.FatalMessages) != 0 {
		t1.Fatal()
	}

	t.Assert((*int)(nil), Equals((*byte)(nil)))

	if len(mock.ErrorMessages) != 1 ||
		len(mock.ErrorMessages[0]) != 1 ||
		mock.ErrorMessages[0][0] != "expected <<nil>(*uint8)> but was <<nil>(*int)>" {
		t1.Fatal(mock.ErrorMessages)
	}

	// Test NotEquals

	mock.ErrorMessages = nil
	mock.FatalMessages = nil

	t.Assert((*byte)(nil), NotEquals((*int)(nil)))

	if len(mock.ErrorMessages) != 0 || len(mock.FatalMessages) != 0 {
		t1.Fatal()
	}

	t.Assert((error)(nil), NotEquals(nil))
	if len(mock.FatalMessages) != 0 {
		t1.Fatal()
	}
	if len(mock.ErrorMessages) != 1 ||
		len(mock.ErrorMessages[0]) != 1 ||
		mock.ErrorMessages[0][0] != "unexpected <<nil>>" {
		t1.Fatal(mock.ErrorMessages)
	}

	mock.ErrorMessages = nil
	mock.FatalMessages = nil
	t.Assert(([]int)(nil), NotEquals(nil))
	if len(mock.FatalMessages) != 0 {
		t1.Fatal()
	}
	if len(mock.ErrorMessages) != 1 ||
		len(mock.ErrorMessages[0]) != 1 ||
		mock.ErrorMessages[0][0] != "unexpected <[]>" {
		t1.Fatal(mock.ErrorMessages)
	}

	mock.ErrorMessages = nil
	mock.FatalMessages = nil
	t.Assert((chan int)(nil), NotEquals(nil))
	if len(mock.FatalMessages) != 0 {
		t1.Fatal()
	}
	if len(mock.ErrorMessages) != 1 ||
		len(mock.ErrorMessages[0]) != 1 ||
		mock.ErrorMessages[0][0] != "unexpected <<nil>>" {
		t1.Fatal(mock.ErrorMessages)
	}

	mock.ErrorMessages = nil
	mock.FatalMessages = nil
	t.Assert((unsafe.Pointer)(nil), NotEquals(nil))
	if len(mock.FatalMessages) != 0 {
		t1.Fatal()
	}
	if len(mock.ErrorMessages) != 1 ||
		len(mock.ErrorMessages[0]) != 1 ||
		mock.ErrorMessages[0][0] != "unexpected <<nil>>" {
		t1.Fatal(mock.ErrorMessages)
	}

	mock.ErrorMessages = nil
	mock.FatalMessages = nil
	t.Assert((func())(nil), NotEquals(nil))
	if len(mock.FatalMessages) != 0 {
		t1.Fatal()
	}
	if len(mock.ErrorMessages) != 1 ||
		len(mock.ErrorMessages[0]) != 1 ||
		mock.ErrorMessages[0][0] != "unexpected <<nil>>" {
		t1.Fatal(mock.ErrorMessages)
	}

	mock.ErrorMessages = nil
	mock.FatalMessages = nil
	t.Assert((map[int]int)(nil), NotEquals(nil))
	if len(mock.FatalMessages) != 0 {
		t1.Fatal()
	}
	if len(mock.ErrorMessages) != 1 ||
		len(mock.ErrorMessages[0]) != 1 ||
		mock.ErrorMessages[0][0] != "unexpected <map[]>" {
		t1.Fatal(mock.ErrorMessages)
	}

	mock.ErrorMessages = nil
	mock.FatalMessages = nil
	t.Assert(([]int)(nil), NotEquals(nil))
	if len(mock.FatalMessages) != 0 {
		t1.Fatal()
	}
	if len(mock.ErrorMessages) != 1 ||
		len(mock.ErrorMessages[0]) != 1 ||
		mock.ErrorMessages[0][0] != "unexpected <[]>" {
		t1.Fatal(mock.ErrorMessages)
	}

	mock.ErrorMessages = nil
	mock.FatalMessages = nil
	t.Assert((*int)(nil), NotEquals(nil))
	if len(mock.FatalMessages) != 0 {
		t1.Fatal()
	}
	if len(mock.ErrorMessages) != 1 ||
		len(mock.ErrorMessages[0]) != 1 ||
		mock.ErrorMessages[0][0] != "unexpected <<nil>>" {
		t1.Fatal(mock.ErrorMessages)
	}

	mock.ErrorMessages = nil
	mock.FatalMessages = nil
	t.Assert(nil, NotEquals(nil))
	if len(mock.FatalMessages) != 0 {
		t1.Fatal()
	}
	if len(mock.ErrorMessages) != 1 ||
		len(mock.ErrorMessages[0]) != 1 ||
		mock.ErrorMessages[0][0] != "unexpected <<nil>>" {
		t1.Fatal(mock.ErrorMessages)
	}

	mock.ErrorMessages = nil
	mock.FatalMessages = nil
	t.Assert(nil, NotEquals((*int)(nil)))
	if len(mock.FatalMessages) != 0 {
		t1.Fatal()
	}
	if len(mock.ErrorMessages) != 1 ||
		len(mock.ErrorMessages[0]) != 1 ||
		mock.ErrorMessages[0][0] != "unexpected <<nil>>" {
		t1.Fatal(mock.ErrorMessages)
	}
}

func TestAssertUntyped(t1 *testing.T) {
	mock := &MockTB{TB: t1}
	t := NewTB(mock)

	t.Assert(100, Equals(UntypedInt(100)))

	if len(mock.ErrorMessages) != 0 || len(mock.FatalMessages) != 0 {
		t1.Fatal(mock.ErrorMessages)
	}

	t.Assert(uint8(100), Equals(UntypedInt(100)))

	if len(mock.ErrorMessages) != 0 || len(mock.FatalMessages) != 0 {
		t1.Fatal(mock.ErrorMessages)
	}

	t.Assert(int16(100), Equals(UntypedInt(100)))

	if len(mock.ErrorMessages) != 0 || len(mock.FatalMessages) != 0 {
		t1.Fatal(mock.ErrorMessages)
	}

	t.Assert(int32(100), Equals(UntypedInt(100)))

	if len(mock.ErrorMessages) != 0 || len(mock.FatalMessages) != 0 {
		t1.Fatal(mock.ErrorMessages)
	}

	t.Assert(int64(100), Equals(UntypedInt(100)))

	if len(mock.ErrorMessages) != 0 || len(mock.FatalMessages) != 0 {
		t1.Fatal(mock.ErrorMessages)
	}

	t.Assert(uint8(100), Equals(UntypedInt(100)))

	if len(mock.ErrorMessages) != 0 || len(mock.FatalMessages) != 0 {
		t1.Fatal(mock.ErrorMessages)
	}

	t.Assert(uint16(100), Equals(UntypedInt(100)))

	if len(mock.ErrorMessages) != 0 || len(mock.FatalMessages) != 0 {
		t1.Fatal(mock.ErrorMessages)
	}

	t.Assert(uint32(100), Equals(UntypedInt(100)))

	if len(mock.ErrorMessages) != 0 || len(mock.FatalMessages) != 0 {
		t1.Fatal(mock.ErrorMessages)
	}

	t.Assert(uint64(100), Equals(UntypedInt(100)))

	if len(mock.ErrorMessages) != 0 || len(mock.FatalMessages) != 0 {
		t1.Fatal(mock.ErrorMessages)
	}

	t.Assert(100, Equals(UntypedUint(100)))

	if len(mock.ErrorMessages) != 0 || len(mock.FatalMessages) != 0 {
		t1.Fatal(mock.ErrorMessages)
	}

	t.Assert(uint8(100), Equals(UntypedUint(100)))

	if len(mock.ErrorMessages) != 0 || len(mock.FatalMessages) != 0 {
		t1.Fatal(mock.ErrorMessages)
	}

	t.Assert(int16(100), Equals(UntypedUint(100)))

	if len(mock.ErrorMessages) != 0 || len(mock.FatalMessages) != 0 {
		t1.Fatal(mock.ErrorMessages)
	}

	t.Assert(int32(100), Equals(UntypedUint(100)))

	if len(mock.ErrorMessages) != 0 || len(mock.FatalMessages) != 0 {
		t1.Fatal(mock.ErrorMessages)
	}

	t.Assert(int64(100), Equals(UntypedUint(100)))

	if len(mock.ErrorMessages) != 0 || len(mock.FatalMessages) != 0 {
		t1.Fatal(mock.ErrorMessages)
	}

	t.Assert(uint8(100), Equals(UntypedUint(100)))

	if len(mock.ErrorMessages) != 0 || len(mock.FatalMessages) != 0 {
		t1.Fatal(mock.ErrorMessages)
	}

	t.Assert(uint16(100), Equals(UntypedUint(100)))

	if len(mock.ErrorMessages) != 0 || len(mock.FatalMessages) != 0 {
		t1.Fatal(mock.ErrorMessages)
	}

	t.Assert(uint32(100), Equals(UntypedUint(100)))

	if len(mock.ErrorMessages) != 0 || len(mock.FatalMessages) != 0 {
		t1.Fatal(mock.ErrorMessages)
	}

	t.Assert(uint64(100), Equals(UntypedUint(100)))

	if len(mock.ErrorMessages) != 0 || len(mock.FatalMessages) != 0 {
		t1.Fatal(mock.ErrorMessages)
	}

	t.Assert(-100, Equals(UntypedInt(-100)))

	if len(mock.ErrorMessages) != 0 || len(mock.FatalMessages) != 0 {
		t1.Fatal(mock.ErrorMessages)
	}

	t.Assert(int8(-100), Equals(UntypedInt(-100)))

	if len(mock.ErrorMessages) != 0 || len(mock.FatalMessages) != 0 {
		t1.Fatal(mock.ErrorMessages)
	}

	t.Assert(int16(-100), Equals(UntypedInt(-100)))

	if len(mock.ErrorMessages) != 0 || len(mock.FatalMessages) != 0 {
		t1.Fatal(mock.ErrorMessages)
	}

	t.Assert(int32(-100), Equals(UntypedInt(-100)))

	if len(mock.ErrorMessages) != 0 || len(mock.FatalMessages) != 0 {
		t1.Fatal(mock.ErrorMessages)
	}

	t.Assert(int64(-100), Equals(UntypedInt(-100)))

	if len(mock.ErrorMessages) != 0 || len(mock.FatalMessages) != 0 {
		t1.Fatal(mock.ErrorMessages)
	}

	t.Assert(float32(-100), Equals(UntypedFloat(-100)))

	if len(mock.ErrorMessages) != 0 || len(mock.FatalMessages) != 0 {
		t1.Fatal(mock.ErrorMessages)
	}

	t.Assert(float64(-100), Equals(UntypedFloat(-100)))

	if len(mock.ErrorMessages) != 0 || len(mock.FatalMessages) != 0 {
		t1.Fatal(mock.ErrorMessages)
	}

	t.Assert(UntypedFloat(-100), Equals(float32(-100)))

	if len(mock.ErrorMessages) != 0 || len(mock.FatalMessages) != 0 {
		t1.Fatal(mock.ErrorMessages)
	}

	t.Assert(UntypedFloat(-100), Equals(float64(-100)))

	if len(mock.ErrorMessages) != 0 || len(mock.FatalMessages) != 0 {
		t1.Fatal(mock.ErrorMessages)
	}

	t.Assert(UntypedInt(100), Equals(uint16(100)))

	if len(mock.ErrorMessages) != 0 || len(mock.FatalMessages) != 0 {
		t1.Fatal(mock.ErrorMessages)
	}

	t.Assert(UntypedInt(100), Equals(UntypedInt(100)))

	if len(mock.ErrorMessages) != 0 || len(mock.FatalMessages) != 0 {
		t1.Fatal(mock.ErrorMessages)
	}

	t.Assert(UntypedInt(100), Equals(UntypedUint(100)))

	if len(mock.ErrorMessages) != 0 || len(mock.FatalMessages) != 0 {
		t1.Fatal(mock.ErrorMessages)
	}

	t.Assert("abc", Equals(UntypedString("abc")))

	if len(mock.ErrorMessages) != 0 || len(mock.FatalMessages) != 0 {
		t1.Fatal(mock.ErrorMessages)
	}

	t.Assert(UntypedString("abc"), Equals("abc"))

	if len(mock.ErrorMessages) != 0 || len(mock.FatalMessages) != 0 {
		t1.Fatal(mock.ErrorMessages)
	}

	t.Assert(UntypedString("abc"), Equals(UntypedString("abc")))

	if len(mock.ErrorMessages) != 0 || len(mock.FatalMessages) != 0 {
		t1.Fatal(mock.ErrorMessages)
	}

	t.Assert(complex(1, 2), Equals(UntypedComplex(complex(1, 2))))

	if len(mock.ErrorMessages) != 0 || len(mock.FatalMessages) != 0 {
		t1.Fatal(mock.ErrorMessages)
	}

	t.Assert(UntypedComplex(complex(1, 2)), Equals(complex(1, 2)))

	if len(mock.ErrorMessages) != 0 || len(mock.FatalMessages) != 0 {
		t1.Fatal(mock.ErrorMessages)
	}

	t.Assert(UntypedComplex(complex(1, 2)), Equals(UntypedComplex(complex(1, 2))))

	if len(mock.ErrorMessages) != 0 || len(mock.FatalMessages) != 0 {
		t1.Fatal(mock.ErrorMessages)
	}

	t.Assert(UntypedComplex(complex(1, 2)), NotEquals(UntypedComplex(complex(1, 3))))

	if len(mock.ErrorMessages) != 0 || len(mock.FatalMessages) != 0 {
		t1.Fatal(mock.ErrorMessages)
	}

	t.Assert(float32(123), Equals(UntypedUint(123)))

	if len(mock.ErrorMessages) != 0 || len(mock.FatalMessages) != 0 {
		t1.Fatal(mock.ErrorMessages)
	}

	t.Assert(UntypedInt(-123), Equals(float32(-123)))

	if len(mock.ErrorMessages) != 0 || len(mock.FatalMessages) != 0 {
		t1.Fatal(mock.ErrorMessages)
	}

	t.Assert(nil, NotEquals(UntypedInt(-123)))

	if len(mock.ErrorMessages) != 0 || len(mock.FatalMessages) != 0 {
		t1.Fatal(mock.ErrorMessages)
	}

	t.Assert(1, Equals(UntypedInt(2)))
	if len(mock.FatalMessages) != 0 {
		t1.Fatal()
	}
	if len(mock.ErrorMessages) != 1 ||
		len(mock.ErrorMessages[0]) != 1 ||
		mock.ErrorMessages[0][0] != "expected <2> but was <1>" {
		t1.Fatal(mock.ErrorMessages)
	}

	mock.ErrorMessages = nil
	mock.FatalMessages = nil

	t.Assert(1, Equals(UntypedFloat(2)))
	if len(mock.FatalMessages) != 0 {
		t1.Fatal()
	}
	if len(mock.ErrorMessages) != 1 ||
		len(mock.ErrorMessages[0]) != 1 ||
		mock.ErrorMessages[0][0] != "expected <2> but was <1>" {
		t1.Fatal(mock.ErrorMessages)
	}

	mock.ErrorMessages = nil
	mock.FatalMessages = nil

	t.Assert("abc", Equals(UntypedString("def")))
	if len(mock.FatalMessages) != 0 {
		t1.Fatal()
	}
	if len(mock.ErrorMessages) != 1 ||
		len(mock.ErrorMessages[0]) != 1 ||
		mock.ErrorMessages[0][0] != "expected <def> but was <abc>" {
		t1.Fatal(mock.ErrorMessages)
	}
}
