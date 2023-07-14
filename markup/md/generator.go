package md

import (
	"fmt"
	"strings"

	"github.com/Scurrra/docsync/markup"
	. "github.com/eminarican/safetypes"
)

// Function that generates empty markdown documentation file template.
func GenerateEmptyDocumentTemplate(plangs []string) markup.Document {
	doc := markup.GenerateEmptyDocumentTemplateIndependent(plangs)

	header := `
# Documentation

This is markdown template. 
	
## Header

This the header part of the document, where you can write the description of the file content.

The next part of the file is the group of documentation blocks.
`

	footer := `
## Footer

This is the footer part of the document.
`

	doc.Content = header + doc.Content + footer

	return doc
}

// Function that generates marfdown documentation from base.
func GenerateDocumentTemplateBase(doc markup.Document) markup.Document {
	// comment Header
	all_blocks := markup.DocumentationBlockKey.FindAllIndex([]byte(doc.Content), -1)
	if len(all_blocks) == 0 {
		return doc
	}
	header := fmt.Sprintf("<!--\n%s\n-->\n\n",
		strings.Trim(doc.Content[:all_blocks[0][0]], " \t\n"),
	)

	// comment footer
	footer := fmt.Sprintf("\n\n<!--\n%s\n-->\n\n",
		strings.Trim(doc.Content[all_blocks[len(all_blocks)-1][1]:], " \t\n"),
	)

	// new content
	doc.Content = header + doc.Content[all_blocks[0][0]:all_blocks[len(all_blocks)-1][1]] + footer

	// each block
	for key, block := range doc.Blocks {
		// snippet description
		if block.Description.IsSome() {
			block.Description = Some(fmt.Sprintf("<!--%s-->", *block.Description.Value))
		}

		// arguments descriptions
		for argi, desc := range block.Arguments.Args {
			block.Arguments.Args[argi] = fmt.Sprintf("<!--%s-->", desc)
		}

		// examples descriptions
		for exi, example := range block.Examples {
			if example.Description.IsSome() {
				block.Examples[exi].Description = Some(fmt.Sprintf("<!--%s-->", *example.Description.Value))
			}
		}

		// comments descriptions
		for comi, comment := range block.Comments {
			if comment.Description.IsSome() {
				block.Comments[comi].Description = Some(fmt.Sprintf("<!--%s-->", *comment.Description.Value))
			}
		}

		doc.Blocks[key] = block
	}

	return doc
}
