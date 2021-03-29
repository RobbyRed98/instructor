package runner

import (
	"github.com/RobbyRed98/instructor/printer"
	"github.com/mattn/go-shellwords"
	"os"
	"os/exec"
	"syscall"
)

type Runner struct {
	printer *printer.Printer
}

func NewRunner(level int) *Runner {
	newPrinter := printer.NewPrinter(&level)
	return &Runner{newPrinter}
}

func (r *Runner) Run(instruction string) {
	instructionFragments, err := shellwords.Parse(instruction)
	if err != nil {
		r.printer.Error("Failed to parse the instruction:", instruction)
		os.Exit(1)
	}

	command := instructionFragments[0]
	commandPath, err := exec.LookPath(command)
	r.printer.Debug("Running:", instruction)
	if err != nil {
		r.printer.Error("Failed to lookup command:", command)
		os.Exit(1)
	}
	syscall.Exec(commandPath, instructionFragments, os.Environ())
}
