// Copyright © 2020 Bitgrip <berlin@bitgrip.de>
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

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/bitgrip/cattlectl/ansible/utils"
	"github.com/bitgrip/cattlectl/internal/pkg/ctl"
	"github.com/bitgrip/cattlectl/internal/pkg/rancher/descriptor"
	"github.com/bitgrip/cattlectl/internal/pkg/template"
)

type moduleArgs struct {
	ApplyFile        string                 `json:"file"`
	ValueFiles       []string               `json:"value_files"`
	Values           map[string]interface{} `json:"values"`
	WorkingDirectory string                 `json:"working_directory"`
	utils.AccessArgs `json:",inline"`
}

type listResponse struct {
	ApplyResult        descriptor.ConvergeResult `json:"apply_result"`
	utils.BaseResponse `json:",inline"`
}

func main() {
	var moduleArgs moduleArgs
	utils.ReadArguments(&moduleArgs)

	var response listResponse
	response.Version = ctl.Version

	if moduleArgs.WorkingDirectory != "" {
		err := os.Chdir(moduleArgs.WorkingDirectory)
		if err != nil {
			response.Msg = fmt.Sprintf("Failed to apply file %s:  - %v", moduleArgs.ApplyFile, err)
			response.Failed = true
			utils.FailJson(response)
		}
	}

	values, err := utils.LoadValues(moduleArgs.Values, moduleArgs.ValueFiles...)
	if err != nil {
		response.Msg = fmt.Sprintf("Failed to apply file %s:  - %v", moduleArgs.ApplyFile, err)
		response.Failed = true
		utils.FailJson(response)
	}
	fileContent, err := ioutil.ReadFile(moduleArgs.ApplyFile)
	if err != nil {
		response.Msg = fmt.Sprintf("Failed to apply file %s:  - %v", moduleArgs.ApplyFile, err)
		response.Failed = true
		utils.FailJson(response)
	}
	projectData, err := template.BuildTemplate(fileContent, values, filepath.Dir(moduleArgs.ApplyFile), false)
	if err != nil {
		response.Msg = fmt.Sprintf("Failed to apply file %s:  - %v", moduleArgs.ApplyFile, err)
		response.Failed = true
		utils.FailJson(response)
	}

	result, err := ctl.ApplyDescriptor(
		moduleArgs.ApplyFile,
		projectData,
		map[string]interface{}{},
		utils.BuildRancherConfig(moduleArgs.AccessArgs),
	)

	if err != nil {
		response.Msg = "Failed to apply descriptor: " + err.Error()
		response.Failed = true
		utils.FailJson(response)
	}
	response.ApplyResult = result
	response.Changed = len(result.CreatedResources) > 0 || len(result.UpgradedResources) > 0
	utils.ExitJson(response)
}
