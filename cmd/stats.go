/*
Copyright © 2026 NBRGLM Developers Private Limited <support@nbrglm.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"os"
	"path"

	"github.com/goccy/go-yaml"
	"github.com/nbrglm/napiway/spec"
	"github.com/spf13/cobra"
)

// statsCmd represents the stats command
var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Prints some basic stats about the specification file and checks if it is valid.",
	Long:  `The stats command reads the specification file and prints some basic stats about it, such as the number of endpoints and schemas. It also checks if the specification file is valid and can be parsed successfully.`,
	Run: func(cmd *cobra.Command, args []string) {
		if !cmd.Flags().Changed("config") {
			cmd.Help()
			return
		}

		configFile, _ := cmd.Flags().GetString("config")

		bytes, err := os.ReadFile(configFile)
		if err != nil {
			panic(err)
		}

		var cfg *spec.Specification = new(spec.Specification)
		err = yaml.Unmarshal(bytes, cfg)
		if err != nil {
			panic(err)
		}

		if err := cfg.Validate(); err != nil {
			panic(err)
		}

		// get absolute path of the config file
		err = os.Chdir(path.Dir(configFile))
		if err != nil {
			cmd.PrintErrf("Failed to change directory to config file directory: %v\nTry passing in the absolute path to the config file.", err)
			os.Exit(1)
		}

		// Print some stats about the spec
		cmd.Printf("Specification loaded successfully!\n")
		cmd.Printf("Number of endpoints: %d\n", len(cfg.Endpoints))
		cmd.Printf("Number of schemas: %d\n", len(cfg.Schemas))

		if cfg.GoServer != nil {
			cmd.Println("Go server helpers generation is enabled!")
			cmd.Println("Output directory:", cfg.GoServer.OutputDir)
		}

		if cfg.GoSDK != nil {
			cmd.Println("Go SDK generation is enabled!")
			cmd.Println("Output directory:", cfg.GoSDK.OutputDir)
		}

		if cfg.TsSDK != nil {
			cmd.Println("TypeScript SDK generation is enabled!")
			cmd.Println("Output directory:", cfg.TsSDK.OutputDir)
		}

		// TODO: docs generation stats
	},
}

func init() {
	statsCmd.Flags().String("config", "", "config file to use (yaml only)")

	rootCmd.AddCommand(statsCmd)
}
