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

func Test_statefulSetClient_Exists(t *testing.T) {
	tests := []struct {
		name      string
		client    *statefulSetClient
		wanted    bool
		wantErr   bool
		wantedErr string
	}{
		{
			name: "Existing",
			client: existingStatefulSetClient(
				t,
				&types.ListOpts{
					Filters: map[string]interface{}{
						"name":        "existing-statefulSet",
						"namespaceId": "test-namespace-id",
					},
				},
			),
			wanted:  true,
			wantErr: false,
		},
		{
			name: "Not_Existing",
			client: notExistingStatefulSetClient(
				t,
				&types.ListOpts{
					Filters: map[string]interface{}{
						"name":        "existing-statefulSet",
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

func Test_statefulSetClient_Create(t *testing.T) {
	tests := []struct {
		name      string
		client    *statefulSetClient
		wantErr   bool
		wantedErr string
	}{
		{
			name: "Create",
			client: notExistingStatefulSetClient(
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

func existingStatefulSetClient(t *testing.T, expectedListOpts *types.ListOpts) *statefulSetClient {
	const (
		projectID       = "test-project-id"
		projectName     = "test-project-name"
		namespaceID     = "test-namespace-id"
		namespace       = "test-namespace"
		clusterID       = "test-cluster-id"
		statefulSetName = "test-statefulSet"
	)
	var (
		statefulSet           = projectModel.StatefulSet{}
		statefulSetDescriptor = projectModel.StatefulSetDescriptor{
			Metadata: projectModel.WorkloadMetadata{
				ProjectName: projectName,
				ProjectID:   projectID,
				Namespace:   namespace,
				NamespaceID: namespaceID,
			},
		}
		testClients = stubs.CreateBackendStubs(t)
	)
	statefulSet.Name = statefulSetName

	statefulSetDescriptor.Spec = statefulSet

	statefulSetOperationsStub := stubs.CreateStatefulSetOperationsStub(t)
	statefulSetOperationsStub.DoList = func(opts *types.ListOpts) (*backendClient.StatefulSetCollection, error) {
		if !reflect.DeepEqual(expectedListOpts, opts) {
			return nil, fmt.Errorf("Unexpected ListOpts %v", opts)
		}
		return &backendClient.StatefulSetCollection{
			Data: []backendClient.StatefulSet{
				backendClient.StatefulSet{
					Name:        "existing-statefulSet",
					NamespaceId: "test-namespace-id",
				},
			},
		}, nil
	}
	testClients.ProjectClient.StatefulSet = statefulSetOperationsStub
	result, err := newStatefulSetClient(
		"existing-statefulSet",
		"test-namespace",
		nil,
		testClients.ProjectClient,
		logrus.New().WithFields(logrus.Fields{}),
	)
	assert.Ok(t, err)
	statefulSetClientResult := result.(*statefulSetClient)
	statefulSetClientResult.namespaceID = "test-namespace-id"
	return statefulSetClientResult
}

func notExistingStatefulSetClient(t *testing.T, expectedListOpts *types.ListOpts) *statefulSetClient {
	const (
		projectID       = "test-project-id"
		projectName     = "test-project-name"
		namespaceID     = "test-namespace-id"
		namespace       = "test-namespace"
		clusterID       = "test-cluster-id"
		statefulSetName = "test-statefulSet"
	)
	var (
		statefulSet           = projectModel.StatefulSet{}
		statefulSetDescriptor = projectModel.StatefulSetDescriptor{
			Metadata: projectModel.WorkloadMetadata{
				ProjectName: projectName,
				ProjectID:   projectID,
				Namespace:   namespace,
				NamespaceID: namespaceID,
			},
		}
		testClients = stubs.CreateBackendStubs(t)
	)
	statefulSet.Name = statefulSetName

	statefulSetDescriptor.Spec = statefulSet

	statefulSetOperationsStub := stubs.CreateStatefulSetOperationsStub(t)
	statefulSetOperationsStub.DoList = func(opts *types.ListOpts) (*backendClient.StatefulSetCollection, error) {
		if !reflect.DeepEqual(expectedListOpts, opts) {
			return nil, fmt.Errorf("Unexpected ListOpts %v", opts)
		}
		return &backendClient.StatefulSetCollection{
			Data: []backendClient.StatefulSet{},
		}, nil
	}
	statefulSetOperationsStub.DoCreate = func(statefulSet *backendClient.StatefulSet) (*backendClient.StatefulSet, error) {
		return statefulSet, nil
	}
	testClients.ProjectClient.StatefulSet = statefulSetOperationsStub
	result, err := newStatefulSetClient(
		"existing-statefulSet",
		"test-namespace",
		nil,
		testClients.ProjectClient,
		logrus.New().WithFields(logrus.Fields{}),
	)
	assert.Ok(t, err)
	statefulSetClientResult := result.(*statefulSetClient)
	statefulSetClientResult.namespaceID = "test-namespace-id"
	return statefulSetClientResult
}
