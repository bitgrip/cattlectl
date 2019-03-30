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
	return projectConverger{
		project: project,
	}
}

type projectConverger struct {
	project projectModel.Project
}

func (converger projectConverger) Converge(client rancher.Client) error {
	rancher.InitCluster(
		converger.project.Metadata.ClusterID,
		converger.project.Metadata.ClusterName,
		client,
	)
	if err := converger.applyProject(client); err != nil {
		return err
	}
	for _, namespace := range converger.project.Namespaces {
		if err := converger.applyNamespace(namespace, client); err != nil {
			return err
		}
	}
	if err := converger.applyResources(client); err != nil {
		return err
	}
	for _, storageClass := range converger.project.StorageClasses {
		if err := converger.applyStorageClass(storageClass, client); err != nil {
			return err
		}
	}
	for _, persistentVolume := range converger.project.PersistentVolumes {
		if err := converger.applyPersistentVolume(persistentVolume, client); err != nil {
			return err
		}
	}
	for _, app := range converger.project.Apps {
		if err := converger.applyApp(app, client); err != nil {
			return err
		}
	}
	return nil
}

func (converger projectConverger) applyProject(client rancher.Client) error {
	if hasProject, projectID, err := client.HasProjectWithName(converger.project.Metadata.Name); hasProject {
		return client.SetProject(converger.project.Metadata.Name, projectID)
	} else if err != nil {
		return fmt.Errorf("Failed to check for project, %v", err)
	}
	projectID, err := client.CreateProject(converger.project.Metadata.Name)
	if err != nil {
		return err
	}
	return client.SetProject(converger.project.Metadata.Name, projectID)
}

func (converger projectConverger) applyNamespace(namespace projectModel.Namespace, client rancher.Client) error {
	if hasNamespace, err := client.HasNamespace(namespace); hasNamespace {
		return nil
	} else if err != nil {
		return fmt.Errorf("Failed to check for namespace, %v", err)
	}
	return client.CreateNamespace(namespace)
}

func (converger projectConverger) applyResources(client rancher.Client) error {
	resources := converger.project.Resources
	for _, certificate := range resources.Certificates {
		if err := converger.applyCertificate(certificate, client); err != nil {
			return err
		}
	}
	for _, configMap := range resources.ConfigMaps {
		if err := converger.applyConfigMap(configMap, client); err != nil {
			return err
		}
	}
	for _, dockerCredential := range resources.DockerCredentials {
		if err := converger.applyDockerCredential(dockerCredential, client); err != nil {
			return err
		}
	}
	for _, secret := range resources.Secrets {
		if err := converger.applySecret(secret, client); err != nil {
			return err
		}
	}
	return nil
}

func (converger projectConverger) applyCertificate(certificate projectModel.Certificate, client rancher.Client) error {
	if hasCertificate, err := client.HasCertificate(certificate); hasCertificate {
		logrus.WithField("certificate", certificate.Name).Debug("Skip change existing certificate")
		return nil
	} else if err != nil {
		return fmt.Errorf("Failed to check for certificate, %v", err)
	}
	err := client.CreateCertificate(certificate)
	if err != nil {
		return err
	}
	return nil
}

func (converger projectConverger) applyConfigMap(configMap projectModel.ConfigMap, client rancher.Client) error {
	if hasConfigMap, err := client.HasConfigMap(configMap); hasConfigMap {
		logrus.WithField("configMap", configMap.Name).Debug("Skip change existing configMap")
		return nil
	} else if err != nil {
		return fmt.Errorf("Failed to check for configMap, %v", err)
	}
	err := client.CreateConfigMap(configMap)
	if err != nil {
		return err
	}
	return nil
}

func (converger projectConverger) applyDockerCredential(dockerCredential projectModel.DockerCredential, client rancher.Client) error {
	if hasDockerCredential, err := client.HasDockerCredential(dockerCredential); hasDockerCredential {
		logrus.WithField("dockerCredential", dockerCredential.Name).Debug("Skip change existing dockerCredential")
		return nil
	} else if err != nil {
		return fmt.Errorf("Failed to check for dockerCredential, %v", err)
	}
	err := client.CreateDockerCredential(dockerCredential)
	if err != nil {
		return err
	}
	return nil
}

func (converger projectConverger) applySecret(secret projectModel.ConfigMap, client rancher.Client) error {
	var hasSecret bool
	var err error
	if secret.Namespace == "" {
		hasSecret, err = client.HasSecret(secret)
	} else {
		hasSecret, err = client.HasNamespacedSecret(secret)
	}
	if hasSecret {
		logrus.WithField("secret", secret.Name).Debug("Skip change existing secret")
		return nil
	} else if err != nil {
		return fmt.Errorf("Failed to check for secret, %v", err)
	}
	if secret.Namespace == "" {
		err = client.CreateSecret(secret)
	} else {
		err = client.CreateNamespacedSecret(secret)
	}
	if err != nil {
		return err
	}
	return nil
}

func (converger projectConverger) applyPersistentVolume(persistentVolume projectModel.PersistentVolume, client rancher.Client) error {
	if hasVolume, err := client.HasPersistentVolume(persistentVolume); hasVolume {
		logrus.WithField("persistent_volume", persistentVolume.Name).Debug("Skip change existing persistent volume")
		return nil
	} else if err != nil {
		return fmt.Errorf("Failed to check for local volume, %v", err)
	}
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
}

func (converger projectConverger) applyStorageClass(storageClass projectModel.StorageClass, client rancher.Client) error {
	if hasStorageClass, err := client.HasStorageClass(storageClass); hasStorageClass {
		logrus.WithField("storage_class", storageClass.Name).Debug("Skip change existing storage class")
		return nil
	} else if err != nil {
		return fmt.Errorf("Failed to check for storage class, %v", err)
	}
	err := client.CreateStorageClass(storageClass)
	if err != nil {
		return err
	}
	return nil
}

func (converger projectConverger) applyApp(app projectModel.App, client rancher.Client) error {
	if hasApp, err := client.HasApp(app); hasApp {
		return client.UpgradeApp(app)
	} else if err == nil {
		return client.CreateApp(app)
	} else {
		return fmt.Errorf("Failed to check for app, %v", err)
	}
}
