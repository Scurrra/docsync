package cmd

import (
	"strings"

	"github.com/Scurrra/docsync/dsconfig"
	. "github.com/eminarican/safetypes"
	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Synchronizes specified language documentation with base. Command should be called from the root documentation directory.",
	RunE: func(cmd *cobra.Command, args []string) error {
		lang = strings.Trim(lang, " \t\n")
		if !noInteract || len(lang) == 0 {
			if !cmd.Flags().Lookup("lang").Changed {
				lang = promptGetInput(
					promptContent{
						"Documentation language for synchronization is not provided",
						"Documentation language for synchronization: ",
						None[string](),
					},
				)
			}
		}

		return dsconfig.SyncLanguage(lang)
	},
}
