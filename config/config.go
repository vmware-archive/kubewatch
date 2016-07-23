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
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var configFileName = "kubewatch.conf.json"
var defaultConfigDir = "/etc/kubewatch"

var (
	slackToken   string
	slackChannel string
	configFile   string
	reasonFlag   string
)

// Config struct contains kubewatch configuration
type Config struct {
	FileName string `json:"-"`
	Handler  struct {
		Slack `json:"slack"`
	} `json:"handler"`
	Reason []string `json:"reason"`
}

// Slack contains slack configuration
type Slack struct {
	Token   string `json:"token"`
	Channel string `json:"channel"`
}

func init() {
	flag.StringVar(&slackToken, "slack-token", "", "Slack token")
	flag.StringVar(&slackChannel, "slack-channel", "", "Slack channel")
	flag.StringVar(&configFile, "config-file", "", "Configuration file")
	flag.StringVar(&reasonFlag, "reason", "", "Filter event by events, comma separated string")
}

// New creates new config object
func New() *Config {
	c := &Config{}
	c.FileName = getConfigFile()

	return c
}

// Load loads configuration from config file
func (c *Config) Load() error {
	defer loadFromFlag(c)

	file, err := os.Open(c.FileName)
	if err != nil {
		return err
	}

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, c)
}

func loadFromFlag(c *Config) {
	if slackToken != "" {
		c.Handler.Slack.Token = slackToken
	}
	if slackChannel != "" {
		c.Handler.Slack.Channel = slackChannel
	}
	if reasonFlag != "" {
		c.Reason = strings.Split(reasonFlag, ",")
	}
}

func getConfigFile() string {
	if configFile != "" {
		return configFile
	}

	curDir, err := os.Getwd()
	if err != nil {
		return ""
	}

	configFiles := []string{
		filepath.Join(curDir, configFileName),
		filepath.Join(homeDir(), "."+configFileName),
		filepath.Join(defaultConfigDir, configFileName),
	}

	for _, f := range configFiles {
		if _, err := os.Stat(f); err == nil {
			return f
		}
	}

	return ""
}

func homeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("USERPROFILE")
		return home
	}
	return os.Getenv("HOME")
}
