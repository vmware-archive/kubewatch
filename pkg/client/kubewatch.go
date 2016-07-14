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
