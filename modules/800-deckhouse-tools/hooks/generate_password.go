/*
Copyright 2024 Flant JSC

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
	"github.com/deckhouse/deckhouse/go_lib/hooks/generate_password"
)

const (
	moduleValuesKey = "deckhouseTools"
	authSecretNS    = "d8-system"
	authSecretName  = "tools-basic-auth"
)

var (
	generatePasswordSettings = generate_password.HookSettings{
		ModuleName: moduleValuesKey,
		Namespace:  authSecretNS,
		SecretName: authSecretName,
	}
)

var _ = generate_password.RegisterHook(generatePasswordSettings)
