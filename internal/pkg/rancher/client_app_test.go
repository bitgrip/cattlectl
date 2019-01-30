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
	"github.com/bitgrip/cattlectl/internal/pkg/rancher/stubs"
	"github.com/rancher/norman/types"
	projectClient "github.com/rancher/types/client/project/v3"
	"github.com/sirupsen/logrus"
)

/*
HasApp(...)
  returns true if app exists
*/
func TestHasApp_AppExisting(t *testing.T) {
	const (
		projectId   = "test-project-id"
		projectName = "test-project-name"
		clusterId   = "test-cluster-id"
		appName     = "test-app"
	)
	var (
		actualListOpts *types.ListOpts
		clientConfig   = ClientConfig{}
		app            = App{
			Name: appName,
		}
		testClients = stubs.CreateTestClients(t)
	)

	appOperationsStub := stubs.CreateAppOperationsStub(t)
	appOperationsStub.DoList = func(opts *types.ListOpts) (*projectClient.AppCollection, error) {
		actualListOpts = opts
		return &projectClient.AppCollection{
			Data: []projectClient.App{
				projectClient.App{
					Name: appName,
				},
			},
		}, nil
	}
	testClients.ProjectClient.App = appOperationsStub

	client := rancherClient{
		clientConfig:  clientConfig,
		projectId:     projectId,
		projectClient: testClients.ProjectClient,
		appCache:      make(map[string]projectClient.App),
		logger:        logrus.WithField("test", true),
	}
	//Act
	result, err := client.HasApp(app)

	//Assert
	assert.Ok(t, err)
	assert.Equals(t, true, result)
	assert.Equals(t, &types.ListOpts{Filters: map[string]interface{}{"name": appName}}, actualListOpts)
}

/*
HasApp(...)
  returns true if local persistent volume exists
*/
func TestHasApp_AppNotExisting(t *testing.T) {
	const (
		projectId   = "test-project-id"
		projectName = "test-project-name"
		clusterId   = "test-cluster-id"
		appName     = "test-app"
	)
	var (
		actualListOpts *types.ListOpts
		clientConfig   = ClientConfig{}
		app            = App{
			Name: appName,
		}
		testClients = stubs.CreateTestClients(t)
	)

	appOperationsStub := stubs.CreateAppOperationsStub(t)
	appOperationsStub.DoList = func(opts *types.ListOpts) (*projectClient.AppCollection, error) {
		actualListOpts = opts
		return &projectClient.AppCollection{
			Data: []projectClient.App{},
		}, nil
	}
	testClients.ProjectClient.App = appOperationsStub

	client := rancherClient{
		clientConfig:  clientConfig,
		projectId:     projectId,
		projectClient: testClients.ProjectClient,
		appCache:      make(map[string]projectClient.App),
		logger:        logrus.WithField("test", true),
	}
	//Act
	result, err := client.HasApp(app)

	//Assert
	assert.Ok(t, err)
	assert.Equals(t, false, result)
	assert.Equals(t, &types.ListOpts{Filters: map[string]interface{}{"name": appName}}, actualListOpts)
}

/*
CreateApp(...)
  uses project client to create app
*/
func TestCreateApp(t *testing.T) {
	const (
		projectId   = "test-project-id"
		projectName = "test-project-name"
		clusterId   = "test-cluster-id"
		appName     = "test-app"
		catalog     = "test-catalog"
		chart       = "test-chart"
		namespace   = "test-namespace"
		version     = "1.2.3"
	)
	var (
		answers      = map[string]string{"first_key": "first_value"}
		actualOpts   *projectClient.App
		clientConfig = ClientConfig{}
		app          = App{
			Name:      appName,
			Catalog:   catalog,
			Chart:     chart,
			Version:   version,
			Namespace: namespace,
			Answers:   answers,
		}
		testClients = stubs.CreateTestClients(t)
	)

	appOperationsStub := stubs.CreateAppOperationsStub(t)
	appOperationsStub.DoCreate = func(opts *projectClient.App) (*projectClient.App, error) {
		actualOpts = opts
		return &projectClient.App{}, nil
	}
	testClients.ProjectClient.App = appOperationsStub

	client := rancherClient{
		clientConfig:  clientConfig,
		projectId:     projectId,
		projectClient: testClients.ProjectClient,
		appCache:      make(map[string]projectClient.App),
		logger:        logrus.WithField("test", true),
	}
	//Act
	err := client.CreateApp(app)

	//Assert
	assert.Ok(t, err)
	assert.Equals(t, &projectClient.App{
		Answers:         answers,
		ExternalID:      "catalog://?catalog=test-catalog&template=test-chart&version=1.2.3",
		Name:            appName,
		OwnerReferences: []projectClient.OwnerReference(nil),
		TargetNamespace: namespace,
	}, actualOpts)
}

/*
UpgradeApp(...)
  uses project client to upgrade app
*/
func TestUpgradeApp(t *testing.T) {
	const (
		projectId   = "test-project-id"
		projectName = "test-project-name"
		clusterId   = "test-cluster-id"
		appName     = "test-app"
		appId       = "abcd"
		catalog     = "test-catalog"
		chart       = "test-chart"
		namespace   = "test-namespace"
		version     = "1.2.3"
	)
	var (
		answers      = map[string]string{"first_key": "first_value"}
		actualOpts   *projectClient.App
		actualInput  *projectClient.AppUpgradeConfig
		clientConfig = ClientConfig{}
		app          = App{
			Name:      appName,
			Catalog:   catalog,
			Chart:     chart,
			Version:   version,
			Namespace: namespace,
			Answers:   answers,
		}
		testClients = stubs.CreateTestClients(t)
	)

	appOperationsStub := stubs.CreateAppOperationsStub(t)
	appOperationsStub.DoActionUpgrade = func(opts *projectClient.App, input *projectClient.AppUpgradeConfig) error {
		actualOpts = opts
		actualInput = input
		return nil
	}
	testClients.ProjectClient.App = appOperationsStub

	client := rancherClient{
		clientConfig:  clientConfig,
		projectId:     projectId,
		projectClient: testClients.ProjectClient,
		appCache: map[string]projectClient.App{
			appName: projectClient.App{Resource: types.Resource{ID: appId}, Name: appName},
		},
		logger: logrus.WithField("test", true),
	}
	//Act
	err := client.UpgradeApp(app)

	//Assert
	assert.Ok(t, err)
	assert.Equals(t, appId, actualOpts.ID)
	assert.Equals(t, appName, actualOpts.Name)
	assert.Equals(t, &projectClient.AppUpgradeConfig{Answers: answers, ExternalID: "catalog://?catalog=test-catalog&template=test-chart&version=1.2.3"}, actualInput)
}
