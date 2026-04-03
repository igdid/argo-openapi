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
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
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

type options struct {
	src         string
	target      string
	author      string
	cobraConfig string
	license     string
	rootDoc     *openapi3.T
}

var opts options

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.argo-openapi.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.PersistentFlags().StringVarP(&opts.src, "src", "s", "./openapi.yaml", "Path to openapi protocol. It can be an URL or a path")
	rootCmd.MarkFlagRequired("src")
	rootCmd.PersistentFlags().StringVarP(&opts.target, "target", "t", "", "Path to a result sensor. It is required if src is an URL")
}
