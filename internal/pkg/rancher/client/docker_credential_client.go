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

func newDockerCredentialClientWithData(
	dockerCredential projectModel.DockerCredential,
	namespace string,
	project ProjectClient,
	backendProjectClient *backendProjectClient.Client,
	logger *logrus.Entry,
) (DockerCredentialClient, error) {
	result, err := newDockerCredentialClient(
		dockerCredential.Name,
		namespace,
		project,
		backendProjectClient,
		logger,
	)
	if err != nil {
		return nil, err
	}
	err = result.SetData(dockerCredential)
	return result, err
}

func newDockerCredentialClient(
	name, namespace string,
	project ProjectClient,
	backendProjectClient *backendProjectClient.Client,
	logger *logrus.Entry,
) (DockerCredentialClient, error) {
	return &dockerCredentialClient{
		namespacedResourceClient: namespacedResourceClient{
			resourceClient: resourceClient{
				name:   name,
				logger: logger.WithField("dockerCredential_name", name).WithField("namespace", namespace),
			},
			namespace: namespace,
			project:   project,
		},
		backendProjectClient: backendProjectClient,
	}, nil
}

type dockerCredentialClient struct {
	namespacedResourceClient
	dockerCredential     projectModel.DockerCredential
	backendProjectClient *backendProjectClient.Client
}

func (client *dockerCredentialClient) init() error {
	namespaceID, err := client.NamespaceID()
	if namespaceID == "" && err == nil {
		return fmt.Errorf("Can not find namespace")
	}
	return err
}

func (client *dockerCredentialClient) Exists() (bool, error) {
	if err := client.init(); err != nil {
		return false, err
	}
	collection, err := client.backendProjectClient.DockerCredential.List(&types.ListOpts{
		Filters: map[string]interface{}{
			"name":        client.name,
			"namespaceId": client.namespaceID,
		},
	})
	if nil != err {
		client.logger.WithError(err).Error("Failed to read dockerCredential list")
		return false, fmt.Errorf("Failed to read dockerCredential list, %v", err)
	}
	for _, item := range collection.Data {
		if item.Name == client.name && item.NamespaceId == client.namespaceID {
			return true, nil
		}
	}
	client.logger.Debug("DockerCredential not found")
	return false, nil
}

func (client *dockerCredentialClient) Create() error {
	if err := client.init(); err != nil {
		return err
	}
	client.logger.Info("Create new DockerCredential")

	registries := make(map[string]backendProjectClient.RegistryCredential)
	for _, registry := range client.dockerCredential.Registries {
		registries[registry.Name] = backendProjectClient.RegistryCredential{
			Username: registry.Username,
			Password: registry.Password,
		}
	}
	labels := make(map[string]string)
	labels["cattlectl.io/hash"] = hashOf(client.dockerCredential)
	projectID, err := client.project.ID()
	if err != nil {
		return err
	}
	newDockerCredential := &backendProjectClient.DockerCredential{
		Name:        client.dockerCredential.Name,
		Labels:      labels,
		Registries:  registries,
		NamespaceId: client.dockerCredential.Namespace,
		ProjectID:   projectID,
	}

	_, err = client.backendProjectClient.DockerCredential.Create(newDockerCredential)
	return err
}

func (client *dockerCredentialClient) Upgrade() error {
	return fmt.Errorf("upgrade statefulset not supported")
}

func (client *dockerCredentialClient) Data() (projectModel.DockerCredential, error) {
	return client.dockerCredential, nil
}

func (client *dockerCredentialClient) SetData(dockerCredential projectModel.DockerCredential) error {
	client.name = dockerCredential.Name
	client.dockerCredential = dockerCredential
	return nil
}
