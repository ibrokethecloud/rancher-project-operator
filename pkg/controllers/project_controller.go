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

	"github.com/ibrokethecloud/rancher-project-operator/pkg/management"

	"k8s.io/apimachinery/pkg/types"

	"github.com/go-logr/logr"
	managementv1alpha1 "github.com/ibrokethecloud/rancher-project-operator/pkg/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ProjectReconciler reconciles a Project object
type ProjectReconciler struct {
	client.Client
	Log       logr.Logger
	Scheme    *runtime.Scheme
	Namespace string
	Secret    string
}

const (
	projectFinalizer = "controller.management.cattle.io"
)

// +kubebuilder:rbac:groups=management.cattle.io,resources=projects,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=management.cattle.io,resources=projects/status,verbs=get;update;patch

func (r *ProjectReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("project", req.NamespacedName)

	project := &managementv1alpha1.Project{}
	if err := r.Get(ctx, req.NamespacedName, project); err != nil {
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

	if project.ObjectMeta.DeletionTimestamp.IsZero() {
		// create / update the project
		status := project.Status.DeepCopy()
		if status != nil && status.Status != management.ProjectSynced {
			status, err = wrapper.ManageProject(project)
		}

		if err != nil {
			return ctrl.Result{}, err
		}
		controllerutil.AddFinalizer(project, projectFinalizer)
		project.Status = *status
		if err = r.Update(ctx, project); err != nil {
			return ctrl.Result{}, err
		}
	} else {
		// Clean up the project
		if controllerutil.ContainsFinalizer(project, projectFinalizer) {
			if err = wrapper.RemoveProject(project); err != nil {
				return ctrl.Result{}, err
			}
		}
		controllerutil.RemoveFinalizer(project, projectFinalizer)
		if err = r.Update(ctx, project); err != nil {
			return ctrl.Result{}, err
		}
	}
	return ctrl.Result{}, nil
}

func (r *ProjectReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&managementv1alpha1.Project{}).
		Complete(r)
}

func (r *ProjectReconciler) GetSecret(ctx context.Context) (secret corev1.Secret, err error) {
	secret = corev1.Secret{}
	err = r.Get(ctx, types.NamespacedName{Name: r.Secret, Namespace: r.Namespace}, &secret)
	return secret, err
}
