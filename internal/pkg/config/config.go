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

package config

import "fmt"

// Config provides rancher access informations
type Config interface {
	RancherURL() string
	InsecureAPI() bool
	CACerts() string
	AccessKey() string
	SecretKey() string
	TokenKey() string
	ClusterName() string
	ClusterID() string
	MergeAnswers() bool
}

func SimpleConfig(
	rancherURL string,
	insecureAPI bool,
	caCerts string,
	accessKey string,
	secretKey string,
	clusterName string,
	clusterID string,
	mergeAnswers bool,
) Config {
	return simpleConfig{
		rancherURL:   rancherURL,
		insecureAPI:  insecureAPI,
		caCerts:      caCerts,
		accessKey:    accessKey,
		secretKey:    secretKey,
		clusterName:  clusterName,
		clusterID:    clusterID,
		mergeAnswers: mergeAnswers,
	}
}

type simpleConfig struct {
	rancherURL   string
	insecureAPI  bool
	caCerts      string
	accessKey    string
	secretKey    string
	clusterName  string
	clusterID    string
	mergeAnswers bool
}

func (config simpleConfig) RancherURL() string {
	return config.rancherURL
}

func (config simpleConfig) InsecureAPI() bool {
	return config.insecureAPI
}

func (config simpleConfig) CACerts() string {
	return config.caCerts
}

func (config simpleConfig) AccessKey() string {
	return config.accessKey
}

func (config simpleConfig) SecretKey() string {
	return config.secretKey
}

func (config simpleConfig) TokenKey() string {
	return fmt.Sprintf("%s:%s", config.AccessKey(), config.SecretKey())
}

func (config simpleConfig) ClusterName() string {
	return config.clusterName
}

func (config simpleConfig) ClusterID() string {
	return config.clusterID
}

func (config simpleConfig) MergeAnswers() bool {
	return config.mergeAnswers
}
