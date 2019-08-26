package main

import (
	"encoding/json"
	"fmt"
	"github.com/nlopes/slack"
	"log"
)

func NowPlayingMessage() ([]byte, error) {
	if spotifyClient.Client == nil {
		return ErrorMessage("Please ask an Admin to log in to Spotify")
	}
	song, err := spotifyClient.NowPlaying()
	if err != nil {
		return nil, err
	}

	var slackMsg *slack.Msg
	if song == nil {
		slackMsg = &slack.Msg{
			Text: "There is no music playing.",
		}
	} else {
		attachments := []slack.Attachment{
			{
				Title:     fmt.Sprintf("%s by %s", song.Title, song.Artist),
				TitleLink: song.Url,
			},
		}
		slackMsg = &slack.Msg{
			Text:        "🎵 Now playing...",
			Attachments: attachments,
		}
	}
	return toJsonBody(slackMsg)
}

func RecentlyPlayedMessage() ([]byte, error) {
	if spotifyClient.Client == nil {
		return ErrorMessage("Please ask an Admin to log in to Spotify")
	}
	songs, err := spotifyClient.RecentlyPlayed()
	if err != nil {
		return nil, err
	}
	attachments := make([]slack.Attachment, 0)
	for _, song := range songs {
		attachment := slack.Attachment{
			Title:     fmt.Sprintf("%s by %s", song.Title, song.Artist),
			TitleLink: song.Url,
		}
		attachments = append(attachments, attachment)
	}
	slackMsg := &slack.Msg{
		Text:        "🎵 Recently Played Songs",
		Attachments: attachments,
	}
	return toJsonBody(slackMsg)
}

func ErrorMessage(errorMessage string) ([]byte, error) {
	slackMsg := &slack.Msg{
		Text: fmt.Sprintf("%s If this error persists contact %s", errorMessage, contactUser),
	}
	return toJsonBody(slackMsg)
}

func toJsonBody(slackMsg *slack.Msg) ([]byte, error) {
	b, err := json.Marshal(slackMsg)
	if err != nil {
		log.Printf("error marshalling to json: %s", err)
		return nil, err
	}
	return b, nil
}

func SendLoginMessage(url string) {
	webhook := getEnv("SLACK_ADMIN_WEBHOOK", "")
	slackMsg := &slack.WebhookMessage{
		Attachments: []slack.Attachment{
			{
				Title:     "Log into Spotify here",
				TitleLink: url,
			},
		},
	}
	err := slack.PostWebhook(webhook, slackMsg)
	if err != nil {
		log.Printf("error posting to webhook: %s", err)
	}
}
