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

package controller

import (
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/skippbox/kubewatch/config"
	"github.com/skippbox/kubewatch/pkg/handlers"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/apis/extensions/v1beta1"
	"k8s.io/kubernetes/pkg/client/cache"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	"k8s.io/kubernetes/pkg/controller/framework"
	"k8s.io/kubernetes/pkg/fields"
	cmdutil "k8s.io/kubernetes/pkg/kubectl/cmd/util"
	"k8s.io/kubernetes/pkg/util/wait"
)

func Controller(conf *config.Config, eventHandler handlers.Handler) {

	factory := cmdutil.NewFactory(nil)
	kubeConfig, err := factory.ClientConfig()
	if err != nil {
		logrus.Fatal(err)
	}

	kubeClient := client.NewOrDie(kubeConfig)
	kubeExtensionsClient := client.NewExtensionsOrDie(kubeConfig)

	if conf.Resource.Pod {
		watchPods(kubeClient, eventHandler)
	}

	if conf.Resource.Services {
		watchServices(kubeClient, eventHandler)
	}

	if conf.Resource.ReplicationController {
		watchReplicationControllers(kubeClient, eventHandler)
	}

	if conf.Resource.Deployment {
		watchDeployments(kubeExtensionsClient, eventHandler)
	}

	if conf.Resource.Job {
		watchJobs(kubeExtensionsClient, eventHandler)
	}

	if conf.Resource.PersistentVolume {
		var servicesStore cache.Store
		servicesStore = watchPersistenVolumes(kubeClient, servicesStore, eventHandler)
	}

	logrus.Fatal(http.ListenAndServe(":8081", nil))
}

func watchPods(client *client.Client, eventHandler handlers.Handler) cache.Store {
	//Define what we want to look for (Pods)
	watchlist := cache.NewListWatchFromClient(client, "pods", api.NamespaceAll, fields.Everything())

	resyncPeriod := 30 * time.Minute

	//Setup an informer to call functions when the watchlist changes
	eStore, eController := framework.NewInformer(
		watchlist,
		&api.Pod{},
		resyncPeriod,
		framework.ResourceEventHandlerFuncs{
			AddFunc:    eventHandler.ObjectCreated,
			DeleteFunc: eventHandler.ObjectDeleted,
		},
	)

	//Run the controller as a goroutine
	go eController.Run(wait.NeverStop)

	return eStore
}

func watchServices(client *client.Client, eventHandler handlers.Handler) cache.Store {
	//Define what we want to look for (Services)
	watchlist := cache.NewListWatchFromClient(client, "services", api.NamespaceAll, fields.Everything())

	resyncPeriod := 30 * time.Minute

	//Setup an informer to call functions when the watchlist changes
	eStore, eController := framework.NewInformer(
		watchlist,
		&api.Service{},
		resyncPeriod,
		framework.ResourceEventHandlerFuncs{
			AddFunc:    eventHandler.ObjectCreated,
			DeleteFunc: eventHandler.ObjectDeleted,
			UpdateFunc: eventHandler.ObjectUpdated,
		},
	)

	//Run the controller as a goroutine
	go eController.Run(wait.NeverStop)

	return eStore
}

func watchReplicationControllers(client *client.Client, eventHandler handlers.Handler) cache.Store {
	//Define what we want to look for (ReplicationControllers)
	watchlist := cache.NewListWatchFromClient(client, "replicationcontrollers", api.NamespaceAll, fields.Everything())

	resyncPeriod := 30 * time.Minute

	//Setup an informer to call functions when the watchlist changes
	eStore, eController := framework.NewInformer(
		watchlist,
		&api.ReplicationController{},
		resyncPeriod,
		framework.ResourceEventHandlerFuncs{
			AddFunc:    eventHandler.ObjectCreated,
			DeleteFunc: eventHandler.ObjectDeleted,
		},
	)

	//Run the controller as a goroutine
	go eController.Run(wait.NeverStop)

	return eStore
}

func watchDeployments(client *client.ExtensionsClient, eventHandler handlers.Handler) cache.Store {
	//Define what we want to look for (Deployments)
	watchlist := cache.NewListWatchFromClient(client, "deployments", api.NamespaceAll, fields.Everything())

	resyncPeriod := 30 * time.Minute

	//Setup an informer to call functions when the watchlist changes
	eStore, eController := framework.NewInformer(
		watchlist,
		&v1beta1.Deployment{},
		resyncPeriod,
		framework.ResourceEventHandlerFuncs{
			AddFunc:    eventHandler.ObjectCreated,
			DeleteFunc: eventHandler.ObjectDeleted,
		},
	)

	//Run the controller as a goroutine
	go eController.Run(wait.NeverStop)

	return eStore
}

func watchJobs(client *client.ExtensionsClient, eventHandler handlers.Handler) cache.Store {
	//Define what we want to look for (Jobs)
	watchlist := cache.NewListWatchFromClient(client, "jobs", api.NamespaceAll, fields.Everything())

	resyncPeriod := 30 * time.Minute

	//Setup an informer to call functions when the watchlist changes
	eStore, eController := framework.NewInformer(
		watchlist,
		&v1beta1.Job{},
		resyncPeriod,
		framework.ResourceEventHandlerFuncs{
			AddFunc:    eventHandler.ObjectCreated,
			DeleteFunc: eventHandler.ObjectDeleted,
		},
	)

	//Run the controller as a goroutine
	go eController.Run(wait.NeverStop)

	return eStore
}

func watchPersistenVolumes(client *client.Client, store cache.Store, eventHandler handlers.Handler) cache.Store {
	//Define what we want to look for (PersistenVolumes)
	watchlist := cache.NewListWatchFromClient(client, "persistentvolumes", api.NamespaceAll, fields.Everything())

	resyncPeriod := 30 * time.Minute

	//Setup an informer to call functions when the watchlist changes
	eStore, eController := framework.NewInformer(
		watchlist,
		&api.PersistentVolume{},
		resyncPeriod,
		framework.ResourceEventHandlerFuncs{
			AddFunc:    eventHandler.ObjectCreated,
			DeleteFunc: eventHandler.ObjectDeleted,
		},
	)

	//Run the controller as a goroutine
	go eController.Run(wait.NeverStop)

	return eStore
}
