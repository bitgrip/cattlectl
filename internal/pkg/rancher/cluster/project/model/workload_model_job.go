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

// JobDescriptor is a workload representing a K8S Job
type JobDescriptor struct {
	APIVersion string `yaml:"api_version"`
	Kind       string
	Metadata   WorkloadMetadata
	Spec       Job
}

type Job struct {
	baseWorkload            `yaml:"baseWorkload,inline"`
	JobConfig               *JobConfig `json:"jobConfig,omitempty" yaml:"jobConfig,omitempty"`
	TTLSecondsAfterFinished *int64     `yaml:"TTLSecondsAfterFinished,omitempty"`
}

type JobConfig struct {
	ActiveDeadlineSeconds *int64 `json:"activeDeadlineSeconds,omitempty" yaml:"activeDeadlineSeconds,omitempty"`
	BackoffLimit          *int64 `json:"backoffLimit,omitempty" yaml:"backoffLimit,omitempty"`
	Completions           *int64 `json:"completions,omitempty" yaml:"completions,omitempty"`
	ManualSelector        *bool  `json:"manualSelector,omitempty" yaml:"manualSelector,omitempty"`
	Parallelism           *int64 `json:"parallelism,omitempty" yaml:"parallelism,omitempty"`
}

func ConvertJobDescriptorToProjectAPI(descriptor JobDescriptor) (projectAPI.Job, error) {
	return ConvertJobToProjectAPI(descriptor.Spec)
}

func ConvertJobToProjectAPI(job Job) (projectAPI.Job, error) {
	result := projectAPI.Job{}
	transferContent, err := json.Marshal(job)
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(transferContent, &result)
	if err != nil {
		return result, err
	}
	return result, nil
}
