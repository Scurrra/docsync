package dsconfig

import (
	"io/ioutil"

	iso6391 "github.com/emvi/iso-639-1"
	"gopkg.in/yaml.v3"
)

type DocType string

const (
	MD DocType = "md"
)

// Configuration settings for markdown documentation
type MarkdownConfig struct {
	Separator string `yaml:"separator"`
}

// Create new `MarkdownConfig`
func newMarkdownConfig(separator string) MarkdownConfig {
	return MarkdownConfig{Separator: separator}
}

// Formatting settings for different markup languages
type FormatConfig struct {
	// Primary markup language for documentation
	MainDocType DocType `yaml:"main_dtype"`

	// Config for markdown
	Markdown MarkdownConfig `yaml:"md"`
}

// Create new `FormatConfig`
func newFormatConfig(md MarkdownConfig) FormatConfig {
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
func newConfig(base string, plangs []string, format *FormatConfig, create_template bool) error {
	// validate primary documentation language code, according to the ISO639-1
	if !iso6391.ValidCode(base) {
		return nil
	}

	// create docsync config
	config := Config{
		Base:   base,
		Langs:  []string{base},
		PLangs: plangs,
		Format: *format,
	}

	// marshall config to yaml
	data, err_yaml := yaml.Marshal(config)
	if err_yaml != nil {
		return err_yaml
	}

	// write config
	err_file := ioutil.WriteFile("docsync.yaml", data, 0)
	if err_file != nil {
		return err_file
	}

	// make template for the base language
	if create_template {
		return createEmptyTemplate(base, plangs, config.Format.MainDocType)
	}

	return nil
}

// Add new documentation language
func addLanguage(lang string, create_template bool) error {
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

	// make template for the base language
	if create_template {
		return createTemplateFromBase(lang)
	}

	return nil
}
