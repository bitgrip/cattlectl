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
	"os"
	"path/filepath"

	"github.com/bitgrip/cattlectl/internal/pkg/rancher/descriptor"
	projectModel "github.com/bitgrip/cattlectl/internal/pkg/rancher/project/model"
	"github.com/bitgrip/cattlectl/internal/pkg/template"
	"github.com/rancher/norman/types/slice"
	"github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

// NewProjectParser creates a Parser that is not printing prettified representations
func NewProjectParser(projectFile string, values map[string]interface{}) descriptor.Parser {
	return newProjectParser(projectFile, values, false, []string{})
}

// NewPrettyProjectParser creates a Parser that is printing prettified representations
func NewPrettyProjectParser(projectFile string, values map[string]interface{}) descriptor.Parser {
	return newProjectParser(projectFile, values, true, []string{})
}

func newProjectParser(projectFile string, values map[string]interface{}, pretty bool, parentProjectFiles []string) descriptor.Parser {
	logger := logrus.WithFields(logrus.Fields{
		"project_file": projectFile,
	})
	return fileParser{
		projectFile:        projectFile,
		pretty:             pretty,
		parentProjectFiles: parentProjectFiles,
		logger:             logger,
		values:             values,
	}
}

type fileParser struct {
	projectFile        string
	pretty             bool
	parentProjectFiles []string
	logger             *logrus.Entry
	values             map[string]interface{}
}

func (parser fileParser) Parse(projectData []byte, target interface{}) error {
	targetProject := target.(*projectModel.Project)
	absProjectFile, err := filepath.Abs(parser.projectFile)
	if err != nil {
		return err
	}
	if slice.ContainsString(parser.parentProjectFiles, absProjectFile) {
		parser.logger.Info("Cycle detected - return empty result", parser.parentProjectFiles, absProjectFile)
		return nil
	}
	allProjectFiles := append(parser.parentProjectFiles, absProjectFile)

	isProject, err := isDescriptor(projectData, "Project", parser.logger)
	if !isProject || err != nil {
		return err
	}

	err = yaml.Unmarshal(projectData, targetProject)
	if err != nil {
		return err
	}
	for _, include := range targetProject.Metadata.Includes {
		includeFiles, err := parser.readIncludeFiles(include)
		if err != nil {
			return err
		}
		for _, includeFile := range includeFiles {
			if err := parser.include(targetProject, includeFile, allProjectFiles); err != nil {
				return err
			}
		}
	}

	return nil
}

func (parser fileParser) readIncludeFiles(include projectModel.Include) ([]string, error) {
	if include.File != "" && include.Files != "" && include.Directory != "" {
		return nil, fmt.Errorf("only one of file, files or directory can have a value")
	}
	if include.File != "" {
		return []string{include.File}, nil
	}
	if include.Files != "" {
		var absFiles string
		if filepath.IsAbs(include.Files) {
			absFiles = include.Files
		} else {
			var err error
			absFiles, err = filepath.Abs(filepath.Clean(fmt.Sprintf("%s/%s", filepath.Dir(parser.projectFile), include.Files)))
			if err != nil {
				return nil, err
			}
		}
		matches, err := filepath.Glob(absFiles)
		if err != nil {
			return nil, err
		}
		return matches, nil
	}
	if include.Directory != "" {
		var absDirectory string
		if filepath.IsAbs(include.Directory) {
			absDirectory = include.Directory
		} else {
			absDirectory = filepath.Clean(fmt.Sprintf("%s/%s", filepath.Dir(parser.projectFile), include.Directory))
		}
		var files []string
		err := filepath.Walk(absDirectory, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}
			if filepath.Ext(path) == ".yaml" || filepath.Ext(path) == ".yml" {
				file, _ := filepath.Abs(path)
				files = append(files, file)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
		return files, nil
	}
	return nil, fmt.Errorf("one of file, files or directory must have a value")
}

func (parser fileParser) include(targetProject *projectModel.Project, file string, allProjectFiles []string) error {
	var childProjectFile string
	if filepath.IsAbs(file) {
		childProjectFile = file
	} else {
		childProjectFile = filepath.Clean(fmt.Sprintf("%s/%s", filepath.Dir(parser.projectFile), file))
	}
	childFileContent, err := ioutil.ReadFile(childProjectFile)
	if err != nil {
		return err
	}
	childProjectData, err := template.BuildTemplate(childFileContent, parser.values, filepath.Dir(childProjectFile), parser.pretty)
	if err != nil {
		return err
	}
	childTarget := projectModel.Project{}
	childParser := newProjectParser(childProjectFile, parser.values, parser.pretty, allProjectFiles)
	err = childParser.Parse(childProjectData, &childTarget)
	if err != nil {
		return err
	}
	err = MergeProject(childTarget, targetProject)
	return err
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
