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

package model

// Cluster is WIP
type Cluster struct {
	APIVersion        string `yaml:"api_version"`
	Kind              string
	Metadata          ClusterMetadata
	StorageClasses    []StorageClass     `yaml:"storage_classes,omitempty"`
	PersistentVolumes []PersistentVolume `yaml:"storage_class,omitempty"`
	Projects          []Project
}

// ClusterMetadata is WIP
type ClusterMetadata struct {
	Name       string
	ID         string `yaml:"id,omitempty"`
	RancherURL string `yaml:"rancher_url,omitempty"`
	AccessKey  string `yaml:"access_key,omitempty"`
	SecretKey  string `yaml:"secret_key,omitempty"`
	TokenKey   string `yaml:"token_key,omitempty"`
}

// Project is a subsection of a cluster
type Project struct {
	APIVersion        string `yaml:"api_version"`
	Kind              string
	Metadata          ProjectMetadata
	Namespaces        []Namespace        `yaml:"namespaces,omitempty"`
	Resources         Resources          `yaml:"resources,omitempty"`
	StorageClasses    []StorageClass     `yaml:"storage_classes,omitempty"`
	PersistentVolumes []PersistentVolume `yaml:"persistent_volumes,omitempty"`
	Apps              []App
}

// ProjectMetadata the meta informations about a Project
type ProjectMetadata struct {
	Name        string
	ID          string    `yaml:"id,omitempty"`
	RancherURL  string    `yaml:"rancher_url,omitempty"`
	AccessKey   string    `yaml:"access_key,omitempty"`
	SecretKey   string    `yaml:"secret_key,omitempty"`
	TokenKey    string    `yaml:"token_key,omitempty"`
	ClusterName string    `yaml:"cluster_name,omitempty"`
	ClusterID   string    `yaml:"cluster_id,omitempty"`
	Includes    []Include `yaml:"includes,omitempty"`
}

// Include is used to merge multiple descriptors into one
type Include struct {
	File      string `yaml:"file,omitempty"`
	Files     string `yaml:"files,omitempty"`
	Directory string `yaml:"directory,omitempty"`
}

// Namespace is a subsection of a Project and is represented in K8S as namespace
type Namespace struct {
	Name string
}

// Resources of a Project
type Resources struct {
	Certificates      []Certificate      `yaml:"certificates,omitempty"`
	ConfigMaps        []ConfigMap        `yaml:"config_maps,omitempty"`
	DockerCredentials []DockerCredential `yaml:"docker_credentials,omitempty"`
	Secrets           []ConfigMap        `yaml:"secrets,omitempty"`
}

// Certificate TLS certs used e.g. for https endpoints
type Certificate struct {
	Name      string
	Key       string
	Certs     string
	Namespace string `yaml:"namespace,omitempty"`
}

// ConfigMap data structure used for K8S configmaps and secrets
type ConfigMap struct {
	Name      string
	Data      map[string]string
	Namespace string `yaml:"namespace,omitempty"`
}

//DockerCredential to access docker registries
type DockerCredential struct {
	Name       string               `yaml:"name,omitempty"`
	Namespace  string               `yaml:"namespace,omitempty"`
	Registries []RegistryCredential `yaml:"registries,omitempty"`
}

// RegistryCredential credentials of one docker registry
type RegistryCredential struct {
	Name     string `yaml:"name,omitempty"`
	Password string `yaml:"password,omitempty"`
	Username string `yaml:"username,omitempty"`
}

// StorageClass represent a K8S StorageClass
type StorageClass struct {
	Name           string
	Provisioner    string
	ReclaimPolicy  string            `yaml:"reclaim_policy"`
	VolumeBindMode string            `yaml:"volume_bind_mode"`
	Parameters     map[string]string `yaml:"parameters,omitempty"`
	MountOptions   []string          `yaml:"mount_options,omitempty"`
}

// PersistentVolume represent a K8S PersistentVolume
type PersistentVolume struct {
	Name             string
	Type             string
	Path             string
	Node             string
	StorageClassName string   `yaml:"storage_class_name"`
	AccessModes      []string `yaml:"access_modes"`
	Capacity         string
	InitScript       string `yaml:"init_script"`
}

// App deployment using a Helm- or Rancher-Chart
type App struct {
	Name        string            `yaml:"name,omitempty"`
	Catalog     string            `yaml:"catalog,omitempty"`
	CatalogType string            `yaml:"catalog_type,omitempty"`
	Chart       string            `yaml:"chart,omitempty"`
	Version     string            `yaml:"version,omitempty"`
	Namespace   string            `yaml:"namespace,omitempty"`
	SkipUpgrade bool              `yaml:"skip_upgrade,omitempty"`
	Answers     map[string]string `yaml:"answers,omitempty"`
	ValuesYaml  string            `yaml:"values_yaml,omitempty"`
}
