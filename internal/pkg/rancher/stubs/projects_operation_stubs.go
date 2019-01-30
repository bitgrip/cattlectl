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

// CreateProjectOperationsStub creates a stub of github.com/rancher/types/client/management/v3/ProjectOperations
func CreateProjectOperationsStub(tb testing.TB) *ProjectOperationsStub {
	return &ProjectOperationsStub{
		tb: tb,
		DoList: func(opts *types.ListOpts) (*managementClient.ProjectCollection, error) {
			assert.FailInStub(tb, 2, "Unexpected call of List")
			return nil, nil
		},
		DoCreate: func(opts *managementClient.Project) (*managementClient.Project, error) {
			assert.FailInStub(tb, 2, "Unexpected call of Create")
			return nil, nil
		},
	}
}

// ProjectOperationsStub structure to hold callbacks used to stub github.com/rancher/types/client/management/v3/ProjectOperations
type ProjectOperationsStub struct {
	tb       testing.TB
	DoList   func(opts *types.ListOpts) (*managementClient.ProjectCollection, error)
	DoCreate func(opts *managementClient.Project) (*managementClient.Project, error)
}

// List implements github.com/rancher/types/client/management/v3/ProjectOperations.List(...)
func (stub ProjectOperationsStub) List(opts *types.ListOpts) (*managementClient.ProjectCollection, error) {
	return stub.DoList(opts)

}

// Create implements github.com/rancher/types/client/management/v3/ProjectOperations.Create(...)
func (stub ProjectOperationsStub) Create(opts *managementClient.Project) (*managementClient.Project, error) {
	return stub.DoCreate(opts)

}

// Update implements github.com/rancher/types/client/management/v3/ProjectOperations.Update(...)
func (stub ProjectOperationsStub) Update(existing *managementClient.Project, updates interface{}) (*managementClient.Project, error) {
	assert.FailInStub(stub.tb, 2, "Unexpected call of Update")
	return nil, nil

}

// Replace implements github.com/rancher/types/client/management/v3/ProjectOperations.Replace(...)
func (stub ProjectOperationsStub) Replace(existing *managementClient.Project) (*managementClient.Project, error) {
	assert.FailInStub(stub.tb, 2, "Unexpected call of Replace")
	return nil, nil

}

// ByID implements github.com/rancher/types/client/management/v3/ProjectOperations.ByID(...)
func (stub ProjectOperationsStub) ByID(id string) (*managementClient.Project, error) {
	assert.FailInStub(stub.tb, 2, "Unexpected call of ByID")
	return nil, nil

}

// Delete implements github.com/rancher/types/client/management/v3/ProjectOperations.Delete(...)
func (stub ProjectOperationsStub) Delete(container *managementClient.Project) error {
	assert.FailInStub(stub.tb, 2, "Unexpected call of Delete")
	return nil

}

// ActionDisableMonitoring implements github.com/rancher/types/client/management/v3/ProjectOperations.ActionDisableMonitoring(...)
func (stub ProjectOperationsStub) ActionDisableMonitoring(resource *managementClient.Project) error {
	assert.FailInStub(stub.tb, 2, "Unexpected call of ActionDisableMonitoring")
	return nil

}

// ActionEnableMonitoring implements github.com/rancher/types/client/management/v3/ProjectOperations.ActionEnableMonitoring(...)
func (stub ProjectOperationsStub) ActionEnableMonitoring(resource *managementClient.Project, input *managementClient.MonitoringInput) error {
	assert.FailInStub(stub.tb, 2, "Unexpected call of ActionEnableMonitoring")
	return nil

}

// ActionExportYaml implements github.com/rancher/types/client/management/v3/ProjectOperations.ActionExportYaml(...)
func (stub ProjectOperationsStub) ActionExportYaml(resource *managementClient.Project) error {
	assert.FailInStub(stub.tb, 2, "Unexpected call of ActionExportYaml")
	return nil

}

// ActionSetpodsecuritypolicytemplate implements github.com/rancher/types/client/management/v3/ProjectOperations.ActionSetpodsecuritypolicytemplate(...)
func (stub ProjectOperationsStub) ActionSetpodsecuritypolicytemplate(resource *managementClient.Project, input *managementClient.SetPodSecurityPolicyTemplateInput) (*managementClient.Project, error) {
	assert.FailInStub(stub.tb, 2, "Unexpected call of ActionSetpodsecuritypolicytemplate")
	return nil, nil
}
