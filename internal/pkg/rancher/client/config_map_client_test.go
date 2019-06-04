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
	projectModel "github.com/bitgrip/cattlectl/internal/pkg/rancher/project/model"
	"github.com/bitgrip/cattlectl/internal/pkg/rancher/stubs"
	"github.com/rancher/norman/types"
	backendProjectClient "github.com/rancher/types/client/project/v3"
	"github.com/sirupsen/logrus"
)

const (
	simpleConfigMapName = "simple-config-map"
)

func Test_configMapClient_Exists(t *testing.T) {
	tests := []struct {
		name      string
		client    *configMapClient
		wanted    bool
		wantErr   bool
		wantedErr string
	}{
		{
			name: "Existing",
			client: existingConfigMapClient(
				t,
				&types.ListOpts{
					Filters: map[string]interface{}{
						"name":        "existing-configMap",
						"namespaceId": "test-namespace-id",
					},
				},
			),
			wanted:  true,
			wantErr: false,
		},
		{
			name: "Not_Existing",
			client: notExistingConfigMapClient(
				t,
				&types.ListOpts{
					Filters: map[string]interface{}{
						"name":        "existing-configMap",
						"namespaceId": "test-namespace-id",
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

func Test_configMapClient_Create(t *testing.T) {
	tests := []struct {
		name      string
		client    *configMapClient
		wantErr   bool
		wantedErr string
	}{
		{
			name: "Create",
			client: notExistingConfigMapClient(
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

func existingConfigMapClient(t *testing.T, expectedListOpts *types.ListOpts) *configMapClient {
	const (
		projectID     = "test-project-id"
		projectName   = "test-project-name"
		namespaceID   = "test-namespace-id"
		namespace     = "test-namespace"
		clusterID     = "test-cluster-id"
		configMapName = "test-configMap"
	)
	var (
		configMap   = projectModel.ConfigMap{}
		testClients = stubs.CreateBackendStubs(t)
	)
	configMap.Name = configMapName

	configMapOperationsStub := stubs.CreateConfigMapOperationsStub(t)
	configMapOperationsStub.DoList = func(opts *types.ListOpts) (*backendProjectClient.ConfigMapCollection, error) {
		if !reflect.DeepEqual(expectedListOpts, opts) {
			return nil, fmt.Errorf("Unexpected ListOpts %v", opts)
		}
		return &backendProjectClient.ConfigMapCollection{
			Data: []backendProjectClient.ConfigMap{
				backendProjectClient.ConfigMap{
					Name:        "existing-configMap",
					NamespaceId: "test-namespace-id",
				},
			},
		}, nil
	}
	testClients.ProjectClient.ConfigMap = configMapOperationsStub
	result, err := newConfigMapClient(
		"existing-configMap",
		"test-namespace",
		&projectClient{
			resourceClient: resourceClient{
				name: projectName,
				id:   projectID,
			},
		},
		testClients.ProjectClient,
		logrus.New().WithFields(logrus.Fields{}),
	)
	assert.Ok(t, err)
	configMapClientResult := result.(*configMapClient)
	configMapClientResult.namespaceID = "test-namespace-id"
	return configMapClientResult
}

func notExistingConfigMapClient(t *testing.T, expectedListOpts *types.ListOpts) *configMapClient {
	const (
		projectID     = "test-project-id"
		projectName   = "test-project-name"
		namespaceID   = "test-namespace-id"
		namespace     = "test-namespace"
		clusterID     = "test-cluster-id"
		configMapName = "test-configMap"
	)
	var (
		configMap   = projectModel.ConfigMap{}
		testClients = stubs.CreateBackendStubs(t)
	)
	configMap.Name = configMapName

	configMapOperationsStub := stubs.CreateConfigMapOperationsStub(t)
	configMapOperationsStub.DoList = func(opts *types.ListOpts) (*backendProjectClient.ConfigMapCollection, error) {
		if !reflect.DeepEqual(expectedListOpts, opts) {
			return nil, fmt.Errorf("Unexpected ListOpts %v", opts)
		}
		return &backendProjectClient.ConfigMapCollection{
			Data: []backendProjectClient.ConfigMap{},
		}, nil
	}
	configMapOperationsStub.DoCreate = func(configMap *backendProjectClient.ConfigMap) (*backendProjectClient.ConfigMap, error) {
		return configMap, nil
	}
	testClients.ProjectClient.ConfigMap = configMapOperationsStub
	result, err := newConfigMapClient(
		"existing-configMap",
		"test-namespace",
		&projectClient{
			resourceClient: resourceClient{
				name: projectName,
				id:   projectID,
			},
		},
		testClients.ProjectClient,
		logrus.New().WithFields(logrus.Fields{}),
	)
	assert.Ok(t, err)
	configMapClientResult := result.(*configMapClient)
	configMapClientResult.namespaceID = "test-namespace-id"
	return configMapClientResult
}
