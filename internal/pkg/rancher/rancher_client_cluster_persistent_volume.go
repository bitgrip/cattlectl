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
	clusterClient "github.com/rancher/types/client/cluster/v3"
)

func (client *rancherClient) HasPersistentVolume(persistentVolume projectModel.PersistentVolume) (bool, error) {
	collection, err := client.clusterClient.PersistentVolume.List(&types.ListOpts{
		Filters: map[string]interface{}{
			"system": "false",
			"name":   persistentVolume.Name,
		},
	})
	if nil != err {
		client.logger.WithError(err).WithField("persistent_volume_name", persistentVolume.Name).Error("Failed to read persistent volume list")
		return false, fmt.Errorf("Failed to read persistent volume list, %v", err)
	}
	for _, item := range collection.Data {
		if item.Name == persistentVolume.Name {
			return true, nil
		}
	}
	client.logger.WithField("persistent_volume_name", persistentVolume.Name).Debug("Persistent volume not found")
	return false, nil
}

func (client *rancherClient) CreatePersistentVolume(persistentVolume projectModel.PersistentVolume) error {
	client.logger.WithField("local_volume_name", persistentVolume.Name).Info("Create new persistent volume")
	newPersistentVolume := &clusterClient.PersistentVolume{
		Name:                          persistentVolume.Name,
		StorageClassID:                persistentVolume.StorageClassName,
		AccessModes:                   persistentVolume.AccessModes,
		Capacity:                      map[string]string{"storage": persistentVolume.Capacity},
		PersistentVolumeReclaimPolicy: "Delete",
		Local: &clusterClient.LocalVolumeSource{
			Path: persistentVolume.Path,
		},
		NodeAffinity: &clusterClient.VolumeNodeAffinity{
			Required: &clusterClient.NodeSelector{
				NodeSelectorTerms: []clusterClient.NodeSelectorTerm{
					clusterClient.NodeSelectorTerm{
						MatchExpressions: []clusterClient.NodeSelectorRequirement{
							clusterClient.NodeSelectorRequirement{
								Key:      "kubernetes.io/hostname",
								Operator: "In",
								Values:   []string{persistentVolume.Node},
							},
						},
					},
				},
			},
		},
	}

	_, err := client.clusterClient.PersistentVolume.Create(newPersistentVolume)
	return err
}
