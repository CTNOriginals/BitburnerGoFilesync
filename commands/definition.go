package commands

type Definition struct {
	// The strings that will activate this command
	Triggers    []string
	Description []string

	Execution func()
}
