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

package flock

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/bitnami-labs/kubewatch/config"
)

func TestFlockInit(t *testing.T) {
	s := &Flock{}
	expectedError := fmt.Errorf(flockErrMsg, "Missing Flock url")

	var Tests = []struct {
		flock config.Flock
		err   error
	}{
		{config.Flock{Url: "foo"}, nil},
		{config.Flock{}, expectedError},
	}

	for _, tt := range Tests {
		c := &config.Config{}
		c.Handler.Flock = tt.flock
		if err := s.Init(c); !reflect.DeepEqual(err, tt.err) {
			t.Fatalf("Init(): %v", err)
		}
	}
}
