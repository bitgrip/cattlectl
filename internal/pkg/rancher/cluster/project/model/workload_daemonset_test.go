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

package model

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/bitgrip/cattlectl/internal/pkg/assert"

	yaml "gopkg.in/yaml.v2"
)

func TestDaemonSet(t *testing.T) {
	// Read cattlectl model descriptor from file
	fileContent, err := ioutil.ReadFile("testdata/input/daemonSet.yaml")
	assert.Ok(t, err)
	daemonSetDescriptor := DaemonSetDescriptor{}
	err = yaml.Unmarshal(fileContent, &daemonSetDescriptor)
	assert.Ok(t, err)

	transferContent, err := yaml.Marshal(daemonSetDescriptor)
	assert.Ok(t, err)
	assert.AssertGoldenFile(t, "descriptor.daemonSet", transferContent)

	// convert cattlectl model to rancher model
	rancherJob, err := ConvertDaemonSetDescriptorToProjectAPI(daemonSetDescriptor)
	assert.Ok(t, err)

	// convert rancher model to cattlectl model
	transferContent, err = json.Marshal(rancherJob)
	assert.Ok(t, err)
	modelJob := DaemonSet{}
	err = json.Unmarshal(transferContent, &modelJob)
	assert.Ok(t, err)

	assert.Equals(t, daemonSetDescriptor.Spec, modelJob)

}
