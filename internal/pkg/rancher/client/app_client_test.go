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
	backendClient "github.com/rancher/types/client/project/v3"
	"github.com/sirupsen/logrus"
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
				&types.ListOpts{
					Filters: map[string]interface{}{
						"name":        "existing-app",
						"namespaceId": "test-namespace-id",
					},
				},
			),
			wanted:  true,
			wantErr: false,
		},
		{
			name: "Not_Existing",
			client: notExistingAppClient(
				t,
				&types.ListOpts{
					Filters: map[string]interface{}{
						"name":        "existing-app",
						"namespaceId": "test-namespace-id",
					},
				},
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
			name: "Create",
			client: notExistingAppClient(
				t,
				nil,
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

func existingAppClient(t *testing.T, expectedListOpts *types.ListOpts) *appClient {
	const (
		projectID   = "test-project-id"
		projectName = "test-project-name"
		namespaceID = "test-namespace-id"
		namespace   = "test-namespace"
		clusterID   = "test-cluster-id"
		appName     = "test-app"
	)
	var (
		app         = projectModel.App{}
		testClients = stubs.CreateBackendStubs(t)
	)
	app.Name = appName

	appOperationsStub := stubs.CreateAppOperationsStub(t)
	appOperationsStub.DoList = func(opts *types.ListOpts) (*backendClient.AppCollection, error) {
		if !reflect.DeepEqual(expectedListOpts, opts) {
			return nil, fmt.Errorf("Unexpected ListOpts %v", opts)
		}
		return &backendClient.AppCollection{
			Data: []backendClient.App{
				backendClient.App{
					Name:        "existing-app",
					NamespaceId: "test-namespace-id",
				},
			},
		}, nil
	}
	testClients.ProjectClient.App = appOperationsStub
	result, err := newAppClient(
		"existing-app",
		"test-namespace",
		nil,
		testClients.ProjectClient,
		logrus.New().WithFields(logrus.Fields{}),
	)
	assert.Ok(t, err)
	appClientResult := result.(*appClient)
	appClientResult.namespaceID = "test-namespace-id"
	return appClientResult
}

func notExistingAppClient(t *testing.T, expectedListOpts *types.ListOpts) *appClient {
	const (
		projectID   = "test-project-id"
		projectName = "test-project-name"
		namespaceID = "test-namespace-id"
		namespace   = "test-namespace"
		clusterID   = "test-cluster-id"
		appName     = "test-app"
	)
	var (
		app         = projectModel.App{}
		testClients = stubs.CreateBackendStubs(t)
	)
	app.Name = appName

	appOperationsStub := stubs.CreateAppOperationsStub(t)
	appOperationsStub.DoList = func(opts *types.ListOpts) (*backendClient.AppCollection, error) {
		if !reflect.DeepEqual(expectedListOpts, opts) {
			return nil, fmt.Errorf("Unexpected ListOpts %v", opts)
		}
		return &backendClient.AppCollection{
			Data: []backendClient.App{},
		}, nil
	}
	appOperationsStub.DoCreate = func(app *backendClient.App) (*backendClient.App, error) {
		return app, nil
	}
	testClients.ProjectClient.App = appOperationsStub
	result, err := newAppClient(
		"existing-app",
		"test-namespace",
		nil,
		testClients.ProjectClient,
		logrus.New().WithFields(logrus.Fields{}),
	)
	assert.Ok(t, err)
	appClientResult := result.(*appClient)
	appClientResult.namespaceID = "test-namespace-id"
	return appClientResult
}
