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

package project

import (
	"testing"

	"github.com/bitgrip/cattlectl/internal/pkg/assert"
	"github.com/bitgrip/cattlectl/internal/pkg/rancher/clientstub"
)

func expectSetClusterForExisting(expectedClusterName, expectedClusterID string, client *clientstub.ClientStub, t *testing.T) {
	client.DoHasClusterWithName = func(clusterName string) (bool, string, error) {
		assert.Equals(t, expectedClusterName, clusterName)
		return true, expectedClusterID, nil
	}
	client.DoSetCluster = func(clusterName, clusterID string) error {
		assert.Equals(t, expectedClusterName, clusterName)
		assert.Equals(t, expectedClusterID, clusterID)
		return nil
	}
}

func expectSetProjectForExisting(expectedProjectName, expectedProjectID string, client *clientstub.ClientStub, t *testing.T) {
	client.DoHasProjectWithName = func(projectName string) (bool, string, error) {
		assert.Equals(t, expectedProjectName, projectName)
		return true, expectedProjectID, nil
	}
	client.DoSetProject = func(projectName, projectID string) error {
		assert.Equals(t, expectedProjectName, projectName)
		assert.Equals(t, expectedProjectID, projectID)
		return nil
	}
}
