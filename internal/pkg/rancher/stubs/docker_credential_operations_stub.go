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

// CreateDockerCredentialOperationsStub creates a stub of github.com/rancher/types/client/project/v3/DockerCredentialOperations
func CreateDockerCredentialOperationsStub(tb testing.TB) *DockerCredentialOperationsStub {
	return &DockerCredentialOperationsStub{
		tb: tb,
		DoList: func(opts *types.ListOpts) (*projectClient.DockerCredentialCollection, error) {
			assert.FailInStub(tb, 2, "Unexpected call of List")
			return nil, nil
		},
		DoCreate: func(opts *projectClient.DockerCredential) (*projectClient.DockerCredential, error) {
			assert.FailInStub(tb, 2, "Unexpected call of Create")
			return nil, nil
		},
		DoUpdate: func(existing *projectClient.DockerCredential, updates interface{}) (*projectClient.DockerCredential, error) {
			assert.FailInStub(tb, 2, "Unexpected call of Update")
			return nil, nil
		},
		DoReplace: func(existing *projectClient.DockerCredential) (*projectClient.DockerCredential, error) {
			assert.FailInStub(tb, 2, "Unexpected call of Replace")
			return nil, nil
		},
		DoByID: func(id string) (*projectClient.DockerCredential, error) {
			assert.FailInStub(tb, 2, "Unexpected call of ByID")
			return nil, nil
		},
		DoDelete: func(container *projectClient.DockerCredential) error {
			assert.FailInStub(tb, 2, "Unexpected call of Delete")
			return nil
		},
	}
}

// DockerCredentialOperationsStub structure to hold callbacks used to stub github.com/rancher/types/client/project/v3/DockerCredentialOperations
type DockerCredentialOperationsStub struct {
	tb        testing.TB
	DoList    func(opts *types.ListOpts) (*projectClient.DockerCredentialCollection, error)
	DoCreate  func(opts *projectClient.DockerCredential) (*projectClient.DockerCredential, error)
	DoUpdate  func(existing *projectClient.DockerCredential, updates interface{}) (*projectClient.DockerCredential, error)
	DoReplace func(existing *projectClient.DockerCredential) (*projectClient.DockerCredential, error)
	DoByID    func(id string) (*projectClient.DockerCredential, error)
	DoDelete  func(container *projectClient.DockerCredential) error
}

// List implements github.com/rancher/types/client/project/v3/DockerCredentialOperations.List(...)
func (stub DockerCredentialOperationsStub) List(opts *types.ListOpts) (*projectClient.DockerCredentialCollection, error) {
	return stub.DoList(opts)
}

// Create implements github.com/rancher/types/client/project/v3/DockerCredentialOperations.Create(...)
func (stub DockerCredentialOperationsStub) Create(opts *projectClient.DockerCredential) (*projectClient.DockerCredential, error) {
	return stub.DoCreate(opts)
}

// Update implements github.com/rancher/types/client/project/v3/DockerCredentialOperations.Update(...)
func (stub DockerCredentialOperationsStub) Update(existing *projectClient.DockerCredential, updates interface{}) (*projectClient.DockerCredential, error) {
	return stub.DoUpdate(existing, updates)
}

// Replace implements github.com/rancher/types/client/project/v3/DockerCredentialOperations.Replace(...)
func (stub DockerCredentialOperationsStub) Replace(existing *projectClient.DockerCredential) (*projectClient.DockerCredential, error) {
	return stub.DoReplace(existing)
}

// ByID implements github.com/rancher/types/client/project/v3/DockerCredentialOperations.ByID(...)
func (stub DockerCredentialOperationsStub) ByID(id string) (*projectClient.DockerCredential, error) {
	return stub.DoByID(id)
}

// Delete implements github.com/rancher/types/client/project/v3/DockerCredentialOperations.Delete(...)
func (stub DockerCredentialOperationsStub) Delete(container *projectClient.DockerCredential) error {
	return stub.DoDelete(container)
}
