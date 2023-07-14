package dsconfig

import (
	"io/ioutil"
	"os"
	"path"

	iso6391 "github.com/emvi/iso-639-1"
	"gopkg.in/yaml.v3"
)

type DocType string

const (
	MD DocType = ".md"
)

// Configuration settings for markdown documentation
type MarkdownConfig struct{}

// Create new `MarkdownConfig`
func NewMarkdownConfig() MarkdownConfig {
	return MarkdownConfig{}
}

// Formatting settings for different markup languages
type FormatConfig struct {
	// Primary markup language for documentation
	MainDocType DocType `yaml:"main_dtype"`

	// Config for markdown
	Markdown MarkdownConfig `yaml:"md"`
}

// Create new `FormatConfig`
func NewFormatConfig(md MarkdownConfig) FormatConfig {
	return FormatConfig{
		MainDocType: MD,
		Markdown:    md,
	}
}

// The main docsync configuration
type Config struct {
	// Primary documentation language
	Base string `yaml:"base"`

	// List of all documentations
	Langs []string `yaml:"langs"`

	// Programming languages, used in the project
	PLangs []string `yaml:"plangs"`

	// The documentation formatting rules
	Format FormatConfig `yaml:"format"`
}

// Create new docsync `Config` and write it to the file
func NewConfig(dir_path, base string, plangs []string, format FormatConfig, create_template bool) error {
	// validate primary documentation language code, according to the ISO639-1
	if !iso6391.ValidCode(base) {
		return nil
	}

	// create docsync config
	config := Config{
		Base:   base,
		Langs:  []string{base},
		PLangs: plangs,
		Format: format,
	}

	// marshall config to yaml
	data, err_yaml := yaml.Marshal(config)
	if err_yaml != nil {
		return err_yaml
	}

	// write config
	if len(dir_path) != 0 && dir_path != "." {
		err_dir := os.Mkdir(path.Join(dir_path), os.ModePerm)
		if err_dir != nil {
			return err_dir
		}
	}

	f, err_file := os.Create(path.Join(dir_path, "docsync.yaml"))
	if err_file != nil {
		return err_file
	}
	defer f.Close()

	_, err_file = f.Write(data)
	if err_file != nil {
		return err_file
	}

	// make template for the base language
	if create_template {
		return CreateEmptyTemplate(dir_path, base, plangs, config.Format.MainDocType)
	}

	return nil
}

// Add new documentation language
func AddLanguage(lang string, create_template_from_base, create_empty_template bool) error {
	// validate language code
	if !iso6391.ValidCode(lang) {
		return nil
	}

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

	config.Langs = append(config.Langs, lang)

	// marshall config to yaml
	data, err_yaml = yaml.Marshal(config)
	if err_yaml != nil {
		return err_yaml
	}

	// write config
	err_file = ioutil.WriteFile("docsync.yaml", data, 0)
	if err_file != nil {
		return err_file
	}

	// make template for the new language
	if create_template_from_base {
		return CreateTemplateFromBase(config.Base, lang, config.Format.MainDocType)
	} else if create_empty_template {
		return CreateEmptyTemplate("", lang, config.PLangs, config.Format.MainDocType)
	}

	return nil
}
