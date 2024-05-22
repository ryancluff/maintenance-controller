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

package controller

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"

	clusterv1 "rcluff.com/maintenance-controller/api/v1"
)

// MaintenanceModeReconciler reconciles a MaintenanceMode object
type MaintenanceModeReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

var (
	MaintenanceModeAnnotation = "maintenance-controller.rcluff.com/maintenance-mode"
)

//+kubebuilder:rbac:groups=cluster.rcluff.com,resources=maintenancemodes,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=cluster.rcluff.com,resources=maintenancemodes/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=cluster.rcluff.com,resources=maintenancemodes/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;update;patch
//+kubebuilder:rbac:groups=apps,resources=deployments/status,verbs=get

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the MaintenanceMode object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.3/pkg/reconcile
func (r *MaintenanceModeReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// 1 - Get the MaintenanceMode object
	// 3 - Get all Deployments
	// 4 - Determine the affected Deployments
	// 4a - have a pvc (of the specified storageclass)
	// 4b - Select deployments that don't have annotations
	// 5 - Set the affected Deployments to 0 replicas

	var maintenanceMode clusterv1.MaintenanceMode
	if err := r.Get(ctx, req.NamespacedName, &maintenanceMode); err != nil {
		log.Error(err, "unable to fetch MaintenanceMode")
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	var deployments appsv1.DeploymentList

	var persistentVolumeClaims corev1.PersistentVolumeClaimList
	if maintenanceMode.Spec.Scope == "namespace" {
		if err := r.List(ctx, &persistentVolumeClaims, client.InNamespace(req.Namespace)); err != nil {
			log.Error(err, "unable to list PersistentVolumeClaims in namespace")
			return ctrl.Result{}, err
		}
		if err := r.List(ctx, &deployments, client.InNamespace(req.Namespace)); err != nil {
			log.Error(err, "unable to list Deployments")
			return ctrl.Result{}, err
		}
	} else {
		if err := r.List(ctx, &persistentVolumeClaims); err != nil {
			log.Error(err, "unable to list PersistentVolumeClaims")
			return ctrl.Result{}, err
		}
		if err := r.List(ctx, &deployments); err != nil {
			log.Error(err, "unable to list Deployments")
			return ctrl.Result{}, err
		}
	}

	for _, deployment := range deployments.Items {

		if deployment.Namespace == maintenanceMode.Namespace {
			if deployment.Spec.Replicas != nil {
				if *deployment.Spec.Replicas > 0 {
					log.Info("Setting deployment to 0", "deployment", deployment.Name)
					deployment.Spec.Replicas = new(int32)
					*deployment.Spec.Replicas = 0
					if err := r.Update(ctx, &deployment); err != nil {
						log.Error(err, "unable to update Deployment", "deployment", deployment.Name)
						return ctrl.Result{}, err
					}
				}
			}
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *MaintenanceModeReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&clusterv1.MaintenanceMode{}).
		Complete(r)
}
