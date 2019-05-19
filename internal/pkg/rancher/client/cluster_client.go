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
	backendClusterClient "github.com/rancher/types/client/cluster/v3"
	backendRancherClient "github.com/rancher/types/client/management/v3"
	"github.com/sirupsen/logrus"
)

func newClusterClient(
	name string,
	config RancherConfig,
	backendRancherClient *backendRancherClient.Client,
	logger *logrus.Entry,
) (ClusterClient, error) {
	return &clusterClient{
		resourceClient: resourceClient{
			name:   name,
			logger: logger.WithField("cluster_name", name),
		},
		config:               config,
		backendRancherClient: backendRancherClient,
	}, nil
}

type clusterClient struct {
	resourceClient
	config               RancherConfig
	backendRancherClient *backendRancherClient.Client
	backendClusterClient *backendClusterClient.Client
	cluster              projectModel.Cluster
}

func (client *clusterClient) init() error {
	if client.backendClusterClient != nil {
		return nil
	}
	var (
		id  string
		err error
	)
	if id, err = client.ID(); err != nil {
		return err
	}
	client.backendClusterClient, err = createBackendClusterClient(client.config, id)
	return err
}

func (client *clusterClient) ID() (string, error) {
	if client.id != "" {
		return client.id, nil
	}
	collection, err := client.backendRancherClient.Cluster.List(&types.ListOpts{
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
	return err != nil, err
}
func (client *clusterClient) Project(projectName string) (ProjectClient, error) {
	if err := client.init(); err != nil {
		return nil, err
	}
	return newProjectClient(projectName, client.config, client.backendRancherClient, client.backendClusterClient, client.logger)
}
func (client *clusterClient) Projects() ([]ProjectClient, error) {
	if err := client.init(); err != nil {
		return nil, err
	}
	collection, err := client.backendRancherClient.Project.List(&types.ListOpts{
		Filters: map[string]interface{}{
			"clusterId": client.id,
		},
	})
	if err != nil {
		return nil, err
	}
	result := make([]ProjectClient, len(collection.Data))
	for i, backendProject := range collection.Data {
		cluster, err := client.Project(backendProject.Name)
		if err != nil {
			return nil, err
		}
		result[i] = cluster
	}
	return result, nil
}
func (client *clusterClient) StorageClass(name string) (StorageClassClient, error) {
	if err := client.init(); err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("upgrade statefulset not supported")
}
func (client *clusterClient) StorageClasses() ([]StorageClassClient, error) {
	if err := client.init(); err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("upgrade statefulset not supported")
}
func (client *clusterClient) PersistentVolume(name string) (PersistentVolumeClient, error) {
	if err := client.init(); err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("upgrade statefulset not supported")
}
func (client *clusterClient) PersistentVolumes() ([]PersistentVolumeClient, error) {
	if err := client.init(); err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("upgrade statefulset not supported")
}
func (client *clusterClient) Namespace(name, projectName string) (NamespaceClient, error) {
	if err := client.init(); err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("upgrade statefulset not supported")
}
func (client *clusterClient) Namespaces(projectName string) ([]NamespaceClient, error) {
	if err := client.init(); err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("upgrade statefulset not supported")
}
