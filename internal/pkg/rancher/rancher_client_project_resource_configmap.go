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
	if _, exists := client.configMapCache[configMap.Namespace+"/"+configMap.Name]; exists {
		return true, nil
	}
	namespaceID, err := client.getNamespaceID(configMap.Namespace)
	if err != nil {
		return false, fmt.Errorf("Failed to read config_map list, %v", err)
	}
	collection, err := client.projectClient.ConfigMap.List(&types.ListOpts{
		Filters: map[string]interface{}{
			"system":      "false",
			"name":        configMap.Name,
			"namespaceId": namespaceID,
		},
	})
	if nil != err {
		client.logger.WithError(err).WithField("config_map_name", configMap.Name).Error("Failed to read config_map list")
		return false, fmt.Errorf("Failed to read config_map list, %v", err)
	}
	for _, item := range collection.Data {
		if item.Name == configMap.Name {
			client.logger.WithField("config_map_name", configMap.Name).WithField("namespace", configMap.Namespace).Debug("ConfigMap found")
			client.configMapCache[configMap.Namespace+"/"+configMap.Name] = item
			return true, nil
		}
	}
	client.logger.WithError(err).WithField("config_map_name", configMap.Name).Debug("ConfigMap not found")
	return false, nil
}

func (client *rancherClient) UpgradeConfigMap(configMap projectModel.ConfigMap) error {
	namespaceID, err := client.getNamespaceID(configMap.Namespace)
	var existingConfigMap projectClient.ConfigMap
	if item, exists := client.configMapCache[configMap.Namespace+"/"+configMap.Name]; exists {
		client.logger.WithField("config_map_name", configMap.Name).WithField("namespace", configMap.Namespace).Trace("Use Cache")
		existingConfigMap = item
	} else {
		collection, err := client.projectClient.ConfigMap.List(&types.ListOpts{
			Filters: map[string]interface{}{
				"system":      "false",
				"name":        configMap.Name,
				"namespaceId": namespaceID,
			},
		})
		if nil != err {
			client.logger.WithError(err).WithField("config_map_name", configMap.Name).WithField("namespace", configMap.Namespace).Error("Failed to read configMap list")
			return fmt.Errorf("Failed to read configMap list, %v", err)
		}

		if len(collection.Data) == 0 {
			return fmt.Errorf("ConfigMap %v not found", configMap.Name)
		}
		existingConfigMap = collection.Data[0]
	}
	if isConfigMapUnchanged(existingConfigMap, configMap) {
		client.logger.WithField("config_map_name", configMap.Name).WithField("namespace", configMap.Namespace).Debug("Skip upgrade ConfigMap - no changes")
		return nil
	}
	client.logger.WithField("config_map_name", configMap.Name).WithField("namespace", configMap.Namespace).Info("Upgrade ConfigMap")
	existingConfigMap.Data = configMap.Data
	existingConfigMap.Labels["cattlectl.io/hash"] = hashOf(configMap)

	_, err = client.projectClient.ConfigMap.Replace(&existingConfigMap)
	return err
}

func (client *rancherClient) CreateConfigMap(configMap projectModel.ConfigMap) error {
	namespaceID, err := client.getNamespaceID(configMap.Namespace)
	if err != nil {
		return fmt.Errorf("Failed to read namespace ID, %v", err)
	}
	client.logger.WithField("config_map_name", configMap.Name).Info("Create new ConfigMap")
	labels := make(map[string]string)
	labels["cattlectl.io/hash"] = hashOf(configMap)
	newConfigMap := &projectClient.ConfigMap{
		Name:        configMap.Name,
		Labels:      labels,
		Data:        configMap.Data,
		NamespaceId: namespaceID,
		ProjectID:   client.projectId,
	}

	_, err = client.projectClient.ConfigMap.Create(newConfigMap)
	return err
}

func isConfigMapUnchanged(existingConfigMap projectClient.ConfigMap, configMap projectModel.ConfigMap) bool {
	hash, hashExists := existingConfigMap.Labels["cattlectl.io/hash"]
	if !hashExists {
		return false
	}
	return hash == hashOf(configMap)
}
