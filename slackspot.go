package main

import (
	"encoding/json"
	"fmt"
	"github.com/nlopes/slack"
	"github.com/zmb3/spotify"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

const redirectURI = "http://localhost:8080/callback" // auth callback
var spotifyClient *SpotifyClient

func main() {
	// todo get from env
	//slackApi := slack.New("Vs5ajYrxthVZ6sYixzVu6yo4")// todo regenerate me

	spotifyClient = &SpotifyClient{
		Client:        nil,
		Authenticator: spotify.NewAuthenticator(redirectURI, spotify.ScopeUserReadCurrentlyPlaying, spotify.ScopeUserReadRecentlyPlayed),
		State:         "bx",
		Channel:       make(chan *spotify.Client),
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		url := spotifyClient.Authenticator.AuthURL(spotifyClient.State)
		_, _ = fmt.Fprintf(w, "Spotify-Slack Integration\n Please log in to Spotify by visiting the following page in your browser: %s", url)
	})
	http.HandleFunc("/callback", completeAuth)
	http.HandleFunc("/nowplaying", nowPlayingHandler)
	http.HandleFunc("/slack", slackHandler)

	go spotifyClient.Login()

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}

}

func slackHandler(w http.ResponseWriter, r *http.Request) {
	const signingSecret = "a1f7c0caf421f4c61def057c4b1c7cf9" // todo regenerate me and get from env
	verifier, err := slack.NewSecretsVerifier(r.Header, signingSecret)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	r.Body = ioutil.NopCloser(io.TeeReader(r.Body, &verifier))
	s, err := slack.SlashCommandParse(r)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = verifier.Ensure(); err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	switch s.Command {
	case "/nowplaying":
		song, err := spotifyClient.NowPlaying()
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		attachments := []slack.Attachment{
			{
				Title:     fmt.Sprintf("%s by %s", song.title, song.artist),
				TitleLink: song.url,
			},
		}
		slackMsg := &slack.Msg{
			Text:        "🎵 Now playing...",
			Attachments: attachments,
		}
		b, err := json.Marshal(slackMsg)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(b)
	case "/lastplayed":
		songs, err := spotifyClient.RecentlyPlayed()
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		attachments := make([]slack.Attachment, 0);
		for _, song := range songs {
			attachment := slack.Attachment{
				Title:     fmt.Sprintf("%s by %s", song.title, song.artist),
				TitleLink: song.url,
			}
			attachments = append(attachments, attachment)
		}
		slackMsg := &slack.Msg{
			Text:        "🎵 Recently Played Songs",
			Attachments: attachments,
		}
		b, err := json.Marshal(slackMsg)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(b)
	default:
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func nowPlayingHandler(w http.ResponseWriter, r *http.Request) {
	if spotifyClient.Client == nil {
		_, _ = fmt.Fprint(w, "Please Log into the BounceX Spotify Account")
		go spotifyClient.Login()
		return
	}

	nowPlaying, _ := spotifyClient.NowPlaying()
	_, _ = fmt.Fprint(w, nowPlaying)
}

func completeAuth(w http.ResponseWriter, r *http.Request) {
	tok, err := spotifyClient.Authenticator.Token(spotifyClient.State, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Print(err)
	}
	if st := r.FormValue("state"); st != spotifyClient.State {
		http.NotFound(w, r)
		log.Printf("State mismatch: %s != %s\n", st, spotifyClient.State)
	}
	// use the token to get an authenticated client
	client := spotifyClient.Authenticator.NewClient(tok)
	w.Header().Set("Content-Type", "text/html")
	_, _ = fmt.Fprintf(w, "Login Completed!")
	spotifyClient.Channel <- &client
}
