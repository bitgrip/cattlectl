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
	"strings"

	projectModel "github.com/bitgrip/cattlectl/internal/pkg/rancher/project/model"
	"github.com/rancher/norman/types"
	projectClient "github.com/rancher/types/client/project/v3"
)

func (client *rancherClient) HasCertificate(certificate projectModel.Certificate) (bool, error) {
	collection, err := client.projectClient.Certificate.List(&types.ListOpts{
		Filters: map[string]interface{}{
			"system": "false",
			"name":   certificate.Name,
		},
	})
	if nil != err {
		client.logger.WithError(err).WithField("certificate_name", certificate.Name).Error("Failed to read certificate list")
		return false, fmt.Errorf("Failed to read certificate list, %v", err)
	}
	for _, item := range collection.Data {
		if item.Name == certificate.Name {
			return true, nil
		}
	}
	client.logger.WithError(err).WithField("certificate_name", certificate.Name).Debug("Certificate not found")
	return false, nil
}

func (client *rancherClient) UpgradeCertificate(certificate projectModel.Certificate) error {
	var existingCertificate projectClient.Certificate
	if item, exists := client.certificateCache[certificate.Name]; exists {
		client.logger.WithField("certificate_name", certificate.Name).Trace("Use Cache")
		existingCertificate = item
	} else {
		collection, err := client.projectClient.Certificate.List(&types.ListOpts{
			Filters: map[string]interface{}{
				"system": "false",
				"name":   certificate.Name,
			},
		})
		if nil != err {
			client.logger.WithError(err).WithField("certificate_name", certificate.Name).Error("Failed to read certificate list")
			return fmt.Errorf("Failed to read certificate list, %v", err)
		}

		if len(collection.Data) == 0 {
			return fmt.Errorf("Certificate %v not found", certificate.Name)
		}
		existingCertificate = collection.Data[0]
	}
	if strings.TrimSpace(existingCertificate.Certs) == strings.TrimSpace(certificate.Certs) {
		client.logger.WithField("certificate_name", certificate.Name).Debug("Skip upgrade certificate - no changes")
		return nil
	}
	client.logger.WithField("certificate_name", certificate.Name).Info("Upgrade Certificate")
	existingCertificate.Key = certificate.Key
	existingCertificate.Certs = certificate.Certs

	_, err := client.projectClient.Certificate.Replace(&existingCertificate)
	return err
}

func (client *rancherClient) CreateCertificate(certificate projectModel.Certificate) error {
	client.logger.WithField("certificate_name", certificate.Name).Info("Create new certificate")
	newCertificate := &projectClient.Certificate{
		Name:        certificate.Name,
		Key:         certificate.Key,
		Certs:       certificate.Certs,
		NamespaceId: certificate.Namespace,
		ProjectID:   client.projectId,
	}

	_, err := client.projectClient.Certificate.Create(newCertificate)
	return err
}

func (client *rancherClient) HasNamespacedCertificate(certificate projectModel.Certificate) (bool, error) {
	if _, exists := client.namespacedCertificateCache[certificate.Name]; exists {
		return true, nil
	}
	namespaceID, err := client.getNamespaceID(certificate.Namespace)
	if err != nil {
		return false, fmt.Errorf("Failed to read config_map list, %v", err)
	}
	collection, err := client.projectClient.NamespacedCertificate.List(&types.ListOpts{
		Filters: map[string]interface{}{
			"system":      "false",
			"name":        certificate.Name,
			"namespaceId": namespaceID,
		},
	})
	if nil != err {
		client.logger.WithError(err).WithField("certificate_name", certificate.Name).Error("Failed to read certificate list")
		return false, fmt.Errorf("Failed to read certificate list, %v", err)
	}
	for _, item := range collection.Data {
		if item.Name == certificate.Name {
			client.logger.WithField("certificate_name", certificate.Name).WithField("namespace", certificate.Namespace).Debug("Certificate found")
			client.namespacedCertificateCache[certificate.Name] = item
			return true, nil
		}
	}
	client.logger.WithError(err).WithField("certificate_name", certificate.Name).WithField("namespace", certificate.Namespace).Debug("Certificate not found")
	return false, nil
}

func (client *rancherClient) UpgradeNamespacedCertificate(certificate projectModel.Certificate) error {
	var existingCertificate projectClient.NamespacedCertificate
	if item, exists := client.namespacedCertificateCache[certificate.Name]; exists {
		client.logger.WithField("certificate_name", certificate.Name).WithField("namespace", certificate.Namespace).Trace("Use Cache")
		existingCertificate = item
	} else {
		namespaceID, err := client.getNamespaceID(certificate.Namespace)
		if err != nil {
			return fmt.Errorf("Failed to read config_map list, %v", err)
		}
		collection, err := client.projectClient.NamespacedCertificate.List(&types.ListOpts{
			Filters: map[string]interface{}{
				"system":      "false",
				"name":        certificate.Name,
				"namespaceId": namespaceID,
			},
		})
		if nil != err {
			client.logger.WithError(err).WithField("certificate_name", certificate.Name).WithField("namespace", certificate.Namespace).Error("Failed to read certificate list")
			return fmt.Errorf("Failed to read certificate list, %v", err)
		}

		if len(collection.Data) == 0 {
			return fmt.Errorf("Certificate %v not found", certificate.Name)
		}
		existingCertificate = collection.Data[0]
	}
	if strings.TrimSpace(existingCertificate.Certs) == strings.TrimSpace(certificate.Certs) {
		client.logger.WithField("certificate_name", certificate.Name).WithField("namespace", certificate.Namespace).Debug("Skip upgrade certificate - no changes")
		return nil
	}
	client.logger.WithField("certificate_name", certificate.Name).WithField("namespace", certificate.Namespace).Info("Upgrade Certificate")
	existingCertificate.Key = certificate.Key
	existingCertificate.Certs = certificate.Certs

	_, err := client.projectClient.NamespacedCertificate.Replace(&existingCertificate)
	return err
}

func (client *rancherClient) CreateNamespacedCertificate(certificate projectModel.Certificate) error {
	client.logger.WithField("certificate_name", certificate.Name).WithField("namespace", certificate.Namespace).Info("Create new Certificate")
	namespaceID, err := client.getNamespaceID(certificate.Namespace)
	if err != nil {
		return fmt.Errorf("Failed to read config_map list, %v", err)
	}
	newCertificate := &projectClient.NamespacedCertificate{
		Name:        certificate.Name,
		Key:         certificate.Key,
		Certs:       certificate.Certs,
		NamespaceId: namespaceID,
		ProjectID:   client.projectId,
	}

	_, err = client.projectClient.NamespacedCertificate.Create(newCertificate)
	return err
}
