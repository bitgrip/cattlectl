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

func (client *rancherClient) HasStorageClass(storageClass projectModel.StorageClass) (bool, error) {
	collection, err := client.clusterClient.StorageClass.List(&types.ListOpts{
		Filters: map[string]interface{}{
			"system": "false",
			"name":   storageClass.Name,
		},
	})
	if nil != err {
		client.logger.WithError(err).WithField("storage_class_name", storageClass.Name).Error("Failed to read storage class list")
		return false, fmt.Errorf("Failed to read storage class list, %v", err)
	}
	for _, item := range collection.Data {
		if item.Name == storageClass.Name {
			return true, nil
		}
	}
	client.logger.WithField("storage_class_name", storageClass.Name).Debug("Storage class not found")
	return false, nil
}

func (client *rancherClient) CreateStorageClass(storageClass projectModel.StorageClass) error {
	client.logger.WithField("storage_class_name", storageClass.Name).Info("Create new storage class")
	newStorageClass := &clusterClient.StorageClass{
		Name:              storageClass.Name,
		VolumeBindingMode: storageClass.VolumeBindMode,
		ReclaimPolicy:     storageClass.ReclaimPolicy,
		Provisioner:       storageClass.Provisioner,
		Parameters:        storageClass.Parameters,
		MountOptions:      storageClass.MountOptions,
	}

	_, err := client.clusterClient.StorageClass.Create(newStorageClass)
	return err
}
