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
	"path/filepath"
)

// clientCmd represents the action after `argo-openapi generate client`
var clientCmd = &cobra.Command{
	Use:   "client [app-name]",
	Short: "Generate a client code for [app-name] integration",
	Args:  cobra.ExactArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		// Fill the project structure
		rootDir := filepath.Join(opts.target, "sensor")
		goPackage, err := utils.GetGoPackage(rootDir)
		if err != nil {
			// if not a git repo
			log.Fatal(err)
		}
		project = Project{
			AppName:   args[0],
			Copyright: copyrightLine(),
			Legal:     getLicense(),
			PkgName:   goPackage,
			GoVersion: utils.CurrentGoVersion(),
			rootDir:   rootDir,
		}
		project.protoTarget = filepath.Join(project.rootDir, "proto")
		log.WithFields(logrus.Fields{
			"Go Package": project.PkgName,
			"App Name":   project.AppName,
			"Copyright":  project.Copyright,
			"License":    project.Legal.Name,
		}).Infof("Project is in %s", project.rootDir)
	},
	Run: generateSensor,
}

func generateSensor(cmd *cobra.Command, args []string) {
	// Start code generation
	project.generateSensorProject()
	project.createOpenAPIFiles("sensor.gen.go")
	project.generateProtobuf()
	log.Infof("%s code generated successfully in %s", project.AppName, opts.target)
}

func init() {
	generateCmd.AddCommand(clientCmd)
}
