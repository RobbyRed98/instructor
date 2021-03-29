package config

import (
	"os"
	"path"
)

func GetInstructionFile() string {
	homeDir, _ := os.UserHomeDir()
	return path.Join(homeDir, ".instructions")
}
