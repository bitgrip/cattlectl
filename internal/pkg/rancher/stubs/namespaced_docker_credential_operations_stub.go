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

// CreateNamespacedDockerCredentialOperationsStub creates a stub of github.com/rancher/types/client/project/v3/NamespacedDockerCredentialOperations
func CreateNamespacedDockerCredentialOperationsStub(tb testing.TB) *NamespacedDockerCredentialOperationsStub {
	return &NamespacedDockerCredentialOperationsStub{
		tb: tb,
		DoList: func(opts *types.ListOpts) (*projectClient.NamespacedDockerCredentialCollection, error) {
			assert.FailInStub(tb, 2, "Unexpected call of List")
			return nil, nil
		},
		DoCreate: func(opts *projectClient.NamespacedDockerCredential) (*projectClient.NamespacedDockerCredential, error) {
			assert.FailInStub(tb, 2, "Unexpected call of Create")
			return nil, nil
		},
		DoUpdate: func(existing *projectClient.NamespacedDockerCredential, updates interface{}) (*projectClient.NamespacedDockerCredential, error) {
			assert.FailInStub(tb, 2, "Unexpected call of Update")
			return nil, nil
		},
		DoReplace: func(existing *projectClient.NamespacedDockerCredential) (*projectClient.NamespacedDockerCredential, error) {
			assert.FailInStub(tb, 2, "Unexpected call of Replace")
			return nil, nil
		},
		DoByID: func(id string) (*projectClient.NamespacedDockerCredential, error) {
			assert.FailInStub(tb, 2, "Unexpected call of ByID")
			return nil, nil
		},
		DoDelete: func(container *projectClient.NamespacedDockerCredential) error {
			assert.FailInStub(tb, 2, "Unexpected call of Delete")
			return nil
		},
	}
}

// NamespacedDockerCredentialOperationsStub structure to hold callbacks used to stub github.com/rancher/types/client/project/v3/DockerCredentialOperations
type NamespacedDockerCredentialOperationsStub struct {
	tb        testing.TB
	DoList    func(opts *types.ListOpts) (*projectClient.NamespacedDockerCredentialCollection, error)
	DoCreate  func(opts *projectClient.NamespacedDockerCredential) (*projectClient.NamespacedDockerCredential, error)
	DoUpdate  func(existing *projectClient.NamespacedDockerCredential, updates interface{}) (*projectClient.NamespacedDockerCredential, error)
	DoReplace func(existing *projectClient.NamespacedDockerCredential) (*projectClient.NamespacedDockerCredential, error)
	DoByID    func(id string) (*projectClient.NamespacedDockerCredential, error)
	DoDelete  func(container *projectClient.NamespacedDockerCredential) error
}

// List implements github.com/rancher/types/client/project/v3/DockerCredentialOperations.List(...)
func (stub NamespacedDockerCredentialOperationsStub) List(opts *types.ListOpts) (*projectClient.NamespacedDockerCredentialCollection, error) {
	return stub.DoList(opts)
}

// Create implements github.com/rancher/types/client/project/v3/DockerCredentialOperations.Create(...)
func (stub NamespacedDockerCredentialOperationsStub) Create(opts *projectClient.NamespacedDockerCredential) (*projectClient.NamespacedDockerCredential, error) {
	return stub.DoCreate(opts)
}

// Update implements github.com/rancher/types/client/project/v3/DockerCredentialOperations.Update(...)
func (stub NamespacedDockerCredentialOperationsStub) Update(existing *projectClient.NamespacedDockerCredential, updates interface{}) (*projectClient.NamespacedDockerCredential, error) {
	return stub.DoUpdate(existing, updates)
}

// Replace implements github.com/rancher/types/client/project/v3/DockerCredentialOperations.Replace(...)
func (stub NamespacedDockerCredentialOperationsStub) Replace(existing *projectClient.NamespacedDockerCredential) (*projectClient.NamespacedDockerCredential, error) {
	return stub.DoReplace(existing)
}

// ByID implements github.com/rancher/types/client/project/v3/DockerCredentialOperations.ByID(...)
func (stub NamespacedDockerCredentialOperationsStub) ByID(id string) (*projectClient.NamespacedDockerCredential, error) {
	return stub.DoByID(id)
}

// Delete implements github.com/rancher/types/client/project/v3/DockerCredentialOperations.Delete(...)
func (stub NamespacedDockerCredentialOperationsStub) Delete(container *projectClient.NamespacedDockerCredential) error {
	return stub.DoDelete(container)
}
