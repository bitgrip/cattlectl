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

package client

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/bitgrip/cattlectl/internal/pkg/assert"
	rancherModel "github.com/bitgrip/cattlectl/internal/pkg/rancher/model"
	"github.com/bitgrip/cattlectl/internal/pkg/rancher/stubs"
	"github.com/rancher/norman/types"
	backendRancherClient "github.com/rancher/types/client/management/v3"
	"github.com/sirupsen/logrus"
)

const (
	simpleCatalogName = "simple-catalog"
	simpleURL         = "http://simple-url"
	simpleBranch      = "simple-branch"
	simpleUsername    = "simple-username"
	simplePassword    = "simplePassword"
)

func Test_rancherCatalogClient_Exists(t *testing.T) {
	tests := []struct {
		name      string
		client    *rancherCatalogClient
		wanted    bool
		wantErr   bool
		wantedErr string
	}{
		{
			name: "Existing",
			client: existingRancherCatalogClient(
				t,
				simpleCatalogName,
				simpleURL,
				simpleBranch,
				simpleUsername,
				simplePassword,
			),
			wanted:  true,
			wantErr: false,
		},
		{
			name: "Not_Existing",
			client: notExistingRancherCatalogClient(
				t,
				simpleCatalogName,
				simpleURL,
				simpleBranch,
				simpleUsername,
				simplePassword,
			),
			wanted:  false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.client.Exists()
			if tt.wantErr {
				assert.NotOk(t, err, tt.wantedErr)
			} else {
				assert.Ok(t, err)
				assert.Equals(t, tt.wanted, got)
			}
		})
	}
}

func Test_rancherCatalogClient_Create(t *testing.T) {
	tests := []struct {
		name      string
		client    *rancherCatalogClient
		wantErr   bool
		wantedErr string
	}{
		{
			name: "Create",
			client: notExistingRancherCatalogClient(
				t,
				simpleCatalogName,
				simpleURL,
				simpleBranch,
				simpleUsername,
				simplePassword,
			),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.client.Create(false)
			if tt.wantErr {
				assert.NotOk(t, err, tt.wantedErr)
			} else {
				assert.Ok(t, err)
			}
		})
	}
}

func Test_rancherCatalogClient_Upgrade(t *testing.T) {
	tests := []struct {
		name      string
		client    *rancherCatalogClient
		wantErr   bool
		wantedErr string
	}{
		{
			name: "Create",
			client: existingRancherCatalogClient(
				t,
				simpleCatalogName,
				simpleURL,
				simpleBranch,
				simpleUsername,
				simplePassword,
			),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.client.Upgrade(false)
			if tt.wantErr {
				assert.NotOk(t, err, tt.wantedErr)
			} else {
				assert.Ok(t, err)
			}
		})
	}
}

func existingRancherCatalogClient(t *testing.T, name, url, branch, username, password string) *rancherCatalogClient {
	testClients := stubs.CreateBackendStubs(t)
	expectedListOpts := &types.ListOpts{
		Filters: map[string]interface{}{
			"name": name,
		},
	}
	rancherCatalogData := rancherModel.Catalog{
		Name:     name,
		URL:      url,
		Branch:   branch,
		Username: username,
		Password: password,
	}
	expectedBackendrCatalog := &backendRancherClient.Catalog{
		Name:     name,
		URL:      url,
		Branch:   branch,
		Username: username,
		Password: password,
		Labels:   map[string]string{"cattlectl.io/hash": "d20875c8c699ed126b385992bf8fc7c384f18e85"},
	}

	rancherCatalogOperationsStub := stubs.CreateRancherCatalogOperationsStub(t)
	rancherCatalogOperationsStub.DoList = func(opts *types.ListOpts) (*backendRancherClient.CatalogCollection, error) {
		if !reflect.DeepEqual(expectedListOpts, opts) {
			return nil, fmt.Errorf("Unexpected ListOpts %v", opts)
		}
		return &backendRancherClient.CatalogCollection{
			Data: []backendRancherClient.Catalog{
				backendRancherClient.Catalog{
					Name:   name,
					Labels: map[string]string{},
				},
			},
		}, nil
	}
	rancherCatalogOperationsStub.DoReplace = func(existing *backendRancherClient.Catalog) (*backendRancherClient.Catalog, error) {
		if !reflect.DeepEqual(expectedBackendrCatalog, existing) {
			return nil, fmt.Errorf("Unexpected ListOpts %v", existing)
		}
		return existing, nil
	}
	testClients.ManagementClient.Catalog = rancherCatalogOperationsStub
	rancherClient := simpleRancherClient()
	rancherClient._backendRancherClient = testClients.ManagementClient
	result, err := newRancherCatalogClient(
		name,
		rancherClient,
		logrus.New().WithFields(logrus.Fields{}),
	)
	assert.Ok(t, err)
	rancherCatalogClientResult := result.(*rancherCatalogClient)
	rancherCatalogClientResult.catalog = rancherCatalogData
	return rancherCatalogClientResult
}

func notExistingRancherCatalogClient(t *testing.T, name, url, branch, username, password string) *rancherCatalogClient {
	testClients := stubs.CreateBackendStubs(t)
	expectedListOpts := &types.ListOpts{
		Filters: map[string]interface{}{
			"name": name,
		},
	}
	rancherCatalogData := rancherModel.Catalog{
		Name:     name,
		URL:      url,
		Branch:   branch,
		Username: username,
		Password: password,
	}
	expectedBackendrCatalog := &backendRancherClient.Catalog{
		Name:     name,
		URL:      url,
		Branch:   branch,
		Username: username,
		Password: password,
		Labels:   map[string]string{"cattlectl.io/hash": hashOf(rancherCatalogData)},
	}

	rancherCatalogOperationsStub := stubs.CreateRancherCatalogOperationsStub(t)
	rancherCatalogOperationsStub.DoList = func(opts *types.ListOpts) (*backendRancherClient.CatalogCollection, error) {
		if !reflect.DeepEqual(expectedListOpts, opts) {
			return nil, fmt.Errorf("Unexpected ListOpts %v", opts)
		}
		return &backendRancherClient.CatalogCollection{
			Data: []backendRancherClient.Catalog{},
		}, nil
	}
	rancherCatalogOperationsStub.DoCreate = func(rancherCatalog *backendRancherClient.Catalog) (*backendRancherClient.Catalog, error) {
		if !reflect.DeepEqual(expectedBackendrCatalog, rancherCatalog) {
			return nil, fmt.Errorf("Unexpected Catalog %v", rancherCatalog)
		}

		return rancherCatalog, nil
	}
	testClients.ManagementClient.Catalog = rancherCatalogOperationsStub
	rancherClient := simpleRancherClient()
	rancherClient._backendRancherClient = testClients.ManagementClient
	result, err := newRancherCatalogClient(
		name,
		rancherClient,
		logrus.New().WithFields(logrus.Fields{}),
	)
	assert.Ok(t, err)
	rancherCatalogClientResult := result.(*rancherCatalogClient)
	rancherCatalogClientResult.catalog = rancherCatalogData
	return rancherCatalogClientResult
}
