package cli

import "testing"

func TestGetInputSegments(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []string
		wantErr bool
	}{
		{name: "empty string", input: "", want: []string{""}},
		{name: "basic space split", input: "config set Port 1234", want: []string{"config", "set", "Port", "1234"}},
		{name: "single quoted segment", input: "config set 'foo bar'", want: []string{"config", "set", "'foo bar'"}},
		{name: "double quoted segment", input: `config set "foo bar"`, want: []string{"config", "set", `"foo bar"`}},
		{name: "backtick quoted segment", input: "config set `foo bar`", want: []string{"config", "set", "`foo bar`"}},
		{name: "mixed nested quotes", input: `"foo 'baz' bar" goo 'drap bah'`, want: []string{`"foo 'baz' bar"`, "goo", "'drap bah'"}},
		{name: "escaped quote inside", input: "goo foo='bar ins\\'t baz' but is", want: []string{"goo", "foo='bar ins\\'t baz'", "but", "is"}},
		{name: "unmatched double quote", input: `"foo 'baz`, wantErr: true},
		{name: "unmatched single quote trailing", input: `something 'like a dog' or cat'`, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got, err = GetInputSegments(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetInputSegments() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if len(got) != len(tt.want) {
				t.Errorf("GetInputSegments() = %v (len %d), want %v (len %d)", got, len(got), tt.want, len(tt.want))
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
