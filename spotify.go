package main

import (
	"bytes"
	"fmt"
	"github.com/zmb3/spotify"
	"log"
)

type SpotifyClient struct {
	Client        *spotify.Client
	Authenticator spotify.Authenticator
	State         string
	Channel       chan *spotify.Client
}

type Song struct {
	title  string
	artist string
	url    string
}

func (c *SpotifyClient) Login() {
	url := c.Authenticator.AuthURL(c.State)
	fmt.Println("Please log in to Spotify by visiting the following page in your browser:", url)

	// wait for auth to complete
	c.Client = <-c.Channel

	// use the client to make calls that require authorization
	user, err := c.Client.CurrentUser()
	if err != nil {
		log.Print(err)
	}
	fmt.Println("logged in as:", user.ID)
}

func (c *SpotifyClient) PrintNowPlaying() (*Song, error) {
	currentlyPlaying, err := c.Client.PlayerCurrentlyPlaying()
	var artists bytes.Buffer

	if err != nil {
		log.Print(err)
		return nil, err
	}

	if !currentlyPlaying.Playing || currentlyPlaying.Item == nil {
		return nil, nil
	}

	for i, artist := range currentlyPlaying.Item.Artists {
		if i > 0 {
			artists.WriteString(", ")
		}
		artists.WriteString(artist.Name)
	}
	song := &Song{
		currentlyPlaying.Item.Name,
		artists.String(),
		currentlyPlaying.Item.ExternalURLs["spotify"],
	}
	return song, nil
}
