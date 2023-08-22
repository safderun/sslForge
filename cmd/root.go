/*
Copyright © 2023 NAME HERE crtforge@burakberk.dev
*/
package cmd

import (
	"crtforge/cmd/services"
	_ "embed"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var caName string
var intermediateCaName string
var version = "v1.0.0"
var commitId = "abcd"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "crtforge",
	Short:   "Be a local cert authority",
	Long:    `With crtforge, you can create root, intermediate and application ca.`,
	Version: version + " " + commitId,
	Run:     rootRun,
}

func rootRun(cmd *cobra.Command, args []string) {
	appName := args[0]
	appDomains := args[1:]

	homeDirectory, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("home directory couldn't find", err)
	}
	configDirectory := homeDirectory + "/.config/crtforge"
	createConfigDir(configDirectory)
	defaultCADir := services.CreateCaDir(configDirectory, caName)

	defaultCARootCACrt, defaultCARootCACnf, defaultCARootCAkey := services.CreateRootCa(defaultCADir)
	_ = defaultCARootCAkey

	defaultCAIntermediateCACrt, defaultCAIntermediateCACnf, defaultCAIntermediateCAkey := services.CreateIntermediateCa(defaultCADir, intermediateCaName, defaultCARootCACnf)

	services.CreateAppCrt(defaultCADir, defaultCAIntermediateCACnf, defaultCAIntermediateCACrt, defaultCAIntermediateCAkey, defaultCARootCACrt, appName, appDomains[0], appDomains)
}

func createConfigDir(configDir string) {
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		err := os.MkdirAll(configDir, 0700)
		if err != nil {
			log.Fatal(err)
		}
	}
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
	rootCmd.Flags().BoolP("version", "v", false, "Print version information")

	// rootCmd.Flags().BoolVar(&versionFlag, "version", false, "Print version")
	rootCmd.Flags().StringVarP(&caName, "root-ca", "r", "default", "Set CA Name")
	rootCmd.Flags().StringVarP(&intermediateCaName, "intermediate-ca", "i", "intermediateCA", "Set Intermediate CA Name")

	// Example usages:
	rootCmd.Example = `Generate a cert under the default root and the default intermediate ca: 
./crtforge crtforgeapp crtforge.com app.crtforge.com api.crtforge.com [flags]

Generate a cert under a root ca named medical and a intermediate ca named frontend:
./crtforge crtforgeapp -r medical -i frontend crtforge.com app.crtforge.com api.crtforge.com [flags]`

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
