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

package clientstub

import (
	"testing"

	"github.com/bitgrip/cattlectl/internal/pkg/assert"
	"github.com/bitgrip/cattlectl/internal/pkg/rancher"
	projectModel "github.com/bitgrip/cattlectl/internal/pkg/rancher/project/model"
)

// NewClientStub creates a rancher.Client which by default
// fails all method calls.
// Set method callbacks for expected methods.
func NewClientStub(tb testing.TB) rancher.Client {
	return ClientStub{
		DoHasClusterWithName: func(name string) (bool, string, error) {
			assert.FailInStub(tb, 2, "Unexpected call of HasClusterWithName")
			return false, "", nil
		},
		DoSetCluster: func(clusterName, clusterID string) error {
			assert.FailInStub(tb, 2, "Unexpected call of SetCluster")
			return nil
		},
		DoHasProjectWithName: func(name string) (bool, string, error) {
			assert.FailInStub(tb, 2, "Unexpected call of HasProjectWithName")
			return false, "", nil
		},
		DoSetProject: func(projectName, projectID string) error {
			assert.FailInStub(tb, 2, "Unexpected call of SetProjec")
			return nil
		},
		DoCreateProject: func(projectName string) (string, error) {
			assert.FailInStub(tb, 2, "Unexpected call of CreateProject")
			return "", nil
		},
		DoHasNamespace: func(namespace projectModel.Namespace) (bool, error) {
			assert.FailInStub(tb, 2, "Unexpected call of HasNamespace")
			return false, nil
		},
		DoCreateNamespace: func(namespace projectModel.Namespace) error {
			assert.FailInStub(tb, 2, "Unexpected call of CreateNamespace")
			return nil
		},
		DoHasCertificate: func(certificate projectModel.Certificate) (bool, error) {
			assert.FailInStub(tb, 2, "Unexpected call of HasCertificate")
			return false, nil
		},
		DoCreateCertificate: func(certificate projectModel.Certificate) error {
			assert.FailInStub(tb, 2, "Unexpected call of CreateCertificate")
			return nil
		},
		DoHasConfigMap: func(configMap projectModel.ConfigMap) (bool, error) {
			assert.FailInStub(tb, 2, "Unexpected call of HasConfigMap")
			return false, nil
		},
		DoCreateConfigMap: func(configMap projectModel.ConfigMap) error {
			assert.FailInStub(tb, 2, "Unexpected call of CreateConfigMap")
			return nil
		},
		DoUpgradeConfigMap: func(configMap projectModel.ConfigMap) error {
			assert.FailInStub(tb, 2, "Unexpected call of CreateConfigMap")
			return nil
		},
		DoHasDockerCredential: func(dockerCredential projectModel.DockerCredential) (bool, error) {
			assert.FailInStub(tb, 2, "Unexpected call of HasDockerCredential")
			return false, nil
		},
		DoCreateDockerCredential: func(dockerCredential projectModel.DockerCredential) error {
			assert.FailInStub(tb, 2, "Unexpected call of CreateDockerCredential")
			return nil
		},
		DoHasSecret: func(secret projectModel.ConfigMap) (bool, error) {
			assert.FailInStub(tb, 2, "Unexpected call of HasSecret")
			return false, nil
		},
		DoUpgradeSecret: func(secret projectModel.ConfigMap) error {
			assert.FailInStub(tb, 2, "Unexpected call of UpgradeSecret")
			return nil
		},
		DoCreateSecret: func(secret projectModel.ConfigMap) error {
			assert.FailInStub(tb, 2, "Unexpected call of CreateSecret")
			return nil
		},
		DoHasNamespacedSecret: func(secret projectModel.ConfigMap) (bool, error) {
			assert.FailInStub(tb, 2, "Unexpected call of HasNamespacedSecret")
			return false, nil
		},
		DoUpgradeNamespacedSecret: func(secret projectModel.ConfigMap) error {
			assert.FailInStub(tb, 2, "Unexpected call of UpgradeNamespacedSecret")
			return nil
		},
		DoCreateNamespacedSecret: func(secret projectModel.ConfigMap) error {
			assert.FailInStub(tb, 2, "Unexpected call of CreateNamespacedSecret")
			return nil
		},
		DoHasStorageClass: func(storageClass projectModel.StorageClass) (bool, error) {
			assert.FailInStub(tb, 2, "Unexpected call of HasStorageClass")
			return false, nil
		},
		DoCreateStorageClass: func(storageClass projectModel.StorageClass) error {
			assert.FailInStub(tb, 2, "Unexpected call of CreateStorageClass")
			return nil
		},
		DoHasPersistentVolume: func(persistentVolume projectModel.PersistentVolume) (bool, error) {
			assert.FailInStub(tb, 2, "Unexpected call of HasPersistentVolume")
			return false, nil
		},
		DoCreatePersistentVolume: func(persistentVolume projectModel.PersistentVolume) error {
			assert.FailInStub(tb, 2, "Unexpected call of CreatePersistentVolume")
			return nil
		},
		DoHasApp: func(app projectModel.App) (bool, error) {
			assert.FailInStub(tb, 2, "Unexpected call of HasApp")
			return false, nil
		},
		DoUpgradeApp: func(app projectModel.App) error {
			assert.FailInStub(tb, 2, "Unexpected call of UpgradeApp")
			return nil
		},
		DoCreateApp: func(app projectModel.App) error {
			assert.FailInStub(tb, 2, "Unexpected call of CreateApp")
			return nil
		},
		DoHasJob: func(namespace string, job projectModel.Job) (bool, error) {
			assert.FailInStub(tb, 2, "Unexpected call of HasJob")
			return false, nil
		},
		DoCreateJob: func(namespace string, job projectModel.Job) error {
			assert.FailInStub(tb, 2, "Unexpected call of CreateJob")
			return nil
		},
		DoHasCronJob: func(namespace string, cronJob projectModel.CronJob) (bool, error) {
			assert.FailInStub(tb, 2, "Unexpected call of HasCronJob")
			return false, nil
		},
		DoCreateCronJob: func(namespace string, cronJob projectModel.CronJob) error {
			assert.FailInStub(tb, 2, "Unexpected call of CreateCronJob")
			return nil
		},
		DoHasDeployment: func(namespace string, deployment projectModel.Deployment) (bool, error) {
			assert.FailInStub(tb, 2, "Unexpected call of HasDeployment")
			return false, nil
		},
		DoCreateDeployment: func(namespace string, deployment projectModel.Deployment) error {
			assert.FailInStub(tb, 2, "Unexpected call of CreateDeployment")
			return nil
		},
		DoHasDaemonSet: func(namespace string, daemonSet projectModel.DaemonSet) (bool, error) {
			assert.FailInStub(tb, 2, "Unexpected call of HasDaemonSet")
			return false, nil
		},
		DoCreateDaemonSet: func(namespace string, daemonSet projectModel.DaemonSet) error {
			assert.FailInStub(tb, 2, "Unexpected call of CreateDaemonSet")
			return nil
		},
		DoHasStatefulSet: func(namespace string, statefulSet projectModel.StatefulSet) (bool, error) {
			assert.FailInStub(tb, 2, "Unexpected call of HasStatefulSet")
			return false, nil
		},
		DoCreateStatefulSet: func(namespace string, statefulSet projectModel.StatefulSet) error {
			assert.FailInStub(tb, 2, "Unexpected call of CreateStatefulSet")
			return nil
		},
	}
}

type ClientStub struct {
	DoHasClusterWithName      func(name string) (bool, string, error)
	DoSetCluster              func(clusterName, clusterID string) error
	DoHasProjectWithName      func(name string) (bool, string, error)
	DoSetProject              func(projectName, projectID string) error
	DoCreateProject           func(projectName string) (string, error)
	DoHasNamespace            func(namespace projectModel.Namespace) (bool, error)
	DoCreateNamespace         func(namespace projectModel.Namespace) error
	DoHasCertificate          func(certificate projectModel.Certificate) (bool, error)
	DoCreateCertificate       func(certificate projectModel.Certificate) error
	DoHasConfigMap            func(configMap projectModel.ConfigMap) (bool, error)
	DoUpgradeConfigMap        func(configMap projectModel.ConfigMap) error
	DoCreateConfigMap         func(configMap projectModel.ConfigMap) error
	DoHasDockerCredential     func(dockerCredential projectModel.DockerCredential) (bool, error)
	DoCreateDockerCredential  func(dockerCredential projectModel.DockerCredential) error
	DoHasSecret               func(secret projectModel.ConfigMap) (bool, error)
	DoUpgradeSecret           func(secret projectModel.ConfigMap) error
	DoCreateSecret            func(secret projectModel.ConfigMap) error
	DoHasNamespacedSecret     func(secret projectModel.ConfigMap) (bool, error)
	DoUpgradeNamespacedSecret func(secret projectModel.ConfigMap) error
	DoCreateNamespacedSecret  func(secret projectModel.ConfigMap) error
	DoHasStorageClass         func(storageClass projectModel.StorageClass) (bool, error)
	DoCreateStorageClass      func(storageClass projectModel.StorageClass) error
	DoHasPersistentVolume     func(persistentVolume projectModel.PersistentVolume) (bool, error)
	DoCreatePersistentVolume  func(persistentVolume projectModel.PersistentVolume) error
	DoHasApp                  func(app projectModel.App) (bool, error)
	DoUpgradeApp              func(app projectModel.App) error
	DoCreateApp               func(app projectModel.App) error
	DoHasJob                  func(namespace string, job projectModel.Job) (bool, error)
	DoCreateJob               func(namespace string, job projectModel.Job) error
	DoHasCronJob              func(namespace string, cronJob projectModel.CronJob) (bool, error)
	DoCreateCronJob           func(namespace string, cronJob projectModel.CronJob) error
	DoHasDeployment           func(namespace string, deployment projectModel.Deployment) (bool, error)
	DoCreateDeployment        func(namespace string, deployment projectModel.Deployment) error
	DoHasDaemonSet            func(namespace string, daemonSet projectModel.DaemonSet) (bool, error)
	DoCreateDaemonSet         func(namespace string, daemonSet projectModel.DaemonSet) error
	DoHasStatefulSet          func(namespace string, statefulSet projectModel.StatefulSet) (bool, error)
	DoCreateStatefulSet       func(namespace string, statefulSet projectModel.StatefulSet) error
}

func (stub ClientStub) HasClusterWithName(name string) (bool, string, error) {
	return stub.DoHasClusterWithName(name)
}
func (stub ClientStub) SetCluster(clusterName, clusterID string) error {
	return stub.DoSetCluster(clusterName, clusterID)
}
func (stub ClientStub) HasProjectWithName(name string) (bool, string, error) {
	return stub.DoHasProjectWithName(name)
}
func (stub ClientStub) SetProject(projectName, projectID string) error {
	return stub.DoSetProject(projectName, projectID)
}
func (stub ClientStub) CreateProject(projectName string) (string, error) {
	return stub.DoCreateProject(projectName)
}
func (stub ClientStub) HasNamespace(namespace projectModel.Namespace) (bool, error) {
	return stub.DoHasNamespace(namespace)
}
func (stub ClientStub) CreateNamespace(namespace projectModel.Namespace) error {
	return stub.DoCreateNamespace(namespace)
}
func (stub ClientStub) HasCertificate(certificate projectModel.Certificate) (bool, error) {
	return stub.DoHasCertificate(certificate)
}
func (stub ClientStub) CreateCertificate(certificate projectModel.Certificate) error {
	return stub.DoCreateCertificate(certificate)
}
func (stub ClientStub) HasConfigMap(configMap projectModel.ConfigMap) (bool, error) {
	return stub.DoHasConfigMap(configMap)
}
func (stub ClientStub) UpgradeConfigMap(configMap projectModel.ConfigMap) error {
	return stub.DoUpgradeConfigMap(configMap)
}
func (stub ClientStub) CreateConfigMap(configMap projectModel.ConfigMap) error {
	return stub.DoCreateConfigMap(configMap)
}
func (stub ClientStub) HasDockerCredential(dockerCredential projectModel.DockerCredential) (bool, error) {
	return stub.DoHasDockerCredential(dockerCredential)
}
func (stub ClientStub) CreateDockerCredential(dockerCredential projectModel.DockerCredential) error {
	return stub.DoCreateDockerCredential(dockerCredential)
}
func (stub ClientStub) HasSecret(secret projectModel.ConfigMap) (bool, error) {
	return stub.DoHasSecret(secret)
}
func (stub ClientStub) UpgradeSecret(secret projectModel.ConfigMap) error {
	return stub.DoUpgradeSecret(secret)
}
func (stub ClientStub) CreateSecret(secret projectModel.ConfigMap) error {
	return stub.DoCreateSecret(secret)
}
func (stub ClientStub) HasNamespacedSecret(secret projectModel.ConfigMap) (bool, error) {
	return stub.DoHasNamespacedSecret(secret)
}
func (stub ClientStub) UpgradeNamespacedSecret(secret projectModel.ConfigMap) error {
	return stub.DoUpgradeNamespacedSecret(secret)
}
func (stub ClientStub) CreateNamespacedSecret(secret projectModel.ConfigMap) error {
	return stub.DoCreateNamespacedSecret(secret)
}
func (stub ClientStub) HasStorageClass(storageClass projectModel.StorageClass) (bool, error) {
	return stub.DoHasStorageClass(storageClass)
}
func (stub ClientStub) CreateStorageClass(storageClass projectModel.StorageClass) error {
	return stub.DoCreateStorageClass(storageClass)
}
func (stub ClientStub) HasPersistentVolume(persistentVolume projectModel.PersistentVolume) (bool, error) {
	return stub.DoHasPersistentVolume(persistentVolume)
}
func (stub ClientStub) CreatePersistentVolume(persistentVolume projectModel.PersistentVolume) error {
	return stub.DoCreatePersistentVolume(persistentVolume)
}
func (stub ClientStub) HasApp(app projectModel.App) (bool, error) {
	return stub.DoHasApp(app)
}
func (stub ClientStub) UpgradeApp(app projectModel.App) error {
	return stub.DoUpgradeApp(app)
}
func (stub ClientStub) CreateApp(app projectModel.App) error {
	return stub.DoCreateApp(app)
}
func (stub ClientStub) HasJob(namespace string, job projectModel.Job) (bool, error) {
	return stub.DoHasJob(namespace, job)
}
func (stub ClientStub) CreateJob(namespace string, job projectModel.Job) error {
	return stub.DoCreateJob(namespace, job)
}
func (stub ClientStub) HasCronJob(namespace string, cronJob projectModel.CronJob) (bool, error) {
	return stub.DoHasCronJob(namespace, cronJob)
}
func (stub ClientStub) CreateCronJob(namespace string, cronJob projectModel.CronJob) error {
	return stub.DoCreateCronJob(namespace, cronJob)
}
func (stub ClientStub) HasDeployment(namespace string, deployment projectModel.Deployment) (bool, error) {
	return stub.DoHasDeployment(namespace, deployment)
}
func (stub ClientStub) CreateDeployment(namespace string, deployment projectModel.Deployment) error {
	return stub.DoCreateDeployment(namespace, deployment)
}
func (stub ClientStub) HasDaemonSet(namespace string, daemonSet projectModel.DaemonSet) (bool, error) {
	return stub.DoHasDaemonSet(namespace, daemonSet)
}
func (stub ClientStub) CreateDaemonSet(namespace string, daemonSet projectModel.DaemonSet) error {
	return stub.DoCreateDaemonSet(namespace, daemonSet)
}
func (stub ClientStub) HasStatefulSet(namespace string, statefulSet projectModel.StatefulSet) (bool, error) {
	return stub.DoHasStatefulSet(namespace, statefulSet)
}
func (stub ClientStub) CreateStatefulSet(namespace string, statefulSet projectModel.StatefulSet) error {
	return stub.DoCreateStatefulSet(namespace, statefulSet)
}
