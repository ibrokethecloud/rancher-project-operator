/*


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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ProjectSpec defines the desired state of Project
type ProjectSpec struct {
	DisplayName                   string                 `json:"displayName"`
	Description                   string                 `json:"description,omitempty"`
	ResourceQuota                 ProjectResourceQuota   `json:"resourceQuota,omitempty"`
	NamespaceDefaultResourceQuota NamespaceResourceQuota `json:"namespaceDefaultResourceQuota,omitempty"`
	ContainerDefaultResourceLimit ContainerResourceLimit `json:"containerDefaultResourceLimit,omitempty"`
	EnableProjectMonitoring       bool                   `json:"enableProjectMonitoring,omitempty"`
	Labels                        map[string]string      `json:"labels,omitempty"`
	PodSecurityPolicyTemplateName string                 `json:"podSecurityPolicyTemplateId,omitempty"`
	MonitoringInput               MonitoringInput        `json:"monitoringInput,omitempty"`
}

// ProjectStatus defines the observed state of Project
type ProjectStatus struct {
	Status string `json:"status"`
	ID     string `json:"id"`
	State  string `json:"state"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Cluster

// Project is the Schema for the projects API
type Project struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ProjectSpec   `json:"spec,omitempty"`
	Status ProjectStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Cluster

// ProjectList contains a list of Project
type ProjectList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Project `json:"items"`
}

type ProjectResourceQuota struct {
	Limit     *ResourceQuotaLimit `json:"limit,omitempty"`
	UsedLimit *ResourceQuotaLimit `json:"usedLimit,omitempty"`
}

type ResourceQuotaLimit struct {
	ConfigMaps             string `json:"configMaps,omitempty"`
	LimitsCPU              string `json:"limitsCpu,omitempty"`
	LimitsMemory           string `json:"limitsMemory,omitempty"`
	PersistentVolumeClaims string `json:"persistentVolumeClaims,omitempty"`
	Pods                   string `json:"pods,omitempty"`
	ReplicationControllers string `json:"replicationControllers,omitempty"`
	RequestsCPU            string `json:"requestsCpu,omitempty"`
	RequestsMemory         string `json:"requestsMemory,omitempty"`
	RequestsStorage        string `json:"requestsStorage,omitempty"`
	Secrets                string `json:"secrets,omitempty"`
	Services               string `json:"services,omitempty"`
	ServicesLoadBalancers  string `json:"servicesLoadBalancers,omitempty"`
	ServicesNodePorts      string `json:"servicesNodePorts,omitempty"`
}

type NamespaceResourceQuota struct {
	Limit *ResourceQuotaLimit `json:"limit,omitempty"`
}

type ContainerResourceLimit struct {
	LimitsCPU      string `json:"limitsCpu,omitempty"`
	LimitsMemory   string `json:"limitsMemory,omitempty"`
	RequestsCPU    string `json:"requestsCpu,omitempty"`
	RequestsMemory string `json:"requestsMemory,omitempty"`
}

type MonitoringInput struct {
	Version string            `json:"version,omitempty"`
	Answers map[string]string `json:"answers,omitempty"`
}

func init() {
	SchemeBuilder.Register(&Project{}, &ProjectList{})
}
