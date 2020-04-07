package resolv

import "encoding/json"

type Upstream struct {
	BasicContext
}

func (u *Upstream) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Value    string   `json:"value"`
		Upstream []Parser `json:"upstream"`
	}{Value: u.Value, Upstream: u.Children})
}

func (u *Upstream) UnmarshalJSON(b []byte) error {
	upstream := struct {
		Value    string   `json:"value"`
		Upstream []Parser `json:"upstream"`
	}{}
	err := json.Unmarshal(b, &upstream)
	if err != nil {
		return err
	}

	u.Name = "upstream"
	u.Value = upstream.Value
	u.Children = upstream.Upstream
	return nil
}

func NewUpstream(value string) *Upstream {
	return &Upstream{BasicContext{
		Name:     "upstream",
		Value:    value,
		depth:    0,
		Children: nil,
	}}
}
