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
	"github.com/rancher/norman/types"
	"github.com/sirupsen/logrus"
)

func (client *rancherClient) HasClusterWithName(name string) (bool, string, error) {
	collection, err := client.managementClient.Cluster.List(&types.ListOpts{
		Filters: map[string]interface{}{
			"name": name,
		},
	})
	if err != nil {
		return false, "", err
	}
	if len(collection.Data) < 1 {
		return false, "", nil
	}
	return true, collection.Data[0].ID, nil
}

func (client *rancherClient) SetCluster(clusterName, clusterId string) error {
	if clusterClient, err := createClusterClient(
		client.clientConfig.RancherURL,
		client.clientConfig.AccessKey,
		client.clientConfig.SecretKey,
		clusterId,
	); err != nil {
		return err
	} else {
		client.logger = client.logger.WithFields(logrus.Fields{
			"cluster_name": clusterName,
			"project_name": "",
		})
		client.clusterId = clusterId
		client.clusterClient = clusterClient
		client.projectId = ""
		client.projectClient = nil
		return nil
	}
}
