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
	"path/filepath"
	"testing"

	"github.com/bitgrip/cattlectl/internal/pkg/assert"
	"github.com/bitgrip/cattlectl/internal/pkg/rancher/project/model"
	"github.com/bitgrip/cattlectl/internal/pkg/template"
	"github.com/sirupsen/logrus"
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

		projectFile := fmt.Sprintf("testdata/input/%s", testCase["file"])
		fileContent, err := ioutil.ReadFile(projectFile)
		assert.Ok(t, err)
		values := make(map[string]interface{})
		values["project_name"] = "test-project"
		values["stage"] = "test-stage"
		projectData, err := template.BuildTemplate(fileContent, values, filepath.Dir(projectFile), false)
		assert.Ok(t, err)

		parser := NewProjectParser(projectFile, projectData, &model.Project{}, values)

		//Act
		err = parser.Parse()

		//Assert
		assert.NotOk(t, err, fmt.Sprintf("Invalid descriptor: %s", testCase[kindKey]))
		assert.Equals(t, fmt.Sprintf("Invalid descriptor: %s", testCase[kindKey]), err.Error())
	}
}

func TestParseValidProjectDescriptor(t *testing.T) {
	testName := "valid-project"
	//Arrange
	projectFile := "testdata/valid-project/project.yaml"
	fileContent, err := ioutil.ReadFile(projectFile)
	assert.Ok(t, err)
	values := make(map[string]interface{})
	if data, err := ioutil.ReadFile(fmt.Sprintf("testdata/%s/values.yaml", testName)); err != nil {
		log.Fatal(err)
	} else if err := yaml.Unmarshal(data, &values); err != nil {
		log.Fatal(err)
	}
	projectData, err := template.BuildTemplate(fileContent, values, filepath.Dir(projectFile), false)
	assert.Ok(t, err)
	project := model.Project{}
	parser := NewProjectParser(projectFile, projectData, &project, values)

	//Act
	err = parser.Parse()

	//Verify
	assert.Ok(t, err)
	actual, err := yaml.Marshal(project)
	assert.Ok(t, err)
	assert.AssertGoldenFile(t, testName, actual)
}

func TestWithGoldenFile(t *testing.T) {
	tests := []string{
		"simple-include",
		"cycle-include",
	}
	for _, test := range tests {
		runTestWithGoldenFile(t, test)
	}
}

func runTestWithGoldenFile(t *testing.T, testName string) {
	//Arrange
	projectFile := fmt.Sprintf("testdata/%s/project.yaml", testName)
	fileContent, err := ioutil.ReadFile(projectFile)
	assert.Ok(t, err)
	values := make(map[string]interface{})
	if data, err := ioutil.ReadFile(fmt.Sprintf("testdata/%s/values.yaml", testName)); err != nil {
		logrus.Info("No values found - use empty")
	} else if err := yaml.Unmarshal(data, &values); err != nil {
		log.Fatal(err)
	}
	projectData, err := template.BuildTemplate(fileContent, values, filepath.Dir(projectFile), false)
	assert.Ok(t, err)
	project := model.Project{}
	parser := NewProjectParser(projectFile, projectData, &project, values)

	//Act
	err = parser.Parse()

	//Verify
	assert.Ok(t, err)
	actual, err := yaml.Marshal(project)
	assert.Ok(t, err)
	assert.AssertGoldenFile(t, testName, actual)
}

func readTestdata(t *testing.T, testdataFile string) model.Project {
	project := model.Project{}
	fileContent, err := ioutil.ReadFile("testdata/" + testdataFile)
	assert.Ok(t, err)

	err = yaml.Unmarshal(fileContent, &project)
	assert.Ok(t, err)
	return project
}
