/*
{{ .Copyright }}
{{ if .Legal.Header }}{{ .Legal.Header }}{{ end }}
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "{{ .AppName }}",
	Short: "Argo events openapi generated sensor",
	Long:  `{{ .AppName }} gets events from argo events and sends them to a remote server as it said in its openapi specification`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

type options struct {
	cfgFile     string
	logLevel    string
	port        int
	metrics     bool
	metricsPort int
}

var opts options

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&opts.cfgFile, "config", "c", "/etc/{{ .AppName }}/config.yaml", "config file path (default is /etc/{{ .AppName }}/config.yaml)")
	rootCmd.PersistentFlags().StringVarP(&opts.logLevel, "log-level", "l", "info", "Log level for the server")
	rootCmd.PersistentFlags().IntVarP(&opts.port, "port", "p", 8081, "GRPC port for incomming events")
	rootCmd.PersistentFlags().BoolVar(&opts.metrics, "enable-metrics", false, "Enable metrics")
	rootCmd.PersistentFlags().IntVar(&opts.metricsPort, "metrics-port", 8080, "HTTP port for serving metrics")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if opts.cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(opts.cfgFile)
	} else {
		viper.AddConfigPath("/etc/{{ .AppName }}")
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
