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

func Test_deploymentClient_Exists(t *testing.T) {
	tests := []struct {
		name      string
		client    *deploymentClient
		wanted    bool
		wantErr   bool
		wantedErr string
	}{
		{
			name: "Existing",
			client: existingDeploymentClient(
				t,
				&types.ListOpts{
					Filters: map[string]interface{}{
						"name":        "existing-deployment",
						"namespaceId": "test-namespace-id",
					},
				},
			),
			wanted:  true,
			wantErr: false,
		},
		{
			name: "Not_Existing",
			client: notExistingDeploymentClient(
				t,
				&types.ListOpts{
					Filters: map[string]interface{}{
						"name":        "existing-deployment",
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

func Test_deploymentClient_Create(t *testing.T) {
	tests := []struct {
		name      string
		client    *deploymentClient
		wantErr   bool
		wantedErr string
	}{
		{
			name: "Create",
			client: notExistingDeploymentClient(
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

func existingDeploymentClient(t *testing.T, expectedListOpts *types.ListOpts) *deploymentClient {
	const (
		projectID      = "test-project-id"
		projectName    = "test-project-name"
		namespaceID    = "test-namespace-id"
		namespace      = "test-namespace"
		clusterID      = "test-cluster-id"
		deploymentName = "test-deployment"
	)
	var (
		deployment           = projectModel.Deployment{}
		deploymentDescriptor = projectModel.DeploymentDescriptor{
			Metadata: projectModel.WorkloadMetadata{
				ProjectName: projectName,
				ProjectID:   projectID,
				Namespace:   namespace,
				NamespaceID: namespaceID,
			},
		}
		testClients = stubs.CreateBackendStubs(t)
	)
	deployment.Name = deploymentName

	deploymentDescriptor.Spec = deployment

	deploymentOperationsStub := stubs.CreateDeploymentOperationsStub(t)
	deploymentOperationsStub.DoList = func(opts *types.ListOpts) (*backendClient.DeploymentCollection, error) {
		if !reflect.DeepEqual(expectedListOpts, opts) {
			return nil, fmt.Errorf("Unexpected ListOpts %v", opts)
		}
		return &backendClient.DeploymentCollection{
			Data: []backendClient.Deployment{
				backendClient.Deployment{
					Name:        "existing-deployment",
					NamespaceId: "test-namespace-id",
				},
			},
		}, nil
	}
	testClients.ProjectClient.Deployment = deploymentOperationsStub
	result, err := newDeploymentClient(
		"existing-deployment",
		"test-namespace",
		nil,
		testClients.ProjectClient,
		logrus.New().WithFields(logrus.Fields{}),
	)
	assert.Ok(t, err)
	deploymentClientResult := result.(*deploymentClient)
	deploymentClientResult.namespaceID = "test-namespace-id"
	return deploymentClientResult
}

func notExistingDeploymentClient(t *testing.T, expectedListOpts *types.ListOpts) *deploymentClient {
	const (
		projectID      = "test-project-id"
		projectName    = "test-project-name"
		namespaceID    = "test-namespace-id"
		namespace      = "test-namespace"
		clusterID      = "test-cluster-id"
		deploymentName = "test-deployment"
	)
	var (
		deployment           = projectModel.Deployment{}
		deploymentDescriptor = projectModel.DeploymentDescriptor{
			Metadata: projectModel.WorkloadMetadata{
				ProjectName: projectName,
				ProjectID:   projectID,
				Namespace:   namespace,
				NamespaceID: namespaceID,
			},
		}
		testClients = stubs.CreateBackendStubs(t)
	)
	deployment.Name = deploymentName

	deploymentDescriptor.Spec = deployment

	deploymentOperationsStub := stubs.CreateDeploymentOperationsStub(t)
	deploymentOperationsStub.DoList = func(opts *types.ListOpts) (*backendClient.DeploymentCollection, error) {
		if !reflect.DeepEqual(expectedListOpts, opts) {
			return nil, fmt.Errorf("Unexpected ListOpts %v", opts)
		}
		return &backendClient.DeploymentCollection{
			Data: []backendClient.Deployment{},
		}, nil
	}
	deploymentOperationsStub.DoCreate = func(deployment *backendClient.Deployment) (*backendClient.Deployment, error) {
		return deployment, nil
	}
	testClients.ProjectClient.Deployment = deploymentOperationsStub
	result, err := newDeploymentClient(
		"existing-deployment",
		"test-namespace",
		nil,
		testClients.ProjectClient,
		logrus.New().WithFields(logrus.Fields{}),
	)
	assert.Ok(t, err)
	deploymentClientResult := result.(*deploymentClient)
	deploymentClientResult.namespaceID = "test-namespace-id"
	return deploymentClientResult
}
