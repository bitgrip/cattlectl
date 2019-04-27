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

package rancher

import (
	"fmt"
	"testing"

	"github.com/bitgrip/cattlectl/internal/pkg/assert"
	projectModel "github.com/bitgrip/cattlectl/internal/pkg/rancher/project/model"
	"github.com/bitgrip/cattlectl/internal/pkg/rancher/stubs"
	"github.com/rancher/norman/types"
	projectClient "github.com/rancher/types/client/project/v3"
	"github.com/sirupsen/logrus"
)

func TestHasDockerCredential_DockerCredentialExisting(t *testing.T) {
	const (
		dockerCredentialID   = "test-dockerCredential-id"
		dockerCredentialName = "test-dockerCredential-name"
	)
	var (
		clientConfig     = ClientConfig{}
		dockerCredential = projectModel.DockerCredential{
			Name: dockerCredentialName,
		}
		testClients           = stubs.CreateBackendStubs(t)
		dockerCredentialCache = make(map[string]projectClient.DockerCredential)
	)

	dockerCredentialOperationsStub := stubs.CreateDockerCredentialOperationsStub(t)
	dockerCredentialOperationsStub.DoList = func(opts *types.ListOpts) (*projectClient.DockerCredentialCollection, error) {
		assert.Equals(t, map[string]interface{}{
			"system": "false",
			"name":   dockerCredentialName,
		}, opts.Filters)
		return &projectClient.DockerCredentialCollection{
			Data: []projectClient.DockerCredential{
				projectClient.DockerCredential{
					Name:     fmt.Sprint(opts.Filters["name"]),
					Resource: types.Resource{ID: dockerCredentialID},
				},
			},
		}, nil
	}
	testClients.ProjectClient.DockerCredential = dockerCredentialOperationsStub

	client := rancherClient{
		clientConfig:          clientConfig,
		projectClient:         testClients.ProjectClient,
		dockerCredentialCache: dockerCredentialCache,
		logger:                logrus.WithField("test", true),
	}

	//Act
	result, err := client.HasDockerCredential(dockerCredential)

	//Assert
	assert.Ok(t, err)
	assert.Equals(t, true, result)
	assert.Equals(t, dockerCredentialID, dockerCredentialCache[dockerCredentialName].ID)

}
