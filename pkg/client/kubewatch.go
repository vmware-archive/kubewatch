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
	"sync"

	"k8s.io/kubernetes/pkg/client/restclient"
	k8sClient "k8s.io/kubernetes/pkg/client/unversioned"

	"github.com/skippbox/kubewatch/config"
)

const userAgent = "kubewatch-client"

// Client represent the kubewatch client
type Client struct {
	client    *k8sClient.Client
	ua        string
	closeChan chan bool
	waitGroup *sync.WaitGroup
	Config    *config.Config
}

// New creates new kubewatch client
func New(conf *config.Config, k8sClientConfig *restclient.Config) (*Client, error) {
	var err error
	if conf == nil {
		conf, err = config.New()
		if err != nil {
			return nil, err
		}
	}

	if k8sClientConfig == nil {
		k8sClientConfig = &restclient.Config{}
	}

	c, err := k8sClient.New(k8sClientConfig)
	if err != nil {
		return nil, err
	}

	kubeWatchClient := &Client{
		client:    c,
		ua:        userAgent,
		closeChan: make(chan bool),
		waitGroup: &sync.WaitGroup{},
		Config:    conf,
	}

	return kubeWatchClient, nil
}

// Events returns k8sClient.EventInterface to work with events resource
func (c *Client) Events(namespace string) k8sClient.EventInterface {
	return c.client.Events(namespace)
}

// Services returns k8sClient.ServiceInterface to work with services resource
func (c *Client) Services(namespace string) k8sClient.ServiceInterface {
	return c.client.Services(namespace)
}

// Stop stops kubewatch client and all its goroutine
func (c *Client) Stop() {
	close(c.closeChan)
	c.waitGroup.Wait()
}
