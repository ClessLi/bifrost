package resolv

import "encoding/json"

type Map struct {
	BasicContext
}

func (m *Map) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Value string `json:"value"`
		Map   []Parser
	}{Value: m.Value, Map: m.Children})
}

func (m *Map) UnmarshalJSON(b []byte) error {
	tMap := struct {
		Value string `json:"value"`
		Map   []Parser
	}{}
	err := json.Unmarshal(b, &tMap)
	if err != nil {
		return err
	}

	m.Name = "map"
	m.Value = tMap.Value
	m.Children = tMap.Map
	return nil
}

func NewMap(value string) *Map {
	return &Map{BasicContext{
		Name:     "map",
		Value:    value,
		depth:    0,
		Children: nil,
	}}
}
