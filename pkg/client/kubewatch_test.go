/*
Copyright 2016 Skippbox, Ltd.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
	c, err := New(nil)
	if err != nil {
		t.Fatal(err)
	}

	assertEqual(t, c.ua, userAgent)
}

func TestEvents(t *testing.T) {
	c, _ := New(nil)

	e := c.Events(api.NamespaceAll)

	_, ok := e.(k8sClient.EventInterface)

	if !ok {
		t.Fatal("Events(): wrong type")
	}
}
