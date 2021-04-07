package main

import (
	"github.com/RobbyRed98/instructor/core"
	"github.com/RobbyRed98/instructor/printer"
	"os"
)

func main() {
	level := printer.INFO
	printy := printer.NewPrinter(&level)

	command := os.Args[1]

	var instructor core.Instructor
	instructor, err := core.NewBasicInstructor(printy)
	if err != nil {
		printy.Error("Fatal error occurred! Please use debug mode!")
		os.Exit(1)
	}

	switch command {
	case "list":
		instructor.List()

	case "add":
		instructor.Add()

	case "rm":
		instructor.Remove()

	case "mv", "rename":
		instructor.Rename()

	case "edit":
		instructor.Edit()

	case "copy":
		instructor.Copy()

	case "reorganize":
		instructor.Reorganize()

	case "Help":
		instructor.Help()

	default:
		instructor.Execute(command)
	}
}

