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
	"fmt"
	"reflect"
	"testing"

)

func TestSlackInit(t *testing.T) {

	expectedError := fmt.Errorf(slackErrMsg, "Missing slack token or channel")

	var Tests = []struct {
		slack Slack
		err   error
	}{
		{Slack{token: "foo", channel: "bar"}, nil},
		{Slack{token: "foo"}, expectedError},
		{Slack{channel: "bar"}, expectedError},
		{Slack{}, expectedError},
	}

	for _, tt := range Tests {
		if err := tt.slack.Init(); !reflect.DeepEqual(err, tt.err) {
			t.Fatalf("Init(): %v", err)
		}
	}
}
