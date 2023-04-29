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

package template_tests

import (
	"os"
	"strings"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/deckhouse/deckhouse/testing/helm"
)

func Test(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "")
}

const globalValues = `
modulesImages:
  registry:
    base: registry.deckhouse.io/deckhouse/ee
discovery:
  d8SpecificNodeCountByRole:
    system: 1
`

var _ = Describe("Module :: multitenancy-manager :: helm template ::", func() {
	f := SetupHelmConfig(``)

	renderMap := make(map[string]string)

	BeforeEach(func() {
		f.ValuesSetFromYaml("global", globalValues)
		f.ValuesSet("global.modulesImages", GetModulesImages())
	})

	AfterSuite(func() {
		builder := strings.Builder{}
		for k, v := range renderMap {
			builder.WriteString("\n# ")
			builder.WriteString(k)
			builder.WriteString("\n")
			builder.WriteString(v)

		}
		err := os.WriteFile("/tmp/rendered_templates.yaml", []byte(builder.String()), 0644)
		Expect(err).ShouldNot(HaveOccurred())
	})

	Context("", func() {
		BeforeEach(func() {
			f.ValuesSet("multitenancyManager.internal.webhookCertificate.ca", "test")
			f.ValuesSet("multitenancyManager.internal.webhookCertificate.crt", "test")
			f.ValuesSet("multitenancyManager.internal.webhookCertificate.key", "test")

			f.HelmRender(WithRenderOutput(renderMap))
		})

		It("Should create a ClusterRoleBinding for each additionalRole", func() {
			Expect(f.RenderError).ShouldNot(HaveOccurred())
		})

	})
})
