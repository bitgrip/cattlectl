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

package rancher

import (
	"github.com/rancher/norman/types"
	managementClient "github.com/rancher/types/client/management/v3"
	"github.com/sirupsen/logrus"
)

func (client *rancherClient) HasProjectWithName(name string) (bool, string, error) {
	collection, err := client.managementClient.Project.List(&types.ListOpts{
		Filters: map[string]interface{}{
			"clusterId": client.clusterId,
			"name":      name,
		},
	})
	if err != nil {
		return false, "", err
	}
	if len(collection.Data) < 1 {
		return false, "", nil
	}
	return true, collection.Data[0].ID, nil
}

func (client *rancherClient) SetProject(projectName, projectId string) error {
	if projectClient, err := createProjectClient(
		client.clientConfig.RancherURL,
		client.clientConfig.AccessKey,
		client.clientConfig.SecretKey,
		client.clusterId,
		projectId,
	); err != nil {
		return err
	} else {
		client.logger = client.logger.WithFields(logrus.Fields{
			"project_name": projectName,
		})
		client.projectId = projectId
		client.projectClient = projectClient
		return nil
	}
}

func (client *rancherClient) CreateProject(projectName string) (string, error) {
	client.logger.WithField("project_name", projectName).Info("Create new project")
	pattern := &managementClient.Project{
		ClusterID: client.clusterId,
		Name:      projectName,
	}
	createdProject, err := client.managementClient.Project.Create(pattern)
	if err != nil {
		client.logger.Warn("Failed to create project")
		return "", err
	}

	return createdProject.ID, nil
}
