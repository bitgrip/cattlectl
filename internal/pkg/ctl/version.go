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
	"github.com/sirupsen/logrus"
)

var minAPIVersion = "2.0"

// Version is the current build version
var Version = "v2.0.0-local"

func isSupportedAPIVersion(apiVersion string) bool {
	v, err := semver.NewVersion(Version)
	if err != nil {
		return false
	}
	minConstraint, err := semver.NewConstraint(fmt.Sprintf(">= %s", minAPIVersion))
	if err != nil {
		logrus.WithField("error", err).Error("Can not build minimum constraint")
		return false
	}
	maxConstraint, err := semver.NewConstraint(fmt.Sprintf("<= %v.%v", v.Major(), v.Minor()))
	if err != nil {
		logrus.WithField("error", err).Error("Can not build maximum constraint")
		return false
	}

	av, err := semver.NewVersion(apiVersion)
	if err != nil {
		return false
	}
	return minConstraint.Check(av) && maxConstraint.Check(av)
}
