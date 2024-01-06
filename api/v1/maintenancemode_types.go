/*
Copyright 2024.

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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// MaintenanceModeSpec defines the desired state of MaintenanceMode
type MaintenanceModeSpec struct {
	// Important: Run "make" to regenerate code after modifying this file

	// The desired state of maintenance mode
	// When enabled, the controller will scale down all pods that are using the storage class specified in the spec.
	Enable bool `json:"enable"`

	// The name of the storage class. Defaults to all if not specified.
	StorageClasses []string `json:"storageClass,omitempty"`
}

// MaintenanceModeStatus defines the observed state of MaintenanceMode
type MaintenanceModeStatus struct {
	// Important: Run "make" to regenerate code after modifying this file

	// Specifies the status of the MaintenanceMode.
	// Valid values are:
	// - "Pending" (default): the controller has not processed the request to enable/disable maintenance mode;
	// - "InProgress": the controller is processing the request to enable/disable maintenance mode;
	// - "Ready": the state is synced to enable/disable maintenance mode.
	// State State `json:"state,omitempty"`

	// Specifies the list of deployments that are currently being targeted.
	// Targets []appsv1.Deployment `json:"targets,omitempty"`
}

// +kubebuilder:validation:Enum=Pending;InProgress;Ready;Conflicted

// State is the current status of maintenance mode
type State string

const (
	PendingState    State = "Pending"
	InProgressState State = "InProgress"
	ReadyState      State = "Ready"
	ConflictedState State = "Conflicted"
)

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// MaintenanceMode is the Schema for the maintenancemodes API
type MaintenanceMode struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MaintenanceModeSpec   `json:"spec,omitempty"`
	Status MaintenanceModeStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// MaintenanceModeList contains a list of MaintenanceMode
type MaintenanceModeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MaintenanceMode `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MaintenanceMode{}, &MaintenanceModeList{})
}
