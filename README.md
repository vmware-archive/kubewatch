# Kubewatch

A Slack bot for Kubernetes

# Installation

```
mkdir -p $GOPATH/github.com/runseb
cd $GOPATH/github.com/runseb
git clone https://github.com/runseb/kubewatch.git
go install -v github.com/runseb/kubewatch
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
