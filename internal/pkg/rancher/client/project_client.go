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
	backendClusterClient "github.com/rancher/types/client/cluster/v3"
	backendRancherClient "github.com/rancher/types/client/management/v3"
	backendProjectClient "github.com/rancher/types/client/project/v3"
	"github.com/sirupsen/logrus"
)

func newProjectClient(
	name string,
	config RancherConfig,
	backendRancherClient *backendRancherClient.Client,
	backendClusterClient *backendClusterClient.Client,
	logger *logrus.Entry,
) (ProjectClient, error) {
	return &projectClient{
		resourceClient: resourceClient{
			name:   name,
			logger: logger.WithField("project_name", name),
		},
		config:               config,
		backendRancherClient: backendRancherClient,
		backendClusterClient: backendClusterClient,
	}, nil
}

type projectClient struct {
	resourceClient
	config               RancherConfig
	clusterClient        ClusterClient
	backendRancherClient *backendRancherClient.Client
	backendClusterClient *backendClusterClient.Client
	backendProjectClient *backendProjectClient.Client
	project              projectModel.Project
}

func (client *projectClient) init() error {
	if client.backendProjectClient != nil {
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
	client.backendProjectClient, err = createProjectClient(client.config, clusterID, projectID)
	return err
}

func (client *projectClient) ID() (string, error) {
	if client.id != "" {
		return client.id, nil
	}
	clusterID, err := client.clusterClient.ID()
	if err != nil {
		return "", err
	}
	collection, err := client.backendRancherClient.Project.List(&types.ListOpts{
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
	if err := client.init(); err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("upgrade statefulset not supported")
}
func (client *projectClient) Namespaces() ([]NamespaceClient, error) {
	if err := client.init(); err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("upgrade statefulset not supported")
}
func (client *projectClient) Certificate(name string) (CertificateClient, error) {
	if err := client.init(); err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("upgrade statefulset not supported")
}
func (client *projectClient) Certificates() ([]CertificateClient, error) {
	if err := client.init(); err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("upgrade statefulset not supported")
}
func (client *projectClient) NamespacedCertificate(name, namespaceName string) (CertificateClient, error) {
	if err := client.init(); err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("upgrade statefulset not supported")
}
func (client *projectClient) NamespacedCertificates(namespaceName string) ([]CertificateClient, error) {
	if err := client.init(); err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("upgrade statefulset not supported")
}
func (client *projectClient) ConfigMap(name, namespaceName string) (ConfigMapClient, error) {
	if err := client.init(); err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("upgrade statefulset not supported")
}
func (client *projectClient) ConfigMaps(namespaceName string) ([]ConfigMapClient, error) {
	if err := client.init(); err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("upgrade statefulset not supported")
}
func (client *projectClient) DockerCredential(name string) (DockerCredentialClient, error) {
	if err := client.init(); err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("upgrade statefulset not supported")
}
func (client *projectClient) DockerCredentials() ([]DockerCredentialClient, error) {
	if err := client.init(); err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("upgrade statefulset not supported")
}
func (client *projectClient) NamespacedDockerCredential(name, namespaceName string) (DockerCredentialClient, error) {
	if err := client.init(); err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("upgrade statefulset not supported")
}
func (client *projectClient) NamespacedDockerCredentials(namespaceName string) ([]DockerCredentialClient, error) {
	if err := client.init(); err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("upgrade statefulset not supported")
}
func (client *projectClient) Secret(name string) (ConfigMapClient, error) {
	if err := client.init(); err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("upgrade statefulset not supported")
}
func (client *projectClient) Secrets() ([]ConfigMapClient, error) {
	if err := client.init(); err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("upgrade statefulset not supported")
}
func (client *projectClient) NamespacedSecret(name, namespaceName string) (ConfigMapClient, error) {
	if err := client.init(); err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("upgrade statefulset not supported")
}
func (client *projectClient) NamespacedSecrets(namespaceName string) ([]ConfigMapClient, error) {
	if err := client.init(); err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("upgrade statefulset not supported")
}
func (client *projectClient) App(name string) (AppClient, error) {
	if err := client.init(); err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("upgrade statefulset not supported")
}
func (client *projectClient) Apps() ([]AppClient, error) {
	if err := client.init(); err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("upgrade statefulset not supported")
}
func (client *projectClient) Job(name, namespaceName string) (JobClient, error) {
	if err := client.init(); err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("upgrade statefulset not supported")
}
func (client *projectClient) Jobs(namespaceName string) ([]JobClient, error) {
	if err := client.init(); err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("upgrade statefulset not supported")
}
func (client *projectClient) CronJob(name, namespaceName string) (CronJobClient, error) {
	if err := client.init(); err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("upgrade statefulset not supported")
}
func (client *projectClient) CronJobs(namespaceName string) ([]CronJobClient, error) {
	if err := client.init(); err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("upgrade statefulset not supported")
}
func (client *projectClient) Deployment(name, namespaceName string) (DeploymentClient, error) {
	if err := client.init(); err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("upgrade statefulset not supported")
}
func (client *projectClient) Deployments(namespaceName string) ([]DeploymentClient, error) {
	if err := client.init(); err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("upgrade statefulset not supported")
}
func (client *projectClient) DaemonSet(name, namespaceName string) (DaemonSetClient, error) {
	if err := client.init(); err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("upgrade statefulset not supported")
}
func (client *projectClient) DaemonSets(namespaceName string) ([]DaemonSetClient, error) {
	if err := client.init(); err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("upgrade statefulset not supported")
}
func (client *projectClient) StatefulSet(name, namespaceName string) (StatefulSetClient, error) {
	if err := client.init(); err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("upgrade statefulset not supported")
}
func (client *projectClient) StatefulSets(namespaceName string) ([]StatefulSetClient, error) {
	if err := client.init(); err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("upgrade statefulset not supported")
}
