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

	clusterModel "github.com/bitgrip/cattlectl/internal/pkg/rancher/cluster/model"
	rancherModel "github.com/bitgrip/cattlectl/internal/pkg/rancher/model"
	"github.com/rancher/norman/types"
	backendClusterClient "github.com/rancher/types/client/cluster/v3"
	backendRancherClient "github.com/rancher/types/client/management/v3"
	"github.com/sirupsen/logrus"
)

func newClusterClient(
	name string,
	config RancherConfig,
	rancherClient RancherClient,
	logger *logrus.Entry,
) (ClusterClient, error) {
	return &clusterClient{
		resourceClient: resourceClient{
			name:   name,
			logger: logger.WithField("cluster_name", name),
		},
		config:            config,
		rancherClient:     rancherClient,
		projectClients:    make(map[string]ProjectClient),
		storageClasses:    make(map[string]StorageClassClient),
		persistentVolumes: make(map[string]PersistentVolumeClient),
		namespaces:        make(map[string]namespaceCacheEntry),
		catalogClients:    make(map[string]CatalogClient),
	}, nil
}

type clusterClient struct {
	resourceClient
	config                RancherConfig
	rancherClient         RancherClient
	_backendClusterClient *backendClusterClient.Client
	cluster               clusterModel.Cluster
	projectClients        map[string]ProjectClient
	storageClasses        map[string]StorageClassClient
	persistentVolumes     map[string]PersistentVolumeClient
	namespaces            map[string]namespaceCacheEntry
	catalogClients        map[string]CatalogClient
}

type namespaceCacheEntry struct {
	projectName string
	namespace   NamespaceClient
}

func (client *clusterClient) Type() string {
	return rancherModel.Cluster
}

func (client *clusterClient) init() error {
	if client._backendClusterClient != nil {
		return nil
	}
	var (
		id  string
		err error
	)
	if id, err = client.ID(); err != nil {
		return err
	}
	client._backendClusterClient, err = createBackendClusterClient(client.config, id)
	return err
}

func (client *clusterClient) ID() (string, error) {
	if client.id != "" {
		return client.id, nil
	}
	backendRancherClient, err := client.backendRancherClient()
	if err != nil {
		return "", err
	}
	collection, err := backendRancherClient.Cluster.List(&types.ListOpts{
		Filters: map[string]interface{}{
			"name": client.name,
		},
	})
	if err != nil {
		return "", err
	}
	if len(collection.Data) < 1 {
		return "", fmt.Errorf("Unknown Cluster [%s]", client.name)
	}
	client.id = collection.Data[0].ID
	return client.id, nil
}

func (client *clusterClient) Exists() (bool, error) {
	_, err := client.ID()
	return err == nil, err
}
func (client *clusterClient) Project(name string) (ProjectClient, error) {
	if cache, exists := client.projectClients[name]; exists {
		return cache, nil
	}
	project, err := newProjectClient(name, client.config, client, client.logger)
	if err != nil {
		return nil, err
	}
	client.projectClients[name] = project
	return project, nil
}
func (client *clusterClient) Projects() ([]ProjectClient, error) {
	backendRancherClient, err := client.backendRancherClient()
	if err != nil {
		return nil, err
	}
	collection, err := backendRancherClient.Project.List(&types.ListOpts{
		Filters: map[string]interface{}{
			"clusterId": client.id,
		},
	})
	if err != nil {
		return nil, err
	}
	result := make([]ProjectClient, len(collection.Data))
	for i, backendProject := range collection.Data {
		project, err := client.Project(backendProject.Name)
		if err != nil {
			return nil, err
		}
		result[i] = project
	}
	return result, nil
}
func (client *clusterClient) StorageClass(name string) (StorageClassClient, error) {
	if cache, exists := client.storageClasses[name]; exists {
		return cache, nil
	}
	storageClass, err := newStorageClassClient(name, client, client.logger)
	if err != nil {
		return nil, err
	}
	client.storageClasses[name] = storageClass
	return storageClass, nil
}
func (client *clusterClient) StorageClasses() ([]StorageClassClient, error) {
	backendClusterClient, err := client.backendClusterClient()
	if err != nil {
		return nil, err
	}
	collection, err := backendClusterClient.StorageClass.List(&types.ListOpts{
		Filters: map[string]interface{}{
			"clusterId": client.id,
		},
	})
	if err != nil {
		return nil, err
	}
	result := make([]StorageClassClient, len(collection.Data))
	for i, backendStorageClass := range collection.Data {
		storageClass, err := client.StorageClass(backendStorageClass.Name)
		if err != nil {
			return nil, err
		}
		result[i] = storageClass
	}
	return result, nil
}
func (client *clusterClient) PersistentVolume(name string) (PersistentVolumeClient, error) {
	if cache, exists := client.persistentVolumes[name]; exists {
		return cache, nil
	}
	persistentVolume, err := newPersistentVolumeClient(name, client, client.logger)
	if err != nil {
		return nil, err
	}
	client.persistentVolumes[name] = persistentVolume
	return persistentVolume, nil
}
func (client *clusterClient) PersistentVolumes() ([]PersistentVolumeClient, error) {
	backendClusterClient, err := client.backendClusterClient()
	if err != nil {
		return nil, err
	}
	collection, err := backendClusterClient.PersistentVolume.List(&types.ListOpts{
		Filters: map[string]interface{}{
			"clusterId": client.id,
		},
	})
	if err != nil {
		return nil, err
	}
	result := make([]PersistentVolumeClient, len(collection.Data))
	for i, backendPersistentVolume := range collection.Data {
		persistentVolume, err := client.PersistentVolume(backendPersistentVolume.Name)
		if err != nil {
			return nil, err
		}
		result[i] = persistentVolume
	}
	return result, nil
}
func (client *clusterClient) Namespace(name, projectName string) (NamespaceClient, error) {
	var (
		projectClient ProjectClient
		logger        *logrus.Entry
		err           error
	)
	logger = client.logger
	if projectName != "" {
		projectClient, err = client.Project(projectName)
		if err != nil {
			return nil, err
		}
		logger = logger.WithField("project_name", projectName)
	}
	if cache, exists := client.namespaces[name]; exists {
		if cache.projectName != projectName {
			return nil, fmt.Errorf("Namespace %s is part of project: %s", name, cache.projectName)
		}
		return cache.namespace, nil
	}
	namespace, err := newNamespaceClient(name, projectClient, client, logger)
	if err != nil {
		return nil, err
	}
	client.namespaces[name] = namespaceCacheEntry{
		projectName: projectName,
		namespace:   namespace,
	}
	return namespace, nil
}
func (client *clusterClient) Namespaces(projectName string) ([]NamespaceClient, error) {
	backendClusterClient, err := client.backendClusterClient()
	if err != nil {
		return nil, err
	}

	filters := map[string]interface{}{
		"clusterId": client.id,
	}
	if projectName != "" {
		project, err := client.Project(projectName)
		if err != nil {
			return nil, err
		}
		projectID, err := project.ID()
		if err != nil {
			return nil, err
		}
		filters["projectId"] = projectID
	}
	collection, err := backendClusterClient.Namespace.List(&types.ListOpts{
		Filters: filters,
	})
	if err != nil {
		return nil, err
	}
	result := make([]NamespaceClient, len(collection.Data))
	for i, backendNamespace := range collection.Data {
		namespace, err := client.Namespace(backendNamespace.Name, projectName)
		if err != nil {
			return nil, err
		}
		result[i] = namespace
	}
	return result, nil
}
func (client *clusterClient) backendRancherClient() (*backendRancherClient.Client, error) {
	return client.rancherClient.backendRancherClient()
}
func (client *clusterClient) backendClusterClient() (*backendClusterClient.Client, error) {
	if err := client.init(); err != nil {
		return nil, err
	}
	return client._backendClusterClient, nil
}

func (client *clusterClient) Catalog(catalogName string) (CatalogClient, error) {
	if cache, exists := client.catalogClients[catalogName]; exists {
		return cache, nil
	}
	result, err := newClusterCatalogClient(catalogName, client, client.logger)
	if err != nil {
		return nil, err
	}
	client.catalogClients[catalogName] = result
	return result, nil
}

func (client *clusterClient) Catalogs() ([]CatalogClient, error) {
	backendRancherClient, err := client.backendRancherClient()
	if err != nil {
		return nil, err
	}
	collection, err := backendRancherClient.ClusterCatalog.List(&types.ListOpts{
		Filters: map[string]interface{}{
			"clusterId": client.id,
		},
	})
	if err != nil {
		return nil, err
	}
	result := make([]CatalogClient, len(collection.Data))
	for i, backendCatalog := range collection.Data {
		catalog, err := client.Catalog(backendCatalog.Name)
		if err != nil {
			return nil, err
		}
		result[i] = catalog
	}
	return result, nil
}
