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

func Test_clusterCatalogClient_Exists(t *testing.T) {
	tests := []struct {
		name      string
		client    *clusterCatalogClient
		wanted    bool
		wantErr   bool
		wantedErr string
	}{
		{
			name: "Existing",
			client: existingClusterCatalogClient(
				t,
				simpleCatalogName,
				simpleClusterID,
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
			client: notExistingClusterCatalogClient(
				t,
				simpleCatalogName,
				simpleClusterID,
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

func Test_clusterCatalogClient_Create(t *testing.T) {
	tests := []struct {
		name      string
		client    *clusterCatalogClient
		wantErr   bool
		wantedErr string
	}{
		{
			name: "Create",
			client: notExistingClusterCatalogClient(
				t,
				simpleCatalogName,
				simpleClusterID,
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

func Test_clusterCatalogClient_Upgrade(t *testing.T) {
	tests := []struct {
		name      string
		client    *clusterCatalogClient
		wantErr   bool
		wantedErr string
	}{
		{
			name: "Create",
			client: existingClusterCatalogClient(
				t,
				simpleCatalogName,
				simpleClusterID,
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

func existingClusterCatalogClient(t *testing.T, name, clusterID, url, branch, username, password string) *clusterCatalogClient {
	testClients := stubs.CreateBackendStubs(t)
	expectedListOpts := &types.ListOpts{
		Filters: map[string]interface{}{
			"name":      name,
			"clusterId": clusterID,
		},
	}
	clusterCatalogData := rancherModel.Catalog{
		Name:     name,
		URL:      url,
		Branch:   branch,
		Username: username,
		Password: password,
	}
	expectedBackendrCatalog := &backendRancherClient.ClusterCatalog{
		Name:      name,
		ClusterID: clusterID,
		URL:       url,
		Branch:    branch,
		Username:  username,
		Password:  password,
		Labels:    map[string]string{"cattlectl.io/hash": hashOf(clusterCatalogData)},
	}

	clusterCatalogOperationsStub := stubs.CreateClusterCatalogOperationsStub(t)
	clusterCatalogOperationsStub.DoList = func(opts *types.ListOpts) (*backendRancherClient.ClusterCatalogCollection, error) {
		if !reflect.DeepEqual(expectedListOpts, opts) {
			return nil, fmt.Errorf("Unexpected ListOpts %v", opts)
		}
		return &backendRancherClient.ClusterCatalogCollection{
			Data: []backendRancherClient.ClusterCatalog{
				backendRancherClient.ClusterCatalog{
					Name:      name,
					ClusterID: clusterID,
					Labels:    map[string]string{},
				},
			},
		}, nil
	}
	clusterCatalogOperationsStub.DoReplace = func(existing *backendRancherClient.ClusterCatalog) (*backendRancherClient.ClusterCatalog, error) {
		if !reflect.DeepEqual(expectedBackendrCatalog, existing) {
			return nil, fmt.Errorf("Unexpected ClusterCatalog %v", existing)
		}
		return existing, nil
	}
	testClients.ManagementClient.ClusterCatalog = clusterCatalogOperationsStub
	rancherClient := simpleRancherClient()
	rancherClient._backendRancherClient = testClients.ManagementClient
	clusterClient := simpleClusterClient()
	clusterClient.rancherClient = rancherClient
	result, err := newClusterCatalogClient(
		name,
		clusterClient,
		logrus.New().WithFields(logrus.Fields{}),
	)
	assert.Ok(t, err)
	clusterCatalogClientResult := result.(*clusterCatalogClient)
	clusterCatalogClientResult.catalog = clusterCatalogData
	return clusterCatalogClientResult
}

func notExistingClusterCatalogClient(t *testing.T, name, clusterID, url, branch, username, password string) *clusterCatalogClient {
	testClients := stubs.CreateBackendStubs(t)
	expectedListOpts := &types.ListOpts{
		Filters: map[string]interface{}{
			"name":      name,
			"clusterId": clusterID,
		},
	}
	clusterCatalogData := rancherModel.Catalog{
		Name:     name,
		URL:      url,
		Branch:   branch,
		Username: username,
		Password: password,
	}
	expectedBackendrCatalog := &backendRancherClient.ClusterCatalog{
		Name:      name,
		ClusterID: clusterID,
		URL:       url,
		Branch:    branch,
		Username:  username,
		Password:  password,
		Labels:    map[string]string{"cattlectl.io/hash": hashOf(clusterCatalogData)},
	}

	clusterCatalogOperationsStub := stubs.CreateClusterCatalogOperationsStub(t)
	clusterCatalogOperationsStub.DoList = func(opts *types.ListOpts) (*backendRancherClient.ClusterCatalogCollection, error) {
		if !reflect.DeepEqual(expectedListOpts, opts) {
			return nil, fmt.Errorf("Unexpected ListOpts %v", opts)
		}
		return &backendRancherClient.ClusterCatalogCollection{
			Data: []backendRancherClient.ClusterCatalog{},
		}, nil
	}
	clusterCatalogOperationsStub.DoCreate = func(clusterCatalog *backendRancherClient.ClusterCatalog) (*backendRancherClient.ClusterCatalog, error) {
		if !reflect.DeepEqual(expectedBackendrCatalog, clusterCatalog) {
			return nil, fmt.Errorf("Unexpected Catalog\n%v\n%v", expectedBackendrCatalog, clusterCatalog)
		}

		return clusterCatalog, nil
	}
	testClients.ManagementClient.ClusterCatalog = clusterCatalogOperationsStub
	rancherClient := simpleRancherClient()
	rancherClient._backendRancherClient = testClients.ManagementClient
	clusterClient := simpleClusterClient()
	clusterClient.rancherClient = rancherClient
	result, err := newClusterCatalogClient(
		name,
		clusterClient,
		logrus.New().WithFields(logrus.Fields{}),
	)
	assert.Ok(t, err)
	clusterCatalogClientResult := result.(*clusterCatalogClient)
	clusterCatalogClientResult.catalog = clusterCatalogData
	return clusterCatalogClientResult
}
