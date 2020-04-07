package resolv

import "encoding/json"

type Location struct {
	BasicContext
}

func (l *Location) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Value    string   `json:"value"`
		Location []Parser `json:"location"`
	}{Value: l.Value, Location: l.Children})
}

func (l *Location) UnmarshalJSON(b []byte) error {
	location := struct {
		Value    string   `json:"value"`
		Location []Parser `json:"location"`
	}{}
	err := json.Unmarshal(b, &location)
	if err != nil {
		return err
	}

	l.Name = "location"
	l.Value = location.Value
	l.Children = location.Location
	return nil
}

func NewLocation(value string) *Location {
	return &Location{BasicContext{
		Name:     "location",
		Value:    value,
		depth:    0,
		Children: nil,
	}}
}
