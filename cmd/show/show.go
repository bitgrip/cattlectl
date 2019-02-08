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

package show

import (
	"fmt"
	"log"

	"github.com/bitgrip/cattlectl/cmd/utils"
	"github.com/bitgrip/cattlectl/internal/pkg/config"
	"github.com/bitgrip/cattlectl/internal/pkg/rancher/project"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
)

var (
	showCmd = &cobra.Command{
		Use:   "show",
		Short: "Show the resulting project descriptor",
		Long:  showLongDescription,
		Run:   show,
	}
	showFile    string
	valuesFile  string
	rootConfig  config.Config
	initCommand func()
)
var (
	newProjectParser = project.NewPrettyParser
)

func BaseCommand(config config.Config, init func()) *cobra.Command {
	rootConfig = config
	initCommand = init
	return showCmd
}

func show(cmd *cobra.Command, args []string) {
	initCommand()
	values, err := utils.LoadValues(valuesFile)
	if err != nil {
		log.Fatal(err)
	}
	if project, err := newProjectParser(showFile).Parse(values); err != nil {
		log.Fatal(err)
	} else if out, err := yaml.Marshal(project); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(string(out))
	}
}

func init() {
	showCmd.Flags().StringVarP(&showFile, "file", "f", "project.yaml", "project file to show")
	showCmd.Flags().StringVar(&valuesFile, "values", "values.yaml", "values file to show")
}
