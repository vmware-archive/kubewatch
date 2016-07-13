# Kubewatch

A Slack watcher for Kubernetes

# Installation

## Manual
```
go get -u github.com/skippbox/kubewatch
```

## Building with Dockerfiles

Buiding builder image:

```
docker build -t kubewatch-builder -f Dockerfile.build .
```

Using the `kubewatch-builder` image to build `kubewatch` binary:

```
$ docker run --rm kubewatch-builder | docker build -t kubewatch -f Dockerfile.run -
$ docker images
REPOSITORY          TAG                 IMAGE ID            CREATED              SIZE
kubewatch           latest              f1ade726c6e2        31 seconds ago       33.08 MB
kubewatch-builder   latest              6b2d325a3b88        About a minute ago   514.2 MB
```

# Environment variables
Preparing your SLACK token, channel.

```
export KW_SLACK_TOKEN='XXXXXXXXXXXXXXXX'
export KW_SLACK_CHANNEL='#channel_name'
```

# Run

```
"$GOPATH"/bin/kubewatch
```

# Testing

```
$ go test -v $(go list ./... | grep -v '/vendor/')
```

# Notes

For now, `kubewatch` watches the `kube-apiserver` on localhost only.
