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

package rancher

import (
	"github.com/bitgrip/cattlectl/internal/pkg/rancher/stubs"
	"github.com/rancher/norman/clientbase"
	clusterClient "github.com/rancher/types/client/cluster/v3"
	managementClient "github.com/rancher/types/client/management/v3"
	projectClient "github.com/rancher/types/client/project/v3"
)

type deferData struct {
	newManagementClient func(opts *clientbase.ClientOpts) (*managementClient.Client, error)
	newClusterClient    func(opts *clientbase.ClientOpts) (*clusterClient.Client, error)
	newProjectClient    func(opts *clientbase.ClientOpts) (*projectClient.Client, error)
}

func initTestClientFactories(testClients *stubs.BackendStubs) deferData {
	result := deferData{
		newManagementClient: newManagementClient,
		newClusterClient:    newClusterClient,
		newProjectClient:    newProjectClient,
	}
	newClusterClient = func(opts *clientbase.ClientOpts) (*clusterClient.Client, error) {
		return testClients.ClusterClient, nil
	}
	newManagementClient = func(opts *clientbase.ClientOpts) (*managementClient.Client, error) {
		return testClients.ManagementClient, nil
	}
	newProjectClient = func(opts *clientbase.ClientOpts) (*projectClient.Client, error) {
		return testClients.ProjectClient, nil
	}
	return result
}

func deferTestClientFactories(data deferData) {
	newManagementClient = data.newManagementClient
	newClusterClient = data.newClusterClient
	newProjectClient = data.newProjectClient
}
