// Copyright © 2018 Bitgrip <berlin@bitgrip.de>
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

package rancher

import (
	"testing"

	"github.com/bitgrip/cattlectl/internal/pkg/assert"
	projectModel "github.com/bitgrip/cattlectl/internal/pkg/rancher/project/model"
	"github.com/bitgrip/cattlectl/internal/pkg/rancher/stubs"
	"github.com/rancher/norman/types"
	clusterClient "github.com/rancher/types/client/cluster/v3"
	projectClient "github.com/rancher/types/client/project/v3"
	"github.com/sirupsen/logrus"
)

/*
HasApp(...)
  returns true if app exists
*/
func TestHasJob_JobExisting(t *testing.T) {
	const (
		projectID   = "test-project-id"
		projectName = "test-project-name"
		namespaceID = "test-namespace-id"
		namespace   = "test-namespace"
		clusterID   = "test-cluster-id"
		jobName     = "test-job"
	)
	var (
		actualListOpts *types.ListOpts
		clientConfig   = ClientConfig{}
		job            = projectModel.Job{}
		jobDescriptor  = projectModel.JobDescriptor{
			Metadata: projectModel.WorkloadMetadata{
				ProjectName: projectName,
				ProjectID:   projectID,
				Namespace:   namespace,
				NamespaceID: namespaceID,
			},
		}
		testClients = stubs.CreateTestClients(t)
	)
	job.Name = jobName

	jobDescriptor.Spec = job

	jobOperationsStub := stubs.CreateJobOperationsStub(t)
	jobOperationsStub.DoList = func(opts *types.ListOpts) (*projectClient.JobCollection, error) {
		actualListOpts = opts
		return &projectClient.JobCollection{
			Data: []projectClient.Job{
				projectClient.Job{
					Name:        jobName,
					NamespaceId: namespaceID,
				},
			},
		}, nil
	}
	testClients.ProjectClient.Job = jobOperationsStub

	client := &rancherClient{
		clientConfig:   clientConfig,
		projectId:      projectID,
		projectClient:  testClients.ProjectClient,
		clusterClient:  testClients.ClusterClient,
		appCache:       make(map[string]projectClient.App),
		namespaceCache: make(map[string]clusterClient.Namespace),
		logger:         logrus.WithField("test", true),
	}
	namespaceItem := clusterClient.Namespace{}
	namespaceItem.ID = namespaceID
	client.namespaceCache[namespace] = namespaceItem
	//Act
	result, err := client.HasJob(namespace, job)

	//Assert
	assert.Ok(t, err)
	assert.Equals(t, true, result)
	assert.Equals(t, &types.ListOpts{Filters: map[string]interface{}{"name": jobName, "namespaceId": namespaceID}}, actualListOpts)
}

/*
HasApp(...)
  returns true if local persistent volume exists
*/
func TestHasJob_JobNotExisting(t *testing.T) {
	const (
		projectID   = "test-project-id"
		projectName = "test-project-name"
		namespaceID = "test-namespace-id"
		namespace   = "test-namespace"
		clusterID   = "test-cluster-id"
		jobName     = "test-job"
	)
	var (
		actualListOpts *types.ListOpts
		clientConfig   = ClientConfig{}
		job            = projectModel.Job{}
		jobDescriptor  = projectModel.JobDescriptor{
			Metadata: projectModel.WorkloadMetadata{
				ProjectName: projectName,
				ProjectID:   projectID,
				Namespace:   namespace,
				NamespaceID: namespaceID,
			},
		}
		testClients = stubs.CreateTestClients(t)
	)
	job.Name = jobName

	jobDescriptor.Spec = job

	jobOperationsStub := stubs.CreateJobOperationsStub(t)
	jobOperationsStub.DoList = func(opts *types.ListOpts) (*projectClient.JobCollection, error) {
		actualListOpts = opts
		return &projectClient.JobCollection{
			Data: []projectClient.Job{},
		}, nil
	}
	testClients.ProjectClient.Job = jobOperationsStub

	client := rancherClient{
		clientConfig:   clientConfig,
		projectId:      projectID,
		projectClient:  testClients.ProjectClient,
		clusterClient:  testClients.ClusterClient,
		appCache:       make(map[string]projectClient.App),
		namespaceCache: make(map[string]clusterClient.Namespace),
		logger:         logrus.WithField("test", true),
	}
	namespaceItem := clusterClient.Namespace{}
	namespaceItem.ID = namespaceID
	client.namespaceCache[namespace] = namespaceItem
	//Act
	result, err := client.HasJob(namespace, job)

	//Assert
	assert.Ok(t, err)
	assert.Equals(t, false, result)
	assert.Equals(t, &types.ListOpts{Filters: map[string]interface{}{"name": jobName, "namespaceId": namespaceID}}, actualListOpts)
}

/*
CreateApp(...)
  uses project client to create app
*/
func TestCreateJob(t *testing.T) {

	const (
		projectID   = "test-project-id"
		projectName = "test-project-name"
		namespaceID = "test-namespace-id"
		namespace   = "test-namespace"
		clusterID   = "test-cluster-id"
		jobName     = "test-job"
	)
	var (
		actualCreateOpts *projectClient.Job
		clientConfig     = ClientConfig{}
		job              = projectModel.Job{}
		jobDescriptor    = projectModel.JobDescriptor{
			Metadata: projectModel.WorkloadMetadata{
				ProjectName: projectName,
				ProjectID:   projectID,
				Namespace:   namespace,
				NamespaceID: namespaceID,
			},
		}
		testClients = stubs.CreateTestClients(t)
	)
	job.Name = jobName

	jobDescriptor.Spec = job

	jobOperationsStub := stubs.CreateJobOperationsStub(t)
	jobOperationsStub.DoCreate = func(opts *projectClient.Job) (*projectClient.Job, error) {
		actualCreateOpts = opts
		return opts, nil
	}
	testClients.ProjectClient.Job = jobOperationsStub

	client := rancherClient{
		clientConfig:   clientConfig,
		projectId:      projectID,
		projectClient:  testClients.ProjectClient,
		clusterClient:  testClients.ClusterClient,
		appCache:       make(map[string]projectClient.App),
		namespaceCache: make(map[string]clusterClient.Namespace),
		logger:         logrus.WithField("test", true),
	}
	namespaceItem := clusterClient.Namespace{}
	namespaceItem.ID = namespaceID
	client.namespaceCache[namespace] = namespaceItem
	//Act
	err := client.CreateJob(namespace, job)

	//Assert
	assert.Ok(t, err)
	assert.Equals(t, &projectClient.Job{Name: jobName, NamespaceId: namespaceID}, actualCreateOpts)
}
