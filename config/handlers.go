package config

type SConfigHandlersNSDefinitions struct {
	GetOnConnect bool
	Destination  string
}

func (this *SConfigHandlersNSDefinitions) ValidateValues() error {
	return nil
}

type SConfigHandlers struct {
	NetscriptDefinitions *SConfigHandlersNSDefinitions
}

func (this *SConfigHandlers) ValidateValues() error {
	var fields = []IConfigGroup{
		this.NetscriptDefinitions,
	}

	for _, field := range fields {
		var err = field.ValidateValues()
		if err != nil {
			return err
		}
	}

	return nil
}
