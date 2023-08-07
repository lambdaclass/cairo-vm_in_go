package main

import (
	"fmt"

	"github.com/lambdaclass/cairo-vm.go/pkg/vm/cairo_run"
)

func main() {
	fibonacciPath := "cairo_programs/fibonacci.json"
	_, err := cairo_run.CairoRun(fibonacciPath)
	if err != nil {
		fmt.Printf("Failed with error: %s", err)
		return
	}
}
