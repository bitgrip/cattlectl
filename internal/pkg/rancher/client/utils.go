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
	"fmt"
	"strings"

	"github.com/rancher/norman/clientbase"
	clusterClient "github.com/rancher/types/client/cluster/v3"
	managementClient "github.com/rancher/types/client/management/v3"
	projectClient "github.com/rancher/types/client/project/v3"
	"github.com/sirupsen/logrus"
)

var (
	newClusterClient    = clusterClient.NewClient
	newManagementClient = managementClient.NewClient
	newProjectClient    = projectClient.NewClient
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

func createClusterClient(config RancherConfig, clusterID string) (*clusterClient.Client, error) {
	logrus.WithFields(logrus.Fields{
		"rancher.url":        config.RancherURL,
		"rancher.cluster_id": clusterID,
	}).Debug("Create Cluster Client")
	options := createClientOpts(config)
	options.URL = options.URL + "/cluster/" + clusterID
	clusterClient, err := newClusterClient(options)
	if err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{
			"rancher.url":        config.RancherURL,
			"rancher.cluster_id": clusterID,
		}).Error("Failed to create cluster client")
		return nil, fmt.Errorf("Failed to create cluster client, %v", err)
	}
	return clusterClient, nil
}

func createManagementClient(config RancherConfig) (*managementClient.Client, error) {
	logrus.WithFields(logrus.Fields{
		"rancher.url": config.RancherURL,
	}).Debug("Create Management Client")
	options := createClientOpts(config)
	managementClientCache, err := newManagementClient(options)
	if err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{
			"rancher.url": config.RancherURL,
		}).Error("Failed to create management client")
		return nil, fmt.Errorf("Failed to create management client, %v", err)
	}
	if managementClientCache == nil || managementClientCache.APIBaseClient.Ops == nil {
		logrus.WithFields(logrus.Fields{
			"rancher.url":      config.RancherURL,
			"rancher.ca_certs": config.CACerts,
		}).Error("Failed to create management client")
		return nil, fmt.Errorf("Failed to create management client")
	}

	return managementClientCache, nil
}

func createProjectClient(config RancherConfig, clusterID string, projectID string) (*projectClient.Client, error) {
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
