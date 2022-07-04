# user

## 项目结构


```
.
├── bin     // 项目生成的可运行文件
├── cmd     // 项目启动的入口文件
│   ├── main.go     // 程序入口
│   ├── wire.go     // 使用wire来维护依赖注入
│   └── wire_gen.go // 使用wire生成的go文件
├── configs     // 可以本地使用的配置文件样例
├── dev         // 测试环境
│   ├── dev-compose.yaml    // 测试环境docker-compose文件
│   └── dev-environment.sh  // 测试环境启动脚本
├── interface   // 该服务所有不对外暴露的代码，通常的业务逻辑都在这下面
│   ├── biz         // 业务逻辑的组装层
│   ├── conf        // 内部使用config的结构定义以及根据结构定义所生成的go文件
│   ├── data        // 业务数据访问
│   ├── server      // http、grpc和mq实例的创建和配置
│   └── service     // 实现了api定义的服务层
├── proto   // 公用的proto文件，从proto子项目导入
│   ├── api         // 微服务使用的proto文件以及根据它们所生成的go文件
│   ├── third_party // api 依赖的第三方proto
│   └── conf        // 通用config结构定义的Proto文件
└── sql     // 数据库sql
```


## Generate other auxiliary files by Makefile
```bash
# 初始化项目并下载和更新依赖项
make init

# 运行
make run

# 构建
make build

# 依赖注入
make wire

# 生成错误文件代码
make error

# 生成配置文件代码
make config

# 生成api文件代码
make api

# 生成所有代码
make all

# 移除所有生成代码
make remove

# 显示帮助
make help
```


## Docker
```bash
# 构建
docker build -t user:<version> .

# 运行
docker run --rm -p 8000:8000 -p 9000:9000 -v </path/to/your/configs>:/data/conf user:<version>

# docker-compose 运行
cd dev
docker-compose -f service-compose.yaml up -d
```
