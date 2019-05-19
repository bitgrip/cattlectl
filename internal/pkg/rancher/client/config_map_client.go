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

	projectModel "github.com/bitgrip/cattlectl/internal/pkg/rancher/project/model"
	"github.com/rancher/norman/types"
	backendProjectClient "github.com/rancher/types/client/project/v3"
	"github.com/sirupsen/logrus"
)

func newConfigMapClientWithData(
	configMap projectModel.ConfigMap,
	namespace string,
	project ProjectClient,
	backendProjectClient *backendProjectClient.Client,
	logger *logrus.Entry,
) (ConfigMapClient, error) {
	result, err := newConfigMapClient(
		configMap.Name,
		namespace,
		project,
		backendProjectClient,
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
	backendProjectClient *backendProjectClient.Client,
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
		backendProjectClient: backendProjectClient,
	}, nil
}

type configMapClient struct {
	namespacedResourceClient
	configMap            projectModel.ConfigMap
	backendProjectClient *backendProjectClient.Client
}

func (client *configMapClient) init() error {
	namespaceID, err := client.NamespaceID()
	if namespaceID == "" && err == nil {
		return fmt.Errorf("Can not find namespace")
	}
	return err
}

func (client *configMapClient) Exists() (bool, error) {
	if err := client.init(); err != nil {
		return false, err
	}
	collection, err := client.backendProjectClient.ConfigMap.List(&types.ListOpts{
		Filters: map[string]interface{}{
			"name":        client.name,
			"namespaceId": client.namespaceID,
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
	if err := client.init(); err != nil {
		return err
	}
	namespaceID, err := client.NamespaceID()
	if err != nil {
		return fmt.Errorf("Failed to read namespace ID, %v", err)
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

	_, err = client.backendProjectClient.ConfigMap.Create(newConfigMap)
	return err
}

func (client *configMapClient) Upgrade() error {
	return fmt.Errorf("upgrade statefulset not supported")
}

func (client *configMapClient) Data() (projectModel.ConfigMap, error) {
	return client.configMap, nil
}

func (client *configMapClient) SetData(configMap projectModel.ConfigMap) error {
	client.name = configMap.Name
	client.configMap = configMap
	return nil
}
