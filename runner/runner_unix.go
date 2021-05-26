// +build !windows

package runner

import (
	"fmt"
	"github.com/mattn/go-shellwords"
	"os"
	"os/exec"
	"syscall"
)

func (r *Runner) Run(instruction string) error {
	instructionFragments, err := shellwords.Parse(instruction)
	if err != nil {
		return fmt.Errorf("Failed to parse the instruction: %s", instruction)
	}

	command := instructionFragments[0]
	commandPath, err := exec.LookPath(command)
	r.printer.Debug("Running:", instruction)
	if err != nil {
		return fmt.Errorf("Failed to lookup command: %s", instruction)
	}
	return syscall.Exec(commandPath, instructionFragments, os.Environ())
}

func (r *Runner) ShellRun(instruction string) error {
	if _, err := os.Stat("/bin/sh"); err == nil {
		return fmt.Errorf("Failed to find the shell: /bin/sh")
	}
	r.printer.Debug("Running:", "sh", "-c", "\"" + instruction + "\"")
	return syscall.Exec("/bin/sh", []string{"sh", "-c" , instruction}, os.Environ())
}
