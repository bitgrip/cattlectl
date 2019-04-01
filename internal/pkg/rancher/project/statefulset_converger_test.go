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

func newTestStatefulSetDescriptor() projectModel.StatefulSetDescriptor {
	statefulSet := projectModel.StatefulSet{}
	statefulSet.Name = "test-statefulset"
	return projectModel.StatefulSetDescriptor{
		Metadata: projectModel.WorkloadMetadata{
			ClusterName: "test-cluster",
			ProjectName: "test-project",
			Namespace:   "test-namespace",
		},
		Spec: statefulSet,
	}
}

func TestConvergeExistingStatefulSetDescriptor(t *testing.T) {
	statefulSetDescriptor := newTestStatefulSetDescriptor()
	client := clientstub.NewClientStub(t).(clientstub.ClientStub)
	expectSetClusterForExisting(statefulSetDescriptor.Metadata.ClusterName, "cluster-id", &client, t)
	expectSetProjectForExisting(statefulSetDescriptor.Metadata.ProjectName, "project-id", &client, t)
	client.DoHasStatefulSet = func(namespace string, part projectModel.StatefulSet) (bool, error) {
		assert.Equals(t, statefulSetDescriptor.Metadata.Namespace, namespace)
		return true, nil
	}
	converger := NewStatefulSetConverger(statefulSetDescriptor)
	converger.Converge(client)
}

func TestConvergeNewStatefulSetDescriptor(t *testing.T) {
	statefulSetDescriptor := newTestStatefulSetDescriptor()
	client := clientstub.NewClientStub(t).(clientstub.ClientStub)
	expectSetClusterForExisting(statefulSetDescriptor.Metadata.ClusterName, "cluster-id", &client, t)
	expectSetProjectForExisting(statefulSetDescriptor.Metadata.ProjectName, "project-id", &client, t)
	client.DoHasStatefulSet = func(namespace string, part projectModel.StatefulSet) (bool, error) {
		assert.Equals(t, statefulSetDescriptor.Metadata.Namespace, namespace)
		return false, nil
	}
	client.DoCreateStatefulSet = func(namespace string, part projectModel.StatefulSet) error {
		assert.Equals(t, statefulSetDescriptor.Metadata.Namespace, namespace)
		assert.Equals(t, statefulSetDescriptor.Spec, part)
		return nil
	}
	converger := NewStatefulSetConverger(statefulSetDescriptor)
	converger.Converge(client)
}
