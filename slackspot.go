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
	"os"
	"strconv"
)

const (
	contentType = "Content-Type"
	jsonType    = "application/json"
	textHtml    = "text/html"
)

var (
	spotifyClient      *SpotifyClient
	slackSigningSecret string
	redirectURI        string
	songHistoryLength  int
	contactUser        string
	host               string
)

func init() {
	slackSigningSecret = getEnv("SLACK_SIGNING_SECRET", "")
	host = getEnv("SLACKSPOT_HOST", "")
	redirectURI = fmt.Sprintf("%s/callback", host)
	var err error
	songHistoryLength, err = strconv.Atoi(getEnv("SONG_HISTORY_LENGTH", "3"))
	if err != nil {
		songHistoryLength = 3
	}
	spotifyClient = &SpotifyClient{
		Client:        nil,
		Authenticator: spotify.NewAuthenticator(redirectURI, spotify.ScopeUserReadCurrentlyPlaying, spotify.ScopeUserReadRecentlyPlayed),
		State:         "bx",
		Channel:       make(chan *spotify.Client),
	}
	contactUser = getEnv("CONTACT_SLACK_USER", "an Administrator")
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/callback", completeAuth)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/nowplaying", nowPlayingHandler)
	http.HandleFunc("/recentlyplayed", recentlyPlayedHandler)
	http.HandleFunc("/slack", slackHandler)

	go spotifyClient.Login()

	go func() {
		log.Print("Listening on port 443")
		err := http.ListenAndServeTLS(":443", "_certs/certificate.crt", "_certs/private.key", nil)
		if err != nil {
			log.Fatal("ListenAndServeTLS: ", err)
		}
	}()

	log.Print("Listening on port 80")
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	go spotifyClient.Login()
	w.Header().Set(contentType, textHtml)
	_, _ = fmt.Fprint(w, "A link to log in to Spotify has been sent to the Slack Admin")
}

func slackHandler(w http.ResponseWriter, r *http.Request) {
	verifier, err := slack.NewSecretsVerifier(r.Header, slackSigningSecret)
	if err != nil {
		log.Printf("new secrets verifier err: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	r.Body = ioutil.NopCloser(io.TeeReader(r.Body, &verifier))
	s, err := slack.SlashCommandParse(r)
	if err != nil {
		log.Printf("error parsing slash command: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Printf("username:%s userID:%s command:%s", s.UserName, s.UserID, s.Command)

	if err = verifier.Ensure(); err != nil {
		log.Printf("error verifying authorization: %v", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	switch s.Command {
	case "/nowplaying":
		b, err := NowPlayingMessage()
		if err != nil {
			w.Header().Set(contentType, jsonType)
			b, _ = ErrorMessage("⛔️ Cannot get currently playing song.")
		}
		w.Header().Set(contentType, jsonType)
		_, _ = w.Write(b)
	case "/recentlyplayed":
		b, err := RecentlyPlayedMessage()
		if err != nil {
			w.Header().Set(contentType, jsonType)
			b, _ = ErrorMessage("⛔️ Cannot get recently played songs.")
		}
		w.Header().Set(contentType, jsonType)
		_, _ = w.Write(b)
	default:
		w.Header().Set(contentType, jsonType)
		b, _ := ErrorMessage("Unrecognized Command")
		_, _ = w.Write(b)
	}
}

func nowPlayingHandler(w http.ResponseWriter, r *http.Request) {
	if spotifyClient.Client == nil {
		_, _ = fmt.Fprint(w, "Please Log in to the BounceX Spotify Account")
		return
	}
	w.Header().Set(contentType, jsonType)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	nowPlaying, _ := spotifyClient.NowPlaying()
	var b []byte
	var err error
	if nowPlaying != nil {
		b, err = json.Marshal(nowPlaying)
	} else {
		w.WriteHeader(http.StatusNoContent)
		b, err = json.Marshal(Song{})
	}
	_, _ = w.Write(b)
	if err != nil {
		log.Print(err)
	}
}

func recentlyPlayedHandler(w http.ResponseWriter, r *http.Request) {
	if spotifyClient.Client == nil {
		_, _ = fmt.Fprint(w, "Please Log in to the BounceX Spotify Account")
		return
	}

	w.Header().Set(contentType, jsonType)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	recentlyPlayed, _ := spotifyClient.RecentlyPlayed()
	b, err := json.Marshal(recentlyPlayed)
	_, _ = w.Write(b)
	if err != nil {
		log.Print(err)
	}
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
	client := spotifyClient.Authenticator.NewClient(tok)
	user, _ := client.CurrentUser()
	w.Header().Set(contentType, textHtml)
	_, _ = fmt.Fprintf(w, "Login Completed!")
	spotifyClient.Channel <- &client
	SendLoginSuccessMessage(user)
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
