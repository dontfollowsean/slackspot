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
	Title  string          `json:"title"`
	Artist string          `json:"artist"`
	Url    string          `json:"url"`
	Images []spotify.Image `json:"images"`
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
		log.Printf("error getting current user: %s", err)
	}
	fmt.Println("logged in as:", user.DisplayName)
}

func (c *SpotifyClient) RecentlyPlayed() ([]*Song, error) {
	recentlyPlayed, err := c.Client.PlayerRecentlyPlayed();
	if err != nil {
		log.Printf("error getting recently played songs: %s", err)
		return nil, err
	}
	songs := make([]*Song, songHistoryLength)
	for i, song := range recentlyPlayed[:songHistoryLength] {
		s := &Song{
			Title:  song.Track.Name,
			Artist: getArtists(song.Track),
			Url:    song.Track.ExternalURLs["spotify"],
		}
		fullTrack, err := c.Client.GetTrack(song.Track.ID)
		if err == nil {
			s.Images = fullTrack.Album.Images
		}
		songs[i] = s
	}
	return songs, nil
}

func (c *SpotifyClient) NowPlaying() (*Song, error) {
	currentlyPlaying, err := c.Client.PlayerCurrentlyPlaying()
	if err != nil {
		log.Printf("error getting currently playing song: %s", err)
		if err.Error() == "EOF" {
			return nil, nil // there's no music playing
		}
		return nil, err
	}

	if !currentlyPlaying.Playing || currentlyPlaying.Item == nil {
		return nil, nil
	}

	song := &Song{
		currentlyPlaying.Item.Name,
		getArtists(currentlyPlaying.Item.SimpleTrack),
		currentlyPlaying.Item.ExternalURLs["spotify"],
		currentlyPlaying.Item.Album.Images,
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
