package storage

import (
	"bufio"
	"fmt"
	"github.com/RobbyRed98/instructor/parser"
	"github.com/RobbyRed98/instructor/printer"
	"github.com/mattn/go-shellwords"
	"os"
	"sort"
	"strings"
)

type Storage struct {
	instructionFilePath       string
	instructionTmpFilePath    string
	instructionSaveFilePath   string
	printy                    *printer.Printer
}

func NewStorage(path string, printer *printer.Printer) *Storage {
	s := Storage{}
	s.instructionFilePath = path
	s.instructionTmpFilePath = path + ".tmp"
	s.instructionSaveFilePath = path + ".bak"
	s.printy = printer
	return &s
}

func (s Storage) Exists() bool {
	_, err := os.Stat(s.instructionFilePath)
	return err == nil
}

func (s Storage) Save() error {
	fileContent, err := os.ReadFile(s.instructionFilePath)
	if err != nil {
		return fmt.Errorf("failed to read instruction file")
	}

	err = os.WriteFile(s.instructionSaveFilePath, fileContent, 0644)
	if err != nil {
		_ = os.Remove(s.instructionSaveFilePath)
		return fmt.Errorf("failed to write instruction save file")
	}
	return nil
}

func (s Storage) DeleteSave() error {
	err := os.Remove(s.instructionSaveFilePath)
	if err != nil {
		return fmt.Errorf("failed to delete instruction save file")
	}
	return nil
}

func (s Storage) Rollback() error {
	err := os.Rename(s.instructionSaveFilePath, s.instructionFilePath)
	if err != nil {
		return fmt.Errorf("failed to replace instructions file by backup")
	}
	return nil
}

func (s Storage) Reorganize() error {
	file, tmpFile, err := s.openInstructionFiles()
	if err != nil {
		return fmt.Errorf("failed to open instructions file or tmp file")
	}

	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		line := scanner.Text() + "\n"
		lines = append(lines, line)
	}
	sort.Strings(lines)

	writer := bufio.NewWriter(tmpFile)
	for _, line := range lines {
		_, err := writer.WriteString(line)
		if err != nil {
			_ = file.Close()
			_ = tmpFile.Close()
			return fmt.Errorf("failed to write to tmp file")
		}
	}
	_ = file.Close()

	err = writer.Flush()
	_ = tmpFile.Close()
	if err != nil {
		return fmt.Errorf("failed to flush lines to tmp file")
	}

	err = os.Rename(s.instructionTmpFilePath, s.instructionFilePath)
	if err != nil {
		_ = os.Remove(s.instructionTmpFilePath)
		return fmt.Errorf("failed to replace the instructions file by the tmp file")
	}

	return nil
}

func (s Storage) AddInstruction(scope string, label string, instruction string) (string, error) {
	file, err := os.OpenFile(s.instructionFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return "", fmt.Errorf("failed to open instructions file in append mode")
	}
	defer file.Close()

	_, err = shellwords.Parse(instruction)

	if err != nil {
		return "", fmt.Errorf("failed to open instructions file")
	}

	entry := parser.GetEntry(scope, label, instruction) + "\n"
	_, err = file.WriteString(entry)
	return entry, err
}

func (s Storage) ListInstructions(scope string, addLineEnds bool) ([]string, error) {
	file, err := os.Open(s.instructionFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open instruction file")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		line := scanner.Text()
		if addLineEnds {
			line = line + "\n"
		}
		if s.hasScope(line, scope) {
			lines = append(lines, line)
		}
	}
	sort.Strings(lines)
	return lines, nil
}

func (s Storage) RenameInstruction(scope string, oldLabel string, newLabel string) error {
	file, tmpFile, err := s.openInstructionFiles()
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(file)
	writer := bufio.NewWriter(tmpFile)

	for scanner.Scan() {
		line := scanner.Text()
		if s.hasScopeAndLabel(line, scope, oldLabel) {
			substrings := strings.SplitAfterN(line, parser.LabelInstructionDelimiter, 2)
			if len(substrings) != 2 {
				_ = file.Close()
				_ = tmpFile.Close()
				return fmt.Errorf("entry is corrupted '%s'", line)
			}
			instruction := substrings[1]
			line = parser.GetEntry(scope, newLabel, instruction)
		}
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			_ = file.Close()
			_ = tmpFile.Close()
			return fmt.Errorf("failed to write to tmp file")
		}
	}
	_ = file.Close()

	err = writer.Flush()
	_ = tmpFile.Close()
	if err != nil {
		return fmt.Errorf("failed to flush lines to tmp file")
	}

	err = os.Rename(s.instructionTmpFilePath, s.instructionFilePath)
	if err != nil {
		_ = os.Remove(s.instructionTmpFilePath)
		return fmt.Errorf("failed to replace the instructions file by the tmp file")
	}

	return nil
}

func (s Storage) RemoveInstruction(scope string, label string) error {
	file, tmpFile, err := s.openInstructionFiles()
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(file)
	writer := bufio.NewWriter(tmpFile)

	for scanner.Scan() {
		line := scanner.Text()
		if s.hasScopeAndLabel(line, scope, label) {
			continue
		}
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			_ = file.Close()
			_ = tmpFile.Close()
			return fmt.Errorf("failed to write line '%s' to tmp file", line)
		}
	}
	_ = file.Close()

	err = writer.Flush()
	_ = tmpFile.Close()
	if err != nil {
		return fmt.Errorf("failed to flush lines to tmp file")
	}

	err = os.Rename(s.instructionTmpFilePath, s.instructionFilePath)
	if err != nil {
		_ = os.Remove(s.instructionTmpFilePath)
		return fmt.Errorf("failed to replace the instructions file by the tmp file")
	}

	return nil
}

func (s Storage) GetInstruction(scope string, label string) (string, error) {
	file, err := os.Open(s.instructionFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to open instruction file")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if s.hasScopeAndLabel(line, scope, label) {
			substrings := strings.SplitAfterN(line, parser.LabelInstructionDelimiter, 2)
			if len(substrings) != 2 {
				return "", fmt.Errorf("invalid state scope-label '%s|%s' command lacks instruction", scope, label)
			}
			instruction := substrings[1]
			return instruction, nil
		}
	}
	return "", fmt.Errorf("scope-label combination '%s|%s' does not exists", scope, label)
}

func (s Storage) HasInstructionFor(scope string, label string) bool {
	file, _ := os.Open(s.instructionFilePath)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if s.hasScopeAndLabel(line, scope, label) {
			return true
		}
	}
	return false
}

func (s Storage) AlterInstructionForNewEntries(srcScope string, destScope string) ([]string, error) {
	instructions, err := s.ListInstructions(srcScope, false)
	if err != nil {
		s.printy.Debug(err.Error())
		return nil, fmt.Errorf("failed to get instructions")
	}

	if len(instructions) < 1 {
		return nil, fmt.Errorf("source scope '%s' has no instructions to copy", srcScope)
	}

	destinationInstructions := make([]string, len(instructions))

	for i, srcInstruction := range instructions {
		destInstruction := strings.Replace(srcInstruction, srcScope + parser.ScopeLabelDelimiter, destScope + parser.ScopeLabelDelimiter, 1)
		if !strings.HasPrefix(destInstruction, destScope) {
			return nil, fmt.Errorf("dest instruction does not begin with destination scope")
		}
		destinationInstructions[i] = destInstruction
	}

	return destinationInstructions, nil
}

func (s Storage) openInstructionFiles() (*os.File, *os.File, error) {
	file, err := os.Open(s.instructionFilePath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open instructions file %s", s.instructionFilePath)
	}

	tmpFile, err := os.OpenFile(s.instructionTmpFilePath, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		_ = file.Close()
		return nil, nil, fmt.Errorf("failed to open instructions tmp file %s", s.instructionTmpFilePath)
	}
	return file, tmpFile, nil
}

func (s Storage) hasScopeAndLabel(entry string, scope string, label string) bool {
	scopeLabelPrefix := parser.GetScopeLabelPair(scope, label) + parser.LabelInstructionDelimiter
	return strings.HasPrefix(entry, scopeLabelPrefix)
}

func (s Storage) hasScope(entry string, scope string) bool {
	if scope == "" {
		return true
	}
	return strings.HasPrefix(entry, scope+parser.ScopeLabelDelimiter)
}


