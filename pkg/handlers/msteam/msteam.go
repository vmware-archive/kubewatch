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

package msteam

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/bitnami-labs/kubewatch/config"
	"github.com/bitnami-labs/kubewatch/pkg/event"
)

var msteamsErrMsg = `
%s

You need to set the MS teams webhook URL,
using --webhookURL, or using environment variables:

export KW_MSTEAMS_WEBHOOKURL=webhook_url

Command line flags will override environment variables

`

var msTeamsColors = map[string]string{
	"Normal":  "2DC72D",
	"Warning": "DEFF22",
	"Danger":  "8C1A1A",
}

// Constants for Sending a Card
const (
	messageType = "MessageCard"
	context     = "http://schema.org/extensions"
)

// TeamsMessageCard is for the Card Fields to send in Teams
// The Documentation is in https://docs.microsoft.com/en-us/outlook/actionable-messages/card-reference#card-fields
type TeamsMessageCard struct {
	Type       string                    `json:"@type"`
	Context    string                    `json:"@context"`
	ThemeColor string                    `json:"themeColor"`
	Summary    string                    `json:"summary"`
	Title      string                    `json:"title"`
	Text       string                    `json:"text,omitempty"`
	Sections   []TeamsMessageCardSection `json:"sections"`
}

// TeamsMessageCardSection is placed under TeamsMessageCard.Sections
// Each element of AlertWebHook.Alerts will the number of elements of TeamsMessageCard.Sections to create
type TeamsMessageCardSection struct {
	ActivityTitle string                         `json:"activityTitle"`
	Facts         []TeamsMessageCardSectionFacts `json:"facts"`
	Markdown      bool                           `json:"markdown"`
}

// TeamsMessageCardSectionFacts is placed under TeamsMessageCardSection.Facts
type TeamsMessageCardSectionFacts struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// Default handler implements Handler interface,
// print each event with JSON format
type MSTeams struct {
	// TeamsWebhookURL is the webhook url of the Teams connector
	TeamsWebhookURL string
}

// sendCard sends the JSON Encoded TeamsMessageCard to the webhook URL
func sendCard(ms *MSTeams, card *TeamsMessageCard) (*http.Response, error) {
	buffer := new(bytes.Buffer)
	if err := json.NewEncoder(buffer).Encode(card); err != nil {
		return nil, fmt.Errorf("Failed encoding message card: %v", err)
	}
	res, err := http.Post(ms.TeamsWebhookURL, "application/json", buffer)
	if err != nil {
		return nil, fmt.Errorf("Failed sending to webhook url %s. Got the error: %v",
			ms.TeamsWebhookURL, err)
	}
	if res.StatusCode != http.StatusOK {
		resMessage, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, fmt.Errorf("Failed reading Teams http response: %v", err)
		}
		return nil, fmt.Errorf("Failed sending to the Teams Channel. Teams http response: %s, %s",
			res.Status, string(resMessage))
	}
	if err := res.Body.Close(); err != nil {
		return nil, err
	}
	return res, nil
}

// notifyMSTeams creates the TeamsMessageCard and send to webhook URL
func notifyMSTeams(ms *MSTeams, obj interface{}, action string) {
	card := &TeamsMessageCard{
		Type:    messageType,
		Context: context,
		Title:   fmt.Sprintf("kubewatch"),
		// Set a default Summary, this is required for Microsoft Teams
		Summary: "kubewatch notification received",
	}

	e := event.New(obj, action)
	card.ThemeColor = msTeamsColors[e.Status]

	var s TeamsMessageCardSection
	s.ActivityTitle = e.Message()
	s.Markdown = true
	card.Sections = append(card.Sections, s)

	if _, err := sendCard(ms, card); err != nil {
		log.Printf("%s\n", err)
		return
	}

	log.Printf("Message successfully sent to MS Teams")
}

// Init initializes handler configuration
func (ms *MSTeams) Init(c *config.Config) error {
	webhookURL := c.Handler.MSTeams.WebhookURL

	if webhookURL == "" {
		webhookURL = os.Getenv("KW_MSTEAMS_WEBHOOKURL")
	}

	if webhookURL == "" {
		return fmt.Errorf(msteamsErrMsg, "Missing MS teams webhook URL")
	}

	ms.TeamsWebhookURL = webhookURL
	return nil
}

// Notify on object creation
func (ms *MSTeams) ObjectCreated(obj interface{}) {
	notifyMSTeams(ms, obj, "created")
}

// Notify on object deletion
func (ms *MSTeams) ObjectDeleted(obj interface{}) {
	notifyMSTeams(ms, obj, "deleted")
}

// Notify on object update
func (ms *MSTeams) ObjectUpdated(oldObj, newObj interface{}) {
	notifyMSTeams(ms, oldObj, "updated")
}

// TestHandler tests the handler configurarion by sending test messages.
func (ms *MSTeams) TestHandler() {
	card := &TeamsMessageCard{
		Type:    messageType,
		Context: context,
		Title:   fmt.Sprintf("kubewatch"),
		// Set a default Summary, this is required for Microsoft Teams
		Summary: "kubewatch notification received",
	}

	var s TeamsMessageCardSection
	s.ActivityTitle = "Testing Handler Configuration. This is a Test message."
	s.Markdown = true
	card.Sections = append(card.Sections, s)

	if _, err := sendCard(ms, card); err != nil {
		log.Printf("%s\n", err)
		return
	}

	log.Printf("Message successfully sent to MS Teams")
}
