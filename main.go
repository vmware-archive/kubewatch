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
Missing slack token or slack channel

Specify by "--slack-token" and "--slack-channel", or using environment variables:

KW_SLACK_TOKEN=slack_token
KW_SLACK_CHANNEL=slack_channel
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
		for _, s := range []string{slackToken, slackChannel} {
			if s == "" {
				log.Fatal(slackErrMsg)
			}
		}
		client.SlackToken = slackToken
		client.SlackChannel = slackChannel
	}

	kubeWatchClient, err := client.New()
	w, err := kubeWatchClient.Events(api.NamespaceAll).Watch(api.ListOptions{Watch: true})
	if err != nil {
		log.Fatal(err)
	}

	kubeWatchClient.EventLoop(w, eventHandler)
}
