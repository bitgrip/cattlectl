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
	backendRancherClient "github.com/rancher/types/client/management/v3"
	managementClient "github.com/rancher/types/client/management/v3"
	backendProjectClient "github.com/rancher/types/client/project/v3"
	"github.com/sirupsen/logrus"
)

const (
	simpleProjectName = "simple-project"
	simpleProjectID   = "simple-project-id"
)

func Test_projectClient_Exists(t *testing.T) {
	tests := []struct {
		name      string
		client    *projectClient
		wanted    bool
		wantErr   bool
		wantedErr string
	}{
		{
			name: "Existing",
			client: existingProjectClient(
				t,
				simpleProjectName,
				simpleProjectID,
				simpleClusterID,
			),
			wanted:  true,
			wantErr: false,
		},
		{
			name: "Not_Existing",
			client: notExistingProjectClient(
				t,
				simpleProjectName,
				simpleProjectID,
				simpleClusterID,
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

func Test_projectClient_Create(t *testing.T) {
	tests := []struct {
		name            string
		client          *projectClient
		wantedProjectID string
		wantErr         bool
		wantedErr       string
	}{
		{
			name: "Create",
			client: notExistingProjectClient(
				t,
				simpleProjectName,
				simpleProjectID,
				simpleClusterID,
			),
			wantedProjectID: simpleProjectID,
			wantErr:         false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.client.Create()
			if tt.wantErr {
				assert.NotOk(t, err, tt.wantedErr)
			} else {
				assert.Ok(t, err)
				projectID, err := tt.client.ID()
				assert.Ok(t, err)
				assert.Equals(t, tt.wantedProjectID, projectID)
			}
		})
	}
}

func Test_projectClient_Upgrade(t *testing.T) {
	tests := []struct {
		name            string
		client          *projectClient
		wantedProjectID string
		wantErr         bool
		wantedErr       string
	}{
		{
			name: "Upgrade",
			client: existingProjectClient(
				t,
				simpleProjectName,
				simpleProjectID,
				simpleClusterID,
			),
			wantedProjectID: simpleProjectID,
			wantErr:         false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.client.Upgrade()
			if tt.wantErr {
				assert.NotOk(t, err, tt.wantedErr)
			} else {
				assert.Ok(t, err)
				projectID, err := tt.client.ID()
				assert.Ok(t, err)
				assert.Equals(t, tt.wantedProjectID, projectID)
			}
		})
	}
}

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
			rancherClient := simpleRancherClient()
			clusterClient := simpleClusterClient()
			rancherClient._backendRancherClient = testClients.ManagementClient
			clusterClient.rancherClient = rancherClient
			clusterClient._backendClusterClient = testClients.ClusterClient
			tt.client.clusterClient = clusterClient
			tt.client._backendProjectClient = testClients.ProjectClient

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

func Test_projectClient_GlobalCertificate(t *testing.T) {
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
			rancherClient := simpleRancherClient()
			clusterClient := simpleClusterClient()
			rancherClient._backendRancherClient = testClients.ManagementClient
			clusterClient.rancherClient = rancherClient
			clusterClient._backendClusterClient = testClients.ClusterClient
			tt.client.clusterClient = clusterClient
			tt.client._backendProjectClient = testClients.ProjectClient

			got, err := tt.client.GlobalCertificate(tt.args.name)
			if tt.wantErr {
				assert.NotOk(t, err, tt.wantedErr)
			} else {
				assert.Ok(t, err)
				gotName, err := got.Name()
				assert.Ok(t, err)
				assert.Equals(t, tt.args.name, gotName)
			}
			got2, err := tt.client.GlobalCertificate(tt.args.name)
			if tt.wantErr {
				assert.NotOk(t, err, tt.wantedErr)
			} else {
				assert.Ok(t, err)
				assert.Assert(t, got == got2, "second call shout returne the same object from cache")
			}
		})
	}
}

func Test_projectClient_GlobalCertificates(t *testing.T) {
	tests := []struct {
		name              string
		client            *projectClient
		foundCertificates []string
		wantedLength      int
		wantErr           bool
		wantedErr         string
	}{
		{
			name:              "success",
			client:            simpleProjectClient(),
			foundCertificates: []string{simpleCertificateName + "1", simpleCertificateName + "2"},
			wantedLength:      2,
			wantErr:           false,
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

			certificateOperationStub := stubs.CreateCertificateOperationsStub(t)
			certificateOperationStub.DoList = foundGlobalCertificates(tt.foundCertificates, simpleProjectID)
			testClients.ProjectClient.Certificate = certificateOperationStub

			got, err := tt.client.GlobalCertificates()
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
		name      string
		namespace string
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
				name:      simpleCertificateName,
				namespace: simpleNamespaceName,
			},
			wantErr: false,
		},
		{
			name:   simpleClusterName,
			client: simpleProjectClient(),
			args: args{
				name:      simpleCertificateName,
				namespace: "",
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

			got, err := tt.client.Certificate(tt.args.name, tt.args.namespace)
			if tt.wantErr {
				assert.NotOk(t, err, tt.wantedErr)
			} else {
				assert.Ok(t, err)
				gotName, err := got.Name()
				assert.Ok(t, err)
				assert.Equals(t, tt.args.name, gotName)
			}
			got2, err := tt.client.Certificate(tt.args.name, tt.args.namespace)
			if tt.wantErr {
				assert.NotOk(t, err, tt.wantedErr)
			} else {
				assert.Ok(t, err)
				assert.Assert(t, got == got2, "second call shout returne the same object from cache")
			}
		})
	}
}

func Test_projectClient_Certificates(t *testing.T) {
	tests := []struct {
		name              string
		client            *projectClient
		foundCertificates []string
		namespaceName     string
		namespaceID       string
		wantedLength      int
		wantErr           bool
		wantedErr         string
	}{
		{
			name:              "success",
			client:            simpleProjectClient(),
			foundCertificates: []string{simpleCertificateName + "1", simpleCertificateName + "2"},
			namespaceName:     simpleNamespaceName,
			namespaceID:       simpleNamespaceID,
			wantedLength:      2,
			wantErr:           false,
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

			certificateOperationStub := stubs.CreateNamespacedCertificateOperationsStub(t)
			certificateOperationStub.DoList = foundCertificates(tt.foundCertificates, tt.namespaceID, simpleProjectID)
			testClients.ProjectClient.NamespacedCertificate = certificateOperationStub

			namespaceOperationsStub := stubs.CreateNamespaceOperationsStub(t)
			namespaceOperationsStub.DoList = foundNamespace(tt.namespaceName, tt.namespaceID, simpleProjectID)
			testClients.ClusterClient.Namespace = namespaceOperationsStub

			got, err := tt.client.Certificates(tt.namespaceName)
			if tt.wantErr {
				assert.NotOk(t, err, tt.wantedErr)
			} else {
				assert.Ok(t, err)
				assert.Equals(t, tt.wantedLength, len(got))
			}
		})
	}
}

func Test_projectClient_ConfigMap(t *testing.T) {
	type args struct {
		name      string
		namespace string
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
				name:      simpleConfigMapName,
				namespace: simpleNamespaceName,
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

			got, err := tt.client.ConfigMap(tt.args.name, tt.args.namespace)
			if tt.wantErr {
				assert.NotOk(t, err, tt.wantedErr)
			} else {
				assert.Ok(t, err)
				gotName, err := got.Name()
				assert.Ok(t, err)
				assert.Equals(t, tt.args.name, gotName)
			}
			got2, err := tt.client.ConfigMap(tt.args.name, tt.args.namespace)
			if tt.wantErr {
				assert.NotOk(t, err, tt.wantedErr)
			} else {
				assert.Ok(t, err)
				assert.Assert(t, got == got2, "second call shout returne the same object from cache")
			}
		})
	}
}

func Test_projectClient_ConfigMaps(t *testing.T) {
	tests := []struct {
		name            string
		client          *projectClient
		foundConfigMaps []string
		namespaceName   string
		namespaceID     string
		wantedLength    int
		wantErr         bool
		wantedErr       string
	}{
		{
			name:            "success",
			client:          simpleProjectClient(),
			foundConfigMaps: []string{simpleConfigMapName + "1", simpleConfigMapName + "2"},
			namespaceName:   simpleNamespaceName,
			namespaceID:     simpleNamespaceID,
			wantedLength:    2,
			wantErr:         false,
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

			configMapOperationStub := stubs.CreateConfigMapOperationsStub(t)
			configMapOperationStub.DoList = foundConfigMaps(tt.foundConfigMaps, tt.namespaceID, simpleProjectID)
			testClients.ProjectClient.ConfigMap = configMapOperationStub

			namespaceOperationsStub := stubs.CreateNamespaceOperationsStub(t)
			namespaceOperationsStub.DoList = foundNamespace(tt.namespaceName, tt.namespaceID, simpleProjectID)
			testClients.ClusterClient.Namespace = namespaceOperationsStub

			got, err := tt.client.ConfigMaps(tt.namespaceName)
			if tt.wantErr {
				assert.NotOk(t, err, tt.wantedErr)
			} else {
				assert.Ok(t, err)
				assert.Equals(t, tt.wantedLength, len(got))
			}
		})
	}
}

func Test_projectClient_GlobalDockerCredential(t *testing.T) {
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
			name:   simpleDockerCredentialName,
			client: simpleProjectClient(),
			args: args{
				name: simpleDockerCredentialName,
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

			got, err := tt.client.GlobalDockerCredential(tt.args.name)
			if tt.wantErr {
				assert.NotOk(t, err, tt.wantedErr)
			} else {
				assert.Ok(t, err)
				gotName, err := got.Name()
				assert.Ok(t, err)
				assert.Equals(t, tt.args.name, gotName)
			}
			got2, err := tt.client.GlobalDockerCredential(tt.args.name)
			if tt.wantErr {
				assert.NotOk(t, err, tt.wantedErr)
			} else {
				assert.Ok(t, err)
				assert.Assert(t, got == got2, "second call shout returne the same object from cache")
			}
		})
	}
}

func Test_projectClient_GlobalDockerCredentials(t *testing.T) {
	tests := []struct {
		name                   string
		client                 *projectClient
		foundDockerCredentials []string
		wantedLength           int
		wantErr                bool
		wantedErr              string
	}{
		{
			name:                   "success",
			client:                 simpleProjectClient(),
			foundDockerCredentials: []string{simpleDockerCredentialName + "1", simpleDockerCredentialName + "2"},
			wantedLength:           2,
			wantErr:                false,
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

			dockerCredentialOperationStub := stubs.CreateDockerCredentialOperationsStub(t)
			dockerCredentialOperationStub.DoList = foundGlobalDockerCredentials(tt.foundDockerCredentials, simpleProjectID)
			testClients.ProjectClient.DockerCredential = dockerCredentialOperationStub

			got, err := tt.client.GlobalDockerCredentials()
			if tt.wantErr {
				assert.NotOk(t, err, tt.wantedErr)
			} else {
				assert.Ok(t, err)
				assert.Equals(t, tt.wantedLength, len(got))
			}
		})
	}
}

func Test_projectClient_DockerCredential(t *testing.T) {
	type args struct {
		name      string
		namespace string
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
				name:      simpleDockerCredentialName,
				namespace: simpleNamespaceName,
			},
			wantErr: false,
		},
		{
			name:   simpleClusterName,
			client: simpleProjectClient(),
			args: args{
				name:      simpleDockerCredentialName,
				namespace: "",
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

			got, err := tt.client.DockerCredential(tt.args.name, tt.args.namespace)
			if tt.wantErr {
				assert.NotOk(t, err, tt.wantedErr)
			} else {
				assert.Ok(t, err)
				gotName, err := got.Name()
				assert.Ok(t, err)
				assert.Equals(t, tt.args.name, gotName)
			}
			got2, err := tt.client.DockerCredential(tt.args.name, tt.args.namespace)
			if tt.wantErr {
				assert.NotOk(t, err, tt.wantedErr)
			} else {
				assert.Ok(t, err)
				assert.Assert(t, got == got2, "second call shout returne the same object from cache")
			}
		})
	}
}

func Test_projectClient_DockerCredentials(t *testing.T) {
	tests := []struct {
		name                   string
		client                 *projectClient
		foundDockerCredentials []string
		namespaceName          string
		namespaceID            string
		wantedLength           int
		wantErr                bool
		wantedErr              string
	}{
		{
			name:                   "success",
			client:                 simpleProjectClient(),
			foundDockerCredentials: []string{simpleDockerCredentialName + "1", simpleDockerCredentialName + "2"},
			namespaceName:          simpleNamespaceName,
			namespaceID:            simpleNamespaceID,
			wantedLength:           2,
			wantErr:                false,
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

			dockerCredentialOperationStub := stubs.CreateNamespacedDockerCredentialOperationsStub(t)
			dockerCredentialOperationStub.DoList = foundDockerCredentials(tt.foundDockerCredentials, tt.namespaceID, simpleProjectID)
			testClients.ProjectClient.NamespacedDockerCredential = dockerCredentialOperationStub

			namespaceOperationsStub := stubs.CreateNamespaceOperationsStub(t)
			namespaceOperationsStub.DoList = foundNamespace(tt.namespaceName, tt.namespaceID, simpleProjectID)
			testClients.ClusterClient.Namespace = namespaceOperationsStub

			got, err := tt.client.DockerCredentials(tt.namespaceName)
			if tt.wantErr {
				assert.NotOk(t, err, tt.wantedErr)
			} else {
				assert.Ok(t, err)
				assert.Equals(t, tt.wantedLength, len(got))
			}
		})
	}
}

func Test_projectClient_GlobalSecret(t *testing.T) {
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
			name:   simpleSecretName,
			client: simpleProjectClient(),
			args: args{
				name: simpleSecretName,
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

			got, err := tt.client.GlobalSecret(tt.args.name)
			if tt.wantErr {
				assert.NotOk(t, err, tt.wantedErr)
			} else {
				assert.Ok(t, err)
				gotName, err := got.Name()
				assert.Ok(t, err)
				assert.Equals(t, tt.args.name, gotName)
			}
			got2, err := tt.client.GlobalSecret(tt.args.name)
			if tt.wantErr {
				assert.NotOk(t, err, tt.wantedErr)
			} else {
				assert.Ok(t, err)
				assert.Assert(t, got == got2, "second call shout returne the same object from cache")
			}
		})
	}
}

func Test_projectClient_GlobalSecrets(t *testing.T) {
	tests := []struct {
		name         string
		client       *projectClient
		foundSecrets []string
		wantedLength int
		wantErr      bool
		wantedErr    string
	}{
		{
			name:         "success",
			client:       simpleProjectClient(),
			foundSecrets: []string{simpleSecretName + "1", simpleSecretName + "2"},
			wantedLength: 2,
			wantErr:      false,
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

			secretOperationsStub := stubs.CreateSecretOperationsStub(t)
			secretOperationsStub.DoList = foundGlobalSecrets(tt.foundSecrets, simpleProjectID)
			testClients.ProjectClient.Secret = secretOperationsStub

			got, err := tt.client.GlobalSecrets()
			if tt.wantErr {
				assert.NotOk(t, err, tt.wantedErr)
			} else {
				assert.Ok(t, err)
				assert.Equals(t, tt.wantedLength, len(got))
			}
		})
	}
}

func Test_projectClient_Secret(t *testing.T) {
	type args struct {
		name      string
		namespace string
	}
	tests := []struct {
		name      string
		client    *projectClient
		args      args
		wantErr   bool
		wantedErr string
	}{
		{
			name:   simpleSecretName,
			client: simpleProjectClient(),
			args: args{
				name:      simpleSecretName,
				namespace: simpleNamespaceName,
			},
			wantErr: false,
		},
		{
			name:   simpleClusterName,
			client: simpleProjectClient(),
			args: args{
				name:      simpleSecretName,
				namespace: "",
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

			got, err := tt.client.Secret(tt.args.name, tt.args.namespace)
			if tt.wantErr {
				assert.NotOk(t, err, tt.wantedErr)
			} else {
				assert.Ok(t, err)
				gotName, err := got.Name()
				assert.Ok(t, err)
				assert.Equals(t, tt.args.name, gotName)
			}
			got2, err := tt.client.Secret(tt.args.name, tt.args.namespace)
			if tt.wantErr {
				assert.NotOk(t, err, tt.wantedErr)
			} else {
				assert.Ok(t, err)
				assert.Assert(t, got == got2, "second call shout returne the same object from cache")
			}
		})
	}
}

func Test_projectClient_Secrets(t *testing.T) {
	tests := []struct {
		name          string
		client        *projectClient
		foundSecrets  []string
		namespaceName string
		namespaceID   string
		wantedLength  int
		wantErr       bool
		wantedErr     string
	}{
		{
			name:          "success",
			client:        simpleProjectClient(),
			foundSecrets:  []string{simpleSecretName + "1", simpleSecretName + "2"},
			namespaceName: simpleNamespaceName,
			namespaceID:   simpleNamespaceID,
			wantedLength:  2,
			wantErr:       false,
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

			secretOperationsStub := stubs.CreateNamespacedSecretOperationsStub(t)
			secretOperationsStub.DoList = foundDockerSecrets(tt.foundSecrets, tt.namespaceID, simpleProjectID)
			testClients.ProjectClient.NamespacedSecret = secretOperationsStub

			namespaceOperationsStub := stubs.CreateNamespaceOperationsStub(t)
			namespaceOperationsStub.DoList = foundNamespace(tt.namespaceName, tt.namespaceID, simpleProjectID)
			testClients.ClusterClient.Namespace = namespaceOperationsStub

			got, err := tt.client.Secrets(tt.namespaceName)
			if tt.wantErr {
				assert.NotOk(t, err, tt.wantedErr)
			} else {
				assert.Ok(t, err)
				assert.Equals(t, tt.wantedLength, len(got))
			}
		})
	}
}

func Test_projectClient_App(t *testing.T) {
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
			name:   simpleSecretName,
			client: simpleProjectClient(),
			args: args{
				name: simpleAppName,
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

			got, err := tt.client.App(tt.args.name)
			if tt.wantErr {
				assert.NotOk(t, err, tt.wantedErr)
			} else {
				assert.Ok(t, err)
				gotName, err := got.Name()
				assert.Ok(t, err)
				assert.Equals(t, tt.args.name, gotName)
			}
			got2, err := tt.client.App(tt.args.name)
			if tt.wantErr {
				assert.NotOk(t, err, tt.wantedErr)
			} else {
				assert.Ok(t, err)
				assert.Assert(t, got == got2, "second call shout returne the same object from cache")
			}
		})
	}
}

func Test_projectClient_Apps(t *testing.T) {
	tests := []struct {
		name         string
		client       *projectClient
		foundApps    []string
		wantedLength int
		wantErr      bool
		wantedErr    string
	}{
		{
			name:         "success",
			client:       simpleProjectClient(),
			foundApps:    []string{simpleAppName + "1", simpleAppName + "2"},
			wantedLength: 2,
			wantErr:      false,
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

			appOperationsStub := stubs.CreateAppOperationsStub(t)
			appOperationsStub.DoList = foundApps(tt.foundApps, simpleProjectID)
			testClients.ProjectClient.App = appOperationsStub

			got, err := tt.client.Apps()
			if tt.wantErr {
				assert.NotOk(t, err, tt.wantedErr)
			} else {
				assert.Ok(t, err)
				assert.Equals(t, tt.wantedLength, len(got))
			}
		})
	}
}

func Test_projectClient_Job(t *testing.T) {
	type args struct {
		name      string
		namespace string
	}
	tests := []struct {
		name      string
		client    *projectClient
		args      args
		wantErr   bool
		wantedErr string
	}{
		{
			name:   simpleJobName,
			client: simpleProjectClient(),
			args: args{
				name:      simpleJobName,
				namespace: simpleNamespaceName,
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

			got, err := tt.client.Job(tt.args.name, tt.args.namespace)
			if tt.wantErr {
				assert.NotOk(t, err, tt.wantedErr)
			} else {
				assert.Ok(t, err)
				gotName, err := got.Name()
				assert.Ok(t, err)
				assert.Equals(t, tt.args.name, gotName)
			}
			got2, err := tt.client.Job(tt.args.name, tt.args.namespace)
			if tt.wantErr {
				assert.NotOk(t, err, tt.wantedErr)
			} else {
				assert.Ok(t, err)
				assert.Assert(t, got == got2, "second call shout returne the same object from cache")
			}
		})
	}
}

func Test_projectClient_Jobs(t *testing.T) {
	tests := []struct {
		name          string
		client        *projectClient
		foundJobs     []string
		namespaceName string
		namespaceID   string
		wantedLength  int
		wantErr       bool
		wantedErr     string
	}{
		{
			name:          "success",
			client:        simpleProjectClient(),
			foundJobs:     []string{simpleJobName + "1", simpleJobName + "2"},
			namespaceName: simpleNamespaceName,
			namespaceID:   simpleNamespaceID,
			wantedLength:  2,
			wantErr:       false,
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

			jobOperationsStub := stubs.CreateJobOperationsStub(t)
			jobOperationsStub.DoList = foundJobs(tt.foundJobs, tt.namespaceID, simpleProjectID)
			testClients.ProjectClient.Job = jobOperationsStub

			namespaceOperationsStub := stubs.CreateNamespaceOperationsStub(t)
			namespaceOperationsStub.DoList = foundNamespace(tt.namespaceName, tt.namespaceID, simpleProjectID)
			testClients.ClusterClient.Namespace = namespaceOperationsStub

			got, err := tt.client.Jobs(tt.namespaceName)
			if tt.wantErr {
				assert.NotOk(t, err, tt.wantedErr)
			} else {
				assert.Ok(t, err)
				assert.Equals(t, tt.wantedLength, len(got))
			}
		})
	}
}

func Test_projectClient_CronJob(t *testing.T) {
	type args struct {
		name      string
		namespace string
	}
	tests := []struct {
		name      string
		client    *projectClient
		args      args
		wantErr   bool
		wantedErr string
	}{
		{
			name:   simpleCronJobName,
			client: simpleProjectClient(),
			args: args{
				name:      simpleCronJobName,
				namespace: simpleNamespaceName,
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

			got, err := tt.client.CronJob(tt.args.name, tt.args.namespace)
			if tt.wantErr {
				assert.NotOk(t, err, tt.wantedErr)
			} else {
				assert.Ok(t, err)
				gotName, err := got.Name()
				assert.Ok(t, err)
				assert.Equals(t, tt.args.name, gotName)
			}
			got2, err := tt.client.CronJob(tt.args.name, tt.args.namespace)
			if tt.wantErr {
				assert.NotOk(t, err, tt.wantedErr)
			} else {
				assert.Ok(t, err)
				assert.Assert(t, got == got2, "second call shout returne the same object from cache")
			}
		})
	}
}

func Test_projectClient_CronJobs(t *testing.T) {
	tests := []struct {
		name          string
		client        *projectClient
		foundCronJobs []string
		namespaceName string
		namespaceID   string
		wantedLength  int
		wantErr       bool
		wantedErr     string
	}{
		{
			name:          "success",
			client:        simpleProjectClient(),
			foundCronJobs: []string{simpleCronJobName + "1", simpleCronJobName + "2"},
			namespaceName: simpleNamespaceName,
			namespaceID:   simpleNamespaceID,
			wantedLength:  2,
			wantErr:       false,
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

			cronJobOperationStub := stubs.CreateCronJobOperationsStub(t)
			cronJobOperationStub.DoList = foundCronJobs(tt.foundCronJobs, tt.namespaceID, simpleProjectID)
			testClients.ProjectClient.CronJob = cronJobOperationStub

			namespaceOperationsStub := stubs.CreateNamespaceOperationsStub(t)
			namespaceOperationsStub.DoList = foundNamespace(tt.namespaceName, tt.namespaceID, simpleProjectID)
			testClients.ClusterClient.Namespace = namespaceOperationsStub

			got, err := tt.client.CronJobs(tt.namespaceName)
			if tt.wantErr {
				assert.NotOk(t, err, tt.wantedErr)
			} else {
				assert.Ok(t, err)
				assert.Equals(t, tt.wantedLength, len(got))
			}
		})
	}
}

func Test_projectClient_Deployment(t *testing.T) {
	type args struct {
		name      string
		namespace string
	}
	tests := []struct {
		name      string
		client    *projectClient
		args      args
		wantErr   bool
		wantedErr string
	}{
		{
			name:   simpleDeploymentName,
			client: simpleProjectClient(),
			args: args{
				name:      simpleDeploymentName,
				namespace: simpleNamespaceName,
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

			got, err := tt.client.Deployment(tt.args.name, tt.args.namespace)
			if tt.wantErr {
				assert.NotOk(t, err, tt.wantedErr)
			} else {
				assert.Ok(t, err)
				gotName, err := got.Name()
				assert.Ok(t, err)
				assert.Equals(t, tt.args.name, gotName)
			}
			got2, err := tt.client.Deployment(tt.args.name, tt.args.namespace)
			if tt.wantErr {
				assert.NotOk(t, err, tt.wantedErr)
			} else {
				assert.Ok(t, err)
				assert.Assert(t, got == got2, "second call shout returne the same object from cache")
			}
		})
	}
}

func Test_projectClient_Deployments(t *testing.T) {
	tests := []struct {
		name             string
		client           *projectClient
		foundDeployments []string
		namespaceName    string
		namespaceID      string
		wantedLength     int
		wantErr          bool
		wantedErr        string
	}{
		{
			name:             "success",
			client:           simpleProjectClient(),
			foundDeployments: []string{simpleDeploymentName + "1", simpleDeploymentName + "2"},
			namespaceName:    simpleNamespaceName,
			namespaceID:      simpleNamespaceID,
			wantedLength:     2,
			wantErr:          false,
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

			deploymentOperationsStub := stubs.CreateDeploymentOperationsStub(t)
			deploymentOperationsStub.DoList = foundDeployments(tt.foundDeployments, tt.namespaceID, simpleProjectID)
			testClients.ProjectClient.Deployment = deploymentOperationsStub

			namespaceOperationsStub := stubs.CreateNamespaceOperationsStub(t)
			namespaceOperationsStub.DoList = foundNamespace(tt.namespaceName, tt.namespaceID, simpleProjectID)
			testClients.ClusterClient.Namespace = namespaceOperationsStub

			got, err := tt.client.Deployments(tt.namespaceName)
			if tt.wantErr {
				assert.NotOk(t, err, tt.wantedErr)
			} else {
				assert.Ok(t, err)
				assert.Equals(t, tt.wantedLength, len(got))
			}
		})
	}
}

func Test_projectClient_DaemonSet(t *testing.T) {
	type args struct {
		name      string
		namespace string
	}
	tests := []struct {
		name      string
		client    *projectClient
		args      args
		wantErr   bool
		wantedErr string
	}{
		{
			name:   simpleDaemonSetName,
			client: simpleProjectClient(),
			args: args{
				name:      simpleDaemonSetName,
				namespace: simpleNamespaceName,
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

			got, err := tt.client.DaemonSet(tt.args.name, tt.args.namespace)
			if tt.wantErr {
				assert.NotOk(t, err, tt.wantedErr)
			} else {
				assert.Ok(t, err)
				gotName, err := got.Name()
				assert.Ok(t, err)
				assert.Equals(t, tt.args.name, gotName)
			}
			got2, err := tt.client.DaemonSet(tt.args.name, tt.args.namespace)
			if tt.wantErr {
				assert.NotOk(t, err, tt.wantedErr)
			} else {
				assert.Ok(t, err)
				assert.Assert(t, got == got2, "second call shout returne the same object from cache")
			}
		})
	}
}

func Test_projectClient_DaemonSets(t *testing.T) {
	tests := []struct {
		name            string
		client          *projectClient
		foundDaemonSets []string
		namespaceName   string
		namespaceID     string
		wantedLength    int
		wantErr         bool
		wantedErr       string
	}{
		{
			name:            "success",
			client:          simpleProjectClient(),
			foundDaemonSets: []string{simpleDaemonSetName + "1", simpleDaemonSetName + "2"},
			namespaceName:   simpleNamespaceName,
			namespaceID:     simpleNamespaceID,
			wantedLength:    2,
			wantErr:         false,
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

			daemonSetOperationsStub := stubs.CreateDaemonSetOperationsStub(t)
			daemonSetOperationsStub.DoList = foundDaemonSets(tt.foundDaemonSets, tt.namespaceID, simpleProjectID)
			testClients.ProjectClient.DaemonSet = daemonSetOperationsStub

			namespaceOperationsStub := stubs.CreateNamespaceOperationsStub(t)
			namespaceOperationsStub.DoList = foundNamespace(tt.namespaceName, tt.namespaceID, simpleProjectID)
			testClients.ClusterClient.Namespace = namespaceOperationsStub

			got, err := tt.client.DaemonSets(tt.namespaceName)
			if tt.wantErr {
				assert.NotOk(t, err, tt.wantedErr)
			} else {
				assert.Ok(t, err)
				assert.Equals(t, tt.wantedLength, len(got))
			}
		})
	}
}

func Test_projectClient_StatefulSet(t *testing.T) {
	type args struct {
		name      string
		namespace string
	}
	tests := []struct {
		name      string
		client    *projectClient
		args      args
		wantErr   bool
		wantedErr string
	}{
		{
			name:   simpleStatefulSetName,
			client: simpleProjectClient(),
			args: args{
				name:      simpleStatefulSetName,
				namespace: simpleNamespaceName,
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

			got, err := tt.client.StatefulSet(tt.args.name, tt.args.namespace)
			if tt.wantErr {
				assert.NotOk(t, err, tt.wantedErr)
			} else {
				assert.Ok(t, err)
				gotName, err := got.Name()
				assert.Ok(t, err)
				assert.Equals(t, tt.args.name, gotName)
			}
			got2, err := tt.client.StatefulSet(tt.args.name, tt.args.namespace)
			if tt.wantErr {
				assert.NotOk(t, err, tt.wantedErr)
			} else {
				assert.Ok(t, err)
				assert.Assert(t, got == got2, "second call shout returne the same object from cache")
			}
		})
	}
}

func Test_projectClient_StatefulSets(t *testing.T) {
	tests := []struct {
		name              string
		client            *projectClient
		foundStatefulSets []string
		namespaceName     string
		namespaceID       string
		wantedLength      int
		wantErr           bool
		wantedErr         string
	}{
		{
			name:              "success",
			client:            simpleProjectClient(),
			foundStatefulSets: []string{simpleStatefulSetName + "1", simpleStatefulSetName + "2"},
			namespaceName:     simpleNamespaceName,
			namespaceID:       simpleNamespaceID,
			wantedLength:      2,
			wantErr:           false,
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

			statefulSetOperationsStub := stubs.CreateStatefulSetOperationsStub(t)
			statefulSetOperationsStub.DoList = foundStatefulSets(tt.foundStatefulSets, tt.namespaceID, simpleProjectID)
			testClients.ProjectClient.StatefulSet = statefulSetOperationsStub

			namespaceOperationsStub := stubs.CreateNamespaceOperationsStub(t)
			namespaceOperationsStub.DoList = foundNamespace(tt.namespaceName, tt.namespaceID, simpleProjectID)
			testClients.ClusterClient.Namespace = namespaceOperationsStub

			got, err := tt.client.StatefulSets(tt.namespaceName)
			if tt.wantErr {
				assert.NotOk(t, err, tt.wantedErr)
			} else {
				assert.Ok(t, err)
				assert.Equals(t, tt.wantedLength, len(got))
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
		certificateClients:      make(map[string]CertificateClient),
		configMapClients:        make(map[string]ConfigMapClient),
		dockerCredentialClients: make(map[string]DockerCredentialClient),
		secretClients:           make(map[string]ConfigMapClient),
		appClients:              make(map[string]AppClient),
		jobClients:              make(map[string]JobClient),
		cronJobClients:          make(map[string]CronJobClient),
		deploymentClients:       make(map[string]DeploymentClient),
		daemonSetClients:        make(map[string]DaemonSetClient),
		statefulSetClients:      make(map[string]StatefulSetClient),
		catalogClients:          make(map[string]CatalogClient),
	}
}

func existingProjectClient(t *testing.T, projectName, projectID, clusterID string) *projectClient {
	var (
		project     = projectModel.Project{}
		testClients = stubs.CreateBackendStubs(t)
	)
	project.Metadata.Name = projectName

	expectedListOpts := &types.ListOpts{
		Filters: map[string]interface{}{
			"name":      simpleProjectName,
			"clusterId": simpleClusterID,
		},
	}

	projectOperationsStub := stubs.CreateProjectOperationsStub(t)
	projectOperationsStub.DoList = func(opts *types.ListOpts) (*backendRancherClient.ProjectCollection, error) {
		if !reflect.DeepEqual(expectedListOpts, opts) {
			return nil, fmt.Errorf("Unexpected ListOpts %v", opts)
		}
		return &backendRancherClient.ProjectCollection{
			Data: []backendRancherClient.Project{
				backendRancherClient.Project{
					Resource: types.Resource{
						ID: projectID,
					},
					Name: simpleProjectName,
				},
			},
		}, nil
	}
	testClients.ManagementClient.Project = projectOperationsStub
	rancherClient := simpleRancherClient()
	rancherClient._backendRancherClient = testClients.ManagementClient
	clusterClient := simpleClusterClient()
	clusterClient.rancherClient = rancherClient
	clusterClient._backendClusterClient = testClients.ClusterClient
	result, err := newProjectClient(
		simpleProjectName,
		RancherConfig{},
		clusterClient,
		logrus.New().WithFields(logrus.Fields{}),
	)
	assert.Ok(t, err)
	projectClientResult := result.(*projectClient)

	return projectClientResult
}

func notExistingProjectClient(t *testing.T, projectName, projectID, clusterID string) *projectClient {
	var (
		project     = projectModel.Project{}
		testClients = stubs.CreateBackendStubs(t)
	)
	project.Metadata.Name = projectName

	expectedListOpts := &types.ListOpts{
		Filters: map[string]interface{}{
			"name":      simpleProjectName,
			"clusterId": simpleClusterID,
		},
	}

	projectOperationsStub := stubs.CreateProjectOperationsStub(t)
	projectOperationsStub.DoList = func(opts *types.ListOpts) (*backendRancherClient.ProjectCollection, error) {
		if !reflect.DeepEqual(expectedListOpts, opts) {
			return nil, fmt.Errorf("Unexpected ListOpts %v", opts)
		}
		return &backendRancherClient.ProjectCollection{
			Data: []backendRancherClient.Project{},
		}, nil
	}
	projectOperationsStub.DoCreate = func(opts *managementClient.Project) (*managementClient.Project, error) {
		if opts.Name != projectName || opts.ClusterID != clusterID {
			return nil, fmt.Errorf("Unexpected Create %v", opts)
		}
		opts.ID = projectID
		return opts, nil
	}

	testClients.ManagementClient.Project = projectOperationsStub
	rancherClient := simpleRancherClient()
	rancherClient._backendRancherClient = testClients.ManagementClient
	clusterClient := simpleClusterClient()
	clusterClient.rancherClient = rancherClient
	clusterClient._backendClusterClient = testClients.ClusterClient
	result, err := newProjectClient(
		simpleProjectName,
		RancherConfig{},
		clusterClient,
		logrus.New().WithFields(logrus.Fields{}),
	)
	assert.Ok(t, err)
	projectClientResult := result.(*projectClient)
	return projectClientResult
}

func foundGlobalCertificates(names []string, projectID string) func(opts *types.ListOpts) (*backendProjectClient.CertificateCollection, error) {
	data := make([]backendProjectClient.Certificate, 0)
	for _, name := range names {
		data = append(data, backendProjectClient.Certificate{Name: name, ProjectID: projectID})
	}
	return func(opts *types.ListOpts) (*backendProjectClient.CertificateCollection, error) {
		return &backendProjectClient.CertificateCollection{
			Data: data,
		}, nil
	}
}

func foundCertificates(names []string, namespaceID string, projectID string) func(opts *types.ListOpts) (*backendProjectClient.NamespacedCertificateCollection, error) {
	data := make([]backendProjectClient.NamespacedCertificate, 0)
	for _, name := range names {
		data = append(data, backendProjectClient.NamespacedCertificate{Name: name, NamespaceId: namespaceID, ProjectID: projectID})
	}
	return func(opts *types.ListOpts) (*backendProjectClient.NamespacedCertificateCollection, error) {
		return &backendProjectClient.NamespacedCertificateCollection{
			Data: data,
		}, nil
	}
}

func foundConfigMaps(names []string, namespaceID string, projectID string) func(opts *types.ListOpts) (*backendProjectClient.ConfigMapCollection, error) {
	data := make([]backendProjectClient.ConfigMap, 0)
	for _, name := range names {
		data = append(data, backendProjectClient.ConfigMap{Name: name, NamespaceId: namespaceID, ProjectID: projectID})
	}
	return func(opts *types.ListOpts) (*backendProjectClient.ConfigMapCollection, error) {
		return &backendProjectClient.ConfigMapCollection{
			Data: data,
		}, nil
	}
}

func foundGlobalDockerCredentials(names []string, projectID string) func(opts *types.ListOpts) (*backendProjectClient.DockerCredentialCollection, error) {
	data := make([]backendProjectClient.DockerCredential, 0)
	for _, name := range names {
		data = append(data, backendProjectClient.DockerCredential{Name: name, ProjectID: projectID})
	}
	return func(opts *types.ListOpts) (*backendProjectClient.DockerCredentialCollection, error) {
		return &backendProjectClient.DockerCredentialCollection{
			Data: data,
		}, nil
	}
}

func foundDockerCredentials(names []string, namespaceID string, projectID string) func(opts *types.ListOpts) (*backendProjectClient.NamespacedDockerCredentialCollection, error) {
	data := make([]backendProjectClient.NamespacedDockerCredential, 0)
	for _, name := range names {
		data = append(data, backendProjectClient.NamespacedDockerCredential{Name: name, NamespaceId: namespaceID, ProjectID: projectID})
	}
	return func(opts *types.ListOpts) (*backendProjectClient.NamespacedDockerCredentialCollection, error) {
		return &backendProjectClient.NamespacedDockerCredentialCollection{
			Data: data,
		}, nil
	}
}

func foundGlobalSecrets(names []string, projectID string) func(opts *types.ListOpts) (*backendProjectClient.SecretCollection, error) {
	data := make([]backendProjectClient.Secret, 0)
	for _, name := range names {
		data = append(data, backendProjectClient.Secret{Name: name, ProjectID: projectID})
	}
	return func(opts *types.ListOpts) (*backendProjectClient.SecretCollection, error) {
		return &backendProjectClient.SecretCollection{
			Data: data,
		}, nil
	}
}

func foundDockerSecrets(names []string, namespaceID string, projectID string) func(opts *types.ListOpts) (*backendProjectClient.NamespacedSecretCollection, error) {
	data := make([]backendProjectClient.NamespacedSecret, 0)
	for _, name := range names {
		data = append(data, backendProjectClient.NamespacedSecret{Name: name, NamespaceId: namespaceID, ProjectID: projectID})
	}
	return func(opts *types.ListOpts) (*backendProjectClient.NamespacedSecretCollection, error) {
		return &backendProjectClient.NamespacedSecretCollection{
			Data: data,
		}, nil
	}
}

func foundApps(names []string, projectID string) func(opts *types.ListOpts) (*backendProjectClient.AppCollection, error) {
	data := make([]backendProjectClient.App, 0)
	for _, name := range names {
		data = append(data, backendProjectClient.App{Name: name, ProjectID: projectID})
	}
	return func(opts *types.ListOpts) (*backendProjectClient.AppCollection, error) {
		return &backendProjectClient.AppCollection{
			Data: data,
		}, nil
	}
}

func foundJobs(names []string, namespaceID string, projectID string) func(opts *types.ListOpts) (*backendProjectClient.JobCollection, error) {
	data := make([]backendProjectClient.Job, 0)
	for _, name := range names {
		data = append(data, backendProjectClient.Job{Name: name, NamespaceId: namespaceID, ProjectID: projectID})
	}
	return func(opts *types.ListOpts) (*backendProjectClient.JobCollection, error) {
		return &backendProjectClient.JobCollection{
			Data: data,
		}, nil
	}
}

func foundCronJobs(names []string, namespaceID string, projectID string) func(opts *types.ListOpts) (*backendProjectClient.CronJobCollection, error) {
	data := make([]backendProjectClient.CronJob, 0)
	for _, name := range names {
		data = append(data, backendProjectClient.CronJob{Name: name, NamespaceId: namespaceID, ProjectID: projectID})
	}
	return func(opts *types.ListOpts) (*backendProjectClient.CronJobCollection, error) {
		return &backendProjectClient.CronJobCollection{
			Data: data,
		}, nil
	}
}

func foundDeployments(names []string, namespaceID string, projectID string) func(opts *types.ListOpts) (*backendProjectClient.DeploymentCollection, error) {
	data := make([]backendProjectClient.Deployment, 0)
	for _, name := range names {
		data = append(data, backendProjectClient.Deployment{Name: name, NamespaceId: namespaceID, ProjectID: projectID})
	}
	return func(opts *types.ListOpts) (*backendProjectClient.DeploymentCollection, error) {
		return &backendProjectClient.DeploymentCollection{
			Data: data,
		}, nil
	}
}

func foundDaemonSets(names []string, namespaceID string, projectID string) func(opts *types.ListOpts) (*backendProjectClient.DaemonSetCollection, error) {
	data := make([]backendProjectClient.DaemonSet, 0)
	for _, name := range names {
		data = append(data, backendProjectClient.DaemonSet{Name: name, NamespaceId: namespaceID, ProjectID: projectID})
	}
	return func(opts *types.ListOpts) (*backendProjectClient.DaemonSetCollection, error) {
		return &backendProjectClient.DaemonSetCollection{
			Data: data,
		}, nil
	}
}

func foundStatefulSets(names []string, namespaceID string, projectID string) func(opts *types.ListOpts) (*backendProjectClient.StatefulSetCollection, error) {
	data := make([]backendProjectClient.StatefulSet, 0)
	for _, name := range names {
		data = append(data, backendProjectClient.StatefulSet{Name: name, NamespaceId: namespaceID, ProjectID: projectID})
	}
	return func(opts *types.ListOpts) (*backendProjectClient.StatefulSetCollection, error) {
		return &backendProjectClient.StatefulSetCollection{
			Data: data,
		}, nil
	}
}
