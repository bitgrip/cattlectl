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
	"fmt"
	"strings"

	"github.com/bitgrip/cattlectl/internal/pkg/rancher"
	"github.com/bitgrip/cattlectl/internal/pkg/rancher/descriptor"
	projectModel "github.com/bitgrip/cattlectl/internal/pkg/rancher/project/model"
	"github.com/sirupsen/logrus"
)

// NewProjectConverger creates a Converger for a given github.com/bitgrip/cattlectl/internal/pkg/projectModel.Project
func NewProjectConverger(project projectModel.Project) descriptor.Converger {
	partConvergers := []descriptor.Converger{
		newProjectPartConverger(project),
	}
	for _, namespace := range project.Namespaces {
		partConvergers = append(partConvergers, newNamespacePartConverger(namespace))
	}

	for _, certificate := range project.Resources.Certificates {
		partConvergers = append(partConvergers, newCertificatePartConverger(certificate))
	}
	for _, configMap := range project.Resources.ConfigMaps {
		partConvergers = append(partConvergers, newConfigMapPartConverger(configMap))
	}
	for _, dockerCredential := range project.Resources.DockerCredentials {
		partConvergers = append(partConvergers, newDockerCredentialPartConverger(dockerCredential))
	}
	for _, secret := range project.Resources.Secrets {
		partConvergers = append(partConvergers, newSecretPartConverger(secret))
	}
	for _, storageClass := range project.StorageClasses {
		partConvergers = append(partConvergers, newStorageClassPartConverger(storageClass))
	}
	for _, persistentVolume := range project.PersistentVolumes {
		partConvergers = append(partConvergers, newPersistentVolumePartConverger(persistentVolume))
	}
	for _, app := range project.Apps {
		partConvergers = append(partConvergers, newAppPartConverger(app))
	}
	return descriptor.ClusterResourceDescriptorConverger(
		project.Metadata.ClusterName,
		partConvergers,
	)
}

func newProjectPartConverger(project projectModel.Project) descriptor.Converger {
	return descriptor.PartConverger{
		PartName: "Project",
		HasPart: func(client rancher.Client) (bool, error) {
			if hasProject, projectID, err := client.HasProjectWithName(project.Metadata.Name); hasProject {
				return true, client.SetProject(project.Metadata.Name, projectID)
			} else if err != nil {
				return false, fmt.Errorf("Failed to check for project, %v", err)
			}
			return false, nil
		},
		UpdatePart: func(client rancher.Client) error {
			logrus.WithField("project", project.Metadata.Name).Debug("Project exists")
			return nil
		},
		CreatePart: func(client rancher.Client) error {
			projectID, err := client.CreateProject(project.Metadata.Name)
			if err != nil {
				return err
			}
			return client.SetProject(project.Metadata.Name, projectID)
		},
	}
}

func newNamespacePartConverger(namespace projectModel.Namespace) descriptor.Converger {
	return descriptor.PartConverger{
		PartName: "Namespace",
		HasPart: func(client rancher.Client) (bool, error) {
			return client.HasNamespace(namespace)
		},
		UpdatePart: func(client rancher.Client) error {
			logrus.WithField("namespace", namespace.Name).Debug("Skip change existing namespace")
			return nil
		},
		CreatePart: func(client rancher.Client) error {
			return client.CreateNamespace(namespace)
		},
	}
}

func newCertificatePartConverger(certificate projectModel.Certificate) descriptor.Converger {
	return descriptor.PartConverger{
		PartName: "Certificate",
		HasPart: func(client rancher.Client) (bool, error) {
			if certificate.Namespace == "" {
				return client.HasCertificate(certificate)
			}
			return client.HasNamespacedCertificate(certificate)
		},
		UpdatePart: func(client rancher.Client) error {
			if certificate.Namespace == "" {
				return client.UpgradeCertificate(certificate)
			}
			return client.UpgradeNamespacedCertificate(certificate)
		},
		CreatePart: func(client rancher.Client) error {
			if certificate.Namespace == "" {
				return client.CreateCertificate(certificate)
			}
			return client.CreateNamespacedCertificate(certificate)
		},
	}
}

func newConfigMapPartConverger(configMap projectModel.ConfigMap) descriptor.Converger {
	return descriptor.PartConverger{
		PartName: "ConfigMap",
		HasPart: func(client rancher.Client) (bool, error) {
			return client.HasConfigMap(configMap)
		},
		UpdatePart: func(client rancher.Client) error {
			return client.UpgradeConfigMap(configMap)
		},
		CreatePart: func(client rancher.Client) error {
			return client.CreateConfigMap(configMap)
		},
	}
}

func newDockerCredentialPartConverger(dockerCredential projectModel.DockerCredential) descriptor.Converger {
	return descriptor.PartConverger{
		PartName: "DockerCredential",
		HasPart: func(client rancher.Client) (bool, error) {
			if dockerCredential.Namespace == "" {
				return client.HasDockerCredential(dockerCredential)
			}
			return client.HasNamespacedDockerCredential(dockerCredential)
		},
		UpdatePart: func(client rancher.Client) error {
			if dockerCredential.Namespace == "" {
				return client.UpgradeDockerCredential(dockerCredential)
			}
			return client.UpgradeNamespacedDockerCredential(dockerCredential)
		},
		CreatePart: func(client rancher.Client) error {
			if dockerCredential.Namespace == "" {
				return client.CreateDockerCredential(dockerCredential)
			}
			return client.CreateNamespacedDockerCredential(dockerCredential)
		},
	}
}

func newSecretPartConverger(secret projectModel.ConfigMap) descriptor.Converger {
	return descriptor.PartConverger{
		PartName: "Secret",
		HasPart: func(client rancher.Client) (bool, error) {
			if secret.Namespace == "" {
				return client.HasSecret(secret)
			}
			return client.HasNamespacedSecret(secret)
		},
		UpdatePart: func(client rancher.Client) error {
			if secret.Namespace == "" {
				return client.UpgradeSecret(secret)
			}
			return client.UpgradeNamespacedSecret(secret)
		},
		CreatePart: func(client rancher.Client) error {
			if secret.Namespace == "" {
				return client.CreateSecret(secret)
			}
			return client.CreateNamespacedSecret(secret)
		},
	}
}

func newPersistentVolumePartConverger(persistentVolume projectModel.PersistentVolume) descriptor.Converger {
	return descriptor.PartConverger{
		PartName: "PersistentVolume",
		HasPart: func(client rancher.Client) (bool, error) {
			return client.HasPersistentVolume(persistentVolume)
		},
		UpdatePart: func(client rancher.Client) error {
			logrus.WithFields(logrus.Fields{
				"persistent_volume": persistentVolume.Name,
			}).Warn("Skip change existing persistent volume")
			return nil
		},
		CreatePart: func(client rancher.Client) error {
			err := client.CreatePersistentVolume(persistentVolume)
			if err != nil {
				return err
			}
			initScript := persistentVolume.InitScript
			initScript = strings.Replace(initScript, "${node}", persistentVolume.Node, -1)
			initScript = strings.Replace(initScript, "${path}", persistentVolume.Path, -1)
			logrus.Info(
				"Make sure host directory exists by running:\n",
				strings.Repeat("-", len(initScript)+4), "\n",
				strings.Repeat(" ", 4), initScript, "\n",
				strings.Repeat("-", len(initScript)+4),
			)
			return nil
		},
	}
}

func newStorageClassPartConverger(storageClass projectModel.StorageClass) descriptor.Converger {
	return descriptor.PartConverger{
		PartName: "StorageClass",
		HasPart: func(client rancher.Client) (bool, error) {
			return client.HasStorageClass(storageClass)
		},
		UpdatePart: func(client rancher.Client) error {
			logrus.WithFields(logrus.Fields{
				"storage_class": storageClass.Name,
			}).Warn("Skip change existing storage class")
			return nil
		},
		CreatePart: func(client rancher.Client) error {
			return client.CreateStorageClass(storageClass)
		},
	}
}

func newAppPartConverger(app projectModel.App) descriptor.Converger {
	return descriptor.PartConverger{
		PartName: "App",
		HasPart: func(client rancher.Client) (bool, error) {
			return client.HasApp(app)
		},
		UpdatePart: func(client rancher.Client) error {
			return client.UpgradeApp(app)
		},
		CreatePart: func(client rancher.Client) error {
			return client.CreateApp(app)
		},
	}
}
