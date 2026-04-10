// Copyright © 2021 Steve Francia <spf@spf13.com>.
// Copyright © 2026 Igor Diakonov <igor@linux.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Modified 05.04.2026 by Igor Diakonov

package cmd

import (
	"context"
	"fmt"
	"github.com/Masterminds/sprig/v3"
	"github.com/fatih/color"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/igdid/argo-openapi/internal/utils"
	"github.com/igdid/argo-openapi/internal/openapiutil"
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
	GoArch      string
	Operations  map[string]openapi3.Parameters
	rootDir     string
	rootDoc     *openapi3.T
	protoTarget string
}

var project Project

func (p *Project) createLicenseFile() error {
	data := map[string]interface{}{
		"copyright": copyrightLine(),
	}
	licenseFile, err := os.Create(fmt.Sprintf("%s/LICENSE", p.rootDir))
	if err != nil {
		return err
	}
	defer licenseFile.Close()

	licenseTemplate := template.Must(template.New("license").Parse(p.Legal.Text))
	return licenseTemplate.Execute(licenseFile, data)
}

// Create OpenAPI protocol files
func (p Project) createOpenAPIFiles(fn string) {
	cfg := codegen.Configuration{
		PackageName: "proto",
		Generate: codegen.GenerateOptions{
			Models:     true,
			Client:     true,
			ServerURLs: true,
			Strict:     true,
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
	l.Info(strings.ToUpper(method), " ", name)
}

func (p Project) validateOpenAPI(cmd *cobra.Command, args []string) {
	ctx := context.Background()

	if err := p.rootDoc.Validate(ctx); err != nil {
		log.Fatalf("invalid spec: %v", err)
	}
	log.Info("OpenAPI spec is valid")
}

func (p *Project) operationInfo(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		for name, path := range p.rootDoc.Paths.Map() {
			showParams("get", name, path.Get)
			showParams("post", name, path.Post)
			showParams("put", name, path.Put)
			showParams("delete", name, path.Delete)
			showParams("options", name, path.Options)
			showParams("patch", name, path.Patch)
			showParams("trace", name, path.Trace)
		}
	} else {
		opId := args[0]
		op, method, path, err := openapiutil.FindOperation(p.rootDoc, opId)
		if err != nil {
			log.Fatalf("FindOperation: %s - %v", opId, err)
		}
		color.New(color.FgGreen).Print(method)
		fmt.Print(" ")
		color.Magenta(path)
		if op.Summary != "" {
			color.New(color.FgYellow).Print("Summary: ")
			fmt.Print(op.Summary)
			fmt.Println()
		}
		if op.Description != "" {
			color.New(color.FgYellow).Print("Description: ")
			fmt.Print(op.Description)
			fmt.Println()
		}
		color.New(color.FgYellow).Print("Deprecated: ")
		fmt.Printf("%v", op.Deprecated)
		fmt.Println()
		if len(op.Parameters) > 0 {
			color.New(color.FgYellow).Println("Parameters: ")
		}
		for _, param := range op.Parameters {
			if param.Ref != "" {
				log.Infof("- %s", param.Ref)
			} else {
				required := ""
				if param.Value.Required {
					required = "(required)"
				}
				fmt.Printf("- %s %s\n", param.Value.Name, required)
			}
		}
	}
}

func isURL(s string) bool {
	u, err := url.ParseRequestURI(s)
	if err != nil {
		return false
	}

	return u.Scheme != "" && u.Host != ""
}

func (p *Project) load(src string) {
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
	p.rootDoc = doc

	p.Operations = map[string]openapi3.Parameters{}
	for _, pathItem := range doc.Paths.Map() {
		for _, op := range pathItem.Operations() {
			if op != nil {
				p.Operations[strings.Title(op.OperationID)] = op.Parameters
			}
		}
	}
}

// Create common files
func (p Project) generateSensorProject() {
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
	p.renderFile("main.go.tpl", true)
	p.renderFile("Dockerfile", true)
	p.renderFile(".dockerignore.tpl", true)
	p.renderFile("cmd/root.go", true)
	p.renderFile("cmd/sender.go", true)
	p.renderFile("cmd/operations.go", true)
	p.renderFile("cmd/config.go", true)
	// Default config shouldn't be rewritten after manual changes
	p.renderFile("config.yaml", false)
}

func (p Project) renderFile(name string, rewrite bool) {
	// Create file
	fn := fmt.Sprintf("%s/%s", p.rootDir, strings.TrimSuffix(name, ".tpl"))
	if _, err := os.Stat(fn); err == nil && !rewrite {
		// no need to continue if file exists and it shouldn't be rewritten
		return
	}
	f, err := os.Create(fn)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Render template
	tpl, _ := utils.Templates.ReadFile(fmt.Sprintf("templates/%s", name))
	ftpl := template.Must(template.New(name).Funcs(sprig.TxtFuncMap()).Parse(string(tpl)))
	err = ftpl.Execute(f, p)
	if err != nil {
		log.Fatal(err)
	}
}

func (p Project) tidy() {
	// go.mod shouldn't be rewritten after manual changes
	p.renderFile("go.mod.tpl", false)
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
