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
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/apis/extensions/v1beta1"
	"k8s.io/client-go/pkg/apis/batch"
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
	if apiService, ok := obj.(*v1.Service); ok {
		namespace = apiService.ObjectMeta.Namespace
		name = apiService.Name
		kind = "service"
		component = string(apiService.Spec.Type)
		reason = action
		status = m[action]
	} else if apiNamespace, ok := obj.(*v1.Namespace); ok {
		name = apiNamespace.Name
		kind = "namespace"
		reason = action
		status = m[action]
	} else if apiPod, ok := obj.(*v1.Pod); ok {
		namespace = apiPod.ObjectMeta.Namespace
		name = apiPod.Name
		kind = "pod"
		reason = action
		host = apiPod.Spec.NodeName
		status = m[action]
	} else if apiRC, ok := obj.(*v1.ReplicationController); ok {
		namespace = apiRC.ObjectMeta.Namespace
		name = apiRC.Name
		kind = "replication controller"
		reason = action
		status = m[action]
	} else if apiDeployment, ok := obj.(*v1beta1.Deployment); ok {
		namespace = apiDeployment.ObjectMeta.Namespace
		name = apiDeployment.Name
		kind = "deployment"
		reason = action
		status = m[action]
	} else if apiJob, ok := obj.(*batch.Job); ok {
		namespace = apiJob.ObjectMeta.Namespace
		name = apiJob.Name
		kind = "job"
		reason = action
		status = m[action]
	} else if apiPV, ok := obj.(*v1.PersistentVolume); ok {
		name = apiPV.Name
		kind = "persistent volume"
		reason = action
		status = m[action]
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
