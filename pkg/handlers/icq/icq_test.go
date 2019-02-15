/*
Copyright azalio.net

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

package icq

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/bitnami-labs/kubewatch/config"
)

func TestIcqInit(t *testing.T) {
	s := &Icq{}
	expectedError := fmt.Errorf(icqErrMsg, "Missing Icq token or uid")

	var Tests = []struct {
		icq config.Icq
		err error
	}{
		{config.Icq{Token: "foo", Uid: "bar"}, nil},
		{config.Icq{Token: "foo"}, expectedError},
		{config.Icq{Uid: "bar"}, expectedError},
		{config.Icq{}, expectedError},
	}

	for _, tt := range Tests {
		c := &config.Config{}
		c.Handler.Icq = tt.icq
		if err := s.Init(c); !reflect.DeepEqual(err, tt.err) {
			t.Fatalf("Init(): %v", err)
		}
	}
}
