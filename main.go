package main

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	SlackSigingSecret  string
	SpotifyRedirectURI string
	ContactUser        string
	SongHistoryLength  int
}

func main() {
	songHistoryLength, err := strconv.Atoi(getEnv("SONG_HISTORY_LENGTH", "3"))
	if err != nil {
		songHistoryLength = 3
	}
	config := &Config{
		SlackSigingSecret:  getEnv("SLACK_SIGNING_SECRET", ""),
		SpotifyRedirectURI: getEnv("SPOTIFY_REDIRECT_URI", ""),
		ContactUser:        getEnv("CONTACT_SLACK_USER", "an Administrator"),
		SongHistoryLength:  songHistoryLength,
	}

	if err := run(config); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
