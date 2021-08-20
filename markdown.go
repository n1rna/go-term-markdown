package markdown

import (
	"bytes"
	"fmt"

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
	var headerBuf bytes.Buffer
	var footerBuf bytes.Buffer

	// renderer.RenderHeader(&headerBuf, doc)
	renderer.CountHeadings(doc)
	ast.WalkFunc(doc, func(node ast.Node, entering bool) ast.WalkStatus {
		return renderer.RenderNode(&headerBuf, node, entering)
	})
	buf := renderer.GetColumnizedBuffer()
	renderer.RenderFooter(&footerBuf, doc)
	fmt.Println(buf.String())
	return headerBuf.Bytes()
}

func Render(source string, lineWidth int, leftPad int, opts ...Options) []byte {
	p := parser.NewWithExtensions(Extensions())
	nodes := md.Parse([]byte(source), p)
	renderer := NewRenderer(lineWidth, leftPad, 2, opts...)

	return FinalRender(nodes, renderer)
}
