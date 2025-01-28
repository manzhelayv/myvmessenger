#!/bin/bash

mkdir ~/mnt/kubernetes/kafka/zookeeper
mkdir ~/mnt/kubernetes/kafka/zookeeper/data
mkdir ~/mnt/kubernetes/kafka/zookeeper/log
mkdir ~/mnt/kubernetes/kafka/kafka
mkdir ~/mnt/kubernetes/kafka/kafka/data
mkdir ~/mnt/kubernetes/kafka/kafka/log
mkdir ~/mnt/kubernetes/kafka/kafka/secrets
mkdir ~/mnt/kubernetes/kafka/connect
mkdir ~/mnt/kubernetes/kafka/akhq
mkdir ~/mnt/kubernetes/kafka/ksqldb

mkdir ~/mnt/kubernetes/db/postgres
mkdir ~/mnt/kubernetes/db/postgres/data
mkdir ~/mnt/kubernetes/db/mongodb
mkdir ~/mnt/kubernetes/db/mongodb/db
mkdir ~/mnt/kubernetes/db/mongodb/backup
mkdir ~/mnt/kubernetes/db/clickhouse
mkdir ~/mnt/kubernetes/db/clickhouse/clickhouse
mkdir ~/mnt/kubernetes/db/clickhouse/clickhouse-server

mkdir ~/mnt/kubernetes/minio/data

mkdir ~/mnt/kubernetes/backend/server
mkdir ~/mnt/kubernetes/backend/server/log
echo > ~/mnt/kubernetes/backend/chat/log/logs

mkdir ~/mnt/kubernetes/backend/f3
mkdir ~/mnt/kubernetes/backend/f3/log
echo > ~/mnt/kubernetes/backend/f3/log/logs

mkdir ~/mnt/kubernetes/backend/chat
mkdir ~/mnt/kubernetes/backend/chat/log
echo > ~/mnt/kubernetes/backend/chat/log/logs

export JOB_DATE=`date +%s`

kubectl apply -f ~/mnt/kubernetes/minio/configmap.yaml
kubectl apply -f ~/mnt/kubernetes/minio/volumes.yaml
kubectl apply -f ~/mnt/kubernetes/minio/minio.yaml

kubectl apply -f ~/mnt/kubernetes/db/volumes.yaml
kubectl apply -f ~/mnt/kubernetes/db/configmap.yaml
kubectl apply -f ~/mnt/kubernetes/db/db.yaml
kubectl apply -f ~/mnt/kubernetes/db/cronjob.yaml

kubectl apply -f ~/mnt/kubernetes/kafka/volumes.yaml
kubectl apply -f ~/mnt/kubernetes/kafka/configmaps.yaml
kubectl apply -f ~/mnt/kubernetes/kafka/kafka.yaml

kubectl apply -f ~/mnt/kubernetes/backend/configmaps.yaml
kubectl apply -f ~/mnt/kubernetes/backend/volumes.yaml
kubectl apply -f ~/mnt/kubernetes/backend/backend.yaml

docker-compose -f ./gitlab/docker-compose.yml up -d























