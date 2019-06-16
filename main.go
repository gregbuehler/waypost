package main

import (
	"runtime"

	"github.com/gregbuehler/waypost/cmd"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	cmd.Execute()
}
