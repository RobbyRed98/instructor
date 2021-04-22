package parser

import (
	"fmt"
	"strings"
)

const ScopeLabelDelimiter = "|"
const LabelInstructionDelimiter = "->"

func Parse(entry string) (string, string, string, error) {
	parts := strings.SplitN(entry, ScopeLabelDelimiter, 2)
	if len(parts) < 2 {
		return "", "", "", fmt.Errorf("failed to parse entry")
	}
	scope := parts[0]
	otherParts := strings.SplitN(parts[1], LabelInstructionDelimiter, 2)
	if len(otherParts) < 2 {
		return "", "", "", fmt.Errorf("failed to parse entry")
	}
	label := otherParts[0]
	instruction := otherParts[1]
	if len(scope) == 0 {
		return "", "", "", fmt.Errorf("scope '%q' has length 0", scope)
	} else if len(label) == 0 {
		return "", "", "", fmt.Errorf("label '%q' has length 0", scope)
	} else if len(instruction) == 0 {
		return "", "", "", fmt.Errorf("instruction '%q' has length 0", scope)
	}
	return scope, label, instruction, nil
}

func GetScopeLabelPair(scope string, label string) string {
	return scope + ScopeLabelDelimiter + label
}

func GetEntry (scope string, label string, instruction string) string {
	return GetScopeLabelPair(scope, label) + LabelInstructionDelimiter + instruction
}


