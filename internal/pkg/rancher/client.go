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
	"strings"

	"github.com/rancher/norman/clientbase"
	"github.com/rancher/norman/types"
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
	HasNamespace(namespace Namespace) (bool, error)
	CreateNamespace(namespace Namespace) error
	HasCertificate(certificate Certificate) (bool, error)
	CreateCertificate(certificate Certificate) error
	HasConfigMap(configMap ConfigMap) (bool, error)
	CreateConfigMap(configMap ConfigMap) error
	HasDockerCredential(dockerCredential DockerCredential) (bool, error)
	CreateDockerCredential(dockerCredential DockerCredential) error
	HasSecret(secret ConfigMap) (bool, error)
	CreateSecret(secret ConfigMap) error
	HasNamespacedSecret(secret ConfigMap) (bool, error)
	CreateNamespacedSecret(secret ConfigMap) error
	HasStorageClass(storageClass StorageClass) (bool, error)
	CreateStorageClass(storageClass StorageClass) error
	HasPersistentVolume(persistentVolume PersistentVolume) (bool, error)
	CreatePersistentVolume(persistentVolume PersistentVolume) error
	HasApp(app App) (bool, error)
	UpgradeApp(app App) error
	CreateApp(app App) error
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
	logger           *logrus.Entry
}

func (client *rancherClient) HasClusterWithName(name string) (bool, string, error) {
	collection, err := client.managementClient.Cluster.List(&types.ListOpts{
		Filters: map[string]interface{}{
			"name": name,
		},
	})
	if err != nil {
		return false, "", err
	}
	if len(collection.Data) < 1 {
		return false, "", nil
	}
	return true, collection.Data[0].ID, nil
}

func (client *rancherClient) SetCluster(clusterName, clusterId string) error {
	if clusterClient, err := createClusterClient(
		client.clientConfig.RancherURL,
		client.clientConfig.AccessKey,
		client.clientConfig.SecretKey,
		clusterId,
	); err != nil {
		return err
	} else {
		client.logger = client.logger.WithFields(logrus.Fields{
			"cluster_name": clusterName,
			"project_name": "",
		})
		client.clusterId = clusterId
		client.clusterClient = clusterClient
		client.projectId = ""
		client.projectClient = nil
		return nil
	}
}

func (client *rancherClient) HasProjectWithName(name string) (bool, string, error) {
	collection, err := client.managementClient.Project.List(&types.ListOpts{
		Filters: map[string]interface{}{
			"clusterId": client.clusterId,
			"name":      name,
		},
	})
	if err != nil {
		return false, "", err
	}
	if len(collection.Data) < 1 {
		return false, "", nil
	}
	return true, collection.Data[0].ID, nil
}

func (client *rancherClient) SetProject(projectName, projectId string) error {
	if projectClient, err := createProjectClient(
		client.clientConfig.RancherURL,
		client.clientConfig.AccessKey,
		client.clientConfig.SecretKey,
		client.clusterId,
		projectId,
	); err != nil {
		return err
	} else {
		client.logger = client.logger.WithFields(logrus.Fields{
			"project_name": projectName,
		})
		client.projectId = projectId
		client.projectClient = projectClient
		return nil
	}
}

func (client *rancherClient) CreateProject(projectName string) (string, error) {
	client.logger.WithField("project_name", projectName).Info("Create new project")
	pattern := &managementClient.Project{
		ClusterID: client.clusterId,
		Name:      projectName,
	}
	createdProject, err := client.managementClient.Project.Create(pattern)
	if err != nil {
		client.logger.Warn("Failed to create project")
		return "", err
	}

	return createdProject.ID, nil
}

func (client *rancherClient) HasNamespace(namespace Namespace) (bool, error) {
	collection, err := client.clusterClient.Namespace.List(&types.ListOpts{
		Filters: map[string]interface{}{
			"system": "false",
			"name":   namespace.Name,
		},
	})
	if nil != err {
		client.logger.WithError(err).WithField("namespace_name", namespace.Name).Error("Failed to read namespace list")
		return false, fmt.Errorf("Failed to read namespace list, %v", err)
	}
	for _, item := range collection.Data {
		if item.Name == namespace.Name {
			return true, nil
		}
	}
	client.logger.WithError(err).WithField("namespace_name", namespace.Name).Debug("Namespace not found")
	return false, nil
}

func (client *rancherClient) CreateNamespace(namespace Namespace) error {
	client.logger.WithField("namespace_name", namespace.Name).Info("Create new namespace")
	newNamespace := &clusterClient.Namespace{
		Name:      namespace.Name,
		ProjectID: client.projectId,
	}

	_, err := client.clusterClient.Namespace.Create(newNamespace)
	return err
}

func (client *rancherClient) HasCertificate(certificate Certificate) (bool, error) {
	collection, err := client.projectClient.Certificate.List(&types.ListOpts{
		Filters: map[string]interface{}{
			"system": "false",
			"name":   certificate.Name,
		},
	})
	if nil != err {
		client.logger.WithError(err).WithField("certificate_name", certificate.Name).Error("Failed to read certificate list")
		return false, fmt.Errorf("Failed to read certificate list, %v", err)
	}
	for _, item := range collection.Data {
		if item.Name == certificate.Name {
			return true, nil
		}
	}
	client.logger.WithError(err).WithField("certificate_name", certificate.Name).Debug("Certificate not found")
	return false, nil
}

func (client *rancherClient) CreateCertificate(certificate Certificate) error {
	client.logger.WithField("certificate_name", certificate.Name).Info("Create new certificate")
	newCertificate := &projectClient.Certificate{
		Name:        certificate.Name,
		Key:         certificate.Key,
		Certs:       certificate.Certs,
		NamespaceId: certificate.Namespace,
		ProjectID:   client.projectId,
	}

	_, err := client.projectClient.Certificate.Create(newCertificate)
	return err
}

func (client *rancherClient) HasConfigMap(configMap ConfigMap) (bool, error) {
	collection, err := client.projectClient.ConfigMap.List(&types.ListOpts{
		Filters: map[string]interface{}{
			"system": "false",
			"name":   configMap.Name,
		},
	})
	if nil != err {
		client.logger.WithError(err).WithField("config_map_name", configMap.Name).Error("Failed to read config_map list")
		return false, fmt.Errorf("Failed to read config_map list, %v", err)
	}
	for _, item := range collection.Data {
		if item.Name == configMap.Name {
			return true, nil
		}
	}
	client.logger.WithError(err).WithField("config_map_name", configMap.Name).Debug("ConfigMap not found")
	return false, nil
}

func (client *rancherClient) CreateConfigMap(configMap ConfigMap) error {
	client.logger.WithField("config_map_name", configMap.Name).Info("Create new ConfigMap")
	newConfigMap := &projectClient.ConfigMap{
		Name:        configMap.Name,
		Data:        configMap.Data,
		NamespaceId: configMap.Namespace,
		ProjectID:   client.projectId,
	}

	_, err := client.projectClient.ConfigMap.Create(newConfigMap)
	return err
}

func (client *rancherClient) HasDockerCredential(dockerCredential DockerCredential) (bool, error) {
	collection, err := client.projectClient.DockerCredential.List(&types.ListOpts{
		Filters: map[string]interface{}{
			"system": "false",
			"name":   dockerCredential.Name,
		},
	})
	if nil != err {
		client.logger.WithError(err).WithField("docker_credential_name", dockerCredential.Name).Error("Failed to read DockerCredential list")
		return false, fmt.Errorf("Failed to read docker_credential list, %v", err)
	}
	for _, item := range collection.Data {
		if item.Name == dockerCredential.Name {
			return true, nil
		}
	}
	client.logger.WithError(err).WithField("docker_credential_name", dockerCredential.Name).Debug("DockerCredential not found")
	return false, nil
}

func (client *rancherClient) CreateDockerCredential(dockerCredential DockerCredential) error {
	client.logger.WithField("docker_credential_name", dockerCredential.Name).Info("Create new DockerCredential")

	registries := make(map[string]projectClient.RegistryCredential)
	for _, registry := range dockerCredential.Registries {
		registries[registry.Name] = projectClient.RegistryCredential{
			Username: registry.Username,
			Password: registry.Password,
		}
	}

	newDockerCredential := &projectClient.DockerCredential{
		Name:        dockerCredential.Name,
		Registries:  registries,
		NamespaceId: dockerCredential.Namespace,
		ProjectID:   client.projectId,
	}

	_, err := client.projectClient.DockerCredential.Create(newDockerCredential)
	return err
}

func (client *rancherClient) HasSecret(secret ConfigMap) (bool, error) {
	collection, err := client.projectClient.Secret.List(&types.ListOpts{
		Filters: map[string]interface{}{
			"system": "false",
			"name":   secret.Name,
		},
	})
	if nil != err {
		client.logger.WithError(err).WithField("secret_name", secret.Name).Error("Failed to read secret list")
		return false, fmt.Errorf("Failed to read secret list, %v", err)
	}
	for _, item := range collection.Data {
		if item.Name == secret.Name {
			client.logger.WithError(err).WithField("secret_name", secret.Name).WithField("namespace", secret.Namespace).Debug("Secret found")
			return true, nil
		}
	}
	client.logger.WithError(err).WithField("secret_name", secret.Name).WithField("namespace", secret.Namespace).Debug("Secret not found")
	return false, nil
}

func (client *rancherClient) CreateSecret(secret ConfigMap) error {
	client.logger.WithField("secret_name", secret.Name).Info("Create new Secret")
	newSecret := &projectClient.Secret{
		Name:      secret.Name,
		Data:      secret.Data,
		ProjectID: client.projectId,
	}

	_, err := client.projectClient.Secret.Create(newSecret)
	return err
}

func (client *rancherClient) HasNamespacedSecret(secret ConfigMap) (bool, error) {
	collection, err := client.projectClient.NamespacedSecret.List(&types.ListOpts{
		Filters: map[string]interface{}{
			"system":      "false",
			"name":        secret.Name,
			"namespaceId": secret.Namespace,
		},
	})
	if nil != err {
		client.logger.WithError(err).WithField("secret_name", secret.Name).Error("Failed to read secret list")
		return false, fmt.Errorf("Failed to read secret list, %v", err)
	}
	for _, item := range collection.Data {
		if item.Name == secret.Name {
			client.logger.WithError(err).WithField("secret_name", secret.Name).WithField("namespace", secret.Namespace).Debug("Secret found")
			return true, nil
		}
	}
	client.logger.WithError(err).WithField("secret_name", secret.Name).WithField("namespace", secret.Namespace).Debug("Secret not found")
	return false, nil
}

func (client *rancherClient) CreateNamespacedSecret(secret ConfigMap) error {
	client.logger.WithField("secret_name", secret.Name).Info("Create new Secret")
	newSecret := &projectClient.NamespacedSecret{
		Name:        secret.Name,
		Data:        secret.Data,
		NamespaceId: secret.Namespace,
		ProjectID:   client.projectId,
	}

	_, err := client.projectClient.NamespacedSecret.Create(newSecret)
	return err
}

func (client *rancherClient) HasStorageClass(storageClass StorageClass) (bool, error) {
	collection, err := client.clusterClient.StorageClass.List(&types.ListOpts{
		Filters: map[string]interface{}{
			"system": "false",
			"name":   storageClass.Name,
		},
	})
	if nil != err {
		client.logger.WithError(err).WithField("storage_class_name", storageClass.Name).Error("Failed to read storage class list")
		return false, fmt.Errorf("Failed to read storage class list, %v", err)
	}
	for _, item := range collection.Data {
		if item.Name == storageClass.Name {
			return true, nil
		}
	}
	client.logger.WithField("storage_class_name", storageClass.Name).Debug("Storage class not found")
	return false, nil
}

func (client *rancherClient) CreateStorageClass(storageClass StorageClass) error {
	client.logger.WithField("storage_class_name", storageClass.Name).Info("Create new storage class")
	newStorageClass := &clusterClient.StorageClass{
		Name:              storageClass.Name,
		VolumeBindingMode: storageClass.VolumeBindMode,
		ReclaimPolicy:     storageClass.ReclaimPolicy,
		Provisioner:       storageClass.Provisioner,
		Parameters:        storageClass.Parameters,
		MountOptions:      storageClass.MountOptions,
	}

	_, err := client.clusterClient.StorageClass.Create(newStorageClass)
	return err
}

func (client *rancherClient) HasPersistentVolume(persistentVolume PersistentVolume) (bool, error) {
	collection, err := client.clusterClient.PersistentVolume.List(&types.ListOpts{
		Filters: map[string]interface{}{
			"system": "false",
			"name":   persistentVolume.Name,
		},
	})
	if nil != err {
		client.logger.WithError(err).WithField("persistent_volume_name", persistentVolume.Name).Error("Failed to read persistent volume list")
		return false, fmt.Errorf("Failed to read persistent volume list, %v", err)
	}
	for _, item := range collection.Data {
		if item.Name == persistentVolume.Name {
			return true, nil
		}
	}
	client.logger.WithField("persistent_volume_name", persistentVolume.Name).Debug("Persistent volume not found")
	return false, nil
}

func (client *rancherClient) CreatePersistentVolume(persistentVolume PersistentVolume) error {
	client.logger.WithField("local_volume_name", persistentVolume.Name).Info("Create new persistent volume")
	newPersistentVolume := &clusterClient.PersistentVolume{
		Name:                          persistentVolume.Name,
		StorageClassID:                persistentVolume.StorageClassName,
		AccessModes:                   persistentVolume.AccessModes,
		Capacity:                      map[string]string{"storage": persistentVolume.Capacity},
		PersistentVolumeReclaimPolicy: "Delete",
		Local: &clusterClient.LocalVolumeSource{
			Path: persistentVolume.Path,
		},
		NodeAffinity: &clusterClient.VolumeNodeAffinity{
			Required: &clusterClient.NodeSelector{
				NodeSelectorTerms: []clusterClient.NodeSelectorTerm{
					clusterClient.NodeSelectorTerm{
						MatchExpressions: []clusterClient.NodeSelectorRequirement{
							clusterClient.NodeSelectorRequirement{
								Key:      "kubernetes.io/hostname",
								Operator: "In",
								Values:   []string{persistentVolume.Node},
							},
						},
					},
				},
			},
		},
	}

	_, err := client.clusterClient.PersistentVolume.Create(newPersistentVolume)
	return err
}

func (client *rancherClient) HasApp(app App) (bool, error) {
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

func (client *rancherClient) UpgradeApp(app App) error {
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

func (client *rancherClient) CreateApp(app App) error {
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
