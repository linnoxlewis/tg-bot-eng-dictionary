package models

type Example struct {
	Text string
}

func (e *Example) getExample() string {
	return e.Text
}