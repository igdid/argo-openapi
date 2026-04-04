package utils

import (
	"embed"
	"os"
)

//go:embed templates/*
var Templates embed.FS

func UnpackFile(fn, targetPath string) error {
	content, err := Templates.ReadFile(fn)
	if err != nil {
		return err
	}
	return os.WriteFile(targetPath, content, 0600)
}
