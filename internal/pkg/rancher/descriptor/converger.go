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
	Converge() error
}

type ResourceClientConverger struct {
	Client   client.ResourceClient
	Children []Converger
}

func (converger *ResourceClientConverger) Converge() error {
	var (
		exists bool
		err    error
	)
	exists, err = converger.Client.Exists()
	if err != nil {
		return err
	}
	if exists {
		err = converger.Client.Upgrade()
	} else {
		err = converger.Client.Create()
	}
	if err != nil {
		return err
	}
	for _, child := range converger.Children {
		err = child.Converge()
		if err != nil {
			return err
		}
	}
	return nil
}
