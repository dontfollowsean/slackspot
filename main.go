package main

import (
	"fmt"
	"github.com/caarlos0/env"
	"github.com/kr/pretty"
	"log"
	"os"
)

type Config struct {
	SlackSigingSecret  string `env:"SLACK_SIGNING_SECRET"`
	SlackAdminWebhook  string `env:"SLACK_ADMIN_WEBHOOK"`
	SpotifyRedirectURI string `env:"SPOTIFY_REDIRECT_URI"`
	ContactUser        string `env:"CONTACT_SLACK_USER"`
	SongHistoryLength  int    `env:"SONG_HISTORY_LENGTH" envDefault:"3"`
}

//var config *Config

func (c Config) getDetails() interface{} {
	return pretty.Sprint(c)
}

func main() {
	config := getAppConfig()

	if err := run(config); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
	}
}

func getAppConfig() *Config {
	config := &Config{}
	err := env.Parse(config)
	if err != nil {
		log.Printf(err.Error())
	}

	log.Printf("Starting app with config : %s", config.getDetails())

	return config
}
