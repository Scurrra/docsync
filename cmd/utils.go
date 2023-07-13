package cmd

import (
	"errors"
	"fmt"
	"os"

	. "github.com/eminarican/safetypes"
	"github.com/manifoldco/promptui"
)

type promptContent struct {
	errorMsg string
	label    string
	dflt     Option[string]
}

func promptGetInput(pc promptContent) string {
	validate := func(input string) error {
		if len(input) <= 0 {
			return errors.New(pc.errorMsg)
		}
		return nil
	}

	templates := &promptui.PromptTemplates{
		Prompt:  "{{ . }} ",
		Valid:   "{{ . | green }} ",
		Invalid: "{{ . | red }} ",
		Success: "{{ . | bold }} ",
	}

	var prompt promptui.Prompt
	if pc.dflt.IsSome() {
		prompt = promptui.Prompt{
			Label:     pc.label,
			Templates: templates,
			Validate:  validate,
			Default:   *pc.dflt.Value,
		}
	} else {
		prompt = promptui.Prompt{
			Label:     pc.label,
			Templates: templates,
			Validate:  validate,
		}
	}

	result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}

	return result
}

func promptGetSelect(pc promptContent, items []string) string {
	index := -1
	var result string
	var err error

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }} ",
		Active:   "{{ . | green }} ",
		Inactive: "{{ . | blue }} ",
		Selected: fmt.Sprintf("%s: {{ . | bold }} ", pc.label),
	}

	for index < 0 {
		prompt := promptui.Select{
			Label:     pc.label,
			Items:     items,
			Templates: templates,
		}

		index, result, err = prompt.Run()
	}

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}

	return result
}
