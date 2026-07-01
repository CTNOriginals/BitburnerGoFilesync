package config

import "testing"

func TestFilePatternsValidateValues_Valid(t *testing.T) {
	var saved = *Values
	t.Cleanup(func() { *Values = saved })

	Values.FilePatterns = SConfigFilePatterns{
		Include: []string{"**/*.js", "*.ts"},
		Exclude: []string{"**/*.d.ts"},
	}

	var err = Values.FilePatterns.ValidateValues()
	if err != nil {
		t.Errorf("ValidateValues() = %v, want nil", err)
	}
}

func TestFilePatternsValidateValues_Invalid(t *testing.T) {
	var saved = *Values
	t.Cleanup(func() { *Values = saved })

	Values.FilePatterns = SConfigFilePatterns{
		Include: []string{"[invalid"},
	}

	var err = Values.FilePatterns.ValidateValues()
	if err == nil {
		t.Fatal("ValidateValues() = nil, want error")
	}
}

func TestFilePatternsValidateValues_Mixed(t *testing.T) {
	var saved = *Values
	t.Cleanup(func() { *Values = saved })

	Values.FilePatterns = SConfigFilePatterns{
		Include: []string{"*.ts", "[invalid"},
		Exclude: []string{"bad[pattern"},
	}

	var err = Values.FilePatterns.ValidateValues()
	if err == nil {
		t.Fatal("ValidateValues() = nil, want error")
	}
}
