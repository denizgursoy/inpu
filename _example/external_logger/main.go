package main

import (
	"github.com/denizgursoy/inpu"
	"github.com/denizgursoy/inpu/loggers/zero"
)

func main() {
	inpu.DefaultLogger = zero.NewInpuLoggerFromZeroLog()
}
