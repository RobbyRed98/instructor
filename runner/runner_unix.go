// +build !windows

package runner

import (
	"github.com/mattn/go-shellwords"
	"os"
	"os/exec"
	"syscall"
)

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
		r.printer.Error("Failed to lookup command:", instruction)
		os.Exit(1)
	}
	syscall.Exec(commandPath, instructionFragments, os.Environ())
}
