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

func newCronJobClientWithData(
	cronJob projectModel.CronJob,
	namespace string,
	project ProjectClient,
	backendClient *projectClient.Client,
	logger *logrus.Entry,
) (CronJobClient, error) {
	result, err := newCronJobClient(
		cronJob.Name,
		namespace,
		project,
		backendClient,
		logger,
	)
	if err != nil {
		return nil, err
	}
	err = result.SetData(cronJob)
	return result, err
}

func newCronJobClient(
	name, namespace string,
	project ProjectClient,
	backendClient *projectClient.Client,
	logger *logrus.Entry,
) (CronJobClient, error) {
	return &cronJobClient{
		namespacedResourceClient: namespacedResourceClient{
			resourceClient: resourceClient{
				name:   name,
				logger: logger.WithField("cronJob_name", name).WithField("namespace", namespace),
			},
			namespace: namespace,
			project:   project,
		},
		backendClient: backendClient,
	}, nil
}

type cronJobClient struct {
	namespacedResourceClient
	cronJob       projectModel.CronJob
	backendClient *projectClient.Client
}

func (client *cronJobClient) init() error {
	namespaceID, err := client.NamespaceID()
	if namespaceID == "" && err == nil {
		return fmt.Errorf("Can not find namespace")
	}
	return err
}

func (client *cronJobClient) Exists() (bool, error) {
	if err := client.init(); err != nil {
		return false, err
	}
	collection, err := client.backendClient.CronJob.List(&types.ListOpts{
		Filters: map[string]interface{}{
			"name":        client.name,
			"namespaceId": client.namespaceID,
		},
	})
	if nil != err {
		client.logger.WithError(err).Error("Failed to read cronJob list")
		return false, fmt.Errorf("Failed to read cronJob list, %v", err)
	}
	for _, item := range collection.Data {
		if item.Name == client.name && item.NamespaceId == client.namespaceID {
			return true, nil
		}
	}
	client.logger.Debug("CronJob not found")
	return false, nil
}

func (client *cronJobClient) Create() error {
	if err := client.init(); err != nil {
		return err
	}
	client.logger.Info("Create new cronJob")
	pattern, err := projectModel.ConvertCronJobToProjectAPI(client.cronJob)
	if err != nil {
		return err
	}
	pattern.NamespaceId = client.namespaceID
	_, err = client.backendClient.CronJob.Create(&pattern)
	return err
}

func (client *cronJobClient) Upgrade() error {
	return fmt.Errorf("upgrade statefulset not supported")
}

func (client *cronJobClient) Data() (projectModel.CronJob, error) {
	return client.cronJob, nil
}

func (client *cronJobClient) SetData(cronJob projectModel.CronJob) error {
	client.name = cronJob.Name
	client.cronJob = cronJob
	return nil
}