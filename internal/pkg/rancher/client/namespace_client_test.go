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
	backendClusterClient "github.com/rancher/types/client/cluster/v3"
	"github.com/sirupsen/logrus"
)

const (
	simpleNamespaceName = "simple-namespace"
	simpleNamespaceID   = "simple-namespace-id"
)

func Test_namespaceClient_Exists(t *testing.T) {
	tests := []struct {
		name      string
		client    *namespaceClient
		wanted    bool
		wantErr   bool
		wantedErr string
	}{
		{
			name: "Existing",
			client: existingNamespaceClient(
				t,
				&types.ListOpts{
					Filters: map[string]interface{}{
						"name": "existing-namespace",
					},
				},
			),
			wanted:  true,
			wantErr: false,
		},
		{
			name: "Not_Existing",
			client: notExistingNamespaceClient(
				t,
				&types.ListOpts{
					Filters: map[string]interface{}{
						"name": "existing-namespace",
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

func Test_namespaceClient_Create(t *testing.T) {
	tests := []struct {
		name      string
		client    *namespaceClient
		wantErr   bool
		wantedErr string
	}{
		{
			name: "Create",
			client: notExistingNamespaceClient(
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

func existingNamespaceClient(t *testing.T, expectedListOpts *types.ListOpts) *namespaceClient {
	const (
		projectID     = "test-project-id"
		projectName   = "test-project-name"
		namespaceID   = "test-namespace-id"
		clusterID     = "test-cluster-id"
		namespaceName = "test-namespace"
	)
	var (
		namespace   = projectModel.Namespace{}
		testClients = stubs.CreateBackendStubs(t)
	)
	namespace.Name = namespaceName

	namespaceOperationsStub := stubs.CreateNamespaceOperationsStub(t)
	namespaceOperationsStub.DoList = func(opts *types.ListOpts) (*backendClusterClient.NamespaceCollection, error) {
		if !reflect.DeepEqual(expectedListOpts, opts) {
			return nil, fmt.Errorf("Unexpected ListOpts %v", opts)
		}
		return &backendClusterClient.NamespaceCollection{
			Data: []backendClusterClient.Namespace{
				backendClusterClient.Namespace{
					Name: "existing-namespace",
				},
			},
		}, nil
	}
	testClients.ClusterClient.Namespace = namespaceOperationsStub
	result, err := newNamespaceClient(
		"existing-namespace",
		nil,
		testClients.ClusterClient,
		logrus.New().WithFields(logrus.Fields{}),
	)
	assert.Ok(t, err)
	namespaceClientResult := result.(*namespaceClient)
	return namespaceClientResult
}

func notExistingNamespaceClient(t *testing.T, expectedListOpts *types.ListOpts) *namespaceClient {
	const (
		projectID     = "test-project-id"
		projectName   = "test-project-name"
		namespaceID   = "test-namespace-id"
		clusterID     = "test-cluster-id"
		namespaceName = "test-namespace"
	)
	var (
		namespace   = projectModel.Namespace{}
		testClients = stubs.CreateBackendStubs(t)
	)
	namespace.Name = namespaceName

	namespaceOperationsStub := stubs.CreateNamespaceOperationsStub(t)
	namespaceOperationsStub.DoList = func(opts *types.ListOpts) (*backendClusterClient.NamespaceCollection, error) {
		if !reflect.DeepEqual(expectedListOpts, opts) {
			return nil, fmt.Errorf("Unexpected ListOpts %v", opts)
		}
		return &backendClusterClient.NamespaceCollection{
			Data: []backendClusterClient.Namespace{},
		}, nil
	}
	namespaceOperationsStub.DoCreate = func(namespace *backendClusterClient.Namespace) (*backendClusterClient.Namespace, error) {
		return namespace, nil
	}
	testClients.ClusterClient.Namespace = namespaceOperationsStub
	result, err := newNamespaceClient(
		"existing-namespace",
		nil,
		testClients.ClusterClient,
		logrus.New().WithFields(logrus.Fields{}),
	)
	assert.Ok(t, err)
	namespaceClientResult := result.(*namespaceClient)
	return namespaceClientResult
}
