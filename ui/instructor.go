package ui

import (
	"github.com/RobbyRed98/instructor/printer"
	"github.com/RobbyRed98/instructor/storage"
)

type Instructor interface {
	List(printy *printer.Printer, instructionStorage *storage.Storage, scope string)
	Add(printy *printer.Printer, instructionStorage *storage.Storage, scope string)
	Remove(printy *printer.Printer, instructionStorage *storage.Storage, scope string)
	Rename(printy *printer.Printer, instructionStorage *storage.Storage, scope string)
	Edit(printy *printer.Printer, instructionStorage *storage.Storage, scope string)
	Copy(printy *printer.Printer, instructionStorage *storage.Storage)
	Reorganize(printy *printer.Printer, instructionStorage *storage.Storage)
}
