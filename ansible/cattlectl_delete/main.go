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

	"github.com/bitgrip/cattlectl/ansible/utils"
	"github.com/bitgrip/cattlectl/internal/pkg/ctl"
)

type moduleArgs struct {
	ProjectName      string   `json:"project_name"`
	Namespace        string   `json:"namespace"`
	Kind             string   `json:"kind"`
	Names            []string `json:"names"`
	utils.AccessArgs `json:",inline"`
}

type listResponse struct {
	Deleted            []string `json:"deleted"`
	utils.BaseResponse `json:",inline"`
}

func main() {
	var moduleArgs moduleArgs
	utils.ReadArguments(&moduleArgs)

	var response listResponse
	response.Version = ctl.Version
	for _, name := range moduleArgs.Names {
		deleted, err := ctl.DeleteProjectResouce(
			moduleArgs.ProjectName,
			moduleArgs.Namespace,
			moduleArgs.Kind,
			name,
			utils.BuildRancherConfig(moduleArgs.AccessArgs),
		)
		if err != nil {
			response.Msg = fmt.Sprintf("Failed to delete %s %s:  - %v", moduleArgs.Kind, name, err)
			response.Failed = true
			utils.FailJson(response)
		}
		if deleted {
			response.Deleted = append(response.Deleted, name)
			response.Changed = true
		}
	}
	response.Msg = fmt.Sprintf("Deleted %v %ss", len(response.Deleted), moduleArgs.Kind)
	utils.ExitJson(response)
}
