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

package v1

import (
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var ytsauruslog = logf.Log.WithName("ytsaurus-resource")

func (r *Ytsaurus) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

//+kubebuilder:webhook:path=/mutate-cluster-ytsaurus-tech-v1-ytsaurus,mutating=true,failurePolicy=fail,sideEffects=None,groups=cluster.ytsaurus.tech,resources=ytsaurus,verbs=create;update,versions=v1,name=mytsaurus.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &Ytsaurus{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *Ytsaurus) Default() {
	ytsauruslog.Info("default", "name", r.Name)

	if r.Spec.Masters.InstanceGroup.EnableAntiAffinity == nil {
		r.Spec.Masters.InstanceGroup.EnableAntiAffinity = new(bool)
		*r.Spec.Masters.InstanceGroup.EnableAntiAffinity = true
	}
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:path=/validate-cluster-ytsaurus-tech-v1-ytsaurus,mutating=false,failurePolicy=fail,sideEffects=None,groups=cluster.ytsaurus.tech,resources=ytsaurus,verbs=create;update,versions=v1,name=vytsaurus.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &Ytsaurus{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *Ytsaurus) ValidateCreate() error {
	ytsauruslog.Info("validate create", "name", r.Name)

	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *Ytsaurus) ValidateUpdate(old runtime.Object) error {
	ytsauruslog.Info("validate update", "name", r.Name)

	// TODO(user): fill in your validation logic upon object update.
	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *Ytsaurus) ValidateDelete() error {
	ytsauruslog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil
}
