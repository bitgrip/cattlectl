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

func TestStatefulSet(t *testing.T) {
	// Read cattlectl model descriptor from file
	fileContent, err := ioutil.ReadFile("testdata/input/statefulSet.yaml")
	assert.Ok(t, err)
	statefulSetDescriptor := StatefulSetDescriptor{}
	err = yaml.Unmarshal(fileContent, &statefulSetDescriptor)
	assert.Ok(t, err)

	transferContent, err := yaml.Marshal(statefulSetDescriptor)
	assert.Ok(t, err)
	assert.AssertGoldenFile(t, "descriptor.statefulSet", transferContent)

	// convert cattlectl model to rancher model
	rancherStatefulSet, err := ConvertStatefulSetDescriptorToProjectAPI(statefulSetDescriptor)
	assert.Ok(t, err)

	// convert rancher model to cattlectl model
	transferContent, err = json.Marshal(rancherStatefulSet)
	assert.Ok(t, err)
	modelStatefulSet := StatefulSet{}
	err = json.Unmarshal(transferContent, &modelStatefulSet)
	assert.Ok(t, err)

	// Rancher model dose not know StorageClass and uses StorageClassID instead.
	// Remove it from the expected side.
	statefulSetDescriptor.Spec.StatefulSetConfig.VolumeClaimTemplates[0].StorageClass = ""
	assert.Equals(t, statefulSetDescriptor.Spec, modelStatefulSet)

}
