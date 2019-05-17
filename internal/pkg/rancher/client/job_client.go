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

func newJobClientWithData(
	job projectModel.Job,
	namespace string,
	project ProjectClient,
	backendClient *projectClient.Client,
	logger *logrus.Entry,
) (JobClient, error) {
	result, err := newJobClient(
		job.Name,
		namespace,
		project,
		backendClient,
		logger,
	)
	if err != nil {
		return nil, err
	}
	err = result.SetData(job)
	return result, err
}

func newJobClient(
	name, namespace string,
	project ProjectClient,
	backendClient *projectClient.Client,
	logger *logrus.Entry,
) (JobClient, error) {
	return &jobClient{
		namespacedResourceClient: namespacedResourceClient{
			resourceClient: resourceClient{
				name:   name,
				logger: logger.WithField("job_name", name).WithField("namespace", namespace),
			},
			namespace: namespace,
			project:   project,
		},
		backendClient: backendClient,
	}, nil
}

type jobClient struct {
	namespacedResourceClient
	job           projectModel.Job
	backendClient *projectClient.Client
}

func (client *jobClient) init() error {
	namespaceID, err := client.NamespaceID()
	if namespaceID == "" && err == nil {
		return fmt.Errorf("Can not find namespace")
	}
	return err
}

func (client *jobClient) Exists() (bool, error) {
	if err := client.init(); err != nil {
		return false, err
	}
	collection, err := client.backendClient.Job.List(&types.ListOpts{
		Filters: map[string]interface{}{
			"name":        client.name,
			"namespaceId": client.namespaceID,
		},
	})
	if nil != err {
		client.logger.WithError(err).Error("Failed to read job list")
		return false, fmt.Errorf("Failed to read job list, %v", err)
	}
	for _, item := range collection.Data {
		if item.Name == client.name && item.NamespaceId == client.namespaceID {
			return true, nil
		}
	}
	client.logger.Debug("Job not found")
	return false, nil
}

func (client *jobClient) Create() error {
	if err := client.init(); err != nil {
		return err
	}
	client.logger.Info("Create new job")
	pattern, err := projectModel.ConvertJobToProjectAPI(client.job)
	if err != nil {
		return err
	}
	pattern.NamespaceId = client.namespaceID
	_, err = client.backendClient.Job.Create(&pattern)
	return err
}

func (client *jobClient) Upgrade() error {
	return fmt.Errorf("upgrade statefulset not supported")
}

func (client *jobClient) Data() (projectModel.Job, error) {
	return client.job, nil
}

func (client *jobClient) SetData(job projectModel.Job) error {
	client.name = job.Name
	client.job = job
	return nil
}
