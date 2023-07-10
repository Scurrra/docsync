package dsconfig

import "os"

// Create an emty template of type `doctype` for the specified `lang` and `plangs`
func createEmptyTemplate(lang string, plangs []string, doctype DocType) error {
	err_dir := os.Mkdir(lang, os.ModePerm)
	if err_dir != nil {
		return err_dir
	}

	// TODO: template files

	return nil
}

// Create template for the specified `lang` and `plangs` using <base>/docs.<doctype> as reference
func createTemplateFromBase(lang string) error {
	err_dir := os.Mkdir(lang, os.ModePerm)
	if err_dir != nil {
		return err_dir
	}

	// TODO: template files

	return nil
}
