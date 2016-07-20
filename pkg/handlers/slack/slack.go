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

package slack

import (
	"fmt"
	"log"
	"os"

	"github.com/nlopes/slack"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/watch"

	"github.com/skippbox/kubewatch/config"
)

var slackColors = map[string]string{
	"Normal":  "good",
	"Warning": "warning",
}

var slackErrMsg = `
%s

You need to set both slack token and channel for slack notify,
using "--slack-token" and "--slack-channel", or using environment variables:

export KW_SLACK_TOKEN=slack_token
export KW_SLACK_CHANNEL=slack_channel

Command line flags will override environment variables

`

// Slack handler implements handler.Handler interface,
// Notify event to slack channel
type Slack struct {
	Token   string
	Channel string
}

// Init prepares slack configuration
func (s *Slack) Init(c *config.Config) error {
	token := c.SlackToken
	channel := c.SlackChannel

	if token == "" {
		token = os.Getenv("KW_SLACK_TOKEN")
	}

	if channel == "" {
		channel = os.Getenv("KW_SLACK_CHANNEL")
	}

	s.Token = token
	s.Channel = channel

	return checkMissingSlackVars(s)
}

// Handle handles event for slack handler,
// send notify event to slack channel
func (s *Slack) Handle(e watch.Event) error {
	err := checkMissingSlackVars(s)
	if err != nil {
		return err
	}

	api := slack.New(s.Token)
	params := slack.PostMessageParameters{}
	attachment := prepareSlackAttachment(e)

	params.Attachments = []slack.Attachment{attachment}
	channelID, timestamp, err := api.PostMessage(s.Channel, "", params)
	if err != nil {
		log.Printf("%s\n", err)
		return err
	}

	log.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)
	return nil
}

func checkMissingSlackVars(s *Slack) error {
	if s.Token == "" || s.Channel == "" {
		return fmt.Errorf(slackErrMsg, "Missing slack token or channel")
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
