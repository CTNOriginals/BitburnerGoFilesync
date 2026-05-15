package utils

import "strings"

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
