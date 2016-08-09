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
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/watch"
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
}

// New create new KubewatchEvent
func New(e watch.Event) Event {
	var namespace, kind, component, host, reason, status string

	if apiEvent, ok := (e.Object).(*api.Event); ok {
		namespace = apiEvent.ObjectMeta.Namespace
		kind = apiEvent.InvolvedObject.Kind
		component = apiEvent.Source.Component
		host = apiEvent.Source.Host
		reason = apiEvent.Reason
		status = apiEvent.Type
	}

	if apiService, ok := (e.Object).(*api.Service); ok {
		namespace = apiService.ObjectMeta.Namespace
		kind = apiService.Kind
		component = string(apiService.Spec.Type)
		reason = string(e.Type)
		status = "Normal"
	}

	kbEvent := Event{
		Namespace: namespace,
		Kind:      kind,
		Component: component,
		Host:      host,
		Reason:    reason,
		Status:    status,
	}

	return kbEvent
}
