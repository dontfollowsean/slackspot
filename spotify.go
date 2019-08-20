package main

import (
	"bytes"
	"fmt"
	"github.com/zmb3/spotify"
	"log"
)

func PrintNowPlaying(client *spotify.Client) string {
	currentlyPlaying, err := client.PlayerCurrentlyPlaying()
	var artists bytes.Buffer

	if err != nil {
		errMsg := fmt.Sprintf("Error getting current song: %s", err.Error())
		log.Print(errMsg)
		return errMsg
	}

	if !currentlyPlaying.Playing || currentlyPlaying.Item == nil{
		return "No music is playing."
	}

	for i, artist := range currentlyPlaying.Item.Artists {
		if i > 0 {
			artists.WriteString(", ")
		}
		artists.WriteString(artist.Name)
	}
	return fmt.Sprintf("%s by %s", currentlyPlaying.Item.Name, artists.String())
}
