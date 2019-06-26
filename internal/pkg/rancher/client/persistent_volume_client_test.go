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

func Test_persistentVolumeClient_Exists(t *testing.T) {
	tests := []struct {
		name      string
		client    *persistentVolumeClient
		wanted    bool
		wantErr   bool
		wantedErr string
	}{
		{
			name: "Existing",
			client: existingPersistentVolumeClient(
				t,
				&types.ListOpts{
					Filters: map[string]interface{}{
						"name": "existing-persistentVolume",
					},
				},
			),
			wanted:  true,
			wantErr: false,
		},
		{
			name: "Not_Existing",
			client: notExistingPersistentVolumeClient(
				t,
				&types.ListOpts{
					Filters: map[string]interface{}{
						"name": "existing-persistentVolume",
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

func Test_persistentVolumeClient_Create(t *testing.T) {
	tests := []struct {
		name      string
		client    *persistentVolumeClient
		wantErr   bool
		wantedErr string
	}{
		{
			name: "Create",
			client: notExistingPersistentVolumeClient(
				t,
				nil,
			),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.client.Create()
			if tt.wantErr {
				assert.NotOk(t, err, tt.wantedErr)
			} else {
				assert.Ok(t, err)
			}
		})
	}
}

func existingPersistentVolumeClient(t *testing.T, expectedListOpts *types.ListOpts) *persistentVolumeClient {
	const (
		projectID            = "test-project-id"
		projectName          = "test-project-name"
		namespaceID          = "test-namespace-id"
		namespace            = "test-namespace"
		clusterID            = "test-cluster-id"
		persistentVolumeName = "test-persistentVolume"
	)
	var (
		persistentVolume = projectModel.PersistentVolume{}
		testClients      = stubs.CreateBackendStubs(t)
	)
	persistentVolume.Name = persistentVolumeName

	persistentVolumeOperationsStub := stubs.CreatePersistentVolumeOperationsStub(t)
	persistentVolumeOperationsStub.DoList = func(opts *types.ListOpts) (*backendClusterClient.PersistentVolumeCollection, error) {
		if !reflect.DeepEqual(expectedListOpts, opts) {
			return nil, fmt.Errorf("Unexpected ListOpts %v", opts)
		}
		return &backendClusterClient.PersistentVolumeCollection{
			Data: []backendClusterClient.PersistentVolume{
				backendClusterClient.PersistentVolume{
					Name: "existing-persistentVolume",
				},
			},
		}, nil
	}
	testClients.ClusterClient.PersistentVolume = persistentVolumeOperationsStub
	clusterClient := simpleClusterClient()
	clusterClient._backendClusterClient = testClients.ClusterClient
	result, err := newPersistentVolumeClient(
		"existing-persistentVolume",
		clusterClient,
		logrus.New().WithFields(logrus.Fields{}),
	)
	assert.Ok(t, err)
	persistentVolumeClientResult := result.(*persistentVolumeClient)
	return persistentVolumeClientResult
}

func notExistingPersistentVolumeClient(t *testing.T, expectedListOpts *types.ListOpts) *persistentVolumeClient {
	const (
		projectID            = "test-project-id"
		projectName          = "test-project-name"
		namespaceID          = "test-namespace-id"
		namespace            = "test-namespace"
		clusterID            = "test-cluster-id"
		persistentVolumeName = "test-persistentVolume"
	)
	var (
		persistentVolume = projectModel.PersistentVolume{}
		testClients      = stubs.CreateBackendStubs(t)
	)
	persistentVolume.Name = persistentVolumeName

	persistentVolumeOperationsStub := stubs.CreatePersistentVolumeOperationsStub(t)
	persistentVolumeOperationsStub.DoList = func(opts *types.ListOpts) (*backendClusterClient.PersistentVolumeCollection, error) {
		if !reflect.DeepEqual(expectedListOpts, opts) {
			return nil, fmt.Errorf("Unexpected ListOpts %v", opts)
		}
		return &backendClusterClient.PersistentVolumeCollection{
			Data: []backendClusterClient.PersistentVolume{},
		}, nil
	}
	persistentVolumeOperationsStub.DoCreate = func(persistentVolume *backendClusterClient.PersistentVolume) (*backendClusterClient.PersistentVolume, error) {
		return persistentVolume, nil
	}
	testClients.ClusterClient.PersistentVolume = persistentVolumeOperationsStub
	clusterClient := simpleClusterClient()
	clusterClient._backendClusterClient = testClients.ClusterClient
	result, err := newPersistentVolumeClient(
		"existing-persistentVolume",
		clusterClient,
		logrus.New().WithFields(logrus.Fields{}),
	)
	assert.Ok(t, err)
	persistentVolumeClientResult := result.(*persistentVolumeClient)
	return persistentVolumeClientResult
}
