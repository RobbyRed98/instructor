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
	newPrinter := printer.NewPrinter(&printLevel)
	if len(os.Args) < 2 {
		newPrinter.Error("No command has been passed.")
		Help(*newPrinter)
		os.Exit(1)
	}

	for i, arg := range os.Args {
		newPrinter.Debug(strconv.Itoa(i), ":", arg)
	}

	command := os.Args[1]
	scope, _ := os.Getwd()

	homeDir, _ := os.UserHomeDir()
	instructionsFilePath := path.Join(homeDir, ".instructions")
	instructionStorage := storage.NewStorage(instructionsFilePath, newPrinter)

	switch command {
	case "list":
		List(newPrinter, scope, instructionStorage)

	case "add":
		Add(newPrinter, instructionStorage, scope)

	case "rm":
		Remove(newPrinter, instructionStorage, scope)

	case "mv", "rename":
		Rename(newPrinter, instructionStorage, scope)

	case "edit":
		Edit(newPrinter, instructionStorage, scope)

	case "reorganize":
		Reorganize(newPrinter, instructionStorage)

	case "Help":
		Help(*newPrinter)

	default:
		Execute(command, instructionStorage, scope, newPrinter, printLevel)
	}
}

func List(newPrinter *printer.Printer, scope string, instructionStorage *storage.Storage) {
	argNum := checkMultiArgs(1, 2, *newPrinter)
	if argNum == 2 && os.Args[2] == "all" {
		scope = ""
		newPrinter.Debug("Using global scope.")
		newPrinter.Debug("Listing all shortcuts.")
	} else if argNum == 2 {
		newPrinter.Error("Invalid argument:", os.Args[2])
		os.Exit(0)
	}
	entries, err := instructionStorage.ListInstructions(scope)
	if err != nil {
		newPrinter.Info("No instructions file exists.")
	}

	for _, entry := range entries {
		entry = "(" + entry
		entry = strings.Replace(entry, "|", " | ", 1)
		entry = strings.Replace(entry, "->", ") -> ", 1)
		entry = strings.Trim(entry, "\n")
		newPrinter.Info(entry)
	}
}

func Add(newPrinter *printer.Printer, instructionStorage *storage.Storage, scope string) {
	checkArgs(3, *newPrinter)
	label := os.Args[2]

	isInstruction := instructionStorage.InstructionExists(scope, label)
	if isInstruction {
		newPrinter.Error("Shortcut already exists!")
		os.Exit(1)
	}

	instruction := os.Args[3]
	entry, err := instructionStorage.AddInstruction(scope, label, instruction)
	if err != nil {
		newPrinter.Error("Failed to create shortcut.")
		newPrinter.Debug(err.Error())
		os.Exit(1)
	}
	newPrinter.Info("Successfully created shortcut.")
	newPrinter.Debug(strings.Trim(entry, "\n"))
}

func Remove(newPrinter *printer.Printer, instructionStorage *storage.Storage, scope string) {
	checkArgs(2, *newPrinter)
	label := os.Args[2]
	isInstruction := instructionStorage.InstructionExists(scope, label)
	if !isInstruction {
		newPrinter.Error("Shortcut does not exist.")
		newPrinter.Debug(scope + "|" + label)
		os.Exit(1)
	}

	err := instructionStorage.RemoveInstruction(scope, label)
	if err != nil {
		newPrinter.Error("Failed to remove shortcut combination.")
		newPrinter.Debug(scope, "|", label)
		newPrinter.Debug(err.Error())
		os.Exit(1)
	}
	newPrinter.Debug("Removed shortcut.")
	newPrinter.Debug(scope, "|", label)
}

func Rename(newPrinter *printer.Printer, instructionStorage *storage.Storage, scope string) {
	checkArgs(3, *newPrinter)
	oldLabel := os.Args[2]
	newLabel := os.Args[3]
	isInstruction := instructionStorage.InstructionExists(scope, oldLabel)
	if !isInstruction {
		newPrinter.Error("No shortcut found.")
		newPrinter.Debug(scope, "|", oldLabel)
		os.Exit(1)
	} else {
		newPrinter.Debug("Shortcut found.")
	}

	err := instructionStorage.RenameInstruction(scope, oldLabel, newLabel)
	if err != nil {
		newPrinter.Error("Failed to rename the shortcut.")
		newPrinter.Debug(scope+"|"+oldLabel, "->", scope+"|"+newLabel)
		os.Exit(1)
	}
	newPrinter.Debug("Successfully renamed shortcut:", oldLabel, "->", newLabel)
}

func Edit(newPrinter *printer.Printer, instructionStorage *storage.Storage, scope string) {
	checkArgs(3, *newPrinter)
	label := os.Args[2]
	instruction := os.Args[3]
	isInstruction := instructionStorage.InstructionExists(scope, label)
	if !isInstruction {
		newPrinter.Error("No shortcut found.")
		newPrinter.Debug(scope, "|", label)
		os.Exit(1)
	} else {
		newPrinter.Debug("Shortcut found.")
	}

	err := instructionStorage.Save()
	if err != nil {
		newPrinter.Error("Failed to edit shortcut.")
		newPrinter.Debug(err.Error())
		os.Exit(1)
	}

	err = instructionStorage.RemoveInstruction(scope, label)
	if err != nil {
		_ = instructionStorage.Rollback()
		newPrinter.Error("Failed to edit shortcut.")
		newPrinter.Debug(scope, "|", label)
		newPrinter.Debug(err.Error())
		os.Exit(1)
	}

	entry, err := instructionStorage.AddInstruction(scope, label, instruction)
	if err != nil {
		_ = instructionStorage.Rollback()
		newPrinter.Error("Failed to edit shortcut.")
		newPrinter.Debug(err.Error())
		os.Exit(1)
	}

	err = instructionStorage.DeleteSave()
	if err != nil {
		newPrinter.Debug(err.Error())
	}
	newPrinter.Info("Successfully edited shortcut.")
	newPrinter.Debug(strings.Trim(entry, "\n"))
}

func Reorganize(newPrinter *printer.Printer, instructionStorage *storage.Storage) {
	checkArgs(1, *newPrinter)
	err := instructionStorage.Reorganize()
	if err != nil {
		newPrinter.Error("Failed to reorganize file.")
		os.Exit(1)
	}
	newPrinter.Debug("Successfully reorganized instructions file.")
}

func Execute(command string, instructionStorage *storage.Storage, scope string, newPrinter *printer.Printer, printLevel int) {
	label := command
	instruction, err := instructionStorage.GetInstruction(scope, label)

	if err != nil {
		newPrinter.Error("Shortcut does not exist.")
		newPrinter.Debug(instruction)

		os.Exit(1)
	}

	instructionRunner := runner.NewRunner(printLevel)
	instructionRunner.Run(instruction)
}

func Help(newPrinter printer.Printer) {
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
		newPrinter.Info(line)
	}
}

func checkArgs(requiredNum int, newPrinter printer.Printer) {
	argsNum := len(os.Args) - 1
	if argsNum != requiredNum {
		newPrinter.Error("Wrong number of arguments.")
		newPrinter.Error("Arguments passed:", strconv.Itoa(argsNum)+",", "arguments required:", strconv.Itoa(requiredNum))
		os.Exit(1)
	}
}

func checkMultiArgs(lowerNum int, upperNum int, newPrinter printer.Printer) int {
	argsNum := len(os.Args) - 1
	if lowerNum > argsNum || argsNum > upperNum {
		newPrinter.Error("Wrong number of arguments.")
		newPrinter.Error("Arguments passed:", strconv.Itoa(argsNum)+",", "allow argument numbers:", strconv.Itoa(lowerNum), "-", strconv.Itoa(upperNum))
		os.Exit(1)
	}
	return argsNum
}
