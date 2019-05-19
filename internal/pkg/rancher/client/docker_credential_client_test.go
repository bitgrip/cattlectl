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
						"name":        "existing-dockerCredential",
						"namespaceId": "test-namespace-id",
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
						"name":        "existing-dockerCredential",
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

func existingDockerCredentialClient(t *testing.T, expectedListOpts *types.ListOpts) *dockerCredentialClient {
	const (
		projectID            = "test-project-id"
		projectName          = "test-project-name"
		namespaceID          = "test-namespace-id"
		namespace            = "test-namespace"
		clusterID            = "test-cluster-id"
		dockerCredentialName = "test-dockerCredential"
	)
	var (
		dockerCredential = projectModel.DockerCredential{}
		testClients      = stubs.CreateBackendStubs(t)
	)
	dockerCredential.Name = dockerCredentialName

	dockerCredentialOperationsStub := stubs.CreateDockerCredentialOperationsStub(t)
	dockerCredentialOperationsStub.DoList = func(opts *types.ListOpts) (*backendProjectClient.DockerCredentialCollection, error) {
		if !reflect.DeepEqual(expectedListOpts, opts) {
			return nil, fmt.Errorf("Unexpected ListOpts %v", opts)
		}
		return &backendProjectClient.DockerCredentialCollection{
			Data: []backendProjectClient.DockerCredential{
				backendProjectClient.DockerCredential{
					Name:        "existing-dockerCredential",
					NamespaceId: "test-namespace-id",
				},
			},
		}, nil
	}
	testClients.ProjectClient.DockerCredential = dockerCredentialOperationsStub
	result, err := newDockerCredentialClient(
		"existing-dockerCredential",
		"test-namespace",
		nil,
		testClients.ProjectClient,
		logrus.New().WithFields(logrus.Fields{}),
	)
	assert.Ok(t, err)
	dockerCredentialClientResult := result.(*dockerCredentialClient)
	dockerCredentialClientResult.namespaceID = "test-namespace-id"
	return dockerCredentialClientResult
}

func notExistingDockerCredentialClient(t *testing.T, expectedListOpts *types.ListOpts) *dockerCredentialClient {
	const (
		projectID            = "test-project-id"
		projectName          = "test-project-name"
		namespaceID          = "test-namespace-id"
		namespace            = "test-namespace"
		clusterID            = "test-cluster-id"
		dockerCredentialName = "test-dockerCredential"
	)
	var (
		dockerCredential = projectModel.DockerCredential{}
		testClients      = stubs.CreateBackendStubs(t)
	)
	dockerCredential.Name = dockerCredentialName

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
	testClients.ProjectClient.DockerCredential = dockerCredentialOperationsStub
	result, err := newDockerCredentialClient(
		"existing-dockerCredential",
		"test-namespace",
		&projectClient{},
		testClients.ProjectClient,
		logrus.New().WithFields(logrus.Fields{}),
	)
	assert.Ok(t, err)
	dockerCredentialClientResult := result.(*dockerCredentialClient)
	dockerCredentialClientResult.namespaceID = "test-namespace-id"
	return dockerCredentialClientResult
}
