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

// NewDeploymentConverger creates a Converger for a given github.com/bitgrip/cattlectl/internal/pkg/projectModel.JobDescriptor
func NewDeploymentConverger(deploymentDescriptor projectModel.DeploymentDescriptor) descriptor.Converger {
	return descriptor.DescriptorConverger{
		InitCluster: func(client rancher.Client) error {
			return rancher.InitCluster(
				deploymentDescriptor.Metadata.ClusterID,
				deploymentDescriptor.Metadata.ClusterName,
				client,
			)
		},
		InitProject: func(client rancher.Client) error {
			return rancher.InitProject(
				deploymentDescriptor.Metadata.ProjectName,
				client,
			)
		},
		PartConvergers: []descriptor.Converger{
			newDeploymentPartConverger(
				deploymentDescriptor.Metadata.ProjectName,
				deploymentDescriptor.Metadata.Namespace,
				deploymentDescriptor.Spec,
			),
		},
	}
}

func newDeploymentPartConverger(projectName, namespace string, deployment projectModel.Deployment) descriptor.Converger {
	return descriptor.PartConverger{
		PartName: "Deployment",
		HasPart: func(client rancher.Client) (bool, error) {
			return client.HasDeployment(namespace, deployment)
		},
		UpdatePart: func(client rancher.Client) error {
			logrus.WithFields(logrus.Fields{
				"project_name":    projectName,
				"namespace":       namespace,
				"deployment_name": deployment.Name,
			}).Warn("Deployment exists need to be removed manually")
			return nil
		},
		CreatePart: func(client rancher.Client) error {
			return client.CreateDeployment(namespace, deployment)
		},
	}
}
