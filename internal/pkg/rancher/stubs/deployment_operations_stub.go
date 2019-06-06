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
	projectClient "github.com/rancher/types/client/project/v3"
)

// CreateDeploymentOperationsStub creates a stub of github.com/rancher/types/client/project/v3/DeploymentOperations
func CreateDeploymentOperationsStub(tb testing.TB) *DeploymentOperationsStub {
	return &DeploymentOperationsStub{
		tb: tb,
		DoList: func(opts *types.ListOpts) (*projectClient.DeploymentCollection, error) {
			assert.FailInStub(tb, 2, "Unexpected call of List")
			return nil, nil
		},
		DoCreate: func(opts *projectClient.Deployment) (*projectClient.Deployment, error) {
			assert.FailInStub(tb, 2, "Unexpected call of Create")
			return nil, nil
		},
		DoUpdate: func(existing *projectClient.Deployment, updates interface{}) (*projectClient.Deployment, error) {
			assert.FailInStub(tb, 2, "Unexpected call of List")
			return nil, nil
		},
		DoReplace: func(existing *projectClient.Deployment) (*projectClient.Deployment, error) {
			assert.FailInStub(tb, 2, "Unexpected call of List")
			return nil, nil
		},
		DoActionPause: func(resource *projectClient.Deployment) error {
			assert.FailInStub(tb, 2, "Unexpected call of List")
			return nil
		},
		DoActionResume: func(resource *projectClient.Deployment) error {
			assert.FailInStub(tb, 2, "Unexpected call of List")
			return nil
		},
		DoActionRollback: func(resource *projectClient.Deployment, input *projectClient.DeploymentRollbackInput) error {
			assert.FailInStub(tb, 2, "Unexpected call of List")
			return nil
		},
	}
}

// DeploymentOperationsStub structure to hold callbacks used to stub github.com/rancher/types/client/project/v3/DeploymentOperations
type DeploymentOperationsStub struct {
	tb               testing.TB
	DoList           func(opts *types.ListOpts) (*projectClient.DeploymentCollection, error)
	DoCreate         func(opts *projectClient.Deployment) (*projectClient.Deployment, error)
	DoUpdate         func(existing *projectClient.Deployment, updates interface{}) (*projectClient.Deployment, error)
	DoReplace        func(existing *projectClient.Deployment) (*projectClient.Deployment, error)
	DoActionPause    func(resource *projectClient.Deployment) error
	DoActionResume   func(resource *projectClient.Deployment) error
	DoActionRollback func(resource *projectClient.Deployment, input *projectClient.DeploymentRollbackInput) error
}

// List implements github.com/rancher/types/client/project/v3/DeploymentOperations.List(...)
func (stub DeploymentOperationsStub) List(opts *types.ListOpts) (*projectClient.DeploymentCollection, error) {
	return stub.DoList(opts)
}

// Create implements github.com/rancher/types/client/project/v3/DeploymentOperations.Create(...)
func (stub DeploymentOperationsStub) Create(opts *projectClient.Deployment) (*projectClient.Deployment, error) {
	return stub.DoCreate(opts)
}

// Update implements github.com/rancher/types/client/project/v3/DeploymentOperations.Update(...)
func (stub DeploymentOperationsStub) Update(existing *projectClient.Deployment, updates interface{}) (*projectClient.Deployment, error) {
	return stub.DoUpdate(existing, updates)
}

// Replace implements github.com/rancher/types/client/project/v3/DeploymentOperations.Replace(...)
func (stub DeploymentOperationsStub) Replace(existing *projectClient.Deployment) (*projectClient.Deployment, error) {
	return stub.DoReplace(existing)
}

// ByID implements github.com/rancher/types/client/project/v3/DeploymentOperations.ByID(...)
func (stub DeploymentOperationsStub) ByID(id string) (*projectClient.Deployment, error) {
	assert.FailInStub(stub.tb, 2, "Unexpected call of List")
	return nil, nil
}

// Delete implements github.com/rancher/types/client/project/v3/DeploymentOperations.Delete(...)
func (stub DeploymentOperationsStub) Delete(container *projectClient.Deployment) error {
	assert.FailInStub(stub.tb, 2, "Unexpected call of List")
	return nil
}

// ActionPause implements github.com/rancher/types/client/project/v3/DeploymentOperations.ActionPause(...)
func (stub DeploymentOperationsStub) ActionPause(resource *projectClient.Deployment) error {
	return stub.DoActionPause(resource)
}

// ActionResume implements github.com/rancher/types/client/project/v3/DeploymentOperations.ActionResume(...)
func (stub DeploymentOperationsStub) ActionResume(resource *projectClient.Deployment) error {
	return stub.DoActionResume(resource)
}

// ActionRollback implements github.com/rancher/types/client/project/v3/DeploymentOperations.ActionRollback(...)
func (stub DeploymentOperationsStub) ActionRollback(resource *projectClient.Deployment, input *projectClient.DeploymentRollbackInput) error {
	return stub.DoActionRollback(resource, input)
}
