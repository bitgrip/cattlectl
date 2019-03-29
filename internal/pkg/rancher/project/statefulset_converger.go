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

// NewStatefulSetConverger creates a Converger for a given github.com/bitgrip/cattlectl/internal/pkg/projectModel.StatefulSetDescriptor
func NewStatefulSetConverger(statefulSetDescriptor projectModel.StatefulSetDescriptor) Converger {
	client, err := newRancherClient(rancher.ClientConfig{
		RancherURL: statefulSetDescriptor.Metadata.RancherURL,
		AccessKey:  statefulSetDescriptor.Metadata.AccessKey,
		SecretKey:  statefulSetDescriptor.Metadata.SecretKey,
	})
	if err != nil {
		logrus.WithError(err).Fatal("Error creating rancher client")
	}
	return statefulSetConverger{
		statefulSetDescriptor: statefulSetDescriptor,
		client:                client,
	}
}

type statefulSetConverger struct {
	statefulSetDescriptor projectModel.StatefulSetDescriptor
	client                rancher.Client
}

func (converger statefulSetConverger) Converge() error {
	if err := converger.initCluster(); err != nil {
		return err
	}
	if err := converger.initProject(); err != nil {
		return err
	}

	if hasStatefulSet, err := converger.client.HasStatefulSet(converger.statefulSetDescriptor.Metadata.Namespace, converger.statefulSetDescriptor.Spec); hasStatefulSet {
		logrus.WithFields(logrus.Fields{
			"project_name":     converger.statefulSetDescriptor.Metadata.ProjectName,
			"namespace":        converger.statefulSetDescriptor.Metadata.Namespace,
			"statefulSet_name": converger.statefulSetDescriptor.Spec.Name,
		}).Warn("StatefulSet exists need to be removed manually")
		return fmt.Errorf("Can not override existing statefulSet %v", converger.statefulSetDescriptor.Spec.Name)
	} else if err != nil {
		return fmt.Errorf("Failed to check for namespace, %v", err)
	}
	return converger.client.CreateStatefulSet(converger.statefulSetDescriptor.Metadata.Namespace, converger.statefulSetDescriptor.Spec)
}

func (converger statefulSetConverger) initCluster() error {
	if converger.statefulSetDescriptor.Metadata.ClusterID != "" {
		if err := converger.client.SetCluster(converger.statefulSetDescriptor.Metadata.ClusterName, converger.statefulSetDescriptor.Metadata.ClusterID); err != nil {
			return fmt.Errorf("Failed to init cluster, %v", err)
		}
		return nil
	}
	if hasCluster, clusterID, err := converger.client.HasClusterWithName(converger.statefulSetDescriptor.Metadata.ClusterName); hasCluster {
		if err = converger.client.SetCluster(converger.statefulSetDescriptor.Metadata.ClusterName, clusterID); err != nil {
			return fmt.Errorf("Failed to init cluster, %v", err)
		}
		return nil
	} else if err != nil {
		return fmt.Errorf("Failed to init cluster, %v", err)
	} else {
		return fmt.Errorf("Cluster not found")
	}
}
func (converger statefulSetConverger) initProject() error {
	if hasProject, projectID, err := converger.client.HasProjectWithName(converger.statefulSetDescriptor.Metadata.ProjectName); hasProject {
		if err = converger.client.SetProject(converger.statefulSetDescriptor.Metadata.ProjectName, projectID); err != nil {
			logrus.WithError(err).WithFields(logrus.Fields{"project_name": converger.statefulSetDescriptor.Metadata.ProjectName}).Warn("Failed to init project")
			return fmt.Errorf("Failed to init project, %v", err)
		}
		return nil
	} else if err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{"project_name": converger.statefulSetDescriptor.Metadata.ProjectName}).Warn("Failed to check for project")
		return fmt.Errorf("Failed to check for project, %v", err)
	}
	logrus.WithFields(logrus.Fields{"project_name": converger.statefulSetDescriptor.Metadata.ProjectName}).Warn("Failed to check for project")
	return fmt.Errorf("Project not found")
}
