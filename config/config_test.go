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

package config

import (
	"io/ioutil"
	"os"
	"testing"
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

func init() {
	configFile = "qqq"
}

func Test_getConfigFile(t *testing.T) {
	if f := getConfigFile(); f != configFile {
		t.Fatalf("getConfigFile(): %+v", f)
	}
}

func TestLoadOK(t *testing.T) {
	content := []byte(configStr)
	tmpConfigFile, err := ioutil.TempFile("", "kubewatch")
	if err != nil {
		t.Fatalf("TestLoad(): %+v", err)
	}

	defer func() {
		_ = os.Remove(tmpConfigFile.Name())
	}()

	if _, err := tmpConfigFile.Write(content); err != nil {
		t.Fatalf("TestLoad(): %+v", err)
	}
	if err := tmpConfigFile.Close(); err != nil {
		t.Fatalf("TestLoad(): %+v", err)
	}

	c := New()
	c.FileName = tmpConfigFile.Name()

	err = c.Load()
	if err != nil {
		t.Fatalf("TestLoad(): %+v", err)
	}
}

func TestLoadNotOK(t *testing.T) {
	var Tests = []struct {
		hasConfigFile bool
		content       []byte
	}{
		{false, []byte(`""`)},
		{true, []byte(`{"invalid json`)},
	}

	for _, tt := range Tests {
		c := New()
		if tt.hasConfigFile {
			tmpConfigFile, err := ioutil.TempFile("", "kubewatch")
			if err != nil {
				t.Fatalf("TestLoadNotOK(): %+v", err)
			}

			defer func() {
				_ = os.Remove(tmpConfigFile.Name())
			}()

			if _, err := tmpConfigFile.Write(tt.content); err != nil {
				t.Fatalf("TestLoadNotOK(): %+v", err)
			}
			if err := tmpConfigFile.Close(); err != nil {
				t.Fatalf("TestLoadNotOK(): %+v", err)
			}
			c.FileName = tmpConfigFile.Name()
		} else {
			c.FileName = ""
		}

		err := c.Load()
		if err == nil {
			t.Fatalf("TestLoadNotOK(): %+v", err)
		}

	}
}

func TestFilterConfig(t *testing.T) {
	var Tests = []struct {
		hasConfigFile bool
		content       []byte
		length        int
	}{
		{false, []byte(`""`), 0},
		{true, []byte(`""`), 0},
		{true, []byte(configStr), 3},
	}

	for _, tt := range Tests {
		c := New()
		if tt.hasConfigFile {
			tmpConfigFile, err := ioutil.TempFile("", "kubewatch")
			if err != nil {
				t.Fatalf("TestFilterConfig(): %+v", err)
			}

			defer func() {
				_ = os.Remove(tmpConfigFile.Name())
			}()

			if _, err := tmpConfigFile.Write(tt.content); err != nil {
				t.Fatalf("TestFilterConfig(): %+v", err)
			}
			if err := tmpConfigFile.Close(); err != nil {
				t.Fatalf("TestFilterConfig(): %+v", err)
			}
			c.FileName = tmpConfigFile.Name()
		} else {
			c.FileName = ""
		}

		_ = c.Load()
		if len(c.Reason) != tt.length {
			t.Fatalf("TestFilterConfig(): %+v", c)
		}

	}
}
