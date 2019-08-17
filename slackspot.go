package main

import (
	"fmt"
	"github.com/zmb3/spotify"
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
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		url := auth.AuthURL(state)
		fmt.Fprintf(w, "Spotify-Slack Integration\n Please log in to Spotify by visiting the following page in your browser: %s", url)
	})
	http.HandleFunc("/callback", completeAuth)
	http.HandleFunc("/nowplaying", nowPlayingHandler)
	//http.HandleFunc("/lastplayed", lastPlatedHandler)

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

func nowPlayingHandler(w http.ResponseWriter, r *http.Request) {
	currentlyPlaying, err := client.PlayerCurrentlyPlaying()
	if err != nil {
		log.Printf("err %s", err)
	}
	if !currentlyPlaying.Playing {
		fmt.Fprint(w, "No music is playing.")
	}
	fmt.Fprintf(w, currentlyPlaying.Item.Name)
}

func initAuth() {
	url := auth.AuthURL(state)
	fmt.Println("Please log in to Spotify by visiting the following page in your browser:", url)

	// wait for auth to complete
	client = <-ch

	// use the client to make calls that require authorization
	user, err := client.CurrentUser()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("You are logged in as:", user.ID)
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
