// Structure of the single documentation file.

package markup

import (
	"fmt"
	"regexp"
	"strings"

	. "github.com/eminarican/safetypes"
)

// Feature status
type FeatureStatus string

const (
	// New feature
	New FeatureStatus = "New"

	// Feature is active
	Active FeatureStatus = "Active"

	// Feature is Deprecated
	Deprecated FeatureStatus = "Deprecated"
)

// Code snippet (function definition, struct, class, module, i.e.)
// to be documented.
type CodeBlock struct {
	// Programming language
	Lang string

	// Code
	Snippet string
}

// Arguments for documentating code.
type Arguments struct {
	// Map where keys are arguments and values are definitions
	Args map[string]string
}

// Usage example.
type Example struct {
	// Optional descripton of the example
	Description Option[string]

	// Code of the example
	Code CodeBlock
}

// Any other documentation paragraphs.
type Comment struct {
	// Name of the comment paragraph, like "Example" in `Example`
	Name string

	// Optional description of the <`Name`> comment
	Description Option[string]

	// Code of the <`Name`> comment
	Code CodeBlock
}

// Single documentation block.
//
// # Fields are in the following order:
//
// 1. `CodeBlock`
//
// 2. Description for the `CodeBlock`
//
// 3. `Arguments` for the `CodeBlock`
//
// 4. List of `Example`s
//
// 5. Other `Comment`s
type DocumentationBlock struct {
	HashKey     string
	Status      FeatureStatus
	Code        CodeBlock
	Description Option[string]
	Arguments   Arguments
	Examples    []Example
	Comments    []Comment
}

var DocumentationBlockKey = regexp.MustCompile(`<\[\w+\]>`) // regex is valid, so the error is ignored

// Single file structure.
type Document struct {
	// File content
	Content string

	// Map of `DocumentationBlock`s
	Blocks map[string]DocumentationBlock
}

// Function that generates empty documentation file template.
func GenerateEmptyDocumentTemplateIndependent(plangs []string) Document {
	args := Arguments{
		map[string]string{
			"arg1": "Description for the arg",
			"arg2": "Description for the arg",
			"arg3": "Description for the arg",
		},
	}

	// generate `DocumentationBlock` per each of `plangs`
	blocks := make(map[string]DocumentationBlock)
	for _, plang := range plangs {
		code := CodeBlock{plang, ""}
		blocks[plang] = DocumentationBlock{
			"", //fmt.Sprintf("%x", sha256.Sum256([]byte(code.Snippet))),
			New,
			code,
			Some("Description for the code snippet"),
			args,
			[]Example{
				{
					Some("Description for the first example"),
					code,
				},
				{
					Some("Description for the second example"),
					code,
				},
			},
			[]Comment{
				{
					"CommentName1",
					Some("Description for the <#CommentName1>"),
					code,
				},
				{
					"CommentName2",
					Some("Description for the <#CommentName2>"),
					code,
				},
			},
		}
	}

	// create links in document content
	for i := 0; i < len(plangs); i++ {
		plangs[i] = fmt.Sprintf("<[%s]>", plangs[i])
	}
	content := strings.Join(plangs, "\n\n")

	return Document{
		content,
		blocks,
	}
}
