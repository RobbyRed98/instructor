// +build windows

package core

import (
	"fmt"
	"github.com/RobbyRed98/instructor/runner"
)

func (bi BasicInstructor) Execute(command string) error {
	label := command
	err := bi.checkLabel(label)
	if err != nil {
		return err
	}

	instruction, err := bi.instructionStorage.GetInstruction(bi.currentScope, label)
	if err != nil {
		bi.printy.Debug(instruction)
		return fmt.Errorf("Shortcut does not exist.")
	}

	instructionRunner := runner.NewRunner(bi.printy.Level())
	instructionRunner.Run(instruction)
	return nil
}
