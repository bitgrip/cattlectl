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
	backendProjectClient "github.com/rancher/types/client/project/v3"
	"github.com/sirupsen/logrus"
)

func newProjectClient(
	name string,
	config RancherConfig,
	clusterClient ClusterClient,
	logger *logrus.Entry,
) (ProjectClient, error) {
	projectLogger := logger.WithField("project_name", name)
	projectLogger.Trace("Create ProjectClient")
	return &projectClient{
		resourceClient: resourceClient{
			name:   name,
			logger: projectLogger,
		},
		config:                  config,
		certificateClients:      make(map[string]CertificateClient),
		configMapClients:        make(map[string]ConfigMapClient),
		dockerCredentialClients: make(map[string]DockerCredentialClient),
		secretClients:           make(map[string]ConfigMapClient),
		appClients:              make(map[string]AppClient),
		jobClients:              make(map[string]JobClient),
		cronJobClients:          make(map[string]CronJobClient),
		deploymentClients:       make(map[string]DeploymentClient),
		daemonSetClients:        make(map[string]DaemonSetClient),
		statefulSetClients:      make(map[string]StatefulSetClient),
	}, nil
}

type projectClient struct {
	resourceClient
	config                  RancherConfig
	clusterClient           ClusterClient
	_backendProjectClient   *backendProjectClient.Client
	project                 projectModel.Project
	certificateClients      map[string]CertificateClient
	configMapClients        map[string]ConfigMapClient
	dockerCredentialClients map[string]DockerCredentialClient
	secretClients           map[string]ConfigMapClient
	appClients              map[string]AppClient
	jobClients              map[string]JobClient
	cronJobClients          map[string]CronJobClient
	deploymentClients       map[string]DeploymentClient
	daemonSetClients        map[string]DaemonSetClient
	statefulSetClients      map[string]StatefulSetClient
}

func (client *projectClient) init() error {
	if client._backendProjectClient != nil {
		return nil
	}
	var (
		clusterID string
		projectID string
		err       error
	)
	if projectID, err = client.ID(); err != nil {
		return err
	}
	if clusterID, err = client.clusterClient.ID(); err != nil {
		return err
	}
	client._backendProjectClient, err = createProjectClient(client.config, clusterID, projectID)
	return err
}

func (client *projectClient) ID() (string, error) {
	if client.id != "" {
		return client.id, nil
	}
	backendRancherClient, err := client.clusterClient.backendRancherClient()
	if err != nil {
		return "", err
	}
	clusterID, err := client.clusterClient.ID()
	if err != nil {
		return "", err
	}
	collection, err := backendRancherClient.Project.List(&types.ListOpts{
		Filters: map[string]interface{}{
			"clusterId": clusterID,
			"name":      client.name,
		},
	})
	if err != nil {
		return "", err
	}
	if len(collection.Data) < 1 {
		return "", fmt.Errorf("Unknown Cluster [%s]", client.name)
	}
	client.id = collection.Data[0].ID
	return client.id, nil
}

func (client *projectClient) Exists() (bool, error) {
	_, err := client.ID()
	return err != nil, err
}
func (client *projectClient) Namespace(name string) (NamespaceClient, error) {
	return client.clusterClient.Namespace(name, client.name)
}
func (client *projectClient) Namespaces() ([]NamespaceClient, error) {
	return client.clusterClient.Namespaces(client.name)
}
func (client *projectClient) GlobalCertificate(name string) (CertificateClient, error) {
	return client.Certificate(name, "")
}
func (client *projectClient) GlobalCertificates() ([]CertificateClient, error) {
	backendProjectClient, err := client.backendProjectClient()
	if err != nil {
		return nil, err
	}

	collection, err := backendProjectClient.Certificate.List(&types.ListOpts{
		Filters: map[string]interface{}{
			"projectId": client.id,
		},
	})
	if err != nil {
		return nil, err
	}
	result := make([]CertificateClient, len(collection.Data))
	for i, backendCertificate := range collection.Data {
		certificate, err := client.GlobalCertificate(backendCertificate.Name)
		if err != nil {
			return nil, err
		}
		result[i] = certificate
	}
	return result, nil
}
func (client *projectClient) Certificate(name, namespaceName string) (CertificateClient, error) {
	if cache, exists := client.certificateClients[fmt.Sprintf("%s::%s", name, namespaceName)]; exists {
		return cache, nil
	}
	certificate, err := newCertificateClient(name, namespaceName, client, client.logger)
	if err != nil {
		return nil, err
	}
	client.certificateClients[fmt.Sprintf("%s::%s", name, namespaceName)] = certificate
	return certificate, nil
}
func (client *projectClient) Certificates(namespaceName string) ([]CertificateClient, error) {
	backendProjectClient, err := client.backendProjectClient()
	if err != nil {
		return nil, err
	}
	if namespaceName == "" {
		return client.GlobalCertificates()
	}

	namespace, err := client.Namespace(namespaceName)
	if err != nil {
		return nil, err
	}
	namespaceID, err := namespace.ID()
	if err != nil {
		return nil, err
	}

	collection, err := backendProjectClient.NamespacedCertificate.List(&types.ListOpts{
		Filters: map[string]interface{}{
			"projectId":   client.id,
			"namespaceId": namespaceID,
		},
	})
	if err != nil {
		return nil, err
	}
	result := make([]CertificateClient, len(collection.Data))
	for i, backendCertificate := range collection.Data {
		certificate, err := client.Certificate(backendCertificate.Name, namespaceName)
		if err != nil {
			return nil, err
		}
		result[i] = certificate
	}
	return result, nil
}
func (client *projectClient) ConfigMap(name, namespaceName string) (ConfigMapClient, error) {
	if cache, exists := client.configMapClients[fmt.Sprintf("%s::%s", name, namespaceName)]; exists {
		return cache, nil
	}
	configMap, err := newConfigMapClient(name, namespaceName, client, client.logger)
	if err != nil {
		return nil, err
	}
	client.configMapClients[fmt.Sprintf("%s::%s", name, namespaceName)] = configMap
	return configMap, nil
}
func (client *projectClient) ConfigMaps(namespaceName string) ([]ConfigMapClient, error) {
	backendProjectClient, err := client.backendProjectClient()
	if err != nil {
		return nil, err
	}

	namespace, err := client.Namespace(namespaceName)
	if err != nil {
		return nil, err
	}
	namespaceID, err := namespace.ID()
	if err != nil {
		return nil, err
	}

	collection, err := backendProjectClient.ConfigMap.List(&types.ListOpts{
		Filters: map[string]interface{}{
			"projectId":   client.id,
			"namespaceId": namespaceID,
		},
	})
	if err != nil {
		return nil, err
	}
	result := make([]ConfigMapClient, len(collection.Data))
	for i, backendConfigMap := range collection.Data {
		configMap, err := client.ConfigMap(backendConfigMap.Name, namespaceName)
		if err != nil {
			return nil, err
		}
		result[i] = configMap
	}
	return result, nil
}
func (client *projectClient) GlobalDockerCredential(name string) (DockerCredentialClient, error) {
	return client.DockerCredential(name, "")
}
func (client *projectClient) GlobalDockerCredentials() ([]DockerCredentialClient, error) {
	backendProjectClient, err := client.backendProjectClient()
	if err != nil {
		return nil, err
	}

	collection, err := backendProjectClient.DockerCredential.List(&types.ListOpts{
		Filters: map[string]interface{}{
			"projectId": client.id,
		},
	})
	if err != nil {
		return nil, err
	}
	result := make([]DockerCredentialClient, len(collection.Data))
	for i, backendDockerCredential := range collection.Data {
		dockerCredential, err := client.GlobalDockerCredential(backendDockerCredential.Name)
		if err != nil {
			return nil, err
		}
		result[i] = dockerCredential
	}
	return result, nil
}
func (client *projectClient) DockerCredential(name, namespaceName string) (DockerCredentialClient, error) {
	if cache, exists := client.dockerCredentialClients[fmt.Sprintf("%s::%s", name, namespaceName)]; exists {
		return cache, nil
	}
	dockerCredential, err := newDockerCredentialClient(name, namespaceName, client, client.logger)
	if err != nil {
		return nil, err
	}
	client.dockerCredentialClients[fmt.Sprintf("%s::%s", name, namespaceName)] = dockerCredential
	return dockerCredential, nil
}
func (client *projectClient) DockerCredentials(namespaceName string) ([]DockerCredentialClient, error) {
	backendProjectClient, err := client.backendProjectClient()
	if err != nil {
		return nil, err
	}
	if namespaceName == "" {
		return client.GlobalDockerCredentials()
	}

	namespace, err := client.Namespace(namespaceName)
	if err != nil {
		return nil, err
	}
	namespaceID, err := namespace.ID()
	if err != nil {
		return nil, err
	}

	collection, err := backendProjectClient.NamespacedDockerCredential.List(&types.ListOpts{
		Filters: map[string]interface{}{
			"projectId":   client.id,
			"namespaceId": namespaceID,
		},
	})
	if err != nil {
		return nil, err
	}
	result := make([]DockerCredentialClient, len(collection.Data))
	for i, backendDockerCredential := range collection.Data {
		dockerCredential, err := client.DockerCredential(backendDockerCredential.Name, namespaceName)
		if err != nil {
			return nil, err
		}
		result[i] = dockerCredential
	}
	return result, nil
}
func (client *projectClient) GlobalSecret(name string) (ConfigMapClient, error) {
	return client.Secret(name, "")
}
func (client *projectClient) GlobalSecrets() ([]ConfigMapClient, error) {
	backendProjectClient, err := client.backendProjectClient()
	if err != nil {
		return nil, err
	}

	collection, err := backendProjectClient.Secret.List(&types.ListOpts{
		Filters: map[string]interface{}{
			"projectId": client.id,
		},
	})
	if err != nil {
		return nil, err
	}
	result := make([]ConfigMapClient, len(collection.Data))
	for i, backendSecret := range collection.Data {
		secret, err := client.GlobalSecret(backendSecret.Name)
		if err != nil {
			return nil, err
		}
		result[i] = secret
	}
	return result, nil
}
func (client *projectClient) Secret(name, namespaceName string) (ConfigMapClient, error) {
	if cache, exists := client.secretClients[fmt.Sprintf("%s::%s", name, namespaceName)]; exists {
		return cache, nil
	}
	secret, err := newSecretClient(name, namespaceName, client, client.logger)
	if err != nil {
		return nil, err
	}
	client.secretClients[fmt.Sprintf("%s::%s", name, namespaceName)] = secret
	return secret, nil
}
func (client *projectClient) Secrets(namespaceName string) ([]ConfigMapClient, error) {
	backendProjectClient, err := client.backendProjectClient()
	if err != nil {
		return nil, err
	}
	if namespaceName == "" {
		return client.GlobalSecrets()
	}

	namespace, err := client.Namespace(namespaceName)
	if err != nil {
		return nil, err
	}
	namespaceID, err := namespace.ID()
	if err != nil {
		return nil, err
	}

	collection, err := backendProjectClient.NamespacedSecret.List(&types.ListOpts{
		Filters: map[string]interface{}{
			"projectId":   client.id,
			"namespaceId": namespaceID,
		},
	})
	if err != nil {
		return nil, err
	}
	result := make([]ConfigMapClient, len(collection.Data))
	for i, backendSecret := range collection.Data {
		secret, err := client.Secret(backendSecret.Name, namespaceName)
		if err != nil {
			return nil, err
		}
		result[i] = secret
	}
	return result, nil
}
func (client *projectClient) App(name string) (AppClient, error) {
	if cache, exists := client.appClients[name]; exists {
		return cache, nil
	}
	app, err := newAppClient(name, client, client.logger)
	if err != nil {
		return nil, err
	}
	client.appClients[name] = app
	return app, nil
}
func (client *projectClient) Apps() ([]AppClient, error) {
	backendProjectClient, err := client.backendProjectClient()
	if err != nil {
		return nil, err
	}

	collection, err := backendProjectClient.App.List(&types.ListOpts{
		Filters: map[string]interface{}{
			"projectId": client.id,
		},
	})
	if err != nil {
		return nil, err
	}
	result := make([]AppClient, len(collection.Data))
	for i, backendApp := range collection.Data {
		app, err := client.App(backendApp.Name)
		if err != nil {
			return nil, err
		}
		result[i] = app
	}
	return result, nil
}
func (client *projectClient) Job(name, namespaceName string) (JobClient, error) {
	if cache, exists := client.jobClients[fmt.Sprintf("%s::%s", name, namespaceName)]; exists {
		return cache, nil
	}
	job, err := newJobClient(name, namespaceName, client, client.logger)
	if err != nil {
		return nil, err
	}
	client.jobClients[fmt.Sprintf("%s::%s", name, namespaceName)] = job
	return job, nil
}
func (client *projectClient) Jobs(namespaceName string) ([]JobClient, error) {
	backendProjectClient, err := client.backendProjectClient()
	if err != nil {
		return nil, err
	}

	namespace, err := client.Namespace(namespaceName)
	if err != nil {
		return nil, err
	}
	namespaceID, err := namespace.ID()
	if err != nil {
		return nil, err
	}

	collection, err := backendProjectClient.Job.List(&types.ListOpts{
		Filters: map[string]interface{}{
			"projectId":   client.id,
			"namespaceId": namespaceID,
		},
	})
	if err != nil {
		return nil, err
	}
	result := make([]JobClient, len(collection.Data))
	for i, backendJob := range collection.Data {
		job, err := client.Job(backendJob.Name, namespaceName)
		if err != nil {
			return nil, err
		}
		result[i] = job
	}
	return result, nil
}
func (client *projectClient) CronJob(name, namespaceName string) (CronJobClient, error) {
	if cache, exists := client.cronJobClients[fmt.Sprintf("%s::%s", name, namespaceName)]; exists {
		return cache, nil
	}
	cronJob, err := newCronJobClient(name, namespaceName, client, client.logger)
	if err != nil {
		return nil, err
	}
	client.cronJobClients[fmt.Sprintf("%s::%s", name, namespaceName)] = cronJob
	return cronJob, nil
}
func (client *projectClient) CronJobs(namespaceName string) ([]CronJobClient, error) {
	backendProjectClient, err := client.backendProjectClient()
	if err != nil {
		return nil, err
	}

	namespace, err := client.Namespace(namespaceName)
	if err != nil {
		return nil, err
	}
	namespaceID, err := namespace.ID()
	if err != nil {
		return nil, err
	}

	collection, err := backendProjectClient.CronJob.List(&types.ListOpts{
		Filters: map[string]interface{}{
			"projectId":   client.id,
			"namespaceId": namespaceID,
		},
	})
	if err != nil {
		return nil, err
	}
	result := make([]CronJobClient, len(collection.Data))
	for i, backendCronJob := range collection.Data {
		cronJob, err := client.CronJob(backendCronJob.Name, namespaceName)
		if err != nil {
			return nil, err
		}
		result[i] = cronJob
	}
	return result, nil
}
func (client *projectClient) Deployment(name, namespaceName string) (DeploymentClient, error) {
	if cache, exists := client.deploymentClients[fmt.Sprintf("%s::%s", name, namespaceName)]; exists {
		return cache, nil
	}
	deployment, err := newDeploymentClient(name, namespaceName, client, client.logger)
	if err != nil {
		return nil, err
	}
	client.deploymentClients[fmt.Sprintf("%s::%s", name, namespaceName)] = deployment
	return deployment, nil
}
func (client *projectClient) Deployments(namespaceName string) ([]DeploymentClient, error) {
	backendProjectClient, err := client.backendProjectClient()
	if err != nil {
		return nil, err
	}

	namespace, err := client.Namespace(namespaceName)
	if err != nil {
		return nil, err
	}
	namespaceID, err := namespace.ID()
	if err != nil {
		return nil, err
	}

	collection, err := backendProjectClient.Deployment.List(&types.ListOpts{
		Filters: map[string]interface{}{
			"projectId":   client.id,
			"namespaceId": namespaceID,
		},
	})
	if err != nil {
		return nil, err
	}
	result := make([]DeploymentClient, len(collection.Data))
	for i, backendDeployment := range collection.Data {
		deployment, err := client.Deployment(backendDeployment.Name, namespaceName)
		if err != nil {
			return nil, err
		}
		result[i] = deployment
	}
	return result, nil
}
func (client *projectClient) DaemonSet(name, namespaceName string) (DaemonSetClient, error) {
	if cache, exists := client.daemonSetClients[fmt.Sprintf("%s::%s", name, namespaceName)]; exists {
		return cache, nil
	}
	daemonSet, err := newDaemonSetClient(name, namespaceName, client, client.logger)
	if err != nil {
		return nil, err
	}
	client.daemonSetClients[fmt.Sprintf("%s::%s", name, namespaceName)] = daemonSet
	return daemonSet, nil
}
func (client *projectClient) DaemonSets(namespaceName string) ([]DaemonSetClient, error) {
	backendProjectClient, err := client.backendProjectClient()
	if err != nil {
		return nil, err
	}

	namespace, err := client.Namespace(namespaceName)
	if err != nil {
		return nil, err
	}
	namespaceID, err := namespace.ID()
	if err != nil {
		return nil, err
	}

	collection, err := backendProjectClient.DaemonSet.List(&types.ListOpts{
		Filters: map[string]interface{}{
			"projectId":   client.id,
			"namespaceId": namespaceID,
		},
	})
	if err != nil {
		return nil, err
	}
	result := make([]DaemonSetClient, len(collection.Data))
	for i, backendDaemonSet := range collection.Data {
		daemonSet, err := client.DaemonSet(backendDaemonSet.Name, namespaceName)
		if err != nil {
			return nil, err
		}
		result[i] = daemonSet
	}
	return result, nil
}
func (client *projectClient) StatefulSet(name, namespaceName string) (StatefulSetClient, error) {
	if cache, exists := client.statefulSetClients[fmt.Sprintf("%s::%s", name, namespaceName)]; exists {
		return cache, nil
	}
	statefulSet, err := newStatefulSetClient(name, namespaceName, client, client.logger)
	if err != nil {
		return nil, err
	}
	client.statefulSetClients[fmt.Sprintf("%s::%s", name, namespaceName)] = statefulSet
	return statefulSet, nil
}
func (client *projectClient) StatefulSets(namespaceName string) ([]StatefulSetClient, error) {
	backendProjectClient, err := client.backendProjectClient()
	if err != nil {
		return nil, err
	}

	namespace, err := client.Namespace(namespaceName)
	if err != nil {
		return nil, err
	}
	namespaceID, err := namespace.ID()
	if err != nil {
		return nil, err
	}

	collection, err := backendProjectClient.StatefulSet.List(&types.ListOpts{
		Filters: map[string]interface{}{
			"projectId":   client.id,
			"namespaceId": namespaceID,
		},
	})
	if err != nil {
		return nil, err
	}
	result := make([]StatefulSetClient, len(collection.Data))
	for i, backendStatefulSet := range collection.Data {
		statefulSet, err := client.StatefulSet(backendStatefulSet.Name, namespaceName)
		if err != nil {
			return nil, err
		}
		result[i] = statefulSet
	}
	return result, nil
}

func (client *projectClient) backendProjectClient() (*backendProjectClient.Client, error) {
	if err := client.init(); err != nil {
		return nil, err
	}
	return client._backendProjectClient, nil
}
