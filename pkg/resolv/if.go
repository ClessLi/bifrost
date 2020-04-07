package resolv

import "encoding/json"

type If struct {
	BasicContext `json:"if"`
}

func (i *If) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Value string   `json:"value"`
		If    []Parser `json:"if"`
	}{Value: i.Value, If: i.Children})
}

func (i *If) UnmarshalJSON(b []byte) error {
	tIf := struct {
		Value string   `json:"value"`
		If    []Parser `json:"if"`
	}{}
	err := json.Unmarshal(b, &tIf)
	if err != nil {
		return err
	}

	i.Name = "if"
	i.Value = tIf.Value
	i.Children = tIf.If
	return nil
}

func NewIf(value string) *If {
	return &If{BasicContext{
		Name:     "if",
		Value:    value,
		depth:    0,
		Children: nil,
	}}
}
