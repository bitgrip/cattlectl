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
	projectModel "github.com/bitgrip/cattlectl/internal/pkg/rancher/project/model"
	"github.com/sirupsen/logrus"
)

// NewProjectConverger creates a Converger for a given github.com/bitgrip/cattlectl/internal/pkg/projectModel.Project
func NewProjectConverger(project projectModel.Project) Converger {
	client, err := newRancherClient(rancher.ClientConfig{
		RancherURL: project.Metadata.RancherURL,
		AccessKey:  project.Metadata.AccessKey,
		SecretKey:  project.Metadata.SecretKey,
	})
	if err != nil {
		logrus.WithError(err).Fatal("Error creating rancher client")
	}
	return projectConverger{
		project: project,
		client:  client,
	}
}

type projectConverger struct {
	project projectModel.Project
	client  rancher.Client
}

func (converger projectConverger) Converge() error {
	if err := converger.initCluster(); err != nil {
		return err
	}
	if err := converger.applyProject(); err != nil {
		return err
	}
	for _, namespace := range converger.project.Namespaces {
		if err := converger.applyNamespace(namespace); err != nil {
			return err
		}
	}
	if err := converger.applyResources(); err != nil {
		return err
	}
	for _, storageClass := range converger.project.StorageClasses {
		if err := converger.applyStorageClass(storageClass); err != nil {
			return err
		}
	}
	for _, persistentVolume := range converger.project.PersistentVolumes {
		if err := converger.applyPersistentVolume(persistentVolume); err != nil {
			return err
		}
	}
	for _, app := range converger.project.Apps {
		if err := converger.applyApp(app); err != nil {
			return err
		}
	}
	return nil
}

func (converger projectConverger) initCluster() error {
	if converger.project.Metadata.ClusterID != "" {
		if err := converger.client.SetCluster(converger.project.Metadata.ClusterName, converger.project.Metadata.ClusterID); err != nil {
			return fmt.Errorf("Failed to init cluster, %v", err)
		}
		return nil
	}
	if hasCluster, clusterID, err := converger.client.HasClusterWithName(converger.project.Metadata.ClusterName); hasCluster {
		if err = converger.client.SetCluster(converger.project.Metadata.ClusterName, clusterID); err != nil {
			return fmt.Errorf("Failed to init cluster, %v", err)
		}
		return nil
	} else if err != nil {
		return fmt.Errorf("Failed to init cluster, %v", err)
	} else {
		return fmt.Errorf("Cluster not found")
	}
}

func (converger projectConverger) applyProject() error {
	if hasProject, projectID, err := converger.client.HasProjectWithName(converger.project.Metadata.Name); hasProject {
		return converger.client.SetProject(converger.project.Metadata.Name, projectID)
	} else if err != nil {
		return fmt.Errorf("Failed to check for project, %v", err)
	}
	projectID, err := converger.client.CreateProject(converger.project.Metadata.Name)
	if err != nil {
		return err
	}
	return converger.client.SetProject(converger.project.Metadata.Name, projectID)
}

func (converger projectConverger) applyNamespace(namespace projectModel.Namespace) error {
	if hasNamespace, err := converger.client.HasNamespace(namespace); hasNamespace {
		return nil
	} else if err != nil {
		return fmt.Errorf("Failed to check for namespace, %v", err)
	}
	return converger.client.CreateNamespace(namespace)
}

func (converger projectConverger) applyResources() error {
	resources := converger.project.Resources
	for _, certificate := range resources.Certificates {
		if err := converger.applyCertificate(certificate); err != nil {
			return err
		}
	}
	for _, configMap := range resources.ConfigMaps {
		if err := converger.applyConfigMap(configMap); err != nil {
			return err
		}
	}
	for _, dockerCredential := range resources.DockerCredentials {
		if err := converger.applyDockerCredential(dockerCredential); err != nil {
			return err
		}
	}
	for _, secret := range resources.Secrets {
		if err := converger.applySecret(secret); err != nil {
			return err
		}
	}
	return nil
}

func (converger projectConverger) applyCertificate(certificate projectModel.Certificate) error {
	if hasCertificate, err := converger.client.HasCertificate(certificate); hasCertificate {
		logrus.WithField("certificate", certificate.Name).Debug("Skip change existing certificate")
		return nil
	} else if err != nil {
		return fmt.Errorf("Failed to check for certificate, %v", err)
	}
	err := converger.client.CreateCertificate(certificate)
	if err != nil {
		return err
	}
	return nil
}

func (converger projectConverger) applyConfigMap(configMap projectModel.ConfigMap) error {
	if hasConfigMap, err := converger.client.HasConfigMap(configMap); hasConfigMap {
		logrus.WithField("configMap", configMap.Name).Debug("Skip change existing configMap")
		return nil
	} else if err != nil {
		return fmt.Errorf("Failed to check for configMap, %v", err)
	}
	err := converger.client.CreateConfigMap(configMap)
	if err != nil {
		return err
	}
	return nil
}

func (converger projectConverger) applyDockerCredential(dockerCredential projectModel.DockerCredential) error {
	if hasDockerCredential, err := converger.client.HasDockerCredential(dockerCredential); hasDockerCredential {
		logrus.WithField("dockerCredential", dockerCredential.Name).Debug("Skip change existing dockerCredential")
		return nil
	} else if err != nil {
		return fmt.Errorf("Failed to check for dockerCredential, %v", err)
	}
	err := converger.client.CreateDockerCredential(dockerCredential)
	if err != nil {
		return err
	}
	return nil
}

func (converger projectConverger) applySecret(secret projectModel.ConfigMap) error {
	var hasSecret bool
	var err error
	if secret.Namespace == "" {
		hasSecret, err = converger.client.HasSecret(secret)
	} else {
		hasSecret, err = converger.client.HasNamespacedSecret(secret)
	}
	if hasSecret {
		logrus.WithField("secret", secret.Name).Debug("Skip change existing secret")
		return nil
	} else if err != nil {
		return fmt.Errorf("Failed to check for secret, %v", err)
	}
	if secret.Namespace == "" {
		err = converger.client.CreateSecret(secret)
	} else {
		err = converger.client.CreateNamespacedSecret(secret)
	}
	if err != nil {
		return err
	}
	return nil
}

func (converger projectConverger) applyPersistentVolume(persistentVolume projectModel.PersistentVolume) error {
	if hasVolume, err := converger.client.HasPersistentVolume(persistentVolume); hasVolume {
		logrus.WithField("persistent_volume", persistentVolume.Name).Debug("Skip change existing persistent volume")
		return nil
	} else if err != nil {
		return fmt.Errorf("Failed to check for local volume, %v", err)
	}
	err := converger.client.CreatePersistentVolume(persistentVolume)
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

func (converger projectConverger) applyStorageClass(storageClass projectModel.StorageClass) error {
	if hasStorageClass, err := converger.client.HasStorageClass(storageClass); hasStorageClass {
		logrus.WithField("storage_class", storageClass.Name).Debug("Skip change existing storage class")
		return nil
	} else if err != nil {
		return fmt.Errorf("Failed to check for storage class, %v", err)
	}
	err := converger.client.CreateStorageClass(storageClass)
	if err != nil {
		return err
	}
	return nil
}

func (converger projectConverger) applyApp(app projectModel.App) error {
	if hasApp, err := converger.client.HasApp(app); hasApp {
		return converger.client.UpgradeApp(app)
	} else if err == nil {
		return converger.client.CreateApp(app)
	} else {
		return fmt.Errorf("Failed to check for app, %v", err)
	}
}
