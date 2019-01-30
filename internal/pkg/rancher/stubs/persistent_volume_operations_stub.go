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

// CreatePersistentVolumeOperationsStub creates a stub of github.com/rancher/types/client/cluster/v3/PersistentVolumeOperations
func CreatePersistentVolumeOperationsStub(tb testing.TB) *PersistentVolumeOperationsStub {
	return &PersistentVolumeOperationsStub{
		tb: tb,
		DoList: func(opts *types.ListOpts) (*clusterClient.PersistentVolumeCollection, error) {
			assert.FailInStub(tb, 2, "Unexpected call of List")
			return nil, nil
		},
		DoCreate: func(opts *clusterClient.PersistentVolume) (*clusterClient.PersistentVolume, error) {
			assert.FailInStub(tb, 2, "Unexpected call of Create")
			return nil, nil
		},
	}
}

// PersistentVolumeOperationsStub structure to hold callbacks used to stub github.com/rancher/types/client/cluster/v3/PersistentVolumeOperations
type PersistentVolumeOperationsStub struct {
	tb       testing.TB
	DoList   func(opts *types.ListOpts) (*clusterClient.PersistentVolumeCollection, error)
	DoCreate func(opts *clusterClient.PersistentVolume) (*clusterClient.PersistentVolume, error)
}

// List implements github.com/rancher/types/client/cluster/v3/PersistentVolumeOperations.List(...)
func (stub PersistentVolumeOperationsStub) List(opts *types.ListOpts) (*clusterClient.PersistentVolumeCollection, error) {
	return stub.DoList(opts)
}

// Create implements github.com/rancher/types/client/cluster/v3/PersistentVolumeOperations.Create(...)
func (stub PersistentVolumeOperationsStub) Create(opts *clusterClient.PersistentVolume) (*clusterClient.PersistentVolume, error) {
	return stub.DoCreate(opts)
}

// Update implements github.com/rancher/types/client/cluster/v3/PersistentVolumeOperations.Update(...)
func (stub PersistentVolumeOperationsStub) Update(existing *clusterClient.PersistentVolume, updates interface{}) (*clusterClient.PersistentVolume, error) {
	assert.FailInStub(stub.tb, 2, "Unexpected call of List")
	return nil, nil
}

// Replace implements github.com/rancher/types/client/cluster/v3/PersistentVolumeOperations.Replace(...)
func (stub PersistentVolumeOperationsStub) Replace(existing *clusterClient.PersistentVolume) (*clusterClient.PersistentVolume, error) {
	assert.FailInStub(stub.tb, 2, "Unexpected call of List")
	return nil, nil
}

// ByID implements github.com/rancher/types/client/cluster/v3/PersistentVolumeOperations.ByID(...)
func (stub PersistentVolumeOperationsStub) ByID(id string) (*clusterClient.PersistentVolume, error) {
	assert.FailInStub(stub.tb, 2, "Unexpected call of List")
	return nil, nil
}

// Delete implements github.com/rancher/types/client/cluster/v3/PersistentVolumeOperations.Delete(...)
func (stub PersistentVolumeOperationsStub) Delete(container *clusterClient.PersistentVolume) error {
	assert.FailInStub(stub.tb, 2, "Unexpected call of List")
	return nil
}
