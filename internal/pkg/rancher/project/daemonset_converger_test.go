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

func newTestDaemonSetDescriptor() projectModel.DaemonSetDescriptor {
	daemonSet := projectModel.DaemonSet{}
	daemonSet.Name = "test-daemonset"
	return projectModel.DaemonSetDescriptor{
		Metadata: projectModel.WorkloadMetadata{
			ClusterName: "test-cluster",
			ProjectName: "test-project",
			Namespace:   "test-namespace",
		},
		Spec: daemonSet,
	}
}

func TestConvergeExistingDaemonSetDescriptor(t *testing.T) {
	daemonSetDescriptor := newTestDaemonSetDescriptor()
	client := clientstub.NewClientStub(t).(clientstub.ClientStub)
	expectSetClusterForExisting(daemonSetDescriptor.Metadata.ClusterName, "cluster-id", &client, t)
	expectSetProjectForExisting(daemonSetDescriptor.Metadata.ProjectName, "project-id", &client, t)
	client.DoHasDaemonSet = func(namespace string, part projectModel.DaemonSet) (bool, error) {
		assert.Equals(t, daemonSetDescriptor.Metadata.Namespace, namespace)
		return true, nil
	}
	converger := NewDaemonSetConverger(daemonSetDescriptor)
	converger.Converge(client)
}

func TestConvergeNewDaemonSetDescriptor(t *testing.T) {
	daemonSetDescriptor := newTestDaemonSetDescriptor()
	client := clientstub.NewClientStub(t).(clientstub.ClientStub)
	expectSetClusterForExisting(daemonSetDescriptor.Metadata.ClusterName, "cluster-id", &client, t)
	expectSetProjectForExisting(daemonSetDescriptor.Metadata.ProjectName, "project-id", &client, t)
	client.DoHasDaemonSet = func(namespace string, part projectModel.DaemonSet) (bool, error) {
		assert.Equals(t, daemonSetDescriptor.Metadata.Namespace, namespace)
		return false, nil
	}
	client.DoCreateDaemonSet = func(namespace string, part projectModel.DaemonSet) error {
		assert.Equals(t, daemonSetDescriptor.Metadata.Namespace, namespace)
		assert.Equals(t, daemonSetDescriptor.Spec, part)
		return nil
	}
	converger := NewDaemonSetConverger(daemonSetDescriptor)
	converger.Converge(client)
}
