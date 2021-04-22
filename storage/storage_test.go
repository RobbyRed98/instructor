// +build !windows

package storage

import (
	"github.com/RobbyRed98/instructor/printer"
	"os"
	"path"
	"testing"
)

const TEST_INSTRUCTION_FILE = "test_instruction"

func root() string {
	return path.Join(os.TempDir(), TEST_INSTRUCTION_FILE)
}

func testee() *Storage {
	level := printer.NONE
	return NewStorage(root(), printer.NewPrinter(&level))
}

func TestStorage_Exists(t *testing.T) {
	testInstructionFile := path.Join(os.TempDir(), TEST_INSTRUCTION_FILE)
	err := os.Remove(testInstructionFile)
	if err != nil {
		t.Errorf("failed to clean up before running tests")
	}

	strg := testee()
	if strg.Exists() {
		t.Errorf("instruction file should not exist, but it does")
	}

	_, err = os.Create(testInstructionFile)
	if err != nil {
		t.Errorf("failed to create instruction file to test on")
	}

	if !strg.Exists() {
		t.Errorf("instruction file should exist, but it does not")
	}
}


