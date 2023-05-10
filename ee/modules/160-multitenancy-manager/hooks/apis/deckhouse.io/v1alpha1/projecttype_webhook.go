/*
Copyright 2023 Flant JSC
Licensed under the Deckhouse Platform Enterprise Edition (EE) license. See https://github.com/deckhouse/deckhouse/blob/main/ee/LICENSE
*/

package v1alpha1

import (
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

var _ webhook.Validator = &ProjectType{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (w *ProjectType) ValidateCreate() error {
	return w.validate()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (w *ProjectType) ValidateUpdate(old runtime.Object) error {
	return w.validate()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (w *ProjectType) ValidateDelete() error {
	return nil
}

func (w *ProjectType) validate() error {
	var allErrs field.ErrorList
	if err := w.validateSpecOpenAPI(); err != nil {
		allErrs = append(allErrs, err)
	}
	if len(allErrs) == 0 {
		return nil
	}
	return apierrors.NewInvalid(
		w.GroupVersionKind().GroupKind(),
		w.Name, allErrs)
}

func (w ProjectType) validateSpecOpenAPI() *field.Error {
	if _, err := w.Spec.LoadOpenAPISchema(); err != nil {
		return field.Invalid(
			field.NewPath("spec").Child("openAPI"),
			w.Spec.OpenAPI,
			"invalid openAPI schema: "+err.Error())
	}
	return nil
}
func (r *ProjectType) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}
