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
	rancherModel "github.com/bitgrip/cattlectl/internal/pkg/rancher/model"
	"github.com/bitgrip/cattlectl/internal/pkg/rancher/stubs"
	"github.com/rancher/norman/types"
	backendRancherClient "github.com/rancher/types/client/management/v3"
	"github.com/sirupsen/logrus"
)

func Test_projectCatalogClient_Exists(t *testing.T) {
	tests := []struct {
		name      string
		client    *projectCatalogClient
		wanted    bool
		wantErr   bool
		wantedErr string
	}{
		{
			name: "Existing",
			client: existingProjectCatalogClient(
				t,
				simpleCatalogName,
				simpleProjectID,
				simpleURL,
				simpleBranch,
				simpleUsername,
				simplePassword,
			),
			wanted:  true,
			wantErr: false,
		},
		{
			name: "Not_Existing",
			client: notExistingProjectCatalogClient(
				t,
				simpleCatalogName,
				simpleProjectID,
				simpleURL,
				simpleBranch,
				simpleUsername,
				simplePassword,
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

func Test_projectCatalogClient_Create(t *testing.T) {
	tests := []struct {
		name      string
		client    *projectCatalogClient
		wantErr   bool
		wantedErr string
	}{
		{
			name: "Create",
			client: notExistingProjectCatalogClient(
				t,
				simpleCatalogName,
				simpleProjectID,
				simpleURL,
				simpleBranch,
				simpleUsername,
				simplePassword,
			),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.client.Create(false)
			if tt.wantErr {
				assert.NotOk(t, err, tt.wantedErr)
			} else {
				assert.Ok(t, err)
			}
		})
	}
}

func Test_projectCatalogClient_Upgrade(t *testing.T) {
	tests := []struct {
		name      string
		client    *projectCatalogClient
		wantErr   bool
		wantedErr string
	}{
		{
			name: "Create",
			client: existingProjectCatalogClient(
				t,
				simpleCatalogName,
				simpleProjectID,
				simpleURL,
				simpleBranch,
				simpleUsername,
				simplePassword,
			),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.client.Upgrade(false)
			if tt.wantErr {
				assert.NotOk(t, err, tt.wantedErr)
			} else {
				assert.Ok(t, err)
			}
		})
	}
}

func existingProjectCatalogClient(t *testing.T, name, projectID, url, branch, username, password string) *projectCatalogClient {
	testClients := stubs.CreateBackendStubs(t)
	expectedListOpts := &types.ListOpts{
		Filters: map[string]interface{}{
			"name":      name,
			"projectID": projectID,
		},
	}
	projectCatalogData := rancherModel.Catalog{
		Name:     name,
		URL:      url,
		Branch:   branch,
		Username: username,
		Password: password,
	}
	expectedBackendrCatalog := &backendRancherClient.ProjectCatalog{
		Name:      name,
		ProjectID: projectID,
		URL:       url,
		Branch:    branch,
		Username:  username,
		Password:  password,
		Labels:    map[string]string{"cattlectl.io/hash": "d20875c8c699ed126b385992bf8fc7c384f18e85"},
	}

	projectCatalogOperationsStub := stubs.CreateProjectCatalogOperationsStub(t)
	projectCatalogOperationsStub.DoList = func(opts *types.ListOpts) (*backendRancherClient.ProjectCatalogCollection, error) {
		if !reflect.DeepEqual(expectedListOpts, opts) {
			return nil, fmt.Errorf("Unexpected ListOpts %v", opts)
		}
		return &backendRancherClient.ProjectCatalogCollection{
			Data: []backendRancherClient.ProjectCatalog{
				backendRancherClient.ProjectCatalog{
					Name:      name,
					ProjectID: projectID,
					Labels:    map[string]string{},
				},
			},
		}, nil
	}
	projectCatalogOperationsStub.DoReplace = func(existing *backendRancherClient.ProjectCatalog) (*backendRancherClient.ProjectCatalog, error) {
		if !reflect.DeepEqual(expectedBackendrCatalog, existing) {
			return nil, fmt.Errorf("Unexpected ProjectCatalog %v", existing)
		}
		return existing, nil
	}
	testClients.ManagementClient.ProjectCatalog = projectCatalogOperationsStub
	rancherClient := simpleRancherClient()
	rancherClient._backendRancherClient = testClients.ManagementClient
	clusterClient := simpleClusterClient()
	clusterClient.rancherClient = rancherClient
	projectClient := simpleProjectClient()
	projectClient.clusterClient = clusterClient
	projectClient._backendProjectClient = testClients.ProjectClient
	result, err := newProjectCatalogClient(
		name,
		projectClient,
		logrus.New().WithFields(logrus.Fields{}),
	)
	assert.Ok(t, err)
	projectCatalogClientResult := result.(*projectCatalogClient)
	projectCatalogClientResult.catalog = projectCatalogData
	return projectCatalogClientResult
}

func notExistingProjectCatalogClient(t *testing.T, name, projectID, url, branch, username, password string) *projectCatalogClient {
	testClients := stubs.CreateBackendStubs(t)
	expectedListOpts := &types.ListOpts{
		Filters: map[string]interface{}{
			"name":      name,
			"projectID": projectID,
		},
	}
	projectCatalogData := rancherModel.Catalog{
		Name:     name,
		URL:      url,
		Branch:   branch,
		Username: username,
		Password: password,
	}
	expectedBackendrCatalog := &backendRancherClient.ProjectCatalog{
		Name:      name,
		ProjectID: projectID,
		URL:       url,
		Branch:    branch,
		Username:  username,
		Password:  password,
		Labels:    map[string]string{"cattlectl.io/hash": hashOf(projectCatalogData)},
	}

	projectCatalogOperationsStub := stubs.CreateProjectCatalogOperationsStub(t)
	projectCatalogOperationsStub.DoList = func(opts *types.ListOpts) (*backendRancherClient.ProjectCatalogCollection, error) {
		if !reflect.DeepEqual(expectedListOpts, opts) {
			return nil, fmt.Errorf("Unexpected ListOpts %v", opts)
		}
		return &backendRancherClient.ProjectCatalogCollection{
			Data: []backendRancherClient.ProjectCatalog{},
		}, nil
	}
	projectCatalogOperationsStub.DoCreate = func(projectCatalog *backendRancherClient.ProjectCatalog) (*backendRancherClient.ProjectCatalog, error) {
		if !reflect.DeepEqual(expectedBackendrCatalog, projectCatalog) {
			return nil, fmt.Errorf("Unexpected Catalog %v", projectCatalog)
		}

		return projectCatalog, nil
	}
	testClients.ManagementClient.ProjectCatalog = projectCatalogOperationsStub
	rancherClient := simpleRancherClient()
	rancherClient._backendRancherClient = testClients.ManagementClient
	clusterClient := simpleClusterClient()
	clusterClient.rancherClient = rancherClient
	projectClient := simpleProjectClient()
	projectClient.clusterClient = clusterClient
	projectClient._backendProjectClient = testClients.ProjectClient
	result, err := newProjectCatalogClient(
		name,
		projectClient,
		logrus.New().WithFields(logrus.Fields{}),
	)
	assert.Ok(t, err)
	projectCatalogClientResult := result.(*projectCatalogClient)
	projectCatalogClientResult.catalog = projectCatalogData
	return projectCatalogClientResult
}
