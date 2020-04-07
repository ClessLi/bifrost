package resolv

import "encoding/json"

type Types struct {
	BasicContext
}

func (t *Types) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Types []Parser `json:"types"`
	}{t.Children})
}

func (t *Types) UnmarshalJSON(b []byte) error {
	types := struct {
		Types []Parser `json:"types"`
	}{}
	err := json.Unmarshal(b, &types)
	if err != nil {
		return err
	}

	t.Name = "types"
	t.Children = types.Types
	return nil
}

func NewTypes() *Types {
	return &Types{BasicContext{
		Name:     "types",
		Value:    "",
		depth:    0,
		Children: nil,
	}}
}
