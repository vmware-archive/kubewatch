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

var ConfigFileName = ".kubewatch.yaml"

type Handler struct {
	Slack      Slack      `json:"slack"`
	Hipchat    Hipchat    `json:"hipchat"`
	Mattermost Mattermost `json:"mattermost"`
	Flock      Flock      `json:"flock"`
	Webhook    Webhook    `json:"webhook"`
}

// Resource contains resource configuration
type Resource struct {
	Deployment            WatchType `json:"deployment"`
	ReplicationController WatchType `json:"rc"`
	ReplicaSet            WatchType `json:"rs"`
	DaemonSet             WatchType `json:"ds"`
	Services              WatchType `json:"svc"`
	Pod                   WatchType `json:"po"`
	Job                   WatchType `json:"job"`
	PersistentVolume      WatchType `json:"pv"`
	Namespace             WatchType `json:"ns"`
	Secret                WatchType `json:"secret"`
	ConfigMap             WatchType `json:"configmap"`
	Ingress               WatchType `json:"ing"`
}

type WatchType struct {
	Watch  bool      `json:"watch"`
	Events EventType `json:"events"`
}

type EventType struct {
	Create             bool `json:"create"`
	Update             bool `json:"update"`
	Delete             bool `json:"delete"`
	LoadBalancerCreate bool `json:"loadbalancercreate"`
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

func (c *Config) CheckMissingResourceEnvvars() {
	if !c.Resource.DaemonSet.Watch && os.Getenv("KW_DAEMONSET") == "true" {
		c.Resource.DaemonSet.Watch = true
		c.Resource.DaemonSet.Events.Create = true
		c.Resource.DaemonSet.Events.Update = true
		c.Resource.DaemonSet.Events.Delete = true
	}
	if !c.Resource.ReplicaSet.Watch && os.Getenv("KW_REPLICASET") == "true" {
		c.Resource.ReplicaSet.Watch = true
		c.Resource.ReplicaSet.Events.Create = true
		c.Resource.ReplicaSet.Events.Update = true
		c.Resource.ReplicaSet.Events.Delete = true
	}
	if !c.Resource.Namespace.Watch && os.Getenv("KW_NAMESPACE") == "true" {
		c.Resource.Namespace.Watch = true
		c.Resource.Namespace.Events.Create = true
		c.Resource.Namespace.Events.Update = true
		c.Resource.Namespace.Events.Delete = true
	}
	if !c.Resource.Deployment.Watch && os.Getenv("KW_DEPLOYMENT") == "true" {
		c.Resource.Deployment.Watch = true
		c.Resource.Deployment.Events.Create = true
		c.Resource.Deployment.Events.Update = true
		c.Resource.Deployment.Events.Delete = true
	}
	if !c.Resource.Pod.Watch && os.Getenv("KW_POD") == "true" {
		c.Resource.Pod.Watch = true
		c.Resource.Pod.Events.Create = true
		c.Resource.Pod.Events.Update = true
		c.Resource.Pod.Events.Delete = true
	}
	if !c.Resource.ReplicationController.Watch && os.Getenv("KW_REPLICATION_CONTROLLER") == "true" {
		c.Resource.ReplicationController.Watch = true
		c.Resource.ReplicationController.Events.Create = true
		c.Resource.ReplicationController.Events.Update = true
		c.Resource.ReplicationController.Events.Delete = true
	}
	if !c.Resource.Services.Watch && os.Getenv("KW_SERVICE") == "true" {
		c.Resource.Services.Watch = true
		c.Resource.Services.Events.Create = true
		c.Resource.Services.Events.Update = true
		c.Resource.Services.Events.Delete = true
	}
	if !c.Resource.Job.Watch && os.Getenv("KW_JOB") == "true" {
		c.Resource.Job.Watch = true
		c.Resource.Job.Events.Create = true
		c.Resource.Job.Events.Update = true
		c.Resource.Job.Events.Delete = true
	}
	if !c.Resource.PersistentVolume.Watch && os.Getenv("KW_PERSISTENT_VOLUME") == "true" {
		c.Resource.PersistentVolume.Watch = true
		c.Resource.PersistentVolume.Events.Create = true
		c.Resource.PersistentVolume.Events.Update = true
		c.Resource.PersistentVolume.Events.Delete = true
	}
	if !c.Resource.Secret.Watch && os.Getenv("KW_SECRET") == "true" {
		c.Resource.Secret.Watch = true
		c.Resource.Secret.Events.Create = true
		c.Resource.Secret.Events.Update = true
		c.Resource.Secret.Events.Delete = true
	}
	if !c.Resource.ConfigMap.Watch && os.Getenv("KW_CONFIGMAP") == "true" {
		c.Resource.ConfigMap.Watch = true
		c.Resource.ConfigMap.Events.Create = true
		c.Resource.ConfigMap.Events.Update = true
		c.Resource.ConfigMap.Events.Delete = true
	}
	if !c.Resource.Ingress.Watch && os.Getenv("KW_INGRESS") == "true" {
		c.Resource.Ingress.Watch = true
		c.Resource.Ingress.Events.Create = true
		c.Resource.Ingress.Events.Update = true
		c.Resource.Ingress.Events.Delete = true
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
