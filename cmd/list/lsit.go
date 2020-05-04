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

package list

import (
	"fmt"

	"github.com/bitgrip/cattlectl/internal/pkg/config"
	"github.com/bitgrip/cattlectl/internal/pkg/ctl"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	validArgs = []string{"app"}
	listCmd   = &cobra.Command{
		Use:       "list TYPE",
		Short:     "Lists an rancher resouce",
		Long:      "Lists an rancher resouce",
		Run:       list,
		ValidArgs: validArgs,
	}
	rootConfig  config.Config
	initCommand = func() {}
)

// BaseCommand is accessor to the package base command
func BaseCommand(config config.Config, init func()) *cobra.Command {
	rootConfig = config
	initCommand = init
	return listCmd
}

func list(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		logrus.Warn(cmd.UsageString())
		return
	}
	resouceType := args[0]
	projectName := viper.GetString("list_cmd.project_name")
	namespace := viper.GetString("list_cmd.namespace")
	pattern := viper.GetString("list_cmd.pattern")
	logrus.
		WithField("project-name", projectName).
		WithField("resouce-type", resouceType).
		WithField("cluster-name", rootConfig.ClusterName()).
		Debug("List project resouces")
	matches, err := ctl.ListProjectResouces(projectName, namespace, resouceType, pattern, rootConfig)
	if err != nil {
		logrus.
			WithField("project-name", projectName).
			WithField("resouce-type", resouceType).
			WithField("cluster-name", rootConfig.ClusterName()).
			Fatal(err)
	}
	for _, match := range matches {
		fmt.Println(match)
	}
}

func init() {

	listCmd.Flags().String("project-name", "", "The name of the project to list resouces from")
	viper.BindPFlag("list_cmd.project_name", listCmd.Flags().Lookup("project-name"))

	listCmd.Flags().String("namespace", "", "The namespace of the project to list resouces from")
	viper.BindPFlag("list_cmd.namespace", listCmd.Flags().Lookup("namespace"))

	listCmd.Flags().String("pattern", "", "Match pattern to filter resouce names")
	viper.BindPFlag("list_cmd.pattern", listCmd.Flags().Lookup("pattern"))
}
