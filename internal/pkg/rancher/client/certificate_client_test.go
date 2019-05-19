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
	projectModel "github.com/bitgrip/cattlectl/internal/pkg/rancher/project/model"
	"github.com/bitgrip/cattlectl/internal/pkg/rancher/stubs"
	"github.com/rancher/norman/types"
	backendProjectClient "github.com/rancher/types/client/project/v3"
	"github.com/sirupsen/logrus"
)

func Test_certificateClient_Exists(t *testing.T) {
	tests := []struct {
		name      string
		client    *certificateClient
		wanted    bool
		wantErr   bool
		wantedErr string
	}{
		{
			name: "Existing",
			client: existingCertificateClient(
				t,
				&types.ListOpts{
					Filters: map[string]interface{}{
						"name":        "existing-certificate",
						"namespaceId": "test-namespace-id",
					},
				},
			),
			wanted:  true,
			wantErr: false,
		},
		{
			name: "Not_Existing",
			client: notExistingCertificateClient(
				t,
				&types.ListOpts{
					Filters: map[string]interface{}{
						"name":        "existing-certificate",
						"namespaceId": "test-namespace-id",
					},
				},
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

func Test_certificateClient_Create(t *testing.T) {
	tests := []struct {
		name      string
		client    *certificateClient
		wantErr   bool
		wantedErr string
	}{
		{
			name: "Create",
			client: notExistingCertificateClient(
				t,
				nil,
			),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.client.Create()
			if tt.wantErr {
				assert.NotOk(t, err, tt.wantedErr)
			} else {
				assert.Ok(t, err)
			}
		})
	}
}

func existingCertificateClient(t *testing.T, expectedListOpts *types.ListOpts) *certificateClient {
	const (
		projectID       = "test-project-id"
		projectName     = "test-project-name"
		namespaceID     = "test-namespace-id"
		namespace       = "test-namespace"
		clusterID       = "test-cluster-id"
		certificateName = "test-certificate"
	)
	var (
		certificate = projectModel.Certificate{}
		testClients = stubs.CreateBackendStubs(t)
	)
	certificate.Name = certificateName

	certificateOperationsStub := stubs.CreateCertificateOperationsStub(t)
	certificateOperationsStub.DoList = func(opts *types.ListOpts) (*backendProjectClient.CertificateCollection, error) {
		if !reflect.DeepEqual(expectedListOpts, opts) {
			return nil, fmt.Errorf("Unexpected ListOpts %v", opts)
		}
		return &backendProjectClient.CertificateCollection{
			Data: []backendProjectClient.Certificate{
				backendProjectClient.Certificate{
					Name:        "existing-certificate",
					NamespaceId: "test-namespace-id",
				},
			},
		}, nil
	}
	testClients.ProjectClient.Certificate = certificateOperationsStub
	result, err := newCertificateClient(
		"existing-certificate",
		"test-namespace",
		&projectClient{
			resourceClient: resourceClient{
				name: projectName,
				id:   projectID,
			},
		},
		testClients.ProjectClient,
		logrus.New().WithFields(logrus.Fields{}),
	)
	assert.Ok(t, err)
	certificateClientResult := result.(*certificateClient)
	certificateClientResult.namespaceID = "test-namespace-id"
	return certificateClientResult
}

func notExistingCertificateClient(t *testing.T, expectedListOpts *types.ListOpts) *certificateClient {
	const (
		projectID       = "test-project-id"
		projectName     = "test-project-name"
		namespaceID     = "test-namespace-id"
		namespace       = "test-namespace"
		clusterID       = "test-cluster-id"
		certificateName = "test-certificate"
	)
	var (
		certificate = projectModel.Certificate{}
		testClients = stubs.CreateBackendStubs(t)
	)
	certificate.Name = certificateName

	certificateOperationsStub := stubs.CreateCertificateOperationsStub(t)
	certificateOperationsStub.DoList = func(opts *types.ListOpts) (*backendProjectClient.CertificateCollection, error) {
		if !reflect.DeepEqual(expectedListOpts, opts) {
			return nil, fmt.Errorf("Unexpected ListOpts %v", opts)
		}
		return &backendProjectClient.CertificateCollection{
			Data: []backendProjectClient.Certificate{},
		}, nil
	}
	certificateOperationsStub.DoCreate = func(certificate *backendProjectClient.Certificate) (*backendProjectClient.Certificate, error) {
		return certificate, nil
	}
	testClients.ProjectClient.Certificate = certificateOperationsStub
	result, err := newCertificateClient(
		"existing-certificate",
		"test-namespace",
		&projectClient{
			resourceClient: resourceClient{
				name: projectName,
				id:   projectID,
			},
		},
		testClients.ProjectClient,
		logrus.New().WithFields(logrus.Fields{}),
	)
	assert.Ok(t, err)
	certificateClientResult := result.(*certificateClient)
	certificateClientResult.namespaceID = "test-namespace-id"
	return certificateClientResult
}
