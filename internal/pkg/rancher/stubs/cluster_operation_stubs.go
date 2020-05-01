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

package stubs

import (
	"testing"

	"github.com/bitgrip/cattlectl/internal/pkg/assert"
	"github.com/rancher/norman/types"
	managementClient "github.com/rancher/types/client/management/v3"
)

// CreateClusterOperationsStub creates a stub of github.com/rancher/types/client/management/v3/ClusterOperations
func CreateClusterOperationsStub(tb testing.TB) *ClusterOperationsStub {
	return &ClusterOperationsStub{
		tb: tb,
		DoList: func(opts *types.ListOpts) (*managementClient.ClusterCollection, error) {
			assert.FailInStub(tb, 2, "Unexpected call of List")
			return nil, nil
		},
		DoCreate: func(opts *managementClient.Cluster) (*managementClient.Cluster, error) {
			assert.FailInStub(tb, 2, "Unexpected call of Create")
			return nil, nil
		},
	}
}

// ClusterOperationsStub structure to hold callbacks used to stub github.com/rancher/types/client/management/v3/ClusterOperations
type ClusterOperationsStub struct {
	tb       testing.TB
	DoList   func(opts *types.ListOpts) (*managementClient.ClusterCollection, error)
	DoCreate func(opts *managementClient.Cluster) (*managementClient.Cluster, error)
}

// List implements github.com/rancher/types/client/management/v3/ClusterOperations.List(...)
func (stub ClusterOperationsStub) List(opts *types.ListOpts) (*managementClient.ClusterCollection, error) {
	return stub.DoList(opts)

}

// Create implements github.com/rancher/types/client/management/v3/ClusterOperations.Create(...)
func (stub ClusterOperationsStub) Create(opts *managementClient.Cluster) (*managementClient.Cluster, error) {
	return stub.DoCreate(opts)

}

// Update implements github.com/rancher/types/client/management/v3/ClusterOperations.Update(...)
func (stub ClusterOperationsStub) Update(existing *managementClient.Cluster, updates interface{}) (*managementClient.Cluster, error) {
	assert.FailInStub(stub.tb, 2, "Unexpected call of Update")
	return nil, nil

}

// Replace implements github.com/rancher/types/client/management/v3/ClusterOperations.Replace(...)
func (stub ClusterOperationsStub) Replace(existing *managementClient.Cluster) (*managementClient.Cluster, error) {
	assert.FailInStub(stub.tb, 2, "Unexpected call of Replace")
	return nil, nil

}

// ByID implements github.com/rancher/types/client/management/v3/ClusterOperations.ByID(...)
func (stub ClusterOperationsStub) ByID(id string) (*managementClient.Cluster, error) {
	assert.FailInStub(stub.tb, 2, "Unexpected call of ByID")
	return nil, nil

}

// Delete implements github.com/rancher/types/client/management/v3/ClusterOperations.Delete(...)
func (stub ClusterOperationsStub) Delete(container *managementClient.Cluster) error {
	assert.FailInStub(stub.tb, 2, "Unexpected call of Delete")
	return nil

}

// ActionBackupEtcd implements github.com/rancher/types/client/management/v3/ClusterOperations.ActionBackupEtcd(...)
func (stub ClusterOperationsStub) ActionBackupEtcd(resource *managementClient.Cluster) error {
	assert.FailInStub(stub.tb, 2, "Unexpected call of ActionBackupEtcd")
	return nil

}

// ActionDisableMonitoring implements github.com/rancher/types/client/management/v3/ClusterOperations.ActionDisableMonitoring(...)
func (stub ClusterOperationsStub) ActionDisableMonitoring(resource *managementClient.Cluster) error {
	assert.FailInStub(stub.tb, 2, "Unexpected call of ActionDisableMonitoring")
	return nil

}

// ActionEditMonitoring implements github.com/rancher/types/client/management/v3/ClusterOperations.ActionEditMonitoring(...)
func (stub ClusterOperationsStub) ActionEditMonitoring(resource *managementClient.Cluster, input *managementClient.MonitoringInput) error {
	assert.FailInStub(stub.tb, 2, "Unexpected call of ActionEnableMonitoring")
	return nil
}

// ActionEnableMonitoring implements github.com/rancher/types/client/management/v3/ClusterOperations.ActionEnableMonitoring(...)
func (stub ClusterOperationsStub) ActionEnableMonitoring(resource *managementClient.Cluster, input *managementClient.MonitoringInput) error {
	assert.FailInStub(stub.tb, 2, "Unexpected call of ActionEnableMonitoring")
	return nil
}

// ActionViewMonitoring implements github.com/rancher/types/client/management/v3/ClusterOperations.ActionViewMonitoring(...)
func (stub ClusterOperationsStub) ActionViewMonitoring(resource *managementClient.Cluster) (*managementClient.MonitoringOutput, error) {
	assert.FailInStub(stub.tb, 2, "Unexpected call of ActionViewMonitoring")
	return nil, nil
}

// ActionExportYaml implements github.com/rancher/types/client/management/v3/ClusterOperations.ActionExportYaml(...)
func (stub ClusterOperationsStub) ActionExportYaml(resource *managementClient.Cluster) (*managementClient.ExportOutput, error) {
	assert.FailInStub(stub.tb, 2, "Unexpected call of ActionExportYaml")
	return nil, nil

}

// ActionGenerateKubeconfig implements github.com/rancher/types/client/management/v3/ClusterOperations.ActionGenerateKubeconfig(...)
func (stub ClusterOperationsStub) ActionGenerateKubeconfig(resource *managementClient.Cluster) (*managementClient.GenerateKubeConfigOutput, error) {
	assert.FailInStub(stub.tb, 2, "Unexpected call of ActionGenerateKubeconfig")
	return nil, nil

}

// ActionImportYaml implements github.com/rancher/types/client/management/v3/ClusterOperations.ActionImportYaml(...)
func (stub ClusterOperationsStub) ActionImportYaml(resource *managementClient.Cluster, input *managementClient.ImportClusterYamlInput) (*managementClient.ImportYamlOutput, error) {
	assert.FailInStub(stub.tb, 2, "Unexpected call of ActionImportYaml")
	return nil, nil

}

// ActionRestoreFromEtcdBackup implements github.com/rancher/types/client/management/v3/ClusterOperations.ActionRestoreFromEtcdBackup(...)
func (stub ClusterOperationsStub) ActionRestoreFromEtcdBackup(resource *managementClient.Cluster, input *managementClient.RestoreFromEtcdBackupInput) error {
	assert.FailInStub(stub.tb, 2, "Unexpected call of ActionRestoreFromEtcdBackup")
	return nil

}

// ActionRotateCertificates implements github.com/rancher/types/client/management/v3/ClusterOperations.ActionRotateCertificates(...)
func (stub ClusterOperationsStub) ActionRotateCertificates(resource *managementClient.Cluster, input *managementClient.RotateCertificateInput) (*managementClient.RotateCertificateOutput, error) {
	assert.FailInStub(stub.tb, 2, "Unexpected call of ActionRotateCertificates")
	return nil, nil

}

// ActionRunSecurityScan implements github.com/rancher/types/client/management/v3/ClusterOperations.ActionRunSecurityScan(...)
func (stub ClusterOperationsStub) ActionRunSecurityScan(resource *managementClient.Cluster) error {
	assert.FailInStub(stub.tb, 2, "Unexpected call of ActionRunSecurityScan")
	return nil

}

// ActionSaveAsTemplate implements github.com/rancher/types/client/management/v3/ClusterOperations.ActionSaveAsTemplate(...)
func (stub ClusterOperationsStub) ActionSaveAsTemplate(resource *managementClient.Cluster, input *managementClient.SaveAsTemplateInput) (*managementClient.SaveAsTemplateOutput, error) {
	assert.FailInStub(stub.tb, 2, "Unexpected call of ActionSaveAsTemplate")
	return nil, nil

}
