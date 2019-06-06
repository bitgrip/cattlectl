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
	simpleDockerCredentialName = "simple-docker-credential"
)

func Test_dockerCredentialClient_Exists(t *testing.T) {
	tests := []struct {
		name      string
		client    *dockerCredentialClient
		wanted    bool
		wantErr   bool
		wantedErr string
	}{
		{
			name: "Existing",
			client: existingDockerCredentialClient(
				t,
				&types.ListOpts{
					Filters: map[string]interface{}{
						"name":        simpleDockerCredentialName,
						"namespaceId": simpleNamespaceID,
					},
				},
			),
			wanted:  true,
			wantErr: false,
		},
		{
			name: "Not_Existing",
			client: notExistingDockerCredentialClient(
				t,
				&types.ListOpts{
					Filters: map[string]interface{}{
						"name":        simpleDockerCredentialName,
						"namespaceId": simpleNamespaceID,
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

func Test_dockerCredentialClient_Create(t *testing.T) {
	tests := []struct {
		name      string
		client    *dockerCredentialClient
		wantErr   bool
		wantedErr string
	}{
		{
			name: "Create",
			client: notExistingDockerCredentialClient(
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

func Test_dockerCredentialClient_Upgrade(t *testing.T) {
	tests := []struct {
		name      string
		client    *dockerCredentialClient
		wantErr   bool
		wantedErr string
	}{
		{
			name: "global",
			client: existingDockerCredentialClient(
				t,
				&types.ListOpts{
					Filters: map[string]interface{}{
						"name":        simpleDockerCredentialName,
						"namespaceId": simpleNamespaceID,
					},
				},
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

func existingDockerCredentialClient(t *testing.T, expectedListOpts *types.ListOpts) *dockerCredentialClient {
	var (
		dockerCredential = projectModel.DockerCredential{}
		testClients      = stubs.CreateBackendStubs(t)
	)
	dockerCredential.Name = simpleDockerCredentialName

	dockerCredentialOperationsStub := stubs.CreateDockerCredentialOperationsStub(t)
	dockerCredentialOperationsStub.DoList = func(opts *types.ListOpts) (*backendProjectClient.DockerCredentialCollection, error) {
		if !reflect.DeepEqual(expectedListOpts, opts) {
			return nil, fmt.Errorf("Unexpected ListOpts %v", opts)
		}
		return &backendProjectClient.DockerCredentialCollection{
			Data: []backendProjectClient.DockerCredential{
				backendProjectClient.DockerCredential{
					Name:        simpleDockerCredentialName,
					NamespaceId: simpleNamespaceID,
					Labels:      map[string]string{},
				},
			},
		}, nil
	}

	namespacedDockerCredentialOperationsStub := stubs.CreateNamespacedDockerCredentialOperationsStub(t)
	namespacedDockerCredentialOperationsStub.DoList = func(opts *types.ListOpts) (*backendProjectClient.NamespacedDockerCredentialCollection, error) {
		if !reflect.DeepEqual(expectedListOpts, opts) {
			return nil, fmt.Errorf("Unexpected ListOpts %v", opts)
		}
		return &backendProjectClient.NamespacedDockerCredentialCollection{
			Data: []backendProjectClient.NamespacedDockerCredential{
				backendProjectClient.NamespacedDockerCredential{
					Name:        simpleDockerCredentialName,
					NamespaceId: simpleNamespaceID,
					Labels:      map[string]string{},
				},
			},
		}, nil
	}
	namespacedDockerCredentialOperationsStub.DoReplace = func(existing *backendProjectClient.NamespacedDockerCredential) (*backendProjectClient.NamespacedDockerCredential, error) {
		return existing, nil
	}

	testClients.ProjectClient.DockerCredential = dockerCredentialOperationsStub
	testClients.ProjectClient.NamespacedDockerCredential = namespacedDockerCredentialOperationsStub
	result, err := newDockerCredentialClient(
		simpleDockerCredentialName,
		simpleNamespaceName,
		&projectClient{
			resourceClient: resourceClient{
				name: simpleProjectName,
				id:   simpleProjectID,
			},
			_backendProjectClient: testClients.ProjectClient,
		},
		logrus.New().WithFields(logrus.Fields{}),
	)
	assert.Ok(t, err)
	dockerCredentialClientResult := result.(*dockerCredentialClient)
	dockerCredentialClientResult.namespaceID = simpleNamespaceID
	return dockerCredentialClientResult
}

func notExistingDockerCredentialClient(t *testing.T, expectedListOpts *types.ListOpts) *dockerCredentialClient {
	var (
		dockerCredential = projectModel.DockerCredential{}
		testClients      = stubs.CreateBackendStubs(t)
	)
	dockerCredential.Name = simpleDockerCredentialName

	dockerCredentialOperationsStub := stubs.CreateDockerCredentialOperationsStub(t)
	dockerCredentialOperationsStub.DoList = func(opts *types.ListOpts) (*backendProjectClient.DockerCredentialCollection, error) {
		if !reflect.DeepEqual(expectedListOpts, opts) {
			return nil, fmt.Errorf("Unexpected ListOpts %v", opts)
		}
		return &backendProjectClient.DockerCredentialCollection{
			Data: []backendProjectClient.DockerCredential{},
		}, nil
	}
	dockerCredentialOperationsStub.DoCreate = func(dockerCredential *backendProjectClient.DockerCredential) (*backendProjectClient.DockerCredential, error) {
		return dockerCredential, nil
	}

	namespacedDockerCredentialOperationsStub := stubs.CreateNamespacedDockerCredentialOperationsStub(t)
	namespacedDockerCredentialOperationsStub.DoList = func(opts *types.ListOpts) (*backendProjectClient.NamespacedDockerCredentialCollection, error) {
		if !reflect.DeepEqual(expectedListOpts, opts) {
			return nil, fmt.Errorf("Unexpected ListOpts %v", opts)
		}
		return &backendProjectClient.NamespacedDockerCredentialCollection{
			Data: []backendProjectClient.NamespacedDockerCredential{},
		}, nil
	}
	namespacedDockerCredentialOperationsStub.DoCreate = func(dockerCredential *backendProjectClient.NamespacedDockerCredential) (*backendProjectClient.NamespacedDockerCredential, error) {
		return dockerCredential, nil
	}
	testClients.ProjectClient.DockerCredential = dockerCredentialOperationsStub
	testClients.ProjectClient.NamespacedDockerCredential = namespacedDockerCredentialOperationsStub
	result, err := newDockerCredentialClient(
		simpleDockerCredentialName,
		simpleNamespaceName,
		&projectClient{
			resourceClient: resourceClient{
				name: simpleProjectName,
				id:   simpleProjectID,
			},
			_backendProjectClient: testClients.ProjectClient,
		},
		logrus.New().WithFields(logrus.Fields{}),
	)
	assert.Ok(t, err)
	dockerCredentialClientResult := result.(*dockerCredentialClient)
	dockerCredentialClientResult.namespaceID = simpleNamespaceID
	return dockerCredentialClientResult
}
