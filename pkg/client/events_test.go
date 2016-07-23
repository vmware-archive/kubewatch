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

package client

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/skippbox/kubewatch/config"

	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/watch"
)

var configStr = `
{
    "handler": {
        "slack": {
            "channel": "slack_channel",
            "token": "slack_token"
        }
    },
    "reason": ["Created", "Pulled", "Started"]
}
`

func TestFilter(t *testing.T) {
	var Tests = []struct {
		hasConfigFile bool
		content       []byte
		reason        string
		expected      bool
	}{
		{false, []byte(`""`), "whatever", true},
		{true, []byte(configStr), "Created", true},
		{true, []byte(configStr), "Createdd", false},
	}

	for _, tt := range Tests {
		c := config.New()
		if tt.hasConfigFile {
			tmpConfigFile, err := ioutil.TempFile("", "kubewatch")
			if err != nil {
				t.Fatalf("TestFilter(): %+v", err)
			}

			defer func() {
				_ = os.Remove(tmpConfigFile.Name())
			}()

			if _, err := tmpConfigFile.Write(tt.content); err != nil {
				t.Fatalf("TestFilter(): %+v", err)
			}
			if err := tmpConfigFile.Close(); err != nil {
				t.Fatalf("TestFilter(): %+v", err)
			}
			c.FileName = tmpConfigFile.Name()
		} else {
			c.FileName = ""
		}

		_ = c.Load()

		client, _ := New(c)

		e := watch.Event{}
		apiEvent := &api.Event{}
		apiEvent.Reason = tt.reason
		e.Object = apiEvent

		if client.Filter(e) != tt.expected {
			t.Fatalf("TestFilter(): %+v", client)
		}
	}
}
