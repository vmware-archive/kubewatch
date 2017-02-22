# Kubewatch
[![Build Status](https://travis-ci.org/skippbox/kubewatch.svg?branch=master)](https://travis-ci.org/skippbox/kubewatch) [![Join us on Slack](https://s3.eu-central-1.amazonaws.com/ngtuna/join-us-on-slack.png)](https://skippbox.herokuapp.com)

`kubewatch` is a Kubernetes watcher that currently publishes notification to the following backend event handlers:

* Slack
* Alertmanager

## Run kubewatch in a Kubernetes cluster

In order to run kubewatch in a Kubernetes cluster quickly, the easiest way is for you to create a [ConfigMap](https://github.com/skippbox/kubewatch/blob/master/kubewatch-configmap.yaml) to hold kubewatch configuration.

An example is provided at [`kubewatch-configmap.yaml`](https://github.com/skippbox/kubewatch/blob/master/kubewatch-configmap.yaml), do not forget to update your own slack channel and token parameters. Alternatively, you could use secrets.

Create k8s configmap:

```console
$ kubectl create -f kubewatch-configmap.yaml
```
Create the [Pod](https://github.com/skippbox/kubewatch/blob/master/kubewatch.yaml) directly, or create your own deployment:

```console
$ kubectl create -f kubewatch.yaml
```

A `kubewatch` container will be created along with `kubectl` sidecar container in order to reach the API server.

Once the Pod is running, you will start seeing Kubernetes events in your configured backend event handler. Here is a screenshot of using the Slack backend handler:

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
```

## Building

### Building with go

- you need go v1.5 or later.
- if your working copy is not in your `GOPATH`, you need to set it accordingly.

```console
$ go build -o kubewatch main.go
```

You can also use the Makefile directly:

```console
$ make build
```

### Building with Docker

Buiding builder image:

```console
$ make builder-image
```

Using the `kubewatch-builder` image to build `kubewatch` binary:

```console
$ make binary-image
$ docker images
REPOSITORY          TAG                 IMAGE ID            CREATED              SIZE
kubewatch           latest              f1ade726c6e2        31 seconds ago       33.08 MB
kubewatch-builder   latest              6b2d325a3b88        About a minute ago   514.2 MB
```

## Download kubewatch package

```console
$ go get -u github.com/skippbox/kubewatch
```

## Configuration
Kubewatch supports `config` command for configuration. Config file will be saved at $HOME/.kubewatch.yaml

### Configure resources to be watched

```console
// rc, po and svc will be watched
$ kubewatch config resource --rc --po --svc

// only svc will be watched
$ kubewatch config resource --svc
```

### Backend event handlers

#### Slack

Provides event handler integrations with your Slack channels.

##### Running

``` 
kubewatch run slack
```

##### Parameters

* `-c, --channel`
* `-t, --token`

##### Environment Variables

* `SLACK_CHANNEL`
* `SLACK_TOKEN`

#### Alertmanager

Provides event handler intergrations with Prometheus's Alertmanager
 
##### Running 
 
``` 
kubewatch run alertmanager
```

##### Parameters

* `-l, --label`
* `-u, --url `

##### Environment Variables

* `ALERTMANAGER_URL`

