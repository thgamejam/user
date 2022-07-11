package conf

import (
	"flag"
	"os"
)

type Env struct {
	// ConsulURL consul的url, 例如: "http://0.0.0.0:8500"
	ConsulURL string
	// ConsulDatacenter 使用的consul数据中心, 例如: "dev"
	ConsulDatacenter string
	// ConsulConfDirectory 使用的consul配置路径, 例如:, "thjam.upload.service/dev/"
	ConsulConfDirectory string

	// ConfigDirectory 配置文件路径
	ConfigDirectory string

	// TemporaryFileDirectory 临时文件存放目录
	TemporaryFileDirectory string
}

var env Env = Env{
	// 添加默认值
	ConsulURL:              "http://0.0.0.0:8500",
	ConsulDatacenter:       "dev",
	ConsulConfDirectory:    "thjam.upload.service/dev/",
	ConfigDirectory:        "../configs",
	TemporaryFileDirectory: "../tmp",
}

func InitEnv() {
	flag.StringVar(&env.ConfigDirectory, "conf", "../configs", "配置路径，例如：-conf config.yaml")
	flag.StringVar(&env.TemporaryFileDirectory, "tmp", "../tmp", "临时文件存放的路径, 例如：-tmp ./tmp/file_edge")
	flag.Parse()

	// 环境变量覆盖输入参数
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

	_temporaryFilePath := os.Getenv("TEMPORARY_FILE_DIRECTORY")
	if _temporaryFilePath != "" {
		env.TemporaryFileDirectory = _temporaryFilePath
	}
}

// GetEnv 获取环境变量储存类型
func GetEnv() *Env {
	return &env
}
