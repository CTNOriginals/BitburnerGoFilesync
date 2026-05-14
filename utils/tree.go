package utils

import (
	"fmt"
	"slices"
	"strings"

	ctnstring "github.com/CTNOriginals/CTNGoUtils/v2/string"
)

type TTreeSymbol rune

type TreeSymbols struct {
	Line   TTreeSymbol
	Split  TTreeSymbol
	Corner TTreeSymbol
	Dash   TTreeSymbol
}

const (
	TreeLine       TTreeSymbol = '┃'
	TreeSplit      TTreeSymbol = '┣'
	TreeCorner     TTreeSymbol = '┗'
	TreeDash       TTreeSymbol = '╸'
	TreeThinLine   TTreeSymbol = '│'
	TreeThinSplit  TTreeSymbol = '├'
	TreeThinCorner TTreeSymbol = '└'
	TreeThinDash   TTreeSymbol = '╴'
)

func TreeSymbolList() []TTreeSymbol {
	return []TTreeSymbol{
		TreeLine, TreeSplit, TreeCorner, TreeDash,
		TreeThinLine, TreeThinSplit, TreeThinCorner, TreeThinDash,
	}
}

func GetTreeSymbols(thin bool) TreeSymbols {
	if thin {
		return TreeSymbols{
			Line:   TreeThinLine,
			Split:  TreeThinSplit,
			Corner: TreeThinCorner,
			Dash:   TreeThinDash,
		}
	}

	return TreeSymbols{
		Line:   TreeLine,
		Split:  TreeSplit,
		Corner: TreeCorner,
		Dash:   TreeDash,
	}

}

type TreeFormatOptions struct {
	// The amount of indent to account for
	IndentCount int

	// If the next line's character at IndentCount is not a Whitespace character
	// the line will be treated as a split
	Whitespace []rune
	// will cause the function to place a
	// corner when its reached and stop formatting any further depth
	MaxDepth int
	Thin     bool
	// TODO: maybe also add a 'Dashed' style option
}

// Formats the string into a folding treelike string
// with special symbols indicating depth, splits and connections.
//
// The resulting string is returned
// and the passed in string builder is also appended to in place
func FormatStringAsTree(
	builder *strings.Builder,
	options TreeFormatOptions,
	lines ...string,
) string {
	if builder == nil {
		builder = &strings.Builder{}
	}

	var symbols = GetTreeSymbols(options.Thin)

	for i := 1; i < len(lines); i++ {
		builder.WriteString("\n")

		if options.MaxDepth <= 0 {
			builder.WriteString(lines[i])
			continue
		}
		var line = []rune(lines[i])
		var char = symbols.Line

		if options.IndentCount < len(line)-1 && !strings.ContainsRune(
			string(TreeSymbolList())+string(options.Whitespace),
			line[options.IndentCount],
		) {
			if options.MaxDepth > 1 {
				char = symbols.Split
			} else {
				char = symbols.Corner
			}

			if options.IndentCount > 1 {
				line[1] = rune(symbols.Dash)
			}
			options.MaxDepth -= 1
		}

		line[0] = rune(char)
		builder.WriteString(string(line))
	}

	return builder.String()
}

type TreeObjectFormatOptions struct {
	validated bool

	// The amount of space to add infront of the content on each line per depth.
	// Default: 2
	IndentCount int

	// Default " "
	IndentString string

	// If the content is already indented, the indent char used can be added in here.
	// if a continues vertical line of ReplaceRunes is found, the tree will connect it.
	//
	// Default: []rune(" ")
	ReplaceRunes []rune
}

var treeObjectFormatOptionDefaults = TreeObjectFormatOptions{
	IndentCount:  2,
	IndentString: " ",
	ReplaceRunes: []rune(" "),
}

func (this *TreeObjectFormatOptions) Validate() {
	if this.validated {
		return
	}

	// Object to compare unset values to
	var compare = TreeObjectFormatOptions{}

	if this.IndentCount == compare.IndentCount {
		this.IndentCount = treeObjectFormatOptionDefaults.IndentCount
	}
	if this.IndentString == compare.IndentString {
		this.IndentString = treeObjectFormatOptionDefaults.IndentString
	}
	if string(this.ReplaceRunes) == string(compare.ReplaceRunes) {
		this.ReplaceRunes = treeObjectFormatOptionDefaults.ReplaceRunes
	}

	this.validated = true
}

type TreeObjectList []*TreeObject

func (this *TreeObjectList) Clean() {
	var list = *this

	for i := 0; i < len(list); i++ {
		var obj = list[i]
		if obj == nil {
			goto cut
		}

		obj.Children.Clean()

		if (obj.Content != "" && obj.Content != nil) || len(obj.Children) > 0 {
			continue
		}
	cut:
		list = append(list[0:i], list[i+1:]...)
	}

	if len(list) == 0 {
		*this = nil
		return
	}

	// fmt.Printf("%d > %d\n", len(*this), len(list))
	*this = list
}

func (this TreeObjectList) String() string {
	var builder strings.Builder

	for _, obj := range this {
		if builder.Len() > 0 {
			builder.WriteRune('\n')
		}

		builder.WriteString(obj.String())
	}

	return builder.String()
}

func (this TreeObjectList) Generate(options *TreeObjectFormatOptions) []string {
	var sections = make([]string, len(this))

	for i, obj := range this {
		sections[i] = obj.Generate(options)
	}

	return sections
}

type TreeObject struct {
	// Value must be either of type string, or have a String() method.
	Content  any
	Children TreeObjectList
}

func (this TreeObject) ContentString() string {

	switch val := this.Content.(type) {
	case string:
		return val
	case fmt.Stringer:
		return val.String()
	case nil:
		return ""
	}

	return fmt.Sprintf("%v", this.Content)
}

func (this TreeObject) Generate(options *TreeObjectFormatOptions) string {
	// var builder strings.Builder
	var content string = this.ContentString()

	this.Children.Clean()
	if len(this.Children) == 0 {
		return content
	}

	var lines = strings.Split(content, "\n")
	var symbols = GetTreeSymbols(false)

	if options == nil {
		options = &treeObjectFormatOptionDefaults
	} else {
		options.Validate()
	}

	for i := 0; i < len(lines); i++ {
		var line = []rune(lines[i])

		if slices.Contains(options.ReplaceRunes, line[0]) {
			line[0] = rune(symbols.Line)
		}

		lines[i] = string(line)
	}

	var children = this.Children.Generate(options)

	for i, child := range children {
		child = ctnstring.Indent(child, options.IndentCount, options.IndentString)
		var clines = strings.Split(child, "\n")

		for j := 0; j < len(clines); j++ {
			var line = []rune(clines[j])

			if i >= len(children)-1 {
				line[0] = rune(symbols.Corner)
				clines[0] = string(line)
				break
			}

			if j == 0 {
				if len(clines) > 1 || j <= len(clines)-1 {
					line[0] = rune(symbols.Split)
				} else {
					line[0] = rune(symbols.Corner)
				}
			} else {
				line[0] = rune(symbols.Line)
			}

			clines[j] = string(line)
		}

		lines = append(lines, clines...)
	}

	// content = fmt.Sprintf("%s\n%s", content, children)

	return strings.Join(lines, "\n")
}

func (this TreeObject) String() string {
	var str = fmt.Sprintf("Content: %v\nChildren: {\n%s\n}", this.Content, ctnstring.Indent(this.Children.String(), 2, " "))
	if len(this.Children) == 0 || this.Children == nil {
		str = fmt.Sprintf("Content: %v", this.Content)
	}
	return str
}
