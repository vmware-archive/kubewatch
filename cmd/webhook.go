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
var webhookConfigCmd = &cobra.Command{
	Use:   "webhook",
	Short: "specific webhook configuration",
	Long:  `specific webhook configuration`,
	Run: func(cmd *cobra.Command, args []string) {
		conf, err := config.New()
		if err != nil {
			logrus.Fatal(err)
		}

		url, err := cmd.Flags().GetString("url")
		if err == nil {
			if len(url) > 0 {
				conf.Handler.Webhook.Url = url
			}
		} else {
			logrus.Fatal(err)
		}

		hmacKey, err := cmd.Flags().GetString("hmac-key")
		if err != nil {
			logrus.Fatal(err)
		}
		if hmacKey != "" {
			conf.Handler.Webhook.HMACKey = hmacKey
		}

		hmacSignatureHeader, err := cmd.Flags().GetString("hmac-signature-header")
		if err != nil {
			logrus.Fatal(err)
		}
		if hmacSignatureHeader != "" {
			conf.Handler.Webhook.HMACSignatureHeader = hmacSignatureHeader
		}

		if conf.Handler.Webhook.HMACSignatureHeader != "" && conf.Handler.Webhook.HMACKey == "" {
			logrus.Warn("HMAC signature header is set but HMAC key is empty")
		}

		if err = conf.Write(); err != nil {
			logrus.Fatal(err)
		}
	},
}

func init() {
	webhookConfigCmd.Flags().StringP("url", "u", "", "Specify Webhook url")
	webhookConfigCmd.Flags().String("hmac-key", "", "A base64 encoded string to generate a webhook signature with")
	webhookConfigCmd.Flags().String("hmac-signature-header", "", "The name of the header to set the hmac signature value to")
}
