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
	"github.com/skippbox/kubewatch/pkg/handlers/alertmanager"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// alertManagerCmd represents the alert manager subcommand
var alertManagerCmd = &cobra.Command{
	Use:   "alertmanager SUBCOMMAND",
	Short: "alertmanager runs kubewatch using the alertmanager handler",
	Long: `alertmanager command allows you to run kubewatch using the alertmanager handler. this is used when
	intergrating with prometheus monitoring systems.`,
	Run: func(cmd *cobra.Command, args []string) {
		conf := &config.Config{}
		if err := conf.Load(); err != nil {
			logrus.Fatal(err)
		}
		url := viper.GetString(urlKey)
		labels := viper.GetStringSlice(labelKey)
		handler := alertmanager.New(conf, url, labels)
		client.Run(handler)
	},
}

func init() {
	viper.SetEnvPrefix("alertmanager")
	alertManagerCmd.Flags().StringP(urlKey, "u", "", "Specify alertmanager host url")
	alertManagerCmd.Flags().StringSliceP(labelKey, "l", []string{}, "Add one or more labels to the alertmanager event, --label mylabel=blah")
	viper.BindPFlag(urlKey, alertManagerCmd.Flags().Lookup(urlKey))
	viper.BindPFlag(labelKey, alertManagerCmd.Flags().Lookup(labelKey))
}
