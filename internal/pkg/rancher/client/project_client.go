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
)

type projectClient struct {
	resourceClient
}

func (client *projectClient) ID() (string, error) {
	return client.id, nil
}

func (client *projectClient) Namespace(name string) (NamespaceClient, error) {
	return nil, fmt.Errorf("upgrade statefulset not supported")
}
func (client *projectClient) Namespaces() ([]NamespaceClient, error) {
	return nil, fmt.Errorf("upgrade statefulset not supported")
}
func (client *projectClient) Certificate(name string) (CertificateClient, error) {
	return nil, fmt.Errorf("upgrade statefulset not supported")
}
func (client *projectClient) Certificates() ([]CertificateClient, error) {
	return nil, fmt.Errorf("upgrade statefulset not supported")
}
func (client *projectClient) NamespacedCertificate(name, namespaceName string) (CertificateClient, error) {
	return nil, fmt.Errorf("upgrade statefulset not supported")
}
func (client *projectClient) NamespacedCertificates(namespaceName string) ([]CertificateClient, error) {
	return nil, fmt.Errorf("upgrade statefulset not supported")
}
func (client *projectClient) ConfigMap(name, namespaceName string) (ConfigMapClient, error) {
	return nil, fmt.Errorf("upgrade statefulset not supported")
}
func (client *projectClient) ConfigMaps(namespaceName string) ([]ConfigMapClient, error) {
	return nil, fmt.Errorf("upgrade statefulset not supported")
}
func (client *projectClient) DockerCredential(name string) (DockerCredentialClient, error) {
	return nil, fmt.Errorf("upgrade statefulset not supported")
}
func (client *projectClient) DockerCredentials() ([]DockerCredentialClient, error) {
	return nil, fmt.Errorf("upgrade statefulset not supported")
}
func (client *projectClient) NamespacedDockerCredential(name, namespaceName string) (DockerCredentialClient, error) {
	return nil, fmt.Errorf("upgrade statefulset not supported")
}
func (client *projectClient) NamespacedDockerCredentials(namespaceName string) ([]DockerCredentialClient, error) {
	return nil, fmt.Errorf("upgrade statefulset not supported")
}
func (client *projectClient) Secret(name string) (ConfigMapClient, error) {
	return nil, fmt.Errorf("upgrade statefulset not supported")
}
func (client *projectClient) Secrets() ([]ConfigMapClient, error) {
	return nil, fmt.Errorf("upgrade statefulset not supported")
}
func (client *projectClient) NamespacedSecret(name, namespaceName string) (ConfigMapClient, error) {
	return nil, fmt.Errorf("upgrade statefulset not supported")
}
func (client *projectClient) NamespacedSecrets(namespaceName string) ([]ConfigMapClient, error) {
	return nil, fmt.Errorf("upgrade statefulset not supported")
}
func (client *projectClient) App(name string) (AppClient, error) {
	return nil, fmt.Errorf("upgrade statefulset not supported")
}
func (client *projectClient) Apps() ([]AppClient, error) {
	return nil, fmt.Errorf("upgrade statefulset not supported")
}
func (client *projectClient) Job(name, namespaceName string) (JobClient, error) {
	return nil, fmt.Errorf("upgrade statefulset not supported")
}
func (client *projectClient) Jobs(namespaceName string) ([]JobClient, error) {
	return nil, fmt.Errorf("upgrade statefulset not supported")
}
func (client *projectClient) CronJob(name, namespaceName string) (CronJobClient, error) {
	return nil, fmt.Errorf("upgrade statefulset not supported")
}
func (client *projectClient) CronJobs(namespaceName string) ([]CronJobClient, error) {
	return nil, fmt.Errorf("upgrade statefulset not supported")
}
func (client *projectClient) Deployment(name, namespaceName string) (DeploymentClient, error) {
	return nil, fmt.Errorf("upgrade statefulset not supported")
}
func (client *projectClient) Deployments(namespaceName string) ([]DeploymentClient, error) {
	return nil, fmt.Errorf("upgrade statefulset not supported")
}
func (client *projectClient) DaemonSet(name, namespaceName string) (DaemonSetClient, error) {
	return nil, fmt.Errorf("upgrade statefulset not supported")
}
func (client *projectClient) DaemonSets(namespaceName string) ([]DaemonSetClient, error) {
	return nil, fmt.Errorf("upgrade statefulset not supported")
}
func (client *projectClient) StatefulSet(name, namespaceName string) (StatefulSetClient, error) {
	return nil, fmt.Errorf("upgrade statefulset not supported")
}
func (client *projectClient) StatefulSets(namespaceName string) ([]StatefulSetClient, error) {
	return nil, fmt.Errorf("upgrade statefulset not supported")
}
