package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Scurrra/docsync/dsconfig"
	. "github.com/eminarican/safetypes"
	"github.com/spf13/cobra"
)

var (
	lang           string
	createFromBase bool
	createEmpty    bool
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create new documentation from base lang. Command should be called from the root documentation directory.",
	RunE: func(cmd *cobra.Command, args []string) error {
		lang = strings.Trim(lang, " \t\n")
		if noInteract && len(lang) != 0 {
			return dsconfig.AddLanguage(lang, createFromBase, createEmpty)
		}

		if !cmd.Flags().Lookup("lang").Changed {
			lang = promptGetInput(
				promptContent{
					"New documentation language is not provided",
					"New documentation language: ",
					None[string](),
				},
			)
		}

		if !cmd.Flags().Lookup("create-from-base").Changed {
			createFromBase, _ = strconv.ParseBool(promptGetSelect(
				promptContent{
					"'create-from-base' is not provided",
					"Create from base: ",
					None[string](),
				},
				[]bool{true, false},
			))
		}

		if !cmd.Flags().Lookup("create-empty").Changed && !createFromBase {
			createEmpty, _ = strconv.ParseBool(promptGetSelect(
				promptContent{
					"'create-empty' is not provided",
					"Create empty template: ",
					None[string](),
				},
				[]bool{true, false},
			))
		}

		if !createFromBase && !createEmpty {
			fmt.Println("'create-from-base' and 'create-empty' can not both be false at the time.")
			return cmd.Context().Err()
		}

		return dsconfig.AddLanguage(lang, createFromBase, createEmpty)
	},
}
