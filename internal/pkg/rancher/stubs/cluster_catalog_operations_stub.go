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

// CreateClusterCatalogOperationsStub creates a stub of github.com/rancher/types/client/project/v3/ClusterCatalogOperations
func CreateClusterCatalogOperationsStub(tb testing.TB) *ClusterCatalogOperationsStub {
	return &ClusterCatalogOperationsStub{
		tb: tb,
		DoList: func(opts *types.ListOpts) (*rancherClient.ClusterCatalogCollection, error) {
			assert.FailInStub(tb, 2, "Unexpected call of List")
			return nil, nil
		},
		DoCreate: func(opts *rancherClient.ClusterCatalog) (*rancherClient.ClusterCatalog, error) {
			assert.FailInStub(tb, 2, "Unexpected call of Create")
			return nil, nil
		},
		DoUpdate: func(existing *rancherClient.ClusterCatalog, updates interface{}) (*rancherClient.ClusterCatalog, error) {
			assert.FailInStub(tb, 2, "Unexpected call of Update")
			return nil, nil
		},
		DoReplace: func(existing *rancherClient.ClusterCatalog) (*rancherClient.ClusterCatalog, error) {
			assert.FailInStub(tb, 2, "Unexpected call of Replace")
			return nil, nil
		},
		DoByID: func(id string) (*rancherClient.ClusterCatalog, error) {
			assert.FailInStub(tb, 2, "Unexpected call of ByID")
			return nil, nil
		},
		DoDelete: func(container *rancherClient.ClusterCatalog) error {
			assert.FailInStub(tb, 2, "Unexpected call of Delete")
			return nil
		},
		DoActionRefresh: func(container *rancherClient.ClusterCatalog) (*rancherClient.CatalogRefresh, error) {
			assert.FailInStub(tb, 2, "Unexpected call of ActionRefresh")
			return nil, nil
		},
		DoCollectionActionRefresh: func(container *rancherClient.ClusterCatalogCollection) (*rancherClient.CatalogRefresh, error) {
			assert.FailInStub(tb, 2, "Unexpected call of CollectionActionRefresh")
			return nil, nil
		},
	}
}

// ClusterCatalogOperationsStub structure to hold callbacks used to stub github.com/rancher/types/client/project/v3/ClusterCatalogOperations
type ClusterCatalogOperationsStub struct {
	tb                        testing.TB
	DoList                    func(opts *types.ListOpts) (*rancherClient.ClusterCatalogCollection, error)
	DoCreate                  func(opts *rancherClient.ClusterCatalog) (*rancherClient.ClusterCatalog, error)
	DoUpdate                  func(existing *rancherClient.ClusterCatalog, updates interface{}) (*rancherClient.ClusterCatalog, error)
	DoReplace                 func(existing *rancherClient.ClusterCatalog) (*rancherClient.ClusterCatalog, error)
	DoByID                    func(id string) (*rancherClient.ClusterCatalog, error)
	DoDelete                  func(container *rancherClient.ClusterCatalog) error
	DoActionRefresh           func(container *rancherClient.ClusterCatalog) (*rancherClient.CatalogRefresh, error)
	DoCollectionActionRefresh func(container *rancherClient.ClusterCatalogCollection) (*rancherClient.CatalogRefresh, error)
}

// List implements github.com/rancher/types/client/project/v3/ClusterCatalogOperations.List(...)
func (stub ClusterCatalogOperationsStub) List(opts *types.ListOpts) (*rancherClient.ClusterCatalogCollection, error) {
	return stub.DoList(opts)
}

// Create implements github.com/rancher/types/client/project/v3/ClusterCatalogOperations.Create(...)
func (stub ClusterCatalogOperationsStub) Create(opts *rancherClient.ClusterCatalog) (*rancherClient.ClusterCatalog, error) {
	return stub.DoCreate(opts)
}

// Update implements github.com/rancher/types/client/project/v3/ClusterCatalogOperations.Update(...)
func (stub ClusterCatalogOperationsStub) Update(existing *rancherClient.ClusterCatalog, updates interface{}) (*rancherClient.ClusterCatalog, error) {
	return stub.DoUpdate(existing, updates)
}

// Replace implements github.com/rancher/types/client/project/v3/ClusterCatalogOperations.Replace(...)
func (stub ClusterCatalogOperationsStub) Replace(existing *rancherClient.ClusterCatalog) (*rancherClient.ClusterCatalog, error) {
	return stub.DoReplace(existing)
}

// ByID implements github.com/rancher/types/client/project/v3/ClusterCatalogOperations.ByID(...)
func (stub ClusterCatalogOperationsStub) ByID(id string) (*rancherClient.ClusterCatalog, error) {
	return stub.DoByID(id)
}

// Delete implements github.com/rancher/types/client/project/v3/ClusterCatalogOperations.Delete(...)
func (stub ClusterCatalogOperationsStub) Delete(container *rancherClient.ClusterCatalog) error {
	return stub.DoDelete(container)
}

// ActionRefresh implements github.com/rancher/types/client/project/v3/ClusterCatalogOperations.ActionRefresh(...)
func (stub ClusterCatalogOperationsStub) ActionRefresh(resource *rancherClient.ClusterCatalog) (*rancherClient.CatalogRefresh, error) {
	return stub.DoActionRefresh(resource)
}

// CollectionActionRefresh implements github.com/rancher/types/client/project/v3/ClusterCatalogOperations.CollectionActionRefresh(...)
func (stub ClusterCatalogOperationsStub) CollectionActionRefresh(resource *rancherClient.ClusterCatalogCollection) (*rancherClient.CatalogRefresh, error) {
	return stub.DoCollectionActionRefresh(resource)
}
