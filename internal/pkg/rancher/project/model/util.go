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

package model

import (
	"fmt"

	"k8s.io/apimachinery/pkg/util/intstr"
)

type IntOrString intstr.IntOrString

func (intorstr *IntOrString) UnmarshalYAML(unmarshal func(interface{}) error) error {
	if err := unmarshal(&intorstr.IntVal); err == nil {
		intorstr.Type = intstr.Int
		return nil
	}
	if err := unmarshal(&intorstr.StrVal); err != nil {
		return err
	}
	intorstr.Type = intstr.String
	return nil
}

func (intorstr IntOrString) MarshalYAML() (interface{}, error) {

	switch intorstr.Type {
	case intstr.Int:
		return intorstr.IntVal, nil
	case intstr.String:
		return intorstr.StrVal, nil
	default:
		return []byte{}, fmt.Errorf("impossible IntOrString.Type")
	}
}

// UnmarshalJSON implements the json.Unmarshaller interface.
func (intorstr *IntOrString) UnmarshalJSON(value []byte) error {
	return (*intstr.IntOrString)(intorstr).UnmarshalJSON(value)
}

// MarshalJSON implements the json.Marshaller interface.
func (intorstr IntOrString) MarshalJSON() ([]byte, error) {
	return (intstr.IntOrString)(intorstr).MarshalJSON()
}
