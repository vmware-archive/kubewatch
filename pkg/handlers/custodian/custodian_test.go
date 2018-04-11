/*
Copyright 2018 Bitnami Inc.

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

package custodian

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/bitnami-labs/kubewatch/config"
)

func TestCustodianInit(t *testing.T) {
	c := &Custodian{}
	expectedError := fmt.Errorf(custodianErrMsg, "Missing custodian foo")

	var Tests = []struct {
		custodian config.Custodian
		err       error
	}{
		{config.Custodian{Foo: "123"}, nil},
		// {config.Custodian{Token: "foo"}, expectedError},
		// {config.Custodian{Channel: "bar"}, expectedError},
		{config.Custodian{}, expectedError},
	}

	for _, tt := range Tests {
		conf := &config.Config{}
		conf.Handler.Custodian = tt.custodian
		if err := c.Init(conf); !reflect.DeepEqual(err, tt.err) {
			t.Fatalf("Init(): %v != %v", err, tt.err)
		}
	}
}
