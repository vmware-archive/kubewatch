# Kubewatch

Kubewatch contains three components: client, config, handler

![Kubewatch Diagram](kubewatch.png?raw=true "Kubewatch Overview")

## Client

The client gets events from `kubeapi-server`, filtering and applying handler to event.

It contains a config object and embedded with a k8s client.

## Config

The config object contains `kubewatch` configuration, like handlers, filters.

A config object is used to creating new client.

## Handler

Handler manages how `kubewatch` handles events.

With each event get from k8s and matched filtering from configuration, it is passed to handler. Currently, `kubewatch` has 2 handlers:

 - `Default`: which just print the event in JSON format
 - `Slack`: which send notification to slack channel based on information from config

More handlers will be added in future.

Each handler must implement the [Handler interface](https://github.com/skippbox/kubewatch/blob/master/pkg/handlers/handler.go#L31)
