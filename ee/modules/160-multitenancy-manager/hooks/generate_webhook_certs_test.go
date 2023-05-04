/*
Copyright 2023 Flant JSC
Licensed under the Deckhouse Platform Enterprise Edition (EE) license. See https://github.com/deckhouse/deckhouse/blob/main/ee/LICENSE
*/

package hooks

/*

User-stories:
Webhook mechanism requires a pair of certificates. This hook generates them and stores in cluster as Secret resource.

*/

import (
	"crypto/x509"
	"encoding/pem"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/deckhouse/deckhouse/testing/hooks"
)

const (
	initValuesString       = `{"multitenancyManager":{}}`
	initConfigValuesString = `{}`
)

const (
	stateSecretCreated = `
apiVersion: v1
kind: Secret
metadata:
  name: multitenancy-manager-webhook
  namespace: d8-multitenancy-manager
data:
  ca.crt: YQo= # a
  tls.crt: Ygo= # b
  tls.key: Ywo= # c
`

	stateSecretChanged = `
apiVersion: v1
kind: Secret
metadata:
  name:  multitenancy-manager-webhook
  namespace: d8-multitenancy-manager
data:
  ca.crt: eAo= # x
  tls.crt: eQo= # y
  tls.key: ego= # z
`
)

var _ = Describe("Multitenancy Manager hooks :: gen webhook certs ::", func() {
	f := HookExecutionConfigInit(initValuesString, initConfigValuesString)

	Context("Empty cluster", func() {
		BeforeEach(func() {
			f.BindingContexts.Set(f.KubeStateSet(``))
			f.RunHook()
		})

		It("", func() {
			Expect(f).To(ExecuteSuccessfully())
		})

		Context("Secret Created", func() {
			BeforeEach(func() {
				f.BindingContexts.Set(f.KubeStateSet(stateSecretCreated))
				f.RunHook()
			})

			It("Cert data must be stored in values", func() {
				Expect(f).To(ExecuteSuccessfully())
				Expect(f.ValuesGet("multitenancyManager.internal.webhookCertificate.ca").String()).To(Equal("a\n"))
				Expect(f.ValuesGet("multitenancyManager.internal.webhookCertificate.crt").String()).To(Equal("b\n"))
				Expect(f.ValuesGet("multitenancyManager.internal.webhookCertificate.key").String()).To(Equal("c\n"))
			})

			Context("Secret Changed", func() {
				BeforeEach(func() {
					f.BindingContexts.Set(f.KubeStateSet(stateSecretChanged))
					f.RunHook()
				})

				It("New cert data must be stored in values", func() {
					Expect(f).To(ExecuteSuccessfully())
					Expect(f.ValuesGet("multitenancyManager.internal.webhookCertificate.ca").String()).To(Equal("x\n"))
					Expect(f.ValuesGet("multitenancyManager.internal.webhookCertificate.crt").String()).To(Equal("y\n"))
					Expect(f.ValuesGet("multitenancyManager.internal.webhookCertificate.key").String()).To(Equal("z\n"))
				})
			})
		})
	})

	Context("Cluster with secret", func() {
		BeforeEach(func() {
			f.BindingContexts.Set(f.KubeStateSet(stateSecretCreated))
			f.RunHook()
		})

		It("Cert data must be stored in values", func() {
			Expect(f).To(ExecuteSuccessfully())
			Expect(f.ValuesGet("userAuthz.internal.webhookCertificate.ca").String()).To(Equal("a\n"))
			Expect(f.ValuesGet("userAuthz.internal.webhookCertificate.crt").String()).To(Equal("b\n"))
			Expect(f.ValuesGet("userAuthz.internal.webhookCertificate.key").String()).To(Equal("c\n"))
		})
	})

	Context("Empty cluster with multitenancy, onBeforeHelm", func() {
		BeforeEach(func() {
			// TODO we need to unset cluster state between contexts.
			f.BindingContexts.Set(f.KubeStateSet(``))
			f.BindingContexts.Set(f.GenerateBeforeHelmContext())
			f.ValuesSet("userAuthz.enableMultiTenancy", true)
			f.RunHook()
		})

		It("New cert data must be generated and stored to values", func() {
			Expect(f).To(ExecuteSuccessfully())
			Expect(f.ValuesGet("userAuthz.internal.webhookCertificate.ca").Exists()).To(BeTrue())
			Expect(f.ValuesGet("userAuthz.internal.webhookCertificate.crt").Exists()).To(BeTrue())
			Expect(f.ValuesGet("userAuthz.internal.webhookCertificate.key").Exists()).To(BeTrue())

			certPool := x509.NewCertPool()
			ok := certPool.AppendCertsFromPEM([]byte(f.ValuesGet("userAuthz.internal.webhookCertificate.ca").String()))
			Expect(ok).To(BeTrue())

			block, _ := pem.Decode([]byte(f.ValuesGet("userAuthz.internal.webhookCertificate.crt").String()))
			Expect(block).ShouldNot(BeNil())

			cert, err := x509.ParseCertificate(block.Bytes)
			Expect(err).ShouldNot(HaveOccurred())

			opts := x509.VerifyOptions{
				DNSName: "127.0.0.1",
				Roots:   certPool,
			}

			_, err = cert.Verify(opts)
			Expect(err).ShouldNot(HaveOccurred())
		})
	})
})
