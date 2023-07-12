package dsconfig

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/Scurrra/docsync/markup/md"
)

// Create an emty template of type `doctype` for the specified `lang` and `plangs`
func CreateEmptyTemplate(path, lang string, plangs []string, doctype DocType) error {
	err_dir := os.Mkdir(path+lang, os.ModePerm)
	if err_dir != nil {
		return err_dir
	}

	switch doctype {
	case MD:
		doc := md.GenerateEmptyDocumentTemplateMarkdown(plangs)
		doc_file := md.RenderDocument(doc)

		err_file := ioutil.WriteFile(fmt.Sprintf("%s.md", lang), []byte(doc_file), 0)
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
