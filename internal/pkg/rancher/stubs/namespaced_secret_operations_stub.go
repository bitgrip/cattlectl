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

// CreateNamespacedSecretOperationsStub creates a stub of github.com/rancher/types/client/project/v3/NamespacedSecretOperations
func CreateNamespacedSecretOperationsStub(tb testing.TB) *NamespacedSecretOperationsStub {
	return &NamespacedSecretOperationsStub{
		tb: tb,
		DoList: func(opts *types.ListOpts) (*projectClient.NamespacedSecretCollection, error) {
			assert.FailInStub(tb, 2, "Unexpected call of List")
			return nil, nil
		},
		DoCreate: func(opts *projectClient.NamespacedSecret) (*projectClient.NamespacedSecret, error) {
			assert.FailInStub(tb, 2, "Unexpected call of Create")
			return nil, nil
		},
		DoUpdate: func(existing *projectClient.NamespacedSecret, updates interface{}) (*projectClient.NamespacedSecret, error) {
			assert.FailInStub(tb, 2, "Unexpected call of Update")
			return nil, nil
		},
		DoReplace: func(existing *projectClient.NamespacedSecret) (*projectClient.NamespacedSecret, error) {
			assert.FailInStub(tb, 2, "Unexpected call of Replace")
			return nil, nil
		},
		DoByID: func(id string) (*projectClient.NamespacedSecret, error) {
			assert.FailInStub(tb, 2, "Unexpected call of ByID")
			return nil, nil
		},
		DoDelete: func(container *projectClient.NamespacedSecret) error {
			assert.FailInStub(tb, 2, "Unexpected call of Delete")
			return nil
		},
	}
}

// NamespacedSecretOperationsStub structure to hold callbacks used to stub github.com/rancher/types/client/project/v3/NamespacedSecretOperations
type NamespacedSecretOperationsStub struct {
	tb        testing.TB
	DoList    func(opts *types.ListOpts) (*projectClient.NamespacedSecretCollection, error)
	DoCreate  func(opts *projectClient.NamespacedSecret) (*projectClient.NamespacedSecret, error)
	DoUpdate  func(existing *projectClient.NamespacedSecret, updates interface{}) (*projectClient.NamespacedSecret, error)
	DoReplace func(existing *projectClient.NamespacedSecret) (*projectClient.NamespacedSecret, error)
	DoByID    func(id string) (*projectClient.NamespacedSecret, error)
	DoDelete  func(container *projectClient.NamespacedSecret) error
}

// List implements github.com/rancher/types/client/project/v3/NamespacedSecretOperations.List(...)
func (stub NamespacedSecretOperationsStub) List(opts *types.ListOpts) (*projectClient.NamespacedSecretCollection, error) {
	return stub.DoList(opts)
}

// Create implements github.com/rancher/types/client/project/v3/NamespacedSecretOperations.Create(...)
func (stub NamespacedSecretOperationsStub) Create(opts *projectClient.NamespacedSecret) (*projectClient.NamespacedSecret, error) {
	return stub.DoCreate(opts)
}

// Update implements github.com/rancher/types/client/project/v3/NamespacedSecretOperations.Update(...)
func (stub NamespacedSecretOperationsStub) Update(existing *projectClient.NamespacedSecret, updates interface{}) (*projectClient.NamespacedSecret, error) {
	return stub.DoUpdate(existing, updates)
}

// Replace implements github.com/rancher/types/client/project/v3/NamespacedSecretOperations.Replace(...)
func (stub NamespacedSecretOperationsStub) Replace(existing *projectClient.NamespacedSecret) (*projectClient.NamespacedSecret, error) {
	return stub.DoReplace(existing)
}

// ByID implements github.com/rancher/types/client/project/v3/NamespacedSecretOperations.ByID(...)
func (stub NamespacedSecretOperationsStub) ByID(id string) (*projectClient.NamespacedSecret, error) {
	return stub.DoByID(id)
}

// Delete implements github.com/rancher/types/client/project/v3/NamespacedSecretOperations.Delete(...)
func (stub NamespacedSecretOperationsStub) Delete(container *projectClient.NamespacedSecret) error {
	return stub.DoDelete(container)
}
