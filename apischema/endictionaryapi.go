package apischema

type ApiResponse []ApiResponseElement

type ApiResponseElement struct {
	Word       string            `json:"word"`
	Phonetic   string            `json:"phonetic"`
	Phonetics  []PhoneticElement `json:"phonetics"`
	Meanings   []Meaning         `json:"meanings"`
	SourceUrls []string          `json:"sourceUrls"`
}

type Meaning struct {
	PartOfSpeech string       `json:"partOfSpeech"`
	Definitions  []Definition `json:"definitions"`
	Synonyms     []string     `json:"synonyms"`
	Antonyms     []string     `json:"antonyms"`
}

type Definition struct {
	Definition string   `json:"definition"`
	Synonyms   []string `json:"synonyms"`
	Antonyms   []string `json:"antonyms"`
	Example    *string  `json:"example,omitempty"`
}

type PhoneticElement struct {
	Text      string `json:"text"`
	Audio     string `json:"audio"`
	SourceURL string `json:"sourceUrl"`
}
