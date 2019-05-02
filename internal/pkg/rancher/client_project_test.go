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
	"fmt"
	"testing"

	"github.com/bitgrip/cattlectl/internal/pkg/assert"
	"github.com/bitgrip/cattlectl/internal/pkg/rancher/stubs"
	"github.com/rancher/norman/clientbase"
	"github.com/rancher/norman/types"
	clusterClient "github.com/rancher/types/client/cluster/v3"
	managementClient "github.com/rancher/types/client/management/v3"
	projectClient "github.com/rancher/types/client/project/v3"
)

/*
HasProject(...)
  returns true if project exists
  sets found project ID to client field
  creates and sets project client to client field
*/
func TestHasProject_ProjectExisting(t *testing.T) {
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

	testClients := stubs.CreateBackendStubs(t)
	newClusterClient = func(opts *clientbase.ClientOpts) (*clusterClient.Client, error) {
		return testClients.ClusterClient, nil
	}
	newManagementClient = func(opts *clientbase.ClientOpts) (*managementClient.Client, error) {
		return testClients.ManagementClient, nil
	}
	newProjectClient = func(opts *clientbase.ClientOpts) (*projectClient.Client, error) {
		return testClients.ProjectClient, nil
	}

	var actualListOpts *types.ListOpts
	const projectId = "test-project-id"
	projectOperationsStub := stubs.CreateProjectOperationsStub(t)
	projectOperationsStub.DoList = func(opts *types.ListOpts) (*managementClient.ProjectCollection, error) {
		actualListOpts = opts
		return &managementClient.ProjectCollection{
			Data: []managementClient.Project{
				managementClient.Project{
					Name:     fmt.Sprint(opts.Filters["name"]),
					Resource: types.Resource{ID: projectId},
				},
			},
		}, nil
	}
	testClients.ManagementClient.Project = projectOperationsStub

	projectName := "test-project-name"
	clientConfig := ClientConfig{}
	client, err := NewClient(clientConfig)
	assert.Ok(t, err)
	client.(*rancherClient).clusterId = "test-cluster-id"

	//Act
	result, foundProjectId, err := client.HasProjectWithName(projectName)

	//Assert
	assert.Ok(t, err)
	assert.Equals(t, true, result)
	assert.Equals(t, projectId, foundProjectId)
	assert.Equals(t, &types.ListOpts{Filters: map[string]interface{}{"clusterId": "test-cluster-id", "name": "test-project-name"}}, actualListOpts)

}

/*
HasProject(...)
  returns false if project dose not exists
  not sets project ID to client field
  not creates project client
*/
func TestHasProject_ProjectNotExisting(t *testing.T) {
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

	testClients := stubs.CreateBackendStubs(t)
	newClusterClient = func(opts *clientbase.ClientOpts) (*clusterClient.Client, error) {
		return testClients.ClusterClient, nil
	}
	newManagementClient = func(opts *clientbase.ClientOpts) (*managementClient.Client, error) {
		return testClients.ManagementClient, nil
	}
	newProjectClient = func(opts *clientbase.ClientOpts) (*projectClient.Client, error) {
		assert.FailInStub(t, 1, "unexpectedCall newProjectClient")
		return nil, nil
	}

	var actualListOpts *types.ListOpts
	const projectId = "test-project-id"
	projectOperationsStub := stubs.CreateProjectOperationsStub(t)
	projectOperationsStub.DoList = func(opts *types.ListOpts) (*managementClient.ProjectCollection, error) {
		actualListOpts = opts
		return &managementClient.ProjectCollection{
			Data: []managementClient.Project{},
		}, nil
	}
	testClients.ManagementClient.Project = projectOperationsStub

	projectName := "test-project-name"
	clientConfig := ClientConfig{}
	client, err := NewClient(clientConfig)
	assert.Ok(t, err)
	client.(*rancherClient).clusterId = "test-cluster-id"

	//Act
	result, foundProjectId, err := client.HasProjectWithName(projectName)

	//Assert
	assert.Ok(t, err)
	assert.Equals(t, false, result)
	assert.Equals(t, "", foundProjectId)
	assert.Equals(t, &types.ListOpts{Filters: map[string]interface{}{"clusterId": "test-cluster-id", "name": "test-project-name"}}, actualListOpts)

}

/*
CreateProject(...)
  uses management client to create project
  sets created project ID to client field
  creates and sets project client to client field
*/
func TestCreateProject(t *testing.T) {
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

	const (
		projectId   = "test-project-id"
		projectName = "test-project-name"
		clusterId   = "test-cluster-id"
	)
	var (
		actualProjectPattern *managementClient.Project
		clientConfig         = ClientConfig{}
		testClients          = stubs.CreateBackendStubs(t)
	)
	newClusterClient = func(opts *clientbase.ClientOpts) (*clusterClient.Client, error) {
		return testClients.ClusterClient, nil
	}
	newManagementClient = func(opts *clientbase.ClientOpts) (*managementClient.Client, error) {
		return testClients.ManagementClient, nil
	}
	newProjectClient = func(opts *clientbase.ClientOpts) (*projectClient.Client, error) {
		return testClients.ProjectClient, nil
	}
	projectOperationsStub := stubs.CreateProjectOperationsStub(t)
	projectOperationsStub.DoCreate = func(project *managementClient.Project) (*managementClient.Project, error) {
		actualProjectPattern = project
		return &managementClient.Project{Resource: types.Resource{ID: projectId}}, nil
	}
	testClients.ManagementClient.Project = projectOperationsStub
	client, err := NewClient(clientConfig)
	assert.Ok(t, err)
	client.(*rancherClient).clusterId = "test-cluster-id"

	//Act
	createdProjectId, err := client.CreateProject(projectName)

	//Assert
	assert.Ok(t, err)
	assert.Equals(t, &managementClient.Project{Name: projectName, ClusterID: clusterId}, actualProjectPattern)
	assert.Equals(t, projectId, createdProjectId)

}
