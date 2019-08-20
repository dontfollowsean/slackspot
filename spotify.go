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

//todo should  return song object and err
func (c *SpotifyClient) PrintNowPlaying() (SongTitle string, SongUrl string) {
	currentlyPlaying, err := c.Client.PlayerCurrentlyPlaying()
	var artists bytes.Buffer

	if err != nil {
		errMsg := fmt.Sprintf("Error getting current song: %s", err.Error())
		log.Print(errMsg)
		return errMsg, ""
	}

	if !currentlyPlaying.Playing || currentlyPlaying.Item == nil{
		return "No music is playing.", ""
	}

	for i, artist := range currentlyPlaying.Item.Artists {
		if i > 0 {
			artists.WriteString(", ")
		}
		artists.WriteString(artist.Name)
	}
	return fmt.Sprintf("%s by %s", currentlyPlaying.Item.Name, artists.String()), currentlyPlaying.Item.ExternalURLs["spotify"]
}
