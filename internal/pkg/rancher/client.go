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
	"strings"

	projectModel "github.com/bitgrip/cattlectl/internal/pkg/rancher/project/model"
	"github.com/rancher/norman/clientbase"
	clusterClient "github.com/rancher/types/client/cluster/v3"
	managementClient "github.com/rancher/types/client/management/v3"
	projectClient "github.com/rancher/types/client/project/v3"
	"github.com/sirupsen/logrus"
)

type Client interface {
	HasClusterWithName(name string) (bool, string, error)
	SetCluster(clusterName, clusterId string) error
	HasProjectWithName(name string) (bool, string, error)
	SetProject(projectName, projectId string) error
	CreateProject(projectName string) (string, error)
	HasNamespace(namespace projectModel.Namespace) (bool, error)
	CreateNamespace(namespace projectModel.Namespace) error
	HasCertificate(certificate projectModel.Certificate) (bool, error)
	CreateCertificate(certificate projectModel.Certificate) error
	HasConfigMap(configMap projectModel.ConfigMap) (bool, error)
	CreateConfigMap(configMap projectModel.ConfigMap) error
	HasDockerCredential(dockerCredential projectModel.DockerCredential) (bool, error)
	CreateDockerCredential(dockerCredential projectModel.DockerCredential) error
	HasSecret(secret projectModel.ConfigMap) (bool, error)
	CreateSecret(secret projectModel.ConfigMap) error
	HasNamespacedSecret(secret projectModel.ConfigMap) (bool, error)
	CreateNamespacedSecret(secret projectModel.ConfigMap) error
	HasStorageClass(storageClass projectModel.StorageClass) (bool, error)
	CreateStorageClass(storageClass projectModel.StorageClass) error
	HasPersistentVolume(persistentVolume projectModel.PersistentVolume) (bool, error)
	CreatePersistentVolume(persistentVolume projectModel.PersistentVolume) error
	HasApp(app projectModel.App) (bool, error)
	UpgradeApp(app projectModel.App) error
	CreateApp(app projectModel.App) error
	// Workload
	HasJob(namespace string, job projectModel.Job) (bool, error)
	CreateJob(namespace string, job projectModel.Job) error
	HasCronJob(namespace string, cronJob projectModel.CronJob) (bool, error)
	CreateCronJob(namespace string, cronJob projectModel.CronJob) error
	HasDeployment(namespace string, deployment projectModel.Deployment) (bool, error)
	CreateDeployment(namespace string, deployment projectModel.Deployment) error
	HasDaemonSet(namespace string, daemonSet projectModel.DaemonSet) (bool, error)
	CreateDaemonSet(namespace string, daemonSet projectModel.DaemonSet) error
}

type ClientConfig struct {
	RancherURL string
	AccessKey  string
	SecretKey  string
}

var (
	newClusterClient    = clusterClient.NewClient
	newManagementClient = managementClient.NewClient
	newProjectClient    = projectClient.NewClient
)

func NewClient(clientConfig ClientConfig) (Client, error) {
	managementClient, err := createManagementClient(clientConfig.RancherURL, clientConfig.AccessKey, clientConfig.SecretKey)
	if err != nil {
		return nil, err
	}
	return &rancherClient{
		clientConfig:     clientConfig,
		managementClient: managementClient,
		appCache:         make(map[string]projectClient.App),
		namespaceCache:   make(map[string]clusterClient.Namespace),
		logger:           logrus.WithFields(logrus.Fields{}),
	}, nil
}

type rancherClient struct {
	clusterId        string
	projectId        string
	clientConfig     ClientConfig
	clusterClient    *clusterClient.Client
	managementClient *managementClient.Client
	projectClient    *projectClient.Client
	appCache         map[string]projectClient.App
	namespaceCache   map[string]clusterClient.Namespace
	logger           *logrus.Entry
}

func createClientOpts(rancherUrl string, accessKey string, secretKey string) *clientbase.ClientOpts {
	serverURL := rancherUrl

	if !strings.HasSuffix(serverURL, "/v3") {
		serverURL = serverURL + "/v3"
	}

	options := &clientbase.ClientOpts{
		URL:       serverURL,
		AccessKey: accessKey,
		SecretKey: secretKey,
		//TODO: Add CACert support
		//CACerts:   config.CACerts,
	}
	return options
}

func createClusterClient(rancherUrl string, accessKey string, secretKey string, clusterId string) (*clusterClient.Client, error) {
	logrus.WithFields(logrus.Fields{
		"rancher.url":        rancherUrl,
		"rancher.cluster_id": clusterId,
	}).Debug("Create Cluster Client")
	options := createClientOpts(rancherUrl, accessKey, secretKey)
	options.URL = options.URL + "/cluster/" + clusterId
	clusterClient, err := newClusterClient(options)
	if err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{
			"rancher.url":        rancherUrl,
			"rancher.cluster_id": clusterId,
		}).Error("Failed to create cluster client")
		return nil, fmt.Errorf("Failed to create cluster client, %v", err)
	}
	return clusterClient, nil
}

func createManagementClient(rancherUrl string, accessKey string, secretKey string) (*managementClient.Client, error) {
	logrus.WithFields(logrus.Fields{
		"rancher.url": rancherUrl,
	}).Debug("Create Management Client")
	options := createClientOpts(rancherUrl, accessKey, secretKey)
	managementClientCache, err := newManagementClient(options)
	if err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{
			"rancher.url": rancherUrl,
		}).Error("Failed to create management client")
		return nil, fmt.Errorf("Failed to create management client, %v", err)
	}

	return managementClientCache, nil
}

func createProjectClient(rancherUrl string, accessKey string, secretKey string, clusterId string, projectId string) (*projectClient.Client, error) {
	logrus.WithFields(logrus.Fields{
		"rancher.url":        rancherUrl,
		"rancher.cluster_id": clusterId,
		"project_id":         projectId,
	}).Debug("Create Project Client")
	options := createClientOpts(rancherUrl, accessKey, secretKey)
	options.URL = options.URL + "/projects/" + projectId

	pc, err := newProjectClient(options)
	if err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{
			"rancher.url":        rancherUrl,
			"rancher.cluster_id": clusterId,
			"project_id":         projectId,
		}).Error("Failed to create project client")
		return nil, fmt.Errorf("Failed to create project client, %v", err)
	}
	return pc, nil
}
