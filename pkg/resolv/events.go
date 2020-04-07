package resolv

import "encoding/json"

type Events struct {
	BasicContext `json:"events"`
}

func (e *Events) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Events []Parser `json:"events"`
	}{Events: e.Children})
}

func (e *Events) UnmarshalJSON(b []byte) error {
	events := struct {
		Events []Parser `json:"events"`
	}{}
	err := json.Unmarshal(b, &events)
	if err != nil {
		return err
	}

	e.Name = "events"
	e.Children = events.Events
	return nil
}

func NewEvents() *Events {
	return &Events{BasicContext{
		Name:     "events",
		Value:    "",
		depth:    0,
		Children: nil,
	}}
}
