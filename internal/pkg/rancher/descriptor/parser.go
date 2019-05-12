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
	"fmt"

	"github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

// Parser is a object that can parse a project file using a map of template values
type Parser interface {
	Parse(data []byte, target interface{}) error
}

func NewLogginParser(expectedType string, logger *logrus.Entry, values map[string]interface{}) Parser {
	return loggingParser{
		expectedType: expectedType,
		logger:       logger,
		values:       values,
	}
}

type loggingParser struct {
	logger       *logrus.Entry
	values       map[string]interface{}
	expectedType string
}

func (parser loggingParser) Parse(data []byte, target interface{}) error {
	isProject, err := isDescriptor(data, parser.expectedType, parser.logger)
	if !isProject || err != nil {
		return err
	}

	return yaml.Unmarshal(data, target)
}

func isDescriptor(data []byte, kind string, logger *logrus.Entry) (bool, error) {
	structure := make(map[string]interface{})
	err := yaml.Unmarshal(data, &structure)
	if err != nil {
		return false, err
	}
	if structure["kind"] != kind {
		logger.
			WithField("expected-kind", kind).
			WithField("actual-kind", structure["kind"]).
			Error("Invalid descriptor")
		return false, fmt.Errorf("Invalid descriptor: %v", structure["kind"])
	}

	return true, nil
}
