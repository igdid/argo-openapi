package cmd

import (
	"github.com/igdid/argo-openapi/internal/licenses"
	"os"
	"text/template"
)

type Project struct {
	// v2
	PkgName   string
	Copyright string
	Legal     License
	AppName   string
}

// Create common files
func generateSensorProject(targetDir string) {
	if !isURL(opts.src) {
		opts.target = filepath.Dir(opts.src)
	} else if opts.target == "" {
		log.Fatal("A target must be specified")
	}
	rootDir := filepath.Join(opts.target, "sensor")
	// Create dir hierarchy
	opts.target = filepath.Join(rootDir, "proto")
	err := os.MkdirAll(opts.target, 0754)
	if err != nil {
		log.Fatal(err)
	}
	cmdDir := filepath.Join(rootDir, "cmd")
	err = os.MkdirAll(cmdDir, 0754)
	if err != nil {
		log.Fatal(err)
	}
	// Create main file
	mainFile, err := os.Create(fmt.Sprintf("%s/main.go", rootDir))
	if err != nil {
		return err
	}
	defer mainFile.Close()

	// Render main template
	mainTemplate := template.Must(template.New("main").Parse(utils.Templates.ReadFile("templates/main.go.tpl")))
	err = mainTemplate.Execute(mainFile, p)
	if err != nil {
		return err
	}
	// Create main file
	rootCmdFile, err := os.Create(fmt.Sprintf("%s/root.go", cmdDir))
	if err != nil {
		return err
	}
	defer rootCmdDir.Close()

	// Render main template
	rootTemplate := template.Must(template.New("root").Parse(utils.Templates.ReadFile("templates/cmd/root.go.tpl")))
	err = rootTemplate.Execute(rootCmdFile, p)
	if err != nil {
		return err
	}

}
