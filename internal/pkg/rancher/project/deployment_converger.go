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
	"fmt"

	"github.com/bitgrip/cattlectl/internal/pkg/rancher"
	projectModel "github.com/bitgrip/cattlectl/internal/pkg/rancher/project/model"
	"github.com/sirupsen/logrus"
)

// NewDeploymentConverger creates a Converger for a given github.com/bitgrip/cattlectl/internal/pkg/projectModel.JobDescriptor
func NewDeploymentConverger(deploymentDescriptor projectModel.DeploymentDescriptor) Converger {
	client, err := newRancherClient(rancher.ClientConfig{
		RancherURL: deploymentDescriptor.Metadata.RancherURL,
		AccessKey:  deploymentDescriptor.Metadata.AccessKey,
		SecretKey:  deploymentDescriptor.Metadata.SecretKey,
	})
	if err != nil {
		logrus.WithError(err).Fatal("Error creating rancher client")
	}
	return deploymentConverger{
		deploymentDescriptor: deploymentDescriptor,
		client:               client,
	}
}

type deploymentConverger struct {
	deploymentDescriptor projectModel.DeploymentDescriptor
	client               rancher.Client
}

func (converger deploymentConverger) Converge() error {
	if err := converger.initCluster(); err != nil {
		return err
	}
	if err := converger.initProject(); err != nil {
		return err
	}

	if hasJob, err := converger.client.HasDeployment(converger.deploymentDescriptor.Metadata.Namespace, converger.deploymentDescriptor.Spec); hasJob {
		logrus.WithFields(logrus.Fields{
			"project_name":    converger.deploymentDescriptor.Metadata.ProjectName,
			"namespace":       converger.deploymentDescriptor.Metadata.Namespace,
			"deployment_name": converger.deploymentDescriptor.Spec.Name,
		}).Warn("Job exists need to be removed manually")
		return fmt.Errorf("Can not override existing job %v", converger.deploymentDescriptor.Spec.Name)
	} else if err != nil {
		return fmt.Errorf("Failed to check for namespace, %v", err)
	}
	return converger.client.CreateDeployment(converger.deploymentDescriptor.Metadata.Namespace, converger.deploymentDescriptor.Spec)
}

func (converger deploymentConverger) initCluster() error {
	if converger.deploymentDescriptor.Metadata.ClusterID != "" {
		if err := converger.client.SetCluster(converger.deploymentDescriptor.Metadata.ClusterName, converger.deploymentDescriptor.Metadata.ClusterID); err != nil {
			return fmt.Errorf("Failed to init cluster, %v", err)
		}
		return nil
	}
	if hasCluster, clusterID, err := converger.client.HasClusterWithName(converger.deploymentDescriptor.Metadata.ClusterName); hasCluster {
		if err = converger.client.SetCluster(converger.deploymentDescriptor.Metadata.ClusterName, clusterID); err != nil {
			return fmt.Errorf("Failed to init cluster, %v", err)
		}
		return nil
	} else if err != nil {
		return fmt.Errorf("Failed to init cluster, %v", err)
	} else {
		return fmt.Errorf("Cluster not found")
	}
}
func (converger deploymentConverger) initProject() error {
	if hasProject, projectID, err := converger.client.HasProjectWithName(converger.deploymentDescriptor.Metadata.ProjectName); hasProject {
		if err = converger.client.SetProject(converger.deploymentDescriptor.Metadata.ProjectName, projectID); err != nil {
			logrus.WithError(err).WithFields(logrus.Fields{"project_name": converger.deploymentDescriptor.Metadata.ProjectName}).Warn("Failed to init project")
			return fmt.Errorf("Failed to init project, %v", err)
		}
		return nil
	} else if err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{"project_name": converger.deploymentDescriptor.Metadata.ProjectName}).Warn("Failed to check for project")
		return fmt.Errorf("Failed to check for project, %v", err)
	}
	logrus.WithFields(logrus.Fields{"project_name": converger.deploymentDescriptor.Metadata.ProjectName}).Warn("Failed to check for project")
	return fmt.Errorf("Project not found")
}
