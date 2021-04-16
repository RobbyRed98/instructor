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

func (bi BasicInstructor) List() {
	argNum := bi.checkMultiArgs(1, 2)
	if argNum == 2 && os.Args[2] == "all" {
		bi.currentScope = ""
		bi.printy.Debug("Using global currentScope.")
		bi.printy.Debug("Listing all shortcuts.")
	} else if argNum == 2 {
		bi.printy.Error("Invalid argument:", os.Args[2])
		os.Exit(0)
	}
	entries, err := bi.instructionStorage.ListInstructions(bi.currentScope, true)
	if err != nil {
		bi.printy.Info("No instructions file exists.")
	}

	for _, entry := range entries {
		entry = "(" + entry
		entry = strings.Replace(entry, "|", " | ", 1)
		entry = strings.Replace(entry, "->", ") -> ", 1)
		entry = strings.Trim(entry, "\n")
		bi.printy.Info(entry)
	}
}

func (bi BasicInstructor) Add() {
	bi.checkArgs(3)
	label := os.Args[2]
	bi.checkLabel(label)

	hasInstructionFor := bi.instructionStorage.HasInstructionFor(bi.currentScope, label)
	if hasInstructionFor {
		bi.printy.Error("Shortcut already exists!")
		os.Exit(1)
	}

	instruction := os.Args[3]
	bi.checkInstruction(instruction)
	if strings.Contains(instruction, "\n") {
		bi.printy.Error("Failed to create shortcut.")
		bi.printy.Error("Instruction cannot contain linebreaks.")
		os.Exit(0)
	}

	entry, err := bi.instructionStorage.AddInstruction(bi.currentScope, label, instruction)
	if err != nil {
		bi.printy.Error("Failed to create shortcut.")
		bi.printy.Debug(err.Error())
		os.Exit(1)
	}
	bi.printy.Info("Successfully created shortcut.")
	bi.printy.Debug(strings.Trim(entry, "\n"))
}

func (bi BasicInstructor) Remove() {
	bi.checkArgs(2)
	label := os.Args[2]
	bi.checkLabel(label)

	hasInstructionFor := bi.instructionStorage.HasInstructionFor(bi.currentScope, label)
	if !hasInstructionFor {
		bi.printy.Error("Shortcut does not exist.")
		bi.printy.Debug(bi.currentScope + "|" + label)
		os.Exit(1)
	}

	err := bi.instructionStorage.RemoveInstruction(bi.currentScope, label)
	if err != nil {
		bi.printy.Error("Failed to remove shortcut combination.")
		bi.printy.Debug(bi.currentScope, "|", label)
		bi.printy.Debug(err.Error())
		os.Exit(1)
	}
	bi.printy.Debug("Removed shortcut.")
	bi.printy.Debug(bi.currentScope, "|", label)
}

func (bi BasicInstructor) Rename() {
	bi.checkArgs(3)
	oldLabel := os.Args[2]
	bi.checkLabel(oldLabel)
	newLabel := os.Args[3]
	bi.checkLabel(newLabel)

	hasInstructionFor := bi.instructionStorage.HasInstructionFor(bi.currentScope, oldLabel)
	if !hasInstructionFor {
		bi.printy.Error("No shortcut found.")
		bi.printy.Debug(bi.currentScope, "|", oldLabel)
		os.Exit(1)
	} else {
		bi.printy.Debug("Shortcut found.")
	}

	err := bi.instructionStorage.RenameInstruction(bi.currentScope, oldLabel, newLabel)
	if err != nil {
		bi.printy.Error("Failed to rename the shortcut.")
		bi.printy.Debug(bi.currentScope+"|"+oldLabel, "->", bi.currentScope+"|"+newLabel)
		os.Exit(1)
	}
	bi.printy.Debug("Successfully renamed shortcut:", oldLabel, "->", newLabel)
}

func (bi BasicInstructor) Edit() {
	bi.checkArgs(3)
	label := os.Args[2]
	bi.checkLabel(label)

	instruction := os.Args[3]
	bi.checkInstruction(instruction)
	hasInstructionFor := bi.instructionStorage.HasInstructionFor(bi.currentScope, label)
	if !hasInstructionFor {
		bi.printy.Error("No shortcut found.")
		bi.printy.Debug(bi.currentScope, "|", label)
		os.Exit(1)
	} else {
		bi.printy.Debug("Shortcut found.")
	}

	err := bi.instructionStorage.Save()
	if err != nil {
		bi.printy.Error("Failed to edit shortcut.")
		bi.printy.Debug(err.Error())
		os.Exit(1)
	}

	err = bi.instructionStorage.RemoveInstruction(bi.currentScope, label)
	if err != nil {
		_ = bi.instructionStorage.Rollback()
		bi.printy.Error("Failed to edit shortcut.")
		bi.printy.Debug(bi.currentScope, "|", label)
		bi.printy.Debug(err.Error())
		os.Exit(1)
	}

	entry, err := bi.instructionStorage.AddInstruction(bi.currentScope, label, instruction)
	if err != nil {
		_ = bi.instructionStorage.Rollback()
		bi.printy.Error("Failed to edit shortcut.")
		bi.printy.Debug(err.Error())
		os.Exit(1)
	}

	err = bi.instructionStorage.DeleteSave()
	if err != nil {
		bi.printy.Debug(err.Error())
	}
	bi.printy.Info("Successfully edited shortcut.")
	bi.printy.Debug(strings.Trim(entry, "\n"))
}

func (bi BasicInstructor) Reorganize() {
	bi.checkArgs(1)
	err := bi.instructionStorage.Reorganize()
	if err != nil {
		bi.printy.Error("Failed to reorganize file.")
		os.Exit(1)
	}
	bi.printy.Debug("Successfully reorganized instructions file.")
}

func (bi BasicInstructor) Copy() {
	bi.checkMultiArgs(2, 3)
	srcScope := os.Args[2]
	var destScope string
	var err error

	srcScope, err = bi.getAbsDirLikePath(srcScope)

	if srcScope == destScope {
		bi.printy.Error("Source and destination currentScope are the same.")
		os.Exit(1)
	}

	if len(os.Args) == 4 {
		destScope, err = bi.getAbsDirPath(os.Args[3])
		if err != nil {
			bi.printy.Debug(err.Error())
			bi.printy.Error(fmt.Sprintf("Destination currentScope '%s' does not exist.", os.Args[4]))
			os.Exit(1)
		}
	} else {
		bi.printy.Debug("Assuming destination directory is the current working directory.")
		destScope, err = os.Getwd()
		if err != nil {
			bi.printy.Error("Could not locate current working directory!")
			os.Exit(1)
		}
	}

	newEntries, err := bi.instructionStorage.AlterInstructionForNewEntries(srcScope, destScope)
	if err != nil {
		bi.printy.Debug(err.Error())
		bi.printy.Error("Failed to copy the instructions from the old currentScope.")
		os.Exit(1)
	}

	err = bi.instructionStorage.Save()
	if err != nil {
		bi.printy.Debug(err.Error())
		bi.printy.Error("Unexpected situation. Failed to save ")
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
			bi.printy.Info(fmt.Sprintf("There is already an shortcut %s in the destination currentScope %s", label, scope))
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
}

func (bi BasicInstructor) Execute(command string) {
	label := command
	bi.checkLabel(label)

	instruction, err := bi.instructionStorage.GetInstruction(bi.currentScope, label)
	if err != nil {
		bi.printy.Error("Shortcut does not exist.")
		bi.printy.Debug(instruction)

		os.Exit(1)
	}

	instructionRunner := runner.NewRunner(bi.printy.Level())
	instructionRunner.Run(instruction)
}

func (bi BasicInstructor) Help() {
	help(&bi.printy)
}

func (bi BasicInstructor) checkArgs(requiredNum int) {
	argsNum := len(os.Args) - 1
	if argsNum != requiredNum {
		bi.printy.Error("Wrong number of arguments.")
		bi.printy.Error("Arguments passed:", strconv.Itoa(argsNum)+",", "arguments required:", strconv.Itoa(requiredNum))
		os.Exit(1)
	}
}

func (bi BasicInstructor) checkMultiArgs(lowerNum int, upperNum int) int {
	argsNum := len(os.Args) - 1
	if lowerNum > argsNum || argsNum > upperNum {
		bi.printy.Error("Wrong number of arguments.")
		bi.printy.Error("Arguments passed:", strconv.Itoa(argsNum)+",", "allow argument numbers:", strconv.Itoa(lowerNum), "-", strconv.Itoa(upperNum))
		os.Exit(1)
	}
	return argsNum
}

func (bi BasicInstructor) checkLabel(label string) {
	if strings.Contains(label, "\n") {
		bi.printy.Error("Invalid argument.")
		bi.printy.Error("Shortcut name cannot contain linebreaks.")
		os.Exit(0)
	}

	if strings.Contains(label, parser.LabelInstructionDelimiter) {
		bi.printy.Error("Invalid argument.")
		bi.printy.Error("Shortcut name cannot contain '" + parser.LabelInstructionDelimiter + "'.")
		os.Exit(0)
	}
}

func (bi BasicInstructor) checkInstruction(instruction string) {
	if strings.Contains(instruction, "\n") {
		bi.printy.Error("Invalid argument.")
		bi.printy.Error("Instruction cannot contain linebreaks.")
		os.Exit(0)
	}
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

	homeDir, err := os.Getwd()

	if err != nil {
		return "", fmt.Errorf("failed to get working directory")
	}

	destPath := path.Join(homeDir, pseudoPath)

	return path.Clean(destPath), nil
}

func help(printy *printer.Printer) {
	helpText := []string{
		"Usage:",
		"ins <command> <args>",
		"",
		"Allows the creation and usage of currentScope-bound shell shortcuts.",
		"",
		"<shortcut>      Executes a created shortcut.",
		"add             Creates a currentScope-bound shortcut for a shell command.",
		"mv              Renames a shortcut.",
		"rename          Also renames a shortcut.",
		"edit 			 Edits the instruction of the shortcut by a replacing it with a new one.",
		"rm              Removes a shortcut.",
		"list            Lists all existing shortcuts.",
		"reorganize      Reorganizes the file in which the shortcuts and commands are stored.",
		"",
		"Help            Prints this Help text.",
	}

	for _, line := range helpText {
		printy.Info(line)
	}
}
