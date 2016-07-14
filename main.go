package main

import (
	"flag"
	"log"

	"k8s.io/kubernetes/pkg/api"

	"github.com/skippbox/kubewatch/pkg/client"
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

	eventHandler, ok := client.Handlers[handlerFlag]
	if !ok {
		log.Fatal("Handler not found")
	}

	if handlerFlag == "slack" {
		if err := client.InitSlack(slackToken, slackChannel); err != nil {
			log.Fatalf(slackErrMsg, err)
		}
	}

	kubeWatchClient, err := client.New()
	w, err := kubeWatchClient.Events(api.NamespaceAll).Watch(api.ListOptions{Watch: true})
	if err != nil {
		log.Fatal(err)
	}

	kubeWatchClient.EventLoop(w, eventHandler)
}
