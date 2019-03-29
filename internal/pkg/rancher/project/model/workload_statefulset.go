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

type StatefulSetDescriptor struct {
	APIVersion string `yaml:"api_version"`
	Kind       string
	Metadata   WorkloadMetadata
	Spec       StatefulSet
}

type StatefulSet struct {
	baseWorkload      `yaml:"baseWorkload,inline"`
	StatefulSetConfig *StatefulSetConfig `json:"statefulSetConfig,omitempty" yaml:"statefulSetConfig,omitempty"`
	Scale             *int64             `json:"scale,omitempty" yaml:"scale,omitempty"`
}

type StatefulSetConfig struct {
	Partition            *int64                  `json:"partition,omitempty" yaml:"partition,omitempty"`
	PodManagementPolicy  string                  `json:"podManagementPolicy,omitempty" yaml:"podManagementPolicy,omitempty"`
	RevisionHistoryLimit *int64                  `json:"revisionHistoryLimit,omitempty" yaml:"revisionHistoryLimit,omitempty"`
	ServiceName          string                  `json:"serviceName,omitempty" yaml:"serviceName,omitempty"`
	Strategy             string                  `json:"strategy,omitempty" yaml:"strategy,omitempty"`
	VolumeClaimTemplates []PersistentVolumeClaim `json:"volumeClaimTemplates,omitempty" yaml:"volumeClaimTemplates,omitempty"`
}

type PersistentVolumeClaim struct {
	AccessModes  []string              `json:"accessModes,omitempty" yaml:"accessModes,omitempty"`
	Annotations  map[string]string     `json:"annotations,omitempty" yaml:"annotations,omitempty"`
	Labels       map[string]string     `json:"labels,omitempty" yaml:"labels,omitempty"`
	Name         string                `json:"name,omitempty" yaml:"name,omitempty"`
	Resources    *ResourceRequirements `json:"resources,omitempty" yaml:"resources,omitempty"`
	Selector     *LabelSelector        `json:"selector,omitempty" yaml:"selector,omitempty"`
	StorageClass string                `json:"storageClass,omitempty" yaml:"storageClass,omitempty"`
}

func ConvertStatefulSetDescriptorToProjectAPI(descriptor StatefulSetDescriptor) (projectAPI.StatefulSet, error) {
	return ConvertStatefulSetToProjectAPI(descriptor.Spec)
}

func ConvertStatefulSetToProjectAPI(statefulSet StatefulSet) (projectAPI.StatefulSet, error) {
	result := projectAPI.StatefulSet{}
	transferContent, err := json.Marshal(statefulSet)
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(transferContent, &result)
	if err != nil {
		return result, err
	}
	return result, nil
}
