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
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/spf13/cobra"
)

// validateCmd represents the validate command
var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate OpenAPI protocol",
	Run: validateOpenAPI,
}

func validateOpenAPI(cmd *cobra.Command, args []string) {
	ctx := context.Background()

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true // если есть $ref на другие файлы

	doc, err := loader.LoadFromFile(opts.filename)
	if err != nil {
		log.Fatalf("load spec: %v", err)
	}
	if err := doc.Validate(ctx); err != nil {
		log.Fatalf("invalid spec: %v", err)
	}
	if doc.Paths != nil {
		for name, path := range doc.Paths.Map() {
			if path.Get != nil {
				if len(path.Get.Parameters) == 0 {
					log.Info("GET ", name)
				}
				for _, p := range path.Get.Parameters {
					log.Info("GET ", name, " ", p.Value.Name, " ", p.Value.In)
				}
			}
			if path.Post != nil {
				if len(path.Post.Parameters) == 0 {
					log.Info("POST ", name)
				}
				for _, p := range path.Post.Parameters {
					log.Info("POST ", name, p)
				}
			}
			if path.Put != nil {
				if len(path.Put.Parameters) == 0 {
					log.Info("PUT ", name)
				}
				for _, p := range path.Put.Parameters {
					log.Info("PUT ", name, p)
				}
			}
		}
	}
	fmt.Println("OpenAPI spec is valid")
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
