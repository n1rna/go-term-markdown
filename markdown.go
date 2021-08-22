package markdown

import (
	"bytes"
	"regexp"

	md "github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/parser"
)

// Extensions returns the bitmask of extensions supported by this renderer.
// The output of this function can be used to instantiate a new markdown
// parser using the `NewWithExtensions` function.
func Extensions() parser.Extensions {
	extensions := parser.NoIntraEmphasis        // Ignore emphasis markers inside words
	extensions |= parser.Tables                 // Parse tables
	extensions |= parser.FencedCode             // Parse fenced code blocks
	extensions |= parser.Autolink               // Detect embedded URLs that are not explicitly marked
	extensions |= parser.Strikethrough          // Strikethrough text using ~~test~~
	extensions |= parser.SpaceHeadings          // Be strict about prefix heading rules
	extensions |= parser.HeadingIDs             // specify heading IDs  with {#id}
	extensions |= parser.BackslashLineBreak     // Translate trailing backslashes into line breaks
	extensions |= parser.DefinitionLists        // Parse definition lists
	extensions |= parser.LaxHTMLBlocks          // more in HTMLBlock, less in HTMLSpan
	extensions |= parser.NoEmptyLineBeforeBlock // no need for new line before a list

	return extensions
}

type ColumnRenderer interface {
	md.Renderer
	CountHeadings(ast ast.Node)
	GetColumnizedBuffer() bytes.Buffer
}

func FinalRender(doc ast.Node, renderer ColumnRenderer) []byte {
	var content bytes.Buffer

	renderer.CountHeadings(doc)
	ast.WalkFunc(doc, func(node ast.Node, entering bool) ast.WalkStatus {
		return renderer.RenderNode(&content, node, entering)
	})

	buf := renderer.GetColumnizedBuffer()
	return buf.Bytes()
}

func CleanupHeaders(doc string) string {
	re := regexp.MustCompile(`-{3}(.*\n)*-{3}`)
	return re.ReplaceAllString(doc, "")
}

func Render(source string, lineWidth int, columns int, leftPad int, opts ...Options) []byte {
	source = CleanupHeaders(source)
	p := parser.NewWithExtensions(Extensions())
	nodes := md.Parse([]byte(source), p)
	renderer := NewRenderer(lineWidth, leftPad, columns, opts...)

	return FinalRender(nodes, renderer)
}
