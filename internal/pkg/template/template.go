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
	"strings"
	"text/template"
)

/*
BuildTemplate builds and executes the given templateData with the given values.

templateData - The data to build the template from
values       - The values to use for execution of the template
trancated    - If base64 content has to be truncated
*/
func BuildTemplate(templateData []byte, values map[string]interface{}, truncated bool) ([]byte, error) {
	projectTemplate := template.New("data-template")
	projectTemplate.Funcs(template.FuncMap{
		"read":   read,
		"indent": indent,
	})
	if truncated {
		projectTemplate.Funcs(template.FuncMap{
			"base64": toTruncatedBase64,
		})
	} else {
		projectTemplate.Funcs(template.FuncMap{
			"base64": toBase64,
		})
	}
	projectTemplate, err := projectTemplate.Parse(string(templateData))
	if err != nil {
		return []byte{}, err
	}
	projectTemplate = projectTemplate.Option("missingkey=error")
	var templateBuffer bytes.Buffer
	if err := projectTemplate.Execute(&templateBuffer, values); err != nil {
		return []byte{}, err
	}

	return templateBuffer.Bytes(), nil

}

func read(fileName string) []byte {
	fileContent, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}
	return fileContent
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
