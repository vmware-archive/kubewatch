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
*/

package slackwebhook

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/bitnami-labs/kubewatch/config"
)

func TestWebhookInit(t *testing.T) {
	s := &SlackWebhook{}

	var Tests = []struct {
		slackwebhook config.SlackWebhook
		err          error
	}{
		{config.SlackWebhook{Channel: "foo", Username: "bar", Slackwebhookurl: "you"}, nil},
		{config.SlackWebhook{Channel: "foo"}, fmt.Errorf(webhookErrMsg, "Missing Slack Webhook Username")},
		{config.SlackWebhook{Username: "bar"}, fmt.Errorf(webhookErrMsg, "Missing Slack Webhook Channel")},
		{config.SlackWebhook{Emoji: ":kubernetes:"}, fmt.Errorf(webhookErrMsg, "Missing Slack Webhook Channel")},
		{config.SlackWebhook{Slackwebhookurl: "you"}, fmt.Errorf(webhookErrMsg, "Missing Slack Webhook Channel")},
		{config.SlackWebhook{}, fmt.Errorf(webhookErrMsg, "Missing Slack Webhook Channel")},
	}

	for _, tt := range Tests {
		c := &config.Config{}
		c.Handler.SlackWebhook = tt.slackwebhook
		if err := s.Init(c); !reflect.DeepEqual(err, tt.err) {
			t.Fatalf("Init(): %v", err)
		}
	}
}
