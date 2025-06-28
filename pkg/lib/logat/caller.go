package logat

import (
	"runtime"
	"strconv"
	"strings"
)

func getCaller() string {
	frame, defined := getCallerFrame(callerSkipOffset)
	if !defined {
		return "undefined"
	}
	return trimmedPath(frame.File, frame.Line)
}

func getCallerFrame(skip int) (frame runtime.Frame, ok bool) {
	const skipOffset = 2 // skip getCallerFrame and Callers

	pc := make([]uintptr, 1)
	numFrames := runtime.Callers(skip+skipOffset, pc)
	if numFrames < 1 {
		return
	}

	frame, _ = runtime.CallersFrames(pc).Next()
	return frame, frame.PC != 0
}

func trimmedPath(file string, line int) string {
	// Find the last separator
	idx := strings.LastIndexByte(file, '/')
	if idx == -1 {
		return fullPath(file, line)
	}

	// Find the penultimate separator
	idx = strings.LastIndexByte(file[:idx], '/')
	if idx == -1 {
		return fullPath(file, line)
	}

	// Keep everything after the penultimate separator
	var caller strings.Builder
	defer caller.Reset()

	caller.WriteString(file[idx+1:])
	caller.WriteByte(':')
	caller.WriteString(strconv.Itoa(line))

	return caller.String()
}

func fullPath(file string, line int) string {
	var caller strings.Builder
	defer caller.Reset()

	caller.WriteString(file)
	caller.WriteByte(':')
	caller.WriteString(strconv.Itoa(line))

	return caller.String()
}
