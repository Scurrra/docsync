package md

import (
	"bytes"
	"fmt"
	"strings"

	. "github.com/Scurrra/docsync/markup"

	. "github.com/eminarican/safetypes"
)

// Function that renders single `CodeBlock` to markdown
func RenderCodeBlock(b *bytes.Buffer, code_block CodeBlock) {
	b.WriteString(fmt.Sprintf("> ```%s\n> ", code_block.Lang))
	b.WriteString(strings.ReplaceAll(code_block.Snippet, "\n", "\n> "))
	b.WriteString("\n> ```\n>\n")
}

// Function that renders single description to markdown
func RenderDescription(b *bytes.Buffer, description Option[string]) {
	if description.IsSome() {
		b.WriteString("> ")
		b.WriteString(strings.ReplaceAll(*description.Value, "\n", "\n> "))
		b.WriteString(">\n")
	}
}

// Function that renders `Argumanets` to markdown
func RenderArguments(b *bytes.Buffer, args Arguments) {
	if len(args.Args) != 0 {
		b.WriteString("> ***Arguments:***\n")
		for arg, desc := range args.Args {
			b.WriteString(
				fmt.Sprintf(
					">  - *`%s`* - %s\n",
					arg,
					strings.ReplaceAll(desc, "\n", "\n> "),
				),
			)
		}
		b.WriteString(">\n")
	}
}

// Function that renders `Exaple`s block
func RenderExamples(b *bytes.Buffer, examples []Example) {
	if examples != nil || len(examples) != 0 {
		b.WriteString("> ***Examples:***\n>\n")
		for _, example := range examples {
			RenderDescription(b, example.Description)
			RenderCodeBlock(b, example.Code)
		}
	}
}

// Function that renders `Comment`s block
func RenderComments(b *bytes.Buffer, comments []Comment) {
	if comments != nil || len(comments) != 0 {
		for _, comment := range comments {
			b.WriteString(fmt.Sprintf("> ***%s:***\n>\n", comment.Name))
			RenderDescription(b, comment.Description)
			RenderCodeBlock(b, comment.Code)
		}
	}
}

// Function that renders single documentation block
func RenderDocumentationBlock(doc_block DocumentationBlock) string {
	var b bytes.Buffer

	b.Write([]byte(fmt.Sprintf("> <!--%x-->\n", doc_block.HashKey)))
	b.Write([]byte(fmt.Sprintf("> <!--%v-->\n", doc_block.Status)))

	RenderCodeBlock(&b, doc_block.Code)
	RenderDescription(&b, doc_block.Description)
	RenderArguments(&b, doc_block.Arguments)
	RenderExamples(&b, doc_block.Examples)
	RenderComments(&b, doc_block.Comments)

	return b.String()
}

func RenderDocument(doc Document) string {
	content := doc.Content
	for doc_key, doc_block := range doc.Blocks {
		content = strings.Replace(
			content,
			"<["+doc_key+"]>",
			RenderDocumentationBlock(doc_block),
			1,
		)
	}

	return content
}
