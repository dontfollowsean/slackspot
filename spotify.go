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
	ID       string          `json:"id"`
	Title    string          `json:"title"`
	Artist   []*Artist       `json:"artist"`
	SongUrl  string          `json:"url"`
	Images   []spotify.Image `json:"images"`
	Progress int             `json:"progress"`
	Duration int             `json:"duration"`
}

type Artist struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Url  string `json:"url"`
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
		artists := getArtist(song.Track)
		s := &Song{
			ID:       song.Track.ID.String(),
			Title:    song.Track.Name,
			Artist:   artists,
			SongUrl:  song.Track.ExternalURLs["spotify"],
			Duration: song.Track.Duration,
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
	artists := getArtist(currentlyPlaying.Item.SimpleTrack)
	song := &Song{
		ID:       currentlyPlaying.Item.ID.String(),
		Title:    currentlyPlaying.Item.Name,
		Artist:   artists,
		SongUrl:  currentlyPlaying.Item.ExternalURLs["spotify"],
		Images:   currentlyPlaying.Item.Album.Images,
		Progress: currentlyPlaying.Progress,
		Duration: currentlyPlaying.Item.Duration,
	}
	return song, nil
}

func getArtistText(artists []*Artist) string {
	numOfArtists := len(artists)
	var artistText bytes.Buffer
	artistText.WriteString("by ")
	for i, artist := range artists {
		if i > 0 && numOfArtists > 2 {
			artistText.WriteString(", ")
		}
		if i > 0 && i == numOfArtists-1 {
			artistText.WriteString(" and ")
		}
		artistText.WriteString(fmt.Sprintf("<%s|%s>", artist.Url, artist.Name))
	}
	return artistText.String()
}

func getArtist(song spotify.SimpleTrack) []*Artist {
	var artists []*Artist
	for _, a := range song.Artists {
		artist := &Artist{
			ID:   a.ID.String(),
			Name: a.Name,
			Url:  a.ExternalURLs["spotify"],
		}
		artists = append(artists, artist)
	}
	return artists
}

func getImageUrl(images []spotify.Image, width int) string {
	var songImgUrl string
	for _, img := range images {
		if img.Width == width {
			songImgUrl = img.URL
		}
	}
	return songImgUrl
}
