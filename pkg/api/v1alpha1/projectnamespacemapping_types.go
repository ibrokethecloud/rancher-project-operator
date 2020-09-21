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

// ProjectNamespaceMappingSpec defines the desired state of ProjectNamespaceMapping
type ProjectNamespaceMappingSpec struct {
	ProjectName string   `json:"projectName"`
	Namespaces  []string `json:"namespaces"`
}

// ProjectNamespaceMappingStatus defines the observed state of ProjectNamespaceMapping
type ProjectNamespaceMappingStatus struct {
	Status string `json:"status"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Cluster

// ProjectNamespaceMapping is the Schema for the projectnamespacemappings API
type ProjectNamespaceMapping struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ProjectNamespaceMappingSpec   `json:"spec,omitempty"`
	Status ProjectNamespaceMappingStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Cluster

// ProjectNamespaceMappingList contains a list of ProjectNamespaceMapping
type ProjectNamespaceMappingList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ProjectNamespaceMapping `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ProjectNamespaceMapping{}, &ProjectNamespaceMappingList{})
}
