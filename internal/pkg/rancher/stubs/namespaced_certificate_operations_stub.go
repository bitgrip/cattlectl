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

// CreateNamespacedCertificateOperationsStub creates a stub of github.com/rancher/types/client/project/v3/NamespacedCertificateOperations
func CreateNamespacedCertificateOperationsStub(tb testing.TB) *NamespacedCertificateOperationsStub {
	return &NamespacedCertificateOperationsStub{
		tb: tb,
		DoList: func(opts *types.ListOpts) (*projectClient.NamespacedCertificateCollection, error) {
			assert.FailInStub(tb, 2, "Unexpected call of List")
			return nil, nil
		},
		DoCreate: func(opts *projectClient.NamespacedCertificate) (*projectClient.NamespacedCertificate, error) {
			assert.FailInStub(tb, 2, "Unexpected call of Create")
			return nil, nil
		},
		DoUpdate: func(existing *projectClient.NamespacedCertificate, updates interface{}) (*projectClient.NamespacedCertificate, error) {
			assert.FailInStub(tb, 2, "Unexpected call of Update")
			return nil, nil
		},
		DoReplace: func(existing *projectClient.NamespacedCertificate) (*projectClient.NamespacedCertificate, error) {
			assert.FailInStub(tb, 2, "Unexpected call of Replace")
			return nil, nil
		},
		DoByID: func(id string) (*projectClient.NamespacedCertificate, error) {
			assert.FailInStub(tb, 2, "Unexpected call of ByID")
			return nil, nil
		},
		DoDelete: func(container *projectClient.NamespacedCertificate) error {
			assert.FailInStub(tb, 2, "Unexpected call of Delete")
			return nil
		},
	}
}

// NamespacedCertificateOperationsStub structure to hold callbacks used to stub github.com/rancher/types/client/project/v3/NamespacedCertificateOperations
type NamespacedCertificateOperationsStub struct {
	tb        testing.TB
	DoList    func(opts *types.ListOpts) (*projectClient.NamespacedCertificateCollection, error)
	DoCreate  func(opts *projectClient.NamespacedCertificate) (*projectClient.NamespacedCertificate, error)
	DoUpdate  func(existing *projectClient.NamespacedCertificate, updates interface{}) (*projectClient.NamespacedCertificate, error)
	DoReplace func(existing *projectClient.NamespacedCertificate) (*projectClient.NamespacedCertificate, error)
	DoByID    func(id string) (*projectClient.NamespacedCertificate, error)
	DoDelete  func(container *projectClient.NamespacedCertificate) error
}

// List implements github.com/rancher/types/client/project/v3/NamespacedCertificateOperations.List(...)
func (stub NamespacedCertificateOperationsStub) List(opts *types.ListOpts) (*projectClient.NamespacedCertificateCollection, error) {
	return stub.DoList(opts)
}

// Create implements github.com/rancher/types/client/project/v3/NamespacedCertificateOperations.Create(...)
func (stub NamespacedCertificateOperationsStub) Create(opts *projectClient.NamespacedCertificate) (*projectClient.NamespacedCertificate, error) {
	return stub.DoCreate(opts)
}

// Update implements github.com/rancher/types/client/project/v3/NamespacedCertificateOperations.Update(...)
func (stub NamespacedCertificateOperationsStub) Update(existing *projectClient.NamespacedCertificate, updates interface{}) (*projectClient.NamespacedCertificate, error) {
	return stub.DoUpdate(existing, updates)
}

// Replace implements github.com/rancher/types/client/project/v3/NamespacedCertificateOperations.Replace(...)
func (stub NamespacedCertificateOperationsStub) Replace(existing *projectClient.NamespacedCertificate) (*projectClient.NamespacedCertificate, error) {
	return stub.DoReplace(existing)
}

// ByID implements github.com/rancher/types/client/project/v3/NamespacedCertificateOperations.ByID(...)
func (stub NamespacedCertificateOperationsStub) ByID(id string) (*projectClient.NamespacedCertificate, error) {
	return stub.DoByID(id)
}

// Delete implements github.com/rancher/types/client/project/v3/NamespacedCertificateOperations.Delete(...)
func (stub NamespacedCertificateOperationsStub) Delete(container *projectClient.NamespacedCertificate) error {
	return stub.DoDelete(container)
}
