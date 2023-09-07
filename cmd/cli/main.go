package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/lambdaclass/cairo-vm.go/pkg/vm/cairo_run"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Wrong argument count: Use go run cmd/cli/main.go COMPILED_JSON")
		return
	}
	cli_args := os.Args[1:]
	programPath := cli_args[0]
	cairoRunConfig := cairo_run.CairoRunConfig{DisableTracePadding: false}

	cairoRunner, err := cairo_run.CairoRun(programPath, cairoRunConfig)
	if err != nil {
		fmt.Printf("Failed with error: %s", err)
		return
	}
	traceFilePath := strings.Replace(programPath, ".json", ".go.trace", 1)
	traceFile, err := os.OpenFile(traceFilePath, os.O_RDWR|os.O_CREATE, 0644)
	defer traceFile.Close()

	memoryFilePath := strings.Replace(programPath, ".json", ".go.memory", 1)
	memoryFile, err := os.OpenFile(memoryFilePath, os.O_RDWR|os.O_CREATE, 0644)
	defer memoryFile.Close()

	cairo_run.WriteEncodedTrace(cairoRunner.Vm.RelocatedTrace, traceFile)
	cairo_run.WriteEncodedMemory(cairoRunner.Vm.RelocatedMemory, memoryFile)

	println("Done!")
}
