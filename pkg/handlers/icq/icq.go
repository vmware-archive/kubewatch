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
	"log"
	"os"
	"time"

	"github.com/bitnami-labs/kubewatch/config"
	kbEvent "github.com/bitnami-labs/kubewatch/pkg/event"
	icq "gopkg.in/icq.v1"
)

var icqErrMsg = `
%s

You need to set both icq token and uid for icq notify,
using "--token/-t" and "--uid/-u", or using environment variables:

export KW_ICQ_TOKEN=icq_token
export KW_ICQ_UID=icq_uid

Command line flags will override environment variables

`

// Icq handler implements handler.Handler interface,
// Notify event to icq uid
type Icq struct {
	Token string
	Uid   string
}

type IcqMessage struct {
	Text string `json:"text"`
}

// Init prepares Icq configuration
func (s *Icq) Init(c *config.Config) error {
	token := c.Handler.Icq.Token
	uid := c.Handler.Icq.Uid

	if token == "" {
		token = os.Getenv("KW_ICQ_TOKEN")
	}

	if uid == "" {
		uid = os.Getenv("KW_ICQ_UID")
	}

	s.Token = token
	s.Uid = uid

	return checkMissingIcqVars(s)
}

func (s *Icq) ObjectCreated(obj interface{}) {
	notifyIcq(s, obj, "created")
}

func (s *Icq) ObjectDeleted(obj interface{}) {
	notifyIcq(s, obj, "deleted")
}

func (s *Icq) ObjectUpdated(oldObj, newObj interface{}) {
	notifyIcq(s, newObj, "updated")
}

func notifyIcq(s *Icq, obj interface{}, action string) {
	e := kbEvent.New(obj, action)
	api := icq.NewAPI(s.Token)

	icqMessage := e.Message()

	r, err := api.SendMessage(s.Uid, icqMessage)
	if err != nil {
		log.Printf("%s\n", err)
		return
	}

	log.Printf("Message successfully sent to channel %s with state %s at %s", s.Uid, r.State, time.Now())
}

func checkMissingIcqVars(s *Icq) error {
	if s.Token == "" || s.Uid == "" {
		return fmt.Errorf(icqErrMsg, "Missing Icq token or uid")
	}

	return nil
}
