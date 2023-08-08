package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/lambdaclass/cairo-vm.go/pkg/vm/cairo_run"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Wrong argument count: Use go run cmd/cli/main.go COMPILED_JSON")
	}
	cli_args := os.Args[1:]
	programPath := cli_args[0]
	cairoRunner, err := cairo_run.CairoRun(programPath)
	if err != nil {
		fmt.Printf("Failed with error: %s", err)
		return
	}
	traceFilePath := strings.Replace(programPath, ".json", ".go.trace", 1)
	traceFile, err := os.OpenFile(traceFilePath, os.O_RDWR|os.O_CREATE, 0644)

	// Dirty trick
	// TODO: Remove once WriteEncodedmemory is merged
	memoryFilePathRs := strings.Replace(programPath, ".json", ".rs.memory", 1)
	memoryFileRs, err := os.Open(memoryFilePathRs)
	if err != nil {
		fmt.Printf("Failed with error: %s", err)
		return
	}
	memoryFilePathGo := strings.Replace(programPath, ".json", ".go.memory", 1)
	memoryFileGo, err := os.Create(memoryFilePathGo)
	if err != nil {
		fmt.Printf("Failed with error: %s", err)
		return
	}
	io.Copy(memoryFileGo, memoryFileRs)

	if err != nil {
		fmt.Println(err)
	}
	defer traceFile.Close()

	cairo_run.WriteEncodedTrace(cairoRunner.Vm.RelocatedTrace, traceFile)

	println("Done!")
}
