/*
Copyright 2024 Flant JSC
Licensed under the Deckhouse Platform Enterprise Edition (EE) license. See https://github.com/deckhouse/deckhouse/blob/main/ee/LICENSE
*/

package hooks

import (
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"time"

	"github.com/flant/addon-operator/pkg/module_manager/go_hook"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/deckhouse/deckhouse/go_lib/certificate"
	"github.com/deckhouse/deckhouse/pkg/log"
	. "github.com/deckhouse/deckhouse/testing/hooks"
)

const clusterDomain = "cluster.local"

func genWebhookCa(logger *log.Logger) (*certificate.Authority, error) {
	ca, err := certificate.GenerateCA(logger, "capcd-controller-manager-webhook", certificate.WithKeyAlgo("ecdsa"),
		certificate.WithKeySize(256),
		certificate.WithCAExpiry("87600h"))
	if err != nil {
		return nil, fmt.Errorf("cannot generate CA: %v", err)
	}

	return &ca, nil
}

func genWebhookTLS(input *go_hook.HookInput, ca *certificate.Authority, cn string, sanPrefix string) (*certificate.Certificate, error) {
	tls, err := certificate.GenerateSelfSignedCert(input.Logger,
		cn,
		*ca,
		certificate.WithKeyAlgo("ecdsa"),
		certificate.WithKeySize(256),
		certificate.WithSigningDefaultExpiry((24*time.Hour)*365*10),
		certificate.WithSigningDefaultUsage([]string{"signing",
			"key encipherment",
			"requestheader-client",
		}),
		certificate.WithSANs(
			sanPrefix+".d8-cloud-provider-vcd",
			sanPrefix+".d8-cloud-provider-vcd.svc",
			sanPrefix+".d8-cloud-provider-vcd."+clusterDomain,
			sanPrefix+".d8-cloud-provider-vcd.svc."+clusterDomain,
		),
	)
	if err != nil {
		return nil, fmt.Errorf("cannot generate TLS: %v", err)
	}

	return &tls, err
}

var _ = Describe("Node Manager hooks :: generate_webhook_certs ::", func() {
	f := HookExecutionConfigInit(`{"global": {"discovery": {"clusterDomain": "`+clusterDomain+`"}},"cloudProviderVcd":{"internal":{"capcdControllerManagerWebhookCert": {}}}}`, "")

	Context("Without secret", func() {
		BeforeEach(func() {
			f.BindingContexts.Set(f.GenerateBeforeHelmContext())
			f.RunHook()
		})

		It("Should add ca and certificate to values", func() {
			Expect(f).To(ExecuteSuccessfully())
			Expect(f.BindingContexts.Array()).ShouldNot(BeEmpty())

			Expect(f.ValuesGet("cloudProviderVcd.internal.capcdControllerManagerWebhookCert.crt").Exists()).To(BeTrue())
			Expect(f.ValuesGet("cloudProviderVcd.internal.capcdControllerManagerWebhookCert.key").Exists()).To(BeTrue())
			Expect(f.ValuesGet("cloudProviderVcd.internal.capcdControllerManagerWebhookCert.crt").Exists()).To(BeTrue())

			blockCA, _ := pem.Decode([]byte(f.ValuesGet("cloudProviderVcd.internal.capcdControllerManagerWebhookCert.ca").String()))
			certCA, err := x509.ParseCertificate(blockCA.Bytes)
			Expect(err).To(BeNil())
			Expect(certCA.IsCA).To(BeTrue())
			Expect(certCA.Subject.CommonName).To(Equal("capcd-controller-manager-webhook"))

			block, _ := pem.Decode([]byte(f.ValuesGet("cloudProviderVcd.internal.capcdControllerManagerWebhookCert.crt").String()))
			cert, err := x509.ParseCertificate(block.Bytes)
			Expect(err).To(BeNil())
			Expect(cert.IsCA).To(BeFalse())
			Expect(cert.Subject.CommonName).To(Equal("capcd-controller-manager-webhook"))
		})
	})
	Context("With secrets", func() {
		caAuthority, _ := genWebhookCa(nil)
		tlsAuthority, _ := genWebhookTLS(&go_hook.HookInput{Logger: log.NewNop()}, caAuthority, "capcd-manager-webhook", "capcd-controller-manager-webhook-service")

		BeforeEach(func() {
			f.BindingContexts.Set(f.KubeStateSet(fmt.Sprintf(`
---
apiVersion: v1
kind: Secret
type: kubernetes.io/tls
metadata:
  name: capcd-controller-manager-webhook-tls
  namespace: d8-cloud-provider-vcd
data:
  ca.crt: %[1]s
  tls.crt: %[2]s
  tls.key: %[3]s
`, base64.StdEncoding.EncodeToString([]byte(caAuthority.Cert)),
				base64.StdEncoding.EncodeToString([]byte(tlsAuthority.Cert)),
				base64.StdEncoding.EncodeToString([]byte(tlsAuthority.Key)))),
			)
			f.RunHook()
		})
		It("Should add existing ca certificate to values", func() {
			Expect(f).To(ExecuteSuccessfully())
			Expect(f.BindingContexts.Array()).ShouldNot(BeEmpty())

			ca := f.ValuesGet("cloudProviderVcd.internal.capcdControllerManagerWebhookCert.ca").String()
			Expect(ca).To(Equal(caAuthority.Cert))
			Expect(f.ValuesGet("cloudProviderVcd.internal.capcdControllerManagerWebhookCert.key").String()).To(Equal(tlsAuthority.Key))
			Expect(f.ValuesGet("cloudProviderVcd.internal.capcdControllerManagerWebhookCert.crt").String()).To(Equal(tlsAuthority.Cert))
		})
	})
})
