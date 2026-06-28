package config

type SConfigLogging struct {
	NoColor   bool
	LogConfig bool
}

func (this SConfigLogging) ValidateValues() error {
	return nil
}
