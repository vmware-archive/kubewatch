/*
Copyright 2018 Bitnami

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

package basecamp

import (
	"fmt"
	"log"
	"os"
	"strings"

	"net/http"
	"net/url"
	"time"

	"github.com/bitnami-labs/kubewatch/config"
	"github.com/bitnami-labs/kubewatch/pkg/event"
)

var basecampErrMsg = `
%s

You need to set BaseCamp url
using "--url/-u" or using environment variables:

export KW_BASECAMP_URL=basecamp_url

Command line flags will override environment variables

`

// BaseCamp handler implements handler.Handler interface,
// Notify event to BaseCamp campfire
type BaseCamp struct {
	Url string
}

// EventMeta contains the metadata about the occurred event
type EventMeta struct {
	Kind      string `json:"kind"`
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Reason    string `json:"reason"`
}

// Init prepares BaseCamp configuration
func (m *BaseCamp) Init(c *config.Config) error {
	url := c.Handler.BaseCamp.Url

	if url == "" {
		url = os.Getenv("KW_BASECAMP_URL")
	}

	m.Url = url

	return checkMissingBaseCampVars(m)
}

// Handle handles an event.
func (m *BaseCamp) Handle(e event.Event) {
	err := postMessage(m.Url, e.Message())
	if err != nil {
		log.Printf("%s\n", err)
		return
	}

	log.Printf("Message successfully sent to %s at %s ", m.Url, time.Now())
}

func checkMissingBaseCampVars(s *BaseCamp) error {
	if s.Url == "" {
		return fmt.Errorf(basecampErrMsg, "Missing BaseCamp url")
	}

	return nil
}

func postMessage(basecampURL string, basecampMessage string) error {
	data := url.Values{}
	data.Set("content", basecampMessage)

	req, err := http.NewRequest("POST", basecampURL, strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		return err
	}

	return nil
}
