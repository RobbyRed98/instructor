package runner

import (
	"github.com/RobbyRed98/instructor/printer"
)

type Runner struct {
	printer *printer.Printer
}

func NewRunner(level int) *Runner {
	newPrinter := printer.NewPrinter(&level)
	return &Runner{newPrinter}
}
