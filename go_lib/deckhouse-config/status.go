/*
Copyright 2022 Flant JSC

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

package deckhouse_config

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/flant/addon-operator/pkg/module_manager/models/modules"
	"github.com/flant/addon-operator/pkg/module_manager/scheduler/extenders"
	dynamic_extender "github.com/flant/addon-operator/pkg/module_manager/scheduler/extenders/dynamically_enabled"
	kube_config_extender "github.com/flant/addon-operator/pkg/module_manager/scheduler/extenders/kube_config"
	script_extender "github.com/flant/addon-operator/pkg/module_manager/scheduler/extenders/script_enabled"
	static_extender "github.com/flant/addon-operator/pkg/module_manager/scheduler/extenders/static"

	"github.com/deckhouse/deckhouse/deckhouse-controller/pkg/apis/deckhouse.io/v1alpha1"
	"github.com/deckhouse/deckhouse/go_lib/deckhouse-config/conversion"
	bootstrapped_extender "github.com/deckhouse/deckhouse/go_lib/dependency/extenders/bootstrapped"
	d7s_version_extender "github.com/deckhouse/deckhouse/go_lib/dependency/extenders/deckhouseversion"
	k8s_version_extender "github.com/deckhouse/deckhouse/go_lib/dependency/extenders/kubernetesversion"
	"github.com/deckhouse/deckhouse/go_lib/set"
)

type ModuleConfigStatus struct {
	Version string
	Message string
}

type ModuleStatus struct {
	Status     string
	Message    string
	HooksState string
}

type StatusReporter struct {
	moduleManager ModuleManager
	possibleNames set.Set
}

func NewModuleInfo(mm ModuleManager, possibleNames set.Set) *StatusReporter {
	return &StatusReporter{
		moduleManager: mm,
		possibleNames: possibleNames,
	}
}

func (s *StatusReporter) ForModule(module *v1alpha1.Module, cfg *v1alpha1.ModuleConfig, bundleName string) ModuleStatus {
	// Figure out additional statuses for known modules.
	statusMsgs := make([]string, 0)
	msgs := make([]string, 0)
	moduleName := module.GetName()

	mod := s.moduleManager.GetModule(module.GetName())
	// return error if module manager doesn't have such a module
	if mod == nil {
		return ModuleStatus{
			Status: "Error: failed to fetch module metadata",
		}
	}

	moduleEnabled := s.moduleManager.IsModuleEnabled(moduleName)
	modeDescriptor := "off"
	// Calculate state and status.
	if moduleEnabled {
		modeDescriptor = "on"
		lastHookErr := mod.GetLastHookError()
		if lastHookErr != nil {
			statusMsgs = append(statusMsgs, fmt.Sprintf("HookError: %v", lastHookErr))
		}
		if mod.GetModuleError() != nil {
			statusMsgs = append(statusMsgs, fmt.Sprintf("ModuleError: %v", mod.GetModuleError()))
		}

		if len(statusMsgs) == 0 { // no errors were added
			// Best effort alarm!
			//
			// Actually, this condition is not correct because the `CanRunHelm` status appears right before the first run.c
			// The right approach is to check the queue for the module run task.
			// However, there are too many addon-operator internals involved.
			// We should consider moving these statuses to the `Module` resource,
			// which is directly controlled by addon-operator.
			switch mod.GetPhase() {
			case modules.CanRunHelm:
				statusMsgs = append(statusMsgs, "Ready")
				// enrich the status message with a notification from the related module config
				if cfg != nil {
					cfgStatus := s.ForConfig(cfg)
					if len(cfgStatus.Message) > 0 {
						msgs = append(msgs, "check module configuration status")
					}
				}

			case modules.Startup:
				statusMsgs = append(statusMsgs, "Enqueued")

			case modules.OnStartupDone:
				statusMsgs = append(statusMsgs, "OnStartUp hooks are completed")

			case modules.WaitForSynchronization:
				statusMsgs = append(statusMsgs, "Synchronization tasks are running")

			case modules.HooksDisabled:
				statusMsgs = append(statusMsgs, "Pending: hooks are disabled")
			}
		}
	}

	updatedBy, updatedByErr := s.moduleManager.GetUpdatedByExtender(moduleName)
	var actor, additionalInfo string

	switch extenders.ExtenderName(updatedBy) {
	case "", static_extender.Name:
		actor = fmt.Sprintf("%s bundle", bundleName)
		if !moduleEnabled && module.Properties.Source != "Embedded" {
			additionalInfo = ", apply module config to enable"
		}

	case dynamic_extender.Name:
		actor = "dynamic global hook"

	case kube_config_extender.Name:
		actor = "module config"

	case script_extender.Name:
		actor = "'enabled'-script"
		additionalInfo = ", refer to the module's documentation"

	case d7s_version_extender.Name:
		actor = "DechouseVersion constraint"
		_, errMsg := d7s_version_extender.Instance().Filter(moduleName, map[string]string{})
		additionalInfo = fmt.Sprintf(", %s", errMsg)

	case bootstrapped_extender.Name:
		actor = "ClusterBoostrapped constraint"
		additionalInfo = ", the cluster isn't bootstrapped yet"

	case k8s_version_extender.Name:
		actor = "KubernetesVersion constraint"
		_, errMsg := k8s_version_extender.Instance().Filter(moduleName, map[string]string{})
		additionalInfo = fmt.Sprintf(", %s", errMsg)
	default:
		actor = updatedBy
	}

	if updatedByErr != nil {
		msgs = append(msgs, fmt.Sprintf("couldn't get the module's additional info: %s", updatedByErr))
	} else {
		msgs = append(msgs, fmt.Sprintf("turned %s by %s%s", modeDescriptor, actor, additionalInfo))
	}

	return ModuleStatus{
		Status:     strings.Join(statusMsgs, ", "),
		Message:    fmt.Sprintf("Info: %s", strings.Join(msgs, ", ")),
		HooksState: mod.GetHookErrorsSummary(),
	}
}

func (s *StatusReporter) ForConfig(cfg *v1alpha1.ModuleConfig) ModuleConfigStatus {
	statusMsgs := make([]string, 0)

	// Special case: unknown module name.
	if !s.possibleNames.Has(cfg.GetName()) {
		return ModuleConfigStatus{
			Version: "",
			Message: "Ignored: unknown module name",
		}
	}

	res := Service().ConfigValidator().Validate(cfg)
	if res.HasError() {
		return ModuleConfigStatus{
			Version: "",
			Message: fmt.Sprintf("Error: %s", res.Error),
		}
	}

	// Fill the 'version' field. The value is a spec.version or the latest version from registered conversions.
	// Also create warning if version is unknown or outdated.
	versionWarning := ""
	version := ""
	converter := conversion.Store().Get(cfg.GetName())
	if cfg.Spec.Version == 0 {
		// Use latest version if spec.version is empty.
		version = strconv.Itoa(converter.LatestVersion())
	}
	if cfg.Spec.Version > 0 {
		version = strconv.Itoa(cfg.Spec.Version)
		if !converter.IsKnownVersion(cfg.Spec.Version) {
			versionWarning = fmt.Sprintf("Error: invalid spec.version, use version %d", converter.LatestVersion())
		} else if cfg.Spec.Version < converter.LatestVersion() {
			// Warn about obsolete version if there is conversion for spec.version.
			versionWarning = fmt.Sprintf("Update available, latest spec.settings schema version is %d", converter.LatestVersion())
		}
	}

	// 'global' config is always enabled.
	if cfg.GetName() == "global" {
		return ModuleConfigStatus{
			Version: version,
			Message: versionWarning,
		}
	}

	if versionWarning != "" {
		statusMsgs = append(statusMsgs, versionWarning)
	}

	return ModuleConfigStatus{
		Version: version,
		Message: strings.Join(statusMsgs, ", "),
	}
}
