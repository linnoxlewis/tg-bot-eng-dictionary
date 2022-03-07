package models

type Definition struct {
	Text string
}

func (d *Definition) GetDefinition() string {
	return d.Text
}
