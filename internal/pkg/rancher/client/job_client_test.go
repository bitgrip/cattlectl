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
	simpleJobName = "simple-job"
)

func Test_jobClient_Exists(t *testing.T) {
	tests := []struct {
		name      string
		client    *jobClient
		wanted    bool
		wantErr   bool
		wantedErr string
	}{
		{
			name: "Existing",
			client: existingJobClient(
				t,
				&types.ListOpts{
					Filters: map[string]interface{}{
						"name":        "existing-job",
						"namespaceId": "test-namespace-id",
					},
				},
			),
			wanted:  true,
			wantErr: false,
		},
		{
			name: "Not_Existing",
			client: notExistingJobClient(
				t,
				&types.ListOpts{
					Filters: map[string]interface{}{
						"name":        "existing-job",
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

func Test_jobClient_Create(t *testing.T) {
	tests := []struct {
		name      string
		client    *jobClient
		wantErr   bool
		wantedErr string
	}{
		{
			name: "Create",
			client: notExistingJobClient(
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

func existingJobClient(t *testing.T, expectedListOpts *types.ListOpts) *jobClient {
	const (
		projectID   = "test-project-id"
		projectName = "test-project-name"
		namespaceID = "test-namespace-id"
		namespace   = "test-namespace"
		clusterID   = "test-cluster-id"
		jobName     = "test-job"
	)
	var (
		job           = projectModel.Job{}
		jobDescriptor = projectModel.JobDescriptor{
			Metadata: projectModel.WorkloadMetadata{
				ProjectName: projectName,
				ProjectID:   projectID,
				Namespace:   namespace,
				NamespaceID: namespaceID,
			},
		}
		testClients = stubs.CreateBackendStubs(t)
	)
	job.Name = jobName

	jobDescriptor.Spec = job

	jobOperationsStub := stubs.CreateJobOperationsStub(t)
	jobOperationsStub.DoList = func(opts *types.ListOpts) (*backendProjectClient.JobCollection, error) {
		if !reflect.DeepEqual(expectedListOpts, opts) {
			return nil, fmt.Errorf("Unexpected ListOpts %v", opts)
		}
		return &backendProjectClient.JobCollection{
			Data: []backendProjectClient.Job{
				backendProjectClient.Job{
					Name:        "existing-job",
					NamespaceId: "test-namespace-id",
				},
			},
		}, nil
	}
	testClients.ProjectClient.Job = jobOperationsStub
	projectClient := simpleProjectClient()
	projectClient._backendProjectClient = testClients.ProjectClient
	result, err := newJobClient(
		"existing-job",
		"test-namespace",
		projectClient,
		logrus.New().WithFields(logrus.Fields{}),
	)
	assert.Ok(t, err)
	jobClientResult := result.(*jobClient)
	jobClientResult.namespaceID = "test-namespace-id"
	return jobClientResult
}

func notExistingJobClient(t *testing.T, expectedListOpts *types.ListOpts) *jobClient {
	const (
		projectID   = "test-project-id"
		projectName = "test-project-name"
		namespaceID = "test-namespace-id"
		namespace   = "test-namespace"
		clusterID   = "test-cluster-id"
		jobName     = "test-job"
	)
	var (
		job           = projectModel.Job{}
		jobDescriptor = projectModel.JobDescriptor{
			Metadata: projectModel.WorkloadMetadata{
				ProjectName: projectName,
				ProjectID:   projectID,
				Namespace:   namespace,
				NamespaceID: namespaceID,
			},
		}
		testClients = stubs.CreateBackendStubs(t)
	)
	job.Name = jobName

	jobDescriptor.Spec = job

	jobOperationsStub := stubs.CreateJobOperationsStub(t)
	jobOperationsStub.DoList = func(opts *types.ListOpts) (*backendProjectClient.JobCollection, error) {
		if !reflect.DeepEqual(expectedListOpts, opts) {
			return nil, fmt.Errorf("Unexpected ListOpts %v", opts)
		}
		return &backendProjectClient.JobCollection{
			Data: []backendProjectClient.Job{},
		}, nil
	}
	jobOperationsStub.DoCreate = func(job *backendProjectClient.Job) (*backendProjectClient.Job, error) {
		return job, nil
	}
	testClients.ProjectClient.Job = jobOperationsStub
	projectClient := simpleProjectClient()
	projectClient._backendProjectClient = testClients.ProjectClient
	result, err := newJobClient(
		"existing-job",
		"test-namespace",
		projectClient,
		logrus.New().WithFields(logrus.Fields{}),
	)
	assert.Ok(t, err)
	jobClientResult := result.(*jobClient)
	jobClientResult.namespaceID = "test-namespace-id"
	return jobClientResult
}
