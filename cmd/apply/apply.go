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
	"io/ioutil"
	"path/filepath"

	"github.com/bitgrip/cattlectl/cmd/utils"
	"github.com/bitgrip/cattlectl/internal/pkg/config"
	"github.com/bitgrip/cattlectl/internal/pkg/ctl"
	"github.com/bitgrip/cattlectl/internal/pkg/rancher/cluster/project"
	"github.com/bitgrip/cattlectl/internal/pkg/template"
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
	valuesFiles []string
	rootConfig  config.Config
	initCommand = func() {}
)

// used services
var (
	doApplyDescriptor = ctl.ApplyDescriptor
	newProjectParser  = project.NewProjectParser
)

// BaseCommand is accessor to the package base command
func BaseCommand(config config.Config, init func()) *cobra.Command {
	rootConfig = config
	initCommand = init
	return applyCmd
}

func apply(cmd *cobra.Command, args []string) {
	initCommand()
	values, err := utils.LoadValues(valuesFiles...)
	if err != nil {
		logrus.WithField("apply_file", applyFile).
			Fatal(err)
	}
	logrus.WithFields(values).Trace("Read descriptor with values")
	fileContent, err := ioutil.ReadFile(applyFile)
	if err != nil {
		logrus.WithField("apply_file", applyFile).
			Fatal(err)
	}
	projectData, err := template.BuildTemplate(fileContent, values, filepath.Dir(applyFile), false)
	if err != nil {
		logrus.WithField("apply_file", applyFile).
			Fatal(err)
	}

	result, err := doApplyDescriptor(applyFile, projectData, values, rootConfig)
	if err != nil {
		logrus.WithFields(values).
			WithField("apply_file", applyFile).
			Fatal(err)
	}
	logrus.
		WithField("upgraded-resouces", len(result.UpgradedResources)).
		WithField("created-resouces", len(result.CreatedResources)).
		Info("Finished Apply")
}

func init() {
	applyCmd.Flags().StringVarP(&applyFile, "file", "f", "project.yaml", "project file to apply")
	applyCmd.Flags().StringSliceVar(&valuesFiles, "values", []string{"values.yaml"}, "values file(s) to apply")

	applyCmd.Flags().Bool("merge-answers", false, "If answers of existing apps should be merged with the new apply answers")
	viper.BindPFlag("rancher.merge_answers", applyCmd.Flags().Lookup("merge-answers"))
	viper.BindEnv("rancher.merge_answers", "RANCHER_MERGE_ANSWERS")
}
