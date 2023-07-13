package cmd

import (
	"strings"

	"github.com/Scurrra/docsync/dsconfig"
	. "github.com/eminarican/safetypes"
	"github.com/spf13/cobra"
)

var (
	noInteract       bool
	docsPath         string
	docsMainType     string
	baseLang         string
	programmingLangs []string
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize new documentation",
	RunE: func(cmd *cobra.Command, args []string) error {
		if noInteract {
			return dsconfig.NewConfig(
				docsPath,
				baseLang,
				programmingLangs,
				dsconfig.NewFormatConfig(dsconfig.MarkdownConfig{}),
				true,
			)
		}

		if !cmd.Flags().Lookup("path").Changed {
			docsPath = promptGetInput(
				promptContent{
					"Documentation path is not provided",
					"Documentation path: ",
					Some("."),
				},
			)
		}

		if !cmd.Flags().Lookup("type").Changed {
			docsMainType = promptGetSelect(
				promptContent{
					"Main type of the documentation files is not provided",
					"Main type of the documentation files: ",
					Some("md"),
				},
				[]string{"md"},
			)
		}

		if !cmd.Flags().Lookup("lang").Changed {
			baseLang = promptGetInput(
				promptContent{
					"Base language of the documentation is not provided",
					"Base language of the documentation (iso639-1 code): ",
					None[string](),
				},
			)
		}

		if !cmd.Flags().Lookup("plangs").Changed {
			programmingLangs = strings.Split(
				promptGetInput(
					promptContent{
						"Programming language of the code snippets is not provided",
						"Programming languages of the code snippets: ",
						None[string](),
					},
				),
				",",
			)
		}

		return dsconfig.NewConfig(
			docsPath,
			baseLang,
			programmingLangs,
			dsconfig.NewFormatConfig(dsconfig.MarkdownConfig{}),
			true,
		)
	},
}
