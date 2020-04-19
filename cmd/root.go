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

package cmd

import (
	"fmt"
	"os"

	"github.com/bitgrip/cattlectl/cmd/apply"
	"github.com/bitgrip/cattlectl/cmd/delete"
	"github.com/bitgrip/cattlectl/cmd/list"
	"github.com/bitgrip/cattlectl/cmd/show"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// rootCmd represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Use:   "cattlectl",
		Short: "controll your cattle on the ranch",
		Long: `cattlectl is a CLI controller to your rancher instance.

It allowes you to describe your full project as code and lett it apply to your
cluster.

cattlectl handles the input in a idempotend way so that it dosn't change your
deployement, if you run cattlectl twice.`,
		DisableAutoGenTag: true,
	}
	cfgFile  string
	logJson  bool
	LogLevel int
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {

	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cattlectl.yaml)")
	rootCmd.PersistentFlags().IntVarP(&LogLevel, "verbosity", "v", 0, "verbosity level to use")
	rootCmd.PersistentFlags().BoolVar(&logJson, "log-json", false, "if to log using json format")
	rootCmd.PersistentFlags().String("rancher-url", "", "The URL to reach the rancher")
	rootCmd.PersistentFlags().Bool("insecure-api", false, "If Rancher uses a self signed certificate")
	rootCmd.PersistentFlags().String("access-key", "", "The access key to access rancher with")
	rootCmd.PersistentFlags().String("secret-key", "", "The secret key to access rancher with")
	rootCmd.PersistentFlags().String("cluster-name", "", "The name of the cluster the project is part of")
	rootCmd.PersistentFlags().String("cluster-id", "", "The ID of the cluster the project is part of")
	viper.BindPFlag("rancher.url", rootCmd.PersistentFlags().Lookup("rancher-url"))
	viper.BindEnv("rancher.url", "RANCHER_URL")

	viper.BindPFlag("rancher.insecure_api", rootCmd.PersistentFlags().Lookup("insecure-api"))
	viper.BindEnv("rancher.insecure_api", "RANCHER_INSECURE_API")

	viper.BindEnv("rancher.ca_certs", "RANCHER_CA_CERTS")

	viper.BindPFlag("rancher.access_key", rootCmd.PersistentFlags().Lookup("access-key"))
	viper.BindEnv("rancher.access_key", "RANCHER_ACCESS_KEY")

	viper.BindPFlag("rancher.secret_key", rootCmd.PersistentFlags().Lookup("secret-key"))
	viper.BindEnv("rancher.secret_key", "RANCHER_SECRET_KEY")

	viper.BindPFlag("rancher.cluster_id", rootCmd.PersistentFlags().Lookup("cluster-id"))
	viper.BindEnv("rancher.cluster_id", "RANCHER_CLUSTER_ID")

	viper.BindPFlag("rancher.cluster_name", rootCmd.PersistentFlags().Lookup("cluster-name"))
	viper.BindEnv("rancher.cluster_name", "RANCHER_CLUSTER_NAME")

	rootCmd.AddCommand(apply.BaseCommand(rancherConfig, initSubCommand))
	rootCmd.AddCommand(delete.BaseCommand(rancherConfig, initSubCommand))
	rootCmd.AddCommand(list.BaseCommand(rancherConfig, initSubCommand))
	rootCmd.AddCommand(show.BaseCommand(rancherConfig, initSubCommand))
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(genDocCmd)
	rootCmd.AddCommand(completionCmd)
}

func initSubCommand() {
	if logJson {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	} else {
		logrus.SetFormatter(&logrus.TextFormatter{
			ForceColors: true,
		})
	}
	if LogLevel > 0 {
		logrus.SetLevel(logrus.DebugLevel)
	}
	if LogLevel >= 2 {
		logrus.SetLevel(logrus.TraceLevel)
		logrus.SetReportCaller(true)
	}
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		viper.AddConfigPath(home)
		viper.SetConfigName(".cattlectl")
	}

	viper.AutomaticEnv()

	viper.ReadInConfig()
}
