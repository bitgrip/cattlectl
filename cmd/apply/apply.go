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

package apply

import (
	"github.com/bitgrip/cattlectl/internal/pkg/config"
	"github.com/bitgrip/cattlectl/internal/pkg/ctl"
	"github.com/bitgrip/cattlectl/internal/pkg/rancher/project"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	applyCmd = &cobra.Command{
		Use:   "apply",
		Short: "Apply a project descriptor to your rancher",
		Long:  applyLongDescription,
		Run:   apply,
	}
	applyFile   string
	valuesFile  string
	rootConfig  config.Config
	initCommand = func() {}
)

// used services
var (
	doApply          = ctl.ApplyProject
	newProjectParser = project.NewParser
	loadValues       = func() map[string]interface{} {
		valuesConfig := viper.New()
		valuesConfig.SetConfigFile(valuesFile)
		valuesConfig.AutomaticEnv()
		valuesConfig.ReadInConfig()
		for _, name := range viper.GetStringSlice("env_value_keys") {
			valuesConfig.BindEnv(name)
		}
		return valuesConfig.AllSettings()
	}
)

func BaseCommand(config config.Config, init func()) *cobra.Command {
	rootConfig = config
	initCommand = init
	return applyCmd
}

func apply(cmd *cobra.Command, args []string) {
	initCommand()
	values := loadValues()
	logrus.WithFields(values).Debug("Read descriptor with values")
	if project, err := newProjectParser(applyFile).Parse(values); err != nil {
		logrus.WithFields(values).
			WithField("apply_file", applyFile).
			Fatal(err)
	} else {
		err := doApply(project, rootConfig)
		if err != nil {
			logrus.WithFields(values).
				WithField("apply_file", applyFile).
				Fatal(err)
		}
	}
}

func init() {
	applyCmd.Flags().StringVarP(&applyFile, "file", "f", "project.yaml", "project file to apply")
	applyCmd.Flags().StringVar(&valuesFile, "values", "values.yaml", "values file to apply")
}
