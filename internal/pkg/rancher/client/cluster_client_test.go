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
	backendClusterClient "github.com/rancher/types/client/cluster/v3"
	managementClient "github.com/rancher/types/client/management/v3"
	"github.com/sirupsen/logrus"
)

const (
	simpleClusterName = "simple-cluster"
	simpleClusterID   = "simple-cluster-id"
)

func Test_clusterClient_Project(t *testing.T) {
	type args struct {
		projectName string
	}
	tests := []struct {
		name      string
		client    *clusterClient
		args      args
		wantErr   bool
		wantedErr string
	}{
		{
			name:   simpleClusterName,
			client: simpleClusterClient(),
			args: args{
				projectName: "simple-project",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Arrange
			testClients := stubs.CreateBackendStubs(t)
			rancherClient := simpleRancherClient()
			rancherClient._backendRancherClient = testClients.ManagementClient
			tt.client.rancherClient = rancherClient
			tt.client._backendClusterClient = testClients.ClusterClient

			got, err := tt.client.Project(tt.args.projectName)
			if tt.wantErr {
				assert.NotOk(t, err, tt.wantedErr)
			} else {
				assert.Ok(t, err)
				gotName, err := got.Name()
				assert.Ok(t, err)
				assert.Equals(t, tt.args.projectName, gotName)
			}
		})
	}
}

func Test_clusterClient_Projects(t *testing.T) {
	tests := []struct {
		name          string
		client        *clusterClient
		foundProjects []string
		wantedLength  int
		wantErr       bool
		wantedErr     string
	}{
		{
			name:          "success",
			client:        simpleClusterClient(),
			foundProjects: []string{"simple-project1", "simple-project2"},
			wantedLength:  2,
			wantErr:       false,
		},
		{
			name:          "only-one-project-per-unique-name",
			client:        simpleClusterClient(),
			foundProjects: []string{"simple-project1", "simple-project1"},
			wantedLength:  2,
			wantErr:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Arrange
			testClients := stubs.CreateBackendStubs(t)
			rancherClient := simpleRancherClient()
			rancherClient._backendRancherClient = testClients.ManagementClient
			tt.client.rancherClient = rancherClient
			tt.client._backendClusterClient = testClients.ClusterClient

			projectOperationsStub := stubs.CreateProjectOperationsStub(t)
			projectOperationsStub.DoList = foundProjects(tt.foundProjects)
			testClients.ManagementClient.Project = projectOperationsStub

			got, err := tt.client.Projects()
			if tt.wantErr {
				assert.NotOk(t, err, tt.wantedErr)
			} else {
				assert.Ok(t, err)
				assert.Equals(t, tt.wantedLength, len(got))
			}
		})
	}
}

func Test_clusterClient_StorageClass(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name      string
		client    *clusterClient
		args      args
		wantErr   bool
		wantedErr string
	}{
		{
			name:   simpleClusterName,
			client: simpleClusterClient(),
			args: args{
				name: "simple-storage-class",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Arrange
			testClients := stubs.CreateBackendStubs(t)
			rancherClient := simpleRancherClient()
			rancherClient._backendRancherClient = testClients.ManagementClient
			tt.client.rancherClient = rancherClient
			tt.client._backendClusterClient = testClients.ClusterClient

			got, err := tt.client.StorageClass(tt.args.name)
			if tt.wantErr {
				assert.NotOk(t, err, tt.wantedErr)
			} else {
				assert.Ok(t, err)
				gotName, err := got.Name()
				assert.Ok(t, err)
				assert.Equals(t, tt.args.name, gotName)
			}
		})
	}
}

func Test_clusterClient_StorageClasses(t *testing.T) {
	tests := []struct {
		name                string
		client              *clusterClient
		foundStorageClasses []string
		wantedLength        int
		wantErr             bool
		wantedErr           string
	}{
		{
			name:                "success",
			client:              simpleClusterClient(),
			foundStorageClasses: []string{"simple-storage-class1", "simple-storage-class2"},
			wantedLength:        2,
			wantErr:             false,
		},
		{
			name:                "only-one-storage-class-per-unique-name",
			client:              simpleClusterClient(),
			foundStorageClasses: []string{"simple-storage-class1", "simple-storage-class1"},
			wantedLength:        2,
			wantErr:             false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Arrange
			testClients := stubs.CreateBackendStubs(t)
			rancherClient := simpleRancherClient()
			rancherClient._backendRancherClient = testClients.ManagementClient
			tt.client.rancherClient = rancherClient
			tt.client._backendClusterClient = testClients.ClusterClient

			storageClassOperationStub := stubs.CreateStorageClassOperationsStub(t)
			storageClassOperationStub.DoList = foundStorageClasses(tt.foundStorageClasses)
			testClients.ClusterClient.StorageClass = storageClassOperationStub

			got, err := tt.client.StorageClasses()
			if tt.wantErr {
				assert.NotOk(t, err, tt.wantedErr)
			} else {
				assert.Ok(t, err)
				assert.Equals(t, tt.wantedLength, len(got))
			}
		})
	}
}

func Test_clusterClient_PersistentVolume(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name      string
		client    *clusterClient
		args      args
		wantErr   bool
		wantedErr string
	}{
		{
			name:   simpleClusterName,
			client: simpleClusterClient(),
			args: args{
				name: "simple-storage-class",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Arrange
			testClients := stubs.CreateBackendStubs(t)
			rancherClient := simpleRancherClient()
			rancherClient._backendRancherClient = testClients.ManagementClient
			tt.client.rancherClient = rancherClient
			tt.client._backendClusterClient = testClients.ClusterClient

			got, err := tt.client.PersistentVolume(tt.args.name)
			if tt.wantErr {
				assert.NotOk(t, err, tt.wantedErr)
			} else {
				assert.Ok(t, err)
				gotName, err := got.Name()
				assert.Ok(t, err)
				assert.Equals(t, tt.args.name, gotName)
			}
		})
	}
}

func Test_clusterClient_PersistentVolumes(t *testing.T) {
	tests := []struct {
		name                   string
		client                 *clusterClient
		foundPersistentVolumes []string
		wantedLength           int
		wantErr                bool
		wantedErr              string
	}{
		{
			name:                   "success",
			client:                 simpleClusterClient(),
			foundPersistentVolumes: []string{"simple-storage-class1", "simple-storage-class2"},
			wantedLength:           2,
			wantErr:                false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Arrange
			testClients := stubs.CreateBackendStubs(t)
			rancherClient := simpleRancherClient()
			rancherClient._backendRancherClient = testClients.ManagementClient
			tt.client.rancherClient = rancherClient
			tt.client._backendClusterClient = testClients.ClusterClient

			persistentVolumeOperationStub := stubs.CreatePersistentVolumeOperationsStub(t)
			persistentVolumeOperationStub.DoList = foundPersistentVolumes(tt.foundPersistentVolumes)
			testClients.ClusterClient.PersistentVolume = persistentVolumeOperationStub

			got, err := tt.client.PersistentVolumes()
			if tt.wantErr {
				assert.NotOk(t, err, tt.wantedErr)
			} else {
				assert.Ok(t, err)
				assert.Equals(t, tt.wantedLength, len(got))
			}
		})
	}
}

func Test_clusterClient_Namespace(t *testing.T) {
	type args struct {
		name        string
		projectName string
	}
	tests := []struct {
		name      string
		client    *clusterClient
		args      args
		wantErr   bool
		wantedErr string
	}{
		{
			name:   simpleClusterName,
			client: simpleClusterClient(),
			args: args{
				name:        simpleNamespaceName,
				projectName: simpleProjectName,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Arrange
			testClients := stubs.CreateBackendStubs(t)
			rancherClient := simpleRancherClient()
			rancherClient._backendRancherClient = testClients.ManagementClient
			tt.client.rancherClient = rancherClient
			tt.client._backendClusterClient = testClients.ClusterClient

			got, err := tt.client.Namespace(tt.args.name, tt.args.projectName)
			if tt.wantErr {
				assert.NotOk(t, err, tt.wantedErr)
			} else {
				assert.Ok(t, err)
				gotName, err := got.Name()
				assert.Ok(t, err)
				assert.Equals(t, tt.args.name, gotName)
			}
		})
	}
}

func Test_clusterClient_Namespaces(t *testing.T) {
	tests := []struct {
		name            string
		client          *clusterClient
		foundNamespaces []string
		wantedLength    int
		wantErr         bool
		wantedErr       string
	}{
		{
			name:            "success",
			client:          simpleClusterClient(),
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
			rancherClient._backendRancherClient = testClients.ManagementClient
			tt.client.rancherClient = rancherClient
			tt.client._backendClusterClient = testClients.ClusterClient

			namespaceOperationStub := stubs.CreateNamespaceOperationsStub(t)
			namespaceOperationStub.DoList = foundNamespaces(tt.foundNamespaces, simpleProjectID)
			testClients.ClusterClient.Namespace = namespaceOperationStub

			got, err := tt.client.Namespaces(simpleProjectName)
			if tt.wantErr {
				assert.NotOk(t, err, tt.wantedErr)
			} else {
				assert.Ok(t, err)
				assert.Equals(t, tt.wantedLength, len(got))
			}
		})
	}
}

func simpleClusterClient() *clusterClient {
	logrus.SetLevel(logrus.TraceLevel)
	return &clusterClient{
		resourceClient: resourceClient{
			name:   simpleClusterName,
			id:     simpleClusterID,
			logger: logrus.WithFields(logrus.Fields{}),
		},
		config: RancherConfig{},
		projectClients: map[string]ProjectClient{
			simpleProjectName: simpleProjectClient(),
		},
		storageClasses:    make(map[string]StorageClassClient),
		persistentVolumes: make(map[string]PersistentVolumeClient),
		namespaces:        make(map[string]namespaceCacheEntry),
	}
}

func foundProjects(names []string) func(opts *types.ListOpts) (*managementClient.ProjectCollection, error) {
	data := make([]managementClient.Project, 0)
	for _, name := range names {
		data = append(data, managementClient.Project{Name: name})
	}
	return func(opts *types.ListOpts) (*managementClient.ProjectCollection, error) {
		return &managementClient.ProjectCollection{
			Data: data,
		}, nil
	}
}

func foundStorageClasses(names []string) func(opts *types.ListOpts) (*backendClusterClient.StorageClassCollection, error) {
	data := make([]backendClusterClient.StorageClass, 0)
	for _, name := range names {
		data = append(data, backendClusterClient.StorageClass{Name: name})
	}
	return func(opts *types.ListOpts) (*backendClusterClient.StorageClassCollection, error) {
		return &backendClusterClient.StorageClassCollection{
			Data: data,
		}, nil
	}
}

func foundPersistentVolumes(names []string) func(opts *types.ListOpts) (*backendClusterClient.PersistentVolumeCollection, error) {
	data := make([]backendClusterClient.PersistentVolume, 0)
	for _, name := range names {
		data = append(data, backendClusterClient.PersistentVolume{Name: name})
	}
	return func(opts *types.ListOpts) (*backendClusterClient.PersistentVolumeCollection, error) {
		return &backendClusterClient.PersistentVolumeCollection{
			Data: data,
		}, nil
	}
}

func foundNamespaces(names []string, projectID string) func(opts *types.ListOpts) (*backendClusterClient.NamespaceCollection, error) {
	data := make([]backendClusterClient.Namespace, 0)
	for _, name := range names {
		data = append(data, backendClusterClient.Namespace{Name: name, ProjectID: projectID})
	}
	return func(opts *types.ListOpts) (*backendClusterClient.NamespaceCollection, error) {
		return &backendClusterClient.NamespaceCollection{
			Data: data,
		}, nil
	}
}
