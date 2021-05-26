// +build windows

package runner

import (
	"fmt"
	"os"
	"os/exec"
)

func (r *Runner) Run(instruction string) error {
	cmd := exec.Command("powershell", instruction)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Env = os.Environ()
	r.printer.Debug("Running:", instruction)
	err := cmd.Run()
	if err != nil {
		fmt.Errorf("Run failed: %s", instruction)
	}
	return nil
}
