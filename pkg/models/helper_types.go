package models

type Changelog struct {
	Entries []ChangelogEntry `yaml:"entries"`
}

type ChangelogEntry struct {
	Version   string `yaml:"version"`
	EntryType string `yaml:"type"`
	EntryDate string `yaml:"date"`
	Message   string `yaml:"message"`
}
