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
	ProjectKind     = "Project"
	JobKind         = "Job"
	CronJobKind     = "CronJob"
	DeploymentKind  = "Deployment"
	DaemonSetKind   = "DaemonSet"
	StatefulSetKind = "StatefulSet"
)

var (
	newProjectConverger     = project.NewProjectConverger
	newJobConverger         = project.NewJobConverger
	newCronJobConverger     = project.NewCronJobConverger
	newDeploymentConverger  = project.NewDeploymentConverger
	newDaemonSetConverger   = project.NewDaemonSetConverger
	newStatefulSetConverger = project.NewStatefulSetConverger

	newProjectParser     = project.NewProjectParser
	newJobParser         = project.NewJobParser
	newCronJobParser     = project.NewCronJobParser
	newDeploymentParser  = project.NewDeploymentParser
	newDaemonSetParser   = project.NewDaemonSetParser
	newStatefulSetParser = project.NewStatefulSetParser
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
	case JobKind:
		jobDescriptor := projectModel.JobDescriptor{}
		if err := newJobParser(file, data, &jobDescriptor, values).Parse(); err != nil {
			return err
		}
		if err := ApplyJob(jobDescriptor, config); err != nil {
			return err
		}
	case CronJobKind:
		cronJobDescriptor := projectModel.CronJobDescriptor{}
		if err := newCronJobParser(file, data, &cronJobDescriptor, values).Parse(); err != nil {
			return err
		}
		if err := ApplyCronJob(cronJobDescriptor, config); err != nil {
			return err
		}
	case DeploymentKind:
		deploymentDescriptor := projectModel.DeploymentDescriptor{}
		if err := newDeploymentParser(file, data, &deploymentDescriptor, values).Parse(); err != nil {
			return err
		}
		if err := ApplyDeployment(deploymentDescriptor, config); err != nil {
			return err
		}
	case DaemonSetKind:
		daemonSetDescriptor := projectModel.DaemonSetDescriptor{}
		if err := newDaemonSetParser(file, data, &daemonSetDescriptor, values).Parse(); err != nil {
			return err
		}
		if err := ApplyDaemonSet(daemonSetDescriptor, config); err != nil {
			return err
		}
	case StatefulSetKind:
		statefulSetDescriptor := projectModel.StatefulSetDescriptor{}
		if err := newStatefulSetParser(file, data, &statefulSetDescriptor, values).Parse(); err != nil {
			return err
		}
		if err := ApplyStatefulSet(statefulSetDescriptor, config); err != nil {
			return err
		}
	}
	return nil
}

func ApplyCronJob(cronJobDescriptor projectModel.CronJobDescriptor, config config.Config) error {
	fillWorkloadMetadata(&cronJobDescriptor.Metadata, config)
	return newCronJobConverger(cronJobDescriptor).Converge()
}

func ApplyJob(jobDescriptor projectModel.JobDescriptor, config config.Config) error {
	fillWorkloadMetadata(&jobDescriptor.Metadata, config)
	return newJobConverger(jobDescriptor).Converge()
}

func ApplyDeployment(deploymentDescriptor projectModel.DeploymentDescriptor, config config.Config) error {
	fillWorkloadMetadata(&deploymentDescriptor.Metadata, config)
	return newDeploymentConverger(deploymentDescriptor).Converge()
}

func ApplyDaemonSet(daemonSetDescriptor projectModel.DaemonSetDescriptor, config config.Config) error {
	fillWorkloadMetadata(&daemonSetDescriptor.Metadata, config)
	return newDaemonSetConverger(daemonSetDescriptor).Converge()
}

func ApplyStatefulSet(statefulSetDescriptor projectModel.StatefulSetDescriptor, config config.Config) error {
	fillWorkloadMetadata(&statefulSetDescriptor.Metadata, config)
	return newStatefulSetConverger(statefulSetDescriptor).Converge()
}

func ApplyProject(project projectModel.Project, config config.Config) error {
	fillProjectMetadata(&project.Metadata, config)
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
	case CronJobKind:
		cronJobDescriptor := projectModel.CronJobDescriptor{}
		if err = newCronJobParser(file, data, &cronJobDescriptor, values).Parse(); err != nil {
			return err
		}
	}
	out, err := yaml.Marshal(descriptor)
	if err != nil {
		return err
	}
	_, err = fmt.Println(string(out))
	return err
}

func fillWorkloadMetadata(metadata *projectModel.WorkloadMetadata, config config.Config) {
	if config.RancherUrl() != "" {
		metadata.RancherURL = config.RancherUrl()
	}
	if config.AccessKey() != "" {
		metadata.AccessKey = config.AccessKey()
	}
	if config.SecretKey() != "" {
		metadata.SecretKey = config.SecretKey()
	}
	if config.TokenKey() != "" {
		metadata.TokenKey = config.TokenKey()
	}
	if config.ClusterName() != "" {
		metadata.ClusterName = config.ClusterName()
	}
	if config.ClusterId() != "" {
		metadata.ClusterID = config.ClusterId()
	}
}

func fillProjectMetadata(metadata *projectModel.ProjectMetadata, config config.Config) {
	if config.RancherUrl() != "" {
		metadata.RancherURL = config.RancherUrl()
	}
	if config.AccessKey() != "" {
		metadata.AccessKey = config.AccessKey()
	}
	if config.SecretKey() != "" {
		metadata.SecretKey = config.SecretKey()
	}
	if config.TokenKey() != "" {
		metadata.TokenKey = config.TokenKey()
	}
	if config.ClusterName() != "" {
		metadata.ClusterName = config.ClusterName()
	}
	if config.ClusterId() != "" {
		metadata.ClusterID = config.ClusterId()
	}
}
