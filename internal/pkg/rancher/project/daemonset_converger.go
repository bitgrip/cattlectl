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

// NewDaemonSetConverger creates a Converger for a given github.com/bitgrip/cattlectl/internal/pkg/projectModel.DaemonSetDescriptor
func NewDaemonSetConverger(daemonSetDescriptor projectModel.DaemonSetDescriptor) Converger {
	client, err := newRancherClient(rancher.ClientConfig{
		RancherURL: daemonSetDescriptor.Metadata.RancherURL,
		AccessKey:  daemonSetDescriptor.Metadata.AccessKey,
		SecretKey:  daemonSetDescriptor.Metadata.SecretKey,
	})
	if err != nil {
		logrus.WithError(err).Fatal("Error creating rancher client")
	}
	return daemonSetConverger{
		daemonSetDescriptor: daemonSetDescriptor,
		client:              client,
	}
}

type daemonSetConverger struct {
	daemonSetDescriptor projectModel.DaemonSetDescriptor
	client              rancher.Client
}

func (converger daemonSetConverger) Converge() error {
	if err := converger.initCluster(); err != nil {
		return err
	}
	if err := converger.initProject(); err != nil {
		return err
	}

	if hasDaemonSet, err := converger.client.HasDaemonSet(converger.daemonSetDescriptor.Metadata.Namespace, converger.daemonSetDescriptor.Spec); hasDaemonSet {
		logrus.WithFields(logrus.Fields{
			"project_name":   converger.daemonSetDescriptor.Metadata.ProjectName,
			"namespace":      converger.daemonSetDescriptor.Metadata.Namespace,
			"daemonSet_name": converger.daemonSetDescriptor.Spec.Name,
		}).Warn("DaemonSet exists need to be removed manually")
		return fmt.Errorf("Can not override existing daemonSet %v", converger.daemonSetDescriptor.Spec.Name)
	} else if err != nil {
		return fmt.Errorf("Failed to check for namespace, %v", err)
	}
	return converger.client.CreateDaemonSet(converger.daemonSetDescriptor.Metadata.Namespace, converger.daemonSetDescriptor.Spec)
}

func (converger daemonSetConverger) initCluster() error {
	if converger.daemonSetDescriptor.Metadata.ClusterID != "" {
		if err := converger.client.SetCluster(converger.daemonSetDescriptor.Metadata.ClusterName, converger.daemonSetDescriptor.Metadata.ClusterID); err != nil {
			return fmt.Errorf("Failed to init cluster, %v", err)
		}
		return nil
	}
	if hasCluster, clusterID, err := converger.client.HasClusterWithName(converger.daemonSetDescriptor.Metadata.ClusterName); hasCluster {
		if err = converger.client.SetCluster(converger.daemonSetDescriptor.Metadata.ClusterName, clusterID); err != nil {
			return fmt.Errorf("Failed to init cluster, %v", err)
		}
		return nil
	} else if err != nil {
		return fmt.Errorf("Failed to init cluster, %v", err)
	} else {
		return fmt.Errorf("Cluster not found")
	}
}
func (converger daemonSetConverger) initProject() error {
	if hasProject, projectID, err := converger.client.HasProjectWithName(converger.daemonSetDescriptor.Metadata.ProjectName); hasProject {
		if err = converger.client.SetProject(converger.daemonSetDescriptor.Metadata.ProjectName, projectID); err != nil {
			logrus.WithError(err).WithFields(logrus.Fields{"project_name": converger.daemonSetDescriptor.Metadata.ProjectName}).Warn("Failed to init project")
			return fmt.Errorf("Failed to init project, %v", err)
		}
		return nil
	} else if err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{"project_name": converger.daemonSetDescriptor.Metadata.ProjectName}).Warn("Failed to check for project")
		return fmt.Errorf("Failed to check for project, %v", err)
	}
	logrus.WithFields(logrus.Fields{"project_name": converger.daemonSetDescriptor.Metadata.ProjectName}).Warn("Failed to check for project")
	return fmt.Errorf("Project not found")
}
