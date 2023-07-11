package md

import "github.com/Scurrra/docsync/markup"

// Function that generates empty markdown documentation file template.
func GenerateEmptyDocumentTemplateMarkdown(plangs []string) markup.Document {
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
