package models

import (
	"fmt"
	"linnoxlewis/tg-bot-eng-dictionary/internal/helpers"
)

type Mean struct {
	Id          string
	WordId      int
	Text        string
	Translation Translation
	Definition  Definition
	Examples    []Example
}

func NewMeanIngot() []Mean {
	return []Mean{}
}

func (m *Mean) ToTranslateMessage(lang string) string {
	result := ""
	if m.Id == "" {
		return "translate for this word not found :( \n"
	}
	if lang == helpers.RuLanguage {
		return fmt.Sprintf("<b>Translate</b> : %s  \n", m.Text)
	}

	translate := fmt.Sprintf("<b>Translate</b> : %s  \n", m.Translation.GetTranslate())
	result = result + translate

	desc := fmt.Sprintf("<b>Descritption</b> : %s \n", m.Definition.GetDefinition())
	result = result + desc

	examples := fmt.Sprintf("<b>Examples</b> : \n")
	result = result + examples
	for _, value := range m.Examples {
		result = result + value.getExample() + "\n"
	}
	return result

}

func (m *Mean) GeneratingWordsToMessage() string {
	return fmt.Sprintf("%s - %s \n", m.Text, m.Translation.GetTranslate())
}
