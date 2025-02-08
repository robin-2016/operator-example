/*
Copyright 2025.

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
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	webv1 "github.com/robin-2016/operator-example/api/v1"
)

type CertManager struct {
}

// WebsiteReconciler reconciles a Website object
type WebsiteReconciler struct {
	client.Client
	Scheme      *runtime.Scheme
	certManager CertManager
}

// +kubebuilder:rbac:groups=web.example.com,resources=websites,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=web.example.com,resources=websites/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=web.example.com,resources=websites/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Website object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.20.0/pkg/reconcile
func (r *WebsiteReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	klog.Infoln(req.Name)
	website := &webv1.Website{}
	if err := r.Get(ctx, req.NamespacedName, website); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// 自动组装Deployment
	if err := r.buildDeployment(website); err != nil {
		return ctrl.Result{}, err
	}

	// 智能HTTPS配置
	if website.Spec.Https {
		if err := r.certManager.IssueCert(website); err != nil {
			// 自动重试机制
			return ctrl.Result{RequeueAfter: 5 * time.Minute}, nil
		}
	}

	return ctrl.Result{RequeueAfter: 10 * time.Minute}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *WebsiteReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&webv1.Website{}).
		Named("website").
		Complete(r)
}

func (r *WebsiteReconciler) buildDeployment(website *webv1.Website) error {
	klog.Infoln(website.Name)
	return nil
}

func (cm *CertManager) IssueCert(website *webv1.Website) error {
	klog.Infoln("test")
	return nil
}
