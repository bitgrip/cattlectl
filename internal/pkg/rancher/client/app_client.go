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

	projectModel "github.com/bitgrip/cattlectl/internal/pkg/rancher/project/model"
	"github.com/rancher/norman/types"
	backendClient "github.com/rancher/types/client/project/v3"
	"github.com/sirupsen/logrus"
)

func newAppClientWithData(
	app projectModel.App,
	namespace string,
	project ProjectClient,
	backendClient *backendClient.Client,
	logger *logrus.Entry,
) (AppClient, error) {
	result, err := newAppClient(
		app.Name,
		namespace,
		project,
		backendClient,
		logger,
	)
	if err != nil {
		return nil, err
	}
	err = result.SetData(app)
	return result, err
}

func newAppClient(
	name, namespace string,
	project ProjectClient,
	backendClient *backendClient.Client,
	logger *logrus.Entry,
) (AppClient, error) {
	return &appClient{
		namespacedResourceClient: namespacedResourceClient{
			resourceClient: resourceClient{
				name:   name,
				logger: logger.WithField("app_name", name).WithField("namespace", namespace),
			},
			namespace: namespace,
			project:   project,
		},
		backendClient: backendClient,
	}, nil
}

type appClient struct {
	namespacedResourceClient
	app           projectModel.App
	backendClient *backendClient.Client
	backendData   *backendClient.App
}

func (client *appClient) init() error {
	namespaceID, err := client.NamespaceID()
	if namespaceID == "" && err == nil {
		return fmt.Errorf("Can not find namespace")
	}
	return err
}

func (client *appClient) Exists() (bool, error) {
	if err := client.init(); err != nil {
		return false, err
	}
	collection, err := client.backendClient.App.List(&types.ListOpts{
		Filters: map[string]interface{}{
			"name":        client.name,
			"namespaceId": client.namespaceID,
		},
	})
	if nil != err {
		client.logger.WithError(err).Error("Failed to read app list")
		return false, fmt.Errorf("Failed to read app list, %v", err)
	}
	for _, item := range collection.Data {
		if item.Name == client.name && item.NamespaceId == client.namespaceID {
			return true, nil
		}
	}
	client.logger.Debug("App not found")
	return false, nil
}

func (client *appClient) Create() error {
	if err := client.init(); err != nil {
		return err
	}

	client.logger.Info("Create new app")
	pattern := &backendClient.App{
		Name:            client.app.Name,
		ExternalID:      fmt.Sprintf("catalog://?catalog=%v&template=%v&version=%v", client.app.Catalog, client.app.Chart, client.app.Version),
		TargetNamespace: client.app.Namespace,
		Answers:         client.app.Answers,
	}
	_, err := client.backendClient.App.Create(pattern)
	return err
}

func (client *appClient) Upgrade() error {
	return fmt.Errorf("upgrade statefulset not supported")
}

func (client *appClient) Data() (projectModel.App, error) {
	return client.app, nil
}

func (client *appClient) SetData(app projectModel.App) error {
	client.name = app.Name
	client.app = app
	return nil
}