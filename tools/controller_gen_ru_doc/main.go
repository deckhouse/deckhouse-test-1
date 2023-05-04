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
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"sigs.k8s.io/yaml"
)

const (
	docRuPrefix = "doc-ru-"
)

func main() {
	for _, crdPath := range os.Args[1:] {
		files, err := ioutil.ReadDir(crdPath)
		if err != nil {
			log.Fatal(err)
		}
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			if strings.HasPrefix(file.Name(), docRuPrefix) {
				continue
			}
			filePath := filepath.Join(crdPath, file.Name())
			if err := processCrdFile(filePath); err != nil {
				log.Fatalln(err)
			}
		}
	}
}

func processCrdFile(fileName string) error {
	contents, err := readFileContent(fileName)
	if err != nil {
		return err
	}

	engCrd, err := readRuCrd(contents)
	if err != nil {
		return err
	}

	ruContents, exists, err := readRuFileIfExists(fileName)
	if err != nil {
		return err
	}
	if !exists {
		if err := writeRuCrd(fileName, engCrd); err != nil {
			return err
		}
		return nil
	}

	existingRuCrd, err := readRuCrd(ruContents)
	if err != nil {
		return err
	}

	newCrd, err := mapExistingCrdToGenerated(existingRuCrd, engCrd)
	if err != nil {
		return err
	}

	if err := writeRuCrd(fileName, newCrd); err != nil {
		return err
	}
	return nil
}

func readFileContent(filename string) ([]byte, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	renderContents, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return renderContents, nil
}

func readRuFileIfExists(engFileName string) ([]byte, bool, error) {
	contents, err := readFileContent(ruFilePath(engFileName))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, false, nil
		}
		return nil, true, err
	}
	return contents, true, nil
}

func writeRuCrd(engFileName string, crd *ruCustomResourceDefinition) error {
	marsh, err := yaml.Marshal(crd)
	if err != nil {
		return err
	}
	return os.WriteFile(ruFilePath(engFileName), marsh, 0644)
}

func ruFilePath(engFileName string) string {
	ruFileName := docRuPrefix + filepath.Base(engFileName)
	return filepath.Join(filepath.Dir(engFileName), ruFileName)
}

func readRuCrd(contents []byte) (*ruCustomResourceDefinition, error) {
	var ruCrd *ruCustomResourceDefinition
	if err := yaml.Unmarshal(contents, &ruCrd); err != nil {
		return nil, err
	}
	return ruCrd, nil
}

func mapExistingCrdToGenerated(existingCrd, newCrd *ruCustomResourceDefinition) (*ruCustomResourceDefinition, error) {
	newVersions := newCrd.Spec.Versions
	existingVersions := existingCrd.Spec.Versions
	crd := &ruCustomResourceDefinition{}
	for i, newVersion := range newVersions {
		if len(existingVersions) < i+1 {
			crd.Spec.Versions = append(crd.Spec.Versions, newVersion)
			continue
		}
		existingVersion := existingVersions[i]
		if newVersion.Name != existingVersion.Name {
			return nil, fmt.Errorf("name of versions for eng and ru doc are not in same order: %s and %s", newVersion.Name, existingVersion.Name)
		}
		newSchema := newVersion.Schema.OpenAPIV3Schema
		oldSchema := existingVersion.Schema.OpenAPIV3Schema

		version := ruCustomResourceDefinitionVersion{
			Name: newVersion.Name,
			Schema: &ruCustomResourceValidation{
				OpenAPIV3Schema: compareSchemasForExistingAndNew(oldSchema, newSchema),
			},
		}
		crd.Spec.Versions = append(crd.Spec.Versions, version)
	}
	return crd, nil
}

func compareSchemasForExistingAndNew(oldSchema, newSchema *JSONSchemaProps) *JSONSchemaProps {
	jsonSchema := &JSONSchemaProps{}
	if newSchema == nil {
		return jsonSchema
	}
	jsonSchema.Description = newSchema.Description
	if oldSchema != nil && oldSchema.Description != "" {
		jsonSchema.Description = oldSchema.Description
	}

	switch {
	case len(newSchema.Properties) > 0:
		m := make(map[string]JSONSchemaProps)
		for key, value := range newSchema.Properties {
			if newSchema == nil || len(newSchema.Properties) < 1 {
				m[key] = value
				continue
			}
			oldValue, f := oldSchema.Properties[key]
			if !f {
				m[key] = value
				continue
			}
			m[key] = *compareSchemasForExistingAndNew(&oldValue, &value)
		}
		jsonSchema.Properties = m
	case newSchema.Items != nil && newSchema.Items.Schema != nil:
		jsonSchema.Items = &JSONSchemaPropsOrArray{}
		if oldSchema.Items == nil || oldSchema.Items.Schema == nil {
			jsonSchema.Items.Schema = newSchema.Items.Schema
			break
		}
		jsonSchema.Items.Schema = compareSchemasForExistingAndNew(oldSchema.Items.Schema, newSchema.Items.Schema)
	case newSchema.Items != nil && len(newSchema.Items.JSONSchemas) > 0:
		jsonSchema.Items = &JSONSchemaPropsOrArray{}
		if oldSchema.Items == nil || len(oldSchema.Items.JSONSchemas) < 1 {
			jsonSchema.Items.JSONSchemas = newSchema.Items.JSONSchemas
			break
		}
		items := make([]JSONSchemaProps, 0, len(newSchema.Items.JSONSchemas))
		oldItems := oldSchema.Items.JSONSchemas
		for i, nsc := range newSchema.Items.JSONSchemas {
			if len(oldItems) < i+1 {
				items = append(items, nsc)
				continue
			}
			oldsc := oldItems[i]
			items = append(items, *compareSchemasForExistingAndNew(&oldsc, &nsc))
		}
		jsonSchema.Items.JSONSchemas = items
	}
	return jsonSchema
}
