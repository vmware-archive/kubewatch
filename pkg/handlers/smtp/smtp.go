/*
Copyright 2020 VMWare

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

/*
Package smtp implements an email notification handler for kubewatch.

See example configuration in the ConfigExample constant.
*/
package smtp

import (
	"fmt"
	"log"
	"time"

	"github.com/bitnami-labs/kubewatch/config"
	kbEvent "github.com/bitnami-labs/kubewatch/pkg/event"
	"github.com/sirupsen/logrus"
)

const (
	defaultSubject = "Kubewatch notification"

	// ConfigExample is an example configuration.
	ConfigExample = `handler:
  smtp:
    to: "myteam@mycompany.com"
    from: "kubewatch@mycluster.com"
    smarthost: smtp.mycompany.com:2525
    subject: Test notification
    auth:
      username: myusername
      password: mypassword
    requireTLS: true
`
)

// SMTP handler implements handler.Handler interface,
// Notify event via email.
type SMTP struct {
	cfg config.SMTP
}

// Init prepares Webhook configuration
func (s *SMTP) Init(c *config.Config) error {
	s.cfg = c.Handler.SMTP

	if s.cfg.To == "" {
		return fmt.Errorf("smtp `to` conf field is required")
	}
	if s.cfg.From == "" {
		return fmt.Errorf("smtp `from` conf field is required")
	}
	if s.cfg.Smarthost == "" {
		return fmt.Errorf("smtp `smarthost` conf field is required")
	}
	return nil
}

// ObjectCreated calls notifyWebhook on event creation
func (s *SMTP) ObjectCreated(obj interface{}) {
	notify(s, obj, "created")
}

// ObjectDeleted calls notifyWebhook on event creation
func (s *SMTP) ObjectDeleted(obj interface{}) {
	notify(s, obj, "deleted")
}

// ObjectUpdated calls notifyWebhook on event creation
func (s *SMTP) ObjectUpdated(oldObj, newObj interface{}) {
	notify(s, newObj, "updated")
}

// TestHandler tests the handler configurarion by sending test messages.
func (s *SMTP) TestHandler() {
	send(s.cfg, "test")
}

func notify(s *SMTP, obj interface{}, action string) {
	e := kbEvent.New(obj, action)
	msg, err := formatEmail(e, action)
	if err != nil {
		logrus.Error(err)
		return
	}
	send(s.cfg, msg)
	log.Printf("Message successfully sent to %s at %s ", s.cfg.To, time.Now())
}

func formatEmail(e kbEvent.Event, action string) (string, error) {
	return e.Message(), nil
}

func send(conf config.SMTP, msg string) {
	if err := sendEmail(conf, msg); err != nil {
		logrus.Error(err)
	}
}
