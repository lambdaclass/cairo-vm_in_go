package main

import (
	"log"
	"os"
	"strings"

	"github.com/lambdaclass/cairo-vm.go/pkg/vm/cairo_run"
	"github.com/urfave/cli/v2"
)

func handleCommands(ctx *cli.Context) error {
	programPath := ctx.Args().First()

	layout := ctx.String("layout")
	if layout == "" {
		layout = "plain"
	}

	proofMode := ctx.Bool("proof_mode")

	secureRun := !proofMode
	if ctx.Bool("secure_run") {
		secureRun = true
	}

	cairoRunConfig := cairo_run.CairoRunConfig{DisableTracePadding: false, ProofMode: proofMode, Layout: layout, SecureRun: secureRun}

	cairoRunner, err := cairo_run.CairoRun(programPath, cairoRunConfig)
	if err != nil {
		return err
	}

	traceFilePath := ctx.String("trace_file")
	if traceFilePath == "" {
		traceFilePath = strings.Replace(programPath, ".json", ".go.trace", 1)
	}
	traceFile, err := os.OpenFile(traceFilePath, os.O_RDWR|os.O_CREATE, 0644)
	defer traceFile.Close()

	memoryFilePath := ctx.String("memory_file")
	if memoryFilePath == "" {
		memoryFilePath = strings.Replace(programPath, ".json", ".go.memory", 1)
	}
	memoryFile, err := os.OpenFile(memoryFilePath, os.O_RDWR|os.O_CREATE, 0644)
	defer memoryFile.Close()

	cairo_run.WriteEncodedTrace(cairoRunner.Vm.RelocatedTrace, traceFile)
	cairo_run.WriteEncodedMemory(cairoRunner.Vm.RelocatedMemory, memoryFile)
	return nil
}

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "proof_mode",
				Aliases: []string{"p"},
				Usage:   "Run in proof mode",
			},
			&cli.BoolFlag{
				Name:    "secure_run",
				Aliases: []string{"p"},
				Usage:   "Run security checks. Default: true unless proof_mode is true",
			},
			&cli.StringFlag{
				Name:    "layout",
				Aliases: []string{"l"},
				Usage:   "Default: plain",
			},
			&cli.StringFlag{
				Name:    "trace_file",
				Aliases: []string{"t"},
				Usage:   "--trace_file <TRACE_FILE>",
			},
			&cli.StringFlag{
				Name:    "memory_file",
				Aliases: []string{"m"},
				Usage:   "--memory_file <MEMORY_FILE>",
			},
		},
		Action: handleCommands,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
