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
	"github.com/rancher/norman/types"
	backendRancherClient "github.com/rancher/types/client/management/v3"
	"github.com/sirupsen/logrus"
)

// NewRancherClient creates a new rancher client
func NewRancherClient(config RancherConfig) (RancherClient, error) {
	return &rancherClient{
		config:         config,
		logger:         logrus.WithFields(logrus.Fields{}),
		clusterClients: make(map[string]ClusterClient),
	}, nil
}

// RancherConfig holds the configuration data to interact with a rancher server
type RancherConfig struct {
	RancherURL   string
	AccessKey    string
	SecretKey    string
	Insecure     bool
	CACerts      string
	MergeAnswers bool
}

type rancherClient struct {
	config                RancherConfig
	_backendRancherClient *backendRancherClient.Client
	logger                *logrus.Entry
	clusterClients        map[string]ClusterClient
}

func (client *rancherClient) init() error {
	if client._backendRancherClient != nil {
		return nil
	}
	backendRancherClient, err := createManagementClient(client.config)
	if err != nil {
		return err
	}
	client._backendRancherClient = backendRancherClient
	return nil
}

func (client *rancherClient) Cluster(clusterName string) (ClusterClient, error) {
	if cache, exists := client.clusterClients[clusterName]; exists {
		return cache, nil
	}
	result, err := newClusterClient(clusterName, client.config, client, client.logger)
	if err != nil {
		return nil, err
	}
	client.clusterClients[clusterName] = result
	return result, nil
}

func (client *rancherClient) Clusters() ([]ClusterClient, error) {
	backendRancherClient, err := client.backendRancherClient()
	if err != nil {
		return nil, err
	}
	collection, err := backendRancherClient.Cluster.List(&types.ListOpts{
		Filters: map[string]interface{}{},
	})
	if err != nil {
		return nil, err
	}
	result := make([]ClusterClient, len(collection.Data))
	for i, backendCluster := range collection.Data {
		cluster, err := client.Cluster(backendCluster.Name)
		if err != nil {
			return nil, err
		}
		result[i] = cluster
	}
	return result, nil
}

func (client *rancherClient) backendRancherClient() (*backendRancherClient.Client, error) {
	if err := client.init(); err != nil {
		return nil, err
	}
	return client._backendRancherClient, nil
}
