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

	managementClient "github.com/rancher/types/client/management/v3"
	"github.com/sirupsen/logrus"
)

// NewRancherClient creates a new rancher client
func NewRancherClient(config RancherConfig) (RancherClient, error) {
	return &rancherClient{
		config: config,
		logger: logrus.WithFields(logrus.Fields{}),
	}, nil
}

// RancherConfig holds the configuration data to interact with a rancher server
type RancherConfig struct {
	RancherURL string
	AccessKey  string
	SecretKey  string
	Insecure   bool
	CACerts    string
}

type rancherClient struct {
	config           RancherConfig
	managementClient *managementClient.Client
	logger           *logrus.Entry
}

func (client *rancherClient) init() error {
	if client.managementClient != nil {
		return nil
	}
	managementClient, err := createManagementClient(client.config)
	if err != nil {
		return err
	}
	client.managementClient = managementClient
	return nil
}

func (client *rancherClient) Cluster(clusterName, clusterID string) (ClusterClient, error) {
	if err := client.init(); err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("not yet implemented")
}

func (client *rancherClient) Clusters() ([]ClusterClient, error) {
	if err := client.init(); err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("not yet implemented")
}
