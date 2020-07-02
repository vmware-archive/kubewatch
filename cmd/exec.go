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

package cmd

import (
	"github.com/Sirupsen/logrus"
	"github.com/bitnami-labs/kubewatch/config"
	"github.com/spf13/cobra"
)

// execConfigCmd represents the exec subcommand
var execConfigCmd = &cobra.Command{
	Use:   "exec FLAG",
	Short: "specific exec configuration",
	Long:  `specific exec configuration`,
	Run: func(cmd *cobra.Command, args []string) {
		conf, err := config.New()
		if err != nil {
			logrus.Fatal(err)
		}

		execCmd, err := cmd.Flags().GetString("cmd")
		if err == nil {
			if len(execCmd) > 0 {
				conf.Handler.Exec.Cmd = execCmd
			}
		} else {
			logrus.Fatal(err)
		}

		if err = conf.Write(); err != nil {
			logrus.Fatal(err)
		}
	},
}

func init() {
	execConfigCmd.Flags().StringP("cmd", "c", "", "Specify local command/shell to run")
}
