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

package rocketchat

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/bitnami-labs/kubewatch/config"
)

func TestRocketChatInit(t *testing.T) {
	s := &Rocketchat{}
	expectedError := fmt.Errorf(rocketChatErrMsg, "Missing rocket configuration paramaters")

	var Tests = []struct {
		rocketChat config.Rocketchat
		err        error
	}{

		{config.Rocketchat{Port: 0}, expectedError},
		{config.Rocketchat{Email: ""}, expectedError},
		{config.Rocketchat{Host: ""}, expectedError},
		{config.Rocketchat{User: ""}, expectedError},
		{config.Rocketchat{Password: ""}, expectedError},
		{config.Rocketchat{Channel: ""}, expectedError},
		{config.Rocketchat{Scheme: ""}, expectedError},
		{config.Rocketchat{}, expectedError},
	}

	for _, tt := range Tests {
		c := &config.Config{}
		c.Handler.Rocketchat = tt.rocketChat
		if err := s.Init(c); !reflect.DeepEqual(err, tt.err) {
			t.Fatalf("Init(): %v", err)
		}
	}
}
