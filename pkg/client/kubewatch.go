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
	"k8s.io/kubernetes/pkg/client/restclient"
	k8sClient "k8s.io/kubernetes/pkg/client/unversioned"
)

const userAgent = "kubewatch-client"

// Client represent the kubewatch client
type Client struct {
	client *k8sClient.Client
	ua     string
}

// New creates new kubewatch client
func New() (*Client, error) {
	config := &restclient.Config{
		Host: "http://127.0.0.1:8080",
	}

	c, err := k8sClient.New(config)
	if err != nil {
		return nil, err
	}

	kubeWatchClient := &Client{
		client: c,
		ua:     userAgent,
	}

	return kubeWatchClient, nil
}

// Events returns k8sClient.EventInterface to work with events resource
func (c *Client) Events(namespace string) k8sClient.EventInterface {
	return c.client.Events(namespace)
}
