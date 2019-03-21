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

// NewJobParser creates a Parser that is printing prettified representations
func NewCronJobParser(descriptorFile string, jobData []byte, target *projectModel.CronJobDescriptor, values map[string]interface{}) Parser {
	logger := logrus.WithFields(logrus.Fields{
		"descriptor_file": descriptorFile,
	})
	return cronJobParser{
		logger:  logger,
		target:  target,
		jobData: jobData,
		values:  values,
	}
}

type cronJobParser struct {
	logger  *logrus.Entry
	target  *projectModel.CronJobDescriptor
	jobData []byte
	values  map[string]interface{}
}

func (parser cronJobParser) Parse() error {
	isProject, err := isDescriptor(parser.jobData, "CronJob", parser.logger)
	if !isProject || err != nil {
		return err
	}

	return yaml.Unmarshal(parser.jobData, parser.target)
}
