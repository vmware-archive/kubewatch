package client

import (
	"encoding/json"
	"errors"
	"log"
	"os"

	"github.com/nlopes/slack"
	"k8s.io/kubernetes/pkg/watch"
)

// NotifySlack sends event to slack channel
func NotifySlack(e watch.Event) error {
	slackAPIToken := os.Getenv("KW_SLACK_TOKEN")
	if slackAPIToken == "" {
		return errors.New("Missing slack token!")
	}

	slackChannel := os.Getenv("KW_SLACK_CHANNEL")
	if slackChannel == "" {
		return errors.New("Missing slack channel!")
	}

	api := slack.New(slackAPIToken)
	params := slack.PostMessageParameters{}
	attachment := slack.Attachment{
		Text: eventFormatter(e),
	}
	params.Attachments = []slack.Attachment{attachment}
	channelID, timestamp, err := api.PostMessage(slackChannel, "", params)
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
