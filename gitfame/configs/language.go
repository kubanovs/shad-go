package language

import (
	_ "embed"
	"encoding/json"
	"strings"
)

//go:embed language_extensions.json
var languageConfigBytes []byte

type LanguageConfig struct {
	languages []LanguageObject
}

type LanguageObject struct {
	Name         string   `json:"name"`
	LanguageType string   `json:"type"`
	Extensions   []string `json:"extensions"`
}

func New() *LanguageConfig {
	languages, err := loadLanguages()
	if err != nil {
		panic("error while read language_extensions.json")
	}
	return &LanguageConfig{languages: languages}
}

func loadLanguages() ([]LanguageObject, error) {
	var langs []LanguageObject
	if err := json.Unmarshal(languageConfigBytes, &langs); err != nil {
		return nil, err
	}
	return langs, nil
}

func (lc *LanguageConfig) GetLanguagesExtensions(languages []string) []string {
	var extensions []string

	for _, lang := range languages {
		for _, configLang := range lc.languages {
			if strings.EqualFold(configLang.Name, lang) {
				extensions = append(extensions, configLang.Extensions...)
			}
		}
	}

	return extensions
}
