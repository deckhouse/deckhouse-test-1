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

package releaseupdater

import (
	"encoding/json"

	"github.com/deckhouse/deckhouse/deckhouse-controller/pkg/apis/deckhouse.io/v1alpha1"
)

// TODO: replace patch marshalling with controller-runtime patch
type StatusPatch v1alpha1.DeckhouseReleaseStatus

func (sp StatusPatch) MarshalJSON() ([]byte, error) {
	m := map[string]interface{}{
		"status": v1alpha1.DeckhouseReleaseStatus(sp),
	}

	return json.Marshal(m)
}

type ByVersion[R v1alpha1.Release] []R

func (a ByVersion[R]) Len() int {
	return len(a)
}

func (a ByVersion[R]) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ByVersion[R]) Less(i, j int) bool {
	return a[i].GetVersion().LessThan(a[j].GetVersion())
}
