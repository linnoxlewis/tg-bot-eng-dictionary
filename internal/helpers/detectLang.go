package helpers

import "regexp"

var (
	ruRegexp   = "[А-Яа-я]+?"
	enRegexp   = "[A-Za-z]+?"
	RuLanguage = "ru"
	EnLanguage = "en"
)

func DetectLang(word string) string {
	match, _ := regexp.MatchString(ruRegexp, word)
	if match {
		return RuLanguage
	}

	return EnLanguage
}
