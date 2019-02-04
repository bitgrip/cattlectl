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

package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"

	"github.com/spf13/viper"
	yaml "gopkg.in/yaml.v2"
)

// LoadValues is reading values from a optional values file (YAML formated)
//
// The values are merged with corresponding environment variables.
func LoadValues(valuesFile string) (map[string]interface{}, error) {
	if err := verifyValuesFile(valuesFile); err != nil {
		return nil, err
	}
	valuesConfig := viper.New()
	valuesConfig.SetConfigFile(valuesFile)
	valuesConfig.AutomaticEnv()
	valuesConfig.ReadInConfig()
	for _, name := range viper.GetStringSlice("env_value_keys") {
		valuesConfig.BindEnv(name)
	}
	return valuesConfig.AllSettings(), nil
}

func verifyValuesFile(valuesFile string) error {
	expected := make(map[string]interface{}, 0)
	fileContent, err := ioutil.ReadFile(valuesFile)
	if err != nil && os.IsNotExist(err) {
		// A not existing values file is valid
		return nil
	} else if err != nil {
		return err
	}
	err = yaml.Unmarshal(fileContent, &expected)
	if err != nil {
		return err
	}
	valuesConfig := viper.New()
	valuesConfig.SetConfigFile(valuesFile)
	valuesConfig.ReadInConfig()
	structureFromViper, _ := yaml.Marshal(valuesConfig.AllSettings())
	actual := make(map[string]interface{}, 0)
	yaml.Unmarshal(structureFromViper, &actual)

	if !reflect.DeepEqual(expected, actual) {
		return fmt.Errorf("uppercase characters are not allowed on value keys")

	}
	return nil
}
