package main

import (
	"encoding/json"
	"fmt"
	"github.com/nlopes/slack"
	"log"
)

func NowPlayingMessage() ([]byte, error) {
	if spotifyClient.Client == nil {
		return ErrorMessage("Please ask an Admin to log into Spotify")
	}
	song, err := spotifyClient.NowPlaying()
	if err != nil {
		log.Print(err)
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
				Title:     fmt.Sprintf("%s by %s", song.title, song.artist),
				TitleLink: song.url,
			},
		}
		slackMsg = &slack.Msg{
			Text:        "🎵 Now playing...",
			Attachments: attachments,
		}
	}
	b, err := json.Marshal(slackMsg)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return b, nil
}

func RecentlyPlayedMessage() ([]byte, error) {
	if spotifyClient.Client == nil {
		return ErrorMessage("Please ask an Admin to log into Spotify")
	}
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
	b, err := toJsonBody(slackMsg)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func ErrorMessage(errorMessage string) ([]byte, error) {
	//contactUser := os.Getenv("ContactUser")
	slackMsg := &slack.Msg{
		Text: errorMessage,
		Attachments: []slack.Attachment{
			{Title: fmt.Sprint("If this error persists contact <@U1055Q4A0|sean>")},//(bx:<@UD0NXF3UY|sean>)
		},
	}
	b, err := toJsonBody(slackMsg)
	if err != nil {
		return b, err
	}
	return b, nil
}

func toJsonBody(slackMsg *slack.Msg) ([]byte, error) {
	b, err := json.Marshal(slackMsg)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return b, nil
}

func SendLoginMessage(url string) {
	webhook := "https://hooks.slack.com/services/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"
	slackMsg := &slack.WebhookMessage{
		Attachments: []slack.Attachment{
			{
				Title:     "Log into Spotify here",
				TitleLink: url,
			},
		},
	}
	//b, err := toJsonBody(slackMsg)
	//if err != nil {
	//	log.Print(err)
	//}
	err := slack.PostWebhook(webhook, slackMsg)
	if err != nil {
		log.Print(err)
	}
}
