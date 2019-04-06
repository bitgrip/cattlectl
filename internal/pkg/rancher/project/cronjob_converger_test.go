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

func newTestCronJobDescriptor() projectModel.CronJobDescriptor {
	cronJob := projectModel.CronJob{}
	cronJob.Name = "test-cron-job"
	return projectModel.CronJobDescriptor{
		Metadata: projectModel.WorkloadMetadata{
			ClusterName: "test-cluster",
			ProjectName: "test-project",
			Namespace:   "test-namespace",
		},
		Spec: cronJob,
	}
}

func TestConvergeExistingCronJobDescriptor(t *testing.T) {
	cronJobDescriptor := newTestCronJobDescriptor()
	client := clientstub.NewClientStub(t).(clientstub.ClientStub)
	expectSetClusterForExisting(cronJobDescriptor.Metadata.ClusterName, "cluster-id", &client, t)
	expectSetProjectForExisting(cronJobDescriptor.Metadata.ProjectName, "project-id", &client, t)
	client.DoHasCronJob = func(namespace string, part projectModel.CronJob) (bool, error) {
		assert.Equals(t, cronJobDescriptor.Metadata.Namespace, namespace)
		assert.Equals(t, cronJobDescriptor.Spec.Name, part.Name)
		return true, nil
	}
	converger := NewCronJobConverger(cronJobDescriptor)
	converger.Converge(&client)
}

func TestConvergeNewCronJobDescriptor(t *testing.T) {
	cronJobDescriptor := newTestCronJobDescriptor()
	client := clientstub.NewClientStub(t).(clientstub.ClientStub)
	expectSetClusterForExisting(cronJobDescriptor.Metadata.ClusterName, "cluster-id", &client, t)
	expectSetProjectForExisting(cronJobDescriptor.Metadata.ProjectName, "project-id", &client, t)
	client.DoHasCronJob = func(namespace string, part projectModel.CronJob) (bool, error) {
		assert.Equals(t, cronJobDescriptor.Metadata.Namespace, namespace)
		return false, nil
	}
	client.DoCreateCronJob = func(namespace string, part projectModel.CronJob) error {
		assert.Equals(t, cronJobDescriptor.Metadata.Namespace, namespace)
		assert.Equals(t, cronJobDescriptor.Spec, part)
		return nil
	}
	converger := NewCronJobConverger(cronJobDescriptor)
	converger.Converge(client)
}
