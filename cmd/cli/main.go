package main

import (
	"bytes"
	"fmt"

	"github.com/lambdaclass/cairo-vm.go/pkg/vm/cairo_run"
)

func main() {
	fibonacciPath := "cairo_programs/fibonacci.json"
	cairoRunner, err := cairo_run.CairoRun(fibonacciPath)
	if err != nil {
		fmt.Printf("Failed with error: %s", err)
		return
	}

	var traceBuffer bytes.Buffer
	cairo_run.WriteEncodedTrace(cairoRunner.Vm.RelocatedTrace, &traceBuffer)

	println("Done!")
}
