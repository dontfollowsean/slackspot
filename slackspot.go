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

const redirectURI = "http://localhost:8080/callback"

var (
	auth   = spotify.NewAuthenticator(redirectURI, spotify.ScopeUserReadCurrentlyPlaying, spotify.ScopeUserReadRecentlyPlayed)
	ch     = make(chan *spotify.Client)
	state  = "bx"
	client *spotify.Client
)

func main() {
	// todo get from env
	//slackApi := slack.New("Vs5ajYrxthVZ6sYixzVu6yo4")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		url := auth.AuthURL(state)
		_, _ = fmt.Fprintf(w, "Spotify-Slack Integration\n Please log in to Spotify by visiting the following page in your browser: %s", url)
	})
	http.HandleFunc("/callback", completeAuth)
	http.HandleFunc("/nowplaying", nowPlayingHandler)
	http.HandleFunc("/slack/nowplaying", slackNowPlayingHandler)
	//http.HandleFunc("/lastplayed", lastPlayedHandler)
	//http.HandleFunc("/slack/lastplayed", lastPlayedHandler)


	// auth
	go func() {
		url := auth.AuthURL(state)
		fmt.Println("Please log in to Spotify by visiting the following page in your browser:", url)

		// wait for auth to complete
		client = <-ch

		// use the client to make calls that require authorization
		user, err := client.CurrentUser()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("logged in as:", user.ID)

	}()

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}

}

func slackNowPlayingHandler(w http.ResponseWriter, r *http.Request) {
	const signingSecret = "a1f7c0caf421f4c61def057c4b1c7cf9"
	verifier, err := slack.NewSecretsVerifier(r.Header, signingSecret)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	r.Body = ioutil.NopCloser(io.TeeReader(r.Body, &verifier))
	s, err := slack.SlashCommandParse(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = verifier.Ensure(); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	song, err := client.PlayerCurrentlyPlaying()
	var titleLink string
	if song != nil && song.Item != nil{
		for k,v := range song.Item.ExternalURLs {
			log.Printf("%s: %s", k, v)
		}
		titleLink = song.Item.ExternalURLs["spotify"]
	}

	switch s.Command {
	case "/nowplaying":
		params := &slack.Msg{
			Text: "🎵 Now playing...",
			Attachments: []slack.Attachment{
				{
					Title: PrintNowPlaying(client),
					TitleLink: titleLink,
				},
			},
		}
		b, err := json.Marshal(params)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	default:
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func nowPlayingHandler(w http.ResponseWriter, r *http.Request) {
	if client == nil {
		_, _ = fmt.Fprint(w, "Please Log into the BounceX Spotify Account")
		return
	}

	nowPlaying := PrintNowPlaying(client)
	_, _ = fmt.Fprint(w, nowPlaying)
}

func completeAuth(w http.ResponseWriter, r *http.Request) {
	tok, err := auth.Token(state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Print(err)
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Printf("State mismatch: %s != %s\n", st, state)
	}
	// use the token to get an authenticated client
	client := auth.NewClient(tok)
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "Login Completed!")
	ch <- &client
}
