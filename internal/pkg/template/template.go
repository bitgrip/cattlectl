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

package template

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	"text/template"

	yaml "gopkg.in/yaml.v2"
)

/*
BuildTemplate builds and executes the given templateData with the given values.

templateData - The data to build the template from
values       - The values to use for execution of the template
trancated    - If base64 content has to be truncated
*/
func BuildTemplate(templateData []byte, values map[string]interface{}, baseDir string, truncated bool) ([]byte, error) {
	absBaseDir, err := filepath.Abs(baseDir)
	if err != nil {
		return []byte{}, err
	}
	descriptorTemplate := template.New("data-template")
	descriptorTemplate.Funcs(template.FuncMap{
		"read":             readFunc(absBaseDir),
		"readWithTemplate": readTemplateFunc(baseDir, values),
		"indent":           indent,
		"trim":             trim,
		"toYaml":           toYaml,
	})
	if truncated {
		descriptorTemplate.Funcs(template.FuncMap{
			"base64": toTruncatedBase64,
		})
	} else {
		descriptorTemplate.Funcs(template.FuncMap{
			"base64": toBase64,
		})
	}
	descriptorTemplate, err = descriptorTemplate.Parse(string(templateData))
	if err != nil {
		return []byte{}, err
	}
	descriptorTemplate = descriptorTemplate.Option("missingkey=error")
	var templateBuffer bytes.Buffer
	if err := descriptorTemplate.Execute(&templateBuffer, values); err != nil {
		return []byte{}, err
	}

	return templateBuffer.Bytes(), nil

}

func readFunc(baseDir string) func(fileName string) []byte {
	return func(fileName string) []byte {
		var absFileName string
		if filepath.IsAbs(fileName) {
			absFileName = fileName
		} else {
			absFileName = filepath.Clean(fmt.Sprintf("%s/%s", baseDir, fileName))
		}
		fileContent, err := ioutil.ReadFile(absFileName)
		if err != nil {
			log.Fatal(err)
		}
		return fileContent
	}
}

func readTemplateFunc(baseDir string, values map[string]interface{}) func(fileName string) string {
	return func(fileName string) string {

		var absFileName string
		if filepath.IsAbs(fileName) {
			absFileName = fileName
		} else {
			absFileName = filepath.Clean(fmt.Sprintf("%s/%s", baseDir, fileName))
		}

		fileContent, err := ioutil.ReadFile(absFileName)
		if err != nil {
			log.Fatal(err)
		}

		descriptorTemplate := template.New(absFileName)
		descriptorTemplate.Funcs(template.FuncMap{
			"read":   readFunc(filepath.Dir(absFileName)),
			"indent": indent,
			"toYaml": toYaml,
			"base64": toBase64,
		})
		descriptorTemplate, err = descriptorTemplate.Parse(string(fileContent))
		if err != nil {
			log.Fatal(err)
		}
		descriptorTemplate = descriptorTemplate.Option("missingkey=error")
		var templateBuffer bytes.Buffer
		if err := descriptorTemplate.Execute(&templateBuffer, values); err != nil {
			log.Fatal(err)
		}
		return string(templateBuffer.Bytes())
	}
}

func toBase64(data interface{}) string {
	var encoded string
	if bytes, isBytes := data.([]byte); isBytes {
		encoded = base64.StdEncoding.EncodeToString(bytes)
	} else {
		encoded = base64.StdEncoding.EncodeToString([]byte(fmt.Sprint(data)))
	}
	return encoded
}

func toTruncatedBase64(data interface{}) string {
	if bytes, isBytes := data.([]byte); isBytes {
		return fmt.Sprintf("< %v bytes base64 encoded >", len(bytes))
	}
	return fmt.Sprintf("< %v bytes base64 encoded >", len([]byte(fmt.Sprint(data))))
}

func indent(indents int, data interface{}) string {
	prefix := strings.Repeat("  ", indents)
	var toIndent string
	if bytes, isBytes := data.([]byte); isBytes {
		toIndent = string(bytes)
	} else {
		toIndent = fmt.Sprint(data)
	}
	result := strings.TrimSpace(strings.Join(strings.Split(toIndent, "\n"), "\n"+prefix))
	return prefix + result
}

func trim(data interface{}) string {
	var toTrim string
	if bytes, isBytes := data.([]byte); isBytes {
		toTrim = string(bytes)
	} else {
		toTrim = fmt.Sprint(data)
	}
	return strings.TrimSpace(toTrim)
}

func toYaml(data interface{}) (string, error) {
	marshalledYaml, err := yaml.Marshal(data)
	if err != nil {
		// Swallow errors inside of a template.
		return "", err
	}
	return string(marshalledYaml), nil
}
