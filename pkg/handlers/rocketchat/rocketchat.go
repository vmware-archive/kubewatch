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

package rocketchat

import (
	"fmt"
	"strconv"


	rcapi "github.com/RocketChat/Rocket.Chat.Go.SDK/rest"
	"github.com/RocketChat/Rocket.Chat.Go.SDK/models"
	"github.com/bitnami-labs/kubewatch/config"
	"github.com/bitnami-labs/kubewatch/pkg/event"
	kbEvent "github.com/bitnami-labs/kubewatch/pkg/event"
	"log"
	"net/url"
	"os"
)

var rocketChatColors = map[string]string{
	"Normal":  "green",
	"Warning": "orange",
	"Danger":  "red",
}

var rocketChatErrMsg = `
%s

You need to set parameters  for rocketchat ,
export RC_HOST=localhost
export RC_PORT=80
export RC_SCHEME=http
export RC_EMAIL=kubewatch@rocketchat.com
export RC_USER=kubewatch
export RC_PASSWORD=kubewatch
export RC_CHANNEL=general

Command line flags will override environment variables

`

// Rocket handler implements handler.Handler interface,
// Notify rc
type Rocketchat struct {
	Host     string
	Port     int
	User     string
	Email    string
	Password string
	Scheme   string
	Channel string
}

// Init prepares slack configuration
func (s *Rocketchat) Init(c *config.Config) error {

	host := c.Handler.Rocketchat.Host
	port := c.Handler.Rocketchat.Port
	user := c.Handler.Rocketchat.User
	password := c.Handler.Rocketchat.Password
	scheme := c.Handler.Rocketchat.Scheme
	channel := c.Handler.Rocketchat.Channel
	email := c.Handler.Rocketchat.Email

	if host == "" {
		host = os.Getenv("RC_HOST")
	}

	if port == 0 {
		port, _ = strconv.Atoi(os.Getenv("RC_PORT"))
	}

	if scheme == "" {
		scheme = os.Getenv("RC_SCHEME")
	}

	if user == "" {
		user= os.Getenv("RC_USER")
	}

	if email == "" {
		email = os.Getenv("RC_EMAIL")
	}
	if password == "" {
		password= os.Getenv("RC_PASSWORD")
	}

	if channel == "" {
		channel =os.Getenv("RC_CHANNEL")
	}


	s.Host = host
	s.Port = port
	s.User = user
	s.Email =email
	s.Channel= channel
	s.Password = password
	s.Scheme = scheme

	return checkMissingRocketChatVars(s)
}

func (s *Rocketchat) ObjectCreated(obj interface{}) {
	notifyRocketChat(s, obj, "created")
}

func (s *Rocketchat) ObjectDeleted(obj interface{}) {
	notifyRocketChat(s, obj, "deleted")
}

func (s *Rocketchat) ObjectUpdated(oldObj, newObj interface{}) {
	notifyRocketChat(s, newObj, "updated")
}

func notifyRocketChat(s *Rocketchat, obj interface{}, action string) {

	client := rcapi.NewClient(&url.URL{Scheme :s.Scheme, Host: s.Host + ":" + strconv.Itoa(s.Port)}, true)
	credentials := &models.UserCredentials{Name: s.User, Email: s.Email, Password: s.Password}
	err := client.Login(credentials)
	if err != nil {
		panic(err)
	}

	e := kbEvent.New(obj, action)
	attachment := prepareRocketChatAttachment(e)

	attachments := []models.Attachment{attachment}
	msg := models.PostMessage{Attachments: attachments,Channel: s.Channel}
	msgResponse, err := client.PostMessage(&msg)
	if err != nil {
		fmt.Printf("failed Message Send %v", err)
	}

	log.Printf("Message successfully sent  %v ", msgResponse)
}

func checkMissingRocketChatVars(s *Rocketchat) error {
	if s.Host == "" || s.Port == 0 || s.Email== "" || s.User == ""  || s.Password == "" || s.Channel== ""{
		return fmt.Errorf(rocketChatErrMsg, "Missing rocket configuration paramaters")
	}



	return nil
}

func prepareRocketChatAttachment(e event.Event) models.Attachment {

	attachment := models.Attachment{
		Fields: []models.AttachmentField{
			{
				Title: e.Kind,
				Value: e.Message(),
			},
		},
	}

	if color, ok := rocketChatColors[e.Status]; ok {
		attachment.Color = color
	}

	return attachment
}
