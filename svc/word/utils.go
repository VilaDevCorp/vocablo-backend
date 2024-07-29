package word

import (
	"vocablo/apischema"
	"vocablo/schema"
)

func ConvertApiResponseToWordForms(response apischema.ApiResponse) (wordFormsResult []CreateForm) {
	for _, wordMeaning := range response {
		wordFormsResult = append(wordFormsResult, ConvertApiWordToWordForm(wordMeaning))
	}
	return
}

func ConvertApiWordToWordForm(apiWord apischema.ApiResponseElement) CreateForm {
	var definitions []schema.Definition
	for _, meaning := range apiWord.Meanings {
		for _, definition := range meaning.Definitions {
			var example string = ""
			if definition.Example != nil {
				example = *definition.Example
			}
			definition := schema.Definition{PartOfSpeech: meaning.PartOfSpeech,
				Definition: definition.Definition, Example: example}
			definitions = append(definitions, definition)
		}
	}
	return CreateForm{
		Term:        apiWord.Word,
		Definitions: definitions,
		Lang:        "en",
	}
}
