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

package ctl

import (
	"fmt"

	"github.com/Masterminds/semver"
)

// Version is the current build version
var Version = "v1.3.0-local"

func isSupportedAPIVersion(apiVersion string) bool {
	v, err := semver.NewVersion(Version)
	if err != nil {
		return false
	}
	c, err := semver.NewConstraint(fmt.Sprintf("<= %v.%v", v.Major(), v.Minor()))
	if err != nil {
		return false
	}

	av, err := semver.NewVersion(apiVersion)
	if err != nil {
		return false
	}
	return c.Check(av)
}
