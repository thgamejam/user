package conf

import (
	"flag"
	"os"
)

type Env struct {
	// 服务ID
	ID string

	// ConsulURL consul的url, 例如: "http://0.0.0.0:8500"
	ConsulURL string
	// ConsulDatacenter 使用的consul数据中心, 例如: "dev"
	ConsulDatacenter string
	// ConsulConfDirectory 使用的consul配置路径, 例如:, "thjam.user.service/dev/"
	ConsulConfDirectory string

	// ConfigDirectory 配置文件路径
	ConfigDirectory string
}

var env Env = Env{}

func InitEnvDefaultValue() {
	// 添加默认值
	env.ID, _ = os.Hostname()
	env.ConsulURL = "http://0.0.0.0:8500"
	env.ConsulDatacenter = "dev"
	env.ConsulConfDirectory = "thjam.user.service/dev/"
	env.ConfigDirectory = "../configs"
}

func InitEnv() {
	InitEnvDefaultValue()

	flag.StringVar(&env.ConfigDirectory, "conf", "../configs", "配置路径，例如：-conf config.yaml")
	flag.Parse()

	// 环境变量覆盖输入参数
	_id := os.Getenv("SERVICE_ID")
	if _id != "" {
		env.ID = _id
	}
	_consulURL := os.Getenv("CONSUL_URL")
	if _consulURL != "" {
		env.ConsulURL = _consulURL
	}
	_consulDatacenter := os.Getenv("CONSUL_DATACENTER")
	if _consulDatacenter != "" {
		env.ConsulDatacenter = _consulDatacenter
	}
	_consulConfDir := os.Getenv("CONSUL_CONF_DIRECTORY")
	if _consulConfDir != "" {
		env.ConsulConfDirectory = _consulConfDir
	}

	_configDirectory := os.Getenv("CONFIG_DIRECTORY")
	if _configDirectory != "" {
		env.ConfigDirectory = _configDirectory
	}
}

// GetEnv 获取环境变量储存类型
func GetEnv() *Env {
	return &env
}
