// Copyright © 2018 Bitgrip <berlin@bitgrip.de>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package model

import (
	"encoding/json"

	projectAPI "github.com/rancher/types/client/project/v3"
)

type DaemonSetDescriptor struct {
	APIVersion string `yaml:"api_version"`
	Kind       string
	Metadata   WorkloadMetadata
	Spec       DaemonSet
}

type DaemonSet struct {
	baseWorkload    `yaml:"baseWorkload,inline"`
	DaemonSetConfig *DaemonSetConfig `json:"daemonSetConfig,omitempty" yaml:"daemonSetConfig,omitempty"`
}

type DaemonSetConfig struct {
	MaxUnavailable       IntOrString `json:"maxUnavailable,omitempty" yaml:"maxUnavailable,omitempty"`
	MinReadySeconds      int64       `json:"minReadySeconds,omitempty" yaml:"minReadySeconds,omitempty"`
	RevisionHistoryLimit *int64      `json:"revisionHistoryLimit,omitempty" yaml:"revisionHistoryLimit,omitempty"`
	Strategy             string      `json:"strategy,omitempty" yaml:"strategy,omitempty"`
}

func ConvertDaemonSetDescriptorToProjectAPI(descriptor DaemonSetDescriptor) (projectAPI.DaemonSet, error) {
	return ConvertDaemonSetToProjectAPI(descriptor.Spec)
}

func ConvertDaemonSetToProjectAPI(daemonSet DaemonSet) (projectAPI.DaemonSet, error) {
	result := projectAPI.DaemonSet{}
	transferContent, err := json.Marshal(daemonSet)
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(transferContent, &result)
	if err != nil {
		return result, err
	}
	return result, nil
}
