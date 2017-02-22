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

package alertmanager

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/skippbox/kubewatch/config"
	"github.com/skippbox/kubewatch/pkg/event"
	"github.com/skippbox/kubewatch/pkg/handlers"
	"log"
	"net/http"
	"strings"
)

// Alert is the data structure which contains the alert information. This structure is stored inside and alerts list.
type Alert struct {
	Labels       map[string]string `json:"labels"`
	Annotations  map[string]string `json:"annotations"`
	GeneratorURL string            `json:"generatorURL"`
}

// Alerts is the end data structure that is converted into JSON and posted to alert managers /api/v1/alerts
// endpoint.
type Alerts []Alert

// AlertManager is the underlying struct used by the alert manager handler receivers
type AlertManager struct {
	url    string
	labels []string
	config *config.Config
}

var alertManagerErrMsg = `
%s

You need to set both alertmanager url for alert manager notify,
using "--url/-u", or using environment variables:

export ALERTMANAGER_URL=alertmanager_url

Command line flags will override environment variables

`

// New returns a alert manager handler interface
func New(conf *config.Config, url string, labels []string) handlers.Handler {
	c := AlertManager{
		url:    url,
		labels: labels,
		config: conf,
	}
	handler := handlers.Handler(&c)
	return handler
}

// Init prepares slack configuration
func (s *AlertManager) Init() error {
	return checkMissingAlertManagerVars(s)
}

// Config returns the config data that will be used by the handler
func (s *AlertManager) Config() *config.Config {
	return s.config
}

func (s *AlertManager) ObjectCreated(obj interface{}) {
	notifyAlertManager(s, obj, "created")
}

func (s *AlertManager) ObjectDeleted(obj interface{}) {
	notifyAlertManager(s, obj, "deleted")
}

func (s *AlertManager) ObjectUpdated(oldObj, newObj interface{}) {
	notifyAlertManager(s, newObj, "updated")
}

func notifyAlertManager(s *AlertManager, obj interface{}, action string) {

	e := event.New(obj, action)

	labels := make(map[string]string)
	annotations := make(map[string]string)

	labels["namespace"] = e.Namespace
	labels["name"] = e.Name
	labels["status"] = e.Status
	labels["reason"] = e.Reason
	labels["component"] = e.Component
	labels["host"] = e.Host
	labels["kind"] = e.Kind
	labels["client"] = "kubewatch"
	labels["action"] = action

	for _, label := range s.labels {
		values := strings.SplitN(label, "=", 1)
		if len(values) > 1 {
			labels[values[0]] = values[1]
		}
	}

	alert := Alert{
		Labels:      labels,
		Annotations: annotations,
	}

	alerts := Alerts{alert}

	url := fmt.Sprintf("%v/api/v1/alerts", s.url)

	jsonBytes, err := json.Marshal(alerts)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBytes))
	if err != nil {
		log.Println(err)
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Printf("Non 200 HTTP response received - %v - %v", resp.StatusCode, resp.Status)
		return
	}

	log.Printf("Message with action \"%v\" was successfully sent to alertmanager (%s)", action, url)
}

func checkMissingAlertManagerVars(s *AlertManager) error {
	if s.url == "" {
		return fmt.Errorf(alertManagerErrMsg, "Missing alertmanager url")
	}

	return nil
}
