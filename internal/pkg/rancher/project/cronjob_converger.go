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
	"github.com/bitgrip/cattlectl/internal/pkg/rancher"
	"github.com/bitgrip/cattlectl/internal/pkg/rancher/descriptor"
	projectModel "github.com/bitgrip/cattlectl/internal/pkg/rancher/project/model"
	"github.com/sirupsen/logrus"
)

// NewCronJobConverger creates a Converger for a given github.com/bitgrip/cattlectl/internal/pkg/projectModel.JobDescriptor
func NewCronJobConverger(cronJobDescriptor projectModel.CronJobDescriptor) descriptor.Converger {
	return descriptor.ProjectResourceDescriptorConverger(
		cronJobDescriptor.Metadata.ClusterName,
		cronJobDescriptor.Metadata.ProjectName,
		[]descriptor.Converger{
			newCronJobPartConverger(
				cronJobDescriptor.Metadata.ProjectName,
				cronJobDescriptor.Metadata.Namespace,
				cronJobDescriptor.Spec,
			),
		},
	)
}

func newCronJobPartConverger(projectName, namespace string, cronJob projectModel.CronJob) descriptor.Converger {
	return descriptor.PartConverger{
		PartName: "CronJob",
		HasPart: func(client rancher.Client) (bool, error) {
			return client.HasCronJob(namespace, cronJob)
		},
		UpdatePart: func(client rancher.Client) error {
			logrus.WithFields(logrus.Fields{
				"project_name": projectName,
				"namespace":    namespace,
				"cronjob_name": cronJob.Name,
			}).Warn("CronJob exists need to be removed manually")
			return nil
		},
		CreatePart: func(client rancher.Client) error {
			return client.CreateCronJob(namespace, cronJob)
		},
	}
}
