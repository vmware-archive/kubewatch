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

package event

import (
	"fmt"

	"github.com/bitnami-labs/kubewatch/pkg/utils"
	apps_v1 "k8s.io/api/apps/v1"
	batch_v1 "k8s.io/api/batch/v1"
	api_v1 "k8s.io/api/core/v1"
	ext_v1beta1 "k8s.io/api/extensions/v1beta1"
	rbac_v1beta1 "k8s.io/api/rbac/v1beta1"
)

// Event represent an event got from k8s api server
// Events from different endpoints need to be casted to KubewatchEvent
// before being able to be handled by handler
type Event struct {
	Namespace string
	Kind      string
	Component string
	Host      string
	Reason    string
	Status    string
	Name      string
}

var m = map[string]string{
	"created": "Normal",
	"deleted": "Danger",
	"updated": "Warning",
}

// New create new KubewatchEvent
func New(obj interface{}, action string) Event {
	var namespace, kind, component, host, reason, status, name string

	objectMeta := utils.GetObjectMetaData(obj)
	namespace = objectMeta.Namespace
	name = objectMeta.Name
	reason = action
	status = m[action]

	switch object := obj.(type) {
	case *ext_v1beta1.DaemonSet:
		kind = "daemon set"
	case *apps_v1.Deployment:
		kind = "deployment"
	case *apps_v1.StatefulSet:
		kind = "stateful set"
	case *batch_v1.Job:
		kind = "job"
	case *api_v1.Namespace:
		kind = "namespace"
	case *ext_v1beta1.Ingress:
		kind = "ingress"
	case *api_v1.PersistentVolume:
		kind = "persistent volume"
	case *api_v1.Pod:
		kind = "pod"
		host = object.Spec.NodeName
	case *api_v1.ReplicationController:
		kind = "replication controller"
	case *ext_v1beta1.ReplicaSet:
		kind = "replica set"
	case *api_v1.Service:
		kind = "service"
		component = string(object.Spec.Type)
	case *api_v1.Secret:
		kind = "secret"
	case *api_v1.ConfigMap:
		kind = "configmap"
	case *api_v1.Node:
		kind = "node"
	case *rbac_v1beta1.ClusterRole:
		kind = "cluster role"
	case *api_v1.ServiceAccount:
		kind = "service account"
	case Event:
		name = object.Name
		kind = object.Kind
		namespace = object.Namespace
	}

	kbEvent := Event{
		Namespace: namespace,
		Kind:      kind,
		Component: component,
		Host:      host,
		Reason:    reason,
		Status:    status,
		Name:      name,
	}
	return kbEvent
}

// Message returns event message in standard format.
// included as a part of event packege to enhance code resuablity across handlers.
func (e *Event) Message() (msg string) {
	// using switch over if..else, since the format could vary based on the kind of the object in future.
	switch e.Kind {
	case "namespace":
		msg = fmt.Sprintf(
			"A namespace `%s` has been `%s`",
			e.Name,
			e.Reason,
		)
	case "node":
		msg = fmt.Sprintf(
			"A node `%s` has been `%s`",
			e.Name,
			e.Reason,
		)
	case "cluster role":
		msg = fmt.Sprintf(
			"A cluster role `%s` has been `%s`",
			e.Name,
			e.Reason,
		)
	case "NodeReady":
		msg = fmt.Sprintf(
			"Node `%s` is Ready : \nNodeReady",
			e.Name,
		)
	case "NodeNotReady":
		msg = fmt.Sprintf(
			"Node `%s` is Not Ready : \nNodeNotReady",
			e.Name,
		)
	case "NodeRebooted":
		msg = fmt.Sprintf(
			"Node `%s` Rebooted : \nNodeRebooted",
			e.Name,
		)
	case "Backoff":
		msg = fmt.Sprintf(
			"Pod `%s` in `%s` Crashed : \nCrashLoopBackOff %s",
			e.Name,
			e.Namespace,
			e.Reason,
		)
	default:
		msg = fmt.Sprintf(
			"A `%s` in namespace `%s` has been `%s`:\n`%s`",
			e.Kind,
			e.Namespace,
			e.Reason,
			e.Name,
		)
	}
	return msg
}
