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

// CreateStatefulSetOperationsStub creates a stub of github.com/rancher/types/client/project/v3/StatefulSetOperations
func CreateStatefulSetOperationsStub(tb testing.TB) *StatefulSetOperationsStub {
	return &StatefulSetOperationsStub{
		tb: tb,
		DoList: func(opts *types.ListOpts) (*projectClient.StatefulSetCollection, error) {
			assert.FailInStub(tb, 2, "Unexpected call of List")
			return nil, nil
		},
		DoCreate: func(opts *projectClient.StatefulSet) (*projectClient.StatefulSet, error) {
			assert.FailInStub(tb, 2, "Unexpected call of Create")
			return nil, nil
		},
		DoUpdate: func(existing *projectClient.StatefulSet, updates interface{}) (*projectClient.StatefulSet, error) {
			assert.FailInStub(tb, 2, "Unexpected call of List")
			return nil, nil
		},
		DoReplace: func(existing *projectClient.StatefulSet) (*projectClient.StatefulSet, error) {
			assert.FailInStub(tb, 2, "Unexpected call of List")
			return nil, nil
		},
	}
}

// StatefulSetOperationsStub structure to hold callbacks used to stub github.com/rancher/types/client/project/v3/StatefulSetOperations
type StatefulSetOperationsStub struct {
	tb        testing.TB
	DoList    func(opts *types.ListOpts) (*projectClient.StatefulSetCollection, error)
	DoCreate  func(opts *projectClient.StatefulSet) (*projectClient.StatefulSet, error)
	DoUpdate  func(existing *projectClient.StatefulSet, updates interface{}) (*projectClient.StatefulSet, error)
	DoReplace func(existing *projectClient.StatefulSet) (*projectClient.StatefulSet, error)
}

// List implements github.com/rancher/types/client/project/v3/StatefulSetOperations.List(...)
func (stub StatefulSetOperationsStub) List(opts *types.ListOpts) (*projectClient.StatefulSetCollection, error) {
	return stub.DoList(opts)
}

// Create implements github.com/rancher/types/client/project/v3/StatefulSetOperations.Create(...)
func (stub StatefulSetOperationsStub) Create(opts *projectClient.StatefulSet) (*projectClient.StatefulSet, error) {
	return stub.DoCreate(opts)
}

// Update implements github.com/rancher/types/client/project/v3/StatefulSetOperations.Update(...)
func (stub StatefulSetOperationsStub) Update(existing *projectClient.StatefulSet, updates interface{}) (*projectClient.StatefulSet, error) {
	return stub.DoUpdate(existing, updates)
}

// Replace implements github.com/rancher/types/client/project/v3/StatefulSetOperations.Replace(...)
func (stub StatefulSetOperationsStub) Replace(existing *projectClient.StatefulSet) (*projectClient.StatefulSet, error) {
	return stub.DoReplace(existing)
}

// ByID implements github.com/rancher/types/client/project/v3/StatefulSetOperations.ByID(...)
func (stub StatefulSetOperationsStub) ByID(id string) (*projectClient.StatefulSet, error) {
	assert.FailInStub(stub.tb, 2, "Unexpected call of List")
	return nil, nil
}

// Delete implements github.com/rancher/types/client/project/v3/StatefulSetOperations.Delete(...)
func (stub StatefulSetOperationsStub) Delete(container *projectClient.StatefulSet) error {
	assert.FailInStub(stub.tb, 2, "Unexpected call of List")
	return nil
}
