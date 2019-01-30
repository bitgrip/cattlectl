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

// CreateNamespaceOperationsStub creates a stub of github.com/rancher/types/client/cluster/v3/NamespaceOperations
func CreateNamespaceOperationsStub(tb testing.TB) *NamespaceOperationsStub {
	return &NamespaceOperationsStub{
		tb: tb,
		DoList: func(opts *types.ListOpts) (*clusterClient.NamespaceCollection, error) {
			assert.FailInStub(tb, 2, "Unexpected call of List")
			return nil, nil
		},
		DoCreate: func(opts *clusterClient.Namespace) (*clusterClient.Namespace, error) {
			assert.FailInStub(tb, 2, "Unexpected call of Create")
			return nil, nil
		},
	}
}

// NamespaceOperationsStub structure to hold callbacks used to stub github.com/rancher/types/client/cluster/v3/NamespaceOperations
type NamespaceOperationsStub struct {
	tb       testing.TB
	DoList   func(opts *types.ListOpts) (*clusterClient.NamespaceCollection, error)
	DoCreate func(opts *clusterClient.Namespace) (*clusterClient.Namespace, error)
}

// List implements github.com/rancher/types/client/cluster/v3/NamespaceOperations.List(...)
func (stub NamespaceOperationsStub) List(opts *types.ListOpts) (*clusterClient.NamespaceCollection, error) {
	return stub.DoList(opts)
}

// Create implements github.com/rancher/types/client/cluster/v3/NamespaceOperations.Create(...)
func (stub NamespaceOperationsStub) Create(opts *clusterClient.Namespace) (*clusterClient.Namespace, error) {
	return stub.DoCreate(opts)
}

// Update implements github.com/rancher/types/client/cluster/v3/NamespaceOperations.Update(...)
func (stub NamespaceOperationsStub) Update(existing *clusterClient.Namespace, updates interface{}) (*clusterClient.Namespace, error) {
	assert.FailInStub(stub.tb, 2, "Unexpected call of List")
	return nil, nil
}

// Replace implements github.com/rancher/types/client/cluster/v3/NamespaceOperations.Replace(...)
func (stub NamespaceOperationsStub) Replace(existing *clusterClient.Namespace) (*clusterClient.Namespace, error) {
	assert.FailInStub(stub.tb, 2, "Unexpected call of List")
	return nil, nil
}

// ByID implements github.com/rancher/types/client/cluster/v3/NamespaceOperations.ByID(...)
func (stub NamespaceOperationsStub) ByID(id string) (*clusterClient.Namespace, error) {
	assert.FailInStub(stub.tb, 2, "Unexpected call of List")
	return nil, nil
}

// Delete implements github.com/rancher/types/client/cluster/v3/NamespaceOperations.Delete(...)
func (stub NamespaceOperationsStub) Delete(container *clusterClient.Namespace) error {
	assert.FailInStub(stub.tb, 2, "Unexpected call of List")
	return nil
}

// ActionMove implements github.com/rancher/types/client/cluster/v3/NamespaceOperations.ActionMove(...)
func (stub NamespaceOperationsStub) ActionMove(resource *clusterClient.Namespace, input *clusterClient.NamespaceMove) error {
	assert.FailInStub(stub.tb, 2, "Unexpected call of List")
	return nil
}
