package utils

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"strings"
)

func Recovery() {
	e := recover()
	if e == nil {
		return
	}

	err := fmt.Errorf("%v", e)
	panicLoc := identifyPanicLoc()
	fmt.Println(err, "\npanic location: ", panicLoc, " \nstacktrace:\n", string(debug.Stack()))
}

func SafeGo(fn func()) {
	go func() {
		defer Recovery()
		fn()
	}()
}

func identifyPanicLoc() string {
	var name, file string
	var line int
	pc := make([]uintptr, 16)

	_ = runtime.Callers(3, pc)
	frames := runtime.CallersFrames(pc)
	for i := 0; i < 16; i++ {
		if frames == nil {
			break
		}
		frame, hasMore := frames.Next()
		fn := runtime.FuncForPC(frame.PC)
		if fn == nil {
			break
		}
		file, line = frame.File, frame.Line
		name = fn.Name()
		if !strings.HasPrefix(name, "runtime.") {
			break
		}
		if !hasMore {
			break
		}
	}
	switch {
	case name != "":
		return fmt.Sprintf("%v:%v", name, line)
	case file != "":
		return fmt.Sprintf("%v:%v", file, line)
	}

	return fmt.Sprintf("pc:%x", pc)
}
