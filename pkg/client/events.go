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

	"k8s.io/kubernetes/pkg/watch"
)

// EventLoop process events in infinitive loop, apply handler function to each event
// Stop when receive interrupt signal
func (c *Client) EventLoop(w watch.Interface, handler func(watch.Event) error) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	defer signal.Stop(signals)

	for {
		select {
		case event, ok := <-w.ResultChan():
			if !ok {
				return
			}
			if err := handler(event); err != nil {
				log.Println(err)
				w.Stop()
			}
		case <-signals:
			log.Println("Catched signal, quit normally.")
			w.Stop()
		}
	}

}
