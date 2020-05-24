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

	projectModel "github.com/bitgrip/cattlectl/internal/pkg/rancher/cluster/project/model"
	rancherModel "github.com/bitgrip/cattlectl/internal/pkg/rancher/model"
	"github.com/rancher/norman/types"
	"github.com/sirupsen/logrus"
)

func newStatefulSetClientWithData(
	statefulSet projectModel.StatefulSet,
	namespace string,
	project ProjectClient,
	logger *logrus.Entry,
) (StatefulSetClient, error) {
	result, err := newStatefulSetClient(
		statefulSet.Name,
		namespace,
		project,
		logger,
	)
	if err != nil {
		return nil, err
	}
	err = result.SetData(statefulSet)
	return result, err
}

func newStatefulSetClient(
	name, namespace string,
	project ProjectClient,
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
	}, nil
}

type statefulSetClient struct {
	namespacedResourceClient
	statefulSet projectModel.StatefulSet
}

func (client *statefulSetClient) Type() string {
	return rancherModel.StatefulSetKind
}

func (client *statefulSetClient) Exists() (bool, error) {
	backendClient, err := client.project.backendProjectClient()
	if err != nil {
		return false, err
	}
	namespaceID, err := client.NamespaceID()
	if err != nil {
		return false, err
	}
	collection, err := backendClient.StatefulSet.List(&types.ListOpts{
		Filters: map[string]interface{}{
			"name":        client.name,
			"namespaceId": namespaceID,
		},
	})
	if nil != err {
		client.logger.WithError(err).Error("Failed to read statefulSet list")
		return false, fmt.Errorf("Failed to read statefulSet list, %v", err)
	}
	for _, item := range collection.Data {
		if item.Name == client.name && item.NamespaceId == namespaceID {
			return true, nil
		}
	}
	client.logger.Debug("StatefulSet not found")
	return false, nil
}

func (client *statefulSetClient) Create(dryRun bool) (changed bool, err error) {
	backendClient, err := client.project.backendProjectClient()
	if err != nil {
		return
	}
	client.logger.Info("Create new statefulSet")
	pattern, err := projectModel.ConvertStatefulSetToProjectAPI(client.statefulSet)
	if err != nil {
		return
	}
	namespaceID, err := client.NamespaceID()
	if err != nil {
		return
	}
	pattern.NamespaceId = namespaceID

	if dryRun {
		client.logger.WithField("object", pattern).Info("Do Dry-Run Create")
	} else {
		_, err = backendClient.StatefulSet.Create(&pattern)
	}
	return err == nil, err
}

func (client *statefulSetClient) Upgrade(dryRun bool) (changed bool, err error) {
	client.logger.Warn("Skip change existing statefulset")
	return
}

func (client *statefulSetClient) Data() (projectModel.StatefulSet, error) {
	return client.statefulSet, nil
}

func (client *statefulSetClient) SetData(statefulSet projectModel.StatefulSet) error {
	client.name = statefulSet.Name
	client.statefulSet = statefulSet
	return nil
}
