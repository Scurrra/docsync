package dsconfig

import (
	"os"
	"path"

	"github.com/Scurrra/docsync/markup/md"
)

// Create an emty template of type `doctype` for the specified `lang` and `plangs`
func CreateEmptyTemplate(dir_path, lang string, plangs []string, doctype DocType) error {
	err_dir := os.Mkdir(path.Join(dir_path, lang), os.ModePerm)
	if err_dir != nil {
		return err_dir
	}

	switch doctype {
	case "md":
		doc := md.GenerateEmptyDocumentTemplateMarkdown(plangs)
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

// Create template for the specified `lang` and `plangs` using <base>/docs.<doctype> as reference
func CreateTemplateFromBase(lang string) error {
	err_dir := os.Mkdir(lang, os.ModePerm)
	if err_dir != nil {
		return err_dir
	}

	// TODO: template files

	return nil
}
