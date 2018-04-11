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
	"log"
	"os"

	// "github.com/Sirupsen/logrus"
	"github.com/bitnami-labs/kubewatch/config"
	// "github.com/bitnami-labs/kubewatch/pkg/event"
	kbEvent "github.com/bitnami-labs/kubewatch/pkg/event"
)

// Custodian handler implements handler.Handler interface,
type Custodian struct {
	Foo string
	Bar string
}

var custodianErrMsg = "booo %s"

// Init prepares custodian configuration
func (c *Custodian) Init(conf *config.Config) error {
	foo := conf.Handler.Custodian.Foo

	if foo == "" {
		foo = os.Getenv("KW_CUSTODIAN_FOO")
	}
	c.Foo = foo
	return checkMissingCustodianVars(c)
}

func (c *Custodian) ObjectCreated(obj interface{}) {
	runCustodian(c, obj, "created")
}

func (c *Custodian) ObjectDeleted(obj interface{}) {
	runCustodian(c, obj, "deleted")
}

func (c *Custodian) ObjectUpdated(oldObj, newObj interface{}) {
	runCustodian(c, newObj, "updated")
}

func runCustodian(c *Custodian, obj interface{}, action string) {
	e := kbEvent.New(obj, action)
	log.Printf("jjo: Got event: %v, event %v", e.Labels, e)
}

func checkMissingCustodianVars(c *Custodian) error {
	if c.Foo == "" {
		return fmt.Errorf(custodianErrMsg, "Missing custodian foo")
	}

	return nil
}
