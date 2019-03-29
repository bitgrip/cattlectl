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
	projectClient "github.com/rancher/types/client/project/v3"
)

func (client *rancherClient) HasConfigMap(configMap projectModel.ConfigMap) (bool, error) {
	collection, err := client.projectClient.ConfigMap.List(&types.ListOpts{
		Filters: map[string]interface{}{
			"system": "false",
			"name":   configMap.Name,
		},
	})
	if nil != err {
		client.logger.WithError(err).WithField("config_map_name", configMap.Name).Error("Failed to read config_map list")
		return false, fmt.Errorf("Failed to read config_map list, %v", err)
	}
	for _, item := range collection.Data {
		if item.Name == configMap.Name {
			return true, nil
		}
	}
	client.logger.WithError(err).WithField("config_map_name", configMap.Name).Debug("ConfigMap not found")
	return false, nil
}

func (client *rancherClient) CreateConfigMap(configMap projectModel.ConfigMap) error {
	client.logger.WithField("config_map_name", configMap.Name).Info("Create new ConfigMap")
	newConfigMap := &projectClient.ConfigMap{
		Name:        configMap.Name,
		Data:        configMap.Data,
		NamespaceId: configMap.Namespace,
		ProjectID:   client.projectId,
	}

	_, err := client.projectClient.ConfigMap.Create(newConfigMap)
	return err
}
