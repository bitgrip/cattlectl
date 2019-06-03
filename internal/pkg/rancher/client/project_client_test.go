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
	"github.com/sirupsen/logrus"
)

const (
	simpleProjectName = "simple-project"
	simpleProjectID   = "simple-project-id"
)

func Test_projectClient_Namespace(t *testing.T) {
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
			tt.client.backendRancherClient = testClients.ManagementClient
			tt.client.backendClusterClient = testClients.ClusterClient
			tt.client.backendProjectClient = testClients.ProjectClient
			clusterClient := simpleClusterClient()
			clusterClient.backendRancherClient = testClients.ManagementClient
			clusterClient.backendClusterClient = testClients.ClusterClient
			tt.client.clusterClient = clusterClient

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

func Test_projectClient_Namespaces(t *testing.T) {
	tests := []struct {
		name            string
		client          *projectClient
		foundNamespaces []string
		wantedLength    int
		wantErr         bool
		wantedErr       string
	}{
		{
			name:            "success",
			client:          simpleProjectClient(),
			foundNamespaces: []string{simpleNamespaceName + "1", simpleNamespaceName + "2"},
			wantedLength:    2,
			wantErr:         false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Arrange
			testClients := stubs.CreateBackendStubs(t)
			tt.client.backendRancherClient = testClients.ManagementClient
			tt.client.backendClusterClient = testClients.ClusterClient
			tt.client.backendProjectClient = testClients.ProjectClient
			clusterClient := simpleClusterClient()
			clusterClient.backendRancherClient = testClients.ManagementClient
			clusterClient.backendClusterClient = testClients.ClusterClient
			tt.client.clusterClient = clusterClient

			namespaceOperationStub := stubs.CreateNamespaceOperationsStub(t)
			namespaceOperationStub.DoList = foundNamespaces(tt.foundNamespaces, simpleProjectID)
			testClients.ClusterClient.Namespace = namespaceOperationStub

			got, err := tt.client.Namespaces()
			if tt.wantErr {
				assert.NotOk(t, err, tt.wantedErr)
			} else {
				assert.Ok(t, err)
				assert.Equals(t, tt.wantedLength, len(got))
			}
		})
	}
}

func Test_projectClient_Certificate(t *testing.T) {
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
				name: simpleCertificateName,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Arrange
			testClients := stubs.CreateBackendStubs(t)
			tt.client.backendRancherClient = testClients.ManagementClient
			tt.client.backendClusterClient = testClients.ClusterClient
			tt.client.backendProjectClient = testClients.ProjectClient
			clusterClient := simpleClusterClient()
			clusterClient.backendRancherClient = testClients.ManagementClient
			clusterClient.backendClusterClient = testClients.ClusterClient
			tt.client.clusterClient = clusterClient

			got, err := tt.client.Certificate(tt.args.name)
			if tt.wantErr {
				assert.NotOk(t, err, tt.wantedErr)
			} else {
				assert.Ok(t, err)
				gotName, err := got.Name()
				assert.Ok(t, err)
				assert.Equals(t, tt.args.name, gotName)
			}
			got2, err := tt.client.Certificate(tt.args.name)
			if tt.wantErr {
				assert.NotOk(t, err, tt.wantedErr)
			} else {
				assert.Ok(t, err)
				assert.Assert(t, got == got2, "second call shout returne the same object from cache")
			}
		})
	}
}

func simpleProjectClient() *projectClient {
	return &projectClient{
		resourceClient: resourceClient{
			name:   simpleProjectName,
			id:     simpleProjectID,
			logger: logrus.WithFields(logrus.Fields{}),
		},
		certificateClients: make(map[string]CertificateClient),
	}
}
