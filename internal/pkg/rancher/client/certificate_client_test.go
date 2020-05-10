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
	projectModel "github.com/bitgrip/cattlectl/internal/pkg/rancher/cluster/project/model"
	"github.com/bitgrip/cattlectl/internal/pkg/rancher/stubs"
	"github.com/rancher/norman/types"
	backendProjectClient "github.com/rancher/types/client/project/v3"
	"github.com/sirupsen/logrus"
)

const (
	simpleCertificateName = "simple-certificate-name"
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
						"name": simpleCertificateName,
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
						"name": simpleCertificateName,
					},
				},
			),
			wanted:  false,
			wantErr: false,
		},
		{
			name: "Existing_Namespaced",
			client: existingNamespacedCertificateClient(
				t,
				&types.ListOpts{
					Filters: map[string]interface{}{
						"name":        simpleCertificateName,
						"namespaceId": simpleNamespaceID,
					},
				},
			),
			wanted:  true,
			wantErr: false,
		},
		{
			name: "Not_Existing_Namespaced",
			client: notExistingNamespacedCertificateClient(
				t,
				&types.ListOpts{
					Filters: map[string]interface{}{
						"name":        simpleCertificateName,
						"namespaceId": simpleNamespaceID,
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
				&types.ListOpts{
					Filters: map[string]interface{}{
						"name":   simpleCertificateName,
						"system": false,
					},
				},
			),
			wantErr: false,
		},
		{
			name: "Create_namespaces",
			client: notExistingNamespacedCertificateClient(
				t,
				&types.ListOpts{
					Filters: map[string]interface{}{
						"name":        simpleCertificateName,
						"namespaceId": simpleNamespaceID,
					},
				},
			),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.client.Create(false)
			if tt.wantErr {
				assert.NotOk(t, err, tt.wantedErr)
			} else {
				assert.Ok(t, err)
			}
		})
	}
}

func Test_certificateClient_Upgrade(t *testing.T) {
	tests := []struct {
		name      string
		client    *certificateClient
		wantErr   bool
		wantedErr string
	}{
		{
			name: "global",
			client: existingCertificateClient(
				t,
				&types.ListOpts{
					Filters: map[string]interface{}{
						"name": simpleCertificateName,
					},
				},
			),
			wantErr: false,
		},
		{
			name: "namespaced",
			client: existingNamespacedCertificateClient(
				t,
				&types.ListOpts{
					Filters: map[string]interface{}{
						"name":        simpleCertificateName,
						"namespaceId": simpleNamespaceID,
					},
				},
			),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.client.Upgrade(false)
			if tt.wantErr {
				assert.NotOk(t, err, tt.wantedErr)
			} else {
				assert.Ok(t, err)
			}
		})
	}
}

func existingCertificateClient(t *testing.T, expectedListOpts *types.ListOpts) *certificateClient {
	var (
		certificate = projectModel.Certificate{}
		testClients = stubs.CreateBackendStubs(t)
	)
	certificate.Name = simpleCertificateName

	certificateOperationsStub := stubs.CreateCertificateOperationsStub(t)
	certificateOperationsStub.DoList = func(opts *types.ListOpts) (*backendProjectClient.CertificateCollection, error) {
		if !reflect.DeepEqual(expectedListOpts, opts) {
			return nil, fmt.Errorf("Unexpected ListOpts\n%v\n%v", opts, expectedListOpts)
		}
		return &backendProjectClient.CertificateCollection{
			Data: []backendProjectClient.Certificate{
				backendProjectClient.Certificate{
					Name: simpleCertificateName,
				},
			},
		}, nil
	}
	testClients.ProjectClient.Certificate = certificateOperationsStub
	result, err := newCertificateClient(
		simpleCertificateName,
		"",
		&projectClient{
			resourceClient: resourceClient{
				name: simpleProjectName,
				id:   simpleProjectID,
			},
			_backendProjectClient: testClients.ProjectClient,
		},
		logrus.New().WithFields(logrus.Fields{}),
	)
	assert.Ok(t, err)
	certificateClientResult := result.(*certificateClient)
	certificateClientResult.namespaceID = simpleNamespaceID
	return certificateClientResult
}

func notExistingCertificateClient(t *testing.T, expectedListOpts *types.ListOpts) *certificateClient {
	var (
		certificate = projectModel.Certificate{}
		testClients = stubs.CreateBackendStubs(t)
	)
	certificate.Name = simpleCertificateName

	certificateOperationsStub := stubs.CreateCertificateOperationsStub(t)
	certificateOperationsStub.DoList = func(opts *types.ListOpts) (*backendProjectClient.CertificateCollection, error) {
		if !reflect.DeepEqual(expectedListOpts, opts) {
			return nil, fmt.Errorf("Unexpected ListOpts\n%v\n%v", opts, expectedListOpts)
		}
		return &backendProjectClient.CertificateCollection{
			Data: []backendProjectClient.Certificate{},
		}, nil
	}
	certificateOperationsStub.DoCreate = func(certificate *backendProjectClient.Certificate) (*backendProjectClient.Certificate, error) {
		if certificate.NamespaceId != "" {
			return nil, fmt.Errorf("Unexpected NamespaceID %v", certificate.NamespaceId)
		}
		return certificate, nil
	}
	testClients.ProjectClient.Certificate = certificateOperationsStub
	result, err := newCertificateClient(
		simpleCertificateName,
		"",
		&projectClient{
			resourceClient: resourceClient{
				name: simpleProjectName,
				id:   simpleProjectID,
			},
			_backendProjectClient: testClients.ProjectClient,
		},
		logrus.New().WithFields(logrus.Fields{}),
	)
	assert.Ok(t, err)
	certificateClientResult := result.(*certificateClient)
	certificateClientResult.namespaceID = simpleNamespaceID
	return certificateClientResult
}

func existingNamespacedCertificateClient(t *testing.T, expectedListOpts *types.ListOpts) *certificateClient {
	var (
		certificate = projectModel.Certificate{}
		testClients = stubs.CreateBackendStubs(t)
	)
	certificate.Name = simpleCertificateName

	certificateOperationsStub := stubs.CreateNamespacedCertificateOperationsStub(t)
	certificateOperationsStub.DoList = func(opts *types.ListOpts) (*backendProjectClient.NamespacedCertificateCollection, error) {
		if !reflect.DeepEqual(expectedListOpts, opts) {
			return nil, fmt.Errorf("Unexpected ListOpts\n%v\n%v", opts, expectedListOpts)
		}
		return &backendProjectClient.NamespacedCertificateCollection{
			Data: []backendProjectClient.NamespacedCertificate{
				backendProjectClient.NamespacedCertificate{
					Name:        simpleCertificateName,
					NamespaceId: simpleNamespaceID,
				},
			},
		}, nil
	}
	testClients.ProjectClient.NamespacedCertificate = certificateOperationsStub
	result, err := newCertificateClient(
		simpleCertificateName,
		simpleNamespaceName,
		&projectClient{
			resourceClient: resourceClient{
				name: simpleProjectName,
				id:   simpleProjectID,
			},
			_backendProjectClient: testClients.ProjectClient,
		},
		logrus.New().WithFields(logrus.Fields{}),
	)
	assert.Ok(t, err)
	certificateClientResult := result.(*certificateClient)
	certificateClientResult.namespaceID = simpleNamespaceID
	return certificateClientResult
}

func notExistingNamespacedCertificateClient(t *testing.T, expectedListOpts *types.ListOpts) *certificateClient {
	var (
		certificate = projectModel.Certificate{}
		testClients = stubs.CreateBackendStubs(t)
	)
	certificate.Name = simpleCertificateName

	certificateOperationsStub := stubs.CreateNamespacedCertificateOperationsStub(t)
	certificateOperationsStub.DoList = func(opts *types.ListOpts) (*backendProjectClient.NamespacedCertificateCollection, error) {
		if !reflect.DeepEqual(expectedListOpts, opts) {
			return nil, fmt.Errorf("Unexpected ListOpts %v", opts)
		}
		return &backendProjectClient.NamespacedCertificateCollection{
			Data: []backendProjectClient.NamespacedCertificate{},
		}, nil
	}
	certificateOperationsStub.DoCreate = func(certificate *backendProjectClient.NamespacedCertificate) (*backendProjectClient.NamespacedCertificate, error) {
		if certificate.NamespaceId != simpleNamespaceID {
			return nil, fmt.Errorf("Unexpected NamespaceID %v", certificate.NamespaceId)
		}
		return certificate, nil
	}
	testClients.ProjectClient.NamespacedCertificate = certificateOperationsStub
	result, err := newCertificateClient(
		simpleCertificateName,
		simpleNamespaceName,
		&projectClient{
			resourceClient: resourceClient{
				name: simpleProjectName,
				id:   simpleProjectID,
			},
			_backendProjectClient: testClients.ProjectClient,
		},
		logrus.New().WithFields(logrus.Fields{}),
	)
	assert.Ok(t, err)
	certificateClientResult := result.(*certificateClient)
	certificateClientResult.namespaceID = simpleNamespaceID
	return certificateClientResult
}
