/*
Copyright 2018 Bitnami

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
	"github.com/bitnami-labs/kubewatch/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// webhookConfigCmd represents the webhook subcommand
var slackwebhookConfigCmd = &cobra.Command{
	Use:   "slackwebhook",
	Short: "specific Slack webhook configuration",
	Long:  `specific Slack webhook configuration`,
	Run: func(cmd *cobra.Command, args []string) {
		conf, err := config.New()
		if err != nil {
			logrus.Fatal(err)
		}

		channel, err := cmd.Flags().GetString("channel")
		if err == nil {
			if len(channel) > 0 {
				conf.Handler.SlackWebhook.Channel = channel
			}
		} else {
			logrus.Fatal(err)
		}

		username, err := cmd.Flags().GetString("username")
		if err == nil {
			if len(username) > 0 {
				conf.Handler.SlackWebhook.Username = username
			}
		} else {
			logrus.Fatal(err)
		}

		emoji, err := cmd.Flags().GetString("emoji")
		if err == nil {
			if len(emoji) > 0 {
				conf.Handler.SlackWebhook.Emoji = emoji
			}
		}

		slackwebhookurl, err := cmd.Flags().GetString("slackwebhookurl")
		if err == nil {
			if len(slackwebhookurl) > 0 {
				conf.Handler.SlackWebhook.Slackwebhookurl = slackwebhookurl
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
	webhookConfigCmd.Flags().StringP("channel", "c", "", "Specify Slack Webhook url Channel")
	webhookConfigCmd.Flags().StringP("username", "n", "", "Specify Slack Webhook url Username")
	webhookConfigCmd.Flags().StringP("emoji", "e", "", "Specify Slack Webhook url Emoji")
	webhookConfigCmd.Flags().StringP("slackwebhookurl", "w", "", "Specify Slack Webhook url")
}
