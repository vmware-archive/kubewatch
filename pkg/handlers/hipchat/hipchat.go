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

package hipchat

import (
	"fmt"
	"log"
	"os"

	hipchat "github.com/tbruyelle/hipchat-go/hipchat"

	"net/url"

	"github.com/bitnami-labs/kubewatch/config"
	"github.com/bitnami-labs/kubewatch/pkg/event"
)

var hipchatColors = map[string]hipchat.Color{
	"Normal":  hipchat.ColorGreen,
	"Warning": hipchat.ColorYellow,
	"Danger":  hipchat.ColorRed,
}

var hipchatErrMsg = `
%s

You need to set both hipchat token and room for hipchat notify,
using "--token/-t", "--room/-r", and "--url/-u" or using environment variables:

export KW_HIPCHAT_TOKEN=hipchat_token
export KW_HIPCHAT_ROOM=hipchat_room
export KW_HIPCHAT_URL=hipchat_url (defaults to https://api.hipchat.com/v2)

Command line flags will override environment variables

`

// Hipchat handler implements handler.Handler interface,
// Notify event to hipchat room
type Hipchat struct {
	Token string
	Room  string
	Url   string
}

// Init prepares hipchat configuration
func (s *Hipchat) Init(c *config.Config) error {
	url := c.Handler.Hipchat.Url
	room := c.Handler.Hipchat.Room
	token := c.Handler.Hipchat.Token

	if token == "" {
		token = os.Getenv("KW_HIPCHAT_TOKEN")
	}

	if room == "" {
		room = os.Getenv("KW_HIPCHAT_ROOM")
	}

	if url == "" {
		url = os.Getenv("KW_HIPCHAT_URL")
	}

	s.Token = token
	s.Room = room
	s.Url = url

	return checkMissingHipchatVars(s)
}

// Handle handles the notification.
func (s *Hipchat) Handle(e event.Event) {
	client := hipchat.NewClient(s.Token)
	if s.Url != "" {
		baseUrl, err := url.Parse(s.Url)
		if err != nil {
			panic(err)
		}
		client.BaseURL = baseUrl
	}

	notificationRequest := prepareHipchatNotification(e)
	_, err := client.Room.Notification(s.Room, &notificationRequest)

	if err != nil {
		log.Printf("%s\n", err)
		return
	}

	log.Printf("Message successfully sent to room %s", s.Room)
}

func checkMissingHipchatVars(s *Hipchat) error {
	if s.Token == "" || s.Room == "" {
		return fmt.Errorf(hipchatErrMsg, "Missing hipchat token or room")
	}

	return nil
}

func prepareHipchatNotification(e event.Event) hipchat.NotificationRequest {
	notification := hipchat.NotificationRequest{
		Message: e.Message(),
		Notify:  true,
		From:    "kubewatch",
	}

	if color, ok := hipchatColors[e.Status]; ok {

		notification.Color = color
	}

	return notification
}
