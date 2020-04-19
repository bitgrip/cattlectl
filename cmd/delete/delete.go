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

package delete

import (
	"github.com/bitgrip/cattlectl/internal/pkg/config"
	"github.com/bitgrip/cattlectl/internal/pkg/ctl"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	validArgs = []string{"app"}
	deleteCmd = &cobra.Command{
		Use:       "delete TYPE NAME",
		Short:     "Deletes an rancher resouce",
		Long:      "Deletes an rancher resouce",
		Run:       delete,
		ValidArgs: validArgs,
	}
	rootConfig  config.Config
	initCommand = func() {}
)

// BaseCommand is accessor to the package base command
func BaseCommand(config config.Config, init func()) *cobra.Command {
	rootConfig = config
	initCommand = init
	return deleteCmd
}

func delete(cmd *cobra.Command, args []string) {
	if len(args) != 2 {
		logrus.Warn(cmd.UsageString())
		return
	}
	resouceType := args[0]
	resourceName := args[1]
	projectName := viper.GetString("rancher.project_name")
	logrus.
		WithField("project-name", projectName).
		WithField("resouce-type", resouceType).
		WithField("resouce-name", resourceName).
		WithField("cluster-name", rootConfig.ClusterName()).
		Info("Delete project resouce")
	err := ctl.DeleteProjectResouce(projectName, resouceType, resourceName, rootConfig)
	if err != nil {
		logrus.
			WithField("project-name", projectName).
			WithField("resouce-type", resouceType).
			WithField("resouce-name", resourceName).
			WithField("cluster-name", rootConfig.ClusterName()).
			Fatal(err)
	}
}

func init() {

	deleteCmd.PersistentFlags().String("project-name", "", "The name of the project to delete resouces from")
	viper.BindPFlag("rancher.project_name", deleteCmd.PersistentFlags().Lookup("project-name"))
	viper.BindEnv("rancher.project_name", "RANCHER_PROJECT_NAME")
}
