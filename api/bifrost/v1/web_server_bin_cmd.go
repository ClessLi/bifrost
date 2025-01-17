package v1

type ExecuteRequest struct {
	ServerName string   `json:"server_name"`
	Args       []string `json:"args"`
}
type ExecuteResponse struct {
	Successful     bool   `json:"successful"`
	StandardOutput []byte `json:"stdout"`
	StandardError  []byte `json:"stderr"`
}
