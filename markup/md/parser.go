package md

import (
	"crypto/sha256"
	"fmt"
	"regexp"
	"strings"

	"github.com/Scurrra/docsync/markup"
	. "github.com/eminarican/safetypes"
)

// Function that parses a single `DocumentationBlock`
func ParseDocumentationBlock(data string) markup.DocumentationBlock {
	var (
		hashKey     string
		status      markup.FeatureStatus
		codeblock   markup.CodeBlock
		description Option[string]
		arguments   markup.Arguments
		examples    []markup.Example
		comments    []markup.Comment
	)

	data = strings.ReplaceAll(data, "\n> ", "\n")
	data = strings.ReplaceAll(data, "\n>", "\n")
	var data_start int // to reduce number of allocations

	// hashkey
	hashKey_start := strings.Index(data, "<!--") + 4
	hashKey_end := strings.Index(data, "-->")
	// useless
	hashKey = strings.Trim(data[hashKey_start:hashKey_end], " \t\n")
	data_start = strings.Index(data, "\n") + 1

	// status
	status_start := data_start + strings.Index(data[data_start:], "<!--") + 4
	status_end := data_start + strings.Index(data[data_start:], "-->")
	status = markup.FeatureStatus(strings.Trim(data[status_start:status_end], " "))
	data_start = status_end + 3

	// codeblock
	lang_start := data_start + strings.Index(data[data_start:], "```") + 3
	lang_end := lang_start + strings.Index(data[lang_start:], "\n")
	lang := data[lang_start:lang_end]
	code_start := lang_end + 1
	code_end := lang_end + strings.Index(data[lang_end:], "\n```\n")
	code := data[code_start:code_end]
	codeblock = markup.CodeBlock{
		Lang:    lang,
		Snippet: code,
	}
	// anyway
	// if len(hashKey) == 0 {
	hashKey = fmt.Sprintf("%x", sha256.Sum256([]byte(code)))
	// }
	data = strings.Trim(data[(code_end+4):], " \n")

	// description
	// but firstly check if there are arguments, examples and other comments
	arguments_start := strings.Index(data, "***Arguments:***\n")
	examples_start := strings.Index(data, "***Examples:***\n")
	r, _ := regexp.Compile(`\*\*\*\S+\*\*\*`)
	var (
		comments_first []int
		comments_start int = -1
	)
	description_start := 0
	description_end := len(data)
	if examples_start != -1 {
		description_end = examples_start
	}
	if arguments_start != -1 {
		description_end = arguments_start
	}
	if description_end == len(data) {
		comments_first = r.FindIndex([]byte(data))
		if len(comments_first) == 2 {
			description_end = comments_first[0]
			comments_start = comments_first[0]
		}
	}
	desc := strings.Trim(data[description_start:description_end], " \n")
	if len(desc) != 0 {
		description = Some(desc)
	} else {
		description = None[string]()
	}

	// arguments
	if arguments_start == -1 {
		arguments = markup.Arguments{Args: nil}
	} else {
		arguments = markup.Arguments{Args: make(map[string]string)}
		data_start = arguments_start + 17 // len("***Arguments:***\n")
		for {
			if data_start >= len(data) {
				break
			}
			// parse `arg`
			arg_start := data_start + strings.Index(data[data_start:], " - *`")
			if arg_start != -1 {
				arg_start += 5
			} else {
				break
			}
			arg_end := data_start + strings.Index(data[data_start:], "`* - ")
			if arg_end < arg_start || arg_end == -1 {
				break
			}
			arg := data[arg_start:arg_end]

			// parse `arg` description
			desc_start := arg_end + 5
			desc_end := desc_start + strings.Index(data[desc_start:], "\n")
			desc := data[desc_start:desc_end]

			data_start = desc_end + 1
			arguments.Args[arg] = desc
		}
	}

	// examples
	if examples_start == -1 {
		examples = nil
	} else {
		examples = []markup.Example{}
		data_start = examples_start + 16 // len("***Examples:***\n")

		comments_first = r.FindIndex([]byte(data[data_start:]))
		if len(comments_first) == 2 {
			comments_first[0] += data_start
			comments_first[1] += data_start
			comments_start = comments_first[0]
		}

		for {
			if data_start >= len(data) {
				break
			}
			if len(comments_first) == 2 && data_start >= comments_first[0] {
				break
			}
			// parse example description
			desc_start := data_start
			desc_end := data_start + strings.Index(data[data_start:], "\n```")
			if desc_end == -1 || (comments_start != -1 && desc_end > comments_start) {
				break
			}
			desc := Some(strings.Trim(data[desc_start:desc_end], " \n"))
			if len(*desc.Value) == 0 {
				desc = None[string]()
			}

			// parse snippet lang
			lang_start := desc_end + 4
			lang_end := lang_start + strings.Index(data[lang_start:], "\n")
			lang := data[lang_start:lang_end]

			// parse snippet
			code_start := lang_end + 1
			code_end := code_start + strings.Index(data[code_start:], "\n```\n")
			code := data[code_start:code_end]

			data_start = code_end + 5
			examples = append(examples, markup.Example{
				Description: desc,
				Code: markup.CodeBlock{
					Lang:    lang,
					Snippet: code,
				},
			})
		}
	}

	comments_first = r.FindIndex([]byte(data[data_start:]))
	if len(comments_first) == 2 {
		comments_first[0] += data_start
		comments_first[1] += data_start
		comments_start = comments_first[0]
	}

	// comments
	if comments_start == -1 {
		comments = nil
	} else {
		comments = []markup.Comment{}
		data_start = comments_start

		for {
			if data_start >= len(data) {
				break
			}
			// parse comment name
			comment_start := data_start + strings.Index(data[data_start:], "***")
			if comment_start == -1 {
				break
			}
			comment_end := comment_start + strings.Index(data[comment_start:], ":***\n")
			name := data[(comment_start + 3):comment_end]

			// parse comment description
			desc_start := comment_end + 4
			desc_end := desc_start + strings.Index(data[desc_start:], "\n```")
			if desc_end == -1 {
				break
			}
			desc := Some(strings.Trim(data[desc_start:desc_end], " \n"))
			if len(*desc.Value) == 0 {
				desc = None[string]()
			}

			// parse snippet lang
			lang_start := desc_end + 4
			lang_end := lang_start + strings.Index(data[lang_start:], "\n")
			lang := data[lang_start:lang_end]

			// parse snippet
			code_start := lang_end + 1
			code_end := code_start + strings.Index(data[code_start:], "\n```")
			code := data[code_start:code_end]

			data_start = code_end + 4
			comments = append(comments, markup.Comment{
				Name:        name,
				Description: desc,
				Code: markup.CodeBlock{
					Lang:    lang,
					Snippet: code,
				},
			})
		}
	}

	return markup.DocumentationBlock{
		HashKey:     hashKey,
		Status:      status,
		Code:        codeblock,
		Description: description,
		Arguments:   arguments,
		Examples:    examples,
		Comments:    comments,
	}
}

// Function that parses document into `Document` struct
func ParseDocument(content string) markup.Document {

	blocks := make(map[string]markup.DocumentationBlock)

	doc_start := strings.Index(content, "> <!--docbegin-->\n") // len("> <!--docbegin-->\n") == 18
	doc_end := strings.Index(content, "> <!--docend-->\n")     // len("> <!--docbegin-->\n") == 16

	for {
		if doc_start == -1 {
			break
		}

		block := ParseDocumentationBlock(content[(doc_start + 18):doc_end])
		blocks[string(block.HashKey[:])] = block
		content = fmt.Sprintf("\n%s<[%s]>%s\n\n", content[:doc_start], block.HashKey, content[(doc_end+16):])

		doc_start = strings.Index(content, "> <!--docbegin-->\n") // len("> <!--docbegin-->\n") == 18
		doc_end = strings.Index(content, "> <!--docend-->\n")     // len("> <!--docbegin-->\n") == 16
	}

	return markup.Document{
		Content: strings.Trim(content, " \t\n"),
		Blocks:  blocks,
	}
}
