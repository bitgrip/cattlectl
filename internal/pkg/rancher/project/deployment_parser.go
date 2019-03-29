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

package project

import (
	projectModel "github.com/bitgrip/cattlectl/internal/pkg/rancher/project/model"
	"github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

// NewDeploymentParser creates a Parser that is printing prettified representations
func NewDeploymentParser(descriptorFile string, deploymentData []byte, target *projectModel.DeploymentDescriptor, values map[string]interface{}) Parser {
	logger := logrus.WithFields(logrus.Fields{
		"descriptor_file": descriptorFile,
	})
	return deploymentParser{
		logger:         logger,
		target:         target,
		deploymentData: deploymentData,
		values:         values,
	}
}

type deploymentParser struct {
	logger         *logrus.Entry
	target         *projectModel.DeploymentDescriptor
	deploymentData []byte
	values         map[string]interface{}
}

func (parser deploymentParser) Parse() error {
	isProject, err := isDescriptor(parser.deploymentData, "Deployment", parser.logger)
	if !isProject || err != nil {
		return err
	}

	return yaml.Unmarshal(parser.deploymentData, parser.target)
}
