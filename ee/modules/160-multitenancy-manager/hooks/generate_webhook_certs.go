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
