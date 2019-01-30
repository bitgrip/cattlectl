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
	"github.com/bitgrip/cattlectl/internal/pkg/config"
	"github.com/bitgrip/cattlectl/internal/pkg/rancher"
	"github.com/bitgrip/cattlectl/internal/pkg/rancher/project"
)

var (
	newConverger = project.NewConverger
)

func ApplyProject(project rancher.Project, config config.Config) error {
	if config.RancherUrl() != "" {
		project.Metadata.RancherURL = config.RancherUrl()
	}
	if config.AccessKey() != "" {
		project.Metadata.AccessKey = config.AccessKey()
	}
	if config.SecretKey() != "" {
		project.Metadata.SecretKey = config.SecretKey()
	}
	if config.TokenKey() != "" {
		project.Metadata.TokenKey = config.TokenKey()
	}
	if config.ClusterName() != "" {
		project.Metadata.ClusterName = config.ClusterName()
	}
	if config.ClusterId() != "" {
		project.Metadata.ClusterID = config.ClusterId()
	}
	return newConverger(project).Converge()
}
