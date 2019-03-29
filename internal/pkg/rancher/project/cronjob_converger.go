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

// NewCronJobConverger creates a Converger for a given github.com/bitgrip/cattlectl/internal/pkg/projectModel.JobDescriptor
func NewCronJobConverger(cronJobDescriptor projectModel.CronJobDescriptor) Converger {
	client, err := newRancherClient(rancher.ClientConfig{
		RancherURL: cronJobDescriptor.Metadata.RancherURL,
		AccessKey:  cronJobDescriptor.Metadata.AccessKey,
		SecretKey:  cronJobDescriptor.Metadata.SecretKey,
	})
	if err != nil {
		logrus.WithError(err).Fatal("Error creating rancher client")
	}
	return cronJobConverger{
		cronJobDescriptor: cronJobDescriptor,
		client:            client,
	}
}

type cronJobConverger struct {
	cronJobDescriptor projectModel.CronJobDescriptor
	client            rancher.Client
}

func (converger cronJobConverger) Converge() error {
	if err := converger.initCluster(); err != nil {
		return err
	}
	if err := converger.initProject(); err != nil {
		return err
	}

	if hasJob, err := converger.client.HasCronJob(converger.cronJobDescriptor.Metadata.Namespace, converger.cronJobDescriptor.Spec); hasJob {
		logrus.WithFields(logrus.Fields{
			"project_name": converger.cronJobDescriptor.Metadata.ProjectName,
			"namespace":    converger.cronJobDescriptor.Metadata.Namespace,
			"job_name":     converger.cronJobDescriptor.Spec.Name,
		}).Warn("Job exists need to be removed manually")
		return fmt.Errorf("Can not override existing job %v", converger.cronJobDescriptor.Spec.Name)
	} else if err != nil {
		return fmt.Errorf("Failed to check for namespace, %v", err)
	}
	return converger.client.CreateCronJob(converger.cronJobDescriptor.Metadata.Namespace, converger.cronJobDescriptor.Spec)
}

func (converger cronJobConverger) initCluster() error {
	if converger.cronJobDescriptor.Metadata.ClusterID != "" {
		if err := converger.client.SetCluster(converger.cronJobDescriptor.Metadata.ClusterName, converger.cronJobDescriptor.Metadata.ClusterID); err != nil {
			return fmt.Errorf("Failed to init cluster, %v", err)
		}
		return nil
	}
	if hasCluster, clusterID, err := converger.client.HasClusterWithName(converger.cronJobDescriptor.Metadata.ClusterName); hasCluster {
		if err = converger.client.SetCluster(converger.cronJobDescriptor.Metadata.ClusterName, clusterID); err != nil {
			return fmt.Errorf("Failed to init cluster, %v", err)
		}
		return nil
	} else if err != nil {
		return fmt.Errorf("Failed to init cluster, %v", err)
	} else {
		return fmt.Errorf("Cluster not found")
	}
}
func (converger cronJobConverger) initProject() error {
	if hasProject, projectID, err := converger.client.HasProjectWithName(converger.cronJobDescriptor.Metadata.ProjectName); hasProject {
		if err = converger.client.SetProject(converger.cronJobDescriptor.Metadata.ProjectName, projectID); err != nil {
			logrus.WithError(err).WithFields(logrus.Fields{"project_name": converger.cronJobDescriptor.Metadata.ProjectName}).Warn("Failed to init project")
			return fmt.Errorf("Failed to init project, %v", err)
		}
		return nil
	} else if err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{"project_name": converger.cronJobDescriptor.Metadata.ProjectName}).Warn("Failed to check for project")
		return fmt.Errorf("Failed to check for project, %v", err)
	}
	logrus.WithFields(logrus.Fields{"project_name": converger.cronJobDescriptor.Metadata.ProjectName}).Warn("Failed to check for project")
	return fmt.Errorf("Project not found")
}
