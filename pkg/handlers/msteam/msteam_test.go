package msteam

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/bitnami-labs/kubewatch/config"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
				ActivityTitle: "A `pod` in namespace `new` has been `created`:\n`foo`",
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
	p := &v1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			UID:       "12345678",
			Name:      "foo",
			Namespace: "new",
		},
	}
	ms.ObjectCreated(p)
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
				ActivityTitle: "A `pod` in namespace `new` has been `deleted`:\n`foo`",
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
	p := &v1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			UID:       "12345678",
			Name:      "foo",
			Namespace: "new",
		},
	}
	ms.ObjectDeleted(p)
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
				ActivityTitle: "A `pod` in namespace `new` has been `updated`:\n`foo`",
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

	oldP := &v1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			UID:       "12345678",
			Name:      "foo",
			Namespace: "new",
		},
	}

	newP := &v1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			UID:       "12345678",
			Name:      "foo-new",
			Namespace: "new",
		},
	}

	ms.ObjectUpdated(oldP, newP)
}
