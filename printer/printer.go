package printer

import (
	"fmt"
	"strings"
)

type Printer struct {
	level *int
}

func NewPrinter(level *int) *Printer {
	return &Printer{level}
}

func (p Printer) Error(args ...string) {
	if *p.level < ERROR {
		return
	}
	fmt.Println(strings.Join(args, " "))
}

func (p Printer) Info(args ...string) {
	if *p.level < INFO {
		return
	}
	fmt.Println(strings.Join(args, " "))
}

func (p Printer) Debug(args ...string) {
	if *p.level < DEBUG {
		return
	}
	fmt.Println(strings.Join(args, " "))
}

func (p Printer) Level() int {
	return *p.level
}
