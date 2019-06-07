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
	projectModel "github.com/bitgrip/cattlectl/internal/pkg/rancher/project/model"
	backendClusterClient "github.com/rancher/types/client/cluster/v3"
	backendRancherClient "github.com/rancher/types/client/management/v3"
	backendProjectClient "github.com/rancher/types/client/project/v3"
)

// RancherClient is comunicating with the Rancher server
type RancherClient interface {
	Cluster(clusterName string) (ClusterClient, error)
	Clusters() ([]ClusterClient, error)

	backendRancherClient() (*backendRancherClient.Client, error)
}

// ResourceClient is a client to any Rancher resource
type ResourceClient interface {
	ID() (string, error)
	Name() (string, error)
	Exists() (bool, error)
	Create() error
	Upgrade() error
}

// NamespacedResourceClient is a client to any Rancher resource belonging to a namespace
type NamespacedResourceClient interface {
	ResourceClient
	NamespaceID() (string, error)
	Namespace() (string, error)
}

// ClusterClient interacts with a Rancher cluster resource
type ClusterClient interface {
	ResourceClient
	Project(projectName string) (ProjectClient, error)
	Projects() ([]ProjectClient, error)
	StorageClass(name string) (StorageClassClient, error)
	StorageClasses() ([]StorageClassClient, error)
	PersistentVolume(name string) (PersistentVolumeClient, error)
	PersistentVolumes() ([]PersistentVolumeClient, error)
	Namespace(name, projectName string) (NamespaceClient, error)
	Namespaces(projectName string) ([]NamespaceClient, error)

	backendRancherClient() (*backendRancherClient.Client, error)
	backendClusterClient() (*backendClusterClient.Client, error)
}

// ProjectClient interacts with a Rancher project resource
type ProjectClient interface {
	ResourceClient
	ClusterID() (string, error)
	Namespace(name string) (NamespaceClient, error)
	Namespaces() ([]NamespaceClient, error)
	GlobalCertificate(name string) (CertificateClient, error)
	GlobalCertificates() ([]CertificateClient, error)
	Certificate(name, namespaceName string) (CertificateClient, error)
	Certificates(namespaceName string) ([]CertificateClient, error)
	ConfigMap(name, namespaceName string) (ConfigMapClient, error)
	ConfigMaps(namespaceName string) ([]ConfigMapClient, error)
	GlobalDockerCredential(name string) (DockerCredentialClient, error)
	GlobalDockerCredentials() ([]DockerCredentialClient, error)
	DockerCredential(name, namespaceName string) (DockerCredentialClient, error)
	DockerCredentials(namespaceName string) ([]DockerCredentialClient, error)
	GlobalSecret(name string) (ConfigMapClient, error)
	GlobalSecrets() ([]ConfigMapClient, error)
	Secret(name, namespaceName string) (ConfigMapClient, error)
	Secrets(namespaceName string) ([]ConfigMapClient, error)
	App(name string) (AppClient, error)
	Apps() ([]AppClient, error)
	Job(name, namespaceName string) (JobClient, error)
	Jobs(namespaceName string) ([]JobClient, error)
	CronJob(name, namespaceName string) (CronJobClient, error)
	CronJobs(namespaceName string) ([]CronJobClient, error)
	Deployment(name, namespaceName string) (DeploymentClient, error)
	Deployments(namespaceName string) ([]DeploymentClient, error)
	DaemonSet(name, namespaceName string) (DaemonSetClient, error)
	DaemonSets(namespaceName string) ([]DaemonSetClient, error)
	StatefulSet(name, namespaceName string) (StatefulSetClient, error)
	StatefulSets(namespaceName string) ([]StatefulSetClient, error)

	backendProjectClient() (*backendProjectClient.Client, error)
}

// NamespaceClient interacts with a Rancher namespace resource
type NamespaceClient interface {
	ResourceClient
	HasProject() (bool, error)
	Data() (projectModel.Namespace, error)
	SetData(storageClass projectModel.Namespace) error
}

// StorageClassClient interacts with a Rancher storage class resource
type StorageClassClient interface {
	ResourceClient
	Data() (projectModel.StorageClass, error)
	SetData(storageClass projectModel.StorageClass) error
}

// PersistentVolumeClient interacts with a Rancher persistent volume resource
type PersistentVolumeClient interface {
	ResourceClient
	Data() (projectModel.PersistentVolume, error)
	SetData(persistentVolume projectModel.PersistentVolume) error
}

// CertificateClient interacts with a Rancher certificate resource
type CertificateClient interface {
	NamespacedResourceClient
	Data() (projectModel.Certificate, error)
	SetData(certificate projectModel.Certificate) error
}

// ConfigMapClient interacts with a Rancher config map or secret resource
type ConfigMapClient interface {
	NamespacedResourceClient
	Data() (projectModel.ConfigMap, error)
	SetData(configMap projectModel.ConfigMap) error
}

// DockerCredentialClient interacts with a Rancher docker credential resource
type DockerCredentialClient interface {
	NamespacedResourceClient
	Data() (projectModel.DockerCredential, error)
	SetData(dockerCredential projectModel.DockerCredential) error
}

// AppClient interacts with a Rancher app resource
type AppClient interface {
	ResourceClient
	Data() (projectModel.App, error)
	SetData(app projectModel.App) error
}

// JobClient interacts with a Rancher job resource
type JobClient interface {
	NamespacedResourceClient
	Data() (projectModel.Job, error)
	SetData(job projectModel.Job) error
}

// CronJobClient interacts with a Rancher cron job resource
type CronJobClient interface {
	NamespacedResourceClient
	Data() (projectModel.CronJob, error)
	SetData(job projectModel.CronJob) error
}

// DeploymentClient interacts with a Rancher deployment resource
type DeploymentClient interface {
	NamespacedResourceClient
	Data() (projectModel.Deployment, error)
	SetData(deployment projectModel.Deployment) error
}

// DaemonSetClient interacts with a Rancher daemon set resource
type DaemonSetClient interface {
	NamespacedResourceClient
	Data() (projectModel.DaemonSet, error)
	SetData(daemonSet projectModel.DaemonSet) error
}

// StatefulSetClient interacts with a Rancher stateful set resource
type StatefulSetClient interface {
	NamespacedResourceClient
	Data() (projectModel.StatefulSet, error)
	SetData(statefulSet projectModel.StatefulSet) error
}
