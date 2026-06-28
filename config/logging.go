package config

type SConfigLogging struct {
	NoColor bool
}

func (this SConfigLogging) ValidateValues() error {
	return nil
}
