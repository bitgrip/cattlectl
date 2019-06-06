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

// CreateSecretOperationsStub creates a stub of github.com/rancher/types/client/project/v3/SecretOperations
func CreateSecretOperationsStub(tb testing.TB) *SecretOperationsStub {
	return &SecretOperationsStub{
		tb: tb,
		DoList: func(opts *types.ListOpts) (*projectClient.SecretCollection, error) {
			assert.FailInStub(tb, 2, "Unexpected call of List")
			return nil, nil
		},
		DoCreate: func(opts *projectClient.Secret) (*projectClient.Secret, error) {
			assert.FailInStub(tb, 2, "Unexpected call of Create")
			return nil, nil
		},
		DoUpdate: func(existing *projectClient.Secret, updates interface{}) (*projectClient.Secret, error) {
			assert.FailInStub(tb, 2, "Unexpected call of Update")
			return nil, nil
		},
		DoReplace: func(existing *projectClient.Secret) (*projectClient.Secret, error) {
			assert.FailInStub(tb, 2, "Unexpected call of Replace")
			return nil, nil
		},
		DoByID: func(id string) (*projectClient.Secret, error) {
			assert.FailInStub(tb, 2, "Unexpected call of ByID")
			return nil, nil
		},
		DoDelete: func(container *projectClient.Secret) error {
			assert.FailInStub(tb, 2, "Unexpected call of Delete")
			return nil
		},
	}
}

// SecretOperationsStub structure to hold callbacks used to stub github.com/rancher/types/client/project/v3/SecretOperations
type SecretOperationsStub struct {
	tb        testing.TB
	DoList    func(opts *types.ListOpts) (*projectClient.SecretCollection, error)
	DoCreate  func(opts *projectClient.Secret) (*projectClient.Secret, error)
	DoUpdate  func(existing *projectClient.Secret, updates interface{}) (*projectClient.Secret, error)
	DoReplace func(existing *projectClient.Secret) (*projectClient.Secret, error)
	DoByID    func(id string) (*projectClient.Secret, error)
	DoDelete  func(container *projectClient.Secret) error
}

// List implements github.com/rancher/types/client/project/v3/SecretOperations.List(...)
func (stub SecretOperationsStub) List(opts *types.ListOpts) (*projectClient.SecretCollection, error) {
	return stub.DoList(opts)
}

// Create implements github.com/rancher/types/client/project/v3/SecretOperations.Create(...)
func (stub SecretOperationsStub) Create(opts *projectClient.Secret) (*projectClient.Secret, error) {
	return stub.DoCreate(opts)
}

// Update implements github.com/rancher/types/client/project/v3/SecretOperations.Update(...)
func (stub SecretOperationsStub) Update(existing *projectClient.Secret, updates interface{}) (*projectClient.Secret, error) {
	return stub.DoUpdate(existing, updates)
}

// Replace implements github.com/rancher/types/client/project/v3/SecretOperations.Replace(...)
func (stub SecretOperationsStub) Replace(existing *projectClient.Secret) (*projectClient.Secret, error) {
	return stub.DoReplace(existing)
}

// ByID implements github.com/rancher/types/client/project/v3/SecretOperations.ByID(...)
func (stub SecretOperationsStub) ByID(id string) (*projectClient.Secret, error) {
	return stub.DoByID(id)
}

// Delete implements github.com/rancher/types/client/project/v3/SecretOperations.Delete(...)
func (stub SecretOperationsStub) Delete(container *projectClient.Secret) error {
	return stub.DoDelete(container)
}
