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
	projectModel "github.com/bitgrip/cattlectl/internal/pkg/rancher/project/model"
	"github.com/bitgrip/cattlectl/internal/pkg/rancher/stubs"
	"github.com/rancher/norman/types"
	backendProjectClient "github.com/rancher/types/client/project/v3"
	"github.com/sirupsen/logrus"
)

const (
	simpleAppName = "simple-app"
	simpleCatalog = "simple-catalog"
)

func Test_appClient_Exists(t *testing.T) {
	tests := []struct {
		name      string
		client    *appClient
		wanted    bool
		wantErr   bool
		wantedErr string
	}{
		{
			name: "Existing",
			client: existingAppClient(
				t,
				simpleAppName,
				simpleNamespaceName,
				simpleCatalog,
				"1.1.1",
				map[string]string{},
				map[string]string{},
				"",
				"",
			),
			wanted:  true,
			wantErr: false,
		},
		{
			name: "Not_Existing",
			client: notExistingAppClient(
				t,
				simpleClusterID,
				simpleProjectID,
				simpleAppName,
				simpleNamespaceName,
				simpleCatalog,
				globalCatalogType,
				"1.1.1",
				map[string]string{},
				"",
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

func Test_appClient_Create(t *testing.T) {
	tests := []struct {
		name      string
		client    *appClient
		wantErr   bool
		wantedErr string
	}{
		{
			name: "simple",
			client: notExistingAppClient(
				t,
				simpleClusterID,
				simpleProjectID,
				simpleAppName,
				simpleNamespaceName,
				simpleCatalog,
				globalCatalogType,
				"1.1.1",
				map[string]string{},
				"",
			),
			wantErr: false,
		},
		{
			name: "cluster-catalog",
			client: notExistingAppClient(
				t,
				simpleClusterID,
				simpleProjectID,
				simpleAppName,
				simpleNamespaceName,
				simpleCatalog,
				clusterCatalogType,
				"1.1.1",
				map[string]string{},
				"",
			),
			wantErr: false,
		},
		{
			name: "project-catalog",
			client: notExistingAppClient(
				t,
				simpleClusterID,
				simpleProjectID,
				simpleAppName,
				simpleNamespaceName,
				simpleCatalog,
				projectCatalogType,
				"1.1.1",
				map[string]string{},
				"",
			),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.client.Create()
			if tt.wantErr {
				assert.NotOk(t, err, tt.wantedErr)
			} else {
				assert.Ok(t, err)
			}
		})
	}
}

func Test_appClient_Upgrade(t *testing.T) {
	tests := []struct {
		name      string
		client    *appClient
		wantErr   bool
		wantedErr string
	}{
		{
			name: "upgrade-answers",
			client: existingAppClient(
				t,
				simpleAppName,
				simpleNamespaceName,
				simpleCatalog,
				"1.1.1",
				map[string]string{
					"key": "value",
				},
				map[string]string{
					"key": "changed-value",
				},
				"",
				"",
			),
			wantErr: false,
		},
		{
			name: "unchanged-answers",
			client: existingAppClient(
				t,
				simpleAppName,
				simpleNamespaceName,
				simpleCatalog,
				"1.1.1",
				map[string]string{
					"key": "value",
				},
				map[string]string{
					"key": "value",
				},
				"",
				"",
			),
			wantErr: false,
		},
		{
			name: "upgrade-values-yaml",
			client: existingAppClient(
				t,
				simpleAppName,
				simpleNamespaceName,
				simpleCatalog,
				"1.1.1",
				nil,
				nil,
				`---
key: value
				`,
				`---
key: changed-value
				`,
			),
			wantErr: false,
		},
		{
			name: "unchanged-values-yaml",
			client: existingAppClient(
				t,
				simpleAppName,
				simpleNamespaceName,
				simpleCatalog,
				"1.1.1",
				nil,
				nil,
				`---
key: value
				`,
				`---
key: value
				`,
			),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.client.Upgrade()
			if tt.wantErr {
				assert.NotOk(t, err, tt.wantedErr)
			} else {
				assert.Ok(t, err)
			}
		})
	}
}

func existingAppClient(t *testing.T, name, namespace, catalog, version string, answers, changedAnswers map[string]string, valuesYaml, changedValuesYaml string) *appClient {
	testClients := stubs.CreateBackendStubs(t)
	expectedListOpts := &types.ListOpts{
		Filters: map[string]interface{}{
			"name": name,
		},
	}
	expectedApp := &backendProjectClient.App{
		Name:            name,
		ExternalID:      fmt.Sprintf("catalog://?catalog=%v&template=%v&version=%v", catalog, name, version),
		TargetNamespace: namespace,
		Answers:         answers,
		ValuesYaml:      valuesYaml,
	}
	expectedUpgradeConfig := &backendProjectClient.AppUpgradeConfig{
		ExternalID: fmt.Sprintf("catalog://?catalog=%v&template=%v&version=%v", catalog, name, version),
		ValuesYaml: changedValuesYaml,
	}
	if changedValuesYaml == "" {
		expectedUpgradeConfig.Answers = changedAnswers
	}

	appOperationsStub := stubs.CreateAppOperationsStub(t)
	appOperationsStub.DoList = func(opts *types.ListOpts) (*backendProjectClient.AppCollection, error) {
		if !reflect.DeepEqual(expectedListOpts, opts) {
			return nil, fmt.Errorf("Unexpected ListOpts %v", opts)
		}
		return &backendProjectClient.AppCollection{
			Data: []backendProjectClient.App{
				backendProjectClient.App{
					Name:            name,
					ExternalID:      fmt.Sprintf("catalog://?catalog=%v&template=%v&version=%v", catalog, name, version),
					TargetNamespace: namespace,
					Answers:         answers,
					ValuesYaml:      valuesYaml,
				},
			},
		}, nil
	}
	if !reflect.DeepEqual(answers, changedAnswers) || valuesYaml != changedValuesYaml {
		appOperationsStub.DoActionUpgrade = func(resource *backendProjectClient.App, input *backendProjectClient.AppUpgradeConfig) error {
			if !reflect.DeepEqual(expectedApp, resource) {
				return fmt.Errorf("Unexpected target App\n%v\n%v", expectedApp, resource)
			}
			if !reflect.DeepEqual(expectedUpgradeConfig, input) {
				return fmt.Errorf("Unexpected upgrade config\n%v\n%v", expectedUpgradeConfig, input)
			}
			return nil
		}
	}
	testClients.ProjectClient.App = appOperationsStub
	result, err := newAppClientWithData(
		projectModel.App{
			Name:       name,
			Namespace:  namespace,
			Catalog:    catalog,
			Version:    version,
			Chart:      name,
			Answers:    changedAnswers,
			ValuesYaml: changedValuesYaml,
		},
		&projectClient{
			_backendProjectClient: testClients.ProjectClient,
		},
		logrus.New().WithFields(logrus.Fields{}),
	)
	assert.Ok(t, err)
	appClientResult := result.(*appClient)
	return appClientResult
}

func notExistingAppClient(t *testing.T, clusterID, projectID, name, namespace, catalog, catalogType, version string, answers map[string]string, valuesYaml string) *appClient {
	testClients := stubs.CreateBackendStubs(t)
	externalID := fmt.Sprintf(globalCatalogExternalID, catalog, name, version)
	switch catalogType {
	case projectCatalogType:
		externalID = fmt.Sprintf(projectCatalogExternalID, projectID, catalog, name, version)
	case clusterCatalogType:
		externalID = fmt.Sprintf(clusterCatalogExternalID, clusterID, catalog, name, version)
	}
	expectedListOpts := &types.ListOpts{
		Filters: map[string]interface{}{
			"name": name,
		},
	}
	expectedApp := &backendProjectClient.App{
		Name:            name,
		ExternalID:      externalID,
		TargetNamespace: namespace,
		Answers:         answers,
		ValuesYaml:      valuesYaml,
	}
	appOperationsStub := stubs.CreateAppOperationsStub(t)
	appOperationsStub.DoList = func(opts *types.ListOpts) (*backendProjectClient.AppCollection, error) {
		if !reflect.DeepEqual(expectedListOpts, opts) {
			return nil, fmt.Errorf("Unexpected ListOpts %v", opts)
		}
		return &backendProjectClient.AppCollection{
			Data: []backendProjectClient.App{},
		}, nil
	}
	appOperationsStub.DoCreate = func(app *backendProjectClient.App) (*backendProjectClient.App, error) {
		if !reflect.DeepEqual(expectedApp, app) {
			return nil, fmt.Errorf("Unexpected App %v", app)
		}
		return app, nil
	}
	testClients.ProjectClient.App = appOperationsStub
	result, err := newAppClientWithData(
		projectModel.App{
			Name:        name,
			Namespace:   namespace,
			Catalog:     catalog,
			CatalogType: catalogType,
			Version:     version,
			Chart:       name,
			Answers:     answers,
			ValuesYaml:  valuesYaml,
		},
		&projectClient{
			resourceClient: resourceClient{
				name: simpleProjectName,
				id:   simpleProjectID,
			},
			clusterClient: &clusterClient{
				resourceClient: resourceClient{
					name: simpleClusterName,
					id:   simpleClusterID,
				},
			},
			_backendProjectClient: testClients.ProjectClient,
		},
		logrus.New().WithFields(logrus.Fields{}),
	)
	assert.Ok(t, err)
	appClientResult := result.(*appClient)
	return appClientResult
}
