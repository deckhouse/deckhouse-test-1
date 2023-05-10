/*
Copyright 2023 Flant JSC
Licensed under the Deckhouse Platform Enterprise Edition (EE) license. See https://github.com/deckhouse/deckhouse/blob/main/ee/LICENSE
*/

package hooks

import (
	"fmt"

	"github.com/deckhouse/deckhouse/ee/modules/160-multitenancy-manager/hooks/internal"
	"github.com/deckhouse/deckhouse/go_lib/hooks/tls_certificate"
)

const (
	certificateSecretName = "multitenancy-manager-webhook"
)

var _ = tls_certificate.RegisterInternalTLSHook(tls_certificate.GenSelfSignedTLSHookConf{
	SANs: tls_certificate.DefaultSANs([]string{
		fmt.Sprintf("%s.%s.svc", certificateSecretName, internal.Namespace),
		fmt.Sprintf("%s.%s", certificateSecretName, internal.Namespace),
		certificateSecretName,
	}),
	CN: certificateSecretName,

	Namespace:            internal.Namespace,
	TLSSecretName:        certificateSecretName,
	FullValuesPathPrefix: "multitenancyManager.internal.webhookCertificate",
})
