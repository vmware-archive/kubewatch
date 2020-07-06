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
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/bitnami-labs/kubewatch/config"
	"github.com/bitnami-labs/kubewatch/pkg/client"
	"github.com/bitnami-labs/kubewatch/pkg/event"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const kubewatchConfigFile = ".kubewatch.yaml"

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "modify kubewatch configuration",
	Long: `
config command allows configuration of ~/.kubewatch.yaml for running kubewatch`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var configAddCmd = &cobra.Command{
	Use:   "add",
	Short: "add webhook config to ~/.kubewatch.yaml",
	Long: `
Adds webhook config to ~/.kubewatch.yaml`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var configTestCmd = &cobra.Command{
	Use:   "test",
	Short: "test handler config present in ~/.kubewatch.yaml",
	Long: `
Tests handler configs present in ~/.kubewatch.yaml by sending test messages`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Testing Handler configs from .kubewatch.yaml")
		conf, err := config.New()
		if err != nil {
			logrus.Fatal(err)
		}
		eventHandler := client.ParseEventHandler(conf)
		e := event.Event{
			Namespace: "testNamespace",
			Name:      "testResource",
			Kind:      "testKind",
			Component: "testComponent",
			Host:      "testHost",
			Reason:    "Tested",
			Status:    "Normal",
		}
		eventHandler.Handle(e)
	},
}

var configSampleCmd = &cobra.Command{
	Use:   "sample",
	Short: "Show a sample config file",
	Long: `
Print a sample config file which can be put in ~/.kubewatch.yaml`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print(config.ConfigSample)
	},
}

var configViewCmd = &cobra.Command{
	Use:   "view",
	Short: "view ~/.kubewatch.yaml",
	Long: `
Display the contents of the contents of ~/.kubewatch.yaml`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintln(os.Stderr, "Contents of ~/.kubewatch.yaml")
		configFile, err := ioutil.ReadFile(filepath.Join(os.Getenv("HOME"), kubewatchConfigFile))
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
		fmt.Print(string(configFile))
	},
}

func init() {
	RootCmd.AddCommand(configCmd)
	configCmd.AddCommand(
		configAddCmd,
		configTestCmd,
		configSampleCmd,
		configViewCmd,
	)

	configAddCmd.AddCommand(
		slackConfigCmd,
		hipchatConfigCmd,
		mattermostConfigCmd,
		flockConfigCmd,
		webhookConfigCmd,
		msteamsConfigCmd,
		smtpConfigCmd,
	)
}
