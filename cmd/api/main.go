package main

import (
	"net/http"
	"strings"

	"backend/configs"
	"backend/internal/api/interfaces"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func main() {
	cobra.CheckErr(serveCmd.Execute())
}

var serveCmd = &cobra.Command{
	Use:   "api",
	Short: "Starts the api server, serving backend APIs",
	Long: `'api' starts the api server, serving the the Backend APIs.

	Env                            Default
	---------------                --------------
	ENV                            dev
	PORT                           8080
	LOG_LEVEL                      INFO
	`,
	Run: serveRun,
}

var bindEnvironments = []string{
	configs.Env,
	configs.EnvPort,
	configs.EnvLogLevel,
}

func init() {
	// set environments
	viper.SetDefault(configs.Env, configs.DefaultEnv)
	viper.SetDefault(configs.EnvPort, configs.DefaultPort)
	viper.SetDefault(configs.EnvLogLevel, configs.DefaultLogLevel)

	for _, e := range bindEnvironments {
		viper.BindEnv(e)
	}
}

func serveRun(cmd *cobra.Command, args []string) {
	// env
	env := strings.ToLower(viper.GetString(configs.Env))
	if !(env == "dev" || env == "test" || env == "staging" || env == "production") {
		logrus.Fatalln(errors.Errorf("unrecognized config variable: %s\n", env))
	}

	// port
	port := viper.GetString(configs.EnvPort)

	// log level
	lvl := viper.GetString(configs.EnvLogLevel)
	level, err := logrus.ParseLevel(lvl)
	if err != nil {
		logrus.Fatalln(errors.Errorf("unrecognized config variable: %s\n", lvl))
	}
	logrus.SetLevel(level)

	for _, e := range bindEnvironments {
		if v := viper.GetString(e); v == "" {
			logrus.Fatalln(errors.Errorf("required " + e + " environment value"))
		}
	}
	logrus.SetFormatter(&logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "timestamp",
			logrus.FieldKeyLevel: "severity",
			logrus.FieldKeyMsg:   "jsonPayload",
		},
	})

	startAPIServer(port)
}

func startAPIServer(port string) {
	app := interfaces.NewApp()
	http.ListenAndServe(":"+port, app.Router())
}
