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

type Cluster struct {
	ApiVersion        string `yaml:"api_version"`
	Kind              string
	Metadata          ClusterMetadata
	StorageClasses    []StorageClass     `yaml:"storage_classes,omitempty"`
	PersistentVolumes []PersistentVolume `yaml:"storage_class,omitempty"`
	Projects          []Project
}

type ClusterMetadata struct {
	Name       string
	ID         string `yaml:"id,omitempty"`
	RancherUrl string `yaml:"rancher_url,omitempty"`
	AccessKey  string `yaml:"access_key,omitempty"`
	SecretKey  string `yaml:"secret_key,omitempty"`
	TokenKey   string `yaml:"token_key,omitempty"`
}

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

type ProjectMetadata struct {
	Name        string
	ID          string `yaml:"id,omitempty"`
	RancherURL  string `yaml:"rancher_url,omitempty"`
	AccessKey   string `yaml:"access_key,omitempty"`
	SecretKey   string `yaml:"secret_key,omitempty"`
	TokenKey    string `yaml:"token_key,omitempty"`
	ClusterName string `yaml:"cluster_name,omitempty"`
	ClusterID   string `yaml:"cluster_id,omitempty"`
}

type Namespace struct {
	Name string
}

type Resources struct {
	Certificates      []Certificate      `yaml:"certificates,omitempty"`
	ConfigMaps        []ConfigMap        `yaml:"config_maps,omitempty"`
	DockerCredentials []DockerCredential `yaml:"docker_credential,omitempty"`
	Secrets           []ConfigMap        `yaml:"secrets,omitempty"`
}

type Certificate struct {
	Name      string
	Key       string
	Certs     string
	Namespace string `yaml:"namespace,omitempty"`
}

type ConfigMap struct {
	Name      string
	Data      map[string]string
	Namespace string `yaml:"namespace,omitempty"`
}

type DockerCredential struct {
	Name       string               `yaml:"name,omitempty"`
	Namespace  string               `yaml:"namespace,omitempty"`
	Registries []RegistryCredential `yaml:"registries,omitempty"`
}

type RegistryCredential struct {
	Name     string `yaml:"name,omitempty"`
	Password string `yaml:"password,omitempty"`
	Username string `yaml:"username,omitempty"`
}

type StorageClass struct {
	Name                    string
	Provisioner             string
	ReclaimPolicy           string                  `yaml:"reclaim_policy"`
	VolumeBindMode          string                  `yaml:"volume_bind_mode"`
	Parameters              map[string]string       `yaml:"parameters,omitempty"`
	MountOptions            []string                `yaml:"mount_options,omitempty"`
	CreatePersistentVolumes bool                    `yaml:"create_persistent_volumes,omitempty"`
	PersistentVolumeGroups  []PersistentVolumeGroup `yaml:"persistent_volume_groups,omitempty"`
}

type PersistentVolumeGroup struct {
	Name         string
	Type         string
	Path         string
	CreateScript string   `yaml:"create_script"`
	AccessModes  []string `yaml:"access_modes"`
	Capacity     string
	Nodes        []string
}

type PersistentVolume struct {
	Name             string
	Path             string
	Node             string
	StorageClassName string   `yaml:"storage_class_name"`
	AccessModes      []string `yaml:"access_modes"`
	Capacity         string
	InitScript       string `yaml:"init_script"`
}

type App struct {
	Name        string
	Catalog     string
	Chart       string
	Version     string
	Namespace   string
	SkipUpgrade bool `yaml:"skip_upgrade,omitempty"`
	Answers     map[string]string
}
