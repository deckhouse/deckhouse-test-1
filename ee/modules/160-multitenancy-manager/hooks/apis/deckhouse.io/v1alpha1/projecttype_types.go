/*
Copyright 2023 Flant JSC
Licensed under the Deckhouse Platform Enterprise Edition (EE) license. See https://github.com/deckhouse/deckhouse/blob/main/ee/LICENSE
*/

// +kubebuilder:object:generate=true
// +kubebuilder:validation:Required
// +groupName=deckhouse.io
// +versionName=v1alpha1

package v1alpha1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/go-openapi/spec"
	apiext "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const OneOfForSubjects = `
- required: [kind, name, namespace, role]
  properties:
    kind:
      enum: [ServiceAccount]
    name: {}
    namespace: {}
    role: {}
- required: [kind, name, role]
  properties:
    kind:
      enum: [User,Group]
    name: {}
    role: {}
`

type ProjectTypeSpec struct {
	// Subjects to apply access to project
	Subjects []AuthorizationRule `json:"subjects,omitempty"`

	// +optional
	// NamespaceMetadata is metadata to apply to Project namespaces
	NamespaceMetadata metav1.ObjectMeta `json:"namespaceMetadata,omitempty"`

	// +kubebuilder:validation:PreserveUnknownFields
	// +kubebuilder:validation:Schemaless
	// OpenApi spec for tempaltes
	OpenAPI *apiext.JSON `json:"openAPI,omitempty"`

	// Resource templates
	ResourcesTemplate ResourceTemplate `json:"resourcesTemplate,omitempty"`
}

func (p *ProjectTypeSpec) LoadOpenAPISchema() (*spec.Schema, error) {
	d, err := json.Marshal(p.OpenAPI)
	if err != nil {
		return nil, fmt.Errorf("json marshal spec.openAPI: %w", err)
	}

	schema := new(spec.Schema)
	if err := json.Unmarshal(d, schema); err != nil {
		return nil, fmt.Errorf("unmarshal spec.openAPI to spec.Schema: %w", err)
	}

	err = spec.ExpandSchema(schema, schema, nil)
	if err != nil {
		return nil, fmt.Errorf("expand the schema in spec.openAPI: %w", err)
	}

	if schema.Type == nil {
		return nil, fmt.Errorf("incorrect type of schema in spec.openAPI")
	}

	return schema, nil
}

// +kubebuilder:validation:OneOf=../ee/modules/160-multitenancy-manager/hooks/apis/deckhouse.io/v1alpha1/projecttype_types.go=OneOfForSubjects
type AuthorizationRule struct {
	// +kubebuilder:validation:Enum=ServiceAccount;User;Group
	Kind string `json:"kind,omitempty"`
	// +kubebuilder:validation:MinLength=1
	Name string `json:"name,omitempty"`
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:Pattern="[a-z0-9]([-a-z0-9]*[a-z0-9])?"
	Namespace string `json:"namespace,omitempty"`
	// +kubebuilder:validation:Enum=User;PrivilegedUser;Editor;Admin
	Role string `json:"role,omitempty"`
}

type ResourceTemplate string

func (r ResourceTemplate) Render(args interface{}) (string, error) {
	var res bytes.Buffer

	tpl, err := template.New("resource-template").
		Funcs(sprig.TxtFuncMap()).
		Parse(string(r))
	if err != nil {
		return "", err
	}

	err = tpl.Execute(&res, args)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(res.String()), nil
}

type ProjectTypeStatus struct {
}

// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced
// +kubebuilder:metadata:labels=module=deckhouse;heritage=deckhouse
// +kubebuilder:resource:shortName=pt
// +kubebuilder:object:root=true
type ProjectType struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ProjectTypeSpec   `json:"spec,omitempty"`
	Status ProjectTypeStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
type ProjectTypeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ProjectType `json:"items"`
}
