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

package utils

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/bitgrip/cattlectl/internal/pkg/assert"
	yaml "gopkg.in/yaml.v2"
)

func TestPlainValuesStructure(t *testing.T) {
	const valuesFile = "testdata/values.yaml"

	expected := make(map[string]interface{}, 0)
	fileContent, err := ioutil.ReadFile(valuesFile)
	assert.Ok(t, err)
	err = yaml.Unmarshal(fileContent, &expected)
	assert.Ok(t, err)

	actual, err := LoadValues(valuesFile)
	assert.Ok(t, err)
	assert.Equals(t, expected, actual)
}

func TestChangeValueByEnvironmentVariable(t *testing.T) {
	const valuesFile = "testdata/values.yaml"
	os.Setenv("KEY1", "altered-by-env")

	expected := make(map[string]interface{}, 0)
	fileContent, err := ioutil.ReadFile(valuesFile)
	assert.Ok(t, err)
	err = yaml.Unmarshal(fileContent, &expected)
	assert.Ok(t, err)
	expected["key1"] = "altered-by-env"

	actual, err := LoadValues(valuesFile)
	assert.Ok(t, err)
	assert.Equals(t, expected, actual)
}

func TestFailWhenValuesKeysAreMissformated(t *testing.T) {
	const valuesFile = "testdata/values-with-camel-case.yaml"

	_, err := LoadValues(valuesFile)
	assert.NotOk(t, err, "uppercase characters are not allowed on value keys")
}

func TestBla(t *testing.T) {
	const valuesFile = "testdata/values-with-structure.yaml"

	_, err := LoadValues(valuesFile)
	assert.Ok(t, err)
}

func TestNotExistingValuesFile(t *testing.T) {
	const valuesFile = "testdata/not-existing-values.yaml"

	expected := make(map[string]interface{}, 0)

	actual, err := LoadValues(valuesFile)
	assert.Ok(t, err)
	assert.Equals(t, expected, actual)
}
