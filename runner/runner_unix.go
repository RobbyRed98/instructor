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

func (r *Runner) BashRun(instruction string) error {
	if _, err := os.Stat("/bin/bash"); err != nil {
		return fmt.Errorf("Failed to find the shell: /bin/bash")
	}
	r.printer.Debug("Running:", "bash", "-c", "\"" + instruction + "\"")
	return syscall.Exec("/bin/bash", []string{"bash", "-c" , instruction}, os.Environ())
}
