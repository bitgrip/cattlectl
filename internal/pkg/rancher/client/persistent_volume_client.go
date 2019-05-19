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

func newPersistentVolumeClientWithData(
	persistentVolume projectModel.PersistentVolume,
	namespace string,
	project ProjectClient,
	backendClusterClient *backendClusterClient.Client,
	logger *logrus.Entry,
) (PersistentVolumeClient, error) {
	result, err := newPersistentVolumeClient(
		persistentVolume.Name,
		namespace,
		project,
		backendClusterClient,
		logger,
	)
	if err != nil {
		return nil, err
	}
	err = result.SetData(persistentVolume)
	return result, err
}

func newPersistentVolumeClient(
	name, namespace string,
	project ProjectClient,
	backendClusterClient *backendClusterClient.Client,
	logger *logrus.Entry,
) (PersistentVolumeClient, error) {
	return &persistentVolumeClient{
		resourceClient: resourceClient{
			name:   name,
			logger: logger.WithField("persistentVolume_name", name).WithField("namespace", namespace),
		},
		backendClusterClient: backendClusterClient,
	}, nil
}

type persistentVolumeClient struct {
	resourceClient
	persistentVolume     projectModel.PersistentVolume
	backendClusterClient *backendClusterClient.Client
}

func (client *persistentVolumeClient) init() error {
	return nil
}

func (client *persistentVolumeClient) Exists() (bool, error) {
	if err := client.init(); err != nil {
		return false, err
	}
	collection, err := client.backendClusterClient.PersistentVolume.List(&types.ListOpts{
		Filters: map[string]interface{}{
			"name": client.name,
		},
	})
	if nil != err {
		client.logger.WithError(err).Error("Failed to read persistentVolume list")
		return false, fmt.Errorf("Failed to read persistentVolume list, %v", err)
	}
	for _, item := range collection.Data {
		if item.Name == client.name {
			return true, nil
		}
	}
	client.logger.Debug("PersistentVolume not found")
	return false, nil
}

func (client *persistentVolumeClient) Create() error {
	if err := client.init(); err != nil {
		return err
	}
	client.logger.Info("Create new persistent volume")
	newPersistentVolume := &backendClusterClient.PersistentVolume{
		Name:                          client.persistentVolume.Name,
		StorageClassID:                client.persistentVolume.StorageClassName,
		AccessModes:                   client.persistentVolume.AccessModes,
		Capacity:                      map[string]string{"storage": client.persistentVolume.Capacity},
		PersistentVolumeReclaimPolicy: "Delete",
		Local: &backendClusterClient.LocalVolumeSource{
			Path: client.persistentVolume.Path,
		},
		NodeAffinity: &backendClusterClient.VolumeNodeAffinity{
			Required: &backendClusterClient.NodeSelector{
				NodeSelectorTerms: []backendClusterClient.NodeSelectorTerm{
					backendClusterClient.NodeSelectorTerm{
						MatchExpressions: []backendClusterClient.NodeSelectorRequirement{
							backendClusterClient.NodeSelectorRequirement{
								Key:      "kubernetes.io/hostname",
								Operator: "In",
								Values:   []string{client.persistentVolume.Node},
							},
						},
					},
				},
			},
		},
	}

	_, err := client.backendClusterClient.PersistentVolume.Create(newPersistentVolume)
	return err
}

func (client *persistentVolumeClient) Upgrade() error {
	return fmt.Errorf("upgrade statefulset not supported")
}

func (client *persistentVolumeClient) Data() (projectModel.PersistentVolume, error) {
	return client.persistentVolume, nil
}

func (client *persistentVolumeClient) SetData(persistentVolume projectModel.PersistentVolume) error {
	client.name = persistentVolume.Name
	client.persistentVolume = persistentVolume
	return nil
}
