<p align="center">
  <img src="./docs/kubewatch-logo.jpeg">
</p>


[![Build Status](https://travis-ci.org/bitnami-labs/kubewatch.svg?branch=master)](https://travis-ci.org/bitnami-labs/kubewatch) [![Go Report Card](https://goreportcard.com/badge/github.com/bitnami-labs/kubewatch)](https://goreportcard.com/report/github.com/bitnami-labs/kubewatch) [![GoDoc](https://godoc.org/github.com/bitnami-labs/kubewatch?status.svg)](https://godoc.org/github.com/bitnami-labs/kubewatch) [![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/bitnami-labs/kubewatch/blob/master/LICENSE)

**kubewatch** is a Kubernetes watcher that currently publishes notification to available collaboration hubs/notification channels. Run it in your k8s cluster, and you will get event notifications through webhooks.

# Usage
```
$ kubewatch -h

Kubewatch: A watcher for Kubernetes

kubewatch is a Kubernetes watcher that publishes notifications
to Slack/hipchat/mattermost/flock channels. It watches the cluster
for resource changes and notifies them through webhooks.

supported webhooks:
 - slack
 - hipchat
 - mattermost
 - flock
 - webhook
 - smtp

Usage:
  kubewatch [flags]
  kubewatch [command]

Available Commands:
  config      modify kubewatch configuration
  resource    manage resources to be watched
  version     print version

Flags:
  -h, --help   help for kubewatch

Use "kubewatch [command] --help" for more information about a command.

```

# Install

### Cluster Installation
#### Using helm:

When you have helm installed in your cluster, use the following setup:

```console
helm install --name kubewatch bitnami/kubewatch --set='rbac.create=true,slack.channel=#YOUR_CHANNEL,slack.token=xoxb-YOUR_TOKEN,resourcesToWatch.pod=true,resourcesToWatch.daemonset=true'
```

You may also provide a values file instead:

```yaml
rbac:
  create: true
resourcesToWatch:
  deployment: false
  replicationcontroller: false
  replicaset: false
  daemonset: false
  services: true
  pod: true
  job: false
  node: false
  clusterrole: true
  serviceaccount: true
  persistentvolume: false
  namespace: false
  secret: false
  configmap: false
  ingress: false
slack:
  channel: '#YOUR_CHANNEL'
  token: 'xoxb-YOUR_TOKEN'
```

And use that:

```console
$ helm upgrade --install kubewatch bitnami/kubewatch --values=values-file.yml
```

#### Using kubectl:

In order to run kubewatch in a Kubernetes cluster quickly, the easiest way is for you to create a [ConfigMap](https://github.com/bitnami-labs/kubewatch/blob/master/kubewatch-configmap.yaml) to hold kubewatch configuration.

An example is provided at [`kubewatch-configmap.yaml`](https://github.com/bitnami-labs/kubewatch/blob/master/kubewatch-configmap.yaml), do not forget to update your own slack channel and token parameters. Alternatively, you could use secrets.

Create k8s configmap:

```console
$ kubectl create -f kubewatch-configmap.yaml
```

Create the [Pod](https://github.com/bitnami-labs/kubewatch/blob/master/kubewatch.yaml) directly, or create your own deployment:

```console
$ kubectl create -f kubewatch.yaml
```

A `kubewatch` container will be created along with `kubectl` sidecar container in order to reach the API server.

Once the Pod is running, you will start seeing Kubernetes events in your configured Slack channel. Here is a screenshot:

![slack](./docs/slack.png)

To modify what notifications you get, update the `kubewatch` ConfigMap and turn on and off (true/false) resources:

```
resource:
  deployment: false
  replicationcontroller: false
  replicaset: false
  daemonset: false
  services: true
  pod: true
  job: false
  node: false
  clusterrole: false
  serviceaccount: false
  persistentvolume: false
  namespace: false
  secret: false
  configmap: false
  ingress: false
```

#### Working with RBAC

Kubernetes Engine clusters running versions 1.6 or higher introduced Role-Based Access Control (RBAC). We can create `ServiceAccount` for it to work with RBAC.

```console
$ kubectl create -f kubewatch-service-account.yaml
```

If you do not have permission to create it, you need to become an admin first. For example, in GKE you would run:

```
$ kubectl create clusterrolebinding cluster-admin-binding --clusterrole=cluster-admin --user=REPLACE_EMAIL_HERE
```

Edit `kubewatch.yaml`, and create a new field under `spec` with `serviceAccountName: kubewatch`, you can achieve this by running:

```console
$ sed -i '/spec:/a\ \ serviceAccountName: kubewatch' kubewatch.yaml
```

Then just create `pod` as usual with:

```console
$ kubectl create -f kubewatch.yaml
```

### Local Installation
#### Using go package installer:

```console
# Download and install kubewatch
$ go get -u github.com/bitnami-labs/kubewatch

# Configure the notification channel
$ kubewatch config add slack --channel <slack_channel> --token <slack_token>

# Add resources to be watched
$ kubewatch resource add --po --svc
INFO[0000] resource svc configured
INFO[0000] resource po configured

# start kubewatch server
$ kubewatch
INFO[0000] Starting kubewatch controller                 pkg=kubewatch-service
INFO[0000] Starting kubewatch controller                 pkg=kubewatch-pod
INFO[0000] Processing add to service: default/kubernetes  pkg=kubewatch-service
INFO[0000] Processing add to service: kube-system/tiller-deploy  pkg=kubewatch-service
INFO[0000] Processing add to pod: kube-system/tiller-deploy-69ffbf64bc-h8zxm  pkg=kubewatch-pod
INFO[0000] Kubewatch controller synced and ready         pkg=kubewatch-service
INFO[0000] Kubewatch controller synced and ready         pkg=kubewatch-pod

```
#### Using Docker:

To Run Kubewatch Container interactively, place the config file in `$HOME/.kubewatch.yaml` location and use the following command. 

```
docker run --rm -it --network host -v $HOME/.kubewatch.yaml:/root/.kubewatch.yaml -v $HOME/.kube/config:/opt/bitnami/kubewatch/.kube/config --name <container-name> bitnami/kubewatch
```

Example:

```
$ docker run --rm -it --network host -v $HOME/.kubewatch.yaml:/root/.kubewatch.yaml -v $HOME/.kube/config:/opt/bitnami/kubewatch/.kube/config --name kubewatch-app bitnami/kubewatch

==> Writing config file...
INFO[0000] Starting kubewatch controller                 pkg=kubewatch-service
INFO[0000] Starting kubewatch controller                 pkg=kubewatch-pod
INFO[0000] Starting kubewatch controller                 pkg=kubewatch-deployment
INFO[0000] Starting kubewatch controller                 pkg=kubewatch-namespace
INFO[0000] Processing add to namespace: kube-node-lease  pkg=kubewatch-namespace
INFO[0000] Processing add to namespace: kube-public      pkg=kubewatch-namespace
INFO[0000] Processing add to namespace: kube-system      pkg=kubewatch-namespace
INFO[0000] Processing add to namespace: default          pkg=kubewatch-namespace
....
```

To Demonise Kubewatch container use

```
$ docker run --rm -d --network host -v $HOME/.kubewatch.yaml:/root/.kubewatch.yaml -v $HOME/.kube/config:/opt/bitnami/kubewatch/.kube/config --name kubewatch-app bitnami/kubewatch
```

# Configure

Kubewatch supports `config` command for configuration. Config file will be saved at `$HOME/.kubewatch.yaml`

```
$ kubewatch config -h

config command allows admin setup his own configuration for running kubewatch

Usage:
  kubewatch config [flags]
  kubewatch config [command]

Available Commands:
  add         add webhook config to .kubewatch.yaml
  test        test handler config present in .kubewatch.yaml
  view        view .kubewatch.yaml

Flags:
  -h, --help   help for config

Use "kubewatch config [command] --help" for more information about a command.
```
### Example:

### slack:

- Create a [slack Bot](https://my.slack.com/services/new/bot)

- Edit the Bot to customize its name, icon and retrieve the API token (it starts with `xoxb-`).

- Invite the Bot into your channel by typing: `/invite @name_of_your_bot` in the Slack message area.

- Add Api token to kubewatch config using the following steps

  ```console
  $ kubewatch config add slack --channel <slack_channel> --token <slack_token>
  ```
  You have an altenative choice to set your SLACK token, channel via environment variables:

  ```console
  $ export KW_SLACK_TOKEN='XXXXXXXXXXXXXXXX'
  $ export KW_SLACK_CHANNEL='#channel_name'
  ```

### flock:

- Create a [flock bot](https://docs.flock.com/display/flockos/Bots).

- Add flock webhook url to config using the following command.
  ```console
  $ kubewatch config add flock --url <flock_webhook_url>
  ```
  You have an altenative choice to set your FLOCK URL

  ```console
  $ export KW_FLOCK_URL='https://api.flock.com/hooks/sendMessage/XXXXXXXX'
  ```

## Testing Config

To test the handler config by send test messages use the following command.
```
$ kubewatch config test -h

Tests handler configs present in .kubewatch.yaml by sending test messages

Usage:
  kubewatch config test [flags]

Flags:
  -h, --help   help for test
```

#### Example:

```
$ kubewatch config test

Testing Handler configs from .kubewatch.yaml
2019/06/03 12:29:23 Message successfully sent to channel ABCD at 1559545162.000100
```

## Viewing config
To view the entire config file `$HOME/.kubewatch.yaml` use the following command.
```
$ kubewatch config view
Contents of .kubewatch.yaml

handler:
  slack:
    token: xoxb-xxxxx-yyyy-zzz
    channel: kube-watch
  hipchat:
    token: ""
    room: ""
    url: ""
  mattermost:
    channel: ""
    url: ""
    username: ""
  flock:
    url: ""
  webhook:
    url: ""
resource:
  deployment: false
  replicationcontroller: false
  replicaset: false
  daemonset: false
  services: false
  pod: true
  job: false
  node: false
  clusterrole: false
  serviceaccount: false
  persistentvolume: false
  namespace: false
  secret: false
  configmap: false
  ingress: false
namespace: ""

```


## Resources

To manage the resources being watched, use the following command, changes will be saved to `$HOME/.kubewatch.yaml`.

```
$ kubewatch resource -h

manage resources to be watched

Usage:
  kubewatch resource [flags]
  kubewatch resource [command]

Available Commands:
  add         adds specific resources to be watched
  remove      remove specific resources being watched

Flags:
      --clusterrole   watch for cluster roles
      --cm            watch for plain configmaps
      --deploy        watch for deployments
      --ds            watch for daemonsets
  -h, --help          help for resource
      --ing           watch for ingresses
      --job           watch for jobs
      --node          watch for Nodes
      --ns            watch for namespaces
      --po            watch for pods
      --pv            watch for persistent volumes
      --rc            watch for replication controllers
      --rs            watch for replicasets
      --sa            watch for service accounts
      --secret        watch for plain secrets
      --svc           watch for services

Use "kubewatch resource [command] --help" for more information about a command.

```

### Add/Remove resource:
```
$ kubewatch resource add -h

adds specific resources to be watched

Usage:
  kubewatch resource add [flags]

Flags:
  -h, --help   help for add

Global Flags:
      --clusterrole   watch for cluster roles
      --cm            watch for plain configmaps
      --deploy        watch for deployments
      --ds            watch for daemonsets
      --ing           watch for ingresses
      --job           watch for jobs
      --node          watch for Nodes
      --ns            watch for namespaces
      --po            watch for pods
      --pv            watch for persistent volumes
      --rc            watch for replication controllers
      --rs            watch for replicasets
      --sa            watch for service accounts
      --secret        watch for plain secrets
      --svc           watch for services

```

### Example:

```console
# rc, po and svc will be watched
$ kubewatch resource add --rc --po --svc

# rc, po and svc will be stopped from being watched
$ kubewatch resource remove --rc --po --svc
```

# Build

### Using go

Clone the repository into your $GOPATH and then build it.
```
$ mkdir -p $GOPATH/src/github.com/bitnami-labs/
$ cd $GOPATH/src/github.com/bitnami-labs/
$ git clone https://github.com/bitnami-labs/kubewatch.git
$ cd kubewatch
$ go build -o kubewatch main.go
```
or

You can also use the Makefile directly:

```console
$ make build
```

#### Prerequisites

- You need to have [Go](http://golang.org) (v1.5 or later)  installed. Make sure to set `$GOPATH`


### Using Docker

```console
$ make docker-image
$ docker images
REPOSITORY          TAG                 IMAGE ID            CREATED              SIZE
kubewatch           latest              919896d3cd90        3 minutes ago       27.9MB
```
#### Prerequisites

- you need to have [docker](https://docs.docker.com/) installed.

# Contribution

Refer to the [contribution guidelines](docs/CONTRIBUTION.md) to get started.