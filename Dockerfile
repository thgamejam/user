FROM golang:1.17 AS builder

COPY . /src
WORKDIR /src

RUN GOPROXY=https://proxy.golang.com.cn,direct make build

FROM debian:stable-slim

RUN echo "deb http://mirrors.aliyun.com/debian stable main" > /etc/apt/sources.list && \
    echo "deb http://mirrors.aliyun.com/debian-security stable-security main" >> /etc/apt/sources.list && \
    echo "deb http://mirrors.aliyun.com/debian stable-updates main" >> /etc/apt/sources.list

RUN apt-get update && apt-get install -y --no-install-recommends \
		ca-certificates  \
        netbase \
        net-tools \
        && rm -rf /var/lib/apt/lists/ \
        && apt-get autoremove -y && apt-get autoclean -y

COPY --from=builder /src/bin /app

RUN mkdir -p /data/conf
COPY ./configs/* /data/conf

WORKDIR /app

EXPOSE 8000
EXPOSE 9000
VOLUME /data/conf

CMD ["./server", "-cloud", "/data/conf/cloud.yaml", "-conf", "/data/conf"]
