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

package googlechat

import (
	"fmt"
	"log"
	"os"

	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/bitnami-labs/kubewatch/config"
	kbEvent "github.com/bitnami-labs/kubewatch/pkg/event"
)

var GoogleChatErrorMsg = `
%s

You need to set GoogleChat url
using "--url/-u" or using environment variables:

export KW_GOOGLECHAT_URL=googleChat url

Command line flags will override environment variables

`

// GoogleChat handler implements handler.Handler interface,
// Notify event to GoogleChat channel
type GoogleChat struct {
	Url string
}

type GoogleChatMessage struct {
	Text string `json:"text"`
	Cards []GoogleCard `json:"cards"`
}

type GoogleCard struct{
	Header GoogleCardHeader `json:header`
	Sections []GoogleCardSection `json:sections`
}

type GoogleCardHeader struct{
	Title string `json:"title"`
}

type GoogleCardSection struct{
	Header string `json:"header"`
	Widgets []GoogleCardWidgetParagraph `json:"widgets"`
}

type GoogleCardWidgetParagraph struct{
	TextParagraph GoogleCardWidgetText `json:"textParagraph"`
}

type GoogleCardWidgetText struct{
	Text string `json:"text"`
}

// Init prepares GoogleChat configuration
func (m *GoogleChat) Init(c *config.Config) error {
	url := c.Handler.GoogleChat.Url

	if url == "" {
		url = os.Getenv("KW_GOOGLECHAT_URL")
	}

	m.Url = url

	return checkMissingGChatVars(m)
}

func (m *GoogleChat) ObjectCreated(obj interface{}) {
	notifyGoogleChat(m, obj, "created")
}

func (m *GoogleChat) ObjectDeleted(obj interface{}) {
	notifyGoogleChat(m, obj, "deleted")
}

func (m *GoogleChat) ObjectUpdated(oldObj, newObj interface{}) {
	notifyGoogleChat(m, newObj, "updated")
}


func notifyGoogleChat(m *GoogleChat, obj interface{}, action string) {
	e := kbEvent.New(obj, action)

	chatMessage := prepareChatMessage(e, m)

	err := postMessage(m.Url, chatMessage)
	if err != nil {
		log.Printf("%s\n", err)
		return
	}

	log.Printf("Message successfully sent to %s at %s ", m.Url, time.Now())
}

func checkMissingGChatVars(s *GoogleChat) error {
	if s.Url == "" {
		return fmt.Errorf(GoogleChatErrorMsg, "Missing GoogleChat url")
	}

	return nil
}

func prepareChatMessage(e kbEvent.Event, m *GoogleChat) *GoogleChatMessage {
	return &GoogleChatMessage{
		Text: "",
		Cards: []GoogleCard{
			{
				Header: GoogleCardHeader{
					Title:	"Kubewatch Notifications",
				},
				Sections: []GoogleCardSection{
					{
						Header: "Message",
						Widgets: []GoogleCardWidgetParagraph{
							{
								TextParagraph: GoogleCardWidgetText{
									Text:	e.Message(),
								},
							},
						},
					},
				},
			},
		},
	}
}

func postMessage(url string, chatMessage *GoogleChatMessage) error {
	message, err := json.Marshal(chatMessage)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(message))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		return err
	}

	return nil
}