package commands

import "testing"

func TestExecuteDeepestNode(t *testing.T) {
	var executed string

	var list = TList{{
		Name: "parent",
		Execution: func(args TInputList, self int) {
			executed = "parent"
		},
		Options: TList{{
			Name: "child",
			Execution: func(args TInputList, self int) {
				executed = "child"
			},
		}},
	}}

	var inputs, err = list.ParseInput("parent", "child")
	if err != nil {
		t.Fatalf("ParseInput error: %v", err)
	}

	inputs.Execute()
	if executed != "child" {
		t.Errorf("Execute() ran %q, want 'child'", executed)
	}
}

func TestExecuteNoExecution(t *testing.T) {
	var list = TList{{
		Name:    "parent",
		Options: TList{{Name: "child"}},
	}}

	var inputs, err = list.ParseInput("parent", "child")
	if err != nil {
		t.Fatalf("ParseInput error: %v", err)
	}

	inputs.Execute() // should not panic
}
