/*
Copyright 2023 Flant JSC
Licensed under the Deckhouse Platform Enterprise Edition (EE) license. See https://github.com/deckhouse/deckhouse/blob/main/ee/LICENSE
*/

package hooks

import (
	"fmt"

	"github.com/flant/addon-operator/pkg/module_manager/go_hook"

	"github.com/deckhouse/deckhouse/ee/modules/160-multitenancy-manager/hooks/internal"
	"github.com/deckhouse/deckhouse/go_lib/hooks/tls_certificate"
)

const (
	certificateSecretName = "multitenancy-manager-webhook"
)

var ErrSkip = fmt.Errorf("skipping")

var _ = tls_certificate.RegisterInternalTLSHook(tls_certificate.GenSelfSignedTLSHookConf{
	BeforeHookCheck: func(input *go_hook.HookInput) bool {
		return len(input.Snapshots[tls_certificate.SnapshotKey]) > 0
	},
	SANs: tls_certificate.DefaultSANs([]string{"127.0.0.1"}),
	CN:   "127.0.0.1",

	Namespace:            internal.Namespace,
	TLSSecretName:        certificateSecretName,
	FullValuesPathPrefix: "multitenancyManager.internal.webhookCertificate",
})
