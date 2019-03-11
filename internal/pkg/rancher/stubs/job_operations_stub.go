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

// CreateJobOperationsStub creates a stub of github.com/rancher/types/client/project/v3/JobOperations
func CreateJobOperationsStub(tb testing.TB) *JobOperationsStub {
	return &JobOperationsStub{
		tb: tb,
		DoList: func(opts *types.ListOpts) (*projectClient.JobCollection, error) {
			assert.FailInStub(tb, 2, "Unexpected call of List")
			return nil, nil
		},
		DoCreate: func(opts *projectClient.Job) (*projectClient.Job, error) {
			assert.FailInStub(tb, 2, "Unexpected call of Create")
			return nil, nil
		},
		DoUpdate: func(existing *projectClient.Job, updates interface{}) (*projectClient.Job, error) {
			assert.FailInStub(tb, 2, "Unexpected call of List")
			return nil, nil
		},
		DoReplace: func(existing *projectClient.Job) (*projectClient.Job, error) {
			assert.FailInStub(tb, 2, "Unexpected call of List")
			return nil, nil
		},
	}
}

// JobOperationsStub structure to hold callbacks used to stub github.com/rancher/types/client/project/v3/JobOperations
type JobOperationsStub struct {
	tb        testing.TB
	DoList    func(opts *types.ListOpts) (*projectClient.JobCollection, error)
	DoCreate  func(opts *projectClient.Job) (*projectClient.Job, error)
	DoUpdate  func(existing *projectClient.Job, updates interface{}) (*projectClient.Job, error)
	DoReplace func(existing *projectClient.Job) (*projectClient.Job, error)
}

// List implements github.com/rancher/types/client/project/v3/JobOperations.List(...)
func (stub JobOperationsStub) List(opts *types.ListOpts) (*projectClient.JobCollection, error) {
	return stub.DoList(opts)
}

// Create implements github.com/rancher/types/client/project/v3/JobOperations.Create(...)
func (stub JobOperationsStub) Create(opts *projectClient.Job) (*projectClient.Job, error) {
	return stub.DoCreate(opts)
}

// Update implements github.com/rancher/types/client/project/v3/JobOperations.Update(...)
func (stub JobOperationsStub) Update(existing *projectClient.Job, updates interface{}) (*projectClient.Job, error) {
	return stub.DoUpdate(existing, updates)
}

// Replace implements github.com/rancher/types/client/project/v3/JobOperations.Replace(...)
func (stub JobOperationsStub) Replace(existing *projectClient.Job) (*projectClient.Job, error) {
	return stub.DoReplace(existing)
}

// ByID implements github.com/rancher/types/client/project/v3/JobOperations.ByID(...)
func (stub JobOperationsStub) ByID(id string) (*projectClient.Job, error) {
	assert.FailInStub(stub.tb, 2, "Unexpected call of List")
	return nil, nil
}

// Delete implements github.com/rancher/types/client/project/v3/JobOperations.Delete(...)
func (stub JobOperationsStub) Delete(container *projectClient.Job) error {
	assert.FailInStub(stub.tb, 2, "Unexpected call of List")
	return nil
}
