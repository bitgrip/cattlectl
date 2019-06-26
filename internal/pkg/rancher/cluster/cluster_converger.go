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

package rancher

import (
	"github.com/bitgrip/cattlectl/internal/pkg/rancher/client"
	clusterModel "github.com/bitgrip/cattlectl/internal/pkg/rancher/cluster/model"
	"github.com/bitgrip/cattlectl/internal/pkg/rancher/descriptor"
)

// NewClusterConverger creates a Converger for a given github.com/bitgrip/cattlectl/internal/pkg/rancher/cluster/model.Cluster
func NewClusterConverger(cluster clusterModel.Cluster, rancherClient client.RancherClient) (descriptor.Converger, error) {
	clusterClient, err := rancherClient.Cluster(cluster.Metadata.Name)
	if err != nil {
		return nil, err
	}
	childConvergers := make([]descriptor.Converger, 0)
	for _, catalog := range cluster.Catalogs {
		catalogClient, err := clusterClient.Catalog(catalog.Name)
		if err != nil {
			return nil, err
		}
		catalogClient.SetData(catalog)
		childConvergers = append(childConvergers, &descriptor.ResourceClientConverger{
			Client: catalogClient,
		})
	}
	return &descriptor.ResourceClientConverger{
		Client:   client.EmptyResourceClient,
		Children: childConvergers,
	}, nil
}
