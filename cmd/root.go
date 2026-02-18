/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"
	"path"

	"github.com/goccy/go-yaml"
	"github.com/nbrglm/napiway/generators/golang"
	"github.com/nbrglm/napiway/generators/typescript"
	"github.com/nbrglm/napiway/spec"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "napiway",
	Short: "NapiWay is a Server and SDK generator for building APIs in Go and SDKs in multiple languages.",
	Long: `NapiWay is a Server and SDK generator for building APIs in Go and SDKs in multiple languages.
	Run this command as 'napiway --config <path-to-config-file>' to generate server and SDK code.

	Find more information at https://github.com/nbrglm/napiway`,
	Run: generate,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.napiway.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().String("config", "", "config file to use (yaml only)")
}

func generate(cmd *cobra.Command, args []string) {
	if !cmd.Flags().Changed("config") {
		cmd.Help()
		return
	}

	configFile, _ := cmd.Flags().GetString("config")

	bytes, err := os.ReadFile(configFile)
	if err != nil {
		panic(err)
	}

	var cfg *spec.Config = new(spec.Config)
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

	if cfg.GoServer != nil {
		err = golang.GenerateServerHelpers(*cfg.GoServer, cfg.Spec)
		if err != nil {
			panic(err)
		}
		cmd.Println("Go server helpers generated successfully!")
	}

	if cfg.GoSDK != nil {
		err = golang.GenerateGoSDK(*cfg.GoSDK, cfg.Spec)
		if err != nil {
			panic(err)
		}
		cmd.Println("Go SDK generated successfully!")
	}

	// TODO: Add TypeScript SDK generation
	if cfg.TsSDK != nil {
		err := typescript.GenerateTSSDK(*cfg.TsSDK, cfg.Spec)
		if err != nil {
			panic(err)
		}
		cmd.Println("TypeScript SDK generated successfully!")
	}

	// TODO: Add Docs generation (including README for the SDKs)
}
