package dsconfig

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/Scurrra/docsync/markup/md"
	"gopkg.in/yaml.v3"
)

// Create an emty template of type `doctype` for the specified `lang` and `plangs`
func CreateEmptyTemplate(dir_path, lang string, plangs []string, doctype DocType) error {
	err_dir := os.Mkdir(path.Join(dir_path, lang), os.ModePerm)
	if err_dir != nil {
		return err_dir
	}

	switch doctype {
	case ".md":
		doc := md.GenerateEmptyDocumentTemplate(plangs)
		doc_file := md.RenderDocument(doc)

		f, err_file := os.Create(path.Join(dir_path, lang, "index.md"))
		if err_file != nil {
			return err_file
		}
		defer f.Close()

		_, err_file = f.Write([]byte(doc_file))
		if err_file != nil {
			return err_file
		}
	}

	return nil
}

// Create template for the specified `lang` and `plangs` using <base>/docs<doctype> as reference
func CreateTemplateFromBase(baseLang, lang string, doctype DocType) error {
	err_dir := os.Mkdir(lang, os.ModePerm)
	if err_dir != nil {
		return err_dir
	}

	err := filepath.Walk(baseLang,
		func(doc_path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if filepath.Ext(doc_path) == string(doctype) {
				doc_f, err_f := ioutil.ReadFile(doc_path)
				if err_f != nil {
					return nil
				}

				buf := strings.Index(doc_path, string(os.PathSeparator)) + 1
				doc_path = doc_path[buf:]

				switch doctype {
				case ".md":
					doc := md.ParseDocument(string(doc_f))

					f, err_f := os.Create(path.Join(lang, doc_path))
					if err_f != nil {
						return err_f
					}
					defer f.Close()

					doc = md.GenerateDocumentTemplateBase(doc)

					_, err_f = f.Write([]byte(md.RenderDocument(doc)))
					if err_f != nil {
						return err_f
					}
				}
			}

			return nil
		})

	return err
}

func SyncLanguage(lang string) error {
	// read config from file
	data, err_file := ioutil.ReadFile("docsync.yaml")
	if err_file != nil {
		return err_file
	}

	// unmarshall config
	config := Config{}
	err_yaml := yaml.Unmarshal(data, &config)
	if err_yaml != nil {
		return err_yaml
	}

	err := filepath.Walk(config.Base,
		func(doc_base_path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			doc_type := filepath.Ext(doc_base_path)
			if len(doc_type) < 2 {
				return nil
			}

			buf := strings.Index(doc_base_path, string(os.PathSeparator)) + 1
			doc_lang_path := path.Join(lang, doc_base_path[buf:])
			if _, err := os.Stat(doc_lang_path); err != nil {
				// create template from base
				switch doc_type {
				case ".md":
					doc_f, err_f := ioutil.ReadFile(doc_base_path)
					if err_f != nil {
						return err_f
					}

					doc := md.ParseDocument(string(doc_f))

					f, err_f := os.Create(doc_lang_path)
					if err_f != nil {
						return err_f
					}
					defer f.Close()

					doc = md.GenerateDocumentTemplateBase(doc)

					_, err_f = f.Write([]byte(md.RenderDocument(doc)))
					if err_f != nil {
						return err_f
					}
				}
			} else {
				// merge files
				switch doc_type {
				case ".md":
					// read base file
					doc_f, err_f := ioutil.ReadFile(doc_base_path)
					if err_f != nil {
						return err_f
					}
					doc_base := md.ParseDocument(string(doc_f))

					// read lang file
					doc_f, err_f = os.ReadFile(doc_lang_path)
					if err_f != nil {
						return err_f
					}
					doc_lang := md.ParseDocument(string(doc_f))

					doc := md.GenerateMergedDocument(doc_base, doc_lang)

					// write file
					f, err_f := os.OpenFile(doc_lang_path, os.O_WRONLY, 0600)
					if err_f != nil {
						fmt.Printf("Error on openning %s\n", doc_lang_path)
						return err_f
					}
					defer f.Close()
					_, err_f = f.Write([]byte(md.RenderDocument(doc)))
					if err_f != nil {
						fmt.Printf("Error on writing %s\n", doc_lang_path)
						return err_f
					}
				}
			}

			return nil
		})

	return err
}
