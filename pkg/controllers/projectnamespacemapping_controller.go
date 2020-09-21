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

package controllers

import (
	"context"
	"reflect"

	json "github.com/json-iterator/go"

	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"github.com/ibrokethecloud/rancher-project-operator/pkg/management"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	managementv1alpha1 "github.com/ibrokethecloud/rancher-project-operator/pkg/api/v1alpha1"
)

// ProjectNamespaceMappingReconciler reconciles a ProjectNamespaceMapping object
type ProjectNamespaceMappingReconciler struct {
	client.Client
	Log       logr.Logger
	Scheme    *runtime.Scheme
	Namespace string
	Secret    string
}

// +kubebuilder:rbac:groups=management.cattle.io,resources=projectnamespacemappings,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=management.cattle.io,resources=projectnamespacemappings/status,verbs=get;update;patch

func (r *ProjectNamespaceMappingReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("projectnamespacemapping", req.NamespacedName)

	pnsMapping := &managementv1alpha1.ProjectNamespaceMapping{}
	if err := r.Get(ctx, req.NamespacedName, pnsMapping); err != nil {
		log.Info("Unable to fetch ProjectNamespaceMapping")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	secret, err := r.GetSecret(ctx)
	if err != nil {
		log.Info("Unable to fetch config secret: ", r.Secret)
		return ctrl.Result{}, err
	}

	wrapper, err := management.NewClient(secret)
	if err != nil {
		log.Info("Error creating a ApiWrapper")
		return ctrl.Result{}, err
	}

	if pnsMapping.ObjectMeta.DeletionTimestamp.IsZero() {
		// create / update the project
		status := pnsMapping.Status.DeepCopy()
		oldSpec := pnsMapping.Annotations["appliedSpec"]
		newSpecByte, err := json.Marshal(pnsMapping.Spec)
		if err != nil {
			return ctrl.Result{}, err
		}

		if !compareSpec(oldSpec, newSpecByte) {
			status, err = wrapper.ManageProjectNameSpaceMapping(pnsMapping)
			if err != nil {
				return ctrl.Result{}, err
			}
			controllerutil.AddFinalizer(pnsMapping, projectFinalizer)
			pnsMapping.Status = *status
			curAnnotation := pnsMapping.Annotations
			curAnnotation["appliedSpec"] = string(newSpecByte)
			pnsMapping.Annotations = curAnnotation
			if err = r.Update(ctx, pnsMapping); err != nil {
				return ctrl.Result{}, err
			}
		}

	} else {
		// cleanup project namespace mapping //
		if controllerutil.ContainsFinalizer(pnsMapping, projectFinalizer) {
			if err = wrapper.RemoveProjectNamespaceMapping(pnsMapping); err != nil {
				return ctrl.Result{}, err
			}
		}
		controllerutil.RemoveFinalizer(pnsMapping, projectFinalizer)
		if err = r.Update(ctx, pnsMapping); err != nil {
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

func (r *ProjectNamespaceMappingReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&managementv1alpha1.ProjectNamespaceMapping{}).
		Complete(r)
}

func (r *ProjectNamespaceMappingReconciler) GetSecret(ctx context.Context) (secret corev1.Secret, err error) {
	secret = corev1.Secret{}
	err = r.Get(ctx, types.NamespacedName{Name: r.Secret, Namespace: r.Namespace}, &secret)
	return secret, err
}

func compareSpec(appliedSpec string, currentSpecByte []byte) bool {
	if len(appliedSpec) == 0 {
		return false
	}
	oldSpec := make(map[string]interface{})
	err := json.Unmarshal([]byte(appliedSpec), &oldSpec)
	if err != nil {
		return false
	}
	newSpec := make(map[string]interface{})
	err = json.Unmarshal(currentSpecByte, &newSpec)
	if err != nil {
		return false
	}

	return reflect.DeepEqual(oldSpec, newSpec)
}
