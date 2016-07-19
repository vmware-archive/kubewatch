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

package slack

import (
	"testing"

	"k8s.io/kubernetes/pkg/watch"
)

func TestSlackInit(t *testing.T) {

	s := &Slack{}

	err := s.Init("foo", "bar")

	if err != nil {
		t.Fatal("NotifySlack(): should return error when missing token or channel")
	}
}

func TestSlackHandle(t *testing.T) {

	w := watch.Event{}
	s := &Slack{}

	err := s.Handle(w)

	if err == nil {
		t.Fatal("NotifySlack(): should return error when missing token or channel")
	}
}
