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

package exec

import (
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/bitnami-labs/kubewatch/config"
	"github.com/bitnami-labs/kubewatch/pkg/event"
)

// Exec handler implements handler.Handler interface,
// Notify event to run local command/shell
type Exec struct {
	Cmd string
}

// Init prepares Exec configuration
func (s *Exec) Init(c *config.Config) error {
	cmd := c.Handler.Exec.Cmd

	if cmd == "" {
		cmd = os.Getenv("KW_EXEC_CMD")
	}

	s.Cmd = cmd
	return nil
}

func (s *Exec) ObjectCreated(obj interface{}) {
	handleEvent(s, obj, "created")
}

func (s *Exec) ObjectDeleted(obj interface{}) {
	handleEvent(s, obj, "deleted")
}

func (s *Exec) ObjectUpdated(oldObj, newObj interface{}) {
	handleEvent(s, newObj, "updated")
}

func handleEvent(s *Exec, obj interface{}, action string) {
	e := event.New(obj, action)

	var cmdAndArgs []string
	cmdCli := strings.Split(s.Cmd, " ")
	for _, c := range cmdCli {
		c = strings.Replace(c, "{reason}", e.Reason, -1)
		c = strings.Replace(c, "{name}", e.Name, -1)

		cmdAndArgs = append(cmdAndArgs, c)
	}

	if err := runCmd(cmdAndArgs); err != nil {
		log.Printf("Exec local command %v failed %s", cmdAndArgs, err)
	} else {
		log.Printf("Exec local command successfully %s", cmdAndArgs)
	}
}

func runCmd(command []string) error {
	cmd := exec.Command(command[0], command[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
