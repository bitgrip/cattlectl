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

package project

import (
	"github.com/bitgrip/cattlectl/internal/pkg/rancher/client"
	"github.com/bitgrip/cattlectl/internal/pkg/rancher/descriptor"
	projectModel "github.com/bitgrip/cattlectl/internal/pkg/rancher/cluster/project/model"
)

// NewCronJobConverger creates a Converger for a given github.com/bitgrip/cattlectl/internal/pkg/projectModel.JobDescriptor
func NewCronJobConverger(cronJobDescriptor projectModel.CronJobDescriptor, projectClient client.ProjectClient) (descriptor.Converger, error) {
	cronJobClient, err := projectClient.CronJob(cronJobDescriptor.Spec.Name, cronJobDescriptor.Metadata.Namespace)
	if err != nil {
		return nil, err
	}
	err = cronJobClient.SetData(cronJobDescriptor.Spec)
	if err != nil {
		return nil, err
	}
	return &descriptor.ResourceClientConverger{
		Client: cronJobClient,
	}, nil
}
