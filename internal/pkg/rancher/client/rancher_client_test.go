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
	"testing"

	"github.com/bitgrip/cattlectl/internal/pkg/assert"
	"github.com/bitgrip/cattlectl/internal/pkg/rancher/stubs"
	"github.com/rancher/norman/types"
	backendRancherClient "github.com/rancher/types/client/management/v3"
	managementClient "github.com/rancher/types/client/management/v3"
	"github.com/sirupsen/logrus"
)

func Test_rancherClient_Cluster(t *testing.T) {

	type args struct {
		clusterName string
	}
	tests := []struct {
		name      string
		client    *rancherClient
		args      args
		wantErr   bool
		wantedErr string
	}{
		{
			name:   "success",
			client: simpleRancherClient(),
			args: args{
				clusterName: "test-cluster",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Arrange
			testClients := stubs.CreateBackendStubs(t)
			tt.client._backendRancherClient = testClients.ManagementClient

			got, err := tt.client.Cluster(tt.args.clusterName)
			if tt.wantErr {
				assert.NotOk(t, err, tt.wantedErr)
			} else {
				assert.Ok(t, err)
				gotName, err := got.Name()
				assert.Ok(t, err)
				assert.Equals(t, tt.args.clusterName, gotName)
			}
		})
	}
}

func Test_rancherClient_Clusters(t *testing.T) {
	tests := []struct {
		name         string
		client       *rancherClient
		doList       func(opts *types.ListOpts) (*managementClient.ClusterCollection, error)
		wantedLength int
		wantErr      bool
		wantedErr    string
	}{
		{
			name:         "success",
			client:       simpleRancherClient(),
			doList:       twoClusters(),
			wantedLength: 2,
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Arrange
			testClients := stubs.CreateBackendStubs(t)
			tt.client._backendRancherClient = testClients.ManagementClient

			clusterOperationsStub := stubs.CreateClusterOperationsStub(t)
			clusterOperationsStub.DoList = tt.doList
			testClients.ManagementClient.Cluster = clusterOperationsStub

			got, err := tt.client.Clusters()
			if tt.wantErr {
				assert.NotOk(t, err, tt.wantedErr)
			} else {
				assert.Ok(t, err)
				assert.Equals(t, tt.wantedLength, len(got))
			}
		})
	}
}

func Test_rancherClient_Catalog(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name      string
		client    *projectClient
		args      args
		wantErr   bool
		wantedErr string
	}{
		{
			name:   simpleClusterName,
			client: simpleProjectClient(),
			args: args{
				name: simpleNamespaceName,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Arrange
			testClients := stubs.CreateBackendStubs(t)
			rancherClient := simpleRancherClient()
			clusterClient := simpleClusterClient()
			rancherClient._backendRancherClient = testClients.ManagementClient
			clusterClient.rancherClient = rancherClient
			clusterClient._backendClusterClient = testClients.ClusterClient
			tt.client.clusterClient = clusterClient
			tt.client._backendProjectClient = testClients.ProjectClient

			got, err := tt.client.Namespace(tt.args.name)
			if tt.wantErr {
				assert.NotOk(t, err, tt.wantedErr)
			} else {
				assert.Ok(t, err)
				gotName, err := got.Name()
				assert.Ok(t, err)
				assert.Equals(t, tt.args.name, gotName)
			}
			got2, err := tt.client.Namespace(tt.args.name)
			if tt.wantErr {
				assert.NotOk(t, err, tt.wantedErr)
			} else {
				assert.Ok(t, err)
				assert.Assert(t, got == got2, "second call shout returne the same object from cache")
			}
		})
	}
}

func Test_rancherClient_Catalogs(t *testing.T) {
	tests := []struct {
		name          string
		client        *rancherClient
		foundCatalogs []string
		wantedLength  int
		wantErr       bool
		wantedErr     string
	}{
		{
			name:          "success",
			client:        simpleRancherClient(),
			foundCatalogs: []string{simpleCatalogName + "1", simpleCatalogName + "2"},
			wantedLength:  2,
			wantErr:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Arrange
			testClients := stubs.CreateBackendStubs(t)
			tt.client._backendRancherClient = testClients.ManagementClient

			catalogOperationStub := stubs.CreateRancherCatalogOperationsStub(t)
			catalogOperationStub.DoList = foundCatalogs(tt.foundCatalogs)
			testClients.ManagementClient.Catalog = catalogOperationStub

			got, err := tt.client.Catalogs()
			if tt.wantErr {
				assert.NotOk(t, err, tt.wantedErr)
			} else {
				assert.Ok(t, err)
				assert.Equals(t, tt.wantedLength, len(got))
			}
		})
	}
}

func simpleRancherClient() *rancherClient {
	return &rancherClient{
		config:         RancherConfig{},
		logger:         logrus.WithFields(logrus.Fields{}),
		clusterClients: make(map[string]ClusterClient),
		catalogClients: make(map[string]CatalogClient),
	}
}

func twoClusters() func(opts *types.ListOpts) (*managementClient.ClusterCollection, error) {
	return func(opts *types.ListOpts) (*managementClient.ClusterCollection, error) {
		return &managementClient.ClusterCollection{
			Data: []managementClient.Cluster{
				managementClient.Cluster{},
				managementClient.Cluster{},
			},
		}, nil
	}
}

func foundCatalogs(names []string) func(opts *types.ListOpts) (*backendRancherClient.CatalogCollection, error) {
	data := make([]backendRancherClient.Catalog, 0)
	for _, name := range names {
		data = append(data, backendRancherClient.Catalog{Name: name})
	}
	return func(opts *types.ListOpts) (*backendRancherClient.CatalogCollection, error) {
		return &backendRancherClient.CatalogCollection{
			Data: data,
		}, nil
	}
}
