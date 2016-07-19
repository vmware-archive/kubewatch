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

package main

import (
	"flag"
	"log"

	"k8s.io/kubernetes/pkg/api"

	"github.com/skippbox/kubewatch/pkg/client"
	"github.com/skippbox/kubewatch/pkg/handlers"
)

var (
	handlerFlag  string
	slackToken   string
	slackChannel string
)

var slackErrMsg = `
%s

You need to set both slack token and channel for slack notify,
using "--slack-token" and "--slack-channel", or using environment variables:

export KW_SLACK_TOKEN=slack_token
export KW_SLACK_CHANNEL=slack_channel

Command line flags will override environment variables

`

func init() {
	flag.StringVar(&handlerFlag, "handler", "default", "Handler for event, can be [slack, default], default handler is printing event")
	flag.StringVar(&slackToken, "slack-token", "", "Slack token")
	flag.StringVar(&slackChannel, "slack-channel", "", "Slack channel")
}

func main() {
	flag.Parse()

	h, ok := handlers.Map[handlerFlag]
	if !ok {
		log.Fatal("Handler not found")
	}

	eventHandler, ok := h.(handlers.Handler)
	if !ok {
		log.Fatal("Not an Handler type")
	}

	if err := eventHandler.Init(slackToken, slackChannel); err != nil {
		log.Fatalf(slackErrMsg, err)
	}

	kubeWatchClient, err := client.New()
	w, err := kubeWatchClient.Events(api.NamespaceAll).Watch(api.ListOptions{Watch: true})
	if err != nil {
		log.Fatal(err)
	}

	kubeWatchClient.EventLoop(w, eventHandler.Handle)
}
