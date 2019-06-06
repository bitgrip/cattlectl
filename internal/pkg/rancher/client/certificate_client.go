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
	"strings"

	projectModel "github.com/bitgrip/cattlectl/internal/pkg/rancher/project/model"
	"github.com/rancher/norman/types"
	backendProjectClient "github.com/rancher/types/client/project/v3"
	"github.com/sirupsen/logrus"
)

func newCertificateClientWithData(
	certificate projectModel.Certificate,
	namespace string,
	project ProjectClient,
	logger *logrus.Entry,
) (CertificateClient, error) {
	result, err := newCertificateClient(
		certificate.Name,
		namespace,
		project,
		logger,
	)
	if err != nil {
		return nil, err
	}
	err = result.SetData(certificate)
	return result, err
}

func newCertificateClient(
	name, namespace string,
	project ProjectClient,
	logger *logrus.Entry,
) (CertificateClient, error) {
	return &certificateClient{
		namespacedResourceClient: namespacedResourceClient{
			resourceClient: resourceClient{
				name:   name,
				logger: logger.WithField("certificate_name", name).WithField("namespace", namespace),
			},
			namespace: namespace,
			project:   project,
		},
	}, nil
}

type certificateClient struct {
	namespacedResourceClient
	certificate projectModel.Certificate
}

func (client *certificateClient) Exists() (bool, error) {
	if client.namespace != "" {
		return client.existsInNamespace()
	}
	return client.existsInProject()
}

func (client *certificateClient) existsInProject() (bool, error) {
	backendClient, err := client.project.backendProjectClient()
	if err != nil {
		return false, err
	}
	collection, err := backendClient.Certificate.List(&types.ListOpts{
		Filters: map[string]interface{}{
			"name": client.name,
		},
	})
	if nil != err {
		client.logger.WithError(err).Error("Failed to read certificate list")
		return false, fmt.Errorf("Failed to read certificate list, %v", err)
	}
	for _, item := range collection.Data {
		if item.Name == client.name {
			return true, nil
		}
	}
	client.logger.Debug("Certificate not found")
	return false, nil
}

func (client *certificateClient) existsInNamespace() (bool, error) {
	backendClient, err := client.project.backendProjectClient()
	if err != nil {
		return false, err
	}
	namespaceID, err := client.NamespaceID()
	if err != nil {
		return false, err
	}
	collection, err := backendClient.NamespacedCertificate.List(&types.ListOpts{
		Filters: map[string]interface{}{
			"name":        client.name,
			"namespaceId": namespaceID,
		},
	})
	if nil != err {
		client.logger.WithError(err).Error("Failed to read certificate list")
		return false, fmt.Errorf("Failed to read certificate list, %v", err)
	}
	for _, item := range collection.Data {
		if item.Name == client.name && item.NamespaceId == client.namespaceID {
			return true, nil
		}
	}
	client.logger.Debug("Certificate not found")
	return false, nil
}

func (client *certificateClient) Create() error {
	projectID, err := client.project.ID()
	if err != nil {
		return fmt.Errorf("Failed to read namespace ID, %v", err)
	}
	client.logger.Info("Create new Certificate")
	labels := make(map[string]string)
	labels["cattlectl.io/hash"] = hashOf(client.certificate)

	if client.namespace != "" {
		return client.createInNamespace(projectID, labels)
	}
	return client.createInProject(projectID, labels)
}

func (client *certificateClient) createInProject(projectID string, labels map[string]string) error {
	backendClient, err := client.project.backendProjectClient()
	if err != nil {
		return err
	}
	newCertificate := &backendProjectClient.Certificate{
		Name:      client.certificate.Name,
		Labels:    labels,
		Key:       client.certificate.Key,
		Certs:     client.certificate.Certs,
		ProjectID: projectID,
	}

	_, err = backendClient.Certificate.Create(newCertificate)
	return err
}

func (client *certificateClient) createInNamespace(projectID string, labels map[string]string) error {
	backendClient, err := client.project.backendProjectClient()
	if err != nil {
		return err
	}
	namespaceID, err := client.NamespaceID()
	if err != nil {
		return err
	}
	newCertificate := &backendProjectClient.NamespacedCertificate{
		Name:        client.certificate.Name,
		Labels:      labels,
		Key:         client.certificate.Key,
		Certs:       client.certificate.Certs,
		NamespaceId: namespaceID,
		ProjectID:   projectID,
	}

	_, err = backendClient.NamespacedCertificate.Create(newCertificate)
	return err
}

func (client *certificateClient) Upgrade() error {
	if client.namespace != "" {
		return client.upgradeInNamespace()
	}
	return client.upgradeInProject()
}

func (client *certificateClient) upgradeInProject() error {
	backendClient, err := client.project.backendProjectClient()
	if err != nil {
		return err
	}
	collection, err := backendClient.Certificate.List(&types.ListOpts{
		Filters: map[string]interface{}{
			"name": client.name,
		},
	})
	if nil != err {
		client.logger.WithError(err).Error("Failed to read certificate list")
		return fmt.Errorf("Failed to read certificate list, %v", err)
	}

	if len(collection.Data) == 0 {
		return fmt.Errorf("Certificate %v not found", client.name)
	}
	existingCertificate := collection.Data[0]
	if strings.TrimSpace(existingCertificate.Certs) == strings.TrimSpace(client.certificate.Certs) {
		client.logger.Debug("Skip upgrade certificate - no changes")
		return nil
	}
	client.logger.Info("Upgrade Certificate")
	existingCertificate.Key = client.certificate.Key
	existingCertificate.Certs = client.certificate.Certs

	_, err = backendClient.Certificate.Replace(&existingCertificate)
	return err
}

func (client *certificateClient) upgradeInNamespace() error {
	backendClient, err := client.project.backendProjectClient()
	if err != nil {
		return err
	}
	namespaceID, err := client.NamespaceID()
	if err != nil {
		return err
	}
	collection, err := backendClient.NamespacedCertificate.List(&types.ListOpts{
		Filters: map[string]interface{}{
			"name":        client.name,
			"namespaceId": namespaceID,
		},
	})
	if nil != err {
		client.logger.WithError(err).Error("Failed to read certificate list")
		return fmt.Errorf("Failed to read certificate list, %v", err)
	}

	if len(collection.Data) == 0 {
		return fmt.Errorf("Certificate %v not found", client.name)
	}
	existingCertificate := collection.Data[0]
	if strings.TrimSpace(existingCertificate.Certs) == strings.TrimSpace(client.certificate.Certs) {
		client.logger.Debug("Skip upgrade certificate - no changes")
		return nil
	}
	client.logger.Info("Upgrade Certificate")
	existingCertificate.Key = client.certificate.Key
	existingCertificate.Certs = client.certificate.Certs

	_, err = backendClient.NamespacedCertificate.Replace(&existingCertificate)
	return err
}

func (client *certificateClient) Data() (projectModel.Certificate, error) {
	return client.certificate, nil
}

func (client *certificateClient) SetData(certificate projectModel.Certificate) error {
	client.name = certificate.Name
	client.namespace = certificate.Namespace
	client.namespaceID = ""
	client.certificate = certificate
	return nil
}
