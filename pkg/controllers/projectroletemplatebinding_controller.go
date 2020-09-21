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

	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"github.com/go-logr/logr"
	json "github.com/json-iterator/go"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	managementv1alpha1 "github.com/ibrokethecloud/rancher-project-operator/pkg/api/v1alpha1"
	"github.com/ibrokethecloud/rancher-project-operator/pkg/management"
)

// ProjectRoleTemplateBindingReconciler reconciles a ProjectRoleTemplateBinding object
type ProjectRoleTemplateBindingReconciler struct {
	client.Client
	Log       logr.Logger
	Scheme    *runtime.Scheme
	Namespace string
	Secret    string
}

// +kubebuilder:rbac:groups=management.cattle.io,resources=projectroletemplatebindings,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=management.cattle.io,resources=projectroletemplatebindings/status,verbs=get;update;patch

func (r *ProjectRoleTemplateBindingReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("projectroletemplatebinding", req.NamespacedName)

	prtb := &managementv1alpha1.ProjectRoleTemplateBinding{}
	if err := r.Get(ctx, req.NamespacedName, prtb); err != nil {
		log.Info("Unable to fetch Project")
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

	if prtb.ObjectMeta.DeletionTimestamp.IsZero() {
		// create / update the prtb
		status := prtb.Status.DeepCopy()
		oldSpec := prtb.Annotations["appliedSpec"]
		newSpecByte, err := json.Marshal(prtb.Spec)
		if err != nil {
			return ctrl.Result{}, err
		}
		if !compareSpec(oldSpec, newSpecByte) {
			log.Info(oldSpec)
			log.Info(string(newSpecByte))
			status, err = wrapper.ManageRoleTemplateBinding(prtb)
			if err != nil {
				return ctrl.Result{}, err
			}

			controllerutil.AddFinalizer(prtb, projectFinalizer)
			prtb.Status = *status
			curAnnotation := prtb.Annotations
			curAnnotation["appliedSpec"] = string(newSpecByte)
			if err = r.Update(ctx, prtb); err != nil {
				return ctrl.Result{}, err
			}
		}

	} else {
		// delete the prtb
		if controllerutil.ContainsFinalizer(prtb, projectFinalizer) {
			if err = wrapper.DeletePRTB(prtb.Status.ID); err != nil {
				return ctrl.Result{}, err
			}
		}
		controllerutil.RemoveFinalizer(prtb, projectFinalizer)
		if err = r.Update(ctx, prtb); err != nil {
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

func (r *ProjectRoleTemplateBindingReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&managementv1alpha1.ProjectRoleTemplateBinding{}).
		Complete(r)
}

func (r *ProjectRoleTemplateBindingReconciler) GetSecret(ctx context.Context) (secret corev1.Secret, err error) {
	secret = corev1.Secret{}
	err = r.Get(ctx, types.NamespacedName{Name: r.Secret, Namespace: r.Namespace}, &secret)
	return secret, err
}
