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

func Test_cronJobClient_Exists(t *testing.T) {
	tests := []struct {
		name      string
		client    *cronJobClient
		wanted    bool
		wantErr   bool
		wantedErr string
	}{
		{
			name: "Existing",
			client: existingCronJobClient(
				t,
				&types.ListOpts{
					Filters: map[string]interface{}{
						"name":        "existing-cronJob",
						"namespaceId": "test-namespace-id",
					},
				},
			),
			wanted:  true,
			wantErr: false,
		},
		{
			name: "Not_Existing",
			client: notExistingCronJobClient(
				t,
				&types.ListOpts{
					Filters: map[string]interface{}{
						"name":        "existing-cronJob",
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

func Test_cronJobClient_Create(t *testing.T) {
	tests := []struct {
		name      string
		client    *cronJobClient
		wantErr   bool
		wantedErr string
	}{
		{
			name: "Create",
			client: notExistingCronJobClient(
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

func existingCronJobClient(t *testing.T, expectedListOpts *types.ListOpts) *cronJobClient {
	const (
		projectID   = "test-project-id"
		projectName = "test-project-name"
		namespaceID = "test-namespace-id"
		namespace   = "test-namespace"
		clusterID   = "test-cluster-id"
		cronJobName = "test-cronJob"
	)
	var (
		cronJob           = projectModel.CronJob{}
		cronJobDescriptor = projectModel.CronJobDescriptor{
			Metadata: projectModel.WorkloadMetadata{
				ProjectName: projectName,
				ProjectID:   projectID,
				Namespace:   namespace,
				NamespaceID: namespaceID,
			},
		}
		testClients = stubs.CreateBackendStubs(t)
	)
	cronJob.Name = cronJobName

	cronJobDescriptor.Spec = cronJob

	cronJobOperationsStub := stubs.CreateCronJobOperationsStub(t)
	cronJobOperationsStub.DoList = func(opts *types.ListOpts) (*backendClient.CronJobCollection, error) {
		if !reflect.DeepEqual(expectedListOpts, opts) {
			return nil, fmt.Errorf("Unexpected ListOpts %v", opts)
		}
		return &backendClient.CronJobCollection{
			Data: []backendClient.CronJob{
				backendClient.CronJob{
					Name:        "existing-cronJob",
					NamespaceId: "test-namespace-id",
				},
			},
		}, nil
	}
	testClients.ProjectClient.CronJob = cronJobOperationsStub
	result, err := newCronJobClient(
		"existing-cronJob",
		"test-namespace",
		nil,
		testClients.ProjectClient,
		logrus.New().WithFields(logrus.Fields{}),
	)
	assert.Ok(t, err)
	cronJobClientResult := result.(*cronJobClient)
	cronJobClientResult.namespaceID = "test-namespace-id"
	return cronJobClientResult
}

func notExistingCronJobClient(t *testing.T, expectedListOpts *types.ListOpts) *cronJobClient {
	const (
		projectID   = "test-project-id"
		projectName = "test-project-name"
		namespaceID = "test-namespace-id"
		namespace   = "test-namespace"
		clusterID   = "test-cluster-id"
		cronJobName = "test-cronJob"
	)
	var (
		cronJob           = projectModel.CronJob{}
		cronJobDescriptor = projectModel.CronJobDescriptor{
			Metadata: projectModel.WorkloadMetadata{
				ProjectName: projectName,
				ProjectID:   projectID,
				Namespace:   namespace,
				NamespaceID: namespaceID,
			},
		}
		testClients = stubs.CreateBackendStubs(t)
	)
	cronJob.Name = cronJobName

	cronJobDescriptor.Spec = cronJob

	cronJobOperationsStub := stubs.CreateCronJobOperationsStub(t)
	cronJobOperationsStub.DoList = func(opts *types.ListOpts) (*backendClient.CronJobCollection, error) {
		if !reflect.DeepEqual(expectedListOpts, opts) {
			return nil, fmt.Errorf("Unexpected ListOpts %v", opts)
		}
		return &backendClient.CronJobCollection{
			Data: []backendClient.CronJob{},
		}, nil
	}
	cronJobOperationsStub.DoCreate = func(cronJob *backendClient.CronJob) (*backendClient.CronJob, error) {
		return cronJob, nil
	}
	testClients.ProjectClient.CronJob = cronJobOperationsStub
	result, err := newCronJobClient(
		"existing-cronJob",
		"test-namespace",
		nil,
		testClients.ProjectClient,
		logrus.New().WithFields(logrus.Fields{}),
	)
	assert.Ok(t, err)
	cronJobClientResult := result.(*cronJobClient)
	cronJobClientResult.namespaceID = "test-namespace-id"
	return cronJobClientResult
}
