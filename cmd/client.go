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
	"os"
	"path/filepath"
	"strings"
	"github.com/spf13/cobra"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/oapi-codegen/oapi-codegen/v2/pkg/codegen"
)

// clientCmd represents the client command
var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Generate a client code",
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
	cfg := codegen.Configuration{
		PackageName: "sensor",
		Generate: codegen.GenerateOptions{
			Models:         true,
			GorillaServer:  true,
			Strict:         true,
		},
		OutputOptions: codegen.OutputOptions{
			SkipPrune: true,
		},
	}
	code, err := codegen.Generate(opts.rootDoc, cfg)
	if err != nil {
		log.Fatal(err)
	}
	dir := filepath.Dir(opts.filename)
	dir = filepath.Join(dir, "sensor")
	err = os.MkdirAll(dir, 0755)
	if err != nil {
	    log.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "sensor.gen.go"), []byte(code), 0644); err != nil {
		log.Fatal(err)
	}
}

func init() {
	generateCmd.AddCommand(clientCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// clientCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// clientCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
