package main

import (
	"github.com/RobbyRed98/instructor/printer"
	"github.com/RobbyRed98/instructor/storage"
	"github.com/RobbyRed98/instructor/ui"
	"os"
	"path"
	"strconv"
)

func main() {
	printLevel := printer.DEBUG
	printy := printer.NewPrinter(&printLevel)
	if len(os.Args) < 2 {
		printy.Error("No command has been passed.")
		ui.Help(*printy)
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
		ui.List(printy, instructionStorage, scope)

	case "add":
		ui.Add(printy, instructionStorage, scope)

	case "rm":
		ui.Remove(printy, instructionStorage, scope)

	case "mv", "rename":
		ui.Rename(printy, instructionStorage, scope)

	case "edit":
		ui.Edit(printy, instructionStorage, scope)

	case "copy":
		ui.Copy(printy, instructionStorage)

	case "reorganize":
		ui.Reorganize(printy, instructionStorage)

	case "Help":
		ui.Help(*printy)

	default:
		ui.Execute(command, instructionStorage, scope, printy, printLevel)
	}
}


