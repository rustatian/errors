// +build debug

package errors

import (
	"bytes"
	"fmt"
	"runtime"
	"strings"
)

type stack struct {
	callers []uintptr
	// TODO(adg): add time of creation
}

func (e *Error) populateStack() {
	e.callers = callers()

	e2, ok := e.Err.(*Error)
	if !ok {
		return
	}

	i := 0

	ok = false
	for ; i < len(e.callers) && i < len(e2.callers); i++ {
		if e.callers[len(e.callers)-1-i] != e.callers[len(e2.callers)-1-i] {
			break
		}
		ok = true

	}

	//ok = false
	//for ; i < len(e.callers) && i < len(e2.callers); i++ {
	//	if e.callers[len(e.callers)-1-i] != e2.callers[len(e2.callers)-1-i] {
	//		break
	//	}
	//	ok = true
	//}
	//if ok {
	//	head := e2.callers[:len(e2.callers)-i]
	//	tail := e.callers
	//	e.callers = make([]uintptr, len(head)+len(tail))
	//	copy(e.callers, head)
	//	copy(e.callers[len(head):], tail)
	//	e2.callers = nil
	//}
}

// frame returns the nth frame, with the frame at top of stack being 0.
func frame(callers []uintptr, n int) runtime.Frame {
	frames := runtime.CallersFrames(callers)
	var f runtime.Frame
	for i := len(callers) - 1; i >= n; i-- {
		var ok bool
		f, ok = frames.Next()
		if !ok {
			break
		}
	}
	return f
}

func (e *Error) printStack(b *bytes.Buffer) {
	printCallers := callers()

	var prev string
	var diff bool
	for i := 0; i < len(e.callers); i++ {
		thisFrame := frame(e.callers, i)

		name := thisFrame.Func.Name()

		if !diff && i < len(printCallers) {
			if name == frame(printCallers, i).Func.Name() {
				continue
			}
			diff = true
		}

		if name == prev {
			continue
		}

		trim := 0
		for {
			j := strings.IndexAny(name[trim:], "./")
			if j < 0 {
				break
			}
			if !strings.HasPrefix(prev, name[:j+trim]) {
				break
			}
			trim += j + 1 // skip over the separator
		}

		// Do the printing.
		appendStrToBuf(b, Separator)
		fmt.Fprintf(b, "%v:%d: ", thisFrame.File, thisFrame.Line)
		if trim > 0 {
			b.WriteString("...")
		}
		b.WriteString(name[trim:])

		prev = name
	}
}

func callers() []uintptr {
	var stk [64]uintptr
	const skip = 4
	n := runtime.Callers(skip, stk[:])
	return stk[:n]
}
