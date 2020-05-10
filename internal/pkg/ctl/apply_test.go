// Copyright © 2020 Bitgrip <berlin@bitgrip.de>
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
	"strings"
	"testing"

	"github.com/bitgrip/cattlectl/internal/pkg/assert"
	"github.com/bitgrip/cattlectl/internal/pkg/config"
	"github.com/bitgrip/cattlectl/internal/pkg/rancher/client"
	rancher_client "github.com/bitgrip/cattlectl/internal/pkg/rancher/client"
	clusterModel "github.com/bitgrip/cattlectl/internal/pkg/rancher/cluster/model"
	projectModel "github.com/bitgrip/cattlectl/internal/pkg/rancher/cluster/project/model"
	"github.com/bitgrip/cattlectl/internal/pkg/rancher/descriptor"
	rancherModel "github.com/bitgrip/cattlectl/internal/pkg/rancher/model"
	yaml "gopkg.in/yaml.v2"
)

var (
	origPrintln              = println
	origRancherClient        = newRancherClient
	origRancherConverger     = newRancherConverger
	origClusterConverger     = newClusterConverger
	origProjectConverger     = newProjectConverger
	origJobConverger         = newJobConverger
	origCronJobConverger     = newCronJobConverger
	origDeploymentConverger  = newDeploymentConverger
	origDaemonSetConverger   = newDaemonSetConverger
	origStatefulSetConverger = newStatefulSetConverger
	origRancherParser        = newRancherParser
	origClusterParser        = newClusterParser
	origProjectParser        = newProjectParser
	origJobParser            = newJobParser
	origCronJobParser        = newCronJobParser
	origDeploymentParser     = newDeploymentParser
	origDaemonSetParser      = newDaemonSetParser
	origStatefulSetParser    = newStatefulSetParser
)

func TestDecodeToApply(t *testing.T) {
	data := `---
description: With all required data
api_version: "1.1"
kind: TestType
---
description: Without kind
api_version: "1.1"
---
description: Without API version
kind: TestType
---`
	decoder := yaml.NewDecoder(strings.NewReader(data))
	tests := []struct {
		name           string
		wantAPIVersion string
		wantKind       string
		wantObject     map[string]interface{}
		wantErr        error
	}{
		{
			name:           "full",
			wantAPIVersion: "1.1",
			wantKind:       "TestType",
			wantObject: map[string]interface{}{
				"description": "With all required data",
				"api_version": "1.1",
				"kind":        "TestType",
			},
			wantErr: nil,
		},
		{
			name:           "without_kind",
			wantAPIVersion: "",
			wantKind:       "",
			wantObject: map[string]interface{}{
				"description": "Without kind",
				"api_version": "1.1",
			},
			wantErr: fmt.Errorf("Kind is undefined"),
		},
		{
			name:           "without_api_version",
			wantAPIVersion: "",
			wantKind:       "TestType",
			wantObject: map[string]interface{}{
				"description": "Without API version",
				"kind":        "TestType",
			},
			wantErr: fmt.Errorf("API version is undefined"),
		},
		{
			name:           "empty_object",
			wantAPIVersion: "",
			wantKind:       "",
			wantObject:     map[string]interface{}{},
			wantErr:        fmt.Errorf("EMPTY"),
		},
		{
			name:           "no_object",
			wantAPIVersion: "",
			wantKind:       "",
			wantObject:     map[string]interface{}{},
			wantErr:        fmt.Errorf("EOF"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAPIVersion, gotKind, gotObject, err := DecodeToApply(decoder)
			if tt.wantErr != nil {
				assert.NotOk(t, err, tt.wantErr.Error())
			} else {
				assert.Ok(t, err)
			}
			assert.Equals(t, tt.wantAPIVersion, gotAPIVersion)
			assert.Equals(t, tt.wantKind, gotKind)
			whantResult, err := yaml.Marshal(tt.wantObject)
			assert.Ok(t, err)
			gotResult, err := yaml.Marshal(gotObject)
			assert.Ok(t, err)
			assert.Equals(t, whantResult, gotResult)
		})
	}
}

func TestApplyDescriptor(t *testing.T) {
	defer resetBackendCalls()

	type args struct {
		file     string
		fullData []byte
		values   map[string]interface{}
		config   config.Config
	}
	tests := []struct {
		name                string
		args                args
		wantErr             error
		setExpectedBackends func(t *testing.T)
	}{
		{
			name: "one_rancher_object",
			args: args{
				file:     "test-descriptor.yaml",
				fullData: []byte("---\napi_version: \"1.1\"\nkind: Rancher"),
				config:   testConfig{},
			},
			setExpectedBackends: func(t *testing.T) {
				newRancherParser = func(descriptorFile string, values map[string]interface{}) descriptor.Parser {
					return testParser{
						expected:     true,
						expectedData: []byte("api_version: \"1.1\"\nkind: Rancher\n"),
						t:            t,
					}
				}
				newRancherConverger = func(rancher rancherModel.Rancher, rancherConfig client.RancherConfig) (descriptor.Converger, error) {
					return testConverger{
						expected: true,
					}, nil
				}
			},
		},
		{
			name: "one_cluster_object",
			args: args{
				file:     "test-descriptor.yaml",
				fullData: []byte("---\napi_version: \"1.1\"\nkind: Cluster"),
				config:   testConfig{},
			},
			setExpectedBackends: func(t *testing.T) {
				newClusterParser = func(descriptorFile string, values map[string]interface{}) descriptor.Parser {
					return testParser{
						expected:     true,
						expectedData: []byte("api_version: \"1.1\"\nkind: Cluster\n"),
						t:            t,
					}
				}
				var expectedRancherClient client.RancherClient
				newRancherClient = func(config client.RancherConfig) (client.RancherClient, error) {
					var err error
					expectedRancherClient, err = rancher_client.NewRancherClient(config)
					return expectedRancherClient, err
				}
				newClusterConverger = func(cluster clusterModel.Cluster, rancherClient client.RancherClient) (descriptor.Converger, error) {
					assert.Assert(t, expectedRancherClient == rancherClient, "Unexpecte rancher client")
					return testConverger{
						expected: true,
					}, nil
				}
			},
		},
		{
			name: "one_project_object",
			args: args{
				file:     "test-descriptor.yaml",
				fullData: []byte("---\napi_version: \"1.1\"\nkind: Project"),
				config: testConfig{
					clusterName: "test-cluster",
				},
			},
			setExpectedBackends: func(t *testing.T) {
				newProjectParser = func(descriptorFile string, values map[string]interface{}) descriptor.Parser {
					return testParser{
						expected:     true,
						expectedData: []byte("api_version: \"1.1\"\nkind: Project\n"),
						t:            t,
					}
				}
				var expectedClusterClient client.ClusterClient
				newRancherClient = func(config client.RancherConfig) (client.RancherClient, error) {
					rancherClient, err := rancher_client.NewRancherClient(config)
					expectedClusterClient, _ = rancherClient.Cluster("test-cluster")
					return rancherClient, err
				}
				newProjectConverger = func(project projectModel.Project, clusterClient client.ClusterClient) (descriptor.Converger, error) {
					assert.Assert(t, expectedClusterClient == clusterClient, "Unexpecte cluster client")
					return testConverger{
						expected: true,
					}, nil
				}
			},
		},
		{
			name: "one_job_object",
			args: args{
				file:     "test-descriptor.yaml",
				fullData: []byte("---\napi_version: \"1.1\"\nkind: Job\nmetadata:\n  project_name: test-project"),
				config: testConfig{
					clusterName: "test-cluster",
				},
			},
			setExpectedBackends: func(t *testing.T) {
				newJobParser = func(descriptorFile string, values map[string]interface{}) descriptor.Parser {
					return testParser{
						expected:     true,
						expectedData: []byte("api_version: \"1.1\"\nkind: Job\nmetadata:\n  project_name: test-project\n"),
						t:            t,
					}
				}
				var expectedProjectClient client.ProjectClient
				newRancherClient = func(config client.RancherConfig) (client.RancherClient, error) {
					rancherClient, err := rancher_client.NewRancherClient(config)
					clusterClient, _ := rancherClient.Cluster("test-cluster")
					expectedProjectClient, _ = clusterClient.Project("test-project")
					return rancherClient, err
				}
				newJobConverger = func(job projectModel.JobDescriptor, projectClient client.ProjectClient) (descriptor.Converger, error) {
					assert.Assert(t, expectedProjectClient == projectClient, "Unexpecte project client")
					return testConverger{
						expected: true,
					}, nil
				}
			},
		},
		{
			name: "one_cronjob_object",
			args: args{
				file:     "test-descriptor.yaml",
				fullData: []byte("---\napi_version: \"1.1\"\nkind: CronJob\nmetadata:\n  project_name: test-project"),
				config: testConfig{
					clusterName: "test-cluster",
				},
			},
			setExpectedBackends: func(t *testing.T) {
				newCronJobParser = func(descriptorFile string, values map[string]interface{}) descriptor.Parser {
					return testParser{
						expected:     true,
						expectedData: []byte("api_version: \"1.1\"\nkind: CronJob\nmetadata:\n  project_name: test-project\n"),
						t:            t,
					}
				}
				var expectedProjectClient client.ProjectClient
				newRancherClient = func(config client.RancherConfig) (client.RancherClient, error) {
					rancherClient, err := rancher_client.NewRancherClient(config)
					clusterClient, _ := rancherClient.Cluster("test-cluster")
					expectedProjectClient, _ = clusterClient.Project("test-project")
					return rancherClient, err
				}
				newCronJobConverger = func(job projectModel.CronJobDescriptor, projectClient client.ProjectClient) (descriptor.Converger, error) {
					assert.Assert(t, expectedProjectClient == projectClient, "Unexpecte project client")
					return testConverger{
						expected: true,
					}, nil
				}
			},
		},
		{
			name: "one_deployment_object",
			args: args{
				file:     "test-descriptor.yaml",
				fullData: []byte("---\napi_version: \"1.1\"\nkind: Deployment\nmetadata:\n  project_name: test-project"),
				config: testConfig{
					clusterName: "test-cluster",
				},
			},
			setExpectedBackends: func(t *testing.T) {
				newDeploymentParser = func(descriptorFile string, values map[string]interface{}) descriptor.Parser {
					return testParser{
						expected:     true,
						expectedData: []byte("api_version: \"1.1\"\nkind: Deployment\nmetadata:\n  project_name: test-project\n"),
						t:            t,
					}
				}
				var expectedProjectClient client.ProjectClient
				newRancherClient = func(config client.RancherConfig) (client.RancherClient, error) {
					rancherClient, err := rancher_client.NewRancherClient(config)
					clusterClient, _ := rancherClient.Cluster("test-cluster")
					expectedProjectClient, _ = clusterClient.Project("test-project")
					return rancherClient, err
				}
				newDeploymentConverger = func(job projectModel.DeploymentDescriptor, projectClient client.ProjectClient) (descriptor.Converger, error) {
					assert.Assert(t, expectedProjectClient == projectClient, "Unexpecte project client")
					return testConverger{
						expected: true,
					}, nil
				}
			},
		},
		{
			name: "one_daemonset_object",
			args: args{
				file:     "test-descriptor.yaml",
				fullData: []byte("---\napi_version: \"1.1\"\nkind: DaemonSet\nmetadata:\n  project_name: test-project"),
				config: testConfig{
					clusterName: "test-cluster",
				},
			},
			setExpectedBackends: func(t *testing.T) {
				newDaemonSetParser = func(descriptorFile string, values map[string]interface{}) descriptor.Parser {
					return testParser{
						expected:     true,
						expectedData: []byte("api_version: \"1.1\"\nkind: DaemonSet\nmetadata:\n  project_name: test-project\n"),
						t:            t,
					}
				}
				var expectedProjectClient client.ProjectClient
				newRancherClient = func(config client.RancherConfig) (client.RancherClient, error) {
					rancherClient, err := rancher_client.NewRancherClient(config)
					clusterClient, _ := rancherClient.Cluster("test-cluster")
					expectedProjectClient, _ = clusterClient.Project("test-project")
					return rancherClient, err
				}
				newDaemonSetConverger = func(job projectModel.DaemonSetDescriptor, projectClient client.ProjectClient) (descriptor.Converger, error) {
					assert.Assert(t, expectedProjectClient == projectClient, "Unexpecte project client")
					return testConverger{
						expected: true,
					}, nil
				}
			},
		},
		{
			name: "one_statefulset_object",
			args: args{
				file:     "test-descriptor.yaml",
				fullData: []byte("---\napi_version: \"1.1\"\nkind: StatefulSet\nmetadata:\n  project_name: test-project"),
				config: testConfig{
					clusterName: "test-cluster",
				},
			},
			setExpectedBackends: func(t *testing.T) {
				newStatefulSetParser = func(descriptorFile string, values map[string]interface{}) descriptor.Parser {
					return testParser{
						expected:     true,
						expectedData: []byte("api_version: \"1.1\"\nkind: StatefulSet\nmetadata:\n  project_name: test-project\n"),
						t:            t,
					}
				}
				var expectedProjectClient client.ProjectClient
				newRancherClient = func(config client.RancherConfig) (client.RancherClient, error) {
					rancherClient, err := rancher_client.NewRancherClient(config)
					clusterClient, _ := rancherClient.Cluster("test-cluster")
					expectedProjectClient, _ = clusterClient.Project("test-project")
					return rancherClient, err
				}
				newStatefulSetConverger = func(job projectModel.StatefulSetDescriptor, projectClient client.ProjectClient) (descriptor.Converger, error) {
					assert.Assert(t, expectedProjectClient == projectClient, "Unexpecte project client")
					return testConverger{
						expected: true,
					}, nil
				}
			},
		},
		{
			name: "two_project_objects",
			args: args{
				file:     "test-descriptor.yaml",
				fullData: []byte("---\napi_version: \"1.1\"\nkind: Project\nmetadata:\n  name: project1\n---\napi_version: \"1.1\"\nkind: Project\nmetadata:\n  name: project2\n---\n"),
				config: testConfig{
					clusterName: "test-cluster",
				},
			},
			setExpectedBackends: func(t *testing.T) {
				newProjectParser = func(descriptorFile string, values map[string]interface{}) descriptor.Parser {
					return testParser{
						expected: true,
						t:        t,
					}
				}
				var expectedClusterClient client.ClusterClient
				newRancherClient = func(config client.RancherConfig) (client.RancherClient, error) {
					rancherClient, err := rancher_client.NewRancherClient(config)
					expectedClusterClient, _ = rancherClient.Cluster("test-cluster")
					return rancherClient, err
				}
				newProjectConverger = func(project projectModel.Project, clusterClient client.ClusterClient) (descriptor.Converger, error) {
					assert.Assert(t, expectedClusterClient == clusterClient, "Unexpecte cluster client")
					assert.Assert(t, "project1" == project.Metadata.Name || "project2" == project.Metadata.Name, fmt.Sprintf("Unexpected project name %s", project.Metadata.Name))
					return testConverger{
						expected: true,
					}, nil
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			unexpectAllBackendCalls()
			tt.setExpectedBackends(t)
			_, err := ApplyDescriptor(tt.args.file, tt.args.fullData, tt.args.values, tt.args.config)
			if tt.wantErr != nil {
				assert.NotOk(t, err, tt.wantErr.Error())
			} else {
				assert.Ok(t, err)
			}
		})
	}
}

func TestParseAndPrintDescriptor(t *testing.T) {
	defer resetBackendCalls()

	type args struct {
		file     string
		fullData []byte
		values   map[string]interface{}
		config   config.Config
	}
	tests := []struct {
		name                string
		args                args
		wantErr             error
		wantStdout          string
		setExpectedBackends func(t *testing.T)
	}{
		{
			name: "one_rancher_object",
			args: args{
				file:     "test-descriptor.yaml",
				fullData: []byte("---\napi_version: \"1.1\"\nkind: Rancher"),
			},
			setExpectedBackends: func(t *testing.T) {
				newRancherParser = func(descriptorFile string, values map[string]interface{}) descriptor.Parser {
					return testParser{
						expected:     true,
						expectedData: []byte("api_version: \"1.1\"\nkind: Rancher\n"),
						t:            t,
					}
				}
			},
			wantStdout: "---\napi_version: \"1.1\"\nkind: Rancher\n",
		},
		{
			name: "one_cluster_object",
			args: args{
				file:     "test-descriptor.yaml",
				fullData: []byte("---\napi_version: \"1.1\"\nkind: Cluster"),
			},
			setExpectedBackends: func(t *testing.T) {
				newClusterParser = func(descriptorFile string, values map[string]interface{}) descriptor.Parser {
					return testParser{
						expected:     true,
						expectedData: []byte("api_version: \"1.1\"\nkind: Cluster\n"),
						t:            t,
					}
				}
			},
			wantStdout: "---\napi_version: \"1.1\"\nkind: Cluster\n",
		},
		{
			name: "one_project_object",
			args: args{
				file:     "test-descriptor.yaml",
				fullData: []byte("---\napi_version: \"1.1\"\nkind: Project"),
			},
			setExpectedBackends: func(t *testing.T) {
				newProjectParser = func(descriptorFile string, values map[string]interface{}) descriptor.Parser {
					return testParser{
						expected:     true,
						expectedData: []byte("api_version: \"1.1\"\nkind: Project\n"),
						t:            t,
					}
				}
			},
			wantStdout: "---\napi_version: \"1.1\"\nkind: Project\n",
		},
		{
			name: "one_job_object",
			args: args{
				file:     "test-descriptor.yaml",
				fullData: []byte("---\napi_version: \"1.1\"\nkind: Job"),
			},
			setExpectedBackends: func(t *testing.T) {
				newJobParser = func(descriptorFile string, values map[string]interface{}) descriptor.Parser {
					return testParser{
						expected:     true,
						expectedData: []byte("api_version: \"1.1\"\nkind: Job\n"),
						t:            t,
					}
				}
			},
			wantStdout: "---\napi_version: \"1.1\"\nkind: Job\nmetadata: {}\nspec: {}\n",
		},
		{
			name: "one_cronjob_object",
			args: args{
				file:     "test-descriptor.yaml",
				fullData: []byte("---\napi_version: \"1.1\"\nkind: CronJob"),
			},
			setExpectedBackends: func(t *testing.T) {
				newCronJobParser = func(descriptorFile string, values map[string]interface{}) descriptor.Parser {
					return testParser{
						expected:     true,
						expectedData: []byte("api_version: \"1.1\"\nkind: CronJob\n"),
						t:            t,
					}
				}
			},
			wantStdout: "---\napi_version: \"1.1\"\nkind: CronJob\nmetadata: {}\nspec: {}\n",
		},
		{
			name: "one_Deployment_object",
			args: args{
				file:     "test-descriptor.yaml",
				fullData: []byte("---\napi_version: \"1.1\"\nkind: Deployment"),
			},
			setExpectedBackends: func(t *testing.T) {
				newDeploymentParser = func(descriptorFile string, values map[string]interface{}) descriptor.Parser {
					return testParser{
						expected:     true,
						expectedData: []byte("api_version: \"1.1\"\nkind: Deployment\n"),
						t:            t,
					}
				}
			},
			wantStdout: "---\napi_version: \"1.1\"\nkind: Deployment\n",
		},
		{
			name: "one_statefulset_object",
			args: args{
				file:     "test-descriptor.yaml",
				fullData: []byte("---\napi_version: \"1.1\"\nkind: StatefulSet"),
			},
			setExpectedBackends: func(t *testing.T) {
				newStatefulSetParser = func(descriptorFile string, values map[string]interface{}) descriptor.Parser {
					return testParser{
						expected:     true,
						expectedData: []byte("api_version: \"1.1\"\nkind: StatefulSet\n"),
						t:            t,
					}
				}
			},
			wantStdout: "---\napi_version: \"1.1\"\nkind: StatefulSet\n",
		},
		{
			name: "one_daemonset_object",
			args: args{
				file:     "test-descriptor.yaml",
				fullData: []byte("---\napi_version: \"1.1\"\nkind: DaemonSet"),
			},
			setExpectedBackends: func(t *testing.T) {
				newDaemonSetParser = func(descriptorFile string, values map[string]interface{}) descriptor.Parser {
					return testParser{
						expected:     true,
						expectedData: []byte("api_version: \"1.1\"\nkind: DaemonSet\n"),
						t:            t,
					}
				}
			},
			wantStdout: "---\napi_version: \"1.1\"\nkind: DaemonSet\n",
		},
		{
			name: "one_project_one_empty_object",
			args: args{
				file:     "test-descriptor.yaml",
				fullData: []byte("---\napi_version: \"1.1\"\nkind: Project\n---\n"),
			},
			setExpectedBackends: func(t *testing.T) {
				newProjectParser = func(descriptorFile string, values map[string]interface{}) descriptor.Parser {
					return testParser{
						expected:     true,
						expectedData: []byte("api_version: \"1.1\"\nkind: Project\n"),
						t:            t,
					}
				}
			},
			wantStdout: "---\napi_version: \"1.1\"\nkind: Project\n",
		},
		{
			name: "two_project_objects",
			args: args{
				file:     "test-descriptor.yaml",
				fullData: []byte("---\napi_version: \"1.1\"\nkind: Project\n---\napi_version: \"1.1\"\nkind: Project"),
			},
			setExpectedBackends: func(t *testing.T) {
				newProjectParser = func(descriptorFile string, values map[string]interface{}) descriptor.Parser {
					return testParser{
						expected:     true,
						expectedData: []byte("api_version: \"1.1\"\nkind: Project\n"),
						t:            t,
					}
				}
			},
			wantStdout: "---\napi_version: \"1.1\"\nkind: Project\n---\napi_version: \"1.1\"\nkind: Project\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			unexpectAllBackendCalls()
			tt.setExpectedBackends(t)
			stdout := ""
			println = func(a ...interface{}) (n int, err error) {
				stdout = fmt.Sprintf("%s%v\n", stdout, a[0])
				return 0, nil
			}
			err := ParseAndPrintDescriptor(tt.args.file, tt.args.fullData, tt.args.values, tt.args.config)
			if tt.wantErr != nil {
				assert.NotOk(t, err, tt.wantErr.Error())
			} else {
				assert.Ok(t, err)
			}
			assert.Equals(t, tt.wantStdout, stdout)
		})
	}
}

func resetBackendCalls() {
	println = origPrintln
	newRancherClient = origRancherClient
	newRancherConverger = origRancherConverger
	newClusterConverger = origClusterConverger
	newProjectConverger = origProjectConverger
	newJobConverger = origJobConverger
	newCronJobConverger = origCronJobConverger
	newDeploymentConverger = origDeploymentConverger
	newDaemonSetConverger = origDaemonSetConverger
	newStatefulSetConverger = origStatefulSetConverger
	newRancherParser = origRancherParser
	newClusterParser = origClusterParser
	newProjectParser = origProjectParser
	newJobParser = origJobParser
	newCronJobParser = origCronJobParser
	newDeploymentParser = origDeploymentParser
	newDaemonSetParser = origDaemonSetParser
	newStatefulSetParser = origStatefulSetParser
}

func unexpectAllBackendCalls() {
	newRancherClient = func(config client.RancherConfig) (client.RancherClient, error) {
		return nil, fmt.Errorf("Unexpected Call newRancherClient(...)")
	}
	newRancherConverger = func(rancher rancherModel.Rancher, rancherConfig client.RancherConfig) (descriptor.Converger, error) {
		return nil, fmt.Errorf("Unexpected Call newRancherConverger(...)")
	}
	newClusterConverger = func(cluster clusterModel.Cluster, rancherClient client.RancherClient) (descriptor.Converger, error) {
		return nil, fmt.Errorf("Unexpected Call newClusterConverger(...)")
	}
	newProjectConverger = func(project projectModel.Project, clusterClient client.ClusterClient) (descriptor.Converger, error) {
		return nil, fmt.Errorf("Unexpected Call newProjectConverger(...)")
	}
	newJobConverger = func(jobDescriptor projectModel.JobDescriptor, projectClient client.ProjectClient) (descriptor.Converger, error) {
		return nil, fmt.Errorf("Unexpected Call newJobConverger(...)")
	}
	newCronJobConverger = func(cronJobDescriptor projectModel.CronJobDescriptor, projectClient client.ProjectClient) (descriptor.Converger, error) {
		return nil, fmt.Errorf("Unexpected Call newCronJobConverger(...)")
	}
	newDeploymentConverger = func(deploymentDescriptor projectModel.DeploymentDescriptor, projectClient client.ProjectClient) (descriptor.Converger, error) {
		return nil, fmt.Errorf("Unexpected Call newDeploymentConverger(...)")
	}
	newDaemonSetConverger = func(daemonSetDescriptor projectModel.DaemonSetDescriptor, projectClient client.ProjectClient) (descriptor.Converger, error) {
		return nil, fmt.Errorf("Unexpected Call newDaemonSetConverger(...)")
	}
	newStatefulSetConverger = func(statefulSetDescriptor projectModel.StatefulSetDescriptor, projectClient client.ProjectClient) (descriptor.Converger, error) {
		return nil, fmt.Errorf("Unexpected Call newStatefulSetConverger(...)")
	}
	newRancherParser = func(descriptorFile string, values map[string]interface{}) descriptor.Parser {
		return testParser{}
	}
	newClusterParser = func(descriptorFile string, values map[string]interface{}) descriptor.Parser {
		return testParser{}
	}
	newProjectParser = func(projectFile string, values map[string]interface{}) descriptor.Parser {
		return testParser{}
	}
	newJobParser = func(descriptorFile string, values map[string]interface{}) descriptor.Parser {
		return testParser{}
	}
	newCronJobParser = func(descriptorFile string, values map[string]interface{}) descriptor.Parser {
		return testParser{}
	}
	newDeploymentParser = func(descriptorFile string, values map[string]interface{}) descriptor.Parser {
		return testParser{}
	}
	newDaemonSetParser = func(descriptorFile string, values map[string]interface{}) descriptor.Parser {
		return testParser{}
	}
	newStatefulSetParser = func(descriptorFile string, values map[string]interface{}) descriptor.Parser {
		return testParser{}
	}
}

type testConverger struct {
	expected bool
}

func (converger testConverger) Converge(bool) (result descriptor.ConvergeResult, err error) {
	if !converger.expected {
		return result, fmt.Errorf("Unexpected Call")
	}
	return
}

type testParser struct {
	expected     bool
	expectedData []byte
	t            *testing.T
}

func (parser testParser) Parse(data []byte, target interface{}) (err error) {
	if !parser.expected {
		return fmt.Errorf("Unexpected Call")
	}
	if len(parser.expectedData) > 0 {
		assert.Equals(parser.t, parser.expectedData, data)
	}
	err = yaml.Unmarshal(data, target)
	return
}

type testConfig struct {
	rancherURL   string
	insecureAPI  bool
	caCerts      string
	accessKey    string
	secretKey    string
	tokenKey     string
	clusterName  string
	clusterID    string
	mergeAnswers bool
	dryRun       bool
}

func (config testConfig) RancherURL() string {
	return config.rancherURL
}
func (config testConfig) InsecureAPI() bool {
	return config.insecureAPI
}
func (config testConfig) CACerts() string {
	return config.caCerts
}
func (config testConfig) AccessKey() string {
	return config.accessKey
}
func (config testConfig) SecretKey() string {
	return config.secretKey
}
func (config testConfig) TokenKey() string {
	return config.tokenKey
}
func (config testConfig) ClusterName() string {
	return config.clusterName
}
func (config testConfig) ClusterID() string {
	return config.clusterID
}
func (config testConfig) MergeAnswers() bool {
	return config.mergeAnswers
}
func (config testConfig) DryRun() bool {
	return config.dryRun
}
