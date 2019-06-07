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

// CreateAppOperationsStub creates a stub of github.com/rancher/types/client/project/v3/AppOperations
func CreateAppOperationsStub(tb testing.TB) *AppOperationsStub {
	return &AppOperationsStub{
		tb: tb,
		DoList: func(opts *types.ListOpts) (*projectClient.AppCollection, error) {
			assert.FailInStub(tb, 2, "Unexpected call of List")
			return nil, nil
		},
		DoCreate: func(opts *projectClient.App) (*projectClient.App, error) {
			assert.FailInStub(tb, 2, "Unexpected call of Create")
			return nil, nil
		},
		DoActionUpgrade: func(resource *projectClient.App, input *projectClient.AppUpgradeConfig) error {
			assert.FailInStub(tb, 2, "Unexpected call of ActionUpgrade")
			return nil
		},
	}
}

// AppOperationsStub structure to hold callbacks used to stub github.com/rancher/types/client/project/v3/AppOperations
type AppOperationsStub struct {
	tb              testing.TB
	DoList          func(opts *types.ListOpts) (*projectClient.AppCollection, error)
	DoCreate        func(opts *projectClient.App) (*projectClient.App, error)
	DoActionUpgrade func(resource *projectClient.App, input *projectClient.AppUpgradeConfig) error
}

// List implements github.com/rancher/types/client/project/v3/AppOperations.List(...)
func (stub AppOperationsStub) List(opts *types.ListOpts) (*projectClient.AppCollection, error) {
	return stub.DoList(opts)
}

// Create implements github.com/rancher/types/client/project/v3/AppOperations.Create(...)
func (stub AppOperationsStub) Create(opts *projectClient.App) (*projectClient.App, error) {
	return stub.DoCreate(opts)
}

// Update implements github.com/rancher/types/client/project/v3/AppOperations.Update(...)
func (stub AppOperationsStub) Update(existing *projectClient.App, updates interface{}) (*projectClient.App, error) {
	assert.FailInStub(stub.tb, 2, "Unexpected call of List")
	return nil, nil
}

// Replace implements github.com/rancher/types/client/project/v3/AppOperations.Replace(...)
func (stub AppOperationsStub) Replace(existing *projectClient.App) (*projectClient.App, error) {
	assert.FailInStub(stub.tb, 2, "Unexpected call of List")
	return nil, nil
}

// ByID implements github.com/rancher/types/client/project/v3/AppOperations.ByID(...)
func (stub AppOperationsStub) ByID(id string) (*projectClient.App, error) {
	assert.FailInStub(stub.tb, 2, "Unexpected call of List")
	return nil, nil
}

// Delete implements github.com/rancher/types/client/project/v3/AppOperations.Delete(...)
func (stub AppOperationsStub) Delete(container *projectClient.App) error {
	assert.FailInStub(stub.tb, 2, "Unexpected call of List")
	return nil
}

// ActionRollback implements github.com/rancher/types/client/project/v3/AppOperations.ActionRollback(...)
func (stub AppOperationsStub) ActionRollback(resource *projectClient.App, input *projectClient.RollbackRevision) error {
	assert.FailInStub(stub.tb, 2, "Unexpected call of List")
	return nil
}

// ActionUpgrade implements github.com/rancher/types/client/project/v3/AppOperations.ActionUpgrade(...)
func (stub AppOperationsStub) ActionUpgrade(resource *projectClient.App, input *projectClient.AppUpgradeConfig) error {
	return stub.DoActionUpgrade(resource, input)
}
