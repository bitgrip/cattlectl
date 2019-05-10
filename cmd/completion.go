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
	"os"

	"github.com/spf13/cobra"
)

// completionCmd represents the completion command
var completionCmd = &cobra.Command{
	Use:   "completion",
	Short: "Generates bash completion scripts",
	Long: `

To configure your bash shell to load completions for each session add to your bashrc

* ~/.bashrc or ~/.profile

        . <(cattlectl completion)

* On Mac (with bash completion installed from brew)

        cattlectl completion > $(brew --prefix)/etc/bash_completion.d/cattlectl

* To load completion run

        . <(cattlectl completion)

This will only temporaly activate completion in the current session.
`,
	Run: func(cmd *cobra.Command, args []string) {
		rootCmd.GenBashCompletion(os.Stdout)
	},
	DisableAutoGenTag: true,
}
