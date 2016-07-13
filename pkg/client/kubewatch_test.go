package client

import (
	"testing"

	"k8s.io/kubernetes/pkg/api"
	k8sClient "k8s.io/kubernetes/pkg/client/unversioned"
)

func assertEqual(t *testing.T, result interface{}, expect interface{}) {
	if result != expect {
		t.Fatalf("Expect (Value: %v) (Type: %T) - Got (Value: %v) (Type: %T)", expect, expect, result, result)
	}
}

func TestNew(t *testing.T) {
	c, err := New()
	if err != nil {
		t.Fatal(err)
	}

	assertEqual(t, c.ua, userAgent)
}

func TestEvents(t *testing.T) {
	c, _ := New()

	e := c.Events(api.NamespaceAll)

	_, ok := e.(k8sClient.EventInterface)

	if !ok {
		t.Fatal("Events(): wrong type")
	}
}
