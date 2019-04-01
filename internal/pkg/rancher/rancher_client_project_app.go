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
	"fmt"
	"reflect"

	projectModel "github.com/bitgrip/cattlectl/internal/pkg/rancher/project/model"
	"github.com/rancher/norman/types"
	projectClient "github.com/rancher/types/client/project/v3"
)

func (client *rancherClient) HasApp(app projectModel.App) (bool, error) {
	collection, err := client.projectClient.App.List(&types.ListOpts{
		Filters: map[string]interface{}{
			"name": app.Name,
		},
	})
	if nil != err {
		client.logger.WithError(err).WithField("app_name", app.Name).Error("Failed to read app list")
		return false, fmt.Errorf("Failed to read app list, %v", err)
	}
	for _, item := range collection.Data {
		if item.Name == app.Name {
			client.appCache[app.Name] = item
			return true, nil
		}
	}
	client.logger.WithField("app_name", app.Name).Debug("App not found")
	return false, nil
}

func (client *rancherClient) UpgradeApp(app projectModel.App) error {
	var installedApp projectClient.App
	if item, exists := client.appCache[app.Name]; exists {
		client.logger.WithField("app_name", app.Name).Trace("Use Cache")
		installedApp = item
	} else {
		client.logger.WithField("app_name", app.Name).Trace("Load from rancher")
		collection, err := client.projectClient.App.List(&types.ListOpts{
			Filters: map[string]interface{}{
				"name": app.Name,
			},
		})
		if nil != err {
			client.logger.WithError(err).WithField("app_name", app.Name).Error("Failed to read app list")
			return fmt.Errorf("Failed to read app list, %v", err)
		}

		if len(collection.Data) == 0 {
			return fmt.Errorf("App %v not found", app.Name)
		}

		installedApp = collection.Data[0]
	}

	if reflect.DeepEqual(installedApp.Answers, app.Answers) {
		client.logger.WithField("app_name", app.Name).Debug("Skip upgrade app - no changes")
		return nil
	}
	if app.SkipUpgrade {
		client.logger.WithField("app_name", app.Name).Info("Suppress upgrade app - by config")
		return nil
	}

	client.logger.WithField("app_name", app.Name).Info("Upgrade app")
	au := &projectClient.AppUpgradeConfig{
		Answers:    app.Answers,
		ExternalID: fmt.Sprintf("catalog://?catalog=%v&template=%v&version=%v", app.Catalog, app.Chart, app.Version),
	}
	return client.projectClient.App.ActionUpgrade(&installedApp, au)
}

func (client *rancherClient) CreateApp(app projectModel.App) error {
	client.logger.WithField("app_name", app.Name).Info("Create new app")
	pattern := &projectClient.App{
		Name:            app.Name,
		ExternalID:      fmt.Sprintf("catalog://?catalog=%v&template=%v&version=%v", app.Catalog, app.Chart, app.Version),
		TargetNamespace: app.Namespace,
		Answers:         app.Answers,
	}
	_, err := client.projectClient.App.Create(pattern)
	return err
}
