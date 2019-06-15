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
	backendProjectClient "github.com/rancher/types/client/project/v3"
	"github.com/sirupsen/logrus"
)

func newConfigMapClientWithData(
	configMap projectModel.ConfigMap,
	namespace string,
	project ProjectClient,
	logger *logrus.Entry,
) (ConfigMapClient, error) {
	result, err := newConfigMapClient(
		configMap.Name,
		namespace,
		project,
		logger,
	)
	if err != nil {
		return nil, err
	}
	err = result.SetData(configMap)
	return result, err
}

func newConfigMapClient(
	name, namespace string,
	project ProjectClient,
	logger *logrus.Entry,
) (ConfigMapClient, error) {
	return &configMapClient{
		namespacedResourceClient: namespacedResourceClient{
			resourceClient: resourceClient{
				name:   name,
				logger: logger.WithField("configMap_name", name).WithField("namespace", namespace),
			},
			namespace: namespace,
			project:   project,
		},
	}, nil
}

type configMapClient struct {
	namespacedResourceClient
	configMap projectModel.ConfigMap
}

func (client *configMapClient) Exists() (bool, error) {
	backendClient, err := client.project.backendProjectClient()
	if err != nil {
		return false, err
	}
	namespaceID, err := client.NamespaceID()
	if err != nil {
		return false, err
	}
	collection, err := backendClient.ConfigMap.List(&types.ListOpts{
		Filters: map[string]interface{}{
			"name":        client.name,
			"namespaceId": namespaceID,
		},
	})
	if nil != err {
		client.logger.WithError(err).Error("Failed to read configMap list")
		return false, fmt.Errorf("Failed to read configMap list, %v", err)
	}
	for _, item := range collection.Data {
		if item.Name == client.name && item.NamespaceId == client.namespaceID {
			return true, nil
		}
	}
	client.logger.Debug("ConfigMap not found")
	return false, nil
}

func (client *configMapClient) Create() error {
	backendClient, err := client.project.backendProjectClient()
	if err != nil {
		return err
	}
	namespaceID, err := client.NamespaceID()
	if err != nil {
		return err
	}
	projectID, err := client.project.ID()
	if err != nil {
		return fmt.Errorf("Failed to read namespace ID, %v", err)
	}
	client.logger.Info("Create new ConfigMap")
	labels := make(map[string]string)
	labels["cattlectl.io/hash"] = hashOf(client.configMap)
	newConfigMap := &backendProjectClient.ConfigMap{
		Name:        client.configMap.Name,
		Labels:      labels,
		Data:        client.configMap.Data,
		NamespaceId: namespaceID,
		ProjectID:   projectID,
	}

	_, err = backendClient.ConfigMap.Create(newConfigMap)
	return err
}

func (client *configMapClient) Upgrade() error {
	backendClient, err := client.project.backendProjectClient()
	if err != nil {
		return err
	}
	namespaceID, err := client.NamespaceID()
	if err != nil {
		return err
	}
	collection, err := backendClient.ConfigMap.List(&types.ListOpts{
		Filters: map[string]interface{}{
			"name":        client.name,
			"namespaceId": namespaceID,
		},
	})
	if nil != err {
		client.logger.WithError(err).Error("Failed to read configMap list")
		return fmt.Errorf("Failed to read configMap list, %v", err)
	}

	if len(collection.Data) == 0 {
		return fmt.Errorf("ConfigMap %v not found", client.name)
	}
	existingConfigMap := collection.Data[0]
	if isConfigMapUnchanged(existingConfigMap, client.configMap) {
		client.logger.Debug("Skip upgrade configMap - no changes")
		return nil
	}
	client.logger.Info("Upgrade ConfigMap")
	existingConfigMap.Labels["cattlectl.io/hash"] = hashOf(client.configMap)
	existingConfigMap.Data = client.configMap.Data

	_, err = backendClient.ConfigMap.Replace(&existingConfigMap)
	return err
}

func (client *configMapClient) Data() (projectModel.ConfigMap, error) {
	return client.configMap, nil
}

func (client *configMapClient) SetData(configMap projectModel.ConfigMap) error {
	client.name = configMap.Name
	client.configMap = configMap
	return nil
}

func isConfigMapUnchanged(existingConfigMap backendProjectClient.ConfigMap, configMap projectModel.ConfigMap) bool {
	hash, hashExists := existingConfigMap.Labels["cattlectl.io/hash"]
	if !hashExists {
		return false
	}
	return hash == hashOf(configMap)
}
