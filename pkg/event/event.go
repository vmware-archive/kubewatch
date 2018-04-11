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

	apps_v1beta1 "k8s.io/api/apps/v1beta1"
	batch_v1 "k8s.io/api/batch/v1"
	api_v1 "k8s.io/api/core/v1"
	ext_v1beta1 "k8s.io/api/extensions/v1beta1"
	api_meta "k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
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
	Labels    map[string]string
}

var m = map[string]string{
	"created": "Normal",
	"deleted": "Danger",
	"updated": "Warning",
}

// New create new KubewatchEvent
func New(obj interface{}, action string) Event {
	var namespace, kind, component, host, reason, status, name string
	if apiService, ok := obj.(*api_v1.Service); ok {
		namespace = apiService.ObjectMeta.Namespace
		name = apiService.Name
		kind = "service"
		component = string(apiService.Spec.Type)
		reason = action
		status = m[action]
	} else if apiNamespace, ok := obj.(*api_v1.Namespace); ok {
		name = apiNamespace.Name
		kind = "namespace"
		reason = action
		status = m[action]
	} else if apiPod, ok := obj.(*api_v1.Pod); ok {
		namespace = apiPod.ObjectMeta.Namespace
		name = apiPod.Name
		kind = "pod"
		reason = action
		host = apiPod.Spec.NodeName
		status = m[action]
	} else if apiRC, ok := obj.(*api_v1.ReplicationController); ok {
		namespace = apiRC.ObjectMeta.Namespace
		name = apiRC.Name
		kind = "replication controller"
		reason = action
		status = m[action]
	} else if apiDeployment, ok := obj.(*apps_v1beta1.Deployment); ok {
		namespace = apiDeployment.ObjectMeta.Namespace
		name = apiDeployment.Name
		kind = "deployment"
		reason = action
		status = m[action]
	} else if apiJob, ok := obj.(*batch_v1.Job); ok {
		namespace = apiJob.ObjectMeta.Namespace
		name = apiJob.Name
		kind = "job"
		reason = action
		status = m[action]
	} else if apiPV, ok := obj.(*api_v1.PersistentVolume); ok {
		name = apiPV.Name
		kind = "persistent volume"
		reason = action
		status = m[action]
	} else if apiDS, ok := obj.(*ext_v1beta1.DaemonSet); ok {
		namespace = apiDS.ObjectMeta.Namespace
		name = apiDS.Name
		kind = "daemon set"
		reason = action
		status = m[action]
	} else if apiRS, ok := obj.(*ext_v1beta1.ReplicaSet); ok {
		namespace = apiRS.ObjectMeta.Namespace
		name = apiRS.Name
		kind = "replica set"
		reason = action
		status = m[action]
	}

	accessor := api_meta.NewAccessor()
	labels, _ := accessor.Labels(obj.(runtime.Object))

	kbEvent := Event{
		Namespace: namespace,
		Kind:      kind,
		Component: component,
		Host:      host,
		Reason:    reason,
		Status:    status,
		Name:      name,
		Labels:    labels,
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
			"A namespace %s has been %s",
			e.Name,
			e.Reason,
		)
	default:
		msg = fmt.Sprintf(
			"A %s in namespace %s has been %s: %s",
			e.Kind,
			e.Namespace,
			e.Reason,
			e.Name,
		)
	}
	return msg
}
