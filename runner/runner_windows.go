// +build windows

package runner

import (
	"os"
	"os/exec"
)

func (r *Runner) Run(instruction string) {
	cmd := exec.Command("powershell", instruction)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Env = os.Environ()
	r.printer.Debug("Running:", instruction)
	err := cmd.Run()
	if err != nil {
		r.printer.Debug("Run failed.")
		os.Exit(1)
	}
}
