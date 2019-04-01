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
	collection, err := client.projectClient.DockerCredential.List(&types.ListOpts{
		Filters: map[string]interface{}{
			"system": "false",
			"name":   dockerCredential.Name,
		},
	})
	if nil != err {
		client.logger.WithError(err).WithField("docker_credential_name", dockerCredential.Name).Error("Failed to read DockerCredential list")
		return false, fmt.Errorf("Failed to read docker_credential list, %v", err)
	}
	for _, item := range collection.Data {
		if item.Name == dockerCredential.Name {
			return true, nil
		}
	}
	client.logger.WithError(err).WithField("docker_credential_name", dockerCredential.Name).Debug("DockerCredential not found")
	return false, nil
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

	newDockerCredential := &projectClient.DockerCredential{
		Name:        dockerCredential.Name,
		Registries:  registries,
		NamespaceId: dockerCredential.Namespace,
		ProjectID:   client.projectId,
	}

	_, err := client.projectClient.DockerCredential.Create(newDockerCredential)
	return err
}
