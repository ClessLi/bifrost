package resolv

import "encoding/json"

type Stream struct {
	BasicContext
}

func (st *Stream) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Stream []Parser `json:"stream"`
	}{st.Children})
}

func (st *Stream) UnmarshalJSON(b []byte) error {
	stream := struct {
		Stream []Parser `json:"stream"`
	}{}
	err := json.Unmarshal(b, &stream)
	if err != nil {
		return err
	}

	st.Name = "stream"
	st.Children = stream.Stream
	return nil
}

func NewStream() *Stream {
	return &Stream{BasicContext{
		Name:     "stream",
		Value:    "",
		depth:    0,
		Children: nil,
	}}
}
