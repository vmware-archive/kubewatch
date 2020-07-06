package msteam

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/bitnami-labs/kubewatch/config"
	"github.com/bitnami-labs/kubewatch/pkg/event"
)

// Tests the Init() function
func TestInit(t *testing.T) {
	s := &MSTeams{}
	expectedError := fmt.Errorf(msteamsErrMsg, "Missing MS teams webhook URL")

	var Tests = []struct {
		ms  config.MSTeams
		err error
	}{
		{config.MSTeams{WebhookURL: "somepath"}, nil},
		{config.MSTeams{}, expectedError},
	}

	for _, tt := range Tests {
		c := &config.Config{}
		c.Handler.MSTeams = tt.ms
		if err := s.Init(c); !reflect.DeepEqual(err, tt.err) {
			t.Fatalf("Init(): %v", err)
		}
	}
}

// Tests ObjectCreated() by passing v1.Pod
func TestObjectCreated(t *testing.T) {
	expectedCard := TeamsMessageCard{
		Type:       messageType,
		Context:    context,
		ThemeColor: msTeamsColors["Normal"],
		Summary:    "kubewatch notification received",
		Title:      "kubewatch",
		Text:       "",
		Sections: []TeamsMessageCardSection{
			{
				ActivityTitle: "A `pod` in namespace `new` has been `Created`:\n`foo`",
				Markdown:      true,
			},
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if r.Method != "POST" {
			t.Errorf("expected a POST request for ObjectCreated()")
		}
		decoder := json.NewDecoder(r.Body)
		var c TeamsMessageCard
		if err := decoder.Decode(&c); err != nil {
			t.Errorf("%v", err)
		}
		if !reflect.DeepEqual(c, expectedCard) {
			t.Errorf("expected %v, got %v", expectedCard, c)
		}
	}))

	ms := &MSTeams{TeamsWebhookURL: ts.URL}
	p := event.Event{
		Name:      "foo",
		Kind:      "pod",
		Namespace: "new",
		Reason:    "Created",
		Status:    "Normal",
	}

	ms.Handle(p)
}

// Tests ObjectDeleted() by passing v1.Pod
func TestObjectDeleted(t *testing.T) {
	expectedCard := TeamsMessageCard{
		Type:       messageType,
		Context:    context,
		ThemeColor: msTeamsColors["Danger"],
		Summary:    "kubewatch notification received",
		Title:      "kubewatch",
		Text:       "",
		Sections: []TeamsMessageCardSection{
			{
				ActivityTitle: "A `pod` in namespace `new` has been `Deleted`:\n`foo`",
				Markdown:      true,
			},
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if r.Method != "POST" {
			t.Errorf("expected a POST request for ObjectDeleted()")
		}
		decoder := json.NewDecoder(r.Body)
		var c TeamsMessageCard
		if err := decoder.Decode(&c); err != nil {
			t.Errorf("%v", err)
		}
		if !reflect.DeepEqual(c, expectedCard) {
			t.Errorf("expected %v, got %v", expectedCard, c)
		}
	}))

	ms := &MSTeams{TeamsWebhookURL: ts.URL}

	p := event.Event{
		Name:      "foo",
		Namespace: "new",
		Kind:      "pod",
		Reason:    "Deleted",
		Status:    "Danger",
	}

	ms.Handle(p)
}

// Tests ObjectUpdated() by passing v1.Pod
func TestObjectUpdated(t *testing.T) {
	expectedCard := TeamsMessageCard{
		Type:       messageType,
		Context:    context,
		ThemeColor: msTeamsColors["Warning"],
		Summary:    "kubewatch notification received",
		Title:      "kubewatch",
		Text:       "",
		Sections: []TeamsMessageCardSection{
			{
				ActivityTitle: "A `pod` in namespace `new` has been `Updated`:\n`foo`",
				Markdown:      true,
			},
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if r.Method != "POST" {
			t.Errorf("expected a POST request for ObjectUpdated()")
		}
		decoder := json.NewDecoder(r.Body)
		var c TeamsMessageCard
		if err := decoder.Decode(&c); err != nil {
			t.Errorf("%v", err)
		}
		if !reflect.DeepEqual(c, expectedCard) {
			t.Errorf("expected %v, got %v", expectedCard, c)
		}
	}))

	ms := &MSTeams{TeamsWebhookURL: ts.URL}

	oldP := event.Event{
		Name:      "foo",
		Namespace: "new",
		Kind:      "pod",
		Reason:    "Updated",
		Status:    "Warning",
	}

	newP := event.Event{
		Name:      "foo-new",
		Namespace: "new",
		Kind:      "pod",
		Reason:    "Updated",
		Status:    "Warning",
	}
	_ = newP

	ms.Handle(oldP)
}
