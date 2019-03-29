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
