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

package project

import (
	"github.com/bitgrip/cattlectl/internal/pkg/rancher/client"
	"github.com/bitgrip/cattlectl/internal/pkg/rancher/descriptor"
	projectModel "github.com/bitgrip/cattlectl/internal/pkg/rancher/project/model"
)

// NewProjectConverger creates a Converger for a given github.com/bitgrip/cattlectl/internal/pkg/projectModel.Project
func NewProjectConverger(project projectModel.Project, clusterClient client.ClusterClient) (descriptor.Converger, error) {
	projectClient, err := clusterClient.Project(project.Metadata.Name)
	if err != nil {
		return nil, err
	}
	childConvergers := make([]descriptor.Converger, 0)
	for _, namespace := range project.Namespaces {
		namespaceClient, err := projectClient.Namespace(namespace.Name)
		if err != nil {
			return nil, err
		}
		namespaceClient.SetData(namespace)
		childConvergers = append(childConvergers, &descriptor.ResourceClientConverger{
			Client: namespaceClient,
		})
	}
	for _, certificate := range project.Resources.Certificates {
		certificateClient, err := projectClient.NamespacedCertificate(certificate.Name, certificate.Namespace)
		if err != nil {
			return nil, err
		}
		certificateClient.SetData(certificate)
		childConvergers = append(childConvergers, &descriptor.ResourceClientConverger{
			Client: certificateClient,
		})
	}
	for _, configMap := range project.Resources.ConfigMaps {
		configMapClient, err := projectClient.ConfigMap(configMap.Name, configMap.Namespace)
		if err != nil {
			return nil, err
		}
		configMapClient.SetData(configMap)
		childConvergers = append(childConvergers, &descriptor.ResourceClientConverger{
			Client: configMapClient,
		})
	}
	for _, dockerCredential := range project.Resources.DockerCredentials {
		dockerCredentialClient, err := projectClient.NamespacedDockerCredential(dockerCredential.Name, dockerCredential.Namespace)
		if err != nil {
			return nil, err
		}
		dockerCredentialClient.SetData(dockerCredential)
		childConvergers = append(childConvergers, &descriptor.ResourceClientConverger{
			Client: dockerCredentialClient,
		})
	}
	for _, secret := range project.Resources.Secrets {
		secretClient, err := projectClient.NamespacedSecret(secret.Name, secret.Namespace)
		if err != nil {
			return nil, err
		}
		secretClient.SetData(secret)
		childConvergers = append(childConvergers, &descriptor.ResourceClientConverger{
			Client: secretClient,
		})
	}
	for _, storageClass := range project.StorageClasses {
		storageClassClient, err := clusterClient.StorageClass(storageClass.Name)
		if err != nil {
			return nil, err
		}
		storageClassClient.SetData(storageClass)
		childConvergers = append(childConvergers, &descriptor.ResourceClientConverger{
			Client: storageClassClient,
		})
	}
	for _, persistentVolume := range project.PersistentVolumes {
		persistentVolumeClient, err := clusterClient.PersistentVolume(persistentVolume.Name)
		if err != nil {
			return nil, err
		}
		persistentVolumeClient.SetData(persistentVolume)
		childConvergers = append(childConvergers, &descriptor.ResourceClientConverger{
			Client: persistentVolumeClient,
		})
	}
	for _, app := range project.Apps {
		appClient, err := projectClient.App(app.Name)
		if err != nil {
			return nil, err
		}
		appClient.SetData(app)
		childConvergers = append(childConvergers, &descriptor.ResourceClientConverger{
			Client: appClient,
		})
	}
	return &descriptor.ResourceClientConverger{
		Client:   projectClient,
		Children: childConvergers,
	}, nil
}
