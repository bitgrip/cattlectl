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

func (client *rancherClient) HasDockerCredential(dockerCredential projectModel.DockerCredential) (bool, error) {
	if _, exists := client.dockerCredentialCache[dockerCredential.Name]; exists {
		return true, nil
	}
	collection, err := client.projectClient.DockerCredential.List(&types.ListOpts{
		Filters: map[string]interface{}{
			"system": "false",
			"name":   dockerCredential.Name,
		},
	})
	if nil != err {
		client.logger.WithError(err).WithField("docker_credential_name", dockerCredential.Name).Error("Failed to read DockerCredential list")
		return false, fmt.Errorf("Failed to read DockerCredential list, %v", err)
	}
	for _, item := range collection.Data {
		if item.Name == dockerCredential.Name {
			client.logger.WithField("docker_credential_name", dockerCredential.Name).Debug("DockerCredential found")
			client.dockerCredentialCache[dockerCredential.Name] = item
			return true, nil
		}
	}
	client.logger.WithError(err).WithField("docker_credential_name", dockerCredential.Name).WithField("namespace", dockerCredential.Namespace).Debug("DockerCredential not found")
	return false, nil
}

func (client *rancherClient) UpgradeDockerCredential(dockerCredential projectModel.DockerCredential) error {
	var existingDockerCredential projectClient.DockerCredential
	if item, exists := client.dockerCredentialCache[dockerCredential.Name]; exists {
		client.logger.WithField("docker_credential_name", dockerCredential.Name).Trace("Use Cache")
		existingDockerCredential = item
	} else {
		collection, err := client.projectClient.DockerCredential.List(&types.ListOpts{
			Filters: map[string]interface{}{
				"system": "false",
				"name":   dockerCredential.Name,
			},
		})
		if nil != err {
			client.logger.WithError(err).WithField("docker_credential_name", dockerCredential.Name).Error("Failed to read DockerCredential list")
			return fmt.Errorf("Failed to read DockerCredential list, %v", err)
		}

		if len(collection.Data) == 0 {
			return fmt.Errorf("DockerCredential %v not found", dockerCredential.Name)
		}
		existingDockerCredential = collection.Data[0]
		client.dockerCredentialCache[dockerCredential.Name] = existingDockerCredential
	}
	if isDockerCredentialUnchanged(existingDockerCredential, dockerCredential) {
		client.logger.WithField("docker_credential_name", dockerCredential.Name).Debug("Skip upgrade DockerCredential - no changes")
		return nil
	}
	client.logger.WithField("docker_credential_name", dockerCredential.Name).Info("Upgrade DockerCredential")
	registries := make(map[string]projectClient.RegistryCredential)
	for _, registry := range dockerCredential.Registries {
		registries[registry.Name] = projectClient.RegistryCredential{
			Username: registry.Username,
			Password: registry.Password,
		}
	}
	existingDockerCredential.Registries = registries
	existingDockerCredential.Labels["cattlectl.io/hash"] = hashOf(dockerCredential)

	_, err := client.projectClient.DockerCredential.Replace(&existingDockerCredential)
	return err
}

func (client *rancherClient) CreateDockerCredential(dockerCredential projectModel.DockerCredential) error {
	client.logger.WithField("docker_credential_name", dockerCredential.Name).Info("Create new DockerCredential")

	registries := make(map[string]projectClient.RegistryCredential)
	for _, registry := range dockerCredential.Registries {
		registries[registry.Name] = projectClient.RegistryCredential{
			Username: registry.Username,
			Password: registry.Password,
		}
	}
	labels := make(map[string]string)
	labels["cattlectl.io/hash"] = hashOf(dockerCredential)
	newDockerCredential := &projectClient.DockerCredential{
		Name:        dockerCredential.Name,
		Labels:      labels,
		Registries:  registries,
		NamespaceId: dockerCredential.Namespace,
		ProjectID:   client.projectId,
	}

	_, err := client.projectClient.DockerCredential.Create(newDockerCredential)
	return err
}

func (client *rancherClient) HasNamespacedDockerCredential(dockerCredential projectModel.DockerCredential) (bool, error) {
	if _, exists := client.namespacedDockerCredentialCache[dockerCredential.Name]; exists {
		return true, nil
	}
	namespaceID, err := client.getNamespaceID(dockerCredential.Namespace)
	if err != nil {
		return false, fmt.Errorf("Failed to read config_map list, %v", err)
	}
	collection, err := client.projectClient.NamespacedDockerCredential.List(&types.ListOpts{
		Filters: map[string]interface{}{
			"system":      "false",
			"name":        dockerCredential.Name,
			"namespaceId": namespaceID,
		},
	})
	if nil != err {
		client.logger.WithError(err).WithField("docker_credential_name", dockerCredential.Name).Error("Failed to read DockerCredential list")
		return false, fmt.Errorf("Failed to read DockerCredential list, %v", err)
	}
	for _, item := range collection.Data {
		if item.Name == dockerCredential.Name {
			client.logger.WithField("docker_credential_name", dockerCredential.Name).WithField("namespace", dockerCredential.Namespace).Debug("DockerCredential found")
			client.namespacedDockerCredentialCache[dockerCredential.Name] = item
			return true, nil
		}
	}
	client.logger.WithError(err).WithField("docker_credential_name", dockerCredential.Name).WithField("namespace", dockerCredential.Namespace).Debug("DockerCredential not found")
	return false, nil
}

func (client *rancherClient) UpgradeNamespacedDockerCredential(dockerCredential projectModel.DockerCredential) error {
	var existingDockerCredential projectClient.NamespacedDockerCredential
	if item, exists := client.namespacedDockerCredentialCache[dockerCredential.Name]; exists {
		client.logger.WithField("docker_credential_name", dockerCredential.Name).WithField("namespace", dockerCredential.Namespace).Trace("Use Cache")
		existingDockerCredential = item
	} else {
		namespaceID, err := client.getNamespaceID(dockerCredential.Namespace)
		if err != nil {
			return fmt.Errorf("Failed to read config_map list, %v", err)
		}
		collection, err := client.projectClient.NamespacedDockerCredential.List(&types.ListOpts{
			Filters: map[string]interface{}{
				"system":      "false",
				"name":        dockerCredential.Name,
				"namespaceId": namespaceID,
			},
		})
		if nil != err {
			client.logger.WithError(err).WithField("docker_credential_name", dockerCredential.Name).WithField("namespace", dockerCredential.Namespace).Error("Failed to read DockerCredential list")
			return fmt.Errorf("Failed to read DockerCredential list, %v", err)
		}

		if len(collection.Data) == 0 {
			return fmt.Errorf("DockerCredential %v not found", dockerCredential.Name)
		}
		existingDockerCredential = collection.Data[0]
	}
	if isNamespacedDockerCredentialUnchanged(existingDockerCredential, dockerCredential) {
		client.logger.WithField("docker_credential_name", dockerCredential.Name).WithField("namespace", dockerCredential.Namespace).Debug("Skip upgrade DockerCredential - no changes")
		return nil
	}
	client.logger.WithField("docker_credential_name", dockerCredential.Name).WithField("namespace", dockerCredential.Namespace).Info("Upgrade DockerCredential")
	registries := make(map[string]projectClient.RegistryCredential)
	for _, registry := range dockerCredential.Registries {
		registries[registry.Name] = projectClient.RegistryCredential{
			Username: registry.Username,
			Password: registry.Password,
		}
	}
	existingDockerCredential.Registries = registries
	existingDockerCredential.Labels["cattlectl.io/hash"] = hashOf(dockerCredential)

	_, err := client.projectClient.NamespacedDockerCredential.Replace(&existingDockerCredential)
	return err
}

func (client *rancherClient) CreateNamespacedDockerCredential(dockerCredential projectModel.DockerCredential) error {
	client.logger.WithField("docker_credential_name", dockerCredential.Name).WithField("namespace", dockerCredential.Namespace).Info("Create new DockerCredential")

	registries := make(map[string]projectClient.RegistryCredential)
	for _, registry := range dockerCredential.Registries {
		registries[registry.Name] = projectClient.RegistryCredential{
			Username: registry.Username,
			Password: registry.Password,
		}
	}
	labels := make(map[string]string)
	labels["cattlectl.io/hash"] = hashOf(dockerCredential)
	namespaceID, err := client.getNamespaceID(dockerCredential.Namespace)
	if err != nil {
		return fmt.Errorf("Failed to read config_map list, %v", err)
	}
	newDockerCredential := &projectClient.NamespacedDockerCredential{
		Name:        dockerCredential.Name,
		Labels:      labels,
		Registries:  registries,
		NamespaceId: namespaceID,
		ProjectID:   client.projectId,
	}

	_, err = client.projectClient.NamespacedDockerCredential.Create(newDockerCredential)
	return err
}

func isDockerCredentialUnchanged(existingDockerCredential projectClient.DockerCredential, dockerCredential projectModel.DockerCredential) bool {
	hash, hashExists := existingDockerCredential.Labels["cattlectl.io/hash"]
	if !hashExists {
		return false
	}
	return hash == hashOf(dockerCredential)
}

func isNamespacedDockerCredentialUnchanged(existingDockerCredential projectClient.NamespacedDockerCredential, dockerCredential projectModel.DockerCredential) bool {
	hash, hashExists := existingDockerCredential.Labels["cattlectl.io/hash"]
	if !hashExists {
		return false
	}
	return hash == hashOf(dockerCredential)
}
