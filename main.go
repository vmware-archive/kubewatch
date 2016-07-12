package main

import (
	"log"

	"k8s.io/kubernetes/pkg/api"

	"github.com/runseb/kubewatch/pkg/client"
)

func main() {
	kubeWatchClient, err := client.New()
	w, err := kubeWatchClient.Events(api.NamespaceAll).Watch(api.ListOptions{Watch: true})
	if err != nil {
		log.Fatal(err)
	}

	kubeWatchClient.EventLoop(w, client.NotifySlack)
}
