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

type CronJobDescriptor struct {
	APIVersion string `yaml:"api_version"`
	Kind       string
	Metadata   WorkloadMetadata
	Spec       CronJob
}

type CronJob struct {
	baseWorkload            `yaml:"baseWorkload,inline"`
	CronJobConfig           *CronJobConfig `json:"cronJobConfig,omitempty" yaml:"cronJobConfig,omitempty"`
	TTLSecondsAfterFinished *int64         `yaml:"ttlSecondsAfterFinished,omitempty"`
}

type CronJobConfig struct {
	ConcurrencyPolicy          string            `json:"concurrencyPolicy,omitempty" yaml:"concurrencyPolicy,omitempty"`
	FailedJobsHistoryLimit     *int64            `json:"failedJobsHistoryLimit,omitempty" yaml:"failedJobsHistoryLimit,omitempty"`
	JobAnnotations             map[string]string `json:"jobAnnotations,omitempty" yaml:"jobAnnotations,omitempty"`
	JobConfig                  *JobConfig        `json:"jobConfig,omitempty" yaml:"jobConfig,omitempty"`
	JobLabels                  map[string]string `json:"jobLabels,omitempty" yaml:"jobLabels,omitempty"`
	Schedule                   string            `json:"schedule,omitempty" yaml:"schedule,omitempty"`
	StartingDeadlineSeconds    *int64            `json:"startingDeadlineSeconds,omitempty" yaml:"startingDeadlineSeconds,omitempty"`
	SuccessfulJobsHistoryLimit *int64            `json:"successfulJobsHistoryLimit,omitempty" yaml:"successfulJobsHistoryLimit,omitempty"`
	Suspend                    *bool             `json:"suspend,omitempty" yaml:"suspend,omitempty"`
}

func ConvertCronJobDescriptorToProjectAPI(descriptor CronJobDescriptor) (projectAPI.CronJob, error) {
	return ConvertCronJobToProjectAPI(descriptor.Spec)
}
func ConvertCronJobToProjectAPI(cronJob CronJob) (projectAPI.CronJob, error) {
	result := projectAPI.CronJob{}
	transferContent, err := json.Marshal(cronJob)
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(transferContent, &result)
	if err != nil {
		return result, err
	}
	return result, nil
}
