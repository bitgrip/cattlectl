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

func newStatefulSetClient(
	name, namespace string,
	project ProjectClient,
	backendClient *projectClient.Client,
	logger *logrus.Entry,
) (StatefulSetClient, error) {
	return &statefulSetClient{
		namespacedResourceClient: namespacedResourceClient{
			resourceClient: resourceClient{
				name:   name,
				logger: logger.WithField("statefulSet_name", name).WithField("namespace", namespace),
			},
			namespace: namespace,
			project:   project,
		},
		backendClient: backendClient,
	}, nil
}

type statefulSetClient struct {
	namespacedResourceClient
	statefulSet   projectModel.StatefulSet
	backendClient *projectClient.Client
}

func (client *statefulSetClient) init() error {
	namespaceID, err := client.NamespaceID()
	if namespaceID == "" && err == nil {
		return fmt.Errorf("Can not find namespace")
	}
	return err
}

func (client *statefulSetClient) Exists() (bool, error) {
	if err := client.init(); err != nil {
		return false, err
	}
	collection, err := client.backendClient.StatefulSet.List(&types.ListOpts{
		Filters: map[string]interface{}{
			"name":        client.name,
			"namespaceId": client.namespaceID,
		},
	})
	if nil != err {
		client.logger.WithError(err).Error("Failed to read statefulSet list")
		return false, fmt.Errorf("Failed to read statefulSet list, %v", err)
	}
	for _, item := range collection.Data {
		if item.Name == client.name && item.NamespaceId == client.namespaceID {
			return true, nil
		}
	}
	client.logger.Debug("StatefulSet not found")
	return false, nil
}

func (client *statefulSetClient) Create() error {
	if err := client.init(); err != nil {
		return err
	}
	client.logger.Info("Create new statefulSet")
	pattern, err := projectModel.ConvertStatefulSetToProjectAPI(client.statefulSet)
	if err != nil {
		return err
	}
	pattern.NamespaceId = client.namespaceID
	_, err = client.backendClient.StatefulSet.Create(&pattern)
	return err
}

func (client *statefulSetClient) Upgrade() error {
	return fmt.Errorf("upgrade statefulset not supported")
}

func (client *statefulSetClient) Data() (projectModel.StatefulSet, error) {
	return client.statefulSet, nil
}

func (client *statefulSetClient) SetData(statefulSet projectModel.StatefulSet) error {
	client.statefulSet = statefulSet
	return nil
}
