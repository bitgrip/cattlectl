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
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

// genDocCmd represents the genDoc command
var genDocCmd = &cobra.Command{
	Use:   "gen-doc [target folder]",
	Short: "genrates the markdown documentation",
	Long:  `Generates the full command tree documantation as markdown files inside the target folder.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := doc.GenMarkdownTree(rootCmd, args[0])
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("find generated documentation at %s/cattlctl.md\n", args[0])
	},
	DisableAutoGenTag: true,
}
