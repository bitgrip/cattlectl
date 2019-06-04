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
	simpleSecretName = "simple-secret"
)

func Test_secretClient_Exists(t *testing.T) {
	tests := []struct {
		name      string
		client    *secretClient
		wanted    bool
		wantErr   bool
		wantedErr string
	}{
		{
			name: "Existing",
			client: existingSecretClient(
				t,
				&types.ListOpts{
					Filters: map[string]interface{}{
						"name":        "existing-secret",
						"namespaceId": "test-namespace-id",
					},
				},
			),
			wanted:  true,
			wantErr: false,
		},
		{
			name: "Not_Existing",
			client: notExistingSecretClient(
				t,
				&types.ListOpts{
					Filters: map[string]interface{}{
						"name":        "existing-secret",
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

func Test_secretClient_Create(t *testing.T) {
	tests := []struct {
		name      string
		client    *secretClient
		wantErr   bool
		wantedErr string
	}{
		{
			name: "Create",
			client: notExistingSecretClient(
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

func existingSecretClient(t *testing.T, expectedListOpts *types.ListOpts) *secretClient {
	const (
		projectID   = "test-project-id"
		projectName = "test-project-name"
		namespaceID = "test-namespace-id"
		namespace   = "test-namespace"
		clusterID   = "test-cluster-id"
		secretName  = "test-secret"
	)
	var (
		secret      = projectModel.ConfigMap{}
		testClients = stubs.CreateBackendStubs(t)
	)
	secret.Name = secretName

	secretOperationsStub := stubs.CreateSecretOperationsStub(t)
	secretOperationsStub.DoList = func(opts *types.ListOpts) (*backendProjectClient.SecretCollection, error) {
		if !reflect.DeepEqual(expectedListOpts, opts) {
			return nil, fmt.Errorf("Unexpected ListOpts %v", opts)
		}
		return &backendProjectClient.SecretCollection{
			Data: []backendProjectClient.Secret{
				backendProjectClient.Secret{
					Name:        "existing-secret",
					NamespaceId: "test-namespace-id",
				},
			},
		}, nil
	}
	testClients.ProjectClient.Secret = secretOperationsStub
	result, err := newSecretClient(
		"existing-secret",
		"test-namespace",
		&projectClient{
			resourceClient: resourceClient{
				name: projectName,
				id:   projectID,
			},
		},
		testClients.ProjectClient,
		logrus.New().WithFields(logrus.Fields{}),
	)
	assert.Ok(t, err)
	secretClientResult := result.(*secretClient)
	secretClientResult.namespaceID = "test-namespace-id"
	return secretClientResult
}

func notExistingSecretClient(t *testing.T, expectedListOpts *types.ListOpts) *secretClient {
	const (
		projectID   = "test-project-id"
		projectName = "test-project-name"
		namespaceID = "test-namespace-id"
		namespace   = "test-namespace"
		clusterID   = "test-cluster-id"
		secretName  = "test-secret"
	)
	var (
		secret      = projectModel.ConfigMap{}
		testClients = stubs.CreateBackendStubs(t)
	)
	secret.Name = secretName

	secretOperationsStub := stubs.CreateSecretOperationsStub(t)
	secretOperationsStub.DoList = func(opts *types.ListOpts) (*backendProjectClient.SecretCollection, error) {
		if !reflect.DeepEqual(expectedListOpts, opts) {
			return nil, fmt.Errorf("Unexpected ListOpts %v", opts)
		}
		return &backendProjectClient.SecretCollection{
			Data: []backendProjectClient.Secret{},
		}, nil
	}
	secretOperationsStub.DoCreate = func(secret *backendProjectClient.Secret) (*backendProjectClient.Secret, error) {
		return secret, nil
	}
	testClients.ProjectClient.Secret = secretOperationsStub
	result, err := newSecretClient(
		"existing-secret",
		"test-namespace",
		&projectClient{
			resourceClient: resourceClient{
				name: projectName,
				id:   projectID,
			},
		},
		testClients.ProjectClient,
		logrus.New().WithFields(logrus.Fields{}),
	)
	assert.Ok(t, err)
	secretClientResult := result.(*secretClient)
	secretClientResult.namespaceID = "test-namespace-id"
	return secretClientResult
}
