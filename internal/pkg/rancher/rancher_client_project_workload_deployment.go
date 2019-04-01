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
	"fmt"

	projectModel "github.com/bitgrip/cattlectl/internal/pkg/rancher/project/model"
	"github.com/rancher/norman/types"
)

func (client *rancherClient) HasDeployment(namespace string, deployment projectModel.Deployment) (bool, error) {
	namespaceID, err := client.getNamespaceID(namespace)
	if err != nil {
		client.logger.WithError(err).WithField("cronjob_name", deployment.Name).WithField("namespace", namespace).Error("Failed to read Deployment list")
		return false, fmt.Errorf("Failed to read Deployment list, %v", err)
	}
	collection, err := client.projectClient.Deployment.List(&types.ListOpts{
		Filters: map[string]interface{}{
			"name":        deployment.Name,
			"namespaceId": namespaceID,
		},
	})
	if nil != err {
		client.logger.WithError(err).WithField("cronjob_name", deployment.Name).WithField("namespaceId", namespaceID).Error("Failed to read Deployment list")
		return false, fmt.Errorf("Failed to read Deployment list, %v", err)
	}
	for _, item := range collection.Data {
		if item.Name == deployment.Name {
			return true, nil
		}
	}
	client.logger.WithField("cronjob_name", deployment.Name).WithField("namespaceId", namespaceID).Debug("Deployment not found")
	return false, nil
}

func (client *rancherClient) CreateDeployment(namespace string, deployment projectModel.Deployment) error {
	client.logger.WithField("cronjob_name", deployment.Name).Info("Create new Deployment")
	namespaceId, err := client.getNamespaceID(namespace)
	if err != nil {
		return err
	}
	pattern, err := projectModel.ConvertDeploymentToProjectAPI(deployment)
	if err != nil {
		return err
	}
	pattern.NamespaceId = namespaceId
	_, err = client.projectClient.Deployment.Create(&pattern)
	return err
}
