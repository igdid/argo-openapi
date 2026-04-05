/*
{{ .Copyright }}
{{ if .Legal.Header }}{{ .Legal.Header }}{{ end }}
*/
package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"{{ .PkgName }}/proto"
	"os"
	"fmt"
	"net"
	"google.golang.org/grpc"
)

var log = logrus.New()

var rootCmd = &cobra.Command{
	Use:   "{{ .AppName }}",
	Short: "Argo events openapi generated sensor",
	Long:  `{{ .AppName }} gets events from argo events and sends them to a remote server as it said in its openapi specification`,
	Run:   runSender,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func runSender(cmd *cobra.Command, args []string) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", opts.port))
	if err != nil {
		log.Fatal(err)
	}
	srv := grpc.NewServer()

	sender := Sender{}
	proto.RegisterTriggerServer(srv, &sender)
	log.Infof("Starting server on :%d, log level: %s, metrics: %v", opts.port, opts.logLevel, opts.metrics)
	if err := srv.Serve(listener); err != nil {
		log.Fatal(err)
	}
	if opts.metrics {
		log.Infof("Listenings metrics on :%d", opts.metricsPort)
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
		log.Infof("Using config file: %s", viper.ConfigFileUsed())
	}
}
