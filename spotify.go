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
	Title  string `json:"title"`
	Artist string `json:"artist"`
	Url    string `json:"url"`
}

func (c *SpotifyClient) Login() {
	url := c.Authenticator.AuthURL(c.State)
	fmt.Println("Please log in to Spotify by visiting the following page in your browser:", url)
	SendLoginMessage(url)
	// wait for auth to complete
	c.Client = <-c.Channel

	// use the client to make calls that require authorization
	user, err := c.Client.CurrentUser()
	if err != nil {
		log.Print(err)
	}
	fmt.Println("logged in as:", user.ID)
}

func (c *SpotifyClient) RecentlyPlayed() ([]*Song, error) {
	recentlyPlayed, err := c.Client.PlayerRecentlyPlayed();
	if err != nil {
		return nil, err
	}
	songs := make([]*Song, songHistoryLength)
	for i, song := range recentlyPlayed[:songHistoryLength] {
		s := &Song{
			Title:  song.Track.Name,
			Artist: getArtists(song.Track),
			Url:    song.Track.ExternalURLs["spotify"],
		}
		songs[i] = s
	}
	return songs, nil
}

func (c *SpotifyClient) NowPlaying() (*Song, error) {
	currentlyPlaying, err := c.Client.PlayerCurrentlyPlaying()

	if err != nil {
		return nil, err
	}

	if !currentlyPlaying.Playing || currentlyPlaying.Item == nil {
		return nil, nil
	}

	song := &Song{
		currentlyPlaying.Item.Name,
		getArtists(currentlyPlaying.Item.SimpleTrack),
		currentlyPlaying.Item.ExternalURLs["spotify"],
	}
	return song, nil
}

func getArtists(song spotify.SimpleTrack) string {
	var artists bytes.Buffer
	for i, artist := range song.Artists {
		if i > 0 {
			artists.WriteString(", ")
		}
		artists.WriteString(artist.Name)
	}
	return artists.String()
}
