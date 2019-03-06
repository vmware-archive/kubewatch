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

// hipchatConfigCmd represents the hipchat subcommand
var hipchatConfigCmd = &cobra.Command{
	Use:   "hipchat",
	Short: "specific hipchat configuration",
	Long:  `specific hipchat configuration`,
	Run: func(cmd *cobra.Command, args []string) {
		conf, err := config.New()
		if err != nil {
			logrus.Fatal(err)
		}

		token, err := cmd.Flags().GetString("token")
		if err == nil {
			if len(token) > 0 {
				conf.Handler.Hipchat.Token = token
			}
		} else {
			logrus.Fatal(err)
		}
		room, err := cmd.Flags().GetString("room")
		if err == nil {
			if len(room) > 0 {
				conf.Handler.Hipchat.Room = room
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
	hipchatConfigCmd.Flags().StringP("room", "r", "", "Specify hipchat room")
	hipchatConfigCmd.Flags().StringP("token", "t", "", "Specify hipchat token")
	hipchatConfigCmd.Flags().StringP("url", "u", "", "Specify hipchat server url")
}
