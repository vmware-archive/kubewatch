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
	"os"
	"os/signal"
	"syscall"

	"github.com/bitnami-labs/kubewatch/config"
	"github.com/bitnami-labs/kubewatch/pkg/controller"
	"github.com/bitnami-labs/kubewatch/pkg/handlers"
	"github.com/bitnami-labs/kubewatch/pkg/handlers/flock"
	"github.com/bitnami-labs/kubewatch/pkg/handlers/hipchat"
	"github.com/bitnami-labs/kubewatch/pkg/handlers/mattermost"
	"github.com/bitnami-labs/kubewatch/pkg/handlers/slack"
	"github.com/bitnami-labs/kubewatch/pkg/handlers/webhook"
)

// Run runs the event loop processing with given handler
func Run(conf *config.Config) {
	stopCh := make(chan struct{})
	defer close(stopCh)

	var eventHandler []handlers.Handler
	if len(conf.Handler.Slack.Channel) > 0 || len(conf.Handler.Slack.Token) > 0 {
		eventHandler = append(eventHandler, new(slack.Slack))
	}
	if len(conf.Handler.Hipchat.Room) > 0 || len(conf.Handler.Hipchat.Token) > 0 {
		eventHandler = append(eventHandler, new(hipchat.Hipchat))
	}
	if len(conf.Handler.Mattermost.Channel) > 0 || len(conf.Handler.Mattermost.Url) > 0 {
		eventHandler = append(eventHandler, new(mattermost.Mattermost))
	}
	if len(conf.Handler.Flock.Url) > 0 {
		eventHandler = append(eventHandler, new(flock.Flock))
	}
	if len(conf.Handler.Webhook.Url) > 0 {
		eventHandler = append(eventHandler, new(webhook.Webhook))
	}

	if len(eventHandler) == 0 {
		eventHandler = append(eventHandler, new(handlers.Default))
	}

	log.Printf("%v\n", eventHandler)
	for _, handler := range eventHandler {
		controller.Start(conf, handler, stopCh)
	}

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGTERM)
	signal.Notify(sigterm, syscall.SIGINT)
	<-sigterm
}
