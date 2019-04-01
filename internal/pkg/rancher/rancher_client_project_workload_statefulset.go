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

func (client *rancherClient) HasStatefulSet(namespace string, statefulSet projectModel.StatefulSet) (bool, error) {
	namespaceID, err := client.getNamespaceID(namespace)
	if err != nil {
		client.logger.WithError(err).WithField("statefulSet_name", statefulSet.Name).WithField("namespace", namespace).Error("Failed to read statefulSet list")
		return false, fmt.Errorf("Failed to read statefulSet list, %v", err)
	}
	collection, err := client.projectClient.StatefulSet.List(&types.ListOpts{
		Filters: map[string]interface{}{
			"name":        statefulSet.Name,
			"namespaceId": namespaceID,
		},
	})
	if nil != err {
		client.logger.WithError(err).WithField("statefulSet_name", statefulSet.Name).WithField("namespaceId", namespaceID).Error("Failed to read statefulSet list")
		return false, fmt.Errorf("Failed to read statefulSet list, %v", err)
	}
	for _, item := range collection.Data {
		if item.Name == statefulSet.Name {
			return true, nil
		}
	}
	client.logger.WithField("statefulSet_name", statefulSet.Name).WithField("namespaceId", namespaceID).Debug("StatefulSet not found")
	return false, nil
}

func (client *rancherClient) CreateStatefulSet(namespace string, statefulSet projectModel.StatefulSet) error {
	client.logger.WithField("statefulSet_name", statefulSet.Name).Info("Create new statefulSet")
	namespaceId, err := client.getNamespaceID(namespace)
	if err != nil {
		return err
	}
	pattern, err := projectModel.ConvertStatefulSetToProjectAPI(statefulSet)
	if err != nil {
		return err
	}
	pattern.NamespaceId = namespaceId
	_, err = client.projectClient.StatefulSet.Create(&pattern)
	return err
}
