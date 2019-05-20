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
			client: testRancherClient(),
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
			tt.client.backendRancherClient = testClients.ManagementClient

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
			client:       testRancherClient(),
			doList:       twoClusters(),
			wantedLength: 2,
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Arrange
			testClients := stubs.CreateBackendStubs(t)
			tt.client.backendRancherClient = testClients.ManagementClient

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

func testRancherClient() *rancherClient {
	return &rancherClient{
		config:         RancherConfig{},
		logger:         logrus.WithFields(logrus.Fields{}),
		clusterClients: make(map[string]ClusterClient),
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
