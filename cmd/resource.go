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

// resourceConfigCmd represents the resource subcommand
var resourceConfigCmd = &cobra.Command{
	Use:   "resource FLAG",
	Short: "specific resources to be watched",
	Long:  `specific resources to be watched`,
	Run: func(cmd *cobra.Command, args []string) {
		conf, err := config.New()
		if err != nil {
			logrus.Fatal(err)
		}

		var b bool
		b, err = cmd.Flags().GetBool("svc")
		if err == nil {
			conf.Resource.Services = b
		} else {
			logrus.Fatal("svc", err)
		}

		b, err = cmd.Flags().GetBool("deployments")
		if err == nil {
			conf.Resource.Deployment = b
		} else {
			logrus.Fatal("deployments", err)
		}

		b, err = cmd.Flags().GetBool("po")
		if err == nil {
			conf.Resource.Pod = b
		} else {
			logrus.Fatal("po", err)
		}

		b, err = cmd.Flags().GetBool("rs")
		if err == nil {
			conf.Resource.ReplicaSet = b
		} else {
			logrus.Fatal("rs", err)
		}

		b, err = cmd.Flags().GetBool("rc")
		if err == nil {
			conf.Resource.ReplicationController = b
		} else {
			logrus.Fatal("rc", err)
		}

		b, err = cmd.Flags().GetBool("ns")
		if err == nil {
			conf.Resource.Namespace = b
		} else {
			logrus.Fatal("ns", err)
		}

		b, err = cmd.Flags().GetBool("jobs")
		if err == nil {
			conf.Resource.Job = b
		} else {
			logrus.Fatal("jobs", err)
		}

		b, err = cmd.Flags().GetBool("pv")
		if err == nil {
			conf.Resource.PersistentVolume = b
		} else {
			logrus.Fatal("pv", err)
		}

		b, err = cmd.Flags().GetBool("ds")
		if err == nil {
			conf.Resource.DaemonSet = b
		} else {
			logrus.Fatal("ds", err)
		}

		b, err = cmd.Flags().GetBool("secret")
		if err == nil {
			conf.Resource.Secret= b
		} else {
			logrus.Fatal("secret", err)
		}

		b, err = cmd.Flags().GetBool("ing")
		if err == nil {
			conf.Resource.Ingress = b
		} else {
			logrus.Fatal("ing", err)
		}


		if err = conf.Write(); err != nil {
			logrus.Fatal(err)
		}
	},
}

func init() {
	resourceConfigCmd.Flags().Bool("svc", false, "watch for services")
	resourceConfigCmd.Flags().Bool("deployments", false, "watch for deployments")
	resourceConfigCmd.Flags().Bool("po", false, "watch for pods")
	resourceConfigCmd.Flags().Bool("rc", false, "watch for replication controllers")
	resourceConfigCmd.Flags().Bool("rs", false, "watch for replicasets")
	resourceConfigCmd.Flags().Bool("ns", false, "watch for namespaces")
	resourceConfigCmd.Flags().Bool("pv", false, "watch for persistent volumes")
	resourceConfigCmd.Flags().Bool("jobs", false, "watch for jobs")
	resourceConfigCmd.Flags().Bool("ds", false, "watch for daemonsets")
	resourceConfigCmd.Flags().Bool("secret", false, "watch for plain secrets")
	resourceConfigCmd.Flags().Bool("ing", false, "watch for ingresses")
}
