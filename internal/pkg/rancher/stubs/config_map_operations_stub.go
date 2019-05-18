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

// CreateConfigMapOperationsStub creates a stub of github.com/rancher/types/client/project/v3/ConfigMapOperations
func CreateConfigMapOperationsStub(tb testing.TB) *ConfigMapOperationsStub {
	return &ConfigMapOperationsStub{
		tb: tb,
		DoList: func(opts *types.ListOpts) (*projectClient.ConfigMapCollection, error) {
			assert.FailInStub(tb, 2, "Unexpected call of List")
			return nil, nil
		},
		DoCreate: func(opts *projectClient.ConfigMap) (*projectClient.ConfigMap, error) {
			assert.FailInStub(tb, 2, "Unexpected call of Create")
			return nil, nil
		},
		DoUpdate: func(existing *projectClient.ConfigMap, updates interface{}) (*projectClient.ConfigMap, error) {
			assert.FailInStub(tb, 2, "Unexpected call of Update")
			return nil, nil
		},
		DoReplace: func(existing *projectClient.ConfigMap) (*projectClient.ConfigMap, error) {
			assert.FailInStub(tb, 2, "Unexpected call of Replace")
			return nil, nil
		},
		DoByID: func(id string) (*projectClient.ConfigMap, error) {
			assert.FailInStub(tb, 2, "Unexpected call of ByID")
			return nil, nil
		},
		DoDelete: func(container *projectClient.ConfigMap) error {
			assert.FailInStub(tb, 2, "Unexpected call of Delete")
			return nil
		},
	}
}

// ConfigMapOperationsStub structure to hold callbacks used to stub github.com/rancher/types/client/project/v3/ConfigMapOperations
type ConfigMapOperationsStub struct {
	tb        testing.TB
	DoList    func(opts *types.ListOpts) (*projectClient.ConfigMapCollection, error)
	DoCreate  func(opts *projectClient.ConfigMap) (*projectClient.ConfigMap, error)
	DoUpdate  func(existing *projectClient.ConfigMap, updates interface{}) (*projectClient.ConfigMap, error)
	DoReplace func(existing *projectClient.ConfigMap) (*projectClient.ConfigMap, error)
	DoByID    func(id string) (*projectClient.ConfigMap, error)
	DoDelete  func(container *projectClient.ConfigMap) error
}

// List implements github.com/rancher/types/client/project/v3/ConfigMapOperations.List(...)
func (stub ConfigMapOperationsStub) List(opts *types.ListOpts) (*projectClient.ConfigMapCollection, error) {
	return stub.DoList(opts)
}

// Create implements github.com/rancher/types/client/project/v3/ConfigMapOperations.Create(...)
func (stub ConfigMapOperationsStub) Create(opts *projectClient.ConfigMap) (*projectClient.ConfigMap, error) {
	return stub.DoCreate(opts)
}

// Update implements github.com/rancher/types/client/project/v3/ConfigMapOperations.Update(...)
func (stub ConfigMapOperationsStub) Update(existing *projectClient.ConfigMap, updates interface{}) (*projectClient.ConfigMap, error) {
	return stub.DoUpdate(existing, updates)
}

// Replace implements github.com/rancher/types/client/project/v3/ConfigMapOperations.Replace(...)
func (stub ConfigMapOperationsStub) Replace(existing *projectClient.ConfigMap) (*projectClient.ConfigMap, error) {
	return stub.DoReplace(existing)
}

// ByID implements github.com/rancher/types/client/project/v3/ConfigMapOperations.ByID(...)
func (stub ConfigMapOperationsStub) ByID(id string) (*projectClient.ConfigMap, error) {
	return stub.DoByID(id)
}

// Delete implements github.com/rancher/types/client/project/v3/ConfigMapOperations.Delete(...)
func (stub ConfigMapOperationsStub) Delete(container *projectClient.ConfigMap) error {
	return stub.DoDelete(container)
}
