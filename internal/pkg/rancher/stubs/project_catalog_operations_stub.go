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

// CreateProjectCatalogOperationsStub creates a stub of github.com/rancher/types/client/project/v3/ProjectCatalogOperations
func CreateProjectCatalogOperationsStub(tb testing.TB) *ProjectCatalogOperationsStub {
	return &ProjectCatalogOperationsStub{
		tb: tb,
		DoList: func(opts *types.ListOpts) (*rancherClient.ProjectCatalogCollection, error) {
			assert.FailInStub(tb, 2, "Unexpected call of List")
			return nil, nil
		},
		DoCreate: func(opts *rancherClient.ProjectCatalog) (*rancherClient.ProjectCatalog, error) {
			assert.FailInStub(tb, 2, "Unexpected call of Create")
			return nil, nil
		},
		DoUpdate: func(existing *rancherClient.ProjectCatalog, updates interface{}) (*rancherClient.ProjectCatalog, error) {
			assert.FailInStub(tb, 2, "Unexpected call of Update")
			return nil, nil
		},
		DoReplace: func(existing *rancherClient.ProjectCatalog) (*rancherClient.ProjectCatalog, error) {
			assert.FailInStub(tb, 2, "Unexpected call of Replace")
			return nil, nil
		},
		DoByID: func(id string) (*rancherClient.ProjectCatalog, error) {
			assert.FailInStub(tb, 2, "Unexpected call of ByID")
			return nil, nil
		},
		DoDelete: func(container *rancherClient.ProjectCatalog) error {
			assert.FailInStub(tb, 2, "Unexpected call of Delete")
			return nil
		},
		DoActionRefresh: func(container *rancherClient.ProjectCatalog) (*rancherClient.CatalogRefresh, error) {
			assert.FailInStub(tb, 2, "Unexpected call of ActionRefresh")
			return nil, nil
		},
		DoCollectionActionRefresh: func(container *rancherClient.ProjectCatalogCollection) (*rancherClient.CatalogRefresh, error) {
			assert.FailInStub(tb, 2, "Unexpected call of CollectionActionRefresh")
			return nil, nil
		},
	}
}

// ProjectCatalogOperationsStub structure to hold callbacks used to stub github.com/rancher/types/client/project/v3/ProjectCatalogOperations
type ProjectCatalogOperationsStub struct {
	tb                        testing.TB
	DoList                    func(opts *types.ListOpts) (*rancherClient.ProjectCatalogCollection, error)
	DoCreate                  func(opts *rancherClient.ProjectCatalog) (*rancherClient.ProjectCatalog, error)
	DoUpdate                  func(existing *rancherClient.ProjectCatalog, updates interface{}) (*rancherClient.ProjectCatalog, error)
	DoReplace                 func(existing *rancherClient.ProjectCatalog) (*rancherClient.ProjectCatalog, error)
	DoByID                    func(id string) (*rancherClient.ProjectCatalog, error)
	DoDelete                  func(container *rancherClient.ProjectCatalog) error
	DoActionRefresh           func(container *rancherClient.ProjectCatalog) (*rancherClient.CatalogRefresh, error)
	DoCollectionActionRefresh func(container *rancherClient.ProjectCatalogCollection) (*rancherClient.CatalogRefresh, error)
}

// List implements github.com/rancher/types/client/project/v3/ProjectCatalogOperations.List(...)
func (stub ProjectCatalogOperationsStub) List(opts *types.ListOpts) (*rancherClient.ProjectCatalogCollection, error) {
	return stub.DoList(opts)
}

// Create implements github.com/rancher/types/client/project/v3/ProjectCatalogOperations.Create(...)
func (stub ProjectCatalogOperationsStub) Create(opts *rancherClient.ProjectCatalog) (*rancherClient.ProjectCatalog, error) {
	return stub.DoCreate(opts)
}

// Update implements github.com/rancher/types/client/project/v3/ProjectCatalogOperations.Update(...)
func (stub ProjectCatalogOperationsStub) Update(existing *rancherClient.ProjectCatalog, updates interface{}) (*rancherClient.ProjectCatalog, error) {
	return stub.DoUpdate(existing, updates)
}

// Replace implements github.com/rancher/types/client/project/v3/ProjectCatalogOperations.Replace(...)
func (stub ProjectCatalogOperationsStub) Replace(existing *rancherClient.ProjectCatalog) (*rancherClient.ProjectCatalog, error) {
	return stub.DoReplace(existing)
}

// ByID implements github.com/rancher/types/client/project/v3/ProjectCatalogOperations.ByID(...)
func (stub ProjectCatalogOperationsStub) ByID(id string) (*rancherClient.ProjectCatalog, error) {
	return stub.DoByID(id)
}

// Delete implements github.com/rancher/types/client/project/v3/ProjectCatalogOperations.Delete(...)
func (stub ProjectCatalogOperationsStub) Delete(container *rancherClient.ProjectCatalog) error {
	return stub.DoDelete(container)
}

// ActionRefresh implements github.com/rancher/types/client/project/v3/ProjectCatalogOperations.ActionRefresh(...)
func (stub ProjectCatalogOperationsStub) ActionRefresh(resource *rancherClient.ProjectCatalog) (*rancherClient.CatalogRefresh, error) {
	return stub.DoActionRefresh(resource)
}

// CollectionActionRefresh implements github.com/rancher/types/client/project/v3/ProjectCatalogOperations.CollectionActionRefresh(...)
func (stub ProjectCatalogOperationsStub) CollectionActionRefresh(resource *rancherClient.ProjectCatalogCollection) (*rancherClient.CatalogRefresh, error) {
	return stub.DoCollectionActionRefresh(resource)
}
