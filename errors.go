package apierrors

import (
	"errors"
	"fmt"
	"runtime"
)

func New(status, code int, message string) *Error {
	return &Error{
		Code:    code,
		Status:  status,
		Message: message,
		frame:   caller(1),
	}
}

type Error struct {
	Code    int
	Status  int
	Message string

	err   error
	base  error
	frame frame
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s", e.Message)
}

func (e *Error) Unwrap() error {
	return e.err
}

func (e *Error) Format(s fmt.State, v rune) {
	if e.base != nil {
		_, _ = fmt.Fprintf(s, "%s: %s\n", e.Error(), e.err)
	}
	e.frame.Format(s, v)
}

func (e *Error) Trace() error {
	err := *e
	err.base = e
	err.frame = caller(1)
	return &err
}

func (e *Error) Wrap(next error) error {
	err := *e
	err.base = e
	err.err = next
	err.frame = caller(1)
	return &err
}

func (e *Error) Swrap(message string) error {
	err := *e
	err.base = e
	err.err = errors.New(message)
	err.frame = caller(1)
	return &err
}

func (e *Error) Swrapf(format string, a ...interface{}) error {
	err := *e
	err.base = e
	err.err = fmt.Errorf(format, a...)
	err.frame = caller(1)
	return &err
}

// A frame contains part of a call stack.
type frame struct {
	// https://go.googlesource.com/go/+/032678e0fb/src/runtime/extern.go#169
	frames [3]uintptr
}

// caller returns a frame that describes a frame on the caller's stack.
func caller(skip int) frame {
	var s frame
	runtime.Callers(skip+1, s.frames[:])
	return s
}

// location reports the file, line, and function of a frame.
func (f frame) location() (function, file string, line int) {
	frames := runtime.CallersFrames(f.frames[:])
	if _, ok := frames.Next(); !ok {
		return "", "", 0
	}
	fr, ok := frames.Next()
	if !ok {
		return "", "", 0
	}
	return fr.Function, fr.File, fr.Line
}

// Format prints the stack as error detail.
func (f frame) Format(p fmt.State, v rune) {
	function, file, line := f.location()
	if function != "" {
		_, _ = fmt.Fprintf(p, "%s\n", function)
	}
	if file != "" {
		_, _ = fmt.Fprintf(p, "%s:%d\n", file, line)
	}
}
