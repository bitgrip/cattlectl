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

func newClusterCatalogClientWithData(
	catalog rancherModel.Catalog,
	clusterClient ClusterClient,
	logger *logrus.Entry,
) (CatalogClient, error) {
	result, err := newClusterCatalogClient(
		catalog.Name,
		clusterClient,
		logger,
	)
	if err != nil {
		return nil, err
	}
	err = result.SetData(catalog)
	return result, err
}

func newClusterCatalogClient(
	name string,
	clusterClient ClusterClient,
	logger *logrus.Entry,
) (CatalogClient, error) {
	return &clusterCatalogClient{
		resourceClient: resourceClient{
			name:   name,
			logger: logger.WithField("catalog_name", name),
		},
		clusterClient: clusterClient,
	}, nil
}

type clusterCatalogClient struct {
	resourceClient
	catalog       rancherModel.Catalog
	clusterClient ClusterClient
}

func (client *clusterCatalogClient) Type() string {
	return rancherModel.ClusterCatalog
}

func (client *clusterCatalogClient) Exists() (bool, error) {
	backendClient, err := client.clusterClient.backendRancherClient()
	if err != nil {
		return false, err
	}
	clusterID, err := client.clusterClient.ID()
	if err != nil {
		return false, err
	}
	collection, err := backendClient.ClusterCatalog.List(&types.ListOpts{
		Filters: map[string]interface{}{
			"name":      client.name,
			"clusterId": clusterID,
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

func (client *clusterCatalogClient) Create(dryRun bool) (changed bool, err error) {
	backendClient, err := client.clusterClient.backendRancherClient()
	if err != nil {
		return
	}
	clusterID, err := client.clusterClient.ID()
	if err != nil {
		return
	}
	client.logger.Info("Create new catalog")
	newClusterCatalog := &backendRancherClient.ClusterCatalog{
		Name:      client.catalog.Name,
		ClusterID: clusterID,
		URL:       client.catalog.URL,
		Branch:    client.catalog.Branch,
		Username:  client.catalog.Username,
		Password:  client.catalog.Password,
		Labels: map[string]string{
			"cattlectl.io/hash": hashOf(client.catalog),
		},
	}

	if dryRun {
		client.logger.WithField("object", newClusterCatalog).Info("Do Dry-Run Create")
	} else {
		_, err = backendClient.ClusterCatalog.Create(newClusterCatalog)
	}
	return err == nil, err
}

func (client *clusterCatalogClient) Upgrade(dryRun bool) (changed bool, err error) {
	backendClient, err := client.clusterClient.backendRancherClient()
	if err != nil {
		return
	}
	clusterID, err := client.clusterClient.ID()
	if err != nil {
		return
	}
	client.logger.Trace("Load from rancher")
	collection, err := backendClient.ClusterCatalog.List(&types.ListOpts{
		Filters: map[string]interface{}{
			"name":      client.name,
			"clusterId": clusterID,
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
	if isClusterCatalogUnchanged(existingCatalog, client.catalog) {
		client.logger.Debug("Skip upgrade catalog - no changes")
		return
	}
	client.logger.Info("Upgrade ClusterCatalog")
	existingCatalog.Labels["cattlectl.io/hash"] = hashOf(client.catalog)
	existingCatalog.URL = client.catalog.URL
	existingCatalog.Branch = client.catalog.Branch
	existingCatalog.Username = client.catalog.Username
	existingCatalog.Password = client.catalog.Password

	if dryRun {
		client.logger.WithField("object", existingCatalog).Info("Do Dry-Run Upgrade")
	} else {
		_, err = backendClient.ClusterCatalog.Replace(&existingCatalog)
	}
	return err == nil, err
}

func (client *clusterCatalogClient) Data() (rancherModel.Catalog, error) {
	return client.catalog, nil
}

func (client *clusterCatalogClient) SetData(catalog rancherModel.Catalog) error {
	client.name = catalog.Name
	client.catalog = catalog
	return nil
}

func isClusterCatalogUnchanged(existingCatalog backendRancherClient.ClusterCatalog, catalog rancherModel.Catalog) bool {
	hash, hashExists := existingCatalog.Labels["cattlectl.io/hash"]
	if !hashExists {
		return false
	}
	return hash == hashOf(catalog)
}
