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

package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/bitgrip/cattlectl/internal/pkg/config"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

type AccessArgs struct {
	RancherURL  string `json:"rancher_url"`
	InsecureAPI bool   `json:"insecure_api"`
	CACerts     string `json:"ca_certs"`
	AccessKey   string `json:"access_key"`
	SecretKey   string `json:"secret_key"`
	ClusterName string `json:"cluster_name"`
	ConfigFile  string `json:"config_file"`
}

type BaseResponse struct {
	Msg     string `json:"msg"`
	Changed bool   `json:"changed"`
	Failed  bool   `json:"failed"`
}

func ExitJson(responseBody interface{}) {
	ReturnResponse(responseBody, false)
}

func FailJson(responseBody interface{}) {
	ReturnResponse(responseBody, true)
}

func ReadArguments(moduleArgs interface{}) {

	var response BaseResponse

	if len(os.Args) != 2 {
		response.Msg = "No argument file provided"
		response.Failed = true
		FailJson(response)
	}

	argsFile := os.Args[1]

	text, err := ioutil.ReadFile(argsFile)
	if err != nil {
		response.Msg = "Could not read configuration file: " + argsFile
		response.Failed = true
		FailJson(response)
	}

	err = json.Unmarshal(text, moduleArgs)
	if err != nil {
		response.Msg = "Configuration file not valid JSON: " + argsFile
		response.Failed = true
		FailJson(response)
	}
}

func ReturnResponse(responseBody interface{}, failed bool) {
	var response []byte
	var err error
	response, err = json.Marshal(responseBody)
	if err != nil {
		response, _ = json.Marshal(BaseResponse{Msg: "Invalid response object", Failed: true})
		failed = true
	}
	fmt.Println(string(response))
	if failed {
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}

func BuildRancherConfig(args AccessArgs) config.Config {
	if args.ConfigFile != "" {
		viper.SetConfigFile(args.ConfigFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			FailJson(BaseResponse{Msg: "Can not read configuration: " + err.Error()})
		}
		viper.AddConfigPath(home)
		viper.SetConfigName(".cattlectl")
	}
	viper.ReadInConfig()
	if args.RancherURL == "" {
		args.RancherURL = viper.GetString("rancher.url")
	}
	if args.InsecureAPI == false {
		args.InsecureAPI = viper.GetBool("rancher.insecure_api")
	}
	if args.CACerts == "" {
		args.CACerts = viper.GetString("rancher.ca_certs")
	}
	if args.AccessKey == "" {
		args.AccessKey = viper.GetString("rancher.access_key")
	}
	if args.SecretKey == "" {
		args.SecretKey = viper.GetString("rancher.secret_key")
	}
	if args.ClusterName == "" {
		args.ClusterName = viper.GetString("rancher.cluster_name")
	}
	return config.SimpleConfig(
		args.RancherURL,
		args.InsecureAPI,
		args.CACerts,
		args.AccessKey,
		args.SecretKey,
		args.ClusterName,
		"",
		false,
	)
}
