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
	backendClient "github.com/rancher/types/client/project/v3"
	"github.com/sirupsen/logrus"
)

func newSecretpClientWithData(
	configMap projectModel.ConfigMap,
	namespace string,
	project ProjectClient,
	backendClient *backendClient.Client,
	logger *logrus.Entry,
) (ConfigMapClient, error) {
	result, err := newConfigMapClient(
		configMap.Name,
		namespace,
		project,
		backendClient,
		logger,
	)
	if err != nil {
		return nil, err
	}
	err = result.SetData(configMap)
	return result, err
}

func newSecretClient(
	name, namespace string,
	project ProjectClient,
	backendClient *backendClient.Client,
	logger *logrus.Entry,
) (ConfigMapClient, error) {
	return &secretClient{
		namespacedResourceClient: namespacedResourceClient{
			resourceClient: resourceClient{
				name:   name,
				logger: logger.WithField("configMap_name", name).WithField("namespace", namespace),
			},
			namespace: namespace,
			project:   project,
		},
		backendClient: backendClient,
	}, nil
}

type secretClient struct {
	namespacedResourceClient
	configMap     projectModel.ConfigMap
	backendClient *backendClient.Client
}

func (client *secretClient) init() error {
	namespaceID, err := client.NamespaceID()
	if namespaceID == "" && err == nil {
		return fmt.Errorf("Can not find namespace")
	}
	return err
}

func (client *secretClient) Exists() (bool, error) {
	if err := client.init(); err != nil {
		return false, err
	}
	collection, err := client.backendClient.Secret.List(&types.ListOpts{
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
	client.logger.Debug("Secret not found")
	return false, nil
}

func (client *secretClient) Create() error {
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
	client.logger.Info("Create new Secret")
	labels := make(map[string]string)
	labels["cattlectl.io/hash"] = hashOf(client.configMap)
	newSecret := &backendClient.Secret{
		Name:        client.configMap.Name,
		Labels:      labels,
		Data:        client.configMap.Data,
		NamespaceId: namespaceID,
		ProjectID:   projectID,
	}

	_, err = client.backendClient.Secret.Create(newSecret)
	return err
}

func (client *secretClient) Upgrade() error {
	return fmt.Errorf("upgrade statefulset not supported")
}

func (client *secretClient) Data() (projectModel.ConfigMap, error) {
	return client.configMap, nil
}

func (client *secretClient) SetData(configMap projectModel.ConfigMap) error {
	client.name = configMap.Name
	client.configMap = configMap
	return nil
}
