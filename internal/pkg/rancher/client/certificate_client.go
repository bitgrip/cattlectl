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

func newCertificateClientWithData(
	certificate projectModel.Certificate,
	namespace string,
	project ProjectClient,
	backendProjectClient *backendProjectClient.Client,
	logger *logrus.Entry,
) (CertificateClient, error) {
	result, err := newCertificateClient(
		certificate.Name,
		namespace,
		project,
		backendProjectClient,
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
	backendProjectClient *backendProjectClient.Client,
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
		backendProjectClient: backendProjectClient,
	}, nil
}

type certificateClient struct {
	namespacedResourceClient
	certificate          projectModel.Certificate
	backendProjectClient *backendProjectClient.Client
}

func (client *certificateClient) init() error {
	namespaceID, err := client.NamespaceID()
	if namespaceID == "" && err == nil {
		return fmt.Errorf("Can not find namespace")
	}
	return err
}

func (client *certificateClient) Exists() (bool, error) {
	if err := client.init(); err != nil {
		return false, err
	}
	collection, err := client.backendProjectClient.Certificate.List(&types.ListOpts{
		Filters: map[string]interface{}{
			"name":        client.name,
			"namespaceId": client.namespaceID,
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
	if err := client.init(); err != nil {
		return err
	}
	namespaceID, err := client.NamespaceID()
	if err != nil {
		return fmt.Errorf("Failed to read namespace ID, %v", err)
	}
	projectID, err := client.project.ID()
	if err != nil {
		return fmt.Errorf("Failed to read namespace ID, %v", err)
	}
	client.logger.Info("Create new Certificate")
	labels := make(map[string]string)
	labels["cattlectl.io/hash"] = hashOf(client.certificate)
	newCertificate := &backendProjectClient.Certificate{
		Name:        client.certificate.Name,
		Labels:      labels,
		Key:         client.certificate.Key,
		Certs:       client.certificate.Certs,
		NamespaceId: namespaceID,
		ProjectID:   projectID,
	}

	_, err = client.backendProjectClient.Certificate.Create(newCertificate)
	return err
}

func (client *certificateClient) Upgrade() error {
	return fmt.Errorf("upgrade statefulset not supported")
}

func (client *certificateClient) Data() (projectModel.Certificate, error) {
	return client.certificate, nil
}

func (client *certificateClient) SetData(certificate projectModel.Certificate) error {
	client.name = certificate.Name
	client.certificate = certificate
	return nil
}
