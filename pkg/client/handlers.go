/*
Copyright 2016 Skippbox, Ltd.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/nlopes/slack"
	"k8s.io/kubernetes/pkg/api"
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

	slackColors = map[string]string{
		"Normal":  "good",
		"Warning": "warning",
	}
)

// InitSlack prepares slack required variables
func InitSlack(token, channel string) error {
	if token == "" {
		token = os.Getenv("KW_SLACK_TOKEN")
	}

	if channel == "" {
		channel = os.Getenv("KW_SLACK_CHANNEL")
	}

	SlackToken = token
	SlackChannel = channel

	return checkMissingSlackVars()
}

// NotifySlack sends event to slack channel
func NotifySlack(e watch.Event) error {
	err := checkMissingSlackVars()
	if err != nil {
		return err
	}

	api := slack.New(SlackToken)
	params := slack.PostMessageParameters{}
	attachment := prepareSlackAttachment(e)

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

func checkMissingSlackVars() error {
	for k, v := range map[string]string{"token": SlackToken, "channel": SlackChannel} {
		if v == "" {
			errMsg := fmt.Sprintf("Missing slack %s!", k)
			return errors.New(errMsg)
		}
	}

	return nil
}

func prepareSlackAttachment(e watch.Event) slack.Attachment {
	apiEvent := (e.Object).(*api.Event)

	msg := fmt.Sprintf(
		"In *Namespace* %s *Kind* %s from *Component* %s on *Host* %s had *Reason* %s",
		apiEvent.ObjectMeta.Namespace,
		apiEvent.InvolvedObject.Kind,
		apiEvent.Source.Component,
		apiEvent.Source.Host,
		apiEvent.Reason,
	)
	attachment := slack.Attachment{
		Fields: []slack.AttachmentField{
			slack.AttachmentField{
				Title: "kubewatch",
				Value: msg,
			},
		},
	}

	if color, ok := slackColors[apiEvent.Type]; ok {
		attachment.Color = color
	}

	attachment.MarkdownIn = []string{"fields"}

	return attachment
}
