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

func newTestDeploymentDescriptor() projectModel.DeploymentDescriptor {
	deployment := projectModel.Deployment{}
	deployment.Name = "test-deployment"
	return projectModel.DeploymentDescriptor{
		Metadata: projectModel.WorkloadMetadata{
			ClusterName: "test-cluster",
			ProjectName: "test-project",
			Namespace:   "test-namespace",
		},
		Spec: deployment,
	}
}

func TestConvergeExistingDeploymentDescriptor(t *testing.T) {
	deploymentDescriptor := newTestDeploymentDescriptor()
	client := clientstub.NewClientStub(t).(clientstub.ClientStub)
	expectSetClusterForExisting(deploymentDescriptor.Metadata.ClusterName, "cluster-id", &client, t)
	expectSetProjectForExisting(deploymentDescriptor.Metadata.ProjectName, "project-id", &client, t)
	client.DoHasDeployment = func(namespace string, part projectModel.Deployment) (bool, error) {
		assert.Equals(t, deploymentDescriptor.Metadata.Namespace, namespace)
		return true, nil
	}
	converger := NewDeploymentConverger(deploymentDescriptor)
	converger.Converge(client)
}

func TestConvergeNewDeploymentDescriptor(t *testing.T) {
	deploymentDescriptor := newTestDeploymentDescriptor()
	client := clientstub.NewClientStub(t).(clientstub.ClientStub)
	expectSetClusterForExisting(deploymentDescriptor.Metadata.ClusterName, "cluster-id", &client, t)
	expectSetProjectForExisting(deploymentDescriptor.Metadata.ProjectName, "project-id", &client, t)
	client.DoHasDeployment = func(namespace string, part projectModel.Deployment) (bool, error) {
		assert.Equals(t, deploymentDescriptor.Metadata.Namespace, namespace)
		return false, nil
	}
	client.DoCreateDeployment = func(namespace string, part projectModel.Deployment) error {
		assert.Equals(t, deploymentDescriptor.Metadata.Namespace, namespace)
		assert.Equals(t, deploymentDescriptor.Spec, part)
		return nil
	}
	converger := NewDeploymentConverger(deploymentDescriptor)
	converger.Converge(client)
}
