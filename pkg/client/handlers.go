package client

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/nlopes/slack"
	"k8s.io/kubernetes/pkg/watch"
)

// Handlers map each event handler function to a name for easily lookup
var Handlers = map[string]func(w watch.Event) error{
	"default": PrintEvent,
	"slack":   NotifySlack,
}

var (
	// SlackToken used by slack handler, for slack authentication
	SlackToken string
	// SlackChannel used by slack handler, specify where the message sent to
	SlackChannel string
)

// NotifySlack sends event to slack channel
func NotifySlack(e watch.Event) error {
	if SlackToken == "" {
		return errors.New("Missing slack token!")
	}

	if SlackChannel == "" {
		return errors.New("Missing slack channel!")
	}

	api := slack.New(SlackToken)
	params := slack.PostMessageParameters{}
	attachment := slack.Attachment{
		Text: eventFormatter(e),
	}
	params.Attachments = []slack.Attachment{attachment}
	channelID, timestamp, err := api.PostMessage(SlackChannel, "", params)
	if err != nil {
		log.Printf("%s\n", err)
		return err
	}

	log.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)
	return nil
}

// PrintEvent print event in json format, for testing or debugging
func PrintEvent(e watch.Event) error {
	b, err := json.Marshal(e)
	if err != nil {
		return err
	}

	log.Println(string(b))

	return nil
}
