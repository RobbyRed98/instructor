package main

import (
	"github.com/RobbyRed98/instructor/config"
	"github.com/RobbyRed98/instructor/printer"
	"github.com/RobbyRed98/instructor/runner"
	"github.com/RobbyRed98/instructor/storage"
	"os"
	"strconv"
	"strings"
)

func main() {
	printLevel := printer.INFO
	newPrinter := printer.NewPrinter(&printLevel)
	if len(os.Args) < 2 {
		newPrinter.Error("No command has been passed.")
		help(*newPrinter)
		os.Exit(1)
	}

	for i, arg := range os.Args {
		newPrinter.Debug(strconv.Itoa(i), ":", arg)
	}

	command := os.Args[1]
	scope, _ := os.Getwd()

	instructionsFilePath := config.GetInstructionFile()
	instructionStorage := storage.NewStorage(instructionsFilePath, newPrinter)

	switch command {
	case "list":
		argNum := checkMultiArgs(1,2, *newPrinter)
		if argNum == 2 && os.Args[2] == "all" {
			scope = ""
			newPrinter.Debug("Using global scope.")
			newPrinter.Debug("Listing all shortcuts.")
		}
		entries, err := instructionStorage.ListInstructions(scope)
		if err != nil {
			newPrinter.Info("No instructions file exists.")
		}

		for _, entry := range entries {
			entry = strings.Replace(entry, "|", " | ", 1)
			entry = strings.Replace(entry, "->", " -> ", 1)
			entry = strings.Trim(entry, "\n")
			newPrinter.Info(entry)
		}

	case "add":
		checkArgs(3, *newPrinter)
		label := os.Args[2]

		isLabel := instructionStorage.LabelExists(scope, label)
		if isLabel {
			newPrinter.Error("Shortcut already exists!")
			os.Exit(1)
		}

		instruction := os.Args[3]
		entry, err := instructionStorage.AddInstruction(scope, label, instruction)
		if err != nil {
			newPrinter.Error(err.Error())
			os.Exit(0)
		}
		newPrinter.Info("Created shortcut.")
		newPrinter.Debug(strings.Trim(entry, "\n"))

	case "rm":
		checkArgs(2, *newPrinter)
		label := os.Args[2]
		isLabel := instructionStorage.LabelExists(scope, label)
		if !isLabel {
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

	case "mv", "rename":
		checkArgs(3, *newPrinter)
		oldLabel := os.Args[2]
		newLabel := os.Args[3]
		isLabel := instructionStorage.LabelExists(scope, oldLabel)
		if !isLabel {
			newPrinter.Error("No shortcut found.")
			newPrinter.Debug(scope, "|", oldLabel)
			os.Exit(1)
		}

		err := instructionStorage.RenameInstruction(scope, oldLabel)
		if err != nil {
			newPrinter.Error("Failed to rename the shortcut.")
			newPrinter.Debug(scope + "|" + oldLabel , "->", scope + "|" + newLabel)
			os.Exit(1)
		}

	case "reorganize":
		checkArgs(1, *newPrinter)
		err := instructionStorage.Reorganize()
		if err != nil {
			newPrinter.Error("Failed to reorganize file.")
			os.Exit(1)
		}
		newPrinter.Debug("Successfully reorganized instructions file.")

	case "help":
		help(*newPrinter)

	default:
		label := command
		instruction, err := instructionStorage.GetInstruction(scope, label)

		if err != nil {
			newPrinter.Error(err.Error())
		}

		instructionRunner := runner.NewRunner(printLevel)
		instructionRunner.Run(instruction)
	}
}

func checkArgs(requiredNum int, newPrinter printer.Printer) {
	argsNum := len(os.Args) - 1
	if argsNum != requiredNum {
		newPrinter.Error("Wrong number of arguments.")
		newPrinter.Error("Arguments passed:", strconv.Itoa(argsNum) + ",", "arguments required:", strconv.Itoa(requiredNum))
		os.Exit(1)
	}
}

func checkMultiArgs(lowerNum int, upperNum int, newPrinter printer.Printer) int {
	argsNum := len(os.Args) - 1
	if lowerNum > argsNum || argsNum > upperNum {
		newPrinter.Error("Wrong number of arguments.")
		newPrinter.Error("Arguments passed:", strconv.Itoa(argsNum) + ",", "allow argument numbers:", strconv.Itoa(lowerNum), "-", strconv.Itoa(upperNum))
		os.Exit(1)
	}
	return argsNum
}

func help(newPrinter printer.Printer) {
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
		"rm              Removes a shortcut.",
		"list            Lists all existing shortcuts.",
		"reorganize      Reorganizes the file in which the shortcuts and commands are stored.",
		"",
		"help            Prints this help text.",
	}

	for _, line := range helpText {
		newPrinter.Info(line)
	}
}
