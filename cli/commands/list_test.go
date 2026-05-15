package commands

import "testing"

func TestPushAddsHelp(t *testing.T) {
	var list TList
	var def = &Definition{
		Name: "testcmd",
	}
	list.Push(def)

	if def.Options.GetDefinitionByName("help") == nil {
		t.Error("Push should auto-attach a 'help' sub-command")
	}
}

func TestParseInputByName(t *testing.T) {
	var list = TList{
		{Name: "alpha"},
		{Name: "beta"},
	}

	var got, err = list.ParseInput("alpha")
	if err != nil {
		t.Fatalf("ParseInput('alpha') unexpected error: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("ParseInput('alpha') returned %d inputs, want 1", len(got))
	}
	if got[0].Def.Name != "alpha" {
		t.Errorf("Def.Name = %q, want 'alpha'", got[0].Def.Name)
	}
}

func TestParseInputUnknown(t *testing.T) {
	var list = TList{{Name: "alpha"}}

	var err error
	_, err = list.ParseInput("unknown")
	if err == nil {
		t.Error("ParseInput('unknown') expected error")
	}
}

func TestParseInputNested(t *testing.T) {
	var list = TList{{
		Name: "config",
		Options: TList{
			{Name: "list"},
			{Name: "set"},
		},
	}}

	var got, err = list.ParseInput("config", "list")
	if err != nil {
		t.Fatalf("ParseInput('config', 'list') unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("got %d inputs, want 2", len(got))
	}
	if got[0].Def.Name != "config" {
		t.Errorf("got[0].Def.Name = %q, want 'config'", got[0].Def.Name)
	}
	if got[1].Def.Name != "list" {
		t.Errorf("got[1].Def.Name = %q, want 'list'", got[1].Def.Name)
	}
}

func TestParseInputByValidator(t *testing.T) {
	var list = TList{
		{Name: "alpha"},
		{
			Name: "numeric",
			Validator: func(input string) bool {
				for _, c := range input {
					if c < '0' || c > '9' {
						return false
					}
				}
				return true
			},
		},
	}

	var got, err = list.ParseInput("42")
	if err != nil {
		t.Fatalf("ParseInput('42') error: %v", err)
	}
	if got[0].Def.Name != "numeric" {
		t.Errorf("matched %q, want 'numeric'", got[0].Def.Name)
	}
}

func TestParseInputByAutocompleteFallback(t *testing.T) {
	var list = TList{
		{Name: "alpha"},
		{
			Name: "file",
			AutoComplete: func(s string) []string {
				return []string{"file.go", "main.go"}
			},
		},
	}

	var got, err = list.ParseInput("file.go")
	if err != nil {
		t.Fatalf("ParseInput('file.go') error: %v", err)
	}
	if got[0].Def.Name != "file" {
		t.Errorf("matched %q, want 'file'", got[0].Def.Name)
	}
}

func TestTryGetValidatedDefinitionNoMatch(t *testing.T) {
	var list = TList{{Name: "alpha"}, {Name: "beta"}}

	var got = list.TryGetValidatedDefinition("unknown")
	if got != nil {
		t.Errorf("got %v, want nil", got)
	}
}

func TestTryGetValidatedDefinitionValidatorMatch(t *testing.T) {
	var list = TList{
		{Name: "alpha"},
		{
			Name: "number",
			Validator: func(input string) bool {
				for _, c := range input {
					if c < '0' || c > '9' {
						return false
					}
				}
				return true
			},
		},
	}

	var got = list.TryGetValidatedDefinition("123")
	if got == nil || got.Name != "number" {
		t.Errorf("got %v, want definition 'number'", got)
	}
}

func TestTryGetValidatedDefinitionExpectValue(t *testing.T) {
	var list = TList{
		{Name: "alpha"},
		{Name: "anything", ExpectValue: true},
	}

	var got = list.TryGetValidatedDefinition("whatever")
	if got == nil || got.Name != "anything" {
		t.Errorf("got %v, want definition 'anything'", got)
	}
}

func TestGetNamesAndGetNamesRecursive(t *testing.T) {
	var list = TList{{
		Name: "parent",
		Options: TList{
			{Name: "child_a"},
			{Name: "child_b"},
		},
	}}

	var names = list.GetNames()
	if len(names) != 1 || names[0] != "parent" {
		t.Errorf("GetNames() = %v, want ['parent']", names)
	}

	var all = list.GetNamesRecursive()
	if len(all) != 3 {
		t.Errorf("GetNamesRecursive() returned %d names, want 3", len(all))
	}
}

func TestGetDefinitionByName(t *testing.T) {
	var list = TList{{Name: "alpha"}, {Name: "beta"}}

	var got = list.GetDefinitionByName("alpha")
	if got == nil || got.Name != "alpha" {
		t.Error("GetDefinitionByName('alpha') should find 'alpha'")
	}

	got = list.GetDefinitionByName("nonexistent")
	if got != nil {
		t.Error("GetDefinitionByName('nonexistent') should return nil")
	}
}
