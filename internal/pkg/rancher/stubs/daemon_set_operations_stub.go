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

// CreateDaemonSetOperationsStub creates a stub of github.com/rancher/types/client/project/v3/DaemonSetOperations
func CreateDaemonSetOperationsStub(tb testing.TB) *DaemonSetOperationsStub {
	return &DaemonSetOperationsStub{
		tb: tb,
		DoList: func(opts *types.ListOpts) (*projectClient.DaemonSetCollection, error) {
			assert.FailInStub(tb, 2, "Unexpected call of List")
			return nil, nil
		},
		DoCreate: func(opts *projectClient.DaemonSet) (*projectClient.DaemonSet, error) {
			assert.FailInStub(tb, 2, "Unexpected call of Create")
			return nil, nil
		},
		DoUpdate: func(existing *projectClient.DaemonSet, updates interface{}) (*projectClient.DaemonSet, error) {
			assert.FailInStub(tb, 2, "Unexpected call of List")
			return nil, nil
		},
		DoReplace: func(existing *projectClient.DaemonSet) (*projectClient.DaemonSet, error) {
			assert.FailInStub(tb, 2, "Unexpected call of List")
			return nil, nil
		},
	}
}

// DaemonSetOperationsStub structure to hold callbacks used to stub github.com/rancher/types/client/project/v3/DaemonSetOperations
type DaemonSetOperationsStub struct {
	tb        testing.TB
	DoList    func(opts *types.ListOpts) (*projectClient.DaemonSetCollection, error)
	DoCreate  func(opts *projectClient.DaemonSet) (*projectClient.DaemonSet, error)
	DoUpdate  func(existing *projectClient.DaemonSet, updates interface{}) (*projectClient.DaemonSet, error)
	DoReplace func(existing *projectClient.DaemonSet) (*projectClient.DaemonSet, error)
}

// List implements github.com/rancher/types/client/project/v3/DaemonSetOperations.List(...)
func (stub DaemonSetOperationsStub) List(opts *types.ListOpts) (*projectClient.DaemonSetCollection, error) {
	return stub.DoList(opts)
}

// Create implements github.com/rancher/types/client/project/v3/DaemonSetOperations.Create(...)
func (stub DaemonSetOperationsStub) Create(opts *projectClient.DaemonSet) (*projectClient.DaemonSet, error) {
	return stub.DoCreate(opts)
}

// Update implements github.com/rancher/types/client/project/v3/DaemonSetOperations.Update(...)
func (stub DaemonSetOperationsStub) Update(existing *projectClient.DaemonSet, updates interface{}) (*projectClient.DaemonSet, error) {
	return stub.DoUpdate(existing, updates)
}

// Replace implements github.com/rancher/types/client/project/v3/DaemonSetOperations.Replace(...)
func (stub DaemonSetOperationsStub) Replace(existing *projectClient.DaemonSet) (*projectClient.DaemonSet, error) {
	return stub.DoReplace(existing)
}

// ByID implements github.com/rancher/types/client/project/v3/DaemonSetOperations.ByID(...)
func (stub DaemonSetOperationsStub) ByID(id string) (*projectClient.DaemonSet, error) {
	assert.FailInStub(stub.tb, 2, "Unexpected call of List")
	return nil, nil
}

// Delete implements github.com/rancher/types/client/project/v3/DaemonSetOperations.Delete(...)
func (stub DaemonSetOperationsStub) Delete(container *projectClient.DaemonSet) error {
	assert.FailInStub(stub.tb, 2, "Unexpected call of List")
	return nil
}
