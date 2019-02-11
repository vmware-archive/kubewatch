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

// rocketChatCmd represents the rocketChat subcommand
var rocketchatConfigCmd = &cobra.Command{
	Use:   "rocketchat FLAG",
	Short: "specific rocketchat configuration",
	Long:  `specific rocketchat configuration`,
	Run: func(cmd *cobra.Command, args []string) {
		conf, err := config.New()
		if err != nil {
			logrus.Fatal(err)
		}

		host, err := cmd.Flags().GetString("host")
		if err == nil {
			if len(host) > 0 {
				conf.Handler.Rocketchat.Host = host
			}
		} else {
			logrus.Fatal(err)
		}
		port, err := cmd.Flags().GetInt("port")
		if err == nil {

			conf.Handler.Rocketchat.Port = port

		} else {
			logrus.Fatal(err)
		}

		user, err := cmd.Flags().GetString("user")
		if err == nil {
			if len(user) > 0 {
				conf.Handler.Rocketchat.User = user
			}
		} else {
			logrus.Fatal(err)
		}
		password, err := cmd.Flags().GetString("password")
		if err == nil {
			if len(password) > 0 {
				conf.Handler.Rocketchat.Host = password
			}
		} else {
			logrus.Fatal(err)
		}

		scheme, err := cmd.Flags().GetString("scheme")
		if err == nil {
			if len(scheme) > 0 {
				conf.Handler.Rocketchat.Scheme = scheme
			}
		} else {
			logrus.Fatal(err)
		}

		channel, err := cmd.Flags().GetString("channel")
		if err == nil {
			if len(channel) > 0 {
				conf.Handler.Rocketchat.Channel = channel
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
	rocketchatConfigCmd.Flags().StringP("scheme", "", "", "Specify rocketchat host scheme eg http")
	rocketchatConfigCmd.Flags().StringP("host", "", "", "Specify rocketchat host rocketchat.com")
	rocketchatConfigCmd.Flags().IntP("port", "", 80, "Specify rocketchat port")
	rocketchatConfigCmd.Flags().StringP("user", "", "kubewatch", "Specify rocketchat user")
	rocketchatConfigCmd.Flags().StringP("password", "", "", "Specify rocketchat password")
	rocketchatConfigCmd.Flags().StringP("channel", "", "", "Specify rocketchat channel")
	rocketchatConfigCmd.Flags().StringP("email", "", "", "Specify rocketchat user email")

}
