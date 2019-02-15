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
	testName := "valid-project"
	//Arrange
	parser := NewParser("testdata/valid-project/project.yaml")

	values := make(map[string]interface{})
	if data, err := ioutil.ReadFile(fmt.Sprintf("testdata/%s/values.yaml", testName)); err != nil {
		log.Fatal(err)
	} else if err := yaml.Unmarshal(data, &values); err != nil {
		log.Fatal(err)
	}

	//Act
	project, err := parser.Parse(values)

	//Verify
	assert.Ok(t, err)
	actual, err := yaml.Marshal(project)
	assert.Ok(t, err)
	assert.AssertGoldenFile(t, testName, actual)
}
