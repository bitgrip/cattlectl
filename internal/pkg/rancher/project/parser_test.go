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
	"fmt"
	"io/ioutil"
	"log"
	"testing"

	"github.com/bitgrip/cattlectl/internal/pkg/assert"
	yaml "gopkg.in/yaml.v2"
)

func TestErrorOnParsingClusterDescriptor(t *testing.T) {
	const (
		fileKey = "file"
		kindKey = "expected-kind"
	)
	var cases = []map[string]string{
		{
			fileKey: "cluster.yaml",
			kindKey: "Cluster",
		},
		{
			fileKey: "missing-kind.yaml",
			kindKey: "<nil>",
		},
	}
	for _, testCase := range cases {
		parser := NewParser(fmt.Sprintf("testdata/input/%s", testCase["file"]))
		values := make(map[string]interface{})
		values["project_name"] = "test-project"
		values["stage"] = "test-stage"

		//Act
		_, err := parser.Parse(values)

		//Assert
		assert.NotOk(t, err, fmt.Sprintf("Invalid descriptor: %s", testCase[kindKey]))
		assert.Equals(t, fmt.Sprintf("Invalid descriptor: %s", testCase[kindKey]), err.Error())
	}
}

func TestParseValidProjectDescriptor(t *testing.T) {
	//Arrange
	parser := NewParser("testdata/valid-project/project.yaml")

	values := make(map[string]interface{})
	if data, err := ioutil.ReadFile("testdata/valid-project/values.yaml"); err != nil {
		log.Fatal(err)
	} else if err := yaml.Unmarshal(data, &values); err != nil {
		log.Fatal(err)
	}

	//Act
	project, err := parser.Parse(values)

	//Verify
	assert.Ok(t, err)
	assert.Equals(t, "v1.0", project.APIVersion)
	assert.Equals(t, "Project", project.Kind)
	assert.Equals(t, "https://ui.rancher.server", project.Metadata.RancherURL)
	assert.Equals(t, "token-12345", project.Metadata.AccessKey)
	assert.Equals(t, "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", project.Metadata.SecretKey)
	assert.Equals(t, "token-12345:xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", project.Metadata.TokenKey)
	assert.Equals(t, "my-wordpress-dev", project.Metadata.Name)
	assert.Equals(t, "j-4444", project.Metadata.ClusterID)
	assert.Equals(t, 1, len(project.Namespaces))
	assert.Equals(t, "my-wordpress-dev-web", project.Namespaces[0].Name)
	assert.Equals(t, 1, len(project.StorageClasses))
	assert.Equals(t, "my-wordpress-dev-local-mariadb", project.StorageClasses[0].Name)
	assert.Equals(t, "kubernetes.io/no-provisioner", project.StorageClasses[0].Provisioner)
	assert.Equals(t, "Delete", project.StorageClasses[0].ReclaimPolicy)
	assert.Equals(t, true, project.StorageClasses[0].CreatePersistentVolumes)
	assert.Equals(t, "WaitForFirstConsumer", project.StorageClasses[0].VolumeBindMode)
	assert.Equals(t, 1, len(project.StorageClasses[0].PersistentVolumeGroups))
	assert.Equals(t, "my-wordpress-dev-mariadb", project.StorageClasses[0].PersistentVolumeGroups[0].Name)
	assert.Equals(t, "local", project.StorageClasses[0].PersistentVolumeGroups[0].Type)
	assert.Equals(t, []string{"ReadWriteOnce"}, project.StorageClasses[0].PersistentVolumeGroups[0].AccessModes)
	assert.Equals(t, "3Gi", project.StorageClasses[0].PersistentVolumeGroups[0].Capacity)
	assert.Equals(t, "/var/data/my-wordpress-dev-mariadb", project.StorageClasses[0].PersistentVolumeGroups[0].Path)
	assert.Equals(t, []string{"node-1", "node-2", "node-3"}, project.StorageClasses[0].PersistentVolumeGroups[0].Nodes)
	assert.Equals(t, "ssh ${node} sudo mkdir -p ${path}", project.StorageClasses[0].PersistentVolumeGroups[0].CreateScript)
}
