package main

import (
	"os"
	"strings"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"

	pkgConf "github.com/thgamejam/pkg/conf"
	pkgConsul "github.com/thgamejam/pkg/consul"
	"user/internal/conf"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string = "thjam.user.service"
	// Version is the version of the compiled software.
	Version string
)

func newApp(logger log.Logger, rr registry.Registrar, hs *http.Server, gs *grpc.Server) *kratos.App {
	return kratos.New(
		kratos.ID(conf.GetEnv().ID),
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
	conf.InitEnv()

	logger := log.With(log.NewStdLogger(os.Stdout),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", conf.GetEnv().ID,
		"service.name", Name,
		"service.version", Version,
		//"trace.id", tracing.TraceID(),
		//"span.id", tracing.SpanID(),
	)

	var pc pkgConf.Consul
	pc = pkgConf.Consul{
		Address:    strings.Split(conf.GetEnv().ConsulURL, "://")[1],
		Scheme:     strings.Split(conf.GetEnv().ConsulURL, "://")[0],
		Datacenter: conf.GetEnv().ConsulDatacenter,
		Path:       conf.GetEnv().ConsulConfDirectory,
	}

	consulUtil := pkgConsul.New(&pc)

	serviceConfig := config.New(
		config.WithSource(
			file.NewSource(conf.GetEnv().ConfigDirectory), // 获取本地的配置文件
			consulUtil.NewConfigSource(),                  // 获取配置中心的配置文件
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
