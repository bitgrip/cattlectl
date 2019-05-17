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

// CreateCronJobOperationsStub creates a stub of github.com/rancher/types/client/project/v3/CronJobOperations
func CreateCronJobOperationsStub(tb testing.TB) *CronJobOperationsStub {
	return &CronJobOperationsStub{
		tb: tb,
		DoList: func(opts *types.ListOpts) (*projectClient.CronJobCollection, error) {
			assert.FailInStub(tb, 2, "Unexpected call of List")
			return nil, nil
		},
		DoCreate: func(opts *projectClient.CronJob) (*projectClient.CronJob, error) {
			assert.FailInStub(tb, 2, "Unexpected call of Create")
			return nil, nil
		},
		DoUpdate: func(existing *projectClient.CronJob, updates interface{}) (*projectClient.CronJob, error) {
			assert.FailInStub(tb, 2, "Unexpected call of List")
			return nil, nil
		},
		DoReplace: func(existing *projectClient.CronJob) (*projectClient.CronJob, error) {
			assert.FailInStub(tb, 2, "Unexpected call of List")
			return nil, nil
		},
	}
}

// CronJobOperationsStub structure to hold callbacks used to stub github.com/rancher/types/client/project/v3/CronJobOperations
type CronJobOperationsStub struct {
	tb        testing.TB
	DoList    func(opts *types.ListOpts) (*projectClient.CronJobCollection, error)
	DoCreate  func(opts *projectClient.CronJob) (*projectClient.CronJob, error)
	DoUpdate  func(existing *projectClient.CronJob, updates interface{}) (*projectClient.CronJob, error)
	DoReplace func(existing *projectClient.CronJob) (*projectClient.CronJob, error)
}

// List implements github.com/rancher/types/client/project/v3/CronJobOperations.List(...)
func (stub CronJobOperationsStub) List(opts *types.ListOpts) (*projectClient.CronJobCollection, error) {
	return stub.DoList(opts)
}

// Create implements github.com/rancher/types/client/project/v3/CronJobOperations.Create(...)
func (stub CronJobOperationsStub) Create(opts *projectClient.CronJob) (*projectClient.CronJob, error) {
	return stub.DoCreate(opts)
}

// Update implements github.com/rancher/types/client/project/v3/CronJobOperations.Update(...)
func (stub CronJobOperationsStub) Update(existing *projectClient.CronJob, updates interface{}) (*projectClient.CronJob, error) {
	return stub.DoUpdate(existing, updates)
}

// Replace implements github.com/rancher/types/client/project/v3/CronJobOperations.Replace(...)
func (stub CronJobOperationsStub) Replace(existing *projectClient.CronJob) (*projectClient.CronJob, error) {
	return stub.DoReplace(existing)
}

// ByID implements github.com/rancher/types/client/project/v3/CronJobOperations.ByID(...)
func (stub CronJobOperationsStub) ByID(id string) (*projectClient.CronJob, error) {
	assert.FailInStub(stub.tb, 2, "Unexpected call of List")
	return nil, nil
}

// Delete implements github.com/rancher/types/client/project/v3/CronJobOperations.Delete(...)
func (stub CronJobOperationsStub) Delete(container *projectClient.CronJob) error {
	assert.FailInStub(stub.tb, 2, "Unexpected call of List")
	return nil
}
