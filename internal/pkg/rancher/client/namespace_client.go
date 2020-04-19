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
	backendClusterClient "github.com/rancher/types/client/cluster/v3"
	"github.com/sirupsen/logrus"
)

func newNamespaceClientWithData(
	namespace projectModel.Namespace,
	projectClient ProjectClient,
	clusterClient ClusterClient,
	logger *logrus.Entry,
) (NamespaceClient, error) {
	result, err := newNamespaceClient(
		namespace.Name,
		projectClient,
		clusterClient,
		logger,
	)
	if err != nil {
		return nil, err
	}
	err = result.SetData(namespace)
	return result, err
}

func newNamespaceClient(
	name string,
	projectClient ProjectClient,
	clusterClient ClusterClient,
	logger *logrus.Entry,
) (NamespaceClient, error) {
	return &namespaceClient{
		resourceClient: resourceClient{
			name:   name,
			logger: logger.WithField("namespace_name", name),
		},
		projectClient: projectClient,
		clusterClient: clusterClient,
	}, nil
}

type namespaceClient struct {
	resourceClient
	namespace     projectModel.Namespace
	projectClient ProjectClient
	clusterClient ClusterClient
}

func (client *namespaceClient) ID() (string, error) {
	if client.id != "" {
		return client.id, nil
	}
	backendClient, err := client.clusterClient.backendClusterClient()
	if err != nil {
		return "", err
	}
	collection, err := backendClient.Namespace.List(&types.ListOpts{
		Filters: map[string]interface{}{
			"name": client.name,
		},
	})
	if err != nil {
		client.logger.WithError(err).Error("Failed to read namespace list")
		return "", fmt.Errorf("Failed to read namespace list, %v", err)
	}
	if len(collection.Data) < 1 {
		client.logger.Debug("Unknown Namespace")
		return "", nil
	}
	client.id = collection.Data[0].ID
	return client.id, nil
}

func (client *namespaceClient) Exists() (bool, error) {
	backendClient, err := client.clusterClient.backendClusterClient()
	if err != nil {
		return false, err
	}
	collection, err := backendClient.Namespace.List(&types.ListOpts{
		Filters: map[string]interface{}{
			"name": client.name,
		},
	})
	if err != nil {
		client.logger.WithError(err).Error("Failed to read namespace list")
		return false, fmt.Errorf("Failed to read namespace list, %v", err)
	}
	for _, item := range collection.Data {
		if item.Name == client.name {
			return true, nil
		}
	}
	client.logger.Debug("Namespace not found")
	return false, nil
}

func (client *namespaceClient) Create() error {
	backendClient, err := client.clusterClient.backendClusterClient()
	if err != nil {
		return err
	}

	client.logger.Info("Create new namespace")
	newNamespace := &backendClusterClient.Namespace{
		Name: client.namespace.Name,
	}
	if hasProjct, _ := client.HasProject(); hasProjct {
		projectID, err := client.projectClient.ID()
		if err != nil {
			return err
		}
		newNamespace.ProjectID = projectID
	}

	_, err = backendClient.Namespace.Create(newNamespace)
	return err
}

func (client *namespaceClient) Upgrade() error {
	client.logger.Debug("Skip change existing namespace")
	return nil
}

func (client *namespaceClient) Delete() (err error) {
	backendClient, err := client.clusterClient.backendClusterClient()
	if err != nil {
		return err
	}
	existingNamespace, err := client.loadExistingNamespace()
	if err != nil {
		return err
	}
	return backendClient.Namespace.Delete(existingNamespace)
}

func (client *namespaceClient) HasProject() (bool, error) {
	return client.projectClient != nil, nil
}

func (client *namespaceClient) Data() (projectModel.Namespace, error) {
	return client.namespace, nil
}

func (client *namespaceClient) SetData(namespace projectModel.Namespace) error {
	client.name = namespace.Name
	client.namespace = namespace
	return nil
}

func (client *namespaceClient) loadExistingNamespace() (existingNamespace *backendClusterClient.Namespace, err error) {
	backendClient, err := client.clusterClient.backendClusterClient()
	if err != nil {
		return
	}
	collection, err := backendClient.Namespace.List(&types.ListOpts{
		Filters: map[string]interface{}{
			"name": client.name,
		},
	})
	if err != nil {
		client.logger.WithError(err).Error("Failed to read namespace list")
		return nil, fmt.Errorf("Failed to read namespace list, %v", err)
	}

	if len(collection.Data) == 0 {
		return
	}

	existingNamespace = &collection.Data[0]
	return
}
