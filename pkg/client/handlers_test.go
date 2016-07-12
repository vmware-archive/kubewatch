package client

import (
	"testing"

	"k8s.io/kubernetes/pkg/watch"
)

func TestNotifySlack(t *testing.T) {

	w := watch.Event{}

	err := NotifySlack(w)

	if err == nil {
		t.Fatal("NotifySlack(): should return error when missing token or channel")
	}
}

func TestPrintEvent(t *testing.T) {

	w := watch.Event{}

	err := PrintEvent(w)

	if err != nil {
		t.Fatal(err)
	}
}
