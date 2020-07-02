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
	"log"

	"github.com/bitnami-labs/kubewatch/config"
	"github.com/bitnami-labs/kubewatch/pkg/controller"
	"github.com/bitnami-labs/kubewatch/pkg/handlers"
	"github.com/bitnami-labs/kubewatch/pkg/handlers/flock"
	"github.com/bitnami-labs/kubewatch/pkg/handlers/hipchat"
	"github.com/bitnami-labs/kubewatch/pkg/handlers/mattermost"
	"github.com/bitnami-labs/kubewatch/pkg/handlers/msteam"
	"github.com/bitnami-labs/kubewatch/pkg/handlers/slack"
	"github.com/bitnami-labs/kubewatch/pkg/handlers/webhook"
)

// Run runs the event loop processing with given handler
func Run(conf *config.Config) {

	var eventHandler = ParseEventHandler(conf)
	controller.Start(conf, eventHandler)
}

// ParseEventHandler returns the respective handler object specified in the config file.
func ParseEventHandler(conf *config.Config) handlers.Handler {

	var eventHandler handlers.Handler
	switch {
	case len(conf.Handler.Slack.Channel) > 0 || len(conf.Handler.Slack.Token) > 0:
		eventHandler = new(slack.Slack)
	case len(conf.Handler.Hipchat.Room) > 0 || len(conf.Handler.Hipchat.Token) > 0:
		eventHandler = new(hipchat.Hipchat)
	case len(conf.Handler.Mattermost.Channel) > 0 || len(conf.Handler.Mattermost.Url) > 0:
		eventHandler = new(mattermost.Mattermost)
	case len(conf.Handler.Flock.Url) > 0:
		eventHandler = new(flock.Flock)
	case len(conf.Handler.Webhook.Url) > 0:
		eventHandler = new(webhook.Webhook)
	case len(conf.Handler.MSTeams.WebhookURL) > 0:
		eventHandler = new(msteam.MSTeams)
	default:
		eventHandler = new(handlers.Default)
	}
	if err := eventHandler.Init(conf); err != nil {
		log.Fatal(err)
	}
	return eventHandler
}
