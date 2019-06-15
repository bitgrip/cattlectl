// Copyright © 2019 Bitgrip <berlin@bitgrip.de>
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

package client

import (
	"fmt"

	projectModel "github.com/bitgrip/cattlectl/internal/pkg/rancher/cluster/project/model"
	"github.com/rancher/norman/types"
	"github.com/sirupsen/logrus"
)

func newDeploymentClientWithData(
	deployment projectModel.Deployment,
	namespace string,
	project ProjectClient,
	logger *logrus.Entry,
) (DeploymentClient, error) {
	result, err := newDeploymentClient(
		deployment.Name,
		namespace,
		project,
		logger,
	)
	if err != nil {
		return nil, err
	}
	err = result.SetData(deployment)
	return result, err
}

func newDeploymentClient(
	name, namespace string,
	project ProjectClient,
	logger *logrus.Entry,
) (DeploymentClient, error) {
	return &deploymentClient{
		namespacedResourceClient: namespacedResourceClient{
			resourceClient: resourceClient{
				name:   name,
				logger: logger.WithField("deployment_name", name).WithField("namespace", namespace),
			},
			namespace: namespace,
			project:   project,
		},
	}, nil
}

type deploymentClient struct {
	namespacedResourceClient
	deployment projectModel.Deployment
}

func (client *deploymentClient) Exists() (bool, error) {
	backendClient, err := client.project.backendProjectClient()
	if err != nil {
		return false, err
	}
	namespaceID, err := client.NamespaceID()
	if err != nil {
		return false, err
	}
	collection, err := backendClient.Deployment.List(&types.ListOpts{
		Filters: map[string]interface{}{
			"name":        client.name,
			"namespaceId": namespaceID,
		},
	})
	if nil != err {
		client.logger.WithError(err).Error("Failed to read deployment list")
		return false, fmt.Errorf("Failed to read deployment list, %v", err)
	}
	for _, item := range collection.Data {
		if item.Name == client.name && item.NamespaceId == namespaceID {
			return true, nil
		}
	}
	client.logger.Debug("Deployment not found")
	return false, nil
}

func (client *deploymentClient) Create() error {
	backendClient, err := client.project.backendProjectClient()
	if err != nil {
		return err
	}
	namespaceID, err := client.NamespaceID()
	if err != nil {
		return err
	}
	client.logger.Info("Create new deployment")
	pattern, err := projectModel.ConvertDeploymentToProjectAPI(client.deployment)
	if err != nil {
		return err
	}
	pattern.NamespaceId = namespaceID
	_, err = backendClient.Deployment.Create(&pattern)
	return err
}

func (client *deploymentClient) Upgrade() error {
	client.logger.Warn("Skip change existing deployment")
	return nil
}

func (client *deploymentClient) Data() (projectModel.Deployment, error) {
	return client.deployment, nil
}

func (client *deploymentClient) SetData(deployment projectModel.Deployment) error {
	client.name = deployment.Name
	client.deployment = deployment
	return nil
}
