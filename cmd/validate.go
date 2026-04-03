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
	"context"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"strings"
)

// validateCmd represents the validate command
var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate OpenAPI protocol",
	Run:   validateOpenAPI,
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

func validateOpenAPI(cmd *cobra.Command, args []string) {
	opts.rootDoc = load(opts.src)
	ctx := context.Background()

	if err := opts.rootDoc.Validate(ctx); err != nil {
		log.Fatalf("invalid spec: %v", err)
	}
	log.Info("OpenAPI spec is valid. The following methods available:")
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
	rootCmd.AddCommand(validateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// validateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// validateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
