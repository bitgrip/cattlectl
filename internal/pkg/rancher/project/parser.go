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
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/bitgrip/cattlectl/internal/pkg/rancher"
	"github.com/bitgrip/cattlectl/internal/pkg/template"
	"github.com/rancher/norman/types/slice"
	"github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

// Parser is a object that can parse a project file using a map of template values
type Parser interface {
	Parse(values map[string]interface{}) (rancher.Project, error)
}

// NewParser creates a Parser that is not printing prettified representations
func NewParser(projectFile string) Parser {
	return newParser(projectFile, false, []string{})
}

// NewPrettyParser creates a Parser that is printing prettified representations
func NewPrettyParser(projectFile string) Parser {
	return newParser(projectFile, true, []string{})
}

func newParser(projectFile string, pretty bool, parentProjectFiles []string) Parser {
	logger := logrus.WithFields(logrus.Fields{
		"project_file": projectFile,
	})
	return fileParser{
		projectFile:        projectFile,
		pretty:             pretty,
		parentProjectFiles: parentProjectFiles,
		logger:             logger,
	}
}

type fileParser struct {
	projectFile        string
	pretty             bool
	parentProjectFiles []string
	logger             *logrus.Entry
}

func (parser fileParser) Parse(values map[string]interface{}) (rancher.Project, error) {
	project := rancher.Project{}
	absProjectFile, err := filepath.Abs(parser.projectFile)
	if err != nil {
		return project, err
	}
	if slice.ContainsString(parser.parentProjectFiles, absProjectFile) {
		parser.logger.Info("Cycle detected - return empty result")
		return project, nil
	}
	allProjectFiles := append(parser.parentProjectFiles, absProjectFile)
	fileContent, err := ioutil.ReadFile(absProjectFile)
	if err != nil {
		return project, err
	}

	projectData, err := template.BuildTemplate(fileContent, values, filepath.Dir(parser.projectFile), parser.pretty)
	if err != nil {
		return project, err
	}

	isProject, err := parser.isProjectDescriptor(projectData)
	if !isProject || err != nil {
		return project, err
	}

	err = yaml.Unmarshal(projectData, &project)
	if err != nil {
		return project, err
	}
	for _, include := range project.Metadata.Includes {
		var childProjectFile string
		if filepath.IsAbs(include.File) {
			childProjectFile = include.File
		} else {
			childProjectFile = filepath.Clean(fmt.Sprintf("%s/%s", filepath.Dir(parser.projectFile), include.File))
		}
		childParser := newParser(childProjectFile, parser.pretty, allProjectFiles)
		child, err := childParser.Parse(values)
		if err != nil {
			return project, err
		}
		project, err = rancher.MergeProject(child, project)
		if err != nil {
			return project, err
		}
	}

	return project, nil
}

func (parser fileParser) isProjectDescriptor(data []byte) (bool, error) {
	structure := make(map[string]interface{})
	err := yaml.Unmarshal(data, &structure)
	if err != nil {
		return false, err
	}
	if structure["kind"] != "Project" {
		parser.logger.WithField("kind", structure["kind"]).Error("Invalid descriptor")
		return false, fmt.Errorf("Invalid descriptor: %v", structure["kind"])
	}

	return true, nil
}
