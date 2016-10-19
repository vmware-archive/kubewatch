# Kubewatch
[![Build Status](https://travis-ci.org/skippbox/kubewatch.svg?branch=master)](https://travis-ci.org/skippbox/kubewatch) [![Join us on Slack](https://s3.eu-central-1.amazonaws.com/ngtuna/join-us-on-slack.png)](https://skippbox.herokuapp.com)

A Slack watcher for Kubernetes.

# Building

## Building with go

- you need go v1.5 or later.
- if your working copy is not in your `GOPATH`, you need to set it accordingly.

```console
$ go build -o kubewatch main.go
```

## Building with Docker

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

# Configuration
Kubewatch supports `config` command for configuration. Config file will be saved at $HOME/.kubewatch.yaml

## Configure slack
```console
$ kubewatch config slack --channel <slack_channel> --token <slack_token>
```

## Configure resources to be watched
```console
// rc, po and svc will be watched
$ kubewatch config resource --rc --po --svc

// only svc will be watched
$ kubewatch config resource --svc
```

## Environment variables
You have an altenative choice to figure out SLACK token, channel via environment variables:

```console
$ export KW_SLACK_TOKEN='XXXXXXXXXXXXXXXX'
$ export KW_SLACK_CHANNEL='#channel_name'
```

# Run kubewatch

```console
$ kubewatch
```

# Run kubewatch in a Kubernetes cluster

Create k8s secrets to hold slack token, channel:
```console
$ kubectl create secret generic kubewatch --from-literal=token=<token> --from-literal=channel=<channel>
```

Create the Pod:
```console
$ kubectl create -f kubewatch.yaml
```

A `kubewatch` sidecar container will be created along with `kubectl` main container in order to reach the API server.

# Testing with make

```console
$ make test
```
