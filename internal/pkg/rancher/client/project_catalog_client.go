// Copyright © 2019 Bitgrip <berlin@bitgrip.de>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by cataloglicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package client

import (
	"fmt"

	rancherModel "github.com/bitgrip/cattlectl/internal/pkg/rancher/model"
	"github.com/rancher/norman/types"
	backendRancherClient "github.com/rancher/types/client/management/v3"
	"github.com/sirupsen/logrus"
)

func newProjectCatalogClientWithData(
	catalog rancherModel.Catalog,
	projectClient ProjectClient,
	logger *logrus.Entry,
) (CatalogClient, error) {
	result, err := newProjectCatalogClient(
		catalog.Name,
		projectClient,
		logger,
	)
	if err != nil {
		return nil, err
	}
	err = result.SetData(catalog)
	return result, err
}

func newProjectCatalogClient(
	name string,
	projectClient ProjectClient,
	logger *logrus.Entry,
) (CatalogClient, error) {
	return &projectCatalogClient{
		resourceClient: resourceClient{
			name:   name,
			logger: logger.WithField("catalog_name", name),
		},
		projectClient: projectClient,
	}, nil
}

type projectCatalogClient struct {
	resourceClient
	catalog       rancherModel.Catalog
	projectClient ProjectClient
}

func (client *projectCatalogClient) Type() string {
	return rancherModel.ProjectCatalog
}

func (client *projectCatalogClient) Exists() (bool, error) {
	backendClient, err := client.projectClient.backendRancherClient()
	if err != nil {
		return false, err
	}
	projectID, err := client.projectClient.ID()
	if err != nil {
		return false, err
	}
	collection, err := backendClient.ProjectCatalog.List(&types.ListOpts{
		Filters: map[string]interface{}{
			"name":      client.name,
			"projectID": projectID,
		},
	})
	if nil != err {
		client.logger.WithError(err).Error("Failed to read catalog list")
		return false, fmt.Errorf("Failed to read catalog list, %v", err)
	}
	for _, item := range collection.Data {
		if item.Name == client.name {
			return true, nil
		}
	}
	client.logger.Debug("Catalog not found")
	return false, nil
}

func (client *projectCatalogClient) Create(dryRun bool) (changed bool, err error) {
	backendClient, err := client.projectClient.backendRancherClient()
	if err != nil {
		return
	}
	projectID, err := client.projectClient.ID()
	if err != nil {
		return
	}
	client.logger.Info("Create new catalog")
	newProjectCatalog := backendRancherClient.ProjectCatalog{
		Name:      client.catalog.Name,
		ProjectID: projectID,
		URL:       client.catalog.URL,
		Branch:    client.catalog.Branch,
		Username:  client.catalog.Username,
		Password:  client.catalog.Password,
		Labels: map[string]string{
			"cattlectl.io/hash": hashOf(client.catalog),
		},
	}

	if dryRun {
		client.logger.WithField("object", newProjectCatalog).Info("Do Dry-Run Create")
	} else {
		_, err = backendClient.ProjectCatalog.Create(&newProjectCatalog)
	}
	return err == nil, err
}

func (client *projectCatalogClient) Upgrade(dryRun bool) (changed bool, err error) {
	backendClient, err := client.projectClient.backendRancherClient()
	if err != nil {
		return
	}
	projectID, err := client.projectClient.ID()
	if err != nil {
		return
	}
	client.logger.Trace("Load from rancher")
	collection, err := backendClient.ProjectCatalog.List(&types.ListOpts{
		Filters: map[string]interface{}{
			"name":      client.name,
			"projectID": projectID,
		},
	})
	if nil != err {
		client.logger.WithError(err).Error("Failed to read catalog list")
		return changed, fmt.Errorf("Failed to read catalog list, %v", err)
	}

	if len(collection.Data) == 0 {
		return changed, fmt.Errorf("Catalog %v not found", client.name)
	}

	existingCatalog := collection.Data[0]
	if isProjectCatalogUnchanged(existingCatalog, client.catalog) {
		client.logger.Debug("Skip upgrade catalog - no changes")
		return
	}
	client.logger.Info("Upgrade ProjectCatalog")
	existingCatalog.Labels["cattlectl.io/hash"] = hashOf(client.catalog)
	existingCatalog.URL = client.catalog.URL
	existingCatalog.Branch = client.catalog.Branch
	existingCatalog.Username = client.catalog.Username
	existingCatalog.Password = client.catalog.Password

	if dryRun {
		client.logger.WithField("object", existingCatalog).Info("Do Dry-Run Upgrade")
	} else {
		_, err = backendClient.ProjectCatalog.Replace(&existingCatalog)
	}
	return err == nil, err
}

func (client *projectCatalogClient) Data() (rancherModel.Catalog, error) {
	return client.catalog, nil
}

func (client *projectCatalogClient) SetData(catalog rancherModel.Catalog) error {
	client.name = catalog.Name
	client.catalog = catalog
	return nil
}

func isProjectCatalogUnchanged(existingCatalog backendRancherClient.ProjectCatalog, catalog rancherModel.Catalog) bool {
	hash, hashExists := existingCatalog.Labels["cattlectl.io/hash"]
	if !hashExists {
		return false
	}
	return hash == hashOf(catalog)
}
