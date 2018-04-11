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

package cmd

import (
	"github.com/Sirupsen/logrus"
	"github.com/bitnami-labs/kubewatch/config"
	"github.com/spf13/cobra"
)

// custodianConfigCmd represents the custodian subcommand
var custodianConfigCmd = &cobra.Command{
	Use:   "custodian FLAG",
	Short: "specific custodian configuration",
	Long:  `specific custodian configuration`,
	Run: func(cmd *cobra.Command, args []string) {
		conf, err := config.New()
		if err != nil {
			logrus.Fatal(err)
		}

		foo, err := cmd.Flags().GetString("foo")
		if err == nil {
			if len(foo) > 0 {
				conf.Handler.Custodian.Foo = foo
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
	custodianConfigCmd.Flags().StringP("foo", "c", "", "Specify custodian foo")
}
