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

	"github.com/bitgrip/cattlectl/internal/pkg/config"
	"github.com/bitgrip/cattlectl/internal/pkg/rancher/project"
	projectModel "github.com/bitgrip/cattlectl/internal/pkg/rancher/project/model"
	yaml "gopkg.in/yaml.v2"
)

const (
	ProjectKind = "Project"
	JobKind     = "Job"
	CronJobKind = "CronJob"
)

var (
	newProjectConverger = project.NewProjectConverger
	newJobConverger     = project.NewJobConverger
	newProjectParser    = project.NewProjectParser
	newJobParser        = project.NewJobParser
)

func ApplyDescriptor(file string, data []byte, values map[string]interface{}, config config.Config) error {
	kind, err := GetKind(data)
	if err != nil {
		return err
	}
	switch kind {
	case ProjectKind:
		project := projectModel.Project{}
		if err := newProjectParser(file, data, &project, values).Parse(); err != nil {
			return err
		}
		if err := ApplyProject(project, config); err != nil {
			return err
		}
		return nil
	case JobKind:
		jobDescriptor := projectModel.JobDescriptor{}
		if err := newJobParser(file, data, &jobDescriptor, values).Parse(); err != nil {
			return err
		}
		if err := ApplyJob(jobDescriptor, config); err != nil {
			return err
		}
		return nil
	}
	return nil
}

func ApplyJob(jobDescriptor projectModel.JobDescriptor, config config.Config) error {
	if config.RancherUrl() != "" {
		jobDescriptor.Metadata.RancherURL = config.RancherUrl()
	}
	if config.AccessKey() != "" {
		jobDescriptor.Metadata.AccessKey = config.AccessKey()
	}
	if config.SecretKey() != "" {
		jobDescriptor.Metadata.SecretKey = config.SecretKey()
	}
	if config.TokenKey() != "" {
		jobDescriptor.Metadata.TokenKey = config.TokenKey()
	}
	if config.ClusterName() != "" {
		jobDescriptor.Metadata.ClusterName = config.ClusterName()
	}
	if config.ClusterId() != "" {
		jobDescriptor.Metadata.ClusterID = config.ClusterId()
	}
	return newJobConverger(jobDescriptor).Converge()
}

func ApplyProject(project projectModel.Project, config config.Config) error {
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
	return newProjectConverger(project).Converge()
}

func GetKind(data []byte) (string, error) {
	structure := make(map[string]interface{})
	if err := yaml.Unmarshal(data, &structure); err != nil {
		return "", err
	}
	if kind, exists := structure["kind"]; exists {
		return fmt.Sprint(kind), nil
	}

	return "UNKNOWN", fmt.Errorf("Kind is undefined")
}

func ParseAndPrintDescriptor(file string, data []byte, values map[string]interface{}, config config.Config) error {
	kind, err := GetKind(data)
	if err != nil {
		return err
	}
	var descriptor interface{}
	switch kind {
	case ProjectKind:
		project := projectModel.Project{}
		if err = newProjectParser(file, data, &project, values).Parse(); err != nil {
			return err
		}
		descriptor = project
	case JobKind:
		jobDescriptor := projectModel.JobDescriptor{}
		if err = newJobParser(file, data, &jobDescriptor, values).Parse(); err != nil {
			return err
		}
		descriptor = jobDescriptor
	}
	out, err := yaml.Marshal(descriptor)
	if err != nil {
		return err
	}
	fmt.Println(string(out))
	return nil
}
