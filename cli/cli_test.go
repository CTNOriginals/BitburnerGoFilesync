package cli

import (
	"testing"
)

func TestGetInputSegments(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{name: "empty string", input: "", want: []string{""}},
		{name: "no space", input: "singleword", want: []string{"singleword"}},
		{name: "no space quoted", input: "'singleword'", want: []string{"'singleword'"}},
		{name: "basic space split", input: "config set Port 1234", want: []string{"config", "set", "Port", "1234"}},
		{name: "single quoted segment", input: "config set 'foo bar'", want: []string{"config", "set", "'foo bar'"}},
		{name: "double quoted segment", input: `config set "foo bar"`, want: []string{"config", "set", `"foo bar"`}},
		{name: "backtick quoted segment", input: "config set `foo bar`", want: []string{"config", "set", "`foo bar`"}},
		{name: "mixed nested quotes", input: `"foo 'baz' bar" goo 'drap bah'`, want: []string{`"foo 'baz' bar"`, "goo", "'drap bah'"}},
		{name: "escaped quote inside", input: "goo foo='bar ins\\'t baz' but is", want: []string{"goo", "foo='bar ins\\'t baz'", "but", "is"}},
		{name: "unmatched double quote", input: `"foo 'baz`, want: []string{`"foo`, `'baz`}},
		{name: "unmatched single quote trailing", input: `something 'like a dog' or cat'`, want: []string{`something`, `'like a dog'`, `or`, `cat'`}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got = GetInputSegments(tt.input)
			if len(got) != len(tt.want) {
				t.Errorf("GetInputSegments()\n got [%s]\t(len %d)\nwant [%s]\t(len %d)", wrapSegments(got), len(got), wrapSegments(tt.want), len(tt.want))
				return
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("GetInputSegments()[%d] = %q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}
