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

func newDockerCredentialClientWithData(
	dockerCredential projectModel.DockerCredential,
	namespace string,
	project ProjectClient,
	logger *logrus.Entry,
) (DockerCredentialClient, error) {
	result, err := newDockerCredentialClient(
		dockerCredential.Name,
		namespace,
		project,
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
	}, nil
}

type dockerCredentialClient struct {
	namespacedResourceClient
	dockerCredential projectModel.DockerCredential
}

func (client *dockerCredentialClient) Exists() (bool, error) {
	if client.namespace != "" {
		return client.existsInNamespace()
	}
	return client.existsInProject()
}

func (client *dockerCredentialClient) existsInProject() (bool, error) {
	backendClient, err := client.project.backendProjectClient()
	if err != nil {
		return false, err
	}
	collection, err := backendClient.DockerCredential.List(&types.ListOpts{
		Filters: map[string]interface{}{
			"name": client.name,
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

func (client *dockerCredentialClient) existsInNamespace() (bool, error) {
	backendClient, err := client.project.backendProjectClient()
	if err != nil {
		return false, err
	}
	namespaceID, err := client.NamespaceID()
	if err != nil {
		return false, err
	}
	collection, err := backendClient.NamespacedDockerCredential.List(&types.ListOpts{
		Filters: map[string]interface{}{
			"name":        client.name,
			"namespaceId": namespaceID,
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
	if client.namespace != "" {
		return client.createInNamespace(registries, labels, projectID)
	}
	return client.createInProject(registries, labels, projectID)
}

func (client *dockerCredentialClient) createInProject(registries map[string]backendProjectClient.RegistryCredential, labels map[string]string, projectID string) error {
	backendClient, err := client.project.backendProjectClient()
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

	_, err = backendClient.DockerCredential.Create(newDockerCredential)
	return err
}

func (client *dockerCredentialClient) createInNamespace(registries map[string]backendProjectClient.RegistryCredential, labels map[string]string, projectID string) error {
	backendClient, err := client.project.backendProjectClient()
	if err != nil {
		return err
	}
	namespaceID, err := client.NamespaceID()
	if err != nil {
		return err
	}
	newDockerCredential := &backendProjectClient.NamespacedDockerCredential{
		Name:        client.dockerCredential.Name,
		Labels:      labels,
		Registries:  registries,
		NamespaceId: namespaceID,
		ProjectID:   projectID,
	}

	_, err = backendClient.NamespacedDockerCredential.Create(newDockerCredential)
	return err
}

func (client *dockerCredentialClient) Upgrade() error {
	if client.namespace != "" {
		return client.upgradeInNamespace()
	}
	return client.upgradeInProject()
}

func (client *dockerCredentialClient) upgradeInProject() error {
	backendClient, err := client.project.backendProjectClient()
	if err != nil {
		return err
	}
	collection, err := backendClient.DockerCredential.List(&types.ListOpts{
		Filters: map[string]interface{}{
			"name": client.name,
		},
	})
	if nil != err {
		client.logger.WithError(err).Error("Failed to read dockerCredential list")
		return fmt.Errorf("Failed to read dockerCredential list, %v", err)
	}

	if len(collection.Data) == 0 {
		return fmt.Errorf("DockerCredential %v not found", client.name)
	}
	existingDockerCredential := collection.Data[0]
	if isProjectDockerCredentialUnchanged(existingDockerCredential, client.dockerCredential) {
		client.logger.Debug("Skip upgrade DockerCredential - no changes")
		return nil
	}
	client.logger.Info("Upgrade DockerCredential")
	registries := make(map[string]backendProjectClient.RegistryCredential)
	for _, registry := range client.dockerCredential.Registries {
		registries[registry.Name] = backendProjectClient.RegistryCredential{
			Username: registry.Username,
			Password: registry.Password,
		}
	}
	existingDockerCredential.Registries = registries
	existingDockerCredential.Labels["cattlectl.io/hash"] = hashOf(client.dockerCredential)

	_, err = backendClient.DockerCredential.Replace(&existingDockerCredential)
	return err
}

func (client *dockerCredentialClient) upgradeInNamespace() error {
	backendClient, err := client.project.backendProjectClient()
	if err != nil {
		return err
	}
	namespaceID, err := client.NamespaceID()
	if err != nil {
		return err
	}
	collection, err := backendClient.NamespacedDockerCredential.List(&types.ListOpts{
		Filters: map[string]interface{}{
			"name":        client.name,
			"namespaceId": namespaceID,
		},
	})
	if nil != err {
		client.logger.WithError(err).Error("Failed to read dockerCredential list")
		return fmt.Errorf("Failed to read dockerCredential list, %v", err)
	}

	if len(collection.Data) == 0 {
		return fmt.Errorf("DockerCredential %v not found", client.name)
	}
	existingDockerCredential := collection.Data[0]
	if isNamespacedDockerCredentialUnchanged(existingDockerCredential, client.dockerCredential) {
		client.logger.Debug("Skip upgrade DockerCredential - no changes")
		return nil
	}
	client.logger.Info("Upgrade DockerCredential")
	registries := make(map[string]backendProjectClient.RegistryCredential)
	for _, registry := range client.dockerCredential.Registries {
		registries[registry.Name] = backendProjectClient.RegistryCredential{
			Username: registry.Username,
			Password: registry.Password,
		}
	}
	existingDockerCredential.Registries = registries
	existingDockerCredential.Labels["cattlectl.io/hash"] = hashOf(client.dockerCredential)

	_, err = backendClient.NamespacedDockerCredential.Replace(&existingDockerCredential)
	return err
}

func (client *dockerCredentialClient) Data() (projectModel.DockerCredential, error) {
	return client.dockerCredential, nil
}

func (client *dockerCredentialClient) SetData(dockerCredential projectModel.DockerCredential) error {
	client.name = dockerCredential.Name
	client.dockerCredential = dockerCredential
	return nil
}

func isProjectDockerCredentialUnchanged(existingDockerCredential backendProjectClient.DockerCredential, dockerCredential projectModel.DockerCredential) bool {
	hash, hashExists := existingDockerCredential.Labels["cattlectl.io/hash"]
	if !hashExists {
		return false
	}
	return hash == hashOf(dockerCredential)
}

func isNamespacedDockerCredentialUnchanged(existingDockerCredential backendProjectClient.NamespacedDockerCredential, dockerCredential projectModel.DockerCredential) bool {
	hash, hashExists := existingDockerCredential.Labels["cattlectl.io/hash"]
	if !hashExists {
		return false
	}
	return hash == hashOf(dockerCredential)
}
