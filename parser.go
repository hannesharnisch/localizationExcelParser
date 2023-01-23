package localizationExcelParser

import (
	"encoding/json"
	"os"
	"reflect"
	"strconv"

	"github.com/xuri/excelize/v2"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func ToExcel(path *string, output *string) {
	letters := []rune{}
	for ch := 'A'; ch <= 'Z'; ch++ {
		letters = append(letters, ch)
	}
	f := excelize.NewFile()
	model := loadModel(path)
	for _, scope := range model.Scopes {
		f.NewSheet(scope.Name)
		f.SetCellValue(scope.Name, "A1", "Beschreibung")
		f.SetCellValue(scope.Name, "B1", "Key")
		for i, lang := range model.Languages {
			f.SetCellValue(scope.Name, string(letters[2:][i])+"1", lang.Code)
		}
	}
	f.DeleteSheet("Sheet1")
	for _, scope := range model.Scopes {
		i := 2
		for _, value := range scope.Entries {
			check(f.SetCellValue(scope.Name, "B"+strconv.Itoa(i), value.Key))
			check(f.SetCellValue(scope.Name, "A"+strconv.Itoa(i), value.Description))
			a := 2
			for _, trans := range value.Translations {
				check(f.SetCellValue(scope.Name, string(letters[a])+strconv.Itoa(i), trans.Value))
				a += 1
			}
			i += 1
		}
	}
	error := f.SaveAs(*output)
	check(error)
	f.Close()
}

func FromExcel(path *string, output *string, defaultLang string) {
	f, error := excelize.OpenFile(*path)
	check(error)

	sheets := f.GetSheetList()

	languages := getLanguages(f, 2)
	localization := Localization{
		Languages: []Language{},
		Scopes:    []LocalizationScope{},
	}
	for _, x := range languages {
		localization.Languages = append(localization.Languages, Language{
			Code:      x,
			IsDefault: x == defaultLang,
		})
	}
	for _, sheet := range sheets {
		rows, error := f.GetRows(sheet)
		check(error)
		scope := LocalizationScope{
			Name: sheet,
		}
		for _, row := range rows[1:] {
			entry := LocalizationEntry{}
			for x, cell := range row {
				if x == 0 {
					entry.Description = cell
				} else if x == 1 {
					entry.Key = cell
				} else {
					entry.Translations = append(entry.Translations, Translation{
						Language: localization.Languages[x-2].Code,
						Value:    cell,
					})
				}
			}
			scope.Entries = append(scope.Entries, entry)
		}
		localization.Scopes = append(localization.Scopes, scope)
	}
	saveModel(output, &localization)
	f.Close()
}

func getLanguages(f *excelize.File, offset int) []string {
	languages := map[string]bool{}
	sheets := f.GetSheetList()
	for _, sheet := range sheets {
		rows, _ := f.GetRows(sheet)
		for _, lang := range rows[0][offset:] {
			languages[lang] = true
		}
	}
	keys := reflect.ValueOf(languages).MapKeys()
	strkeys := make([]string, len(keys))
	for i := 0; i < len(keys); i++ {
		strkeys[i] = keys[i].String()
	}
	return strkeys
}

func loadModel(modelPath *string) *Localization {
	jsonData, err := os.ReadFile(*modelPath)
	check(err)
	var localization Localization
	error := json.Unmarshal(jsonData, &localization)
	check(error)
	return &localization
}

func saveModel(modelPath *string, localization *Localization) {
	json, err := json.MarshalIndent(localization, "", "  ")
	check(err)
	os.WriteFile(*modelPath, json, 0644)
}
