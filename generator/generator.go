package generator

import (
	"embed"
	"os"
	"path/filepath"
	"text/template"
)

//go:embed templates/*
var templateFS embed.FS

type ProjectConfig struct {
	ProjectName string
	Database    string
	IncludeSQLAlchemy bool
}

func CreateProject(cfg ProjectConfig) error {
	
	err := os.MkdirAll(cfg.ProjectName, 0755)
	if err != nil {
		return err
	}

	tmpl, err := template.ParseFS(templateFS, "templates/main.py.tmpl")
	if err != nil {
		return err
	}

	filePath := filepath.Join(cfg.ProjectName, "main.py")
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	return tmpl.Execute(f, cfg)
}