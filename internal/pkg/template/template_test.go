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

package template

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/bitgrip/cattlectl/internal/pkg/assert"
)

func TestGoldenFileTestCases(t *testing.T) {
	testdirs, err := ioutil.ReadDir("testdata")
	assert.Ok(t, err)
	for _, testdir := range testdirs {
		if !testdir.IsDir() {
			continue
		}
		runTestCaseVariants(testdir.Name(), t)
	}
}

func runTestCaseVariants(testName string, t *testing.T) {
	defer testChdir(t, "testdata/"+testName)()
	runTestCase(testName, false, t)
	runTestCase(testName, true, t)
}

func runTestCase(testName string, truncated bool, t *testing.T) {

	templateData, err := ioutil.ReadFile("input.txt")
	assert.Ok(t, err)

	// test structs for testing toYaml with structs
	type testInnerStruct struct {
		Key1 string
		Key2 string
	}
	type testStruct struct {
		Foo  testInnerStruct
		Key3 string
	}
	values := map[string]interface{}{
		"key1":            "value1",
		"file_name1":      "read-file1.txt",
		"my_array_value":  []string{"a1", "a2", "a3"},
		"my_struct_value": testStruct{testInnerStruct{"value1", "value2"}, "value3"},
	}
	actual, err := BuildTemplate(templateData, values, ".", truncated)
	assert.Ok(t, err)
	golden := fmt.Sprintf("%s-%v.golden", testName, truncated)
	if *assert.Update {
		ioutil.WriteFile(golden, actual, 0644)
	}
	expected, err := ioutil.ReadFile(golden)
	assert.Ok(t, err)
	assert.Equals(t, string(expected), string(actual))
}

func testChdir(t *testing.T, dir string) func() {
	old, err := os.Getwd()
	assert.Ok(t, err)
	err = os.Chdir(dir)
	assert.Ok(t, err)
	return func() {
		if err := os.Chdir(old); err != nil {
			t.Fatalf("err: %s", err)
		}
	}
}

func TestErrorOnParsingClusterDescriptor(t *testing.T) {
	type testCase struct {
		value          interface{}
		expectedResult string
	}

	testCases := []testCase{
		testCase{[]byte("admin"), "YWRtaW4="},
		testCase{"admin", "YWRtaW4="},
	}
	for _, test := range testCases {
		//Act
		result := toBase64(test.value)

		//Assert
		assert.Equals(t, test.expectedResult, result)
	}
}

//func TestErrorWhenYamlDataIsWrong(t *testing.T) {
//	invalidYaml, err := ioutil.ReadFile("testdata/invalid_yaml.yaml.zip")
//Act
//	result, err := toYaml(invalidYaml)

//Assert
//	assert.NotOk(t, err, "")
//	assert.Equals(t, "", result)
//}
