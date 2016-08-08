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
	"flag"
	"log"
	"os"
	"os/signal"

	"k8s.io/kubernetes/pkg/api"
	cmdutil "k8s.io/kubernetes/pkg/kubectl/cmd/util"

	"github.com/skippbox/kubewatch/config"
	"github.com/skippbox/kubewatch/pkg/handlers"
)

var handlerFlag string

func init() {
	flag.StringVar(&handlerFlag, "handler", "default", "Handler for event, can be [slack, default], default handler is printing event")
}

// Run runs the event loop processing with given handler
func Run(conf *config.Config) {
	factory := cmdutil.NewFactory(nil)
	k8sClientConfig, err := factory.ClientConfig()
	if err != nil {
		log.Fatal(err)
	}

	client, err := New(conf, k8sClientConfig)
	if err != nil {
		log.Fatal(err)
	}

	h, ok := handlers.Map[handlerFlag]
	if !ok {
		log.Fatal("Handler not found")
	}

	eventHandler, ok := h.(handlers.Handler)
	if !ok {
		log.Fatal("Not an Handler type")
	}

	if err := eventHandler.Init(conf); err != nil {
		log.Fatal(err)
	}

	eventList, err := client.Events(api.NamespaceAll).List(api.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}

	watchEvents, err := client.Events(api.NamespaceAll).Watch(api.ListOptions{
		Watch:           true,
		ResourceVersion: eventList.ResourceVersion,
	})
	if err != nil {
		log.Fatal(err)
	}

	serviceList, err := client.Services(api.NamespaceAll).List(api.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}

	watchServices, err := client.Services(api.NamespaceAll).Watch(api.ListOptions{
		Watch:           true,
		ResourceVersion: serviceList.ResourceVersion,
	})
	if err != nil {
		log.Fatal(err)
	}

	client.waitGroup.Add(2)
	go client.EventLoop(watchEvents, eventHandler.Handle)
	go client.EventLoop(watchServices, eventHandler.Handle)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	defer signal.Stop(signals)

	log.Println("Press Ctrl+C to quit...")
	<-signals
	log.Println("Exiting...")

	client.Stop()
	log.Println("Exited normally.")
}
