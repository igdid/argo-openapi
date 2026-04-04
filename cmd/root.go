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
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"net/url"
	"os"
)

var log = logrus.New()

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "argo-openapi",
	Short: "Argo-OpenAPI creates Argo Events Custom Sensors from OpenAPI specification",
	Long: `Examples:

# Create a Sensor for sending messages to Telegram
argo-openapi -s https://raw.githubusercontent.com/alserom/telegram-bot-api-spec/refs/heads/main/openapi.json -t ./examples/telegram generate client

# Validate an OpenAPI specification and show endpoint list
argo-openapi -s https://raw.githubusercontent.com/alserom/telegram-bot-api-spec/refs/heads/main/openapi.json validate
`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// Common argo-openapi options 
type options struct {
	// From cmd/root.go
	src         string
	target      string
	// From cmd/generate.go
	license     string
	author      string
	cobraConfig string
}

var opts options

func init() {
	rootCmd.PersistentFlags().StringVarP(&opts.src, "src", "s", "./openapi.yaml", "Path to openapi protocol. It can be an URL or a path")
	rootCmd.MarkFlagRequired("src")
	rootCmd.PersistentFlags().StringVarP(&opts.target, "target", "t", "", "Path to a result sensor. It is required if src is an URL")
}
