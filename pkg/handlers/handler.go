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

package handlers

import (
	"github.com/bitnami-labs/kubewatch/config"
	"github.com/bitnami-labs/kubewatch/pkg/handlers/flock"
	"github.com/bitnami-labs/kubewatch/pkg/handlers/hipchat"
	"github.com/bitnami-labs/kubewatch/pkg/handlers/mattermost"
	"github.com/bitnami-labs/kubewatch/pkg/handlers/msteam"
	"github.com/bitnami-labs/kubewatch/pkg/handlers/slack"
	"github.com/bitnami-labs/kubewatch/pkg/handlers/webhook"
)

// Handler is implemented by any handler.
// The Handle method is used to process event
type Handler interface {
	Init(c *config.Config) error
	ObjectCreated(obj interface{})
	ObjectDeleted(obj interface{})
	ObjectUpdated(oldObj, newObj interface{})
	TestHandler()
}

// Map maps each event handler function to a name for easily lookup
var Map = map[string]interface{}{
	"default":    &Default{},
	"slack":      &slack.Slack{},
	"hipchat":    &hipchat.Hipchat{},
	"mattermost": &mattermost.Mattermost{},
	"flock":      &flock.Flock{},
	"webhook":    &webhook.Webhook{},
	"ms-teams":   &msteam.MSTeams{},
}

// Default handler implements Handler interface,
// print each event with JSON format
type Default struct {
}

// Init initializes handler configuration
// Do nothing for default handler
func (d *Default) Init(c *config.Config) error {
	return nil
}

// ObjectCreated sends events on object creation
func (d *Default) ObjectCreated(obj interface{}) {

}

// ObjectDeleted sends events on object deletion
func (d *Default) ObjectDeleted(obj interface{}) {

}

// ObjectUpdated sends events on object updation
func (d *Default) ObjectUpdated(oldObj, newObj interface{}) {

}

// TestHandler tests the handler configurarion by sending test messages.
func (d *Default) TestHandler() {

}
