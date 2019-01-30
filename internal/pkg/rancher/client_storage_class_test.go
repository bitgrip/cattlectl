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
	"github.com/bitgrip/cattlectl/internal/pkg/rancher/stubs"
	"github.com/rancher/norman/types"
	clusterClient "github.com/rancher/types/client/cluster/v3"
	"github.com/sirupsen/logrus"
)

/*
HasStorageClass(...)
  returns true if storage class exists
*/
func TestHasStorageClass_StorageClassExisting(t *testing.T) {
	const (
		projectId        = "test-project-id"
		projectName      = "test-project-name"
		clusterId        = "test-cluster-id"
		storageClassName = "test-storage-class"
	)
	var (
		actualListOpts *types.ListOpts
		clientConfig   = ClientConfig{}
		storageClass   = StorageClass{
			Name: storageClassName,
		}
		testClients = stubs.CreateTestClients(t)
	)

	storageClassOperationsStub := stubs.CreateStorageClassOperationsStub(t)
	storageClassOperationsStub.DoList = func(opts *types.ListOpts) (*clusterClient.StorageClassCollection, error) {
		actualListOpts = opts
		return &clusterClient.StorageClassCollection{
			Data: []clusterClient.StorageClass{
				clusterClient.StorageClass{
					Name: storageClassName,
				},
			},
		}, nil

	}
	testClients.ClusterClient.StorageClass = storageClassOperationsStub

	client := rancherClient{
		clientConfig:  clientConfig,
		projectId:     projectId,
		clusterClient: testClients.ClusterClient,
		logger:        logrus.WithField("test", true),
	}
	//Act
	result, err := client.HasStorageClass(storageClass)

	//Assert
	assert.Ok(t, err)
	assert.Equals(t, true, result)
	assert.Equals(t, &types.ListOpts{Filters: map[string]interface{}{"system": "false", "name": "test-storage-class"}}, actualListOpts)
}

/*
HasNamespace(...)
  returns true if namespace exists
*/
func TestHasStorageClass_NamespaceNotExisting(t *testing.T) {
	const (
		projectId        = "test-project-id"
		projectName      = "test-project-name"
		clusterId        = "test-cluster-id"
		storageClassName = "test-storage-class"
	)
	var (
		actualListOpts *types.ListOpts
		clientConfig   = ClientConfig{}
		storageClass   = StorageClass{
			Name: storageClassName,
		}
		testClients = stubs.CreateTestClients(t)
	)

	storageClassOperationsStub := stubs.CreateStorageClassOperationsStub(t)
	storageClassOperationsStub.DoList = func(opts *types.ListOpts) (*clusterClient.StorageClassCollection, error) {
		actualListOpts = opts
		return &clusterClient.StorageClassCollection{
			Data: []clusterClient.StorageClass{},
		}, nil

	}
	testClients.ClusterClient.StorageClass = storageClassOperationsStub

	client := rancherClient{
		clientConfig:  clientConfig,
		projectId:     projectId,
		clusterClient: testClients.ClusterClient,
		logger:        logrus.WithField("test", true),
	}
	//Act
	result, err := client.HasStorageClass(storageClass)

	//Assert
	assert.Ok(t, err)
	assert.Equals(t, false, result)
	assert.Equals(t, &types.ListOpts{Filters: map[string]interface{}{"system": "false", "name": "test-storage-class"}}, actualListOpts)
}

/*
HasNamespace(...)
  returns true if namespace exists
*/
func TestCreateStorageClass(t *testing.T) {
	const (
		projectId        = "test-project-id"
		projectName      = "test-project-name"
		clusterId        = "test-cluster-id"
		storageClassName = "test-storage-class"
	)
	var (
		actualOpts   *clusterClient.StorageClass
		clientConfig = ClientConfig{}
		storageClass = StorageClass{
			Name: storageClassName,
		}
		testClients = stubs.CreateTestClients(t)
	)

	storageClassOperationsStub := stubs.CreateStorageClassOperationsStub(t)
	storageClassOperationsStub.DoCreate = func(opts *clusterClient.StorageClass) (*clusterClient.StorageClass, error) {
		actualOpts = opts
		return &clusterClient.StorageClass{}, nil

	}
	testClients.ClusterClient.StorageClass = storageClassOperationsStub

	client := rancherClient{
		clientConfig:  clientConfig,
		projectId:     projectId,
		clusterClient: testClients.ClusterClient,
		logger:        logrus.WithField("test", true),
	}
	//Act
	err := client.CreateStorageClass(storageClass)

	//Assert
	assert.Ok(t, err)
	assert.Equals(t, &clusterClient.StorageClass{Name: "test-storage-class"}, actualOpts)
}
