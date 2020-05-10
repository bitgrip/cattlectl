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
	"reflect"
	"testing"

	"github.com/bitgrip/cattlectl/internal/pkg/assert"
	projectModel "github.com/bitgrip/cattlectl/internal/pkg/rancher/cluster/project/model"
	"github.com/bitgrip/cattlectl/internal/pkg/rancher/stubs"
	"github.com/rancher/norman/types"
	backendClusterClient "github.com/rancher/types/client/cluster/v3"
	"github.com/sirupsen/logrus"
)

func Test_storageClassClient_Exists(t *testing.T) {
	tests := []struct {
		name      string
		client    *storageClassClient
		wanted    bool
		wantErr   bool
		wantedErr string
	}{
		{
			name: "Existing",
			client: existingStorageClassClient(
				t,
				&types.ListOpts{
					Filters: map[string]interface{}{
						"name": "existing-storageClass",
					},
				},
			),
			wanted:  true,
			wantErr: false,
		},
		{
			name: "Not_Existing",
			client: notExistingStorageClassClient(
				t,
				&types.ListOpts{
					Filters: map[string]interface{}{
						"name": "existing-storageClass",
					},
				},
			),
			wanted:  false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.client.Exists()
			if tt.wantErr {
				assert.NotOk(t, err, tt.wantedErr)
			} else {
				assert.Ok(t, err)
				assert.Equals(t, tt.wanted, got)
			}
		})
	}
}

func Test_storageClassClient_Create(t *testing.T) {
	tests := []struct {
		name      string
		client    *storageClassClient
		wantErr   bool
		wantedErr string
	}{
		{
			name: "Create",
			client: notExistingStorageClassClient(
				t,
				nil,
			),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.client.Create(false)
			if tt.wantErr {
				assert.NotOk(t, err, tt.wantedErr)
			} else {
				assert.Ok(t, err)
			}
		})
	}
}

func existingStorageClassClient(t *testing.T, expectedListOpts *types.ListOpts) *storageClassClient {
	const (
		projectID        = "test-project-id"
		projectName      = "test-project-name"
		namespaceID      = "test-namespace-id"
		namespace        = "test-namespace"
		clusterID        = "test-cluster-id"
		storageClassName = "test-storageClass"
	)
	var (
		storageClass = projectModel.StorageClass{}
		testClients  = stubs.CreateBackendStubs(t)
	)
	storageClass.Name = storageClassName

	storageClassOperationsStub := stubs.CreateStorageClassOperationsStub(t)
	storageClassOperationsStub.DoList = func(opts *types.ListOpts) (*backendClusterClient.StorageClassCollection, error) {
		if !reflect.DeepEqual(expectedListOpts, opts) {
			return nil, fmt.Errorf("Unexpected ListOpts %v", opts)
		}
		return &backendClusterClient.StorageClassCollection{
			Data: []backendClusterClient.StorageClass{
				backendClusterClient.StorageClass{
					Name: "existing-storageClass",
				},
			},
		}, nil
	}
	testClients.ClusterClient.StorageClass = storageClassOperationsStub
	clusterClient := simpleClusterClient()
	clusterClient._backendClusterClient = testClients.ClusterClient
	result, err := newStorageClassClient(
		"existing-storageClass",
		clusterClient,
		logrus.New().WithFields(logrus.Fields{}),
	)
	assert.Ok(t, err)
	storageClassClientResult := result.(*storageClassClient)
	return storageClassClientResult
}

func notExistingStorageClassClient(t *testing.T, expectedListOpts *types.ListOpts) *storageClassClient {
	const (
		projectID        = "test-project-id"
		projectName      = "test-project-name"
		namespaceID      = "test-namespace-id"
		namespace        = "test-namespace"
		clusterID        = "test-cluster-id"
		storageClassName = "test-storageClass"
	)
	var (
		storageClass = projectModel.StorageClass{}
		testClients  = stubs.CreateBackendStubs(t)
	)
	storageClass.Name = storageClassName

	storageClassOperationsStub := stubs.CreateStorageClassOperationsStub(t)
	storageClassOperationsStub.DoList = func(opts *types.ListOpts) (*backendClusterClient.StorageClassCollection, error) {
		if !reflect.DeepEqual(expectedListOpts, opts) {
			return nil, fmt.Errorf("Unexpected ListOpts %v", opts)
		}
		return &backendClusterClient.StorageClassCollection{
			Data: []backendClusterClient.StorageClass{},
		}, nil
	}
	storageClassOperationsStub.DoCreate = func(storageClass *backendClusterClient.StorageClass) (*backendClusterClient.StorageClass, error) {
		return storageClass, nil
	}
	testClients.ClusterClient.StorageClass = storageClassOperationsStub
	clusterClient := simpleClusterClient()
	clusterClient._backendClusterClient = testClients.ClusterClient
	result, err := newStorageClassClient(
		"existing-storageClass",
		clusterClient,
		logrus.New().WithFields(logrus.Fields{}),
	)
	assert.Ok(t, err)
	storageClassClientResult := result.(*storageClassClient)
	return storageClassClientResult
}
