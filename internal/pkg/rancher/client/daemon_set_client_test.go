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
	projectModel "github.com/bitgrip/cattlectl/internal/pkg/rancher/cluster/project/model"
	"github.com/bitgrip/cattlectl/internal/pkg/rancher/stubs"
	"github.com/rancher/norman/types"
	backendProjectClient "github.com/rancher/types/client/project/v3"
	"github.com/sirupsen/logrus"
)

const (
	simpleDaemonSetName = "simple-daemon-set"
)

func Test_daemonSetClient_Exists(t *testing.T) {
	tests := []struct {
		name      string
		client    *daemonSetClient
		wanted    bool
		wantErr   bool
		wantedErr string
	}{
		{
			name: "Existing",
			client: existingDaemonSetClient(
				t,
				&types.ListOpts{
					Filters: map[string]interface{}{
						"name":        "existing-daemonSet",
						"namespaceId": "test-namespace-id",
					},
				},
			),
			wanted:  true,
			wantErr: false,
		},
		{
			name: "Not_Existing",
			client: notExistingDaemonSetClient(
				t,
				&types.ListOpts{
					Filters: map[string]interface{}{
						"name":        "existing-daemonSet",
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

func Test_daemonSetClient_Create(t *testing.T) {
	tests := []struct {
		name      string
		client    *daemonSetClient
		wantErr   bool
		wantedErr string
	}{
		{
			name: "Create",
			client: notExistingDaemonSetClient(
				t,
				nil,
			),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.client.Create(false)
			if tt.wantErr {
				assert.NotOk(t, err, tt.wantedErr)
			} else {
				assert.Ok(t, err)
			}
		})
	}
}

func existingDaemonSetClient(t *testing.T, expectedListOpts *types.ListOpts) *daemonSetClient {
	const (
		projectID     = "test-project-id"
		projectName   = "test-project-name"
		namespaceID   = "test-namespace-id"
		namespace     = "test-namespace"
		clusterID     = "test-cluster-id"
		daemonSetName = "test-daemonSet"
	)
	var (
		daemonSet           = projectModel.DaemonSet{}
		daemonSetDescriptor = projectModel.DaemonSetDescriptor{
			Metadata: projectModel.WorkloadMetadata{
				ProjectName: projectName,
				ProjectID:   projectID,
				Namespace:   namespace,
				NamespaceID: namespaceID,
			},
		}
		testClients = stubs.CreateBackendStubs(t)
	)
	daemonSet.Name = daemonSetName

	daemonSetDescriptor.Spec = daemonSet

	daemonSetOperationsStub := stubs.CreateDaemonSetOperationsStub(t)
	daemonSetOperationsStub.DoList = func(opts *types.ListOpts) (*backendProjectClient.DaemonSetCollection, error) {
		if !reflect.DeepEqual(expectedListOpts, opts) {
			return nil, fmt.Errorf("Unexpected ListOpts %v", opts)
		}
		return &backendProjectClient.DaemonSetCollection{
			Data: []backendProjectClient.DaemonSet{
				backendProjectClient.DaemonSet{
					Name:        "existing-daemonSet",
					NamespaceId: "test-namespace-id",
				},
			},
		}, nil
	}
	testClients.ProjectClient.DaemonSet = daemonSetOperationsStub
	projectClient := simpleProjectClient()
	projectClient._backendProjectClient = testClients.ProjectClient
	result, err := newDaemonSetClient(
		"existing-daemonSet",
		"test-namespace",
		projectClient,
		logrus.New().WithFields(logrus.Fields{}),
	)
	assert.Ok(t, err)
	daemonSetClientResult := result.(*daemonSetClient)
	daemonSetClientResult.namespaceID = "test-namespace-id"
	return daemonSetClientResult
}

func notExistingDaemonSetClient(t *testing.T, expectedListOpts *types.ListOpts) *daemonSetClient {
	const (
		projectID     = "test-project-id"
		projectName   = "test-project-name"
		namespaceID   = "test-namespace-id"
		namespace     = "test-namespace"
		clusterID     = "test-cluster-id"
		daemonSetName = "test-daemonSet"
	)
	var (
		daemonSet           = projectModel.DaemonSet{}
		daemonSetDescriptor = projectModel.DaemonSetDescriptor{
			Metadata: projectModel.WorkloadMetadata{
				ProjectName: projectName,
				ProjectID:   projectID,
				Namespace:   namespace,
				NamespaceID: namespaceID,
			},
		}
		testClients = stubs.CreateBackendStubs(t)
	)
	daemonSet.Name = daemonSetName

	daemonSetDescriptor.Spec = daemonSet

	daemonSetOperationsStub := stubs.CreateDaemonSetOperationsStub(t)
	daemonSetOperationsStub.DoList = func(opts *types.ListOpts) (*backendProjectClient.DaemonSetCollection, error) {
		if !reflect.DeepEqual(expectedListOpts, opts) {
			return nil, fmt.Errorf("Unexpected ListOpts %v", opts)
		}
		return &backendProjectClient.DaemonSetCollection{
			Data: []backendProjectClient.DaemonSet{},
		}, nil
	}
	daemonSetOperationsStub.DoCreate = func(daemonSet *backendProjectClient.DaemonSet) (*backendProjectClient.DaemonSet, error) {
		return daemonSet, nil
	}
	testClients.ProjectClient.DaemonSet = daemonSetOperationsStub
	projectClient := simpleProjectClient()
	projectClient._backendProjectClient = testClients.ProjectClient
	result, err := newDaemonSetClient(
		"existing-daemonSet",
		"test-namespace",
		projectClient,
		logrus.New().WithFields(logrus.Fields{}),
	)
	assert.Ok(t, err)
	daemonSetClientResult := result.(*daemonSetClient)
	daemonSetClientResult.namespaceID = "test-namespace-id"
	return daemonSetClientResult
}
