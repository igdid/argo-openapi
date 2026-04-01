package utils

import (
	"embed"
	"os"
)

//go:embed templates/*
var templates embed.FS

func UnpackFile(fn, targetPath string) error {
	content, err := templates.ReadFile(fn)
	if err != nil {
		return err
	}
	return os.WriteFile(targetPath, content, 0600)
}
