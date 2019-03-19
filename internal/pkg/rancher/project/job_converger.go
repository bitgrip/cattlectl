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

// NewJobConverger creates a Converger for a given github.com/bitgrip/cattlectl/internal/pkg/projectModel.JobDescriptor
func NewJobConverger(jobDescriptor projectModel.JobDescriptor) Converger {
	client, err := newRancherClient(rancher.ClientConfig{
		RancherURL: jobDescriptor.Metadata.RancherURL,
		AccessKey:  jobDescriptor.Metadata.AccessKey,
		SecretKey:  jobDescriptor.Metadata.SecretKey,
	})
	if err != nil {
		logrus.WithError(err).Fatal("Error creating rancher client")
	}
	return jobConverger{
		jobDescriptor: jobDescriptor,
		client:        client,
	}
}

type jobConverger struct {
	jobDescriptor projectModel.JobDescriptor
	client        rancher.Client
}

func (converger jobConverger) Converge() error {
	if err := converger.initCluster(); err != nil {
		return err
	}
	if err := converger.initProject(); err != nil {
		return err
	}

	if hasJob, err := converger.client.HasJob(converger.jobDescriptor.Metadata.Namespace, converger.jobDescriptor.Spec); hasJob {
		logrus.WithFields(logrus.Fields{
			"project_name": converger.jobDescriptor.Metadata.ProjectName,
			"namespace":    converger.jobDescriptor.Metadata.Namespace,
			"job_name":     converger.jobDescriptor.Spec.Name,
		}).Warn("Job exists need to be removed manually")
		return fmt.Errorf("Can not override existing job %v", converger.jobDescriptor.Spec.Name)
	} else if err != nil {
		return fmt.Errorf("Failed to check for namespace, %v", err)
	}
	return converger.client.CreateJob(converger.jobDescriptor.Metadata.Namespace, converger.jobDescriptor.Spec)
}

func (converger jobConverger) initCluster() error {
	if converger.jobDescriptor.Metadata.ClusterID != "" {
		if err := converger.client.SetCluster(converger.jobDescriptor.Metadata.ClusterName, converger.jobDescriptor.Metadata.ClusterID); err != nil {
			return fmt.Errorf("Failed to init cluster, %v", err)
		}
		return nil
	}
	if hasCluster, clusterID, err := converger.client.HasClusterWithName(converger.jobDescriptor.Metadata.ClusterName); hasCluster {
		if err = converger.client.SetCluster(converger.jobDescriptor.Metadata.ClusterName, clusterID); err != nil {
			return fmt.Errorf("Failed to init cluster, %v", err)
		}
		return nil
	} else if err != nil {
		return fmt.Errorf("Failed to init cluster, %v", err)
	} else {
		return fmt.Errorf("Cluster not found")
	}
}
func (converger jobConverger) initProject() error {
	if hasProject, projectID, err := converger.client.HasProjectWithName(converger.jobDescriptor.Metadata.ProjectName); hasProject {
		if err = converger.client.SetProject(converger.jobDescriptor.Metadata.ProjectName, projectID); err != nil {
			logrus.WithError(err).WithFields(logrus.Fields{"project_name": converger.jobDescriptor.Metadata.ProjectName}).Warn("Failed to init project")
			return fmt.Errorf("Failed to init project, %v", err)
		}
		return nil
	} else if err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{"project_name": converger.jobDescriptor.Metadata.ProjectName}).Warn("Failed to check for project")
		return fmt.Errorf("Failed to check for project, %v", err)
	}
	logrus.WithFields(logrus.Fields{"project_name": converger.jobDescriptor.Metadata.ProjectName}).Warn("Failed to check for project")
	return fmt.Errorf("Project not found")
}
