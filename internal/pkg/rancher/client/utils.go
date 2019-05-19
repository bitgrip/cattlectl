// Copyright © 2019 Bitgrip <berlin@bitgrip.de>
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

package client

import (
	"crypto/sha1"
	"fmt"
	"strings"

	"github.com/rancher/norman/clientbase"
	backendClusterClient "github.com/rancher/types/client/cluster/v3"
	backendRancherClient "github.com/rancher/types/client/management/v3"
	backendProjectClient "github.com/rancher/types/client/project/v3"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var (
	newBackendClusterClient = backendClusterClient.NewClient
	newManagementClient     = backendRancherClient.NewClient
	newProjectClient        = backendProjectClient.NewClient
)

func createClientOpts(config RancherConfig) *clientbase.ClientOpts {
	serverURL := config.RancherURL

	if !strings.HasSuffix(serverURL, "/v3") {
		serverURL = serverURL + "/v3"
	}

	if config.Insecure {
		logrus.WithFields(logrus.Fields{"rancher.url": serverURL}).
			Warn("Without SSL verification")
	}
	if config.CACerts != "" {
		logrus.WithFields(logrus.Fields{"rancher.url": serverURL, "rancher.ca_certs": config.CACerts}).
			Trace("Whith custom CA")
	}

	options := &clientbase.ClientOpts{
		URL:       serverURL,
		AccessKey: config.AccessKey,
		SecretKey: config.SecretKey,
		Insecure:  config.Insecure,
		CACerts:   config.CACerts,
	}
	return options
}

func createBackendClusterClient(config RancherConfig, clusterID string) (*backendClusterClient.Client, error) {
	logrus.WithFields(logrus.Fields{
		"rancher.url":        config.RancherURL,
		"rancher.cluster_id": clusterID,
	}).Debug("Create Cluster Client")
	options := createClientOpts(config)
	options.URL = options.URL + "/cluster/" + clusterID
	backendClusterClient, err := newBackendClusterClient(options)
	if err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{
			"rancher.url":        config.RancherURL,
			"rancher.cluster_id": clusterID,
		}).Error("Failed to create cluster client")
		return nil, fmt.Errorf("Failed to create cluster client, %v", err)
	}
	return backendClusterClient, nil
}

func createManagementClient(config RancherConfig) (*backendRancherClient.Client, error) {
	logrus.WithFields(logrus.Fields{
		"rancher.url": config.RancherURL,
	}).Debug("Create Management Client")
	options := createClientOpts(config)
	backendRancherClientCache, err := newManagementClient(options)
	if err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{
			"rancher.url": config.RancherURL,
		}).Error("Failed to create management client")
		return nil, fmt.Errorf("Failed to create management client, %v", err)
	}
	if backendRancherClientCache == nil || backendRancherClientCache.APIBaseClient.Ops == nil {
		logrus.WithFields(logrus.Fields{
			"rancher.url":      config.RancherURL,
			"rancher.ca_certs": config.CACerts,
		}).Error("Failed to create management client")
		return nil, fmt.Errorf("Failed to create management client")
	}

	return backendRancherClientCache, nil
}

func createProjectClient(config RancherConfig, clusterID string, projectID string) (*backendProjectClient.Client, error) {
	logrus.WithFields(logrus.Fields{
		"rancher.url":        config.RancherURL,
		"rancher.cluster_id": clusterID,
		"project_id":         projectID,
	}).Debug("Create Project Client")
	options := createClientOpts(config)
	options.URL = options.URL + "/projects/" + projectID

	pc, err := newProjectClient(options)
	if err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{
			"rancher.url":        config.RancherURL,
			"rancher.cluster_id": clusterID,
			"project_id":         projectID,
		}).Error("Failed to create project client")
		return nil, fmt.Errorf("Failed to create project client, %v", err)
	}
	return pc, nil
}

func hashOf(data interface{}) string {

	dataBytes, _ := yaml.Marshal(data)
	h := sha1.New()
	h.Write(dataBytes)
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}
