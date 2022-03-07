package models

type Translation struct {
	Text string
}

func (t *Translation) GetTranslate() string {
	return t.Text
}