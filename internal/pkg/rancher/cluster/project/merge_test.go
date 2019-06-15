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

package project

import (
	"io/ioutil"
	"testing"

	"github.com/bitgrip/cattlectl/internal/pkg/assert"
	projectModel "github.com/bitgrip/cattlectl/internal/pkg/rancher/cluster/project/model"
	yaml "gopkg.in/yaml.v2"
)

func TestMergeProjects(t *testing.T) {
	testName := "simple-include"
	child := readTestdataProject(t, "include/simple/child.yaml")
	parent := readTestdataProject(t, "include/simple/parent.yaml")
	err := MergeProject(child, &parent)
	assert.Ok(t, err)

	actual, err := yaml.Marshal(parent)
	assert.Ok(t, err)
	assert.AssertGoldenFile(t, testName, actual)
}

func readTestdataProject(t *testing.T, testdataFile string) projectModel.Project {
	project := projectModel.Project{}
	fileContent, err := ioutil.ReadFile("testdata/" + testdataFile)
	assert.Ok(t, err)

	err = yaml.Unmarshal(fileContent, &project)
	assert.Ok(t, err)
	return project
}
