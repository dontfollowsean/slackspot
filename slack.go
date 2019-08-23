package main

import (
	"encoding/json"
	"fmt"
	"github.com/nlopes/slack"
	"log"
)

func NowPlayingMessage() ([]byte, error) {
	song, err := spotifyClient.NowPlaying()
	if err != nil {
		log.Print(err)
		return nil, err
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
		return nil, err
	}
	return b, nil
}

func RecentlyPlayedMessage() ([]byte, error) {
	songs, err := spotifyClient.RecentlyPlayed()
	if err != nil {
		log.Print(err)
		return nil, err
	}
	attachments := make([]slack.Attachment, 0)
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
		return nil, err
	}
	return b, nil
}