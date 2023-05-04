/*
Copyright 2023 Flant JSC

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

package main

import (
	"encoding/json"

	apiext "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

type ruCustomResourceDefinition struct {
	Spec ruCustomResourceDefinitionSpec `json:"spec"`
}

type ruCustomResourceDefinitionSpec struct {
	Versions []ruCustomResourceDefinitionVersion
}

type ruCustomResourceDefinitionVersion struct {
	Name   string                      `json:"name"`
	Schema *ruCustomResourceValidation `json:"schema,omitempty"`
}

type ruCustomResourceValidation struct {
	OpenAPIV3Schema *JSONSchemaProps `json:"openAPIV3Schema,omitempty"`
}

type JSONSchemaProps struct {
	Description string                     `json:"description"`
	Default     *apiext.JSON               `json:"default,omitempty"`
	Items       *JSONSchemaPropsOrArray    `json:"items,omitempty"`
	Properties  map[string]JSONSchemaProps `json:"properties,omitempty"`
	Example     *apiext.JSON               `json:"example,omitempty"`
}

type JSONSchemaPropsOrArray struct {
	Schema      *JSONSchemaProps
	JSONSchemas []JSONSchemaProps
}

func (s JSONSchemaPropsOrArray) MarshalJSON() ([]byte, error) {
	if len(s.JSONSchemas) > 0 {
		return json.Marshal(s.JSONSchemas)
	}
	return json.Marshal(s.Schema)
}

func (s *JSONSchemaPropsOrArray) UnmarshalJSON(data []byte) error {
	var nw JSONSchemaPropsOrArray
	var first byte
	if len(data) > 1 {
		first = data[0]
	}
	if first == '{' {
		var sch JSONSchemaProps
		if err := json.Unmarshal(data, &sch); err != nil {
			return err
		}
		nw.Schema = &sch
	}
	if first == '[' {
		if err := json.Unmarshal(data, &nw.JSONSchemas); err != nil {
			return err
		}
	}
	*s = nw
	return nil
}
