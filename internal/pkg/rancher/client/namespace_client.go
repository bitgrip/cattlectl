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
	backendClusterClient "github.com/rancher/types/client/cluster/v3"
	"github.com/sirupsen/logrus"
)

func newNamespaceClientWithData(
	namespace projectModel.Namespace,
	projectClient ProjectClient,
	backendClusterClient *backendClusterClient.Client,
	logger *logrus.Entry,
) (NamespaceClient, error) {
	result, err := newNamespaceClient(
		namespace.Name,
		projectClient,
		backendClusterClient,
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
	backendClusterClient *backendClusterClient.Client,
	logger *logrus.Entry,
) (NamespaceClient, error) {
	return &namespaceClient{
		resourceClient: resourceClient{
			name:   name,
			logger: logger.WithField("namespace_name", name),
		},
		projectClient:        projectClient,
		backendClusterClient: backendClusterClient,
	}, nil
}

type namespaceClient struct {
	resourceClient
	namespace            projectModel.Namespace
	projectClient        ProjectClient
	backendClusterClient *backendClusterClient.Client
}

func (client *namespaceClient) init() error {
	return nil
}

func (client *namespaceClient) Exists() (bool, error) {
	if err := client.init(); err != nil {
		return false, err
	}
	collection, err := client.backendClusterClient.Namespace.List(&types.ListOpts{
		Filters: map[string]interface{}{
			"name": client.name,
		},
	})
	if nil != err {
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
	if err := client.init(); err != nil {
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

	_, err := client.backendClusterClient.Namespace.Create(newNamespace)
	return err
}

func (client *namespaceClient) Upgrade() error {
	return fmt.Errorf("upgrade statefulset not supported")
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
