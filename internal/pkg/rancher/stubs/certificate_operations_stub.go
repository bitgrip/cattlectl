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

// CreateCertificateOperationsStub creates a stub of github.com/rancher/types/client/project/v3/CertificateOperations
func CreateCertificateOperationsStub(tb testing.TB) *CertificateOperationsStub {
	return &CertificateOperationsStub{
		tb: tb,
		DoList: func(opts *types.ListOpts) (*projectClient.CertificateCollection, error) {
			assert.FailInStub(tb, 2, "Unexpected call of List")
			return nil, nil
		},
		DoCreate: func(opts *projectClient.Certificate) (*projectClient.Certificate, error) {
			assert.FailInStub(tb, 2, "Unexpected call of Create")
			return nil, nil
		},
		DoUpdate: func(existing *projectClient.Certificate, updates interface{}) (*projectClient.Certificate, error) {
			assert.FailInStub(tb, 2, "Unexpected call of Update")
			return nil, nil
		},
		DoReplace: func(existing *projectClient.Certificate) (*projectClient.Certificate, error) {
			assert.FailInStub(tb, 2, "Unexpected call of Replace")
			return nil, nil
		},
		DoByID: func(id string) (*projectClient.Certificate, error) {
			assert.FailInStub(tb, 2, "Unexpected call of ByID")
			return nil, nil
		},
		DoDelete: func(container *projectClient.Certificate) error {
			assert.FailInStub(tb, 2, "Unexpected call of Delete")
			return nil
		},
	}
}

// CertificateOperationsStub structure to hold callbacks used to stub github.com/rancher/types/client/project/v3/CertificateOperations
type CertificateOperationsStub struct {
	tb        testing.TB
	DoList    func(opts *types.ListOpts) (*projectClient.CertificateCollection, error)
	DoCreate  func(opts *projectClient.Certificate) (*projectClient.Certificate, error)
	DoUpdate  func(existing *projectClient.Certificate, updates interface{}) (*projectClient.Certificate, error)
	DoReplace func(existing *projectClient.Certificate) (*projectClient.Certificate, error)
	DoByID    func(id string) (*projectClient.Certificate, error)
	DoDelete  func(container *projectClient.Certificate) error
}

// List implements github.com/rancher/types/client/project/v3/CertificateOperations.List(...)
func (stub CertificateOperationsStub) List(opts *types.ListOpts) (*projectClient.CertificateCollection, error) {
	return stub.DoList(opts)
}

// Create implements github.com/rancher/types/client/project/v3/CertificateOperations.Create(...)
func (stub CertificateOperationsStub) Create(opts *projectClient.Certificate) (*projectClient.Certificate, error) {
	return stub.DoCreate(opts)
}

// Update implements github.com/rancher/types/client/project/v3/CertificateOperations.Update(...)
func (stub CertificateOperationsStub) Update(existing *projectClient.Certificate, updates interface{}) (*projectClient.Certificate, error) {
	return stub.DoUpdate(existing, updates)
}

// Replace implements github.com/rancher/types/client/project/v3/CertificateOperations.Replace(...)
func (stub CertificateOperationsStub) Replace(existing *projectClient.Certificate) (*projectClient.Certificate, error) {
	return stub.DoReplace(existing)
}

// ByID implements github.com/rancher/types/client/project/v3/CertificateOperations.ByID(...)
func (stub CertificateOperationsStub) ByID(id string) (*projectClient.Certificate, error) {
	return stub.DoByID(id)
}

// Delete implements github.com/rancher/types/client/project/v3/CertificateOperations.Delete(...)
func (stub CertificateOperationsStub) Delete(container *projectClient.Certificate) error {
	return stub.DoDelete(container)
}
