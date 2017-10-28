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

	"github.com/skippbox/kubewatch/config"
	"github.com/skippbox/kubewatch/pkg/event"
	kbEvent "github.com/skippbox/kubewatch/pkg/event"
)

var hipchatColors = map[string]string{
	"Normal":  "good",
	"Warning": "warning",
	"Danger":  "danger",
}

var hipchatErrMsg = `
%s

You need to set both hipchat token and channel for hipchat notify,
using "--token/-t" and "--channel/-c", or using environment variables:

export KW_HIPCHAT_TOKEN=hipchat_token
export KW_HIPCHAT_CHANNEL=hipchat_channel

Command line flags will override environment variables

`

// Hipchat handler implements handler.Handler interface,
// Notify event to hipchat channel
type Hipchat struct {
	Token   string
	Channel string
}

// Init prepares hipchat configuration
func (s *Hipchat) Init(c *config.Config) error {
	token := c.Handler.Hipchat.Token
	channel := c.Handler.Hipchat.Channel

	if token == "" {
		token = os.Getenv("KW_HIPCHAT_TOKEN")
	}

	if channel == "" {
		channel = os.Getenv("KW_HIPCHAT_CHANNEL")
	}

	s.Token = token
	s.Channel = channel

	return checkMissingHipchatVars(s)
}

func (s *Hipchat) ObjectCreated(obj interface{}) {
	notifyHipchat(s, obj, "created")
}

func (s *Hipchat) ObjectDeleted(obj interface{}) {
	notifyHipchat(s, obj, "deleted")
}

func (s *Hipchat) ObjectUpdated(oldObj, newObj interface{}) {
	notifyHipchat(s, newObj, "updated")
}

func notifyHipchat(s *Hipchat, obj interface{}, action string) {
	e := kbEvent.New(obj, action)
	api := hipchat.NewClient(s.Token)
	params := hipchat.PostMessageParameters{}
	attachment := prepareHipchatAttachment(e)

	params.Attachments = []hipchat.Attachment{attachment}
	params.AsUser = true
	channelID, timestamp, err := api.PostMessage(s.Channel, "", params)
	if err != nil {
		log.Printf("%s\n", err)
		return
	}

	log.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)
}

func checkMissingHipchatVars(s *Hipchat) error {
	if s.Token == "" || s.Channel == "" {
		return fmt.Errorf(hipchatErrMsg, "Missing hipchat token or channel")
	}

	return nil
}

func prepareHipchatAttachment(e event.Event) hipchat.Attachment {
	msg := fmt.Sprintf(
		"A %s in namespace %s has been %s: %s",
		e.Kind,
		e.Namespace,
		e.Reason,
		e.Name,
	)

	attachment := hipchat.Attachment{
		Fields: []hipchat.AttachmentField{
			hipchat.AttachmentField{
				Title: "kubewatch",
				Value: msg,
			},
		},
	}

	if color, ok := hipchatColors[e.Status]; ok {
		attachment.Color = color
	}

	attachment.MarkdownIn = []string{"fields"}

	return attachment
}
