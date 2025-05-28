package plscli

import (
	"net/http"
	"time"
)

type PlsClient struct {
	Http       *http.Client
	ClientId   string
	DeployName string
	Url        string
}

type ClientCfg struct {
	Host       string
	Port       string
	DeployName string
}

// Config 함수는 호스트, 포트, 디플로이 이름을 인자로 받아
// ClientCfg 객체를 리턴합니다.
func Config(host, port, deployName string) ClientCfg {
	return ClientCfg{
		Host:       host,
		Port:       port,
		DeployName: deployName,
	}
}

// NewClient 함수는 Config 함수로 생성한 설정을 인자로 받아
// 클라이언트 객체를 리턴합니다.
func NewClient(cfg ClientCfg) *PlsClient {
	return &PlsClient{
		Http:       &http.Client{Timeout: 5 * time.Second},
		DeployName: cfg.DeployName,
		Url:        cfg.Host + ":" + cfg.Port,
	}
}
