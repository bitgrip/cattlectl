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

package project

import (
	"testing"

	"github.com/bitgrip/cattlectl/internal/pkg/assert"
	"github.com/bitgrip/cattlectl/internal/pkg/rancher/clientstub"
	projectModel "github.com/bitgrip/cattlectl/internal/pkg/rancher/project/model"
)

func newTestJobDescriptor() projectModel.JobDescriptor {
	job := projectModel.Job{}
	job.Name = "test-job"
	return projectModel.JobDescriptor{
		Metadata: projectModel.WorkloadMetadata{
			ClusterName: "test-cluster",
			ProjectName: "test-project",
			Namespace:   "test-namespace",
		},
		Spec: job,
	}
}

func TestConvergeExistingJobDescriptor(t *testing.T) {
	jobDescriptor := newTestJobDescriptor()
	client := clientstub.NewClientStub(t).(clientstub.ClientStub)
	expectSetClusterForExisting(jobDescriptor.Metadata.ClusterName, "cluster-id", &client, t)
	expectSetProjectForExisting(jobDescriptor.Metadata.ProjectName, "project-id", &client, t)
	client.DoHasJob = func(namespace string, part projectModel.Job) (bool, error) {
		assert.Equals(t, jobDescriptor.Metadata.Namespace, namespace)
		return true, nil
	}
	converger := NewJobConverger(jobDescriptor)
	converger.Converge(client)
}

func TestConvergeNewJobDescriptor(t *testing.T) {
	jobDescriptor := newTestJobDescriptor()
	client := clientstub.NewClientStub(t).(clientstub.ClientStub)
	expectSetClusterForExisting(jobDescriptor.Metadata.ClusterName, "cluster-id", &client, t)
	expectSetProjectForExisting(jobDescriptor.Metadata.ProjectName, "project-id", &client, t)
	client.DoHasJob = func(namespace string, part projectModel.Job) (bool, error) {
		assert.Equals(t, jobDescriptor.Metadata.Namespace, namespace)
		return false, nil
	}
	client.DoCreateJob = func(namespace string, part projectModel.Job) error {
		assert.Equals(t, jobDescriptor.Metadata.Namespace, namespace)
		assert.Equals(t, jobDescriptor.Spec, part)
		return nil
	}
	converger := NewJobConverger(jobDescriptor)
	converger.Converge(client)
}
