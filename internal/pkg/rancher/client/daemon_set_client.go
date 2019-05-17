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
	projectClient "github.com/rancher/types/client/project/v3"
	"github.com/sirupsen/logrus"
)

func newDaemonSetClientWithData(
	daemonSet projectModel.DaemonSet,
	namespace string,
	project ProjectClient,
	backendClient *projectClient.Client,
	logger *logrus.Entry,
) (DaemonSetClient, error) {
	result, err := newDaemonSetClient(
		daemonSet.Name,
		namespace,
		project,
		backendClient,
		logger,
	)
	if err != nil {
		return nil, err
	}
	err = result.SetData(daemonSet)
	return result, err
}

func newDaemonSetClient(
	name, namespace string,
	project ProjectClient,
	backendClient *projectClient.Client,
	logger *logrus.Entry,
) (DaemonSetClient, error) {
	return &daemonSetClient{
		namespacedResourceClient: namespacedResourceClient{
			resourceClient: resourceClient{
				name:   name,
				logger: logger.WithField("daemonSet_name", name).WithField("namespace", namespace),
			},
			namespace: namespace,
			project:   project,
		},
		backendClient: backendClient,
	}, nil
}

type daemonSetClient struct {
	namespacedResourceClient
	daemonSet     projectModel.DaemonSet
	backendClient *projectClient.Client
}

func (client *daemonSetClient) init() error {
	namespaceID, err := client.NamespaceID()
	if namespaceID == "" && err == nil {
		return fmt.Errorf("Can not find namespace")
	}
	return err
}

func (client *daemonSetClient) Exists() (bool, error) {
	if err := client.init(); err != nil {
		return false, err
	}
	collection, err := client.backendClient.DaemonSet.List(&types.ListOpts{
		Filters: map[string]interface{}{
			"name":        client.name,
			"namespaceId": client.namespaceID,
		},
	})
	if nil != err {
		client.logger.WithError(err).Error("Failed to read daemonSet list")
		return false, fmt.Errorf("Failed to read daemonSet list, %v", err)
	}
	for _, item := range collection.Data {
		if item.Name == client.name && item.NamespaceId == client.namespaceID {
			return true, nil
		}
	}
	client.logger.Debug("DaemonSet not found")
	return false, nil
}

func (client *daemonSetClient) Create() error {
	if err := client.init(); err != nil {
		return err
	}
	client.logger.Info("Create new daemonSet")
	pattern, err := projectModel.ConvertDaemonSetToProjectAPI(client.daemonSet)
	if err != nil {
		return err
	}
	pattern.NamespaceId = client.namespaceID
	_, err = client.backendClient.DaemonSet.Create(&pattern)
	return err
}

func (client *daemonSetClient) Upgrade() error {
	return fmt.Errorf("upgrade statefulset not supported")
}

func (client *daemonSetClient) Data() (projectModel.DaemonSet, error) {
	return client.daemonSet, nil
}

func (client *daemonSetClient) SetData(daemonSet projectModel.DaemonSet) error {
	client.name = daemonSet.Name
	client.daemonSet = daemonSet
	return nil
}