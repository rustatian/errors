package errors

import (
	"bytes"
	"encoding"
	"fmt"
	"log"
	"runtime"
)

type Error struct {
	Op   Op
	Kind Kind
	Err  error

	// Stack information
	stack
}

func (e *Error) isZero() bool {
	return e.Op == "" && e.Kind == 0 && e.Err == nil
}

var (
	_ error                      = (*Error)(nil)
	_ encoding.BinaryUnmarshaler = (*Error)(nil)
	_ encoding.BinaryMarshaler   = (*Error)(nil)
)

// Op describes an operation
type Op string

var Separator = ":\n\t"

type Kind uint8

// Kinds of errors.
const (
	Undefined Kind = iota // Undefined error.
	Timeout
	IncorrectAuthToken
	ActionForbidden
	EntityNotFound
	MissingFields
	FieldValidation
	IncorrectFieldType
	CachingInProgress
	IncorrectHTTPRequest
	ServiceUnavailable
)

func (k Kind) String() string {
	switch k {
	case Undefined:
		return "UNDEF"
	case Timeout:
		return "TO001"
	case IncorrectAuthToken:
		return "AUTH001"
	case ActionForbidden:
		return "AUTH002"
	case EntityNotFound:
		return "F001"
	case MissingFields:
		return "VAL001"
	case FieldValidation:
		return "VAL002"
	case IncorrectFieldType:
		return "VAL003"
	case CachingInProgress:
		return "CIP001"
	case IncorrectHTTPRequest:
		return "HTTP001"
	case ServiceUnavailable:
		return "SVC001"
	}
	return "unknown error kind"
}

// E builds an error value from its arguments.
func E(args ...interface{}) error {
	if len(args) == 0 {
		panic("call to errors.E with no arguments")
	}
	e := &Error{}
	for _, arg := range args {
		switch arg := arg.(type) {
		case Op:
			e.Op = arg
		case string:
			e.Err = Str(arg)
		case Kind:
			e.Kind = arg
		case *Error:
			// Make a copy
			eCopy := *arg
			e.Err = &eCopy
		case error:
			e.Err = arg
			// add map map[string]string
		default:
			_, file, line, _ := runtime.Caller(1)
			log.Printf("errors.E: bad call from %s:%d: %v", file, line, args)
			return Errorf("unknown type %T, value %v in error call", arg, arg)
		}
	}

	// Populate stack information
	e.populateStack()

	prev, ok := e.Err.(*Error)
	if !ok {
		return e
	}

	if prev.Kind == e.Kind {
		prev.Kind = Undefined
	}

	if e.Kind == Undefined {
		e.Kind = prev.Kind
		prev.Kind = Undefined
	}
	return e
}

func appendStrToBuf(b *bytes.Buffer, str string) {
	if b.Len() == 0 {
		return
	}
	b.WriteString(str)
}

func (e *Error) Error() string {
	b := new(bytes.Buffer)
	e.printStack(b)
	if e.Op != "" {
		appendStrToBuf(b, ": ")
		b.WriteString(string(e.Op))
	}

	if e.Kind != 0 {
		appendStrToBuf(b, ": ")
		b.WriteString(e.Kind.String())
	}
	if e.Err != nil {
		if prevErr, ok := e.Err.(*Error); ok {
			if !prevErr.isZero() {
				// indent - separator
				appendStrToBuf(b, Separator)
				b.WriteString(e.Err.Error())
			}
		} else {
			appendStrToBuf(b, ": ")
			b.WriteString(e.Err.Error())
		}
	}
	if b.Len() == 0 {
		return "no error"
	}
	return b.String()
}

// errors.New
func Str(text string) error {
	return errorString{text}
}

type errorString struct {
	s string
}

func (e errorString) Error() string {
	return e.s
}

func Errorf(format string, args ...interface{}) error {
	return errorString{fmt.Sprintf(format, args...)}
}

func Match(err1, err2 error) bool {
	e1, ok := err1.(*Error)
	if !ok {
		return false
	}
	e2, ok := err2.(*Error)
	if !ok {
		return false
	}
	if e1.Op != "" && e2.Op != e1.Op {
		return false
	}
	if e1.Kind != Undefined && e2.Kind != e1.Kind {
		return false
	}
	if e1.Err != nil {
		if _, ok := e1.Err.(*Error); ok {
			return Match(e1.Err, e2.Err)
		}
		if e2.Err == nil || e2.Err.Error() != e1.Err.Error() {
			return false
		}
	}
	return true
}

// Is reports whether err is an *Error of the given Kind.
func Is(kind Kind, err error) bool {
	e, ok := err.(*Error)
	if !ok {
		return false
	}
	if e.Kind != Undefined {
		return e.Kind == kind
	}
	if e.Err != nil {
		return Is(kind, e.Err)
	}
	return false
}

// Do smt with no care about result (and panics)
func SafelyDo(work func()) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("work failed: %s", err)
		}
	}()

	work()
}
