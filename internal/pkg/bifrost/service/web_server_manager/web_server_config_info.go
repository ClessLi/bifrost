package web_server_manager

// WebServerType, web服务器类型对象，定义web服务器所属类型
type WebServerType string

const (
	// Web服务类型
	NGINX WebServerType = "nginx"
	HTTPD WebServerType = "httpd"
)

type WebServerConfigInfo struct {
	Name           string        `yaml:"name"`
	Type           WebServerType `yaml:"type"`
	BackupCycle    int           `yaml:"backupCycle"`
	BackupSaveTime int           `yaml:"backupSaveTime"`
	BackupDir      string        `yaml:"backupDir,omitempty"`
	ConfPath       string        `yaml:"confPath"`
	VerifyExecPath string        `yaml:"verifyExecPath"`
}
