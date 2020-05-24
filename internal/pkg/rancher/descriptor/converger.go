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

package descriptor

import (
	"github.com/bitgrip/cattlectl/internal/pkg/rancher/client"
)

type Converger interface {
	Converge(bool) (ConvergeResult, error)
}

type ConvergeResult struct {
	CreatedResources  []ResourceDescriptor `json:"created_resources"`
	UpgradedResources []ResourceDescriptor `json:"upgraded_resources"`
}

type ResourceDescriptor struct {
	Type string
	Name string
}

type ResourceClientConverger struct {
	Client   client.ResourceClient
	Children []Converger
}

func (converger *ResourceClientConverger) Converge(dryRun bool) (result ConvergeResult, err error) {
	var (
		name    string
		exists  bool
		changed bool
	)
	if name, err = converger.Client.Name(); err != nil {
		return
	}
	if exists, err = converger.Client.Exists(); err != nil {
		return
	}
	if exists {
		changed, err = converger.Client.Upgrade(dryRun)
		if err != nil {
			return result, err
		}
		if changed {
			result.UpgradedResources = append(result.UpgradedResources, ResourceDescriptor{Type: converger.Client.Type(), Name: name})
		}
	} else {
		changed, err = converger.Client.Create(dryRun)
		if err != nil {
			return
		}
		if changed {
			result.CreatedResources = append(result.CreatedResources, ResourceDescriptor{Type: converger.Client.Type(), Name: name})
		}
	}
	for _, child := range converger.Children {
		childResult, err := child.Converge(dryRun)
		result.CreatedResources = append(result.CreatedResources, childResult.CreatedResources...)
		result.UpgradedResources = append(result.UpgradedResources, childResult.UpgradedResources...)
		if err != nil {
			return result, err
		}
	}
	return
}
