package main

import (
	"net/http"

	"github.com/drone/drone-go/plugin/config"
	"github.com/fbcbarbosa/drone-ignore-config/plugin"

	_ "github.com/joho/godotenv/autoload"
	"github.com/kelseyhightower/envconfig"
	"github.com/onrik/logrus/filename"
	log "github.com/sirupsen/logrus"
)

type spec struct {
	Debug   bool   `envconfig:"PLUGIN_DEBUG"`
	Address string `envconfig:"PLUGIN_ADDRESS" default:":3000"`
	Secret  string `envconfig:"PLUGIN_SECRET"`
	Token   string `envconfig:"GITHUB_TOKEN"`
	Server  string `envconfig:"GITHUB_SERVER"`
}

func main() {
	spec := new(spec)
	err := envconfig.Process("", spec)
	if err != nil {
		log.Fatal(err)
	}

	log.AddHook(filename.NewHook())

	if spec.Debug {
		log.SetLevel(log.DebugLevel)
	}
	if spec.Secret == "" {
		log.Fatalln("missing secret key")
	}
	if spec.Token == "" {
		log.Warnln("missing github token")
	}
	if spec.Address == "" {
		spec.Address = ":3000"
	}

	handler := config.Handler(
		plugin.New(
			spec.Server,
			spec.Token,
		),
		spec.Secret,
		log.StandardLogger(),
	)

	log.Infof("server listening on address %s", spec.Address)

	http.Handle("/", handler)
	log.Fatal(http.ListenAndServe(spec.Address, nil))
}
