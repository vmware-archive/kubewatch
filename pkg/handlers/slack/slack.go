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

	"github.com/nlopes/slack"

	"github.com/skippbox/kubewatch/config"
	"github.com/skippbox/kubewatch/pkg/event"
	"github.com/skippbox/kubewatch/pkg/handlers"
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

export SLACK_TOKEN=slack_token
export SLACK_CHANNEL=slack_channel

Command line flags will override environment variables

`

// Slack is the underlying struct used by the slack handler receivers
type Slack struct {
	token   string
	channel string
	config  *config.Config
}

// New returns a slack handler interface
func New(conf *config.Config, channel string, token string) handlers.Handler {
	c := Slack{
		token:   token,
		channel: channel,
		config:  conf,
	}
	handler := handlers.Handler(&c)
	return handler
}

// Config returns the config data that will be used by the handler
func (s *Slack) Config() *config.Config {
	return s.config
}

// Init prepares slack configuration
func (s *Slack) Init() error {
	fmt.Println(s.channel, s.token)
	return checkMissingSlackVars(s)
}

func (s *Slack) ObjectCreated(obj interface{}) {
	notifySlack(s, obj, "created")
}

func (s *Slack) ObjectDeleted(obj interface{}) {
	notifySlack(s, obj, "deleted")
}

func (s *Slack) ObjectUpdated(oldObj, newObj interface{}) {
	notifySlack(s, newObj, "updated")
}

func notifySlack(s *Slack, obj interface{}, action string) {
	e := event.New(obj, action)
	api := slack.New(s.token)
	params := slack.PostMessageParameters{}
	attachment := prepareSlackAttachment(e)

	params.Attachments = []slack.Attachment{attachment}
	params.AsUser = true
	channelID, timestamp, err := api.PostMessage(s.channel, "", params)
	if err != nil {
		log.Printf("%s\n", err)
		return
	}

	log.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)
}

func checkMissingSlackVars(s *Slack) error {
	if s.token == "" || s.channel == "" {
		return fmt.Errorf(slackErrMsg, "Missing slack token or channel")
	}

	return nil
}

func prepareSlackAttachment(e event.Event) slack.Attachment {
	msg := fmt.Sprintf(
		"A %s in namespace %s has been %s: %s",
		e.Kind,
		e.Namespace,
		e.Reason,
		e.Name,
	)

	attachment := slack.Attachment{
		Fields: []slack.AttachmentField{
			slack.AttachmentField{
				Title: "kubewatch",
				Value: msg,
			},
		},
	}

	if color, ok := slackColors[e.Status]; ok {
		attachment.Color = color
	}

	attachment.MarkdownIn = []string{"fields"}

	return attachment
}
