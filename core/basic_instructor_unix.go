// +build !windows

package core

import (
	"fmt"
	"github.com/RobbyRed98/instructor/runner"
	"os"
)

func (bi BasicInstructor) Execute(command string) error {
	argsNum, err := bi.checkMultiArgs(1, 2)
	if err != nil {
		return err
	}

	if argsNum == 2 && !(os.Args[1] == "-s" || os.Args[1] == "--shell") {
		return fmt.Errorf("The passed flag is unknown: " + os.Args[1])
	}

	label := command
	err = bi.checkLabel(label)
	if err != nil {
		return err
	}

	instruction, err := bi.instructionStorage.GetInstruction(bi.currentScope, label)
	if err != nil {
		bi.printy.Debug(instruction)
		return fmt.Errorf("Shortcut does not exist.")
	}

	instructionRunner := runner.NewRunner(bi.printy.Level())
	if argsNum == 1 {
		return instructionRunner.Run(instruction)
	} else if argsNum != 2 {
		return instructionRunner.ShellRun(instruction)
	}
	return fmt.Errorf("Unexpected situation: No appropriate runner found! That should not happen!")
}
