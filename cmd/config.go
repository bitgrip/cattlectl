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

package cmd

import (
	"fmt"

	"github.com/spf13/viper"
)

var rancherConfig = config{}

type config struct {
}

func (config) RancherURL() string {
	return viper.GetString("rancher.url")
}

func (config) InsecureAPI() bool {
	return viper.GetBool("rancher.insecure_api")
}

func (config) CACerts() string {
	return viper.GetString("rancher.ca_certs")
}

func (config) AccessKey() string {
	return viper.GetString("rancher.access_key")
}

func (config) SecretKey() string {
	return viper.GetString("rancher.secret_key")
}

func (c config) TokenKey() string {
	return fmt.Sprintf("%s:%s", c.AccessKey(), c.SecretKey())
}

func (config) ClusterName() string {
	return viper.GetString("rancher.cluster_name")
}

func (config) ClusterID() string {
	return viper.GetString("rancher.cluster_id")
}

func (config) MergeAnswers() bool {
	return viper.GetBool("rancher.merge_answers")
}
