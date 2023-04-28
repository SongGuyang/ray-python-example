/*
Copyright 2023.

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
	"github.com/go-logr/logr"
	automlv1 "github.com/ray-automl/apis/automl/v1"
	"github.com/ray-automl/controllers/utils"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"time"
)

// TrainerReconciler reconciles a Trainer object
type TrainerReconciler struct {
	client.Client
	Config *rest.Config
	Log    logr.Logger
	Scheme *runtime.Scheme
}

var _ reconcile.Reconciler = &TrainerReconciler{}

//+kubebuilder:rbac:groups=automl.my.domain,resources=trainers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=automl.my.domain,resources=trainers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=automl.my.domain,resources=trainers/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// the Trainer object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *TrainerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var err error

	// Try to fetch the Trainer instance
	instance := &automlv1.Trainer{}
	if err = r.Get(context.TODO(), req.NamespacedName, instance); err == nil {
		r.Log.Info("reconcile for instance", "req", req.NamespacedName)
		return r.trainerReconcile(req, instance)
	}

	// No match found
	if errors.IsNotFound(err) {
		r.Log.Info("Read request instance not found error!", "name", req.NamespacedName)
	} else {
		r.Log.Error(err, "Read request instance error!")
	}

	return ctrl.Result{}, nil
}

func (r *TrainerReconciler) trainerReconcile(req ctrl.Request, instance *automlv1.Trainer) (ctrl.Result, error) {
	var err error

	if err = r.reconcileServices(instance); err != nil && !errors.IsAlreadyExists(err) {
		r.Log.Error(err, "failed to create service for instance", "instance", req.NamespacedName)
		return ctrl.Result{Requeue: true}, nil
	}

	if err := r.reconcileTrainerDeploy(instance); err != nil && !errors.IsAlreadyExists(err) {
		r.Log.Error(err, "failed to create trainer deploy for instance", "instance", req.NamespacedName)
		return ctrl.Result{Requeue: true}, nil
	}

	isReady := false
	if isReady, err = r.CheckDeploymentStatus(instance.Namespace, instance.Name); err != nil && !errors.IsNotFound(err) {
		r.Log.Error(err, "failed to check trainer deploy status for instance", "instance", req.NamespacedName)
		return ctrl.Result{Requeue: true}, nil
	}

	r.Log.Info("trainer deploy instance checkDeploymentStatus", "isReady", isReady)

	if isReady || !isReady {
		if err := r.reconcileTrainerWorkerDeploy(instance); err != nil && !errors.IsAlreadyExists(err) {
			r.Log.Error(err, "failed to create trainer worker deploy for instance", "instance", req.NamespacedName)
			return ctrl.Result{Requeue: true}, nil
		}
	} else {
		return ctrl.Result{RequeueAfter: 10 * time.Second}, nil
	}

	return ctrl.Result{}, nil
}

func (r *TrainerReconciler) reconcileServices(instance *automlv1.Trainer) error {
	service := utils.NewService(instance, r.Log)
	utils.SetTrainerOwnerReference(service, instance)
	if err := r.Create(context.TODO(), service); err != nil {
		return err
	}
	return nil
}

func (r *TrainerReconciler) reconcileTrainerDeploy(instance *automlv1.Trainer) error {
	deployment := utils.NewDeployment(instance, r.Log)
	utils.SetTrainerOwnerReference(deployment, instance)
	if err := r.Create(context.TODO(), deployment); err != nil {
		return err
	}
	return nil
}

func (r *TrainerReconciler) reconcileTrainerWorkerDeploy(instance *automlv1.Trainer) error {
	deployments := utils.NewDeploymentInstanceWorker(instance, r.Log)
	for _, deployment := range deployments {
		utils.SetTrainerOwnerReference(deployment, instance)
		r.Log.Info("reconcileTrainerWorkerDeploy", "deployment", deployment)
		if err := r.Create(context.TODO(), deployment); err != nil {
			return err
		}
	}
	return nil
}

// CheckDeploymentStatus checks if a deployment is running and ready
func (r *TrainerReconciler) CheckDeploymentStatus(namespace string, deploymentName string) (bool, error) {
	// Create Kubernetes clientset
	clientset, err := kubernetes.NewForConfig(r.Config)
	if err != nil {
		return false, err
	}

	// Wait for the deployment to be ready
	err = wait.PollImmediate(time.Second, time.Minute*5, func() (bool, error) {

		deployment, err := clientset.AppsV1().Deployments(namespace).Get(context.TODO(), deploymentName, metav1.GetOptions{})
		if err != nil {
			return false, err
		}
		if deployment.Status.ReadyReplicas == deployment.Status.Replicas {
			return true, nil
		}
		return false, nil
	})

	if err != nil {
		return false, err
	}
	return false, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *TrainerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&automlv1.Trainer{}).
		Complete(r)
}