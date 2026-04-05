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
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"path/filepath"
)

// clientCmd represents the client command
var clientCmd = &cobra.Command{
	Use:   "client [app-name]",
	Short: "Generate a client code for [app-name] integration",
	Run:   generateSensor,
	Args:  cobra.ExactArgs(1),
}

func generateSensor(cmd *cobra.Command, args []string) {
	goPackage, err := utils.GetGoPackage(opts.target)
	if err != nil {
		// if not a git repo
		log.Fatal(err)
	}
	// Fill the project structure
	project = Project{
		AppName:   args[0],
		Copyright: copyrightLine(),
		Legal:     getLicense(),
		PkgName:   goPackage,
	}
	
	project.rootDir = filepath.Join(project.protoTarget, "sensor")
	project.protoTarget = filepath.Join(project.rootDir, "proto")

	// Start code generation
	log.WithFields(logrus.Fields{
		"Go Package": project.PkgName,
		"App Name":   project.AppName,
		"Copyright":  project.Copyright,
		"License":    project.Legal.Name,
	}).Info("Starting code generation")
	project.generateSensorProject()
	project.createOpenAPIFiles("sensor.gen.go")
	project.generateProtobuf(opts.target)
	log.Infof("%s code generated successfully in %s", project.AppName, opts.target)
}

func init() {
	generateCmd.AddCommand(clientCmd)
}

// Create Protobuf protocol files
func (p Project) generateProtobuf(targetDir string) {
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
		"--go_opt=Mtrigger.proto="+opts.target,
		"--go-grpc_opt=Mtrigger.proto="+opts.target,
		"/tmp/trigger.proto",
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	log.Infof("Generating Protobuf protocol in %s", targetDir)
	if err := cmd.Run(); err != nil {
		log.Fatalf("protoc failed: %v", err)
	}
}
