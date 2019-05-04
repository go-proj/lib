package common

import (
	"os"
	"runtime"
)

func PrintStack() {
	var buf [4096]byte
	n := runtime.Stack(buf[:], false)
	os.Stdout.Write(buf[:n])
}
