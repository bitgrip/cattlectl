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

package stubs

import (
	"testing"

	"github.com/rancher/norman/clientbase"
	clusterClient "github.com/rancher/types/client/cluster/v3"
	managementClient "github.com/rancher/types/client/management/v3"
	projectClient "github.com/rancher/types/client/project/v3"
)

// CreateBackendStubs builds a set of testClients
func CreateBackendStubs(tb testing.TB) *BackendStubs {
	testClients := &BackendStubs{
		ClusterClient: &clusterClient.Client{},
		ManagementClient: &managementClient.Client{
			APIBaseClient: clientbase.APIBaseClient{
				Ops: &clientbase.APIOperations{},
			},
		},
		ProjectClient: &projectClient.Client{},
	}
	testClients.ClusterClient.Namespace = CreateNamespaceOperationsStub(tb)
	testClients.ClusterClient.StorageClass = CreateStorageClassOperationsStub(tb)
	testClients.ClusterClient.PersistentVolume = CreatePersistentVolumeOperationsStub(tb)
	testClients.ManagementClient.Cluster = CreateClusterOperationsStub(tb)
	testClients.ManagementClient.Project = CreateProjectOperationsStub(tb)
	testClients.ProjectClient.App = CreateAppOperationsStub(tb)
	testClients.ProjectClient.Job = CreateJobOperationsStub(tb)
	return testClients
}

// BackendStubs is a grouping structure for the rancher clients
type BackendStubs struct {
	ClusterClient    *clusterClient.Client
	ManagementClient *managementClient.Client
	ProjectClient    *projectClient.Client
}
