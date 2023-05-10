/*
Copyright 2023 Flant JSC
Licensed under the Deckhouse Platform Enterprise Edition (EE) license. See https://github.com/deckhouse/deckhouse/blob/main/ee/LICENSE
*/

package template_tests

import (
	"encoding/base64"
	"fmt"
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
discovery:
  d8SpecificNodeCountByRole:
    system: 1
`

var _ = Describe("Module :: multitenancy-manager :: helm template ::", func() {
	f := SetupHelmConfig(``)

	BeforeEach(func() {
		f.ValuesSetFromYaml("global", globalValues)
		f.ValuesSet("global.modulesImages", GetModulesImages())
	})

	Context("webhook secret", func() {
		BeforeEach(func() {
			f.ValuesSet("multitenancyManager.internal.webhookCertificate.ca", "test-ca")
			f.ValuesSet("multitenancyManager.internal.webhookCertificate.crt", "test-crt")
			f.ValuesSet("multitenancyManager.internal.webhookCertificate.key", "test-key")

			f.HelmRender()
		})

		It("Should create a Secret for webhook", func() {
			Expect(f.RenderError).ShouldNot(HaveOccurred())
			secret := f.KubernetesResource("Secret", "d8-multitenancy-manager", "multitenancy-manager-webhook")
			Expect(secret.Exists()).To(BeTrue())
			b64SecretData := fmt.Sprintf(
				`{"ca.crt":"%s","tls.crt":"%s","tls.key":"%s"}`,
				base64.StdEncoding.EncodeToString([]byte("test-ca")),
				base64.StdEncoding.EncodeToString([]byte("test-crt")),
				base64.StdEncoding.EncodeToString([]byte("test-key")),
			)
			Expect(secret.Field("data").String()).To(Equal(b64SecretData))
		})
	})
})
