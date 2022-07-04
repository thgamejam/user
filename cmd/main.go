package main

import (
	"flag"
	"os"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"

	pkgConsul "github.com/thgamejam/pkg/consul"
	"user/internal/conf"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string = "user.service"
	// Version is the version of the compiled software.
	Version string
	// flagConfigPath is the config flag.
	flagConfigPath string
	// cloudConfigFile 服务注册，配置中心的配置文件
	cloudConfigFile string

	id, _ = os.Hostname()
)

func Init() {
	flag.StringVar(&flagConfigPath, "conf", "../configs", "配置路径，例如：-conf config.yaml")
	flag.StringVar(&cloudConfigFile, "cloud", "../configs/cloud.yaml", "配置&发现服务的配置路径，例如：-cloud cloud.yaml")
}

func newApp(logger log.Logger, rr registry.Registrar, hs *http.Server, gs *grpc.Server) *kratos.App {
	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			hs,
			gs,
		),
		kratos.Registrar(rr),
	)
}

func main() {
	Init()
	flag.Parse()
	logger := log.With(log.NewStdLogger(os.Stdout),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", id,
		"service.name", Name,
		"service.version", Version,
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID(),
	)

	cloudConfig := config.New(
		config.WithSource(
			file.NewSource(cloudConfigFile), // 获取本地的配置文件
		),
	)
	defer cloudConfig.Close()
	// 必须进行一次合并
	if err := cloudConfig.Load(); err != nil {
		panic(err)
	}

	var consulConfig conf.CloudBootstrap
	if err := cloudConfig.Scan(&consulConfig); err != nil {
		panic(err)
	}

	consulUtil := pkgConsul.New(consulConfig.Consul)

	serviceConfig := config.New(
		config.WithSource(
			file.NewSource(flagConfigPath), // 获取本地的配置文件
			consulUtil.NewConfigSource(),   // 获取配置中心的配置文件
		),
	)
	defer serviceConfig.Close()
	if err := serviceConfig.Load(); err != nil {
		panic(err)
	}

	// 读取配置到结构体
	var bc conf.Bootstrap
	if err := serviceConfig.Scan(&bc); err != nil {
		panic(err)
	}

	// 服务注册
	rr := consulUtil.NewRegistrar()
	// 服务发现
	rd := consulUtil.NewDiscovery()

	app, cleanup, err := wireApp(bc.Server, bc.Data, bc.User, rr, rd, logger)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}
