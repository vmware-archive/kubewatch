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
	"path/filepath"
	"runtime"

	"gopkg.in/yaml.v2"
)

// ConfigFileName stores file of config
var ConfigFileName = ".kubewatch.yaml"

// Handler contains handler configuration
type Handler struct {
	Slack      Slack      `json:"slack"`
	Hipchat    Hipchat    `json:"hipchat"`
	Mattermost Mattermost `json:"mattermost"`
	Flock      Flock      `json:"flock"`
	Webhook    Webhook    `json:"webhook"`
	MSTeams    MSTeams    `json:"msteams"`
}

// Resource contains resource configuration
type Resource struct {
	Deployment            bool `json:"deployment"`
	ReplicationController bool `json:"rc"`
	ReplicaSet            bool `json:"rs"`
	DaemonSet             bool `json:"ds"`
	Services              bool `json:"svc"`
	Pod                   bool `json:"po"`
	Job                   bool `json:"job"`
	PersistentVolume      bool `json:"pv"`
	Namespace             bool `json:"ns"`
	Secret                bool `json:"secret"`
	ConfigMap             bool `json:"configmap"`
	Ingress               bool `json:"ing"`
}

// Config struct contains kubewatch configuration
type Config struct {
	Handler Handler `json:"handler"`
	//Reason   []string `json:"reason"`
	Resource Resource `json:"resource"`
	// for watching specific namespace, leave it empty for watching all.
	// this config is ignored when watching namespaces
	Namespace string `json:"namespace,omitempty"`
}

// Slack contains slack configuration
type Slack struct {
	Token   string `json:"token"`
	Channel string `json:"channel"`
}

// Hipchat contains hipchat configuration
type Hipchat struct {
	Token string `json:"token"`
	Room  string `json:"room"`
	Url   string `json:"url"`
}

// Mattermost contains mattermost configuration
type Mattermost struct {
	Channel  string `json:"room"`
	Url      string `json:"url"`
	Username string `json:"username"`
}

// Flock contains flock configuration
type Flock struct {
	Url string `json:"url"`
}

// Webhook contains webhook configuration
type Webhook struct {
	Url string `json:"url"`
}

// MSTeams contains MSTeams configuration
type MSTeams struct {
	WebhookURL string `json:"webhookurl"`
}

// New creates new config object
func New() (*Config, error) {
	c := &Config{}
	if err := c.Load(); err != nil {
		return c, err
	}

	return c, nil
}

func createIfNotExist() error {
	// create file if not exist
	configFile := filepath.Join(configDir(), ConfigFileName)
	_, err := os.Stat(configFile)
	if err != nil {
		if os.IsNotExist(err) {
			file, err := os.Create(configFile)
			if err != nil {
				return err
			}
			file.Close()
		} else {
			return err
		}
	}
	return nil
}

// Load loads configuration from config file
func (c *Config) Load() error {
	err := createIfNotExist()
	if err != nil {
		return err
	}

	file, err := os.Open(getConfigFile())
	if err != nil {
		return err
	}

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	if len(b) != 0 {
		return yaml.Unmarshal(b, c)
	}

	return nil
}

// CheckMissingResourceEnvvars will read the environment for equivalent config variables to set
func (c *Config) CheckMissingResourceEnvvars() {
	if !c.Resource.DaemonSet && os.Getenv("KW_DAEMONSET") == "true" {
		c.Resource.DaemonSet = true
	}
	if !c.Resource.ReplicaSet && os.Getenv("KW_REPLICASET") == "true" {
		c.Resource.ReplicaSet = true
	}
	if !c.Resource.Namespace && os.Getenv("KW_NAMESPACE") == "true" {
		c.Resource.Namespace = true
	}
	if !c.Resource.Deployment && os.Getenv("KW_DEPLOYMENT") == "true" {
		c.Resource.Deployment = true
	}
	if !c.Resource.Pod && os.Getenv("KW_POD") == "true" {
		c.Resource.Pod = true
	}
	if !c.Resource.ReplicationController && os.Getenv("KW_REPLICATION_CONTROLLER") == "true" {
		c.Resource.ReplicationController = true
	}
	if !c.Resource.Services && os.Getenv("KW_SERVICE") == "true" {
		c.Resource.Services = true
	}
	if !c.Resource.Job && os.Getenv("KW_JOB") == "true" {
		c.Resource.Job = true
	}
	if !c.Resource.PersistentVolume && os.Getenv("KW_PERSISTENT_VOLUME") == "true" {
		c.Resource.PersistentVolume = true
	}
	if !c.Resource.Secret && os.Getenv("KW_SECRET") == "true" {
		c.Resource.Secret = true
	}
	if !c.Resource.ConfigMap && os.Getenv("KW_CONFIGMAP") == "true" {
		c.Resource.ConfigMap = true
	}
	if !c.Resource.Ingress && os.Getenv("KW_INGRESS") == "true" {
		c.Resource.Ingress = true
	}
	if (c.Handler.Slack.Channel == "") && (os.Getenv("SLACK_CHANNEL") != "") {
		c.Handler.Slack.Channel = os.Getenv("SLACK_CHANNEL")
	}
	if (c.Handler.Slack.Token == "") && (os.Getenv("SLACK_TOKEN") != "") {
		c.Handler.Slack.Token = os.Getenv("SLACK_TOKEN")
	}
}

func (c *Config) Write() error {
	b, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(getConfigFile(), b, 0644)
	if err != nil {
		return err
	}

	return nil
}

func getConfigFile() string {
	configFile := filepath.Join(configDir(), ConfigFileName)
	if _, err := os.Stat(configFile); err == nil {
		return configFile
	}

	return ""
}

func configDir() string {
	if configDir := os.Getenv("KW_CONFIG"); configDir != "" {
		return configDir
	}

	if runtime.GOOS == "windows" {
		home := os.Getenv("USERPROFILE")
		return home
	}
	return os.Getenv("HOME")
	//path := "/etc/kubewatch"
	//if _, err := os.Stat(path); os.IsNotExist(err) {
	//	os.Mkdir(path, 755)
	//}
	//return path
}
