// Copyright © 2019 Bitgrip <berlin@bitgrip.de>
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
	rancherClient "github.com/rancher/types/client/management/v3"
)

// CreateRancherCatalogOperationsStub creates a stub of github.com/rancher/types/client/project/v3/RancherCatalogOperations
func CreateRancherCatalogOperationsStub(tb testing.TB) *RancherCatalogOperationsStub {
	return &RancherCatalogOperationsStub{
		tb: tb,
		DoList: func(opts *types.ListOpts) (*rancherClient.CatalogCollection, error) {
			assert.FailInStub(tb, 2, "Unexpected call of List")
			return nil, nil
		},
		DoCreate: func(opts *rancherClient.Catalog) (*rancherClient.Catalog, error) {
			assert.FailInStub(tb, 2, "Unexpected call of Create")
			return nil, nil
		},
		DoUpdate: func(existing *rancherClient.Catalog, updates interface{}) (*rancherClient.Catalog, error) {
			assert.FailInStub(tb, 2, "Unexpected call of Update")
			return nil, nil
		},
		DoReplace: func(existing *rancherClient.Catalog) (*rancherClient.Catalog, error) {
			assert.FailInStub(tb, 2, "Unexpected call of Replace")
			return nil, nil
		},
		DoByID: func(id string) (*rancherClient.Catalog, error) {
			assert.FailInStub(tb, 2, "Unexpected call of ByID")
			return nil, nil
		},
		DoDelete: func(container *rancherClient.Catalog) error {
			assert.FailInStub(tb, 2, "Unexpected call of Delete")
			return nil
		},
		DoActionRefresh: func(container *rancherClient.Catalog) error {
			assert.FailInStub(tb, 2, "Unexpected call of ActionRefresh")
			return nil
		},
		DoCollectionActionRefresh: func(container *rancherClient.CatalogCollection) error {
			assert.FailInStub(tb, 2, "Unexpected call of CollectionActionRefresh")
			return nil
		},
	}
}

// RancherCatalogOperationsStub structure to hold callbacks used to stub github.com/rancher/types/client/project/v3/RancherCatalogOperations
type RancherCatalogOperationsStub struct {
	tb                        testing.TB
	DoList                    func(opts *types.ListOpts) (*rancherClient.CatalogCollection, error)
	DoCreate                  func(opts *rancherClient.Catalog) (*rancherClient.Catalog, error)
	DoUpdate                  func(existing *rancherClient.Catalog, updates interface{}) (*rancherClient.Catalog, error)
	DoReplace                 func(existing *rancherClient.Catalog) (*rancherClient.Catalog, error)
	DoByID                    func(id string) (*rancherClient.Catalog, error)
	DoDelete                  func(container *rancherClient.Catalog) error
	DoActionRefresh           func(container *rancherClient.Catalog) error
	DoCollectionActionRefresh func(container *rancherClient.CatalogCollection) error
}

// List implements github.com/rancher/types/client/project/v3/RancherCatalogOperations.List(...)
func (stub RancherCatalogOperationsStub) List(opts *types.ListOpts) (*rancherClient.CatalogCollection, error) {
	return stub.DoList(opts)
}

// Create implements github.com/rancher/types/client/project/v3/RancherCatalogOperations.Create(...)
func (stub RancherCatalogOperationsStub) Create(opts *rancherClient.Catalog) (*rancherClient.Catalog, error) {
	return stub.DoCreate(opts)
}

// Update implements github.com/rancher/types/client/project/v3/RancherCatalogOperations.Update(...)
func (stub RancherCatalogOperationsStub) Update(existing *rancherClient.Catalog, updates interface{}) (*rancherClient.Catalog, error) {
	return stub.DoUpdate(existing, updates)
}

// Replace implements github.com/rancher/types/client/project/v3/RancherCatalogOperations.Replace(...)
func (stub RancherCatalogOperationsStub) Replace(existing *rancherClient.Catalog) (*rancherClient.Catalog, error) {
	return stub.DoReplace(existing)
}

// ByID implements github.com/rancher/types/client/project/v3/RancherCatalogOperations.ByID(...)
func (stub RancherCatalogOperationsStub) ByID(id string) (*rancherClient.Catalog, error) {
	return stub.DoByID(id)
}

// Delete implements github.com/rancher/types/client/project/v3/RancherCatalogOperations.Delete(...)
func (stub RancherCatalogOperationsStub) Delete(container *rancherClient.Catalog) error {
	return stub.DoDelete(container)
}

func (stub RancherCatalogOperationsStub) ActionRefresh(resource *rancherClient.Catalog) error {
	return stub.DoActionRefresh(resource)
}

func (stub RancherCatalogOperationsStub) CollectionActionRefresh(resource *rancherClient.CatalogCollection) error {
	return stub.DoCollectionActionRefresh(resource)
}
