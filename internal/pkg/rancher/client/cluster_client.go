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
	backendClusterClient "github.com/rancher/types/client/cluster/v3"
)

func newClusterClient() (ClusterClient, error) {
	return &clusterClient{}, nil
}

type clusterClient struct {
	resourceClient
	backendClusterClient backendClusterClient.Client
	cluster              projectModel.Cluster
}

func (client *clusterClient) Project(projectName, projectID string) (ProjectClient, error) {
	return nil, fmt.Errorf("upgrade statefulset not supported")
}
func (client *clusterClient) Projects() ([]ProjectClient, error) {
	return nil, fmt.Errorf("upgrade statefulset not supported")
}
func (client *clusterClient) StorageClass(name string) (StorageClassClient, error) {
	return nil, fmt.Errorf("upgrade statefulset not supported")
}
func (client *clusterClient) StorageClasses() ([]StorageClassClient, error) {
	return nil, fmt.Errorf("upgrade statefulset not supported")
}
func (client *clusterClient) PersistentVolume(name string) (PersistentVolumeClient, error) {
	return nil, fmt.Errorf("upgrade statefulset not supported")
}
func (client *clusterClient) PersistentVolumes() ([]PersistentVolumeClient, error) {
	return nil, fmt.Errorf("upgrade statefulset not supported")
}
func (client *clusterClient) Namespace(name, projectName string) (NamespaceClient, error) {
	return nil, fmt.Errorf("upgrade statefulset not supported")
}
func (client *clusterClient) Namespaces(projectName string) ([]NamespaceClient, error) {
	return nil, fmt.Errorf("upgrade statefulset not supported")
}
