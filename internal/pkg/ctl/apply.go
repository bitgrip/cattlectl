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
	"bytes"
	"fmt"
	"strings"

	"github.com/bitgrip/cattlectl/internal/pkg/config"
	rancher "github.com/bitgrip/cattlectl/internal/pkg/rancher"
	rancher_client "github.com/bitgrip/cattlectl/internal/pkg/rancher/client"
	cluster "github.com/bitgrip/cattlectl/internal/pkg/rancher/cluster"
	clusterModel "github.com/bitgrip/cattlectl/internal/pkg/rancher/cluster/model"
	"github.com/bitgrip/cattlectl/internal/pkg/rancher/cluster/project"
	projectModel "github.com/bitgrip/cattlectl/internal/pkg/rancher/cluster/project/model"
	"github.com/bitgrip/cattlectl/internal/pkg/rancher/descriptor"
	rancherModel "github.com/bitgrip/cattlectl/internal/pkg/rancher/model"
	yaml "gopkg.in/yaml.v2"
)

var (
	println = fmt.Println

	newRancherClient = rancher_client.NewRancherClient

	newRancherConverger     = rancher.NewRancherConverger
	newClusterConverger     = cluster.NewClusterConverger
	newProjectConverger     = project.NewProjectConverger
	newJobConverger         = project.NewJobConverger
	newCronJobConverger     = project.NewCronJobConverger
	newDeploymentConverger  = project.NewDeploymentConverger
	newDaemonSetConverger   = project.NewDaemonSetConverger
	newStatefulSetConverger = project.NewStatefulSetConverger

	newRancherParser     = rancher.NewRancherParser
	newClusterParser     = cluster.NewClusterParser
	newProjectParser     = project.NewProjectParser
	newJobParser         = project.NewJobParser
	newCronJobParser     = project.NewCronJobParser
	newDeploymentParser  = project.NewDeploymentParser
	newDaemonSetParser   = project.NewDaemonSetParser
	newStatefulSetParser = project.NewStatefulSetParser
)

// ApplyDescriptor the the CTL perform a apply action
func ApplyDescriptor(file string, fullData []byte, values map[string]interface{}, config config.Config) (result descriptor.ConvergeResult, err error) {
	decoder := yaml.NewDecoder(bytes.NewReader(fullData))
	for {
		apiVersion, kind, object, decodeErr := DecodeToApply(decoder)
		if decodeErr != nil {
			if decodeErr.Error() == "EMPTY" {
				continue
			} else if decodeErr.Error() == "EOF" {
				break
			}
			err = decodeErr
			return
		}
		var (
			singleObjectData []byte
			singleResult     descriptor.ConvergeResult
		)
		if singleObjectData, err = yaml.Marshal(object); err != nil {
			return
		}
		if !isSupportedAPIVersion(apiVersion) {
			return result, fmt.Errorf("Unsupported api version %s", apiVersion)
		}
		switch kind {
		case rancherModel.RancherKind:
			descriptor := rancherModel.Rancher{}
			if err = newRancherParser(file, values).Parse(singleObjectData, &descriptor); err != nil {
				return
			}
			if singleResult, err = ApplyRancher(descriptor, config); err != nil {
				return
			}
		case rancherModel.ClusterKind:
			descriptor := clusterModel.Cluster{}
			if err = newClusterParser(file, values).Parse(singleObjectData, &descriptor); err != nil {
				return
			}
			if singleResult, err = ApplyCluster(descriptor, config); err != nil {
				return
			}
		case rancherModel.ProjectKind:
			project := projectModel.Project{}
			if err = newProjectParser(file, values).Parse(singleObjectData, &project); err != nil {
				return
			}
			if singleResult, err = ApplyProject(project, config); err != nil {
				return
			}
		case rancherModel.JobKind:
			jobDescriptor := projectModel.JobDescriptor{}
			if err = newJobParser(file, values).Parse(singleObjectData, &jobDescriptor); err != nil {
				return
			}
			if singleResult, err = ApplyJob(jobDescriptor, config); err != nil {
				return
			}
		case rancherModel.CronJobKind:
			cronJobDescriptor := projectModel.CronJobDescriptor{}
			if err = newCronJobParser(file, values).Parse(singleObjectData, &cronJobDescriptor); err != nil {
				return
			}
			if singleResult, err = ApplyCronJob(cronJobDescriptor, config); err != nil {
				return
			}
		case rancherModel.DeploymentKind:
			deploymentDescriptor := projectModel.DeploymentDescriptor{}
			if err = newDeploymentParser(file, values).Parse(singleObjectData, &deploymentDescriptor); err != nil {
				return
			}
			if singleResult, err = ApplyDeployment(deploymentDescriptor, config); err != nil {
				return
			}
		case rancherModel.DaemonSetKind:
			daemonSetDescriptor := projectModel.DaemonSetDescriptor{}
			if err = newDaemonSetParser(file, values).Parse(singleObjectData, &daemonSetDescriptor); err != nil {
				return
			}
			if singleResult, err = ApplyDaemonSet(daemonSetDescriptor, config); err != nil {
				return
			}
		case rancherModel.StatefulSetKind:
			statefulSetDescriptor := projectModel.StatefulSetDescriptor{}
			if err = newStatefulSetParser(file, values).Parse(singleObjectData, &statefulSetDescriptor); err != nil {
				return
			}
			if singleResult, err = ApplyStatefulSet(statefulSetDescriptor, config); err != nil {
				return
			}
		default:
			return result, fmt.Errorf("Unknown descriptor %s", kind)
		}
		result.CreatedResources = append(result.CreatedResources, singleResult.CreatedResources...)
		result.UpgradedResources = append(result.UpgradedResources, singleResult.UpgradedResources...)
	}
	return
}

// ApplyCronJob the the CTL perform a apply action to a cronjob descriptor
func ApplyCronJob(cronJobDescriptor projectModel.CronJobDescriptor, config config.Config) (result descriptor.ConvergeResult, err error) {
	_, _, projectClient, err := fillWorkloadMetadata(&cronJobDescriptor.Metadata, config)
	if err != nil {
		return
	}
	converger, err := newCronJobConverger(cronJobDescriptor, projectClient)
	if err != nil {
		return
	}
	return converger.Converge(config.DryRun())
}

// ApplyJob the the CTL perform a apply action to a job descriptor
func ApplyJob(jobDescriptor projectModel.JobDescriptor, config config.Config) (result descriptor.ConvergeResult, err error) {
	_, _, projectClient, err := fillWorkloadMetadata(&jobDescriptor.Metadata, config)
	if err != nil {
		return
	}
	converger, err := newJobConverger(jobDescriptor, projectClient)
	if err != nil {
		return
	}
	return converger.Converge(config.DryRun())
}

// ApplyDeployment the the CTL perform a apply action to a deployment descriptor
func ApplyDeployment(deploymentDescriptor projectModel.DeploymentDescriptor, config config.Config) (result descriptor.ConvergeResult, err error) {
	_, _, projectClient, err := fillWorkloadMetadata(&deploymentDescriptor.Metadata, config)
	if err != nil {
		return
	}
	converger, err := newDeploymentConverger(deploymentDescriptor, projectClient)
	if err != nil {
		return
	}
	return converger.Converge(config.DryRun())
}

// ApplyDaemonSet the the CTL perform a apply action to a daemon set descriptor
func ApplyDaemonSet(daemonSetDescriptor projectModel.DaemonSetDescriptor, config config.Config) (result descriptor.ConvergeResult, err error) {
	_, _, projectClient, err := fillWorkloadMetadata(&daemonSetDescriptor.Metadata, config)
	if err != nil {
		return
	}
	converger, err := newDaemonSetConverger(daemonSetDescriptor, projectClient)
	if err != nil {
		return
	}
	return converger.Converge(config.DryRun())
}

// ApplyStatefulSet the the CTL perform a apply action to a stateful set descriptor
func ApplyStatefulSet(statefulSetDescriptor projectModel.StatefulSetDescriptor, config config.Config) (result descriptor.ConvergeResult, err error) {
	_, _, projectClient, err := fillWorkloadMetadata(&statefulSetDescriptor.Metadata, config)
	if err != nil {
		return
	}
	converger, err := newStatefulSetConverger(statefulSetDescriptor, projectClient)
	if err != nil {
		return
	}
	return converger.Converge(config.DryRun())
}

// ApplyProject the the CTL perform a apply action to a project descriptor
func ApplyProject(project projectModel.Project, config config.Config) (result descriptor.ConvergeResult, err error) {
	_, clusterClient, err := fillProjectMetadata(&project.Metadata, config)
	if err != nil {
		return
	}
	converger, err := newProjectConverger(project, clusterClient)
	if err != nil {
		return
	}
	return converger.Converge(config.DryRun())
}

// ApplyCluster the the CTL perform a apply action to a cluster descriptor
func ApplyCluster(cluster clusterModel.Cluster, config config.Config) (result descriptor.ConvergeResult, err error) {
	rancherClient, err := fillClusterMetadata(&cluster.Metadata, config)
	if err != nil {
		return
	}
	converger, err := newClusterConverger(cluster, rancherClient)
	if err != nil {
		return
	}
	return converger.Converge(config.DryRun())
}

// ApplyRancher the the CTL perform a apply action to a cluster descriptor
func ApplyRancher(rancher rancherModel.Rancher, config config.Config) (result descriptor.ConvergeResult, err error) {
	rancherConfig, err := fillRancherMetadata(&rancher.Metadata, config)
	if err != nil {
		return
	}
	converger, err := newRancherConverger(rancher, rancherConfig)
	if err != nil {
		return
	}
	return converger.Converge(config.DryRun())
}

// DecodeToApply will decode the next object from the decoder
func DecodeToApply(decoder *yaml.Decoder) (apiVersion, kind string, object map[string]interface{}, err error) {
	err = decoder.Decode(&object)
	if err != nil {
		return
	}
	if len(object) == 0 {
		err = fmt.Errorf("EMPTY")
		return
	}
	if foundKind, exists := object["kind"]; exists {
		kind = fmt.Sprint(foundKind)
	} else {
		err = fmt.Errorf("Kind is undefined")
		return
	}
	if foundAPIVersion, exists := object["api_version"]; exists {
		apiVersion = fmt.Sprint(foundAPIVersion)
	} else {
		err = fmt.Errorf("API version is undefined")
		return
	}
	return
}

// ParseAndPrintDescriptor parse and print the 'data'
func ParseAndPrintDescriptor(file string, fullData []byte, values map[string]interface{}, config config.Config) error {
	decoder := yaml.NewDecoder(bytes.NewReader(fullData))
	for {
		apiVersion, kind, object, err := DecodeToApply(decoder)
		if err != nil {
			if err.Error() == "EMPTY" {
				continue
			} else if err.Error() == "EOF" {
				break
			}
			return err
		}
		singleObjectData, err := yaml.Marshal(object)
		if err != nil {
			return err
		}
		if !isSupportedAPIVersion(apiVersion) {
			return fmt.Errorf("Unsupported api version %s", apiVersion)
		}
		var descriptor interface{}
		switch kind {
		case rancherModel.RancherKind:
			descriptor = rancherModel.Rancher{}
			if err = newRancherParser(file, values).Parse(singleObjectData, &descriptor); err != nil {
				return err
			}
		case rancherModel.ClusterKind:
			descriptor = clusterModel.Cluster{}
			if err = newClusterParser(file, values).Parse(singleObjectData, &descriptor); err != nil {
				return err
			}
		case rancherModel.ProjectKind:
			project := projectModel.Project{}
			if err = newProjectParser(file, values).Parse(singleObjectData, &project); err != nil {
				return err
			}
			descriptor = project
		case rancherModel.JobKind:
			jobDescriptor := projectModel.JobDescriptor{}
			if err = newJobParser(file, values).Parse(singleObjectData, &jobDescriptor); err != nil {
				return err
			}
			descriptor = jobDescriptor
		case rancherModel.CronJobKind:
			cronJobDescriptor := projectModel.CronJobDescriptor{}
			if err = newCronJobParser(file, values).Parse(singleObjectData, &cronJobDescriptor); err != nil {
				return err
			}
			descriptor = cronJobDescriptor
		case rancherModel.DeploymentKind:
			descriptor = projectModel.Deployment{}
			if err = newDeploymentParser(file, values).Parse(singleObjectData, &descriptor); err != nil {
				return err
			}
		case rancherModel.StatefulSetKind:
			descriptor = projectModel.StatefulSet{}
			if err = newStatefulSetParser(file, values).Parse(singleObjectData, &descriptor); err != nil {
				return err
			}
		case rancherModel.DaemonSetKind:
			descriptor = projectModel.DaemonSet{}
			if err = newDaemonSetParser(file, values).Parse(singleObjectData, &descriptor); err != nil {
				return err
			}
		default:
			return fmt.Errorf("Unknown descriptor %s", kind)
		}
		out, err := yaml.Marshal(descriptor)
		if err != nil {
			return err
		}
		println("---")
		println(strings.TrimSpace(string(out)))

	}
	return nil
}

func fillWorkloadMetadata(metadata *projectModel.WorkloadMetadata, config config.Config) (rancher_client.RancherClient, rancher_client.ClusterClient, rancher_client.ProjectClient, error) {
	if config.RancherURL() != "" {
		metadata.RancherURL = config.RancherURL()
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
	if config.ClusterID() != "" {
		metadata.ClusterID = config.ClusterID()
	}

	rancherClient, err := newRancherClient(rancher_client.RancherConfig{
		RancherURL:   metadata.RancherURL,
		AccessKey:    metadata.AccessKey,
		SecretKey:    metadata.SecretKey,
		Insecure:     config.InsecureAPI(),
		CACerts:      config.CACerts(),
		MergeAnswers: config.MergeAnswers(),
	})
	if err != nil {
		return nil, nil, nil, err
	}
	clusterClient, err := rancherClient.Cluster(metadata.ClusterName)
	if err != nil {
		return nil, nil, nil, err
	}
	projectClient, err := clusterClient.Project(metadata.ProjectName)
	if err != nil {
		return nil, nil, nil, err
	}
	return rancherClient, clusterClient, projectClient, nil
}

func fillProjectMetadata(metadata *projectModel.ProjectMetadata, config config.Config) (rancher_client.RancherClient, rancher_client.ClusterClient, error) {
	if config.RancherURL() != "" {
		metadata.RancherURL = config.RancherURL()
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
	if config.ClusterID() != "" {
		metadata.ClusterID = config.ClusterID()
	}

	rancherClient, err := newRancherClient(rancher_client.RancherConfig{
		RancherURL:   metadata.RancherURL,
		AccessKey:    metadata.AccessKey,
		SecretKey:    metadata.SecretKey,
		Insecure:     config.InsecureAPI(),
		CACerts:      config.CACerts(),
		MergeAnswers: config.MergeAnswers(),
	})
	if err != nil {
		return nil, nil, err
	}
	clusterClient, err := rancherClient.Cluster(metadata.ClusterName)
	if err != nil {
		return nil, nil, err
	}
	return rancherClient, clusterClient, nil
}

func fillClusterMetadata(metadata *clusterModel.ClusterMetadata, config config.Config) (rancher_client.RancherClient, error) {
	if config.RancherURL() != "" {
		metadata.RancherURL = config.RancherURL()
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
		metadata.Name = config.ClusterName()
	}
	if config.ClusterID() != "" {
		metadata.ID = config.ClusterID()
	}

	rancherClient, err := newRancherClient(rancher_client.RancherConfig{
		RancherURL: metadata.RancherURL,
		AccessKey:  metadata.AccessKey,
		SecretKey:  metadata.SecretKey,
		Insecure:   config.InsecureAPI(),
		CACerts:    config.CACerts(),
	})
	if err != nil {
		return nil, err
	}
	return rancherClient, nil
}

func fillRancherMetadata(metadata *rancherModel.RancherMetadata, config config.Config) (rancher_client.RancherConfig, error) {
	if config.RancherURL() != "" {
		metadata.RancherURL = config.RancherURL()
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

	return rancher_client.RancherConfig{
		RancherURL: metadata.RancherURL,
		AccessKey:  metadata.AccessKey,
		SecretKey:  metadata.SecretKey,
		Insecure:   config.InsecureAPI(),
		CACerts:    config.CACerts(),
	}, nil
}
