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
	"github.com/skippbox/kubewatch/config"
	"github.com/skippbox/kubewatch/pkg/client"
	"github.com/skippbox/kubewatch/pkg/handlers/slack"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// slackCmd represents the slack subcommand
var slackCmd = &cobra.Command{
	Use:   "slack SUBCOMMAND",
	Short: "slack runs kubewatch using the slack handler",
	Long: `slack command allows you to run kubewatch using the slack handler. this is used when
	intergrating with slack channels.`,
	Run: func(cmd *cobra.Command, args []string) {
		conf, err := config.New()
		if err != nil {
			logrus.Fatal(err)
		}
		channel := viper.GetString(channelKey)
		token := viper.GetString(tokenKey)
		handler := slack.New(conf, channel, token)
		client.Run(handler)

	},
}

func init() {
	viper.SetEnvPrefix("slack")
	slackCmd.Flags().StringP(channelKey, "c", "", "Specify slack channel")
	slackCmd.Flags().StringP(tokenKey, "t", "", "Specify slack token")
	viper.BindPFlag(channelKey, slackCmd.Flags().Lookup(channelKey))
	viper.BindPFlag(tokenKey, slackCmd.Flags().Lookup(tokenKey))
}
