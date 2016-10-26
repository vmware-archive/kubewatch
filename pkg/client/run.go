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

	"github.com/skippbox/kubewatch/config"
	"github.com/skippbox/kubewatch/pkg/handlers"
	"github.com/skippbox/kubewatch/pkg/handlers/slack"
	"github.com/skippbox/kubewatch/pkg/controller"
)

// Run runs the event loop processing with given handler
func Run(conf *config.Config) {
	var eventHandler handlers.Handler
	//switch {
	//case len(conf.Handler.Slack.Channel) > 0 || len(conf.Handler.Slack.Token) > 0:
	//	eventHandler = new(slack.Slack)
	//default:
	//	eventHandler = new(handlers.Default)
	//}
	//TODO: temporary fix eventHandler = slack. Will add more later.
	eventHandler = new(slack.Slack)

	if err := eventHandler.Init(conf); err != nil {
		log.Fatal(err)
	}

	controller.Controller(conf, eventHandler)
}
