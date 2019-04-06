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
	clusterClient "github.com/rancher/types/client/cluster/v3"
)

func (client *rancherClient) HasNamespace(namespace projectModel.Namespace) (bool, error) {
	if _, exists := client.namespaceCache[namespace.Name]; exists {
		return true, nil
	}
	collection, err := client.clusterClient.Namespace.List(&types.ListOpts{
		Filters: map[string]interface{}{
			"system": "false",
			"name":   namespace.Name,
		},
	})
	if nil != err {
		client.logger.WithError(err).WithField("namespace_name", namespace.Name).Error("Failed to read namespace list")
		return false, fmt.Errorf("Failed to read namespace list, %v", err)
	}
	for _, item := range collection.Data {
		if item.Name == namespace.Name {
			client.namespaceCache[item.Name] = item
			return true, nil
		}
	}
	client.logger.WithError(err).WithField("namespace_name", namespace.Name).Debug("Namespace not found")
	return false, nil
}

func (client *rancherClient) CreateNamespace(namespace projectModel.Namespace) error {
	client.logger.WithField("namespace_name", namespace.Name).Info("Create new namespace")
	newNamespace := &clusterClient.Namespace{
		Name:      namespace.Name,
		ProjectID: client.projectId,
	}

	_, err := client.clusterClient.Namespace.Create(newNamespace)
	return err
}

func (client *rancherClient) getNamespaceID(namespace string) (string, error) {
	if namespace == "" {
		return "", fmt.Errorf("missing required namespace name")
	}
	hasNamespace, err := client.HasNamespace(projectModel.Namespace{Name: namespace})
	if err != nil {
		return "", err
	}
	if hasNamespace {
		return client.namespaceCache[namespace].ID, nil
	}
	return "", fmt.Errorf("Namespace not existing %v", namespace)
}
