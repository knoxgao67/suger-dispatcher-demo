#!/bin/bash

# 获取脚本文件所在目录
script_dir=$(dirname "$0")



# prepare
# 1. kafka
#helm repo add rhcharts https://ricardo-aires.github.io/helm-charts/
#helm upgrade --install kafka rhcharts/kafka

# 2. temporal
# cd /tmp
# git clone https://github.com/temporalio/helm-charts.git
# cd /tmp/helm-charts/charts/temporal
#  helm install \
#    --set server.replicaCount=1 \
#    --set cassandra.config.cluster_size=1 \
#    --set elasticsearch.replicas=1 \
#    --set prometheus.enabled=false \
#    --set grafana.enabled=false \
#    temporaltest . --timeout 15m


# 进入脚本文件所在目录
cd "$script_dir"
# make image
docker build  -t dispatcher-demo:latest .

echo "uninstall dispatcher-demo if exist"
helm uninstall dispatcher-demo

echo "install dispatcher-demo"
helm install dispatcher-demo ./helm_chart

