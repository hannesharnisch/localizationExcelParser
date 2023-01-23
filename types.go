package localizationExcelParser

type Localization struct {
	Scopes    []LocalizationScope `json:"scopes"`
	Languages []Language          `json:"langs"`
}

type LocalizationScope struct {
	Name    string              `json:"name"`
	Entries []LocalizationEntry `json:"entries"`
}

type LocalizationEntry struct {
	Key          string        `json:"key"`
	Description  string        `json:"description"`
	Translations []Translation `json:"translations"`
}

type Translation struct {
	Language string `json:"lang"`
	Value    string `json:"value"`
}

type Language struct {
	Code      string `json:"code"`
	IsDefault bool   `json:"isDefault"`
}
