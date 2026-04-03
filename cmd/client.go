/*
Copyright © 2026 Igor Diakonov <igor@linux.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"github.com/igdid/argo-openapi/internal/utils"
	"github.com/oapi-codegen/oapi-codegen/v2/pkg/codegen"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"path/filepath"
)

// clientCmd represents the client command
var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Generate a client code",
	Run:   generateSensor,
}

func generateProtobuf(goPackage, targetDir string) {
	triggerFn := "/tmp/trigger.proto"
	err := utils.UnpackFile("templates/trigger.proto", triggerFn)
	if err != nil {
		log.Fatal(err, triggerFn)
	}
	defer os.Remove(triggerFn)

	cmd := exec.Command(
		"protoc",
		"-I/tmp",
		"--go_out="+targetDir,
		"--go-grpc_out="+targetDir,
		"--go_opt=paths=source_relative",
		"--go-grpc_opt=paths=source_relative",
		"--go_opt=Mtrigger.proto="+goPackage,
		"--go-grpc_opt=Mtrigger.proto="+goPackage,
		"/tmp/trigger.proto",
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	log.Info("Generating Protobuf protocol")
	if err := cmd.Run(); err != nil {
		log.Fatalf("protoc failed: %v", err)
	}
}

func generateSensor(cmd *cobra.Command, args []string) {
	// Create project dirs
	if !isURL(opts.src) {
		opts.target = filepath.Dir(opts.src)
	} else if opts.target == "" {
		log.Fatal("A target must be specified")
	}
	opts.target = filepath.Join(opts.target, "sensor", "proto")
	err := os.MkdirAll(opts.target, 0755)
	if err != nil {
		log.Fatal(err)
	}

	// Create OpenAPI protocol files
	opts.rootDoc = load(opts.src)
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
	code, err := codegen.Generate(opts.rootDoc, cfg)
	if err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(opts.target, "sensor.gen.go"), []byte(code), 0644); err != nil {
		log.Fatal(err)
	}
	// Create Protobuf protocol files
	goPackage, err := utils.GetGoPackage(opts.target)
	if err != nil {
		log.Fatal(err)
	}
	generateProtobuf(goPackage, opts.target)
}

func init() {
	generateCmd.AddCommand(clientCmd)
}
