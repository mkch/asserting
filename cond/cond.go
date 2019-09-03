// Package cond defines the assertion condition.
package cond

// Condition is a condition with failure message.
type Condition interface {
	// Test returns whether the condition is met.
	Test(v interface{}) bool
	// Message returns the failure message.
	// Message will be called only when Test returns false.
	Message(v interface{}) string
}

// Cond is a condition used by assert.TB.Assert.
// The assertion succeeds if Condition.Test returns true, fails otherwise.
// If the assertion fails, the failure message will be reported
// with testing.TB.Error or testing.TB.Fatal, see the document of SetFatal.
type Cond interface {
	Condition
	// SetMessage replaces the default failure message, overwriting function set by
	// SetMessageFunc if any.
	SetMessage(msg string) Cond
	// SetMessageFunc sets f as the failure message generator, overwriting message set
	// by SetMessage if any.
	// If necessary, the failure message will be retrieved lazily from f.
	SetMessageFunc(f func() string) Cond
	// SetFatal indicates the assertion to use TB.Fatal() instead of TB.Error() in the testing package
	// of go standard library to report failures.
	SetFatal() Cond
	fatal() bool
	message(v interface{}) string
}

type cond struct {
	Condition
	userMsg func() string
	isFatal bool
}

func (c *cond) SetMessage(msg string) Cond {
	c.userMsg = func() string { return msg }
	return c
}

func (c *cond) SetMessageFunc(f func() string) Cond {
	c.userMsg = f
	return c
}

func (c *cond) SetFatal() Cond {
	c.isFatal = true
	return c
}

func (c *cond) fatal() bool {
	return c.isFatal
}

func (c *cond) message(v interface{}) string {
	if c.userMsg != nil {
		return c.userMsg()
	}
	return c.Message(v)
}

// Fatal returns whether cond.Fatal has been called.
func Fatal(cond Cond) bool {
	return cond.fatal()
}

// Message returns the failure message.
// If a user defined message has been set with cond.SetMessage(msg) or cond.SetMessageFunc(f),
// returns the msg or f(). Returns cond.Message(v) otherwise.
func Message(cond Cond, v interface{}) string {
	return cond.message(v)
}

// New creates a Cond with c.
func New(c Condition) Cond {
	return &cond{Condition: c}
}
