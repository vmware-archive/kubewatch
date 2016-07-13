# Kubewatch

A Slack watcher for Kubernetes

# Installation

```
go get -u github.com/skippbox/kubewatch
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
