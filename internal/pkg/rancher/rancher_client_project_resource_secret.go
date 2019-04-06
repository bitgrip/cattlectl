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
	"reflect"

	projectModel "github.com/bitgrip/cattlectl/internal/pkg/rancher/project/model"
	"github.com/rancher/norman/types"
	projectClient "github.com/rancher/types/client/project/v3"
)

func (client *rancherClient) HasSecret(secret projectModel.ConfigMap) (bool, error) {
	if _, exists := client.secretCache[secret.Name]; exists {
		return true, nil
	}
	collection, err := client.projectClient.Secret.List(&types.ListOpts{
		Filters: map[string]interface{}{
			"system": "false",
			"name":   secret.Name,
		},
	})
	if nil != err {
		client.logger.WithError(err).WithField("secret_name", secret.Name).Error("Failed to read secret list")
		return false, fmt.Errorf("Failed to read secret list, %v", err)
	}
	for _, item := range collection.Data {
		if item.Name == secret.Name {
			client.logger.WithField("secret_name", secret.Name).Debug("Secret found")
			client.secretCache[secret.Name] = item
			return true, nil
		}
	}
	client.logger.WithError(err).WithField("secret_name", secret.Name).WithField("namespace", secret.Namespace).Debug("Secret not found")
	return false, nil
}

func (client *rancherClient) UpgradeSecret(secret projectModel.ConfigMap) error {
	var existingSecret projectClient.Secret
	if item, exists := client.secretCache[secret.Name]; exists {
		client.logger.WithField("secret_name", secret.Name).Trace("Use Cache")
		existingSecret = item
	} else {
		collection, err := client.projectClient.Secret.List(&types.ListOpts{
			Filters: map[string]interface{}{
				"system": "false",
				"name":   secret.Name,
			},
		})
		if nil != err {
			client.logger.WithError(err).WithField("secret_name", secret.Name).Error("Failed to read secret list")
			return fmt.Errorf("Failed to read secret list, %v", err)
		}

		if len(collection.Data) == 0 {
			return fmt.Errorf("Secret %v not found", secret.Name)
		}
		existingSecret = collection.Data[0]
	}
	if reflect.DeepEqual(existingSecret.Data, secret.Data) {
		client.logger.WithField("secret_name", secret.Name).Debug("Skip upgrade secret - no changes")
		return nil
	}
	client.logger.WithField("secret_name", secret.Name).Info("Upgrade Secret")
	existingSecret.Data = secret.Data

	_, err := client.projectClient.Secret.Replace(&existingSecret)
	return err
}

func (client *rancherClient) CreateSecret(secret projectModel.ConfigMap) error {
	client.logger.WithField("secret_name", secret.Name).Info("Create new Secret")
	newSecret := &projectClient.Secret{
		Name:      secret.Name,
		Data:      secret.Data,
		ProjectID: client.projectId,
	}

	_, err := client.projectClient.Secret.Create(newSecret)
	return err
}

func (client *rancherClient) HasNamespacedSecret(secret projectModel.ConfigMap) (bool, error) {
	if _, exists := client.namespacedSecretCache[secret.Name]; exists {
		return true, nil
	}
	collection, err := client.projectClient.NamespacedSecret.List(&types.ListOpts{
		Filters: map[string]interface{}{
			"system":      "false",
			"name":        secret.Name,
			"namespaceId": secret.Namespace,
		},
	})
	if nil != err {
		client.logger.WithError(err).WithField("secret_name", secret.Name).Error("Failed to read secret list")
		return false, fmt.Errorf("Failed to read secret list, %v", err)
	}
	for _, item := range collection.Data {
		if item.Name == secret.Name {
			client.logger.WithField("secret_name", secret.Name).WithField("namespace", secret.Namespace).Debug("Secret found")
			client.namespacedSecretCache[secret.Name] = item
			return true, nil
		}
	}
	client.logger.WithError(err).WithField("secret_name", secret.Name).WithField("namespace", secret.Namespace).Debug("Secret not found")
	return false, nil
}

func (client *rancherClient) UpgradeNamespacedSecret(secret projectModel.ConfigMap) error {
	var existingSecret projectClient.NamespacedSecret
	if item, exists := client.namespacedSecretCache[secret.Name]; exists {
		client.logger.WithField("secret_name", secret.Name).WithField("namespace", secret.Namespace).Trace("Use Cache")
		existingSecret = item
	} else {
		collection, err := client.projectClient.NamespacedSecret.List(&types.ListOpts{
			Filters: map[string]interface{}{
				"system":      "false",
				"name":        secret.Name,
				"namespaceId": secret.Namespace,
			},
		})
		if nil != err {
			client.logger.WithError(err).WithField("secret_name", secret.Name).WithField("namespace", secret.Namespace).Error("Failed to read secret list")
			return fmt.Errorf("Failed to read secret list, %v", err)
		}

		if len(collection.Data) == 0 {
			return fmt.Errorf("Secret %v not found", secret.Name)
		}
		existingSecret = collection.Data[0]
	}
	if reflect.DeepEqual(existingSecret.Data, secret.Data) {
		client.logger.WithField("secret_name", secret.Name).WithField("namespace", secret.Namespace).Debug("Skip upgrade secret - no changes")
		return nil
	}
	client.logger.WithField("secret_name", secret.Name).WithField("namespace", secret.Namespace).Info("Upgrade Secret")
	existingSecret.Data = secret.Data

	_, err := client.projectClient.NamespacedSecret.Replace(&existingSecret)
	return err
}

func (client *rancherClient) CreateNamespacedSecret(secret projectModel.ConfigMap) error {
	client.logger.WithField("secret_name", secret.Name).WithField("namespace", secret.Namespace).Info("Create new Secret")
	newSecret := &projectClient.NamespacedSecret{
		Name:        secret.Name,
		Data:        secret.Data,
		NamespaceId: secret.Namespace,
		ProjectID:   client.projectId,
	}

	_, err := client.projectClient.NamespacedSecret.Create(newSecret)
	return err
}
