package main

import (
	"github.com/RobbyRed98/instructor/printer"
	"github.com/RobbyRed98/instructor/runner"
	"github.com/RobbyRed98/instructor/storage"
	"os"
	"path"
	"strconv"
	"strings"
)

func main() {
	printLevel := printer.INFO
	printy := printer.NewPrinter(&printLevel)
	if len(os.Args) < 2 {
		printy.Error("No command has been passed.")
		Help(*printy)
		os.Exit(1)
	}

	for i, arg := range os.Args {
		printy.Debug(strconv.Itoa(i), ":", arg)
	}

	command := os.Args[1]
	scope, _ := os.Getwd()

	homeDir, _ := os.UserHomeDir()
	instructionsFilePath := path.Join(homeDir, ".instructions")
	instructionStorage := storage.NewStorage(instructionsFilePath, printy)

	switch command {
	case "list":
		List(printy, scope, instructionStorage)

	case "add":
		Add(printy, instructionStorage, scope)

	case "rm":
		Remove(printy, instructionStorage, scope)

	case "mv", "rename":
		Rename(printy, instructionStorage, scope)

	case "edit":
		Edit(printy, instructionStorage, scope)

	case "reorganize":
		Reorganize(printy, instructionStorage)

	case "Help":
		Help(*printy)

	default:
		Execute(command, instructionStorage, scope, printy, printLevel)
	}
}

func List(printy *printer.Printer, scope string, instructionStorage *storage.Storage) {
	argNum := checkMultiArgs(1, 2, *printy)
	if argNum == 2 && os.Args[2] == "all" {
		scope = ""
		printy.Debug("Using global scope.")
		printy.Debug("Listing all shortcuts.")
	} else if argNum == 2 {
		printy.Error("Invalid argument:", os.Args[2])
		os.Exit(0)
	}
	entries, err := instructionStorage.ListInstructions(scope)
	if err != nil {
		printy.Info("No instructions file exists.")
	}

	for _, entry := range entries {
		entry = "(" + entry
		entry = strings.Replace(entry, "|", " | ", 1)
		entry = strings.Replace(entry, "->", ") -> ", 1)
		entry = strings.Trim(entry, "\n")
		printy.Info(entry)
	}
}

func Add(printy *printer.Printer, instructionStorage *storage.Storage, scope string) {
	checkArgs(3, *printy)
	label := os.Args[2]
	checkLabel(label, *printy)

	hasInstructionFor := instructionStorage.HasInstructionFor(scope, label)
	if hasInstructionFor {
		printy.Error("Shortcut already exists!")
		os.Exit(1)
	}

	instruction := os.Args[3]
	checkInstruction(instruction, *printy)
	if strings.Contains(instruction, "\n") {
		printy.Error("Failed to create shortcut.")
		printy.Error("Instruction cannot contain linebreaks.")
		os.Exit(0)
	}
	
	entry, err := instructionStorage.AddInstruction(scope, label, instruction)
	if err != nil {
		printy.Error("Failed to create shortcut.")
		printy.Debug(err.Error())
		os.Exit(1)
	}
	printy.Info("Successfully created shortcut.")
	printy.Debug(strings.Trim(entry, "\n"))
}

func Remove(printy *printer.Printer, instructionStorage *storage.Storage, scope string) {
	checkArgs(2, *printy)
	label := os.Args[2]
	checkLabel(label, *printy)

	hasInstructionFor := instructionStorage.HasInstructionFor(scope, label)
	if !hasInstructionFor {
		printy.Error("Shortcut does not exist.")
		printy.Debug(scope + "|" + label)
		os.Exit(1)
	}

	err := instructionStorage.RemoveInstruction(scope, label)
	if err != nil {
		printy.Error("Failed to remove shortcut combination.")
		printy.Debug(scope, "|", label)
		printy.Debug(err.Error())
		os.Exit(1)
	}
	printy.Debug("Removed shortcut.")
	printy.Debug(scope, "|", label)
}

func Rename(printy *printer.Printer, instructionStorage *storage.Storage, scope string) {
	checkArgs(3, *printy)
	oldLabel := os.Args[2]
	checkLabel(oldLabel, *printy)
	newLabel := os.Args[3]
	checkLabel(newLabel, *printy)

	hasInstructionFor := instructionStorage.HasInstructionFor(scope, oldLabel)
	if !hasInstructionFor {
		printy.Error("No shortcut found.")
		printy.Debug(scope, "|", oldLabel)
		os.Exit(1)
	} else {
		printy.Debug("Shortcut found.")
	}

	err := instructionStorage.RenameInstruction(scope, oldLabel, newLabel)
	if err != nil {
		printy.Error("Failed to rename the shortcut.")
		printy.Debug(scope+"|"+oldLabel, "->", scope+"|"+newLabel)
		os.Exit(1)
	}
	printy.Debug("Successfully renamed shortcut:", oldLabel, "->", newLabel)
}

func Edit(printy *printer.Printer, instructionStorage *storage.Storage, scope string) {
	checkArgs(3, *printy)
	label := os.Args[2]
	checkLabel(label, *printy)

	instruction := os.Args[3]
	checkInstruction(instruction, *printy)
	hasInstructionFor := instructionStorage.HasInstructionFor(scope, label)
	if !hasInstructionFor {
		printy.Error("No shortcut found.")
		printy.Debug(scope, "|", label)
		os.Exit(1)
	} else {
		printy.Debug("Shortcut found.")
	}

	err := instructionStorage.Save()
	if err != nil {
		printy.Error("Failed to edit shortcut.")
		printy.Debug(err.Error())
		os.Exit(1)
	}

	err = instructionStorage.RemoveInstruction(scope, label)
	if err != nil {
		_ = instructionStorage.Rollback()
		printy.Error("Failed to edit shortcut.")
		printy.Debug(scope, "|", label)
		printy.Debug(err.Error())
		os.Exit(1)
	}

	entry, err := instructionStorage.AddInstruction(scope, label, instruction)
	if err != nil {
		_ = instructionStorage.Rollback()
		printy.Error("Failed to edit shortcut.")
		printy.Debug(err.Error())
		os.Exit(1)
	}

	err = instructionStorage.DeleteSave()
	if err != nil {
		printy.Debug(err.Error())
	}
	printy.Info("Successfully edited shortcut.")
	printy.Debug(strings.Trim(entry, "\n"))
}

func Reorganize(printy *printer.Printer, instructionStorage *storage.Storage) {
	checkArgs(1, *printy)
	err := instructionStorage.Reorganize()
	if err != nil {
		printy.Error("Failed to reorganize file.")
		os.Exit(1)
	}
	printy.Debug("Successfully reorganized instructions file.")
}

func Execute(command string, instructionStorage *storage.Storage, scope string, printy *printer.Printer, printLevel int) {
	label := command
	checkLabel(label, *printy)

	instruction, err := instructionStorage.GetInstruction(scope, label)
	if err != nil {
		printy.Error("Shortcut does not exist.")
		printy.Debug(instruction)

		os.Exit(1)
	}

	instructionRunner := runner.NewRunner(printLevel)
	instructionRunner.Run(instruction)
}

func Help(printy printer.Printer) {
	helpText := []string{
		"Usage:",
		"ins <command> <args>",
		"",
		"Allows the creation and usage of scope-bound shell shortcuts.",
		"",
		"<shortcut>      Executes a created shortcut.",
		"add             Creates a scope-bound shortcut for a shell command.",
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

func checkArgs(requiredNum int, printy printer.Printer) {
	argsNum := len(os.Args) - 1
	if argsNum != requiredNum {
		printy.Error("Wrong number of arguments.")
		printy.Error("Arguments passed:", strconv.Itoa(argsNum)+",", "arguments required:", strconv.Itoa(requiredNum))
		os.Exit(1)
	}
}

func checkMultiArgs(lowerNum int, upperNum int, printy printer.Printer) int {
	argsNum := len(os.Args) - 1
	if lowerNum > argsNum || argsNum > upperNum {
		printy.Error("Wrong number of arguments.")
		printy.Error("Arguments passed:", strconv.Itoa(argsNum)+",", "allow argument numbers:", strconv.Itoa(lowerNum), "-", strconv.Itoa(upperNum))
		os.Exit(1)
	}
	return argsNum
}

func checkLabel(label string, printy printer.Printer) {
	if strings.Contains(label, "\n") {
		printy.Error("Invalid argument.")
		printy.Error("Shortcut name cannot contain linebreaks.")
		os.Exit(0)
	}
}

func checkInstruction(instruction string, printy printer.Printer) {
	if strings.Contains(instruction, "\n") {
		printy.Error("Invalid argument.")
		printy.Error("Instruction cannot contain linebreaks.")
		os.Exit(0)
	}
}
