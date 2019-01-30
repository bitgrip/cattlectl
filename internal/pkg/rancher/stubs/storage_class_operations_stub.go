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
	clusterClient "github.com/rancher/types/client/cluster/v3"
)

// CreateStorageClassOperationsStub creates a stub of github.com/rancher/types/client/cluster/v3/StorageClassOperations
func CreateStorageClassOperationsStub(tb testing.TB) *StorageClassOperationsStub {
	return &StorageClassOperationsStub{
		tb: tb,
		DoList: func(opts *types.ListOpts) (*clusterClient.StorageClassCollection, error) {
			assert.FailInStub(tb, 2, "Unexpected call of List")
			return nil, nil
		},
		DoCreate: func(opts *clusterClient.StorageClass) (*clusterClient.StorageClass, error) {
			assert.FailInStub(tb, 2, "Unexpected call of Create")
			return nil, nil
		},
	}
}

// StorageClassOperationsStub structure to hold callbacks used to stub github.com/rancher/types/client/cluster/v3/StorageClassOperations
type StorageClassOperationsStub struct {
	tb       testing.TB
	DoList   func(opts *types.ListOpts) (*clusterClient.StorageClassCollection, error)
	DoCreate func(opts *clusterClient.StorageClass) (*clusterClient.StorageClass, error)
}

// List implements github.com/rancher/types/client/cluster/v3/StorageClassOperations.List(...)
func (stub StorageClassOperationsStub) List(opts *types.ListOpts) (*clusterClient.StorageClassCollection, error) {
	return stub.DoList(opts)
}

// Create implements github.com/rancher/types/client/cluster/v3/StorageClassOperations.Create(...)
func (stub StorageClassOperationsStub) Create(opts *clusterClient.StorageClass) (*clusterClient.StorageClass, error) {
	return stub.DoCreate(opts)
}

// Update implements github.com/rancher/types/client/cluster/v3/StorageClassOperations.Update(...)
func (stub StorageClassOperationsStub) Update(existing *clusterClient.StorageClass, updates interface{}) (*clusterClient.StorageClass, error) {
	assert.FailInStub(stub.tb, 2, "Unexpected call of List")
	return nil, nil
}

// Replace implements github.com/rancher/types/client/cluster/v3/StorageClassOperations.Replace(...)
func (stub StorageClassOperationsStub) Replace(existing *clusterClient.StorageClass) (*clusterClient.StorageClass, error) {
	assert.FailInStub(stub.tb, 2, "Unexpected call of List")
	return nil, nil
}

// ByID implements github.com/rancher/types/client/cluster/v3/StorageClassOperations.ByID(...)
func (stub StorageClassOperationsStub) ByID(id string) (*clusterClient.StorageClass, error) {
	assert.FailInStub(stub.tb, 2, "Unexpected call of List")
	return nil, nil
}

// Delete implements github.com/rancher/types/client/cluster/v3/StorageClassOperations.Delete(...)
func (stub StorageClassOperationsStub) Delete(container *clusterClient.StorageClass) error {
	assert.FailInStub(stub.tb, 2, "Unexpected call of List")
	return nil
}
