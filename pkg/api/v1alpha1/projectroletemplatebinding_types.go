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

// ProjectRoleTemplateBindingSpec defines the desired state of ProjectRoleTemplateBinding
type ProjectRoleTemplateBindingSpec struct {
	UserName           string `json:"userName,omitempty"`
	UserPrincipalName  string `json:"userPrincipalName,omitempty"`
	GroupName          string `json:"groupName,omitempty"`
	GroupPrincipalName string `json:"groupPrincipalName,omitempty"`
	ProjectName        string `json:"projectName"`
	RoleTemplateName   string `json:"roleTemplateName,omitempty"`
}

// ProjectRoleTemplateBindingStatus defines the observed state of ProjectRoleTemplateBinding
type ProjectRoleTemplateBindingStatus struct {
	Status string `json:"status"`
	ID     string `json:"id"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Cluster

// ProjectRoleTemplateBinding is the Schema for the projectroletemplatebindings API
type ProjectRoleTemplateBinding struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ProjectRoleTemplateBindingSpec   `json:"spec,omitempty"`
	Status ProjectRoleTemplateBindingStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Cluster

// ProjectRoleTemplateBindingList contains a list of ProjectRoleTemplateBinding
type ProjectRoleTemplateBindingList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ProjectRoleTemplateBinding `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ProjectRoleTemplateBinding{}, &ProjectRoleTemplateBindingList{})
}
