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
	"github.com/sirupsen/logrus"
)

/*
HasNamespace(...)
  returns true if namespace exists
*/
func TestHasNamespace_NamespaceExisting(t *testing.T) {
	const (
		projectId     = "test-project-id"
		projectName   = "test-project-name"
		clusterId     = "test-cluster-id"
		namespaceName = "test-namespace"
	)
	var (
		actualListOpts *types.ListOpts
		clientConfig   = ClientConfig{}
		namespace      = projectModel.Namespace{
			Name: namespaceName,
		}
		testClients = stubs.CreateTestClients(t)
	)

	namespaceOperationsStub := stubs.CreateNamespaceOperationsStub(t)
	namespaceOperationsStub.DoList = func(opts *types.ListOpts) (*clusterClient.NamespaceCollection, error) {
		actualListOpts = opts
		return &clusterClient.NamespaceCollection{
			Data: []clusterClient.Namespace{
				clusterClient.Namespace{
					Name: namespaceName,
				},
			},
		}, nil

	}
	testClients.ClusterClient.Namespace = namespaceOperationsStub

	client := rancherClient{
		clientConfig:   clientConfig,
		projectId:      projectId,
		clusterClient:  testClients.ClusterClient,
		namespaceCache: make(map[string]clusterClient.Namespace),
		logger:         logrus.WithField("test", true),
	}
	//Act
	result, err := client.HasNamespace(namespace)

	//Assert
	assert.Ok(t, err)
	assert.Equals(t, true, result)
	assert.Equals(t, &types.ListOpts{Filters: map[string]interface{}{"system": "false", "name": "test-namespace"}}, actualListOpts)
}

/*
HasNamespace(...)
  returns true if namespace exists
*/
func TestHasNamespace_NamespaceNotExisting(t *testing.T) {
	const (
		projectId     = "test-project-id"
		projectName   = "test-project-name"
		clusterId     = "test-cluster-id"
		namespaceName = "test-namespace"
	)
	var (
		actualListOpts *types.ListOpts
		clientConfig   = ClientConfig{}
		namespace      = projectModel.Namespace{
			Name: namespaceName,
		}
		testClients = stubs.CreateTestClients(t)
	)

	namespaceOperationsStub := stubs.CreateNamespaceOperationsStub(t)
	namespaceOperationsStub.DoList = func(opts *types.ListOpts) (*clusterClient.NamespaceCollection, error) {
		actualListOpts = opts
		return &clusterClient.NamespaceCollection{
			Data: []clusterClient.Namespace{},
		}, nil

	}
	testClients.ClusterClient.Namespace = namespaceOperationsStub

	client := rancherClient{
		clientConfig:  clientConfig,
		projectId:     projectId,
		clusterClient: testClients.ClusterClient,
		logger:        logrus.WithField("test", true),
	}
	//Act
	result, err := client.HasNamespace(namespace)

	//Assert
	assert.Ok(t, err)
	assert.Equals(t, false, result)
	assert.Equals(t, &types.ListOpts{Filters: map[string]interface{}{"system": "false", "name": "test-namespace"}}, actualListOpts)
}

/*
CreateNamespace(...)
  uses cluster client to create namespace
*/
func TestCreateNamespace(t *testing.T) {
	const (
		projectId     = "test-project-id"
		projectName   = "test-project-name"
		clusterId     = "test-cluster-id"
		namespaceName = "test-namespace"
	)
	var (
		actualOpts   *clusterClient.Namespace
		clientConfig = ClientConfig{}
		namespace    = projectModel.Namespace{
			Name: namespaceName,
		}
		testClients = stubs.CreateTestClients(t)
	)

	namespaceOperationsStub := stubs.CreateNamespaceOperationsStub(t)
	namespaceOperationsStub.DoCreate = func(opts *clusterClient.Namespace) (*clusterClient.Namespace, error) {
		actualOpts = opts
		return nil, nil
	}
	testClients.ClusterClient.Namespace = namespaceOperationsStub

	client := rancherClient{
		clientConfig:  clientConfig,
		projectId:     projectId,
		clusterClient: testClients.ClusterClient,
		logger:        logrus.WithField("test", true),
	}
	//Act
	err := client.CreateNamespace(namespace)

	//Assert
	assert.Ok(t, err)
	assert.Equals(t, &clusterClient.Namespace{Name: "test-namespace", ProjectID: "test-project-id"}, actualOpts)
}
