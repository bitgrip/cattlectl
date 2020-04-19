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
	"strings"

	projectModel "github.com/bitgrip/cattlectl/internal/pkg/rancher/cluster/project/model"
	"github.com/rancher/norman/types"
	backendProjectClient "github.com/rancher/types/client/project/v3"
	"github.com/sirupsen/logrus"
)

const (
	projectCatalogType       = "projectCatalog"
	clusterCatalogType       = "clusterCatalog"
	globalCatalogType        = ""
	projectCatalogExternalID = "catalog://?catalog=%s/%s&type=projectCatalog&template=%s&version=%s"
	clusterCatalogExternalID = "catalog://?catalog=%s/%s&type=clusterCatalog&template=%s&version=%s"
	globalCatalogExternalID  = "catalog://?catalog=%s&template=%s&version=%s"
)

func newAppClientWithData(
	app projectModel.App,
	project ProjectClient,
	logger *logrus.Entry,
) (AppClient, error) {
	result, err := newAppClient(
		app.Name,
		project,
		logger,
	)
	if err != nil {
		return nil, err
	}
	err = result.SetData(app)
	return result, err
}

func newAppClient(
	name string,
	project ProjectClient,
	logger *logrus.Entry,
) (AppClient, error) {
	return &appClient{
		resourceClient: resourceClient{
			name:   name,
			logger: logger.WithField("app_name", name),
		},
		projectClient: project,
	}, nil
}

type appClient struct {
	resourceClient
	app           projectModel.App
	backendData   *backendProjectClient.App
	projectClient ProjectClient
}

func (client *appClient) Exists() (bool, error) {
	backendClient, err := client.projectClient.backendProjectClient()
	if err != nil {
		return false, err
	}
	collection, err := backendClient.App.List(&types.ListOpts{
		Filters: map[string]interface{}{
			"name": client.name,
		},
	})
	if nil != err {
		client.logger.WithError(err).Error("Failed to read app list")
		return false, fmt.Errorf("Failed to read app list, %v", err)
	}
	for _, item := range collection.Data {
		if item.Name == client.name {
			return true, nil
		}
	}
	client.logger.Debug("App not found")
	return false, nil
}

func (client *appClient) Create() error {
	backendClient, err := client.projectClient.backendProjectClient()
	if err != nil {
		return err
	}

	externalID, err := client.externalID(client.app.Catalog, client.app.CatalogType, client.app.Chart, client.app.Version)
	if err != nil {
		return err
	}
	client.logger.Info("Create new app")
	pattern := &backendProjectClient.App{
		Name:            client.app.Name,
		ExternalID:      externalID,
		TargetNamespace: client.app.Namespace,
	}
	if client.app.ValuesYaml != "" {
		pattern.ValuesYaml = client.app.ValuesYaml
	} else {
		pattern.Answers = client.app.Answers
	}
	_, err = backendClient.App.Create(pattern)
	return err
}

func (client *appClient) Upgrade() (err error) {

	backendClient, err := client.projectClient.backendProjectClient()
	if err != nil {
		return
	}
	installedApp, err := client.loadInstalledApp()

	if installedApp == nil {
		return fmt.Errorf("App %v not found", client.name)
	}

	externalID, err := client.externalID(client.app.Catalog, client.app.CatalogType, client.app.Chart, client.app.Version)
	if err != nil {
		return err
	}
	au := &backendProjectClient.AppUpgradeConfig{
		ExternalID: externalID,
	}
	if client.app.ValuesYaml != "" {
		if strings.TrimSpace(installedApp.ValuesYaml) == strings.TrimSpace(client.app.ValuesYaml) {
			client.logger.Debug("Skip upgrade app - no changes")
			return nil
		}
		au.ValuesYaml = client.app.ValuesYaml
	} else {
		resultAnswers := map[string]string{}
		if client.projectClient.config().MergeAnswers {
			for key, value := range installedApp.Answers {
				resultAnswers[key] = value
			}
		}
		for key, value := range client.app.Answers {
			resultAnswers[key] = value
		}
		if reflect.DeepEqual(installedApp.Answers, resultAnswers) {
			client.logger.Debug("Skip upgrade app - no changes")
			return nil
		}
		au.Answers = resultAnswers
	}
	if client.app.SkipUpgrade {
		client.logger.Info("Suppress upgrade app - by config")
		return nil
	}

	client.logger.Info("Upgrade app")
	return backendClient.App.ActionUpgrade(installedApp, au)
}

func (client *appClient) Delete() (err error) {
	backendClient, err := client.projectClient.backendProjectClient()
	if err != nil {
		return
	}
	installedApp, err := client.loadInstalledApp()
	return backendClient.App.Delete(installedApp)
}

func (client *appClient) Data() (projectModel.App, error) {
	return client.app, nil
}

func (client *appClient) SetData(app projectModel.App) error {
	if len(app.Answers) != 0 && app.ValuesYaml != "" {
		return fmt.Errorf("Answers AND ValuesYaml is not supported")
	}
	client.name = app.Name
	client.app = app
	return nil
}

func (client *appClient) loadInstalledApp() (installedApp *backendProjectClient.App, err error) {
	backendClient, err := client.projectClient.backendProjectClient()
	if err != nil {
		return
	}
	client.logger.Trace("Load from rancher")
	collection, err := backendClient.App.List(&types.ListOpts{
		Filters: map[string]interface{}{
			"name": client.name,
		},
	})
	if nil != err {
		client.logger.WithError(err).Error("Failed to read app list")
		err = fmt.Errorf("Failed to read app list, %v", err)
		return
	}

	if len(collection.Data) == 0 {
		return
	}

	installedApp = &collection.Data[0]
	return
}

func (client *appClient) externalID(catalog, catalogType, chart, version string) (string, error) {
	switch catalogType {
	case projectCatalogType:
		projectID, err := client.projectClient.ID()
		if err != nil {
			return "", err
		}
		return fmt.Sprintf(projectCatalogExternalID, projectID, catalog, chart, version), nil
	case clusterCatalogType:
		clusterID, err := client.projectClient.ClusterID()
		if err != nil {
			return "", err
		}
		return fmt.Sprintf(clusterCatalogExternalID, clusterID, catalog, chart, version), nil
	case globalCatalogType:
		return fmt.Sprintf(globalCatalogExternalID, catalog, chart, version), nil
	default:
		return "", fmt.Errorf("Unknown catalog type %s", catalogType)
	}
}
