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
	"strings"
	"github.com/spf13/cobra"
	"github.com/getkin/kin-openapi/openapi3"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate Argo Events Sensor",
	Run: generateSensor,
}

func showParams(method, name string, operation *openapi3.Operation) {
	if operation == nil {
		return
	}
	if len(operation.Parameters) == 0 {
		log.Info(strings.ToUpper(method), " ", name)
	}
	for _, p := range operation.Parameters {
		log.Info(strings.ToUpper(method), " ", name, " - ", p.Value.Name, " in ", p.Value.In)
	}
}

func generateSensor(cmd *cobra.Command, args []string) {
	opts.rootDoc = load(opts.filename)
	for name, path := range opts.rootDoc.Paths.Map() {
		showParams("get", name, path.Get)
		showParams("post", name, path.Post)
		showParams("put", name, path.Put)
		showParams("delete", name, path.Delete)
		showParams("options", name, path.Options)
		showParams("patch", name, path.Patch)
		showParams("trace", name, path.Trace)
	}
}
func init() {
	rootCmd.AddCommand(generateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// generateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// generateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
