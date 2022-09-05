#!/bin/bash

yaml="dev-compose.yaml"

# 判断Docker的环境

# wls2 环境
docker="docker.exe"
docker_compose="docker-compose.exe"
docker_path=$(command -v docker.exe)
docker_compose_path=$(command -v docker-compose.exe)

# 正常 docker 环境
if [ -z "$docker_path" ]; then
  docker_path=$(command -v docker)
  docker="docker"
  docker_compose_path=$(command -v docker-compose)
  docker_compose="docker-compose"
fi

# podman 环境
if [ -z "$docker_path" ]; then
  docker_path=$(command -v podman)
  docker="podman"
  docker_compose_path=$(command -v podman-compose)
  docker_compose="podman-compose"
fi

echo "docker path: $docker_path"
echo "docker-compose path: $docker_compose_path"

# 创建网络
network_name="dev_bridge"
filter_name=$($docker network ls | grep $network_name | awk '{ print $2 }')
if [[ -z $filter_name ]]; then
  echo '创建测速环境网络...'
  $docker network create --subnet=172.88.0.0/24 --gateway=172.88.0.1 $network_name
fi

$docker_compose -f $yaml up -d
