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

	"github.com/rancher/types/client/project/v3"
	yaml "gopkg.in/yaml.v2"
)

func TestWorkload(t *testing.T) {
	// Read cattlectl model descriptor from file
	fileContent, err := ioutil.ReadFile("testdata/input/job.yaml")
	assert.Ok(t, err)
	jobDescriptor := JobDescriptor{}
	err = yaml.Unmarshal(fileContent, &jobDescriptor)
	assert.Ok(t, err)

	transferContent, err := yaml.Marshal(jobDescriptor)
	assert.Ok(t, err)
	assert.AssertGoldenFile(t, "descriptor.job", transferContent)

	// convert cattlectl model spec to rancher model
	transferContent, err = json.Marshal(jobDescriptor.Spec)
	assert.Ok(t, err)
	rancherJob := client.Job{}
	err = json.Unmarshal(transferContent, &rancherJob)
	assert.Ok(t, err)

	// convert rancher model to cattlectl model
	transferContent, err = json.Marshal(rancherJob)
	assert.Ok(t, err)
	modelJob := Job{}
	err = json.Unmarshal(transferContent, &modelJob)
	assert.Ok(t, err)

	assert.Equals(t, jobDescriptor.Spec, modelJob)

}
