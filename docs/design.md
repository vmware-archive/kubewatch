# Kubewatch

Kubewatch contains three components: controller, config, handler

![Kubewatch Diagram](kubewatch.png?raw=true "Kubewatch Overview")

## Config

The config object contains `kubewatch` configuration, like handlers, filters.

A config object is used to creating new client.

## Controller

The controller initializes using the config object by reading the `.kubewatch.yaml` or command line arguments.
If the parameters are not fully mentioned, the config falls back to read a set of standard environment variables.

Controller creates necessary `SharedIndexInformer`s provided by `kubernetes/client-go` for listening and watching
resource changes. Controller updates this subscription information with Kubernetes API Server.

Whenever, the Kubernetes Controller Manager gets events related to the subscribed resources, it pushes the events to
`SharedIndexInformer`. This in-turn puts the events onto a rate-limiting queue for better handling of the events.

Controller picks the events from the queue and hands over the events to the appropriate handler after
necessary filtering.

## Handler

Handler manages how `kubewatch` handles events.

With each event get from k8s and matched filtering from configuration, it is passed to handler. Currently, `kubewatch` has 7 handlers:

 - `Default`: which just print the event in JSON format
 - `Flock`: which send notification to Flock channel based on information from config
 - `Hipchat`: which send notification to Hipchat room based on information from config
 - `Mattermost`: which send notification to Mattermost channel based on information from config
 - `MS Teams`: which send notification to MS Team incoming webhook based on information from config
 - `Slack`: which send notification to Slack channel based on information from config
 - `Smtp`: which sends notifications to email recipients using a SMTP server obtained from config

More handlers will be added in future.

Each handler must implement the [Handler interface](https://github.com/bitnami-labs/kubewatch/blob/master/pkg/handlers/handler.go#L31)
