package core

import (
	"errors"
	"fmt"
	"github.com/RobbyRed98/instructor/parser"
	"github.com/RobbyRed98/instructor/printer"
	"github.com/RobbyRed98/instructor/runner"
	"github.com/RobbyRed98/instructor/storage"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

type BasicInstructor struct {
	printy             printer.Printer
	instructionStorage storage.Storage
	currentScope       string
}

func NewBasicInstructor(printy *printer.Printer) (*BasicInstructor, error) {
	for i, arg := range os.Args {
		printy.Debug(strconv.Itoa(i) + ":", arg)
	}

	scope, err := os.Getwd()
	if err != nil {
		printy.Debug("Failed to get the current working directory.")
		return nil, err
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		printy.Debug("Failed to get the home directory.")
		return nil, err
	}

	instructionsFilePath := path.Join(homeDir, ".instructions")
	instructionStorage := storage.NewStorage(instructionsFilePath, printy)

	basicInstructor := BasicInstructor{*printy, *instructionStorage, scope}
	return &basicInstructor, nil
}

func (bi BasicInstructor) List() error {
	argNum, err := bi.checkMultiArgs(1, 2)
	if err != nil {
		return err
	} else if argNum == 2 && os.Args[2] == "all" {
		bi.currentScope = ""
		bi.printy.Debug("Using global scope.")
		bi.printy.Debug("Listing all shortcuts.")
	} else if argNum == 2 {
		return fmt.Errorf("Invalid argument: %s", os.Args[2])
	}

	entries, err := bi.instructionStorage.ListInstructions(bi.currentScope, true)
	if err != nil {
		bi.printy.Info("No instructions file exists.")
		return nil
	}

	for _, entry := range entries {
		entry = "(" + entry
		entry = strings.Replace(entry, "|", " | ", 1)
		entry = strings.Replace(entry, "->", ") -> ", 1)
		entry = strings.Trim(entry, "\n")
		bi.printy.Info(entry)
	}

	return nil
}

func (bi BasicInstructor) Add() error {
	err := bi.checkArgs(3)
	if err != nil {
		return err
	}

	label := os.Args[2]
	err = bi.checkLabel(label)
	if err != nil {
		return err
	}

	hasInstructionFor := bi.instructionStorage.HasInstructionFor(bi.currentScope, label)
	if hasInstructionFor {
		return fmt.Errorf("Shortcut already exists.")
	}

	instruction := os.Args[3]
	err = bi.checkInstruction(instruction)
	if err != nil {
		return err
	}

	if strings.Contains(instruction, "\n") {
		return fmt.Errorf("Failed to create shortcut.\nInstruction cannot contain linebreaks.")
	}

	entry, err := bi.instructionStorage.AddInstruction(bi.currentScope, label, instruction)
	if err != nil {
		bi.printy.Debug(err.Error())
		return fmt.Errorf("Failed to create shortcut.")
	}
	bi.printy.Info("Successfully created shortcut.")
	bi.printy.Debug(strings.Trim(entry, "\n"))
	return nil
}

func (bi BasicInstructor) Remove() error {
	err := bi.checkArgs(2)
	if err != nil {
		return err
	}

	label := os.Args[2]
	err = bi.checkLabel(label)
	if err != nil {
		return err
	}

	hasInstructionFor := bi.instructionStorage.HasInstructionFor(bi.currentScope, label)
	if !hasInstructionFor {
		bi.printy.Debug(bi.currentScope, "|", label)
		return fmt.Errorf("Shortcut does not exist.")
	}

	err = bi.instructionStorage.RemoveInstruction(bi.currentScope, label)
	if err != nil {
		bi.printy.Debug(err.Error())
		bi.printy.Debug(bi.currentScope, "|", label)
		return fmt.Errorf("Failed to remove shortcut combination.")
	}

	bi.printy.Debug("Removed shortcut.")
	bi.printy.Debug(bi.currentScope, "|", label)
	return nil
}

func (bi BasicInstructor) Rename() error {
	err := bi.checkArgs(3)
	if err != nil {
		return err
	}

	oldLabel := os.Args[2]
	err = bi.checkLabel(oldLabel)
	if err != nil {
		return err
	}

	newLabel := os.Args[3]
	err = bi.checkLabel(newLabel)
	if err != nil {
		return err
	}

	hasInstructionFor := bi.instructionStorage.HasInstructionFor(bi.currentScope, oldLabel)
	if !hasInstructionFor {
		bi.printy.Debug(bi.currentScope, "|", oldLabel)
		return fmt.Errorf("No shortcut found.")
	} else {
		bi.printy.Debug("Shortcut found.")
	}

	err = bi.instructionStorage.RenameInstruction(bi.currentScope, oldLabel, newLabel)
	if err != nil {
		bi.printy.Debug(bi.currentScope+"|"+oldLabel, "->", bi.currentScope+"|"+newLabel)
		return fmt.Errorf("Failed to rename the shortcut.")
	}
	bi.printy.Debug("Successfully renamed shortcut:", oldLabel, "->", newLabel)
	return nil
}

func (bi BasicInstructor) Edit() error {
	err := bi.checkArgs(3)
	if err != nil {
		return err
	}

	label := os.Args[2]
	err = bi.checkLabel(label)
	if err != nil {
		return err
	}

	instruction := os.Args[3]
	err = bi.checkInstruction(instruction)
	if err != nil {
		return err
	}
	hasInstructionFor := bi.instructionStorage.HasInstructionFor(bi.currentScope, label)
	if !hasInstructionFor {
		bi.printy.Debug(bi.currentScope, "|", label)
		return fmt.Errorf("No shortcut found.")
	} else {
		bi.printy.Debug("Shortcut found.")
	}

	err = bi.instructionStorage.Save()
	if err != nil {
		bi.printy.Debug(err.Error())
		return fmt.Errorf("Failed to edit shortcut.")
	}

	err = bi.instructionStorage.RemoveInstruction(bi.currentScope, label)
	if err != nil {
		_ = bi.instructionStorage.Rollback()
		bi.printy.Debug(err.Error())
		bi.printy.Debug(bi.currentScope, "|", label)
		return fmt.Errorf("Failed to edit shortcut.")
	}

	entry, err := bi.instructionStorage.AddInstruction(bi.currentScope, label, instruction)
	if err != nil {
		_ = bi.instructionStorage.Rollback()
		bi.printy.Debug(err.Error())
		return fmt.Errorf("Failed to edit shortcut.")
	}

	err = bi.instructionStorage.DeleteSave()
	if err != nil {
		bi.printy.Debug(err.Error())
	}
	bi.printy.Info("Successfully edited shortcut.")
	bi.printy.Debug(strings.Trim(entry, "\n"))
	return nil
}

func (bi BasicInstructor) Reorganize() error {
	err := bi.checkArgs(1)
	if err != nil {
		return err
	}

	err = bi.instructionStorage.Reorganize()
	if err != nil {
		return fmt.Errorf("Failed to reorganize file.")
	}

	bi.printy.Debug("Successfully reorganized instructions file.")
	return nil
}

func (bi BasicInstructor) Copy() error {
	_, err := bi.checkMultiArgs(2, 3)
	if err != nil {
		return err
	}

	srcScope := os.Args[2]
	var destScope string
	srcScope, err = bi.getAbsDirLikePath(srcScope)
	if err != nil {
		return fmt.Errorf("Failed to resolve source scope.")
	}

	if srcScope == destScope {
		return fmt.Errorf("Source and destination scope are the same.")
	}

	if len(os.Args) == 3 {
		destScope, err = bi.getAbsDirPath(os.Args[3])
		if err != nil {
			bi.printy.Debug(err.Error())
			return fmt.Errorf("Destination scope '%s' does not exist.", os.Args[3])
		}
	} else {
		bi.printy.Debug("Assuming destination directory is the current working directory.")
		destScope, err = os.Getwd()
		if err != nil {
			return fmt.Errorf("Could not locate current working directory!")
		}
	}

	newEntries, err := bi.instructionStorage.AlterInstructionForNewEntries(srcScope, destScope)
	if err != nil {
		bi.printy.Debug(err.Error())
		return fmt.Errorf("Failed to copy the instructions from the old scope.")
	}

	err = bi.instructionStorage.Save()
	if err != nil {
		bi.printy.Debug(err.Error())
		return fmt.Errorf("Unexpected situation. Failed to save.")
	}

	for _, entry := range newEntries {
		scope, label, instruction, err := parser.Parse(entry)
		if err != nil {
			bi.printy.Debug(err.Error())
			bi.printy.Error("Failed to parse the new entry:", entry)
			bi.printy.Info("Skipping the entry!")
			continue
		}

		if bi.instructionStorage.HasInstructionFor(scope, label) {
			bi.printy.Info(fmt.Sprintf("There is already an shortcut %s in the destination scope %s", label, scope))
			bi.printy.Info("Skipping the entry!")
			continue
		}

		_, err = bi.instructionStorage.AddInstruction(scope, label, instruction)
		if err != nil {
			bi.printy.Debug(err.Error())
			bi.printy.Error(fmt.Sprintf("Failed to add the shortcut (%s|%s)->%s.", scope, label, instruction))
			bi.printy.Info("Skipping the entry!")
			continue
		}
	}

	bi.printy.Info("Successfully copied the shortcuts.")
	return nil
}

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

func (bi BasicInstructor) Help() {
	helpText := []string{
		"Usage:",
		"ins <command> <args>",
		"",
		"Allows the creation and usage of scope-bound shortcuts for shell instructions.",
		"",
		"<shortcut>      Executes a created shortcut.",
		"add             Creates a shortcut for shell commands which is bound to the current scope.",
		"mv              Renames a shortcut.",
		"rename          Also renames a shortcut.",
		"edit 			 Edits the instruction of the shortcut by a replacing it with a new one.",
		"rm              Removes a shortcut.",
		"list            Lists all existing shortcuts.",
		"reorganize      Reorganizes the file in which the shortcuts and commands are stored.",
		"",
		"help            Prints this Help text.",
	}

	for _, line := range helpText {
		println(line)
	}
}

func (bi BasicInstructor) checkArgs(requiredNum int) error {
	argsNum := len(os.Args) - 1
	if argsNum != requiredNum {
		return fmt.Errorf("Wrong number of arguments.\nArguments passed: %d arguments required: %d",
			argsNum, requiredNum)
	}
	return nil
}

func (bi BasicInstructor) checkMultiArgs(lowerNum int, upperNum int) (int, error) {
	argsNum := len(os.Args) - 1
	if lowerNum > argsNum || argsNum > upperNum {
		return -1, fmt.Errorf("Wrong number of arguments.\nArguments passed: %d, allow argument numbers: %d", lowerNum, upperNum)
	}
	return argsNum, nil
}

func (bi BasicInstructor) checkLabel(label string) error {
	if strings.Contains(label, "\n") {
		return fmt.Errorf("Invalid argument.\nShortcut name cannot contain linebreaks.")
	}

	if strings.Contains(label, parser.LabelInstructionDelimiter) {
		return fmt.Errorf("Invalid argument.\nShortcut name cannot contain '%s'",
			parser.LabelInstructionDelimiter)
	}
	return nil
}

func (bi BasicInstructor) checkInstruction(instruction string) error {
	if strings.Contains(instruction, "\n") {
		return fmt.Errorf("Invalid argument.\nInstruction cannot contain linebreaks.")
	}
	return nil
}

func (bi BasicInstructor) getAbsDirPath(path string) (string, error)  {
	if strings.HasPrefix(path, "~") {
		homeDir, _ := os.UserHomeDir()
		path = strings.Replace(path, "~", homeDir, 1)
	}
	stat, err := os.Stat(path)
	exists := !errors.Is(err, os.ErrNotExist)

	if !exists {
		return "", fmt.Errorf("path does not exist: %s", path)
	}

	if !stat.IsDir() {
		return "", fmt.Errorf("path is not a directory: %s", path)
	}

	destScope, err := filepath.Abs(path)
	return destScope, nil
}

func (bi BasicInstructor) getAbsDirLikePath(pseudoPath string) (string, error)  {
	if strings.HasPrefix(pseudoPath, "/") {
		return path.Clean(pseudoPath), nil
	}

	if strings.HasPrefix(pseudoPath, "~") {
		homeDir, _ := os.UserHomeDir()
		pseudoPath = strings.Replace(pseudoPath, "~", homeDir, 1)
		return path.Clean(pseudoPath), nil
	}

	cwd, err := os.Getwd()

	if err != nil {
		return "", fmt.Errorf("failed to get working directory")
	}

	destPath := path.Join(cwd, pseudoPath)

	return path.Clean(destPath), nil
}
