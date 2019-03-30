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
	"errors"
	"testing"

	"github.com/bitgrip/cattlectl/internal/pkg/assert"
	"github.com/rancher/norman/clientbase"
	clusterClient "github.com/rancher/types/client/cluster/v3"
	managementClient "github.com/rancher/types/client/management/v3"
	projectClient "github.com/rancher/types/client/project/v3"
)

/*
CreateClient(...)
  returns error if create management client fails
*/
func TestCreateClient_ManagementClientFails(t *testing.T) {
	// Arrange
	oldNewClusterClient := newClusterClient
	defer func() {
		newClusterClient = oldNewClusterClient
	}()
	oldNewProjectClient := newProjectClient
	defer func() {
		newProjectClient = oldNewProjectClient
	}()
	oldNewManagementClient := newManagementClient
	defer func() {
		newManagementClient = oldNewManagementClient
	}()

	//testClients := stubs.CreateBackendStubs()
	newClusterClient = func(opts *clientbase.ClientOpts) (*clusterClient.Client, error) {
		return nil, nil
	}
	newManagementClient = func(opts *clientbase.ClientOpts) (*managementClient.Client, error) {
		return nil, errors.New("Test-Error")
	}
	newProjectClient = func(opts *clientbase.ClientOpts) (*projectClient.Client, error) {
		assert.FailInStub(t, 1, "unexpectedCall newProjectClient")
		return nil, nil
	}

	clientConfig := ClientConfig{}

	// Act
	_, err := NewClient(clientConfig)
	assert.NotOk(t, err, "Failed to create management client, Test-Error")

	// Assert
}
