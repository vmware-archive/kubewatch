/*
Copyright 2018 Bitnami

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

/*
Author: Richard Knechtel
Company: Blast Motion
Date: 12/27/2021

Info:
Example Message send to slack webhook (In Python):

  # Example Webhook URL:
  url = "https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX""
  # Message to Slack
  msg = {
      "channel": "#my-alerts",
      "username": "Webhook_Username",
      "text": "Pod startup failed",
      "icon_emoji": ""
  }
  msg = json.dumps(msg).encode('utf-8')
  resp = http.request('POST',url, body=msg)
*/

package slackwebhook

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/slack-go/slack"

	"github.com/bitnami-labs/kubewatch/config"
	"github.com/bitnami-labs/kubewatch/pkg/event"
)

var webhookErrMsg = `
%s

You need to set Webhook url
using "--channel/-c, --username/-n, --emoji/-e, --slackwebhookurl/-u" or using environment variables:

export KW_SLACK_CHANNEL=slack_channel
export KW_SLACK_USERNAME=slack_username
export KW_SLACK_EMOJI=slack_emoji
export KW_SLACK_WEBHOOK_URL=slack_webhook_url

Command line flags will override environment variables

`

// Webhook handler implements handler.Handler interface,
// Notify event to Webhook channel
type SlackWebhook struct {
	Channel         string
	Username        string
	Emoji           string
	Slackwebhookurl string
}

// Init prepares Webhook configuration
func (m *SlackWebhook) Init(c *config.Config) error {
	channel := c.Handler.SlackWebhook.Channel
	username := c.Handler.SlackWebhook.Username
	emoji := c.Handler.SlackWebhook.Emoji
	slackwebhookurl := c.Handler.SlackWebhook.Slackwebhookurl

	if channel == "" {
		channel = os.Getenv("KW_SLACK_CHANNEL")
	}
	if username == "" {
		username = os.Getenv("KW_SLACK_USERNAME")
	}
	if emoji == "" {
		emoji = os.Getenv("KW_SLACK_EMOJI")
	}
	if slackwebhookurl == "" {
		slackwebhookurl = os.Getenv("KW_SLACK_WEBHOOK_URL")
	}

	m.Channel = channel
	m.Username = username
	m.Emoji = emoji
	m.Slackwebhookurl = slackwebhookurl

	return checkMissingWebhookVars(m)
}

// Handle handles an event.
func (m *SlackWebhook) Handle(e event.Event) {

	webhookMessage := slack.WebhookMessage{
		Channel:   m.Channel,
		Username:  m.Username,
		Text:      e.Message(),
		IconEmoji: m.Emoji,
	}

	log.Printf("slackwebhook-handle():Slackwebhook WebHookMessage: %s", webhookMessage.Text)

	err := slack.PostWebhook(m.Slackwebhookurl, &webhookMessage)

	if err != nil {
		log.Printf("slackwebhook-handle() Error: %s\n", err)
		return
	}

	log.Printf("Message successfully sent to %s at %s. Message: %s", m.Slackwebhookurl, time.Now(), webhookMessage.Text)
}

func checkMissingWebhookVars(s *SlackWebhook) error {
	if s.Channel == "" {
		return fmt.Errorf(webhookErrMsg, "Missing Slack Webhook Channel")
	}
	if s.Username == "" {
		return fmt.Errorf(webhookErrMsg, "Missing Slack Webhook Username")
	}
	if s.Slackwebhookurl == "" {
		return fmt.Errorf(webhookErrMsg, "Missing Slack Webhook url")
	}

	return nil
}
