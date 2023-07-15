package md

import (
	"fmt"
	"strings"

	"github.com/Scurrra/docsync/markup"
	. "github.com/eminarican/safetypes"
	"golang.org/x/exp/slices"
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
	all_blocks := markup.DocumentationBlockKey.FindAllIndex([]byte(doc.Content), -1)
	if len(all_blocks) == 0 {
		return doc
	}

	// // comment Header
	// header := fmt.Sprintf("<!--\n%s\n-->\n\n",
	// 	strings.Trim(doc.Content[:all_blocks[0][0]], " \t\n"),
	// )

	// // comment footer
	// footer := fmt.Sprintf("\n\n<!--\n%s\n-->\n\n",
	// 	strings.Trim(doc.Content[all_blocks[len(all_blocks)-1][1]:], " \t\n"),
	// )

	// // new content
	// doc.Content = header + doc.Content[all_blocks[0][0]:all_blocks[len(all_blocks)-1][1]] + footer

	contents := make([]string, 2*len(all_blocks)+1)
	// comment Header
	contents[0] = fmt.Sprintf("<!--\n%s\n-->\n\n\n",
		strings.Trim(doc.Content[:all_blocks[0][0]], " \t\n"),
	)

	contents_i := 1
	for i := 0; i < len(all_blocks)-1; i++ {
		buf := strings.Trim(doc.Content[all_blocks[i][1]:all_blocks[i+1][0]], " \t\n")
		if len(buf) != 0 {
			buf = fmt.Sprintf("\n<!--\n%s\n-->\n\n\n", buf)
		} else {
			buf = doc.Content[all_blocks[i][1]:all_blocks[i+1][0]]
		}
		contents[contents_i] = doc.Content[all_blocks[i][0]:all_blocks[i][1]]
		contents_i += 1
		contents[contents_i] = buf
		contents_i += 1
	}
	contents[contents_i] = doc.Content[all_blocks[len(all_blocks)-1][0]:all_blocks[len(all_blocks)-1][1]]

	// comment footer
	contents[len(contents)-1] = fmt.Sprintf("\n\n<!--\n%s\n-->\n\n\n",
		strings.Trim(doc.Content[all_blocks[len(all_blocks)-1][1]:], " \t\n"),
	)

	// new content
	doc.Content = strings.Join(contents, "")

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

// Function that merges two docs into the new one
func GenerateMergedDocument(doc_base, doc_lang markup.Document) markup.Document {
	// works only with blocks
	// if block already exists, it's status sets to Active (except if it's deprecated)
	// if it's a new block, it's copied to the end of the file
	// if block is deleted, it's status set to Deprecated

	var (
		lang_blocks_seen = []string{}
		new_blocks       = false
	)

	for key_base, block_base := range doc_base.Blocks {
		lang_blocks_seen = append(lang_blocks_seen, key_base)
		block_lang, ok := doc_lang.Blocks[key_base]
		if ok {
			if block_base.Status != markup.Deprecated {
				block_lang.Status = markup.Active
			} // else {
			// 	block_lang.Status = markup.Deprecated
			// }

			doc_lang.Blocks[key_base] = block_lang
		} else {
			if !new_blocks {
				new_blocks = true
				doc_lang.Content += "\n\n# New blocks\n"
			}
			doc_lang.Blocks[key_base] = block_base
			doc_lang.Content = fmt.Sprintf("%s\n<[%s]>\n", doc_lang.Content, block_base.HashKey)
		}
	}
	for key_lang, block_lang := range doc_lang.Blocks {
		if !slices.Contains(lang_blocks_seen, key_lang) {
			block_lang.Status = markup.Deprecated
			doc_lang.Blocks[key_lang] = block_lang
		}
	}

	return doc_lang
}
