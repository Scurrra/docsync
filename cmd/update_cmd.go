package cmd

import (
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	"github.com/Scurrra/docsync/dsconfig"
	"github.com/Scurrra/docsync/markup"
	"github.com/Scurrra/docsync/markup/md"
	. "github.com/eminarican/safetypes"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	updatStatus bool
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update hashkeys of documentation blocks in base documentation. Command should be called from the root documentation directory.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if !noInteract {
			if !cmd.Flags().Lookup("update-status").Changed {
				updatStatus, _ = strconv.ParseBool(promptGetSelect(
					promptContent{
						"'update-status' is not provided",
						"Update status: ",
						None[string](),
					},
					[]bool{true, false},
				))
			}
		}

		// read config from file
		data, err_file := ioutil.ReadFile("docsync.yaml")
		if err_file != nil {
			return err_file
		}

		// unmarshall config
		config := dsconfig.Config{}
		err_yaml := yaml.Unmarshal(data, &config)
		if err_yaml != nil {
			return err_yaml
		}

		err := filepath.Walk(config.Base,
			func(doc_path string, info fs.FileInfo, err error) error {
				if err != nil {
					return nil
				}

				if len(filepath.Ext(doc_path)) != 0 {
					format := dsconfig.DocType(filepath.Ext(doc_path))
					doc_f, err_f := ioutil.ReadFile(doc_path)
					if err_f != nil {
						return nil
					}

					switch format {
					case ".md":
						// here hashkeys will be recomputed here
						doc := md.ParseDocument(string(doc_f))

						if updatStatus {
							for key, block := range doc.Blocks {
								if block.Status == "New" {
									block.Status = markup.Active
								}
								doc.Blocks[key] = block
							}
						}

						// so we need just rewrite
						err_f = os.WriteFile(doc_path, []byte(md.RenderDocument(doc)), os.ModePerm)
						if err_f != nil {
							return err_f
						}
					}
				}

				return nil
			},
		)

		return err
	},
}
