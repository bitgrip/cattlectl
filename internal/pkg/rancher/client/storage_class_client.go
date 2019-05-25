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
	"github.com/sirupsen/logrus"
)

func newStorageClassClientWithData(
	storageClass projectModel.StorageClass,
	backendClusterClient *backendClusterClient.Client,
	logger *logrus.Entry,
) (StorageClassClient, error) {
	result, err := newStorageClassClient(
		storageClass.Name,
		backendClusterClient,
		logger,
	)
	if err != nil {
		return nil, err
	}
	err = result.SetData(storageClass)
	return result, err
}

func newStorageClassClient(
	name string,
	backendClusterClient *backendClusterClient.Client,
	logger *logrus.Entry,
) (StorageClassClient, error) {
	return &storageClassClient{
		resourceClient: resourceClient{
			name:   name,
			logger: logger.WithField("storageClass_name", name),
		},
		backendClusterClient: backendClusterClient,
	}, nil
}

type storageClassClient struct {
	resourceClient
	storageClass         projectModel.StorageClass
	backendClusterClient *backendClusterClient.Client
}

func (client *storageClassClient) init() error {
	return nil
}

func (client *storageClassClient) Exists() (bool, error) {
	if err := client.init(); err != nil {
		return false, err
	}
	collection, err := client.backendClusterClient.StorageClass.List(&types.ListOpts{
		Filters: map[string]interface{}{
			"name": client.name,
		},
	})
	if nil != err {
		client.logger.WithError(err).Error("Failed to read storageClass list")
		return false, fmt.Errorf("Failed to read storageClass list, %v", err)
	}
	for _, item := range collection.Data {
		if item.Name == client.name {
			return true, nil
		}
	}
	client.logger.Debug("StorageClass not found")
	return false, nil
}

func (client *storageClassClient) Create() error {
	if err := client.init(); err != nil {
		return err
	}
	client.logger.Info("Create new storage class")
	newStorageClass := &backendClusterClient.StorageClass{
		Name:              client.storageClass.Name,
		VolumeBindingMode: client.storageClass.VolumeBindMode,
		ReclaimPolicy:     client.storageClass.ReclaimPolicy,
		Provisioner:       client.storageClass.Provisioner,
		Parameters:        client.storageClass.Parameters,
		MountOptions:      client.storageClass.MountOptions,
	}

	_, err := client.backendClusterClient.StorageClass.Create(newStorageClass)
	return err
}

func (client *storageClassClient) Upgrade() error {
	return fmt.Errorf("upgrade statefulset not supported")
}

func (client *storageClassClient) Data() (projectModel.StorageClass, error) {
	return client.storageClass, nil
}

func (client *storageClassClient) SetData(storageClass projectModel.StorageClass) error {
	client.name = storageClass.Name
	client.storageClass = storageClass
	return nil
}
