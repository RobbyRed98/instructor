package main

import (
	"github.com/RobbyRed98/instructor/core"
	"github.com/RobbyRed98/instructor/printer"
	"os"
)

func main() {
	if len(os.Args) <= 1 {
		core.BasicInstructor{}.Help()
		os.Exit(0)
	}

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
		err = instructor.List()

	case "add":
		err = instructor.Add()

	case "rm":
		err = instructor.Remove()

	case "mv", "rename":
		err = instructor.Rename()

	case "edit":
		err = instructor.Edit()

	case "reorganize":
		err = instructor.Reorganize()

	case "help":
		instructor.Help()
		err = nil

	default:
		err = instructor.Execute(command)
	}

	if err != nil {
		printy.Error(err.Error())
		os.Exit(1)
	}
}

