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

	"github.com/bitnami-labs/kubewatch/config"
	"github.com/bitnami-labs/kubewatch/pkg/event"
	kbEvent "github.com/bitnami-labs/kubewatch/pkg/event"
)

var slackColors = map[string]string{
	"Normal":  "good",
	"Warning": "warning",
	"Danger":  "danger",
}

var slackErrMsg = `
%s

You need to set both slack token and channel for slack notify,
using "--token/-t" and "--channel/-c", or using environment variables:

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
	token := c.Handler.Slack.Token
	channel := c.Handler.Slack.Channel

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

// ObjectCreated calls notifySlack on event creation
func (s *Slack) ObjectCreated(obj interface{}) {
	notifySlack(s, obj, "created")
}

// ObjectDeleted calls notifySlack on event creation
func (s *Slack) ObjectDeleted(obj interface{}) {
	notifySlack(s, obj, "deleted")
}

// ObjectUpdated calls notifySlack on event creation
func (s *Slack) ObjectUpdated(oldObj, newObj interface{}) {
	notifySlack(s, newObj, "updated")
}

// TestHandler tests the handler configurarion by sending test messages.
func (s *Slack) TestHandler() {
	api := slack.New(s.Token)
	params := slack.PostMessageParameters{}
	attachment := slack.Attachment{
		Fields: []slack.AttachmentField{
			{
				Title: "kubewatch",
				Value: "Testing Handler Configuration. This is a Test message.",
			},
		},
	}
	params.Attachments = []slack.Attachment{attachment}
	params.AsUser = true
	channelID, timestamp, err := api.PostMessage(s.Channel, "", params)
	if err != nil {
		log.Printf("%s\n", err)
		return
	}

	log.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)
}

func notifySlack(s *Slack, obj interface{}, action string) {
	e := kbEvent.New(obj, action)
	api := slack.New(s.Token)
	params := slack.PostMessageParameters{}
	attachment := prepareSlackAttachment(e)

	params.Attachments = []slack.Attachment{attachment}
	params.AsUser = true
	channelID, timestamp, err := api.PostMessage(s.Channel, "", params)
	if err != nil {
		log.Printf("%s\n", err)
		return
	}

	log.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)
}

func checkMissingSlackVars(s *Slack) error {
	if s.Token == "" || s.Channel == "" {
		return fmt.Errorf(slackErrMsg, "Missing slack token or channel")
	}

	return nil
}

func prepareSlackAttachment(e event.Event) slack.Attachment {

	attachment := slack.Attachment{
		Fields: []slack.AttachmentField{
			{
				Title: "kubewatch",
				Value: e.Message(),
			},
		},
	}

	if color, ok := slackColors[e.Status]; ok {
		attachment.Color = color
	}

	attachment.MarkdownIn = []string{"fields"}

	return attachment
}
