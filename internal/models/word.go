package models

import "linnoxlewis/tg-bot-eng-dictionary/internal/helpers"

type Word struct {
	Id       int
	Text     string
	Meanings []Meanings
}

func NewWordIngot() []Word {
	return []Word{}
}

func (w *Word) ToTranslateMessage(lang string) string {
	result := ""
	if w.Id == 0 {
		return "translate for this word not found :( \n"
	}
	result = "translate for this word is :\n"
	if lang == helpers.RuLanguage {
		return result + w.Text
	} else {
		for _, value := range w.Meanings {
			result = result + value.Translation.Text + "\n"
		}

		return result
	}
}
