/*
Copyright © 2023 Scurrra (Ilya Borowski) <iscurrra@gmail.com>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	version    = "0.0.1"
	noInteract bool
)

var rootCmd = &cobra.Command{
	Use:     "docsync",
	Version: version,
	Short:   "docsync is a tool helping with translating your documentation into other languages",
	Long: `docsync is a documentation synchronization tool.
	
docsync creates a template for the new translation of documentation to your project into another language.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(updateCmd)

	rootCmd.PersistentFlags().BoolVar(&noInteract, "no-interact", false, "Ask for missing flags in interactive mod or not.")

	initCmd.Flags().StringVar(&docsPath, "path", ".", "Path where docs will be placed. '.' means the current directory.")
	initCmd.Flags().StringVar(&docsMainType, "type", "md", "The main documentation files' type.")
	initCmd.Flags().StringVar(&baseLang, "lang", "en", "The base language of the documentation. Please, specify ISO639-1 code.")
	initCmd.Flags().StringSliceVar(&programmingLangs, "plangs", []string{}, "Programming languages of code from the documentations.")

	createCmd.Flags().StringVarP(&lang, "lang", "l", "", "New documentation language.")
	createCmd.Flags().BoolVar(&createFromBase, "create-from-base", true, "Create documentation using base as template.")
	createCmd.Flags().BoolVar(&createEmpty, "create-empty", false, "Create an empty documentation.")
}
