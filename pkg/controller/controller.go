package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/client/cache"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	"k8s.io/kubernetes/pkg/controller/framework"
	"k8s.io/kubernetes/pkg/fields"
	cmdutil "k8s.io/kubernetes/pkg/kubectl/cmd/util"
	"k8s.io/kubernetes/pkg/util/wait"
)

func Controller() {

	factory := cmdutil.NewFactory(nil)
	kubeClient, err := factory.ClientConfig()
	if err != nil {
		logrus.Fatal(err)
	}

	var podsStore cache.Store
	podsStore = watchPods(kubeClient, podsStore)

	var servicesStore cache.Store
	servicesStore = watchServices(kubeClient, servicesStore)

	var rcStore cache.Store
	rcStore = watchReplicationControllers(kubeClient, rcStore)

	logrus.Fatal(http.ListenAndServe(":8080", nil))
}

func watchPods(client *client.Client, store cache.Store) cache.Store {
	//Define what we want to look for (Pods)
	watchlist := cache.NewListWatchFromClient(client, "pods", api.NamespaceAll, fields.Everything())

	resyncPeriod := 30 * time.Minute

	//Setup an informer to call functions when the watchlist changes
	eStore, eController := framework.NewInformer(
		watchlist,
		&api.Pod{},
		resyncPeriod,
		framework.ResourceEventHandlerFuncs{
			AddFunc:    podCreated,
			DeleteFunc: podDeleted,
		},
	)

	//Run the controller as a goroutine
	go eController.Run(wait.NeverStop)

	return eStore
}

func podCreated(obj interface{}) {
	pod := obj.(*api.Pod)
	fmt.Println("Pod created: " + pod.ObjectMeta.Name)
}

func podDeleted(obj interface{}) {
	pod := obj.(*api.Pod)
	fmt.Println("Pod deleted: " + pod.ObjectMeta.Name)
}

func watchServices(client *client.Client, store cache.Store) cache.Store {
	//Define what we want to look for (Services)
	watchlist := cache.NewListWatchFromClient(client, "services", api.NamespaceAll, fields.Everything())

	resyncPeriod := 30 * time.Minute

	//Setup an informer to call functions when the watchlist changes
	eStore, eController := framework.NewInformer(
		watchlist,
		&api.Service{},
		resyncPeriod,
		framework.ResourceEventHandlerFuncs{
			AddFunc:    serviceCreated,
			DeleteFunc: serviceDeleted,
		},
	)

	//Run the controller as a goroutine
	go eController.Run(wait.NeverStop)

	return eStore
}

func serviceCreated(obj interface{}) {
	service := obj.(*api.Service)
	fmt.Println("Service created: " + service.ObjectMeta.Name)
}

func serviceDeleted(obj interface{}) {
	service := obj.(*api.Service)
	fmt.Println("Service deleted: " + service.ObjectMeta.Name)
}

func watchReplicationControllers(client *client.Client, store cache.Store) cache.Store {
	//Define what we want to look for (ReplicationControllers)
	watchlist := cache.NewListWatchFromClient(client, "replicationcontrollers", api.NamespaceAll, fields.Everything())

	resyncPeriod := 30 * time.Minute

	//Setup an informer to call functions when the watchlist changes
	eStore, eController := framework.NewInformer(
		watchlist,
		&api.ReplicationController{},
		resyncPeriod,
		framework.ResourceEventHandlerFuncs{
			AddFunc:    replicationcontrollerCreated,
			DeleteFunc: replicationcontrollerDeleted,
		},
	)

	//Run the controller as a goroutine
	go eController.Run(wait.NeverStop)

	return eStore
}

func replicationcontrollerCreated(obj interface{}) {
	replicationcontroller := obj.(*api.ReplicationController)
	fmt.Println("ReplicationController created: " + replicationcontroller.ObjectMeta.Name)
}

func replicationcontrollerDeleted(obj interface{}) {
	replicationcontroller := obj.(*api.ReplicationController)
	fmt.Println("ReplicationController deleted: " + replicationcontroller.ObjectMeta.Name)
}
