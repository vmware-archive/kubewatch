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

	"k8s.io/kubernetes/pkg/watch"

	kbEvent "github.com/skippbox/kubewatch/pkg/event"
)

// EventLoop process events in infinitive loop, apply handler function to each event
// Stop when receive interrupt signal
func (c *Client) EventLoop(w watch.Interface, handler func(kbEvent.Event) error) {
	defer c.waitGroup.Done()

	for {
		select {
		case event, ok := <-w.ResultChan():
			if !ok {
				return
			}
			e := kbEvent.New(event)
			if c.Filter(e) {
				if err := handler(e); err != nil {
					log.Println(err)
					w.Stop()
				}
			}
		case <-c.closeChan:
			log.Printf("Stopping watching events from %+v...\n", w)
			w.Stop()
			return
		}
	}
}

// Filter checks whether event matches configuration or not
func (c *Client) Filter(e kbEvent.Event) bool {
	reason := e.Reason

	if len(c.Config.Reason) == 0 {
		return true
	}

	for _, r := range c.Config.Reason {
		if r == reason || r == "*" {
			return true
		}
	}

	return false
}
