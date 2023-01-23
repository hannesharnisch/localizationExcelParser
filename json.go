package localizationExcelParser

import (
	"encoding/json"
	"os"
)

func SaveAsJson(filePath *string, localization *Localization) {
	json, err := json.MarshalIndent(localization, "", "  ")
	check(err)
	os.WriteFile(*filePath, json, 0644)
}

func LoadfromJson(filePath *string) *Localization {
	jsonData, err := os.ReadFile(*filePath)
	check(err)
	var localization Localization
	error := json.Unmarshal(jsonData, &localization)
	check(error)
	return &localization
}
