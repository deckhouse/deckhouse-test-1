/*
Copyright 2021 Flant JSC

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
	"encoding/json"
	"fmt"

	"github.com/flant/addon-operator/pkg/module_manager/go_hook"
	"github.com/flant/addon-operator/sdk"

	"github.com/deckhouse/deckhouse/go_lib/pwgen"
)

var _ = sdk.RegisterFunc(&go_hook.HookConfig{
	OnBeforeHelm: &go_hook.OrderedConfig{Order: 10},
}, generatePassword)

const defaultLocationTemplate = `
[ {
  "users": {"admin": "%s"},
  "location": "/"
} ]
`

func generatePassword(input *go_hook.HookInput) error {
	_, ok := input.Values.GetOk("basicAuth.locations")
	if ok {
		return nil
	}

	rawLocations := make([]map[string]interface{}, 0)

	locations := fmt.Sprintf(defaultLocationTemplate, pwgen.AlphaNum(20))
	err := json.Unmarshal([]byte(locations), &rawLocations)
	if err != nil {
		return err
	}

	input.ConfigValues.Set("basicAuth.locations", rawLocations)
	return nil
}
