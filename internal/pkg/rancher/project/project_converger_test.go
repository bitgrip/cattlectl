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
	"io/ioutil"
	"testing"

	"github.com/bitgrip/cattlectl/internal/pkg/assert"
	"github.com/bitgrip/cattlectl/internal/pkg/rancher/clientstub"
	projectModel "github.com/bitgrip/cattlectl/internal/pkg/rancher/project/model"
	yaml "gopkg.in/yaml.v2"
)

const expectedProjectID = "test-project-id"

func newTestProjectDescriptor(t *testing.T) projectModel.Project {
	fileContent, err := ioutil.ReadFile("testdata/projects/all-resources.yaml")
	assert.Ok(t, err)
	projectDescriptor := projectModel.Project{}
	err = yaml.Unmarshal(fileContent, &projectDescriptor)
	assert.Ok(t, err)
	return projectDescriptor
}

func TestConvergeExistingProjectDescriptor(t *testing.T) {
	projectDescriptor := newTestProjectDescriptor(t)
	projectDescriptor.Namespaces = make([]projectModel.Namespace, 0)
	projectDescriptor.Resources = projectModel.Resources{}
	projectDescriptor.StorageClasses = make([]projectModel.StorageClass, 0)
	projectDescriptor.PersistentVolumes = make([]projectModel.PersistentVolume, 0)
	projectDescriptor.Apps = make([]projectModel.App, 0)

	client := clientstub.NewClientStub(t).(clientstub.ClientStub)
	expectSetClusterForExisting(projectDescriptor.Metadata.ClusterName, "cluster-id", &client, t)
	expectSetProjectForExisting(projectDescriptor.Metadata.Name, expectedProjectID, &client, t)

	converger := NewProjectConverger(projectDescriptor)
	converger.Converge(&client)
}

func TestConvergeNonExistingProjectDescriptor(t *testing.T) {
	var projectCreated = false
	var projectSet = false

	projectDescriptor := newTestProjectDescriptor(t)
	projectDescriptor.Namespaces = make([]projectModel.Namespace, 0)
	projectDescriptor.Resources = projectModel.Resources{}
	projectDescriptor.StorageClasses = make([]projectModel.StorageClass, 0)
	projectDescriptor.PersistentVolumes = make([]projectModel.PersistentVolume, 0)
	projectDescriptor.Apps = make([]projectModel.App, 0)

	client := clientstub.NewClientStub(t).(clientstub.ClientStub)
	expectSetClusterForExisting(projectDescriptor.Metadata.ClusterName, "cluster-id", &client, t)
	client.DoHasProjectWithName = func(projectName string) (bool, string, error) {
		assert.Equals(t, projectDescriptor.Metadata.Name, projectName)
		return false, "", nil
	}
	client.DoCreateProject = func(projectName string) (string, error) {
		assert.Equals(t, projectDescriptor.Metadata.Name, projectName)
		projectCreated = true
		return expectedProjectID, nil
	}
	client.DoSetProject = func(projectName, projectID string) error {
		assert.Equals(t, projectDescriptor.Metadata.Name, projectName)
		assert.Equals(t, expectedProjectID, projectID)
		projectSet = true
		return nil
	}

	converger := NewProjectConverger(projectDescriptor)
	converger.Converge(&client)
	assert.Assert(t, projectCreated, "project not created")
	assert.Assert(t, projectSet, "project not set")
}

func TestConvergeProjectDescriptorExistingNamespace(t *testing.T) {
	projectDescriptor := newTestProjectDescriptor(t)
	projectDescriptor.Resources = projectModel.Resources{}
	projectDescriptor.StorageClasses = make([]projectModel.StorageClass, 0)
	projectDescriptor.PersistentVolumes = make([]projectModel.PersistentVolume, 0)
	projectDescriptor.Apps = make([]projectModel.App, 0)

	client := clientstub.NewClientStub(t).(clientstub.ClientStub)
	expectSetClusterForExisting(projectDescriptor.Metadata.ClusterName, "cluster-id", &client, t)
	expectSetProjectForExisting(projectDescriptor.Metadata.Name, expectedProjectID, &client, t)
	client.DoHasNamespace = func(namespace projectModel.Namespace) (bool, error) {
		assert.Equals(t, projectDescriptor.Namespaces[0], namespace)
		return true, nil
	}

	converger := NewProjectConverger(projectDescriptor)
	converger.Converge(&client)
}

func TestConvergeProjectDescriptorNonExistingNamespace(t *testing.T) {
	var namespaceCreated = false
	projectDescriptor := newTestProjectDescriptor(t)
	projectDescriptor.Resources = projectModel.Resources{}
	projectDescriptor.StorageClasses = make([]projectModel.StorageClass, 0)
	projectDescriptor.PersistentVolumes = make([]projectModel.PersistentVolume, 0)
	projectDescriptor.Apps = make([]projectModel.App, 0)

	client := clientstub.NewClientStub(t).(clientstub.ClientStub)
	expectSetClusterForExisting(projectDescriptor.Metadata.ClusterName, "cluster-id", &client, t)
	expectSetProjectForExisting(projectDescriptor.Metadata.Name, expectedProjectID, &client, t)
	client.DoHasNamespace = func(namespace projectModel.Namespace) (bool, error) {
		assert.Equals(t, projectDescriptor.Namespaces[0], namespace)
		return false, nil
	}
	client.DoCreateNamespace = func(namespace projectModel.Namespace) error {
		assert.Equals(t, projectDescriptor.Namespaces[0], namespace)
		namespaceCreated = true
		return nil
	}

	converger := NewProjectConverger(projectDescriptor)
	converger.Converge(&client)
	assert.Assert(t, namespaceCreated, "namespace not created")
}

func TestConvergeProjectDescriptorExistingCertificate(t *testing.T) {
	var certificateUpgraded = false
	projectDescriptor := newTestProjectDescriptor(t)
	projectDescriptor.Namespaces = make([]projectModel.Namespace, 0)
	//projectDescriptor.Resources.Certificates = make([]projectModel.Certificate, 0)
	projectDescriptor.Resources.ConfigMaps = make([]projectModel.ConfigMap, 0)
	projectDescriptor.Resources.DockerCredentials = make([]projectModel.DockerCredential, 0)
	projectDescriptor.Resources.Secrets = make([]projectModel.ConfigMap, 0)
	projectDescriptor.StorageClasses = make([]projectModel.StorageClass, 0)
	projectDescriptor.PersistentVolumes = make([]projectModel.PersistentVolume, 0)
	projectDescriptor.Apps = make([]projectModel.App, 0)

	client := clientstub.NewClientStub(t).(clientstub.ClientStub)
	expectSetClusterForExisting(projectDescriptor.Metadata.ClusterName, "cluster-id", &client, t)
	expectSetProjectForExisting(projectDescriptor.Metadata.Name, expectedProjectID, &client, t)
	client.DoHasCertificate = func(certificate projectModel.Certificate) (bool, error) {
		assert.Equals(t, projectDescriptor.Resources.Certificates[0], certificate)
		return true, nil
	}
	client.DoUpgradeCertificate = func(certificate projectModel.Certificate) error {
		assert.Equals(t, projectDescriptor.Resources.Certificates[0], certificate)
		certificateUpgraded = true
		return nil
	}

	converger := NewProjectConverger(projectDescriptor)
	converger.Converge(&client)
	assert.Assert(t, certificateUpgraded, "certificate not upgraded")
}

func TestConvergeProjectDescriptorNonExistingCertificate(t *testing.T) {
	var certificateCreated = false
	projectDescriptor := newTestProjectDescriptor(t)
	projectDescriptor.Namespaces = make([]projectModel.Namespace, 0)
	//projectDescriptor.Resources.Certificates = make([]projectModel.Certificate, 0)
	projectDescriptor.Resources.ConfigMaps = make([]projectModel.ConfigMap, 0)
	projectDescriptor.Resources.DockerCredentials = make([]projectModel.DockerCredential, 0)
	projectDescriptor.Resources.Secrets = make([]projectModel.ConfigMap, 0)
	projectDescriptor.StorageClasses = make([]projectModel.StorageClass, 0)
	projectDescriptor.PersistentVolumes = make([]projectModel.PersistentVolume, 0)
	projectDescriptor.Apps = make([]projectModel.App, 0)

	client := clientstub.NewClientStub(t).(clientstub.ClientStub)
	expectSetClusterForExisting(projectDescriptor.Metadata.ClusterName, "cluster-id", &client, t)
	expectSetProjectForExisting(projectDescriptor.Metadata.Name, expectedProjectID, &client, t)
	client.DoHasCertificate = func(certificate projectModel.Certificate) (bool, error) {
		assert.Equals(t, projectDescriptor.Resources.Certificates[0], certificate)
		return false, nil
	}
	client.DoCreateCertificate = func(certificate projectModel.Certificate) error {
		assert.Equals(t, projectDescriptor.Resources.Certificates[0], certificate)
		certificateCreated = true
		return nil
	}

	converger := NewProjectConverger(projectDescriptor)
	converger.Converge(&client)
	assert.Assert(t, certificateCreated, "certificate not created")
}

func TestConvergeProjectDescriptorExistingConfigMap(t *testing.T) {
	var configMapUpgraded = false
	projectDescriptor := newTestProjectDescriptor(t)
	projectDescriptor.Namespaces = make([]projectModel.Namespace, 0)
	projectDescriptor.Resources.Certificates = make([]projectModel.Certificate, 0)
	//projectDescriptor.Resources.ConfigMaps = make([]projectModel.ConfigMap, 0)
	projectDescriptor.Resources.DockerCredentials = make([]projectModel.DockerCredential, 0)
	projectDescriptor.Resources.Secrets = make([]projectModel.ConfigMap, 0)
	projectDescriptor.StorageClasses = make([]projectModel.StorageClass, 0)
	projectDescriptor.PersistentVolumes = make([]projectModel.PersistentVolume, 0)
	projectDescriptor.Apps = make([]projectModel.App, 0)

	client := clientstub.NewClientStub(t).(clientstub.ClientStub)
	expectSetClusterForExisting(projectDescriptor.Metadata.ClusterName, "cluster-id", &client, t)
	expectSetProjectForExisting(projectDescriptor.Metadata.Name, expectedProjectID, &client, t)
	client.DoHasConfigMap = func(configMap projectModel.ConfigMap) (bool, error) {
		assert.Equals(t, projectDescriptor.Resources.ConfigMaps[0], configMap)
		return true, nil
	}
	client.DoUpgradeConfigMap = func(configMap projectModel.ConfigMap) error {
		assert.Equals(t, projectDescriptor.Resources.ConfigMaps[0], configMap)
		configMapUpgraded = true
		return nil
	}

	converger := NewProjectConverger(projectDescriptor)
	converger.Converge(&client)
	assert.Assert(t, configMapUpgraded, "configMap not upgraded")
}

func TestConvergeProjectDescriptorNonExistingConfigMap(t *testing.T) {
	var configMapCreated = false
	projectDescriptor := newTestProjectDescriptor(t)
	projectDescriptor.Namespaces = make([]projectModel.Namespace, 0)
	projectDescriptor.Resources.Certificates = make([]projectModel.Certificate, 0)
	//projectDescriptor.Resources.ConfigMaps = make([]projectModel.ConfigMap, 0)
	projectDescriptor.Resources.DockerCredentials = make([]projectModel.DockerCredential, 0)
	projectDescriptor.Resources.Secrets = make([]projectModel.ConfigMap, 0)
	projectDescriptor.StorageClasses = make([]projectModel.StorageClass, 0)
	projectDescriptor.PersistentVolumes = make([]projectModel.PersistentVolume, 0)
	projectDescriptor.Apps = make([]projectModel.App, 0)

	client := clientstub.NewClientStub(t).(clientstub.ClientStub)
	expectSetClusterForExisting(projectDescriptor.Metadata.ClusterName, "cluster-id", &client, t)
	expectSetProjectForExisting(projectDescriptor.Metadata.Name, expectedProjectID, &client, t)
	client.DoHasConfigMap = func(configMap projectModel.ConfigMap) (bool, error) {
		assert.Equals(t, projectDescriptor.Resources.ConfigMaps[0], configMap)
		return false, nil
	}
	client.DoCreateConfigMap = func(configMap projectModel.ConfigMap) error {
		assert.Equals(t, projectDescriptor.Resources.ConfigMaps[0], configMap)
		configMapCreated = true
		return nil
	}

	converger := NewProjectConverger(projectDescriptor)
	converger.Converge(&client)
	assert.Assert(t, configMapCreated, "configMap not created")
}

func TestConvergeProjectDescriptorExistingDockerCredential(t *testing.T) {
	var dockerCredentialUpgraded = false
	projectDescriptor := newTestProjectDescriptor(t)
	projectDescriptor.Namespaces = make([]projectModel.Namespace, 0)
	projectDescriptor.Resources.Certificates = make([]projectModel.Certificate, 0)
	projectDescriptor.Resources.ConfigMaps = make([]projectModel.ConfigMap, 0)
	//projectDescriptor.Resources.DockerCredentials = make([]projectModel.DockerCredential, 0)
	projectDescriptor.Resources.Secrets = make([]projectModel.ConfigMap, 0)
	projectDescriptor.StorageClasses = make([]projectModel.StorageClass, 0)
	projectDescriptor.PersistentVolumes = make([]projectModel.PersistentVolume, 0)
	projectDescriptor.Apps = make([]projectModel.App, 0)

	client := clientstub.NewClientStub(t).(clientstub.ClientStub)
	expectSetClusterForExisting(projectDescriptor.Metadata.ClusterName, "cluster-id", &client, t)
	expectSetProjectForExisting(projectDescriptor.Metadata.Name, expectedProjectID, &client, t)
	client.DoHasDockerCredential = func(dockerCredential projectModel.DockerCredential) (bool, error) {
		assert.Equals(t, projectDescriptor.Resources.DockerCredentials[0], dockerCredential)
		return true, nil
	}
	client.DoUpgradeDockerCredential = func(dockerCredential projectModel.DockerCredential) error {
		assert.Equals(t, projectDescriptor.Resources.DockerCredentials[0], dockerCredential)
		dockerCredentialUpgraded = true
		return nil
	}

	converger := NewProjectConverger(projectDescriptor)
	converger.Converge(&client)
	assert.Assert(t, dockerCredentialUpgraded, "dockerCredential not upgraded")
}

func TestConvergeProjectDescriptorNonExistingDockerCredential(t *testing.T) {
	var dockerCredentialCreated = false
	projectDescriptor := newTestProjectDescriptor(t)
	projectDescriptor.Namespaces = make([]projectModel.Namespace, 0)
	projectDescriptor.Resources.Certificates = make([]projectModel.Certificate, 0)
	projectDescriptor.Resources.ConfigMaps = make([]projectModel.ConfigMap, 0)
	//projectDescriptor.Resources.DockerCredentials = make([]projectModel.DockerCredential, 0)
	projectDescriptor.Resources.Secrets = make([]projectModel.ConfigMap, 0)
	projectDescriptor.StorageClasses = make([]projectModel.StorageClass, 0)
	projectDescriptor.PersistentVolumes = make([]projectModel.PersistentVolume, 0)
	projectDescriptor.Apps = make([]projectModel.App, 0)

	client := clientstub.NewClientStub(t).(clientstub.ClientStub)
	expectSetClusterForExisting(projectDescriptor.Metadata.ClusterName, "cluster-id", &client, t)
	expectSetProjectForExisting(projectDescriptor.Metadata.Name, expectedProjectID, &client, t)
	client.DoHasDockerCredential = func(dockerCredential projectModel.DockerCredential) (bool, error) {
		assert.Equals(t, projectDescriptor.Resources.DockerCredentials[0], dockerCredential)
		return false, nil
	}
	client.DoCreateDockerCredential = func(dockerCredential projectModel.DockerCredential) error {
		assert.Equals(t, projectDescriptor.Resources.DockerCredentials[0], dockerCredential)
		dockerCredentialCreated = true
		return nil
	}

	converger := NewProjectConverger(projectDescriptor)
	converger.Converge(&client)
	assert.Assert(t, dockerCredentialCreated, "dockerCredential not created")
}

func TestConvergeProjectDescriptorExistingSecret(t *testing.T) {
	var secretUpgraded = false
	projectDescriptor := newTestProjectDescriptor(t)
	projectDescriptor.Namespaces = make([]projectModel.Namespace, 0)
	projectDescriptor.Resources.Certificates = make([]projectModel.Certificate, 0)
	projectDescriptor.Resources.ConfigMaps = make([]projectModel.ConfigMap, 0)
	projectDescriptor.Resources.DockerCredentials = make([]projectModel.DockerCredential, 0)
	//projectDescriptor.Resources.Secrets = make([]projectModel.ConfigMap, 0)
	projectDescriptor.StorageClasses = make([]projectModel.StorageClass, 0)
	projectDescriptor.PersistentVolumes = make([]projectModel.PersistentVolume, 0)
	projectDescriptor.Apps = make([]projectModel.App, 0)

	client := clientstub.NewClientStub(t).(clientstub.ClientStub)
	expectSetClusterForExisting(projectDescriptor.Metadata.ClusterName, "cluster-id", &client, t)
	expectSetProjectForExisting(projectDescriptor.Metadata.Name, expectedProjectID, &client, t)
	client.DoHasSecret = func(secret projectModel.ConfigMap) (bool, error) {
		assert.Equals(t, projectDescriptor.Resources.Secrets[0], secret)
		return true, nil
	}
	client.DoUpgradeSecret = func(secret projectModel.ConfigMap) error {
		assert.Equals(t, projectDescriptor.Resources.Secrets[0], secret)
		secretUpgraded = true
		return nil
	}

	converger := NewProjectConverger(projectDescriptor)
	converger.Converge(&client)
	assert.Assert(t, secretUpgraded, "secret not upgraded")
}

func TestConvergeProjectDescriptorNonExistingSecret(t *testing.T) {
	var secretCreated = false
	projectDescriptor := newTestProjectDescriptor(t)
	projectDescriptor.Namespaces = make([]projectModel.Namespace, 0)
	projectDescriptor.Resources.Certificates = make([]projectModel.Certificate, 0)
	projectDescriptor.Resources.ConfigMaps = make([]projectModel.ConfigMap, 0)
	projectDescriptor.Resources.DockerCredentials = make([]projectModel.DockerCredential, 0)
	//projectDescriptor.Resources.Secrets = make([]projectModel.ConfigMap, 0)
	projectDescriptor.StorageClasses = make([]projectModel.StorageClass, 0)
	projectDescriptor.PersistentVolumes = make([]projectModel.PersistentVolume, 0)
	projectDescriptor.Apps = make([]projectModel.App, 0)

	client := clientstub.NewClientStub(t).(clientstub.ClientStub)
	expectSetClusterForExisting(projectDescriptor.Metadata.ClusterName, "cluster-id", &client, t)
	expectSetProjectForExisting(projectDescriptor.Metadata.Name, expectedProjectID, &client, t)
	client.DoHasSecret = func(secret projectModel.ConfigMap) (bool, error) {
		assert.Equals(t, projectDescriptor.Resources.Secrets[0], secret)
		return false, nil
	}
	client.DoCreateSecret = func(secret projectModel.ConfigMap) error {
		assert.Equals(t, projectDescriptor.Resources.Secrets[0], secret)
		secretCreated = true
		return nil
	}

	converger := NewProjectConverger(projectDescriptor)
	converger.Converge(&client)
	assert.Assert(t, secretCreated, "secret not created")
}

func TestConvergeProjectDescriptorExistingStorageClass(t *testing.T) {
	projectDescriptor := newTestProjectDescriptor(t)
	projectDescriptor.Namespaces = make([]projectModel.Namespace, 0)
	projectDescriptor.Resources.Certificates = make([]projectModel.Certificate, 0)
	projectDescriptor.Resources.ConfigMaps = make([]projectModel.ConfigMap, 0)
	projectDescriptor.Resources.DockerCredentials = make([]projectModel.DockerCredential, 0)
	projectDescriptor.Resources.Secrets = make([]projectModel.ConfigMap, 0)
	//projectDescriptor.StorageClasses = make([]projectModel.StorageClass, 0)
	projectDescriptor.PersistentVolumes = make([]projectModel.PersistentVolume, 0)
	projectDescriptor.Apps = make([]projectModel.App, 0)

	client := clientstub.NewClientStub(t).(clientstub.ClientStub)
	expectSetClusterForExisting(projectDescriptor.Metadata.ClusterName, "cluster-id", &client, t)
	expectSetProjectForExisting(projectDescriptor.Metadata.Name, expectedProjectID, &client, t)
	client.DoHasStorageClass = func(storageClass projectModel.StorageClass) (bool, error) {
		assert.Equals(t, projectDescriptor.StorageClasses[0], storageClass)
		return true, nil
	}

	converger := NewProjectConverger(projectDescriptor)
	converger.Converge(&client)
}

func TestConvergeProjectDescriptorNonExistingStorageClass(t *testing.T) {
	var storageClassCreated = false
	projectDescriptor := newTestProjectDescriptor(t)
	projectDescriptor.Namespaces = make([]projectModel.Namespace, 0)
	projectDescriptor.Resources.Certificates = make([]projectModel.Certificate, 0)
	projectDescriptor.Resources.ConfigMaps = make([]projectModel.ConfigMap, 0)
	projectDescriptor.Resources.DockerCredentials = make([]projectModel.DockerCredential, 0)
	projectDescriptor.Resources.Secrets = make([]projectModel.ConfigMap, 0)
	//projectDescriptor.StorageClasses = make([]projectModel.StorageClass, 0)
	projectDescriptor.PersistentVolumes = make([]projectModel.PersistentVolume, 0)
	projectDescriptor.Apps = make([]projectModel.App, 0)

	client := clientstub.NewClientStub(t).(clientstub.ClientStub)
	expectSetClusterForExisting(projectDescriptor.Metadata.ClusterName, "cluster-id", &client, t)
	expectSetProjectForExisting(projectDescriptor.Metadata.Name, expectedProjectID, &client, t)
	client.DoHasStorageClass = func(storageClass projectModel.StorageClass) (bool, error) {
		assert.Equals(t, projectDescriptor.StorageClasses[0], storageClass)
		return false, nil
	}
	client.DoCreateStorageClass = func(storageClass projectModel.StorageClass) error {
		assert.Equals(t, projectDescriptor.StorageClasses[0], storageClass)
		storageClassCreated = true
		return nil
	}

	converger := NewProjectConverger(projectDescriptor)
	converger.Converge(&client)
	assert.Assert(t, storageClassCreated, "storageClass not created")
}

func TestConvergeProjectDescriptorExistingPersistentVolume(t *testing.T) {
	projectDescriptor := newTestProjectDescriptor(t)
	projectDescriptor.Namespaces = make([]projectModel.Namespace, 0)
	projectDescriptor.Resources.Certificates = make([]projectModel.Certificate, 0)
	projectDescriptor.Resources.ConfigMaps = make([]projectModel.ConfigMap, 0)
	projectDescriptor.Resources.DockerCredentials = make([]projectModel.DockerCredential, 0)
	projectDescriptor.Resources.Secrets = make([]projectModel.ConfigMap, 0)
	projectDescriptor.StorageClasses = make([]projectModel.StorageClass, 0)
	//projectDescriptor.PersistentVolumes = make([]projectModel.PersistentVolume, 0)
	projectDescriptor.Apps = make([]projectModel.App, 0)

	client := clientstub.NewClientStub(t).(clientstub.ClientStub)
	expectSetClusterForExisting(projectDescriptor.Metadata.ClusterName, "cluster-id", &client, t)
	expectSetProjectForExisting(projectDescriptor.Metadata.Name, expectedProjectID, &client, t)
	client.DoHasPersistentVolume = func(persistentVolume projectModel.PersistentVolume) (bool, error) {
		assert.Equals(t, projectDescriptor.PersistentVolumes[0], persistentVolume)
		return true, nil
	}

	converger := NewProjectConverger(projectDescriptor)
	converger.Converge(&client)
}

func TestConvergeProjectDescriptorNonExistingPersistentVolume(t *testing.T) {
	var persistentVolumeCreated = false
	projectDescriptor := newTestProjectDescriptor(t)
	projectDescriptor.Namespaces = make([]projectModel.Namespace, 0)
	projectDescriptor.Resources.Certificates = make([]projectModel.Certificate, 0)
	projectDescriptor.Resources.ConfigMaps = make([]projectModel.ConfigMap, 0)
	projectDescriptor.Resources.DockerCredentials = make([]projectModel.DockerCredential, 0)
	projectDescriptor.Resources.Secrets = make([]projectModel.ConfigMap, 0)
	projectDescriptor.StorageClasses = make([]projectModel.StorageClass, 0)
	//projectDescriptor.PersistentVolumes = make([]projectModel.PersistentVolume, 0)
	projectDescriptor.Apps = make([]projectModel.App, 0)

	client := clientstub.NewClientStub(t).(clientstub.ClientStub)
	expectSetClusterForExisting(projectDescriptor.Metadata.ClusterName, "cluster-id", &client, t)
	expectSetProjectForExisting(projectDescriptor.Metadata.Name, expectedProjectID, &client, t)
	client.DoHasPersistentVolume = func(persistentVolume projectModel.PersistentVolume) (bool, error) {
		assert.Equals(t, projectDescriptor.PersistentVolumes[0], persistentVolume)
		return false, nil
	}
	client.DoCreatePersistentVolume = func(persistentVolume projectModel.PersistentVolume) error {
		assert.Equals(t, projectDescriptor.PersistentVolumes[0], persistentVolume)
		persistentVolumeCreated = true
		return nil
	}

	converger := NewProjectConverger(projectDescriptor)
	converger.Converge(&client)
	assert.Assert(t, persistentVolumeCreated, "persistentVolume not created")
}

func TestConvergeProjectDescriptorExistingApp(t *testing.T) {
	var appUpgraded = false
	projectDescriptor := newTestProjectDescriptor(t)
	projectDescriptor.Namespaces = make([]projectModel.Namespace, 0)
	projectDescriptor.Resources.Certificates = make([]projectModel.Certificate, 0)
	projectDescriptor.Resources.ConfigMaps = make([]projectModel.ConfigMap, 0)
	projectDescriptor.Resources.DockerCredentials = make([]projectModel.DockerCredential, 0)
	projectDescriptor.Resources.Secrets = make([]projectModel.ConfigMap, 0)
	projectDescriptor.StorageClasses = make([]projectModel.StorageClass, 0)
	projectDescriptor.PersistentVolumes = make([]projectModel.PersistentVolume, 0)
	//projectDescriptor.Apps = make([]projectModel.App, 0)

	client := clientstub.NewClientStub(t).(clientstub.ClientStub)
	expectSetClusterForExisting(projectDescriptor.Metadata.ClusterName, "cluster-id", &client, t)
	expectSetProjectForExisting(projectDescriptor.Metadata.Name, expectedProjectID, &client, t)
	client.DoHasApp = func(app projectModel.App) (bool, error) {
		assert.Equals(t, projectDescriptor.Apps[0], app)
		return true, nil
	}
	client.DoUpgradeApp = func(app projectModel.App) error {
		assert.Equals(t, projectDescriptor.Apps[0], app)
		appUpgraded = true
		return nil
	}

	converger := NewProjectConverger(projectDescriptor)
	converger.Converge(&client)
	assert.Assert(t, appUpgraded, "app not upgraded")
}

func TestConvergeProjectDescriptorNonExistingApp(t *testing.T) {
	var appCreated = false
	projectDescriptor := newTestProjectDescriptor(t)
	projectDescriptor.Namespaces = make([]projectModel.Namespace, 0)
	projectDescriptor.Resources.Certificates = make([]projectModel.Certificate, 0)
	projectDescriptor.Resources.ConfigMaps = make([]projectModel.ConfigMap, 0)
	projectDescriptor.Resources.DockerCredentials = make([]projectModel.DockerCredential, 0)
	projectDescriptor.Resources.Secrets = make([]projectModel.ConfigMap, 0)
	projectDescriptor.StorageClasses = make([]projectModel.StorageClass, 0)
	projectDescriptor.PersistentVolumes = make([]projectModel.PersistentVolume, 0)
	//projectDescriptor.Apps = make([]projectModel.App, 0)

	client := clientstub.NewClientStub(t).(clientstub.ClientStub)
	expectSetClusterForExisting(projectDescriptor.Metadata.ClusterName, "cluster-id", &client, t)
	expectSetProjectForExisting(projectDescriptor.Metadata.Name, expectedProjectID, &client, t)
	client.DoHasApp = func(app projectModel.App) (bool, error) {
		assert.Equals(t, projectDescriptor.Apps[0], app)
		return false, nil
	}
	client.DoCreateApp = func(app projectModel.App) error {
		assert.Equals(t, projectDescriptor.Apps[0], app)
		appCreated = true
		return nil
	}

	converger := NewProjectConverger(projectDescriptor)
	converger.Converge(&client)
	assert.Assert(t, appCreated, "app not created")
}
