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
	"testing"

	"github.com/bitgrip/cattlectl/internal/pkg/assert"
	projectModel "github.com/bitgrip/cattlectl/internal/pkg/rancher/project/model"
	"github.com/bitgrip/cattlectl/internal/pkg/rancher/stubs"
	"github.com/rancher/norman/types"
	clusterClient "github.com/rancher/types/client/cluster/v3"
	"github.com/sirupsen/logrus"
)

/*
HasPersistentVolume(...)
  returns true if local persisternt volume exists
*/
func TestHasPersistentVolume_PersistentVolumeExisting(t *testing.T) {
	const (
		projectId            = "test-project-id"
		projectName          = "test-project-name"
		clusterId            = "test-cluster-id"
		persistentVolumeName = "test-local-volume"
	)
	var (
		actualListOpts   *types.ListOpts
		clientConfig     = ClientConfig{}
		persistentVolume = projectModel.PersistentVolume{
			Name: persistentVolumeName,
		}
		testClients = stubs.CreateBackendStubs(t)
	)

	persistentVolumeOperationsStub := stubs.CreatePersistentVolumeOperationsStub(t)
	persistentVolumeOperationsStub.DoList = func(opts *types.ListOpts) (*clusterClient.PersistentVolumeCollection, error) {
		actualListOpts = opts
		return &clusterClient.PersistentVolumeCollection{
			Data: []clusterClient.PersistentVolume{
				clusterClient.PersistentVolume{
					Name: persistentVolumeName,
				},
			},
		}, nil
	}
	testClients.ClusterClient.PersistentVolume = persistentVolumeOperationsStub

	client := rancherClient{
		clientConfig:  clientConfig,
		projectId:     projectId,
		clusterClient: testClients.ClusterClient,
		logger:        logrus.WithField("test", true),
	}
	//Act
	result, err := client.HasPersistentVolume(persistentVolume)

	//Assert
	assert.Ok(t, err)
	assert.Equals(t, true, result)
	assert.Equals(t, &types.ListOpts{Filters: map[string]interface{}{"system": "false", "name": "test-local-volume"}}, actualListOpts)
}

/*
HasPersistentVolume(...)
  returns true if local persistent volume exists
*/
func TestHasPersistentVolume_PersistentVolumeNotExisting(t *testing.T) {
	const (
		projectId            = "test-project-id"
		projectName          = "test-project-name"
		clusterId            = "test-cluster-id"
		persistentVolumeName = "test-local-volume"
	)
	var (
		actualListOpts   *types.ListOpts
		clientConfig     = ClientConfig{}
		persistentVolume = projectModel.PersistentVolume{
			Name: persistentVolumeName,
		}
		testClients = stubs.CreateBackendStubs(t)
	)

	persistentVolumeOperationsStub := stubs.CreatePersistentVolumeOperationsStub(t)
	persistentVolumeOperationsStub.DoList = func(opts *types.ListOpts) (*clusterClient.PersistentVolumeCollection, error) {
		actualListOpts = opts
		return &clusterClient.PersistentVolumeCollection{
			Data: []clusterClient.PersistentVolume{},
		}, nil
	}
	testClients.ClusterClient.PersistentVolume = persistentVolumeOperationsStub

	client := rancherClient{
		clientConfig:  clientConfig,
		projectId:     projectId,
		clusterClient: testClients.ClusterClient,
		logger:        logrus.WithField("test", true),
	}
	//Act
	result, err := client.HasPersistentVolume(persistentVolume)

	//Assert
	assert.Ok(t, err)
	assert.Equals(t, false, result)
	assert.Equals(t, &types.ListOpts{Filters: map[string]interface{}{"system": "false", "name": "test-local-volume"}}, actualListOpts)
}

/*
CreatePersistentVolume(...)
  uses cluster client to create local persistent volume
*/
func TestCreatePersistentVolume(t *testing.T) {
	const (
		projectId            = "test-project-id"
		projectName          = "test-project-name"
		clusterId            = "test-cluster-id"
		persistentVolumeName = "test-local-volume"
		storageClassName     = "test-storage-class"
		capacity             = "3Gi"
		path                 = "/test/path"
		node                 = "test-node"
	)
	var (
		actualOpts       *clusterClient.PersistentVolume
		clientConfig     = ClientConfig{}
		persistentVolume = projectModel.PersistentVolume{
			Name:             persistentVolumeName,
			StorageClassName: storageClassName,
			AccessModes:      []string{"ReadWriteOnce"},
			Capacity:         capacity,
			Path:             path,
			Node:             node,
		}
		testClients = stubs.CreateBackendStubs(t)
	)

	persistentVolumeOperationsStub := stubs.CreatePersistentVolumeOperationsStub(t)
	persistentVolumeOperationsStub.DoCreate = func(opts *clusterClient.PersistentVolume) (*clusterClient.PersistentVolume, error) {
		actualOpts = opts
		return &clusterClient.PersistentVolume{}, nil

	}

	testClients.ClusterClient.PersistentVolume = persistentVolumeOperationsStub

	client := rancherClient{
		clientConfig:  clientConfig,
		projectId:     projectId,
		clusterClient: testClients.ClusterClient,
		logger:        logrus.WithField("test", true),
	}
	//Act
	err := client.CreatePersistentVolume(persistentVolume)

	//Assert
	assert.Ok(t, err)
	assert.Equals(t, &clusterClient.PersistentVolume{
		AccessModes: []string{"ReadWriteOnce"},
		Local:       &clusterClient.LocalVolumeSource{Path: path},
		NodeAffinity: &clusterClient.VolumeNodeAffinity{
			Required: &clusterClient.NodeSelector{
				NodeSelectorTerms: []clusterClient.NodeSelectorTerm{
					clusterClient.NodeSelectorTerm{
						MatchExpressions: []clusterClient.NodeSelectorRequirement{
							clusterClient.NodeSelectorRequirement{
								Key:      "kubernetes.io/hostname",
								Operator: "In",
								Values:   []string{node},
							},
						},
					},
				},
			},
		},
		Capacity:                      map[string]string{"storage": "3Gi"},
		Name:                          "test-local-volume",
		PersistentVolumeReclaimPolicy: "Delete",
		StorageClassID:                "test-storage-class",
	}, actualOpts)
}
