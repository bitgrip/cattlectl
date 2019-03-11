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

	clusterClient "github.com/rancher/types/client/cluster/v3"
	managementClient "github.com/rancher/types/client/management/v3"
	projectClient "github.com/rancher/types/client/project/v3"
)

// CreateTestClients builds a set of testClients
func CreateTestClients(tb testing.TB) *TestClients {
	testClients := &TestClients{
		ClusterClient:    &clusterClient.Client{},
		ManagementClient: &managementClient.Client{},
		ProjectClient:    &projectClient.Client{},
	}
	testClients.ClusterClient.Namespace = CreateNamespaceOperationsStub(tb)
	testClients.ClusterClient.StorageClass = CreateStorageClassOperationsStub(tb)
	testClients.ClusterClient.PersistentVolume = CreatePersistentVolumeOperationsStub(tb)
	testClients.ManagementClient.Project = CreateProjectOperationsStub(tb)
	testClients.ProjectClient.App = CreateAppOperationsStub(tb)
	testClients.ProjectClient.Job = CreateJobOperationsStub(tb)
	return testClients
}

// TestClients is a grouping structure for the rancher clients
type TestClients struct {
	ClusterClient    *clusterClient.Client
	ManagementClient *managementClient.Client
	ProjectClient    *projectClient.Client
}
