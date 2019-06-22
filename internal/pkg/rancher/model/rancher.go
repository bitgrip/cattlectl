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

package model

// Types of Descriptors which are expected values of the field 'kind'
const (
	RancherKind     = "Rancher"
	ClusterKind     = "Cluster"
	ProjectKind     = "Project"
	JobKind         = "Job"
	CronJobKind     = "CronJob"
	DeploymentKind  = "Deployment"
	DaemonSetKind   = "DaemonSet"
	StatefulSetKind = "StatefulSet"
)

// Rancher represents global members
type Rancher struct {
	APIVersion string          `yaml:"api_version"`
	Kind       string          `yaml:"kind"`
	Metadata   RancherMetadata `yaml:"metadata"`
	Catalogs   []Catalog       `yaml:"catalogs,omitempty"`
}

// RancherMetadata are global meta informations
type RancherMetadata struct {
	RancherURL string `yaml:"rancher_url,omitempty"`
	AccessKey  string `yaml:"access_key,omitempty"`
	SecretKey  string `yaml:"secret_key,omitempty"`
	TokenKey   string `yaml:"token_key,omitempty"`
}

// Catalog is a template source to be used by Apps
type Catalog struct {
	Name     string `yaml:"name"`
	URL      string `yaml:"url"`
	Branch   string `yaml:"branch,omitempty"`
	Username string `yaml:"username,omitempty"`
	Password string `yaml:"password,omitempty"`
}
