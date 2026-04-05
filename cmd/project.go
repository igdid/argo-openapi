package cmd

import (
	"context"
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/igdid/argo-openapi/internal/utils"
	"github.com/oapi-codegen/oapi-codegen/v2/pkg/codegen"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

type Project struct {
	PkgName     string
	Copyright   string
	Legal       License
	AppName     string
	GoVersion   string
	rootDir     string
	rootDoc     *openapi3.T
	protoTarget string
}

var project Project

// Create OpenAPI protocol files
func (p *Project) createOpenAPIFiles(fn string) {
	p.rootDoc = load(opts.src)
	cfg := codegen.Configuration{
		PackageName: "proto",
		Generate: codegen.GenerateOptions{
			Models: true,
			Client: true,
			//GorillaServer:  true,
			Strict: true,
		},
		OutputOptions: codegen.OutputOptions{
			SkipPrune: true,
		},
	}
	code, err := codegen.Generate(p.rootDoc, cfg)
	if err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(p.protoTarget, fn), []byte(code), 0644); err != nil {
		log.Fatal(err)
	}
}

func showParams(method, name string, operation *openapi3.Operation) {
	if operation == nil {
		return
	}
	l := log.WithFields(logrus.Fields{
		"operationID": operation.OperationID,
	})
	if len(operation.Parameters) == 0 {
		l.Info(strings.ToUpper(method), " ", name)
	}
	for _, p := range operation.Parameters {
		l.WithFields(logrus.Fields{
			"param": p.Value.Name,
			"in":    p.Value.In,
		}).Info(strings.ToUpper(method), " ", name)
	}
}

func (p *Project) validateOpenAPI(cmd *cobra.Command, args []string) {
	p.rootDoc = load(opts.src)
	ctx := context.Background()

	if err := p.rootDoc.Validate(ctx); err != nil {
		log.Fatalf("invalid spec: %v", err)
	}
	log.Info("OpenAPI spec is valid. The following methods available:")
	for name, path := range p.rootDoc.Paths.Map() {
		showParams("get", name, path.Get)
		showParams("post", name, path.Post)
		showParams("put", name, path.Put)
		showParams("delete", name, path.Delete)
		showParams("options", name, path.Options)
		showParams("patch", name, path.Patch)
		showParams("trace", name, path.Trace)
	}
}

func isURL(s string) bool {
	u, err := url.ParseRequestURI(s)
	if err != nil {
		return false
	}

	return u.Scheme != "" && u.Host != ""
}

func load(src string) *openapi3.T {
	var doc *openapi3.T
	var err error
	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	if isURL(src) {
		url, _ := url.ParseRequestURI(src)
		doc, err = loader.LoadFromURI(url)
	} else {
		doc, err = loader.LoadFromFile(src)
	}
	if err != nil {
		log.Fatalf("load spec: %v - %s", err, src)
	}
	if doc.Paths == nil {
		log.Fatal("No paths are inside the file. Nothing to do")
	}
	return doc
}

// Create common files
func (p *Project) generateSensorProject() {
	// Create dir hierarchy
	err := os.MkdirAll(p.protoTarget, 0754)
	if err != nil {
		log.Fatal(err)
	}
	cmdDir := filepath.Join(p.rootDir, "cmd")
	err = os.MkdirAll(cmdDir, 0754)
	if err != nil {
		log.Fatal(err)
	}
	p.renderFile("main.go.tpl")
	p.renderFile("Dockerfile")
	p.renderFile("cmd/root.go")
	p.renderFile("cmd/sender.go")
}

func (p Project) renderFile(name string) {
	// Create file
	fn := strings.TrimSuffix(name, ".tpl")
	f, err := os.Create(fmt.Sprintf("%s/%s", p.rootDir, fn))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Render template
	tpl, _ := utils.Templates.ReadFile(fmt.Sprintf("templates/%s", name))
	ftpl := template.Must(template.New(name).Parse(string(tpl)))
	err = ftpl.Execute(f, p)
	if err != nil {
		log.Fatal(err)
	}
}

func (p Project) tidy() {
	p.renderFile("go.mod.tpl")
	log.Info("Running `go mod tidy`")
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = p.rootDir
	//cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatalf("`go mod tidy` failed: %v", err)
	}
}

// Create Protobuf protocol files
func (p Project) generateProtobuf() {
	triggerFn := "/tmp/trigger.proto"
	err := utils.UnpackFile("templates/trigger.proto", triggerFn)
	if err != nil {
		log.Fatal(err, triggerFn)
	}
	defer os.Remove(triggerFn)
	cmd := exec.Command(
		"protoc",
		"-I/tmp",
		"--go_out="+p.protoTarget,
		"--go-grpc_out="+p.protoTarget,
		"--go_opt=paths=source_relative",
		"--go-grpc_opt=paths=source_relative",
		"--go_opt=Mtrigger.proto="+p.protoTarget,
		"--go-grpc_opt=Mtrigger.proto="+p.protoTarget,
		"/tmp/trigger.proto",
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatalf("protoc failed: %v", err)
	}
}
