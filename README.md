# Kubewatch
[![Build Status](https://travis-ci.org/skippbox/kubewatch.svg?branch=master)](https://travis-ci.org/skippbox/kubewatch) [![Join us on Slack](https://s3.eu-central-1.amazonaws.com/ngtuna/join-us-on-slack.png)](https://skippbox.herokuapp.com)

A Slack watcher for Kubernetes

# Installation

## Manual
```
go get -u github.com/skippbox/kubewatch
```

## Building with Dockerfiles

Buiding builder image:

```
$ make builder-image
```

Using the `kubewatch-builder` image to build `kubewatch` binary:

```
$ make binary-image
$ docker images
REPOSITORY          TAG                 IMAGE ID            CREATED              SIZE
kubewatch           latest              f1ade726c6e2        31 seconds ago       33.08 MB
kubewatch-builder   latest              6b2d325a3b88        About a minute ago   514.2 MB
```

# Configuration
You can use configuration to specify `kubewatch` configuration, see example in `examples/conf/kubewatch.conf.json`

# Environment variables
Preparing your SLACK token, channel.

```
export KW_SLACK_TOKEN='XXXXXXXXXXXXXXXX'
export KW_SLACK_CHANNEL='#channel_name'
```

# Run Locally

```
"$GOPATH"/bin/kubewatch
```

# Run in a Kubernetes cluster

Create k8s secrets to hold slack token, channel:
```sh
kubectl create secret generic kubewatch --from-literal=token=<token> --from-literal=channel=<channel>
```

Create the Pod:
```sh
kubectl create -f kubewatch.yaml
```

It uses a `kubectl` side car container to reach the API server.


# Testing

```
$ make test
```
