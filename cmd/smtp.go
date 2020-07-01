/*
Copyright 2020 VMware

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
	"os"

	"github.com/bitnami-labs/kubewatch/pkg/handlers/smtp"
	"github.com/spf13/cobra"
)

// smtpConfigCmd represents the smtp subcommand
var smtpConfigCmd = &cobra.Command{
	Use:   "smtp",
	Short: "specific smtp configuration",
	Long:  `specific smtp configuration`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(os.Stderr, "CLI setters not implemented yet, please edit ~/.kubewatch.yaml directly. Example:\n\n%s", smtp.ConfigExample)
	},
}

func init() {
}
